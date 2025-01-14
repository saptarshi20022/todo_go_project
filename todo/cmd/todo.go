package todo

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"time"
)

type item struct {
	Task      string
	Done      bool
	CreatedAt time.Time
	CompletedAt time.Time
}

type Todos [] item // a slice of type item

func (t *Todos)  Add(task string){
	todo := item{
		Task: task,
		Done: false,
		CreatedAt: time.Now(),
		CompletedAt: time.Time{},
	}

	*t = append(*t, todo) // appending todo to the list
}

func (t *Todos) Complete(index int) error{
	list := *t
	if index <= 0 || index > len(list) {
		return errors.New("index invalid")
	}

	list[index-1].CompletedAt = time.Now()
	list[index-1].Done = true

	return nil
}

func (t *Todos) Delete(index int) error{
	list := *t
	if index <= 0 || index > len(list) {
		return errors.New("index invalid")
	}

	*t = append(list[:index-1], list[index:]...) // appends everything together except the index item

	return nil
}

func (t *Todos) List() {
	
	for i, item := range *t {
		i++
		fmt.Printf("%d - %s \n", i, item.Task)
	}

}

// Loading it into a file, from the file system
func (t *Todos) Load(filename string) error{
	file, err := os.ReadFile(filename)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return nil
		}
		return err
	}

	if len(file) == 0 {
		return err
	}

	err = json.Unmarshal(file, t) // stores the result of file into t
	if err != nil {
		return err
	}

	return nil
}

// to store the todo value and write it 

func (t *Todos) Store(fileaname string) error{
	
	data, err := json.Marshal(t) // adds the value in t to JSON encoding
	if err != nil {
		return err
	}

	return os.WriteFile(fileaname, data, 0644) // 0644 is the permission
}

