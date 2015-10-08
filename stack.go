package main

import (
	"fmt"
	"strconv"
)

type stack []string

func (s *stack) Push(v ...string) {
	*s = append(*s, v...)
}
func (s *stack) Peek() string {
	N := len(*s)
	if N == 0 {
		return ""
	}
	return (*s)[N-1]
}
func (s *stack) Pop() string {
	N := len(*s)
	if N == 0 {
		return ""
	}
	newS, val := (*s)[:N-1], (*s)[N-1]
	*s = newS
	return val
}
func (s *stack) PopNumbers(n int) ([]float64, error) {
	N := len(*s)
	if N < n {
		return nil, fmt.Errorf("Not enough entries on stack (need %d, have %d)", n, N)
	}
	nums := make([]float64, 0, n)
	for i := 1; i <= n; i++ {
		str := (*s)[N-i]
		num, err := strconv.ParseFloat(str, 64)
		if err != nil {
			return nil, fmt.Errorf("Arg %d (%q) can't be parsed as a number: %v", i, str, err)
		}
		nums = append(nums, num)
	}
	*s = (*s)[:N-n]
	return nums, nil
}
