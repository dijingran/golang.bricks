package sort

import "fmt"

func SelectSort(line string) (s string) {
	b := []byte(line)
	n := len(b)
	for i := 0; i < n-1 ; i++ {
		for j := i + 1 ; j < n  ; j++ {
			if b[j] > b[i] {
				t := b[j]
				b[j] = b[i]
				b[i] = t
			}
		}
	}
	return string(b)
}

func BubbleSort(line string) (s string) {
	b := []byte(line)
	n := len(b)
	for i := 0; i < n-1 ; i++ {
		for j := 0 ; j < n-i-1  ; j++ {
			if b[j] < b[j+1] {
				t := b[j]
				b[j] = b[j+1]
				b[j+1] = t
			}
		}
	}
	return string(b)
}


func BubbleSort2(line string) (s string) {
	b := []byte(line)
	for i := len(line) - 1; i > 0 ; {
		pos := 0
		for j := 0 ; j < i; j++ {
			if b[j] < b[j+1] {
				t := b[j]
				b[j] = b[j+1]
				b[j+1] = t
				pos = j
			}
		}
		i = pos
	}
	return string(b)
}


func QuickSort(line string) (s string) {
	return line
	//	return string(b)
}

// TODO
func QSort(b []int, l, r int) {
	i := l;
	j := r
	x := b[i]
	for ; i < j ; {
		for ; j > i; {
			if b[j] > x {
				b[i] = b[j]
				i++
				break
			}
			j--
		}
		fmt.Println(j)

		for ; i < j; {
			if b[i] < x {
				b[j] = b[i]
				j--
				break
			}
			i++
		}
		fmt.Println(i)
	}
}
