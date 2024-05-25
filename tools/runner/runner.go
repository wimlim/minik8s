package runner

import (
	"time"
)

type Runner struct {
}

func NewRunner() *Runner {
	return &Runner{}
}

func (r *Runner) RunLoop(delay time.Duration, period time.Duration, fun func()) {

	<-time.After(delay)

	for {
		fun()
		<-time.After(period)
	}
}
