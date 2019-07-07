package scarab

import (
	"fmt"
	"sync"

	"github.com/golang/protobuf/proto"
	"github.com/syndtr/goleveldb/leveldb"
	leveldb_util "github.com/syndtr/goleveldb/leveldb/util"

	spb "github.com/q3k/scarab/proto/storage"
)

type key interface {
	Get() []byte
}

type keyRunningJob struct {
	id int64
}

func (k keyRunningJob) Get() []byte {
	s := "running-job/"
	if k.id != 0 {
		s += fmt.Sprintf("%016x", k.id)
	}

	return []byte(s)
}

func (r *RunningJob) Key() key {
	return keyRunningJob{
		id: r.id,
	}
}

type levelDBStorage struct {
	db *leveldb.DB

	// Used for operations like Create which require atomic changes
	mu sync.Mutex
}

func NewLevelDBStorage(path string) (Storage, error) {
	db, err := leveldb.OpenFile(path, nil)
	if err != nil {
		return nil, fmt.Errorf("Could not create or open database: %v", err)
	}

	return &levelDBStorage{
		db: db,
	}, nil
}
func (l *levelDBStorage) UpdateBatch(jobs []*RunningJob) error {
	marshaled := make([][]byte, len(jobs))

	// Marshal all jobs first
	for i, j := range jobs {
		if j.id == 0 {
			return fmt.Errorf("cannot batch store new job %v", j)
		}

		p := j.Proto()
		b, err := proto.Marshal(p)
		if err != nil {
			return fmt.Errorf("could not marshal job %v: %v", j, err)
		}
		marshaled[i] = b
	}

	// Update them in the database
	for i, j := range jobs {
		b := marshaled[i]
		err := l.db.Put(j.Key().Get(), b, nil)
		if err != nil {
			return fmt.Errorf("could not store job %v: %v", j, err)
		}
	}

	return nil
}

func (l *levelDBStorage) Create(job *RunningJob) error {
	// Fallback to batch update when job already has ID
	if job.id != 0 {
		return l.UpdateBatch([]*RunningJob{job})
	}

	l.mu.Lock()
	defer l.mu.Unlock()

	// Range over jobs to find highest ID.
	iter := l.db.NewIterator(leveldb_util.BytesPrefix(job.Key().Get()), nil)

	var id int64 = 1
	if iter.Last() {
		b := iter.Value()
		job := spb.RunningJob{}
		err := proto.Unmarshal(b, &job)
		if err != nil {
			return fmt.Errorf("could not retrieve last job: %v", err)
		}
		id = job.Id + 1
	}

	// Save job.
	job.id = id
	return l.UpdateBatch([]*RunningJob{job})
}
