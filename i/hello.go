package i

import (
	"fmt"
)

// Hell hell
func Hell() {
	fmt.Println("Hell")
	panic("haha")
}

// Hello hello
func Hello(s string) {
	fmt.Printf("Hell%s World", s)
}
