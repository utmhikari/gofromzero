package xi

import "testing"

func TestMake(t *testing.T) {
	var m = make([]int, 5, 10)
	t.Logf("[Make] m is %v, len: %v, cap: %v", m, len(m), cap(m))
}

func TestNoMake(t *testing.T) {
	var m [5]int
	t.Logf("[NoMake] m is %v, len: %v, cap: %v", m, len(m), cap(m))
}

func TestNew(t *testing.T) {
	var i = new(int)
	t.Logf("[New] value of i is %v", *i)
}

func TestNoNew(t *testing.T) {
	var i *int
	t.Logf("[NoNew] value of i is %v", *i)
}
