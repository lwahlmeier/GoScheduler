package GoScheduler

import (
	"context"
	"fmt"
	"math/rand"
	"runtime"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestGetDefault(t *testing.T) {
	ds := GetDefaultScheduler()
	assert.NotNil(t, ds)
	assert.True(t, ds.IsRunning())
	assert.Equal(t, ds, GetDefaultScheduler())
	ds.(*DynamicScheduler).Stop()
	assert.False(t, ds.IsRunning())
	assert.True(t, GetDefaultScheduler().IsRunning())
}

func TestGetDefaultMany(t *testing.T) {
	defaultScheduler = atomic.Pointer[DynamicScheduler]{}
	dsChan := make(chan Scheduler, 1000)
	wait := sync.WaitGroup{}
	for i := 0; i < 100; i++ {
		wait.Add(1)
		go func() {
			dsChan <- GetDefaultScheduler()
			runtime.Gosched()
			wait.Done()
		}()
	}
	wait.Wait()
	ds := GetDefaultScheduler()
	for i := 0; i < 100; i++ {
		assert.Equal(t, ds, <-dsChan)
	}
}

func isAprox(value, compare time.Duration) bool {
	return (value - compare).Abs() <= time.Millisecond*10
}

func BasicSchedule(sch Scheduler, t *testing.T) {

	wg1 := sync.WaitGroup{}
	wg1.Add(1)
	fmt.Printf("Start\n")
	st := time.Now()
	sch.Schedule(time.Millisecond*5, false, func() {
		wg1.Done()
	})
	fmt.Printf("Scheduled\n")
	wg1.Wait()
	ts := time.Since(st)
	fmt.Printf("time:%s\n", ts)
	assert.True(t, ts >= time.Millisecond*5, "Shoud wait 5ms or more")
	assert.True(t, ts <= time.Millisecond*7, "Shoud not wait to long")
	fmt.Printf("Done\n")
}

func CancelSchedule(sch Scheduler, t *testing.T) {
	RAN := false
	lock := sync.Mutex{}
	ctx, cf := context.WithCancel(context.Background())
	sch.ScheduleWithContext(time.Millisecond*5, false, func() {
		lock.Lock()
		defer lock.Unlock()
		RAN = true
		fmt.Println("TEST")
	}, ctx)
	cf()
	time.Sleep(time.Millisecond * 10)
	lock.Lock()
	defer lock.Unlock()
	assert.False(t, RAN, "Should not be set")
}

func BasicSchedule2(sch Scheduler, t *testing.T) {
	var wg sync.WaitGroup
	slock := sync.Mutex{}
	addMap := make(map[int]time.Duration)
	addList := make([]int, 0)

	for i := 15; i < 100; i++ {
		wg.Add(1)
		st := time.Now()
		ci := i
		sch.ScheduleWithContext(time.Millisecond*time.Duration(i), false, func() {
			slock.Lock()
			defer wg.Done()
			defer slock.Unlock()
			ss := time.Since(st)
			addList = append(addList, ci)
			addMap[ci] = ss
		}, context.Background())
	}

	st := time.Now()
	wg.Wait()
	assert.True(t, isAprox(time.Since(st), time.Millisecond*100))

	// fmt.Println(time.Since(st))
	for _, v := range addList {
		k := v
		d := addMap[k]
		fmt.Printf("%d-%s-%t\n", k, d, isAprox(d, time.Millisecond*time.Duration(k)))
		assert.True(t, isAprox(d, time.Millisecond*time.Duration(k)))
	}

}

func RandomSchedule(sch Scheduler, t *testing.T) {
	var wg sync.WaitGroup
	slock := sync.Mutex{}
	addMap := make(map[int]time.Duration)
	addList := make([]int, 0)

	for i := 5; i < 100; i++ {
		wg.Add(1)
		st := time.Now()
		ci := rand.Intn(95) + 5
		sch.ScheduleWithContext(time.Millisecond*time.Duration(ci), false, func() {
			slock.Lock()
			defer wg.Done()
			defer slock.Unlock()
			ss := time.Since(st)
			addList = append(addList, ci)
			addMap[ci] = ss
		}, context.Background())
	}

	st := time.Now()
	wg.Wait()
	assert.True(t, isAprox(time.Since(st), time.Millisecond*100))

	// fmt.Println(time.Since(st))
	for _, v := range addList {
		k := v
		d := addMap[k]
		// fmt.Printf("%d-%s-%t\n", k, d, isAprox(d, time.Millisecond*time.Duration(k)))
		assert.True(t, isAprox(d, time.Millisecond*time.Duration(k)))
	}
}
