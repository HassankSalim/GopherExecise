package main

import (
	"fmt"
	"math/rand"
	"sync"
)

const (
	size = 1000000
)

func mergeSort(l, r int, a []int, w *sync.WaitGroup) {
	if l >= r {
		w.Done()
		return 
	}
	m := (l + r) / 2
	var waitgroup sync.WaitGroup
	waitgroup.Add(1)
	go mergeSort(l, m, a, &waitgroup)
	waitgroup.Add(1)
	go mergeSort(m + 1, r, a, &waitgroup)
	waitgroup.Wait()
	merge(l, m, r, a[:])
	w.Done()
	return 
}


func merge(s, m, e int, a []int) {
	i, j := s, m + 1
	k := 0
	temp := make([]int, e - s + 1)
	for ;i < m + 1 && j < e + 1; k++ {
		if a[i] < a[j] {
			temp[k] = a[i]
			i++
		} else {
			temp[k] = a[j]
			j++			
		}
	}
	for ;i < m + 1; i, k = i+1, k+1 {
		temp[k] = a[i]
	}
	for ;j < e + 1; j, k = j+1, k+1 {
		temp[k] = a[j]
	}
	
	for j, i = 0, s; i <= e; i, j = i+1, j+1 {
		a[i] = temp[j]
	}
}

func main() {
	arr := make([]int, size)
	for i := range arr {
		arr[i] = rand.Intn(10000)
	}
	var w sync.WaitGroup
	w.Add(1)
	go mergeSort(0, len(arr) - 1, arr[0:], &w)
	w.Wait()
	fmt.Printf("Sorted arr %v \n", arr)
}