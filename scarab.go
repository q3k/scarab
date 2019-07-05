package scarab

import (
	"sync"

	"github.com/golang/glog"
	"github.com/golang/protobuf/proto"

	cpb "github.com/q3k/scarab/proto/common"
)

type JobDefinition struct {
	Name        string
	Description string

	ArgsDescriptor []*cpb.ArgumentDefinition

	Steps []StepDefinition
}

type StepDefinition struct {
	Name        string
	Description string
	DependsStep string
}

type RunningJob struct {
	definition *JobDefinition

	Args  proto.Message
	State proto.Message
}

type RunningStep struct {
	definition *StepDefinition
	job        *RunningJob
}

type Service struct {
	Definitions map[string]*JobDefinition
	jobsMu      sync.RWMutex
	jobs        []*RunningJob
	storage     Storage
}

func NewService(definitions []*JobDefinition, storage Storage) *Service {
	// Create scarab structures.
	// TODO(q3k): Restore state.

	s := &Service{
		Definitions: make(map[string]*JobDefinition),
		jobs:        []*RunningJob{},
		storage:     storage,
	}

	for _, def := range definitions {
		glog.Infof("Loaded Job Definition %q", def.Name)
		s.Definitions[def.Name] = def
	}

	return s
}

func (s *Service) RunningJobs() []*RunningJob {
	s.jobsMu.Lock()
	defer s.jobsMu.Unlock()

	res := make([]*RunningJob, len(s.jobs))
	for i, j := range s.jobs {
		res[i] = &RunningJob{
			definition: j.definition,
			Args:       j.Args,
			State:      j.State,
		}
	}

	return res
}

func (s *Service) Save() error {
	return nil
}
