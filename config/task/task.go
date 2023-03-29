package task

type (
	Crontab struct {
		Name     string
		Explain  string
		Rule     string
		Callback func()
	}
	Task struct {
		Crontabs []Crontab
	}
)

var task = &Task{}

func Set(config *Task) {
	task = config
}

func Get() *Task {
	return task
}
