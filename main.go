package main

import (
	"github.com/gofromzero/v"
	"time"
)

func main() {
	v.TestServerAsync()
	go func() {
		time.Sleep(5 * time.Second)
		v.TestClient()
	}()
	lb, err := v.NewLoadBalancerOnConfig("./v/config.json")
	if err != nil {
		panic(err)
	}
	lb.Run()
}
