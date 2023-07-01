package xii

import (
	"errors"
	"fmt"
	"golang.org/x/sync/singleflight"
	"log"
	"sync"
	"testing"
	"time"
)

const (
	RegionAmerica = "America"
	RegionEurope  = "Europe"
	RegionAsia    = "Asia"
	RegionAfrica  = "Africa"
)

var mpRegionWaitTime = map[string]time.Duration{
	RegionAmerica: 5 * time.Second,
	RegionEurope:  3 * time.Second,
	RegionAfrica:  4 * time.Second,
	RegionAsia:    2 * time.Second,
}

func getToken(region string, business string) (string, error) {
	waitTime, ok := mpRegionWaitTime[region]
	if !ok || waitTime == 0 {
		return "", errors.New("unsupported region: " + region)
	}
	log.Printf("[getToken] region: %s, business: %s, wait-time: %v", region, business, waitTime)
	time.Sleep(waitTime)
	return fmt.Sprintf("%s|%s|%d", region, business, time.Now().UnixMilli()), nil
}

var getTokenGroup singleflight.Group

type GetTokenTask struct {
	Region   string
	Business string
	callback func() (interface{}, error)
}

func (t *GetTokenTask) key() string {
	return fmt.Sprintf("%s|%s", t.Region, t.Business)
}

func (t *GetTokenTask) Do() string {
	key := t.key()
	v, err, _ := getTokenGroup.Do(key, t.callback)
	if err != nil {
		log.Printf("[GetTokenTask] [%s] get token err: %v", key, err)
		return ""
	}
	token, ok := v.(string)
	if !ok {
		log.Printf("[GetTokenTask] [%s] convert token to string err", key)
		return ""
	}
	log.Printf("[GetTokenTask] [%s] got token: %s", key, token)
	return token
}

func newGetTokenTask(region string, business string) *GetTokenTask {
	return &GetTokenTask{
		Region:   region,
		Business: business,
		callback: func() (interface{}, error) {
			return getToken(region, business)
		},
	}
}

func GetToken(region string, business string) string {
	if region == "" || business == "" {
		return ""
	}
	task := newGetTokenTask(region, business)
	return task.Do()
}

func TestGetToken(t *testing.T) {
	numTasks := 1000

	// 1 business, random region
	business := "gofromzero"
	regions := []string{RegionAmerica, RegionEurope, RegionAsia, RegionAfrica}

	// run multiple tasks
	t.Logf("start %d tasks...", numTasks)
	var wg sync.WaitGroup
	for i := 0; i < numTasks; i++ {
		num := i + 1
		wg.Add(1)
		go func(n int) {
			defer wg.Done()
			idx := n % len(regions)
			region := regions[idx]
			token := GetToken(region, business)
			if token == "" {
				t.Logf("[task:%d] get token failed!", n)
			} else {
				t.Logf("[task:%d] got token -> %s", n, token)
			}
		}(num)
		time.Sleep(1 * time.Millisecond)
	}
	wg.Wait()

	t.Logf("finish!")
}
