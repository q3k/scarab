package scarab

import (
	"github.com/golang/protobuf/proto"

	spb "github.com/q3k/scarab/proto/state"
)

type JobDefinition struct {
	Name        string
	Description string

	ArgsDescriptor []*spb.ArgumentDefinition

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
	Jobs        []*RunningJob
}
