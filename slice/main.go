package main

import "fmt"

func main() {
	var arr1 [3]int
	arr2 := [3]int{1, 2, 3}

	// The assignment of array is the copy of deep value, not the pointer
	arr1 = arr2
	arr2[1] = 10

	// The assignment of slice is the copy of pointer, not the value
	arr3 := arr1[0:2]
	arr3[0] = 100

	fmt.Println(arr1, arr1[0:2], arr2[0:2])

	// according to the difference, we can judge array and slice
	arr4 := [...]int{1, 2, 3}
	arr5 := arr4
	arr4[1] = 10
	fmt.Println(arr4, arr5)

	// append will cover the value of the slice's tail
	arr6 := append(arr5[0:2], 666)
	arr6[0] = 20
	fmt.Println(arr5, arr6)
}
