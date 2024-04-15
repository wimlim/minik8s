package scheduler

import (
	"testing"
)

func TestTmp(t *testing.T) {
	str := Scheduler()
	if str != "Scheduler" {
		t.Errorf("Scheduler() = %v, want %v", str, "Scheduler")
	}
}
