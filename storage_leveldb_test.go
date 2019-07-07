package scarab

import (
	"fmt"
	"io/ioutil"
	"os"
	"testing"

	"github.com/go-test/deep"
	"github.com/golang/protobuf/proto"
	cpb "github.com/q3k/scarab/proto/common"
)

func makeTestDB() (Storage, func()) {
	dir, err := ioutil.TempDir("", "scarab-storage-test")
	if err != nil {
		panic(err)
	}

	l, err := NewLevelDBStorage(dir)
	if err != nil {
		panic(err)
	}
	cleanup := func() {
		os.RemoveAll(dir)
	}
	return l, cleanup
}

func compareRunningJob(a, b *RunningJob, t *testing.T) {
	// First compare all proto fields manually

	if want, got := len(a.Arguments), len(b.Arguments); want != got {
		t.Fatalf("job.Arguments size diff, want %d got %d", want, got)
	}
	for i, aa := range a.Arguments {
		ba := b.Arguments[i]
		if !proto.Equal(aa, ba) {
			t.Errorf("job.Arguments[%d] diff", i)
		}
	}

	if want, got := len(a.definition.Arguments), len(b.definition.Arguments); want != got {
		t.Fatalf("job.definition.Arguments size diff, want %d got %d", want, got)
	}
	for i, aa := range a.definition.Arguments {
		ba := b.definition.Arguments[i]
		if !proto.Equal(aa, ba) {
			t.Errorf("job.definition.Arguments[%d] diff", i)
		}
	}

	// Then, unset them and compare the rest.
	aa, ba := a.Arguments, b.Arguments
	ada, bda := a.definition.Arguments, b.definition.Arguments
	a.Arguments = nil
	b.Arguments = nil
	a.definition.Arguments = nil
	b.definition.Arguments = nil

	deep.CompareUnexportedFields = true
	if diff := deep.Equal(a, b); diff != nil {
		t.Error(diff)
	}

	a.Arguments = aa
	b.Arguments = ba
	a.definition.Arguments = ada
	b.definition.Arguments = bda
}

func jobDefinition() *JobDefinition {
	return &JobDefinition{
		Name:        "test-job",
		Description: "test job",
		Arguments: []*cpb.ArgumentDefinition{
			&cpb.ArgumentDefinition{
				Name:        "foo",
				Description: "Foo",
				Type:        cpb.ArgumentDefinition_TYPE_ONE_LINE_STRING,
			},
		},
		Steps: []StepDefinition{
			{Name: "bar", Description: "Bar"},
			{Name: "baz", Description: "Baz", DependsSteps: []string{"bar"}},
			{Name: "barfoo", Description: "Barfoo", DependsSteps: []string{"bar"}},
			{Name: "done", Description: "Done", DependsSteps: []string{"bar", "barfoo"}},
		},
	}
}

func TestCreate(t *testing.T) {
	db, cleanup := makeTestDB()
	defer cleanup()

	j := &RunningJob{
		definition: jobDefinition(),
		Arguments: []*cpb.Argument{
			&cpb.Argument{Name: "foo", Value: "foo!"},
		},
	}

	err := db.Create(j)
	if err != nil {
		t.Fatalf("db.CreateJob(%v): %v", err)
	}

	if j.id == 0 {
		t.Fatalf("created job did not get an ID")
	}

	r, err := db.Load()
	if err != nil {
		t.Fatalf("db.Load(): %v", err)
	}

	if want, got := 1, len(r); want != got {
		t.Fatalf("Got %d jobs, want %d", got, want)
	}

	lj := r[0]
	compareRunningJob(j, lj, t)
}

func TestCreateMultiple(t *testing.T) {
	db, cleanup := makeTestDB()
	defer cleanup()

	// Create jobs

	njobs := 10000
	createdJobs := make(map[int64]*RunningJob)

	for i := 0; i < njobs; i += 1 {
		argFoo := fmt.Sprintf("foo! %d", i)
		j := &RunningJob{
			definition: jobDefinition(),
			Arguments: []*cpb.Argument{
				&cpb.Argument{Name: "foo", Value: argFoo},
			},
		}
		err := db.Create(j)
		if err != nil {
			t.Fatalf("db.CreateJob(%v): %v", err)
		}

		if j.id == 0 {
			t.Fatalf("db.CreateJob(%v): did not get an id")
		}

		createdJobs[j.id] = j
	}

	if want, got := njobs, len(createdJobs); want != got {
		t.Fatalf("wanted to create %d jobs, created  %d", want, got)
	}

	// Check created jobs

	r, err := db.Load()
	if err != nil {
		t.Fatalf("db.Load(): %v", err)
	}

	if want, got := njobs, len(r); want != got {
		t.Fatalf("Got %d jobs, want %d", got, want)
	}

	seenJobs := make(map[int64]bool)
	for _, j := range r {
		compareRunningJob(createdJobs[j.id], j, t)
		seenJobs[j.id] = true
	}

	if want, got := njobs, len(seenJobs); want != got {
		t.Fatalf("Got %d unique jobs, want %d", got, want)
	}
}

func TestUpdateMultiple(t *testing.T) {
	db, cleanup := makeTestDB()
	defer cleanup()

	// Create jobs

	njobs := 10000
	createdJobs := make(map[int64]*RunningJob)

	for i := 0; i < njobs; i += 1 {
		argFoo := fmt.Sprintf("foo! %d", i)
		j := &RunningJob{
			definition: jobDefinition(),
			Arguments: []*cpb.Argument{
				&cpb.Argument{Name: "foo", Value: argFoo},
			},
		}
		err := db.Create(j)
		if err != nil {
			t.Fatalf("db.CreateJob(%v): %v", err)
		}

		createdJobs[j.id] = j
	}

	if want, got := njobs, len(createdJobs); want != got {
		t.Fatalf("wanted to create %d jobs, created  %d", want, got)
	}

	// Update all job arguments

	batch := make([]*RunningJob, len(createdJobs))
	i := 0
	for _, job := range createdJobs {
		argFoo := fmt.Sprintf("bar! %d", i)
		job.Arguments[0].Value = argFoo
		batch[i] = job
		i += 1
	}

	err := db.UpdateBatch(batch)
	if err != nil {
		t.Fatalf("db.Update(): %v", err)
	}

	// Check updated jobs

	r, err := db.Load()
	if err != nil {
		t.Fatalf("db.Load(): %v", err)
	}

	if want, got := njobs, len(r); want != got {
		t.Fatalf("Got %d jobs, want %d", got, want)
	}

	seenJobs := make(map[int64]bool)
	for _, j := range r {
		compareRunningJob(createdJobs[j.id], j, t)
		seenJobs[j.id] = true
	}

	if want, got := njobs, len(seenJobs); want != got {
		t.Fatalf("Got %d unique jobs, want %d", got, want)
	}
}
