package x

import (
	"context"
	"errors"
	"testing"
	"time"
)

func TestCtxBase(t *testing.T) {
	t.Skip()
	ctxBg := context.Background()
	ctxTodo := context.TODO()
	t.Logf("context.Background: %s, %d", ctxBg, ctxBg)
	t.Logf("context.TODO: %s, %d", ctxBg, ctxTodo)
	select {
	case <-ctxTodo.Done():
		t.Logf("context.TODO is done")
	case <-time.After(1 * time.Second):
		t.Logf("timeout")
	}
}

func TestCtxWithCancel(t *testing.T) {
	t.Skip()
	cancelCause := errors.New("debug")
	ctxCancel, cancel := context.WithCancelCause(context.Background())
	t.Logf("context.WithCancel: %v, %p -> cause: %v", ctxCancel, cancel, cancelCause)

	sleepTimeout := 1 * time.Second
	waiterTimeout := 2 * time.Second

	// waiter for cancel
	join := make(chan string)
	go func(ctx context.Context, timeout time.Duration, retChan chan string) {
		var ret string
		select {
		case <-ctx.Done():
			t.Logf("ctx done! -> err: %v, cause: %v", ctx.Err(), context.Cause(ctx))
			ret = "done"
		case <-time.After(waiterTimeout):
			t.Logf("waiter timeout")
			ret = "timeout"
		}
		retChan <- ret
	}(ctxCancel, waiterTimeout, join)

	// sleep & do cancel with cause
	time.Sleep(sleepTimeout)
	cancel(cancelCause)
	t.Logf("cancel done!")

	// join waiter
	ret := <-join
	t.Logf("waiter ret: %s", ret)
}

func TestCtxWithDeadline(t *testing.T) {
	t.Skip()
	timeout := 3 * time.Second

	deadline := time.Now().Add(timeout)
	ctxDeadline, cancel := context.WithDeadline(context.Background(), deadline)
	t.Logf("context.WithDeadline: %v, %p", ctxDeadline, cancel)

	// deadline/cancel detector
	join := make(chan string)
	go func(ctx context.Context, retChan chan string) {
		var ret string
		select {
		case <-ctx.Done():
			ddl, ok := ctx.Deadline()
			if !ok {
				t.Logf("ctx deadline not set")
				ret = "nothing"
			} else if time.Now().After(ddl) {
				t.Logf("ctx reached deadline: %v -> err: %v", ddl, ctx.Err())
				ret = "deadline"
			} else {
				t.Logf("ctx early canceled! -> err: %v", ctx.Err())
				ret = "cancel"
			}
		}
		retChan <- ret
	}(ctxDeadline, join)

	// manually cancel after cancelTimeout
	cancelTimeout := 1 * time.Second
	time.Sleep(cancelTimeout)
	cancel()
	t.Logf("cancel done!")

	ret := <-join
	t.Logf("ret: %s", ret)
}

func TestCtxWithValue(t *testing.T) {
	//t.Skip()
	key1, value1 := "hello", "world"
	ctxValue1 := context.WithValue(context.Background(), key1, value1)
	key2, value2 := "foo", "bar"
	ctxValue2 := context.WithValue(ctxValue1, key2, value2)

	t.Logf("ctxValue1: %s", ctxValue1)
	t.Logf("ctxValue1.%s = %v", key1, ctxValue1.Value(key1))
	t.Logf("ctxValue1.%s = %v", key2, ctxValue1.Value(key2))
	t.Logf("ctxValue2: %s", ctxValue2)
	t.Logf("ctxValue2.%s = %v", key1, ctxValue2.Value(key1))
	t.Logf("ctxValue2.%s = %v", key2, ctxValue2.Value(key2))
}
