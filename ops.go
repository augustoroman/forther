package main

import (
	"fmt"
	"math"
	"strings"
)

type PushToStack struct{}

func (a PushToStack) Match(string) bool { return true }
func (a PushToStack) Run(f *Forther, val string) error {
	f.stack.Push(val)
	return nil
}

type SimpleMathOps struct{}

func (a SimpleMathOps) Match(s string) bool {
	return s == "+" || s == "-" || s == "*" || s == "/" || s == "mod"
}
func (a SimpleMathOps) Run(f *Forther, op string) error {
	nums, err := f.stack.PopNumbers(2)
	if err != nil {
		return err
	}
	switch op {
	case "+":
		f.stack.Push(fmt.Sprint(nums[0] + nums[1]))
	case "-":
		f.stack.Push(fmt.Sprint(nums[0] - nums[1]))
	case "*":
		f.stack.Push(fmt.Sprint(nums[0] * nums[1]))
	case "/":
		f.stack.Push(fmt.Sprint(nums[0] / nums[1]))
	case "mod":
		f.stack.Push(fmt.Sprint(math.Mod(nums[0], nums[1])))
	}
	return nil
}

type Keyword struct {
	name string
	fn   func(*Forther, string) error
}

func (k Keyword) Match(cmd string) bool           { return cmd == k.name }
func (k Keyword) Run(f *Forther, op string) error { return k.fn(f, op) }
func (k Keyword) Complete(line string) []string {
	if strings.HasPrefix(k.name, line) {
		return []string{k.name}
	}
	return nil
}
