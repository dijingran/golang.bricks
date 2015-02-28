package sort

import (
	"testing"
	"strings"
	"time"
	"fmt"
)

const (
	TIMES = 1000 * 100
)

func TestQuickSort(t *testing.T) {
	a:=[]int{1,3,4,5,2,8,7}
	QSort(a, 0, 6)
	fmt.Printf("%v \n",a)
}

func TestBubbleSort(t *testing.T) {
	s1 := "12346789"
	e := "98764321"
	s2 := strings.Repeat("12346789a", 10)
	if e != BubbleSort(s1) {
		t.Fatalf("Wrong result.")
	}

	start := time.Now().UnixNano()
	for i := 0; i < TIMES; i++ {
		BubbleSort(s2)
	}
	trace("BubbleSort", start)

	if e != SelectSort(s1) {
		t.Fatalf("Wrong result.")
	}
	start = time.Now().UnixNano()
	for i := 0; i < TIMES; i++ {
		SelectSort(s2)
	}
	trace("SelectSort", start)

	if e != BubbleSort2(s1) {
		t.Fatalf("Wrong result.")
	}
	start = time.Now().UnixNano()
	for i := 0; i < TIMES; i++ {
		BubbleSort2(s2)
	}
	trace("BubbleSort2", start)
}

func trace(name string, start int64) {
	fmt.Printf("%s cost %d ms.\n", name, (time.Now().UnixNano()-start)/1000/1000)
}
