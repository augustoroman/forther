package main

import (
	"fmt"
	"io"
	"log"
	"strings"

	"github.com/peterh/liner"
)

func main() {
	if err := commandLoop(); err != nil {
		log.Fatalln(err)
	}
}

func commandLoop() error {
	in := liner.NewLiner()
	defer in.Close()

	in.SetCtrlCAborts(true)

	f := NewForther()
	in.SetCompleter(f.Complete)
	for {
		line, err := in.Prompt(f.Prompt())
		if err == io.EOF || err == liner.ErrPromptAborted {
			return nil
		} else if err != nil {
			return err
		}
		in.AppendHistory(line)

		for i, cmd := range strings.Fields(line) {
			if err := f.Process(cmd); err == io.EOF {
				return nil
			} else if err != nil {
				fmt.Printf("Cannot run op %d (%q): %v\n", i+1, cmd, err)
				break
			}
		}
	}
}

type Forther struct {
	stack     stack
	showStack bool
	ops       []Operation
}

type Operation interface {
	Match(op string) bool
	Run(f *Forther, op string) error
}

// Some operations may also support command-line completion.  If so, they
// implement the Completer interface.
type Completer interface {
	Complete(line string) []string
}

func NewForther() Forther {
	return Forther{
		ops: []Operation{
			Keyword{"bye", func(*Forther, string) error { return io.EOF }},
			Keyword{"help", func(f *Forther, _ string) error { f.PrintHelp(); return nil }},

			// Prompt display:
			Keyword{"showstack",
				func(f *Forther, _ string) error { f.showStack = true; return nil }},
			Keyword{"noshowstack",
				func(f *Forther, _ string) error { f.showStack = false; return nil }},

			// Show top or all of stack.
			Keyword{".", func(f *Forther, _ string) error { f.PrintTop(); return nil }},
			Keyword{".s", func(f *Forther, _ string) error { f.PrintStack(); return nil }},

			// Stack operations
			Keyword{"dup", func(f *Forther, _ string) error {
				a := f.stack.Pop()
				f.stack.Push(a, a)
				return nil
			}},
			Keyword{"drop", func(f *Forther, _ string) error { f.stack.Pop(); return nil }},
			Keyword{"swap", func(f *Forther, _ string) error {
				b, a := f.stack.Pop(), f.stack.Pop()
				f.stack.Push(b, a)
				return nil
			}},
			Keyword{"over", func(f *Forther, _ string) error {
				b, a := f.stack.Pop(), f.stack.Pop()
				f.stack.Push(a, b, a)
				return nil
			}},

			SimpleMathOps{}, // Handle most math operations
			PushToStack{},   // Anything we don't recognize, just push onto the stack.
		},
	}
}

func (f *Forther) Prompt() string {
	const prompt = "» "
	if f.showStack {
		return strings.Join(f.stack, " ") + " " + prompt
	}
	return prompt
}

func (f *Forther) Process(cmd string) error {
	for _, op := range f.ops {
		if op.Match(cmd) {
			return op.Run(f, cmd)
		}
	}
	panic("Unknown operation")
}

func (f *Forther) Complete(line string) []string {
	var options []string
	for _, op := range f.ops {
		if completer, ok := op.(Completer); ok {
			options = append(options, completer.Complete(line)...)
		}
	}
	return options
}

func (f *Forther) PrintStack() {
	for _, s := range f.stack {
		fmt.Println(s)
	}
}
func (f *Forther) PrintTop() {
	fmt.Println(f.stack.Peek())
}
func (f *Forther) PrintHelp() {
	var cmds []string
	for _, op := range f.ops {
		if kw, ok := op.(Keyword); ok {
			cmds = append(cmds, kw.name)
		}
	}
	fmt.Println("Known commands: ")
	fmt.Println(" ", strings.Join(cmds, " "))
	fmt.Println("  + - * / mod ")
	fmt.Println("(anything is pushed onto the stack)")
}
