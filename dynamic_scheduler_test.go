package GoScheduler

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestBasicSchedule(t *testing.T) {
	waitTime = time.Millisecond * 50
	sch := CreateDynamicScheduler()
	BasicSchedule(sch, t)
	assert.True(t, sch.IsRunning())
	sch.Stop()
	assert.False(t, sch.IsRunning())
}

func TestCancelSchedule(t *testing.T) {
	waitTime = time.Millisecond * 50
	sch := CreateDynamicScheduler()
	CancelSchedule(sch, t)
	assert.True(t, sch.IsRunning())
	sch.Stop()
	assert.False(t, sch.IsRunning())
}

func TestBasicSchedule2(t *testing.T) {
	waitTime = time.Millisecond * 50
	sch := CreateDynamicScheduler()
	BasicSchedule2(sch, t)
	assert.True(t, sch.IsRunning())
	sch.Stop()
	assert.False(t, sch.IsRunning())
}

func TestRandomSchedule(t *testing.T) {
	waitTime = time.Millisecond * 50
	sch := CreateDynamicScheduler()
	RandomSchedule(sch, t)
	assert.True(t, sch.IsRunning())
	sch.WaitForStop()
	assert.False(t, sch.IsRunning())
}

func TestWaitStop(t *testing.T) {
	waitTime = time.Millisecond * 50
	sch := CreateDynamicScheduler()
	sch.Schedule(time.Millisecond, false, func() {
		time.Sleep(time.Millisecond * 100)
	})
	st := time.Now()
	time.Sleep(time.Millisecond * 5)

	sch.WaitForStop()
	assert.False(t, sch.IsRunning())
	assert.True(t, time.Since(st) >= (time.Millisecond*100))
}
