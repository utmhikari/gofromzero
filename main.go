package main

import "github.com/gofromzero/iiii"

func main() {
	err := iiii.StartMongo()
	// err := iiii.RollBack()
	if err != nil {
		panic(err)
	}
}
