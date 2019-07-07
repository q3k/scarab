package scarab

type Storage interface {
	// Update the given jobs. Given jobs must already exist in storage.
	UpdateBatch(jobs []*RunningJob) error
	// Create a job in storage, and set its new unique ID.
	Create(job *RunningJob) error
	// Load jobs from storage.
	Load() ([]*RunningJob, error)
}
