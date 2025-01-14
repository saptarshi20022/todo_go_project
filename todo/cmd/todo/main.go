package main

import (
	"bufio"
	todo "cmd-todo-app/cmd"
	"errors"
	"flag"
	"fmt"
	"io"
	"os"
	"strings"
)

const (
	todoFile = ".todos.json"
)

func main() {
	add := flag.Bool("add", false, "add a new todo") // declaration of command flag add
	complete := flag.Int("complete", 0, "mark a todo file as complete") // 
	delete := flag.Int("delete", 0, "delete a todo file")
	list := flag.Bool("list", false, "list all todos")

	flag.Parse()

	todos := &todo.Todos{}

	if err := todos.Load(todoFile); err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(1)
	}

	switch{
	case *add:
		task, err := getInput(os.Stdin, flag.Args()...)
		if err != nil {
			fmt.Fprintln(os.Stderr, err.Error())
			os.Exit(1)
		}

		todos.Add(task)
		saveAndCheck(todoFile, todos)
		
	case *complete > 0:
		err := todos.Complete(*complete)
		if err != nil {
			fmt.Fprintln(os.Stderr, err.Error())
			os.Exit(1)
		}
		saveAndCheck(todoFile, todos)

	case *delete > 0:
		err := todos.Delete(*delete)
		if err != nil {
			fmt.Fprintln(os.Stderr, err.Error())
			os.Exit(1)
		}
		saveAndCheck(todoFile, todos)

	case *list:
		todos.List()


	default:
		fmt.Fprintln(os.Stdout, "invalid command")
		os.Exit(1)
	}

}


func saveAndCheck(todoFile string, todos *todo.Todos) {
	err := todos.Store(todoFile)
		if err != nil {
			fmt.Fprintln(os.Stderr, err.Error())
			os.Exit(1)
		}
}

func getInput(r io.Reader, args ...string)(string, error) {
	
	if len(args)>0{
		return strings.Join(args, " "), nil
	}

	scanner := bufio.NewScanner(r)
	scanner.Scan()

	if err := scanner.Err(); err != nil {
		return "", err
	}

	text := scanner.Text()

	if len(text) == 0 {
		return "", errors.New("empty todo is not allowed")
	}
	
	return text, nil
	
}