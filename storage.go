package scarab

type Storage interface {
	UpdateBatch(jobs []*RunningJob) error
	Create(job *RunningJob) error
}
