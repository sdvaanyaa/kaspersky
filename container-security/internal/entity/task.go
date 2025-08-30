package entity

const (
	TaskStateQueued  = "queued"
	TaskStateRunning = "running"
	TaskStateDone    = "done"
	TaskStateFailed  = "failed"
)

type Task struct {
	ID         string `json:"id"`
	Payload    string `json:"payload"`
	MaxRetries int    `json:"max_retries"`
	Retries    int
	State      string
}
