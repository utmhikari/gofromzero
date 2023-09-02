package xv

import (
	"encoding/json"
	"sync"
	"testing"
)

type User struct {
	Name string `json:"name"`
	Age  int    `json:"age"`
}

var user *User
var userOnce sync.Once

func initUser() {
	user = &User{}
	cfgStr := `{"name":"foobar","age":18}`
	if err := json.Unmarshal([]byte(cfgStr), user); err != nil {
		panic("load user err: " + err.Error())
	}
}

func getUser() *User {
	userOnce.Do(initUser)
	return user
}

func TestSyncOnce(t *testing.T) {
	var wg sync.WaitGroup
	for i := 1; i < 1000; i++ {
		wg.Add(1)
		go func(n int) {
			defer wg.Done()
			curUser := getUser()
			t.Logf("[%d] got user: %+v", n, curUser)
		}(i)
	}
	wg.Wait()
}
