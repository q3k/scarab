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

func (r *RunningJob) ProtoStorage() *spb.RunningJob {
	return &spb.RunningJob{
		Id:         r.id,
		Definition: r.definition.Proto(),
		Arguments:  r.Arguments,
	}
}

func (r *RunningJob) Proto() *cpb.RunningJob {
	return &cpb.RunningJob{
		Id:         r.id,
		Definition: r.definition.Proto(),
		Arguments:  r.Arguments,
	}
}

func UnmarshalRunningJobStorage(j *spb.RunningJob) *RunningJob {
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

func NewService(definitions []*JobDefinition, storage Storage) (*Service, error) {
	// Create scarab structures.

	r, err := storage.Load()
	if err != nil {
		return nil, err
	}

	s := &Service{
		Definitions: make(map[string]*JobDefinition),
		jobs:        r,
		storage:     storage,
	}

	for _, def := range definitions {
		glog.Infof("Loaded Job Definition %q", def.Name)
		s.Definitions[def.Name] = def
	}

	return s, nil
}

func (s *Service) Save() error {
	return nil
}
