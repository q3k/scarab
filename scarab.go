package scarab

import (
	"sync"

	"github.com/golang/glog"
	"github.com/golang/protobuf/proto"

	cpb "github.com/q3k/scarab/proto/common"
	spb "github.com/q3k/scarab/proto/storage"
)

type JobDefinition struct {
	Name        string
	Description string

	Arguments []*cpb.ArgumentDefinition

	Steps []StepDefinition
}

func (j *JobDefinition) Proto() *cpb.JobDefinition {
	res := &cpb.JobDefinition{
		Name:        j.Name,
		Description: j.Description,
		Arguments:   j.Arguments,
		Steps:       make([]*cpb.StepDefinition, len(j.Steps)),
	}
	for j, step := range j.Steps {
		res.Steps[j] = &cpb.StepDefinition{
			Name:         step.Name,
			Description:  step.Description,
			DependsSteps: step.DependsSteps,
		}
	}
	return res
}

func UnmarshalJobDefinition(j *cpb.JobDefinition) *JobDefinition {
	res := &JobDefinition{
		Name:        j.Name,
		Description: j.Description,
		Arguments:   j.Arguments,
		Steps:       make([]StepDefinition, len(j.Steps)),
	}
	for i, step := range j.Steps {
		res.Steps[i] = StepDefinition{
			Name:         step.Name,
			Description:  step.Description,
			DependsSteps: step.DependsSteps,
		}
	}
	return res
}

type StepDefinition struct {
	Name         string
	Description  string
	DependsSteps []string
}

type RunningJob struct {
	id int64

	definition *JobDefinition

	Arguments []*cpb.Argument
	State     proto.Message
}

func (r *RunningJob) Proto() *spb.RunningJob {
	return &spb.RunningJob{
		Id:         r.id,
		Definition: r.definition.Proto(),
		Arguments:  r.Arguments,
	}
}

func UnmarshalRunningJob(j *spb.RunningJob) *RunningJob {
	return &RunningJob{
		id:         j.Id,
		definition: UnmarshalJobDefinition(j.Definition),
		Arguments:  j.Arguments,
	}
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
			Arguments:  j.Arguments,
			State:      j.State,
		}
	}

	return res
}

func (s *Service) Save() error {
	return nil
}
