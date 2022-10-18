package GoScheduler

import (
	"context"
	"sync"
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

func TestCancel(t *testing.T) {
	waitTime = time.Millisecond * 50000
	sch := CreateDynamicScheduler()
	ctx, cf := context.WithCancel(context.Background())
	sch.ScheduleWithContext(time.Millisecond, true, func() {
		assert.Fail(t, "Bad")
	}, ctx)
	cf()
	waiter := sync.WaitGroup{}
	waiter.Add(1)
	ctx, cf = context.WithCancel(context.Background())
	sch.ScheduleWithContext(time.Millisecond*5, true, func() {
		waiter.Done()
	}, ctx)
	waiter.Wait()
	cf()
	time.Sleep(time.Millisecond * 10)
}
