package viii

import (
	"reflect"
	"testing"
	"time"
)

func TestCache(t *testing.T) {
	k, v := "hello", "world"
	var curCache Cache = Current()

	// set & get & delete
	curCache.Set(k, v)
	cached, ok := curCache.Get(k)
	if !ok {
		t.Fatalf("cannot cache %s:%s", k, v)
	} else {
		t.Logf("got cached %s:%v (type: %s)", k, cached, reflect.TypeOf(cached).Name())
	}
	curCache.Delete(k)
	_, ok = curCache.Get(k)
	if ok {
		t.Fatalf("cannot delete %s:%s", k, v)
	} else {
		t.Logf("delete cached %s:%s", k, v)
	}

	// set expire
	curCache.SetExpire(k, v, 1*time.Second)
	cached, ok = curCache.Get(k)
	if !ok {
		t.Fatalf("cannot cache %s:%s", k, v)
	} else {
		t.Logf("got cached %s:%v (type: %s)", k, cached, reflect.TypeOf(cached).Name())
	}
	time.Sleep(3 * time.Second)
	_, ok = curCache.Get(k)
	if ok {
		t.Fatalf("cannot expire %s:%s", k, v)
	} else {
		t.Logf("expired %s:%s", k, v)
	}
}
