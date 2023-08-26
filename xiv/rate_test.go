package xiv

import (
	"context"
	"errors"
	"go.uber.org/ratelimit"
	"golang.org/x/time/rate"
	"sync"
	"testing"
	"time"
)

func generateData(num int) []any {
	var data []any
	for i := 0; i < num; i++ {
		data = append(data, i)
	}
	return data
}

func process(obj any) (any, error) {
	integer, ok := obj.(int)
	if !ok {
		return nil, errors.New("invalid integer")
	}
	time.Sleep(1)
	nextInteger := integer * 10
	if integer%99 == 0 {
		return nextInteger, errors.New("not a happy number")
	}
	return nextInteger, nil
}

func TestRate(t *testing.T) {
	limit := rate.Limit(50)
	burst := 25
	limiter := rate.NewLimiter(limit, burst)
	size := 500

	data := generateData(size)
	var wg sync.WaitGroup
	startTime := time.Now()
	for i, item := range data {
		wg.Add(1)
		go func(idx int, obj any) {
			defer wg.Done()
			if err := limiter.Wait(context.Background()); err != nil {
				t.Logf("[%d] [EXCEPTION] wait err: %v", idx, err)
			}
			processed, err := process(obj)
			if err != nil {
				t.Logf("[%d] [ERROR] processed: %v, err: %v", idx, processed, err)
			} else {
				t.Logf("[%d] [OK] processed: %v", idx, processed)
			}
		}(i, item)
	}
	wg.Wait()
	endTime := time.Now()
	t.Logf("start: %v, end: %v, seconds: %v", startTime, endTime, endTime.Sub(startTime).Seconds())
}

func TestRateLimit(t *testing.T) {
	limiter := ratelimit.New(50)
	size := 500

	data := generateData(size)
	var wg sync.WaitGroup
	startTime := time.Now()
	for i, item := range data {
		wg.Add(1)
		go func(idx int, obj any) {
			defer wg.Done()
			limiter.Take()
			processed, err := process(obj)
			if err != nil {
				t.Logf("[%d] [ERROR] processed: %v, err: %v", idx, processed, err)
			} else {
				t.Logf("[%d] [OK] processed: %v", idx, processed)
			}
		}(i, item)
	}
	wg.Wait()
	endTime := time.Now()
	t.Logf("start: %v, end: %v, seconds: %v", startTime, endTime, endTime.Sub(startTime).Seconds())
}
