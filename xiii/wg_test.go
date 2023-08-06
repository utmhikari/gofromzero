package xiii

import (
	"sync"
	"testing"
	"time"
)

func TestWaitGroup(t *testing.T) {
	var wg sync.WaitGroup

	startTime := time.Now()
	for i := 0; i < 5; i++ {
		n := i + 1
		sleepTime := time.Duration(n) * time.Second
		wg.Add(1)

		go func() {
			defer wg.Done()

			t.Logf("task %d started", n)
			time.Sleep(sleepTime)
			t.Logf("task %d ended", n)
		}()
	}
	t.Logf("waiting for all tasks done...")
	wg.Wait()
	endTime := time.Now()
	t.Logf("all tasks done! elapsed time: %v", endTime.Sub(startTime))
}
