package scarab

import (
	"fmt"

	"github.com/syndtr/goleveldb/leveldb"
)

type Storage interface {
	Store(jobs []*RunningJob) error
}

type levelDBStorage struct {
	db *leveldb.DB
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

func (l *levelDBStorage) Store(jobs []*RunningJob) error {
	return nil
}
