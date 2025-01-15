package todo

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/alexeyco/simpletable"
	
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
	
	table := simpletable.New()

	table.Header = &simpletable.Header{
		Cells: []*simpletable.Cell{
			{Align: simpletable.AlignCenter, Text: "#"},
			{Align: simpletable.AlignCenter, Text: "Task"},
			{Align: simpletable.AlignCenter, Text: "Done?"},
			{Align: simpletable.AlignRight, Text: "Created At"},
			{Align: simpletable.AlignRight, Text: "Completed At"},
		},
	}

	var cells [][]*simpletable.Cell

	for index,item := range *t {
		index++ // we want the index to start from 1

		task := blue(item.Task)
		done := red("NO")
		completeTime := red("--")
		if item.Done {
			task = green(fmt.Sprintf("\u2705 %s", item.Task)) //adds tick mark
			done = green("YES")
			completeTime = green(item.CompletedAt.Format(time.RFC822))
		}

		cells = append(cells, *&[]*simpletable.Cell{ // appending the pointer to slice of simpletable cell
			{Text: fmt.Sprintf("%d", index)}, //Sprintf formats according to a format specifier and returns the resulting string.
			{Text: task},
			{Text: done},
			{Text: item.CreatedAt.Format(time.RFC822)},
			{Text: completeTime},
		})
	}

	table.Body = &simpletable.Body{Cells: cells}

	table.Footer = &simpletable.Footer{Cells: []*simpletable.Cell{
		{Align: simpletable.AlignCenter, Span: 5, Text: red(fmt.Sprintf("We have %d pending to-dos",t.countPending()))},
	}}

	table.SetStyle(simpletable.StyleUnicode)

	table.Println()

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

func (t *Todos) countPending() int{
	total := 0
	for _, item := range *t {
		if !item.Done {
			total++
		}
	}

	return total
}

