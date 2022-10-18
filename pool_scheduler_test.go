package GoScheduler

import (
	"runtime"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestBasicPoolSchedule(t *testing.T) {
	waitTime = time.Millisecond * 50
	sch := CreatePoolScheduler(10)
	BasicSchedule(sch, t)
	assert.True(t, sch.IsRunning())
	sch.Stop()
	assert.False(t, sch.IsRunning())
}

func TestCancelPoolSchedule(t *testing.T) {
	waitTime = time.Millisecond * 50
	sch := CreatePoolScheduler(10)
	CancelSchedule(sch, t)
	assert.True(t, sch.IsRunning())
	sch.Stop()
	assert.False(t, sch.IsRunning())
}

func TestBasicPoolSchedule2(t *testing.T) {
	waitTime = time.Millisecond * 50
	sch := CreatePoolScheduler(10)
	BasicSchedule2(sch, t)
	assert.True(t, sch.IsRunning())
	sch.Stop()
	assert.False(t, sch.IsRunning())
}

func TestRandomPoolSchedule(t *testing.T) {
	waitTime = time.Millisecond * 50
	sch := CreatePoolScheduler(10)
	RandomSchedule(sch, t)
	assert.True(t, sch.IsRunning())
	sch.Stop()
	assert.False(t, sch.IsRunning())
}

func TestClear(t *testing.T) {
	waitTime = time.Millisecond * 50
	sch := CreatePoolScheduler(10)
	for i := 0; i < 100; i++ {
		rt := time.Minute
		sch.Schedule(rt, false, func() {})
	}
	runtime.Gosched()
	time.Sleep(time.Millisecond)
	l := sync.Mutex{}
	l.Lock()
	assert.Equal(t, 100, len(sch.jobs))
	l.Unlock()
	sch.Clear()
	runtime.Gosched()
	time.Sleep(time.Millisecond)
	l.Lock()
	assert.Equal(t, 0, len(sch.jobs))
	l.Unlock()
	assert.True(t, sch.IsRunning())
	sch.Stop()
	assert.False(t, sch.IsRunning())
}

func TestDelay(t *testing.T) {
	sch := CreatePoolScheduler(10)
	waiter := sync.WaitGroup{}
	waiter.Add(10)
	count := 0
	sch.Schedule(time.Millisecond*10, true, func() {
		count += 1
		if count <= 10 {
			waiter.Done()
		}
	})
	waiter.Wait()
	assert.True(t, sch.IsRunning())
	sch.Stop()
	assert.False(t, sch.IsRunning())
}

func TestWaitStopPool(t *testing.T) {
	waitTime = time.Millisecond * 50
	sch := CreatePoolScheduler(10)
	sch.Schedule(time.Millisecond, false, func() {
		time.Sleep(time.Millisecond * 100)
	})
	st := time.Now()
	time.Sleep(time.Millisecond * 5)

	sch.WaitForStop()
	assert.False(t, sch.IsRunning())
	assert.True(t, time.Since(st) >= (time.Millisecond*100))
}

func TestNegativePool(t *testing.T) {
	waitTime = time.Millisecond * 50
	sch := CreatePoolScheduler(10)
	sch.Schedule(-time.Millisecond, false, func() {
		time.Sleep(time.Millisecond * 10)
	})
	sch.Schedule(-time.Second, false, func() {
		time.Sleep(time.Millisecond * 10)
	})
	time.Sleep(time.Millisecond * 5)
	sch.WaitForStop()
	assert.False(t, sch.IsRunning())
}
