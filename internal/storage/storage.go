package storage

type Storage interface {
	AddTask(t *Task) error
	RemoveTask(t *Task) error
	ListTasks(userId string) ([]string, error)
	IfExistsTask(t *Task) (bool, error)
}

type Task struct {
	TaskID string
	UserID string
	Name   string
}
