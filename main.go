package main

import (
	"bytes"
	"os"
	"strings"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

type TodoApp struct {
	app      *tview.Application
	todoFile string
	todoList *tview.List
	todos    []string
}

func (t *TodoApp) PopulateTodos() error {
	b, err := os.ReadFile("todo.txt")
	if err != nil {
		return err
	}

	todos := tview.NewList()
	byte_a := byte('a')
	for line := range bytes.SplitSeq(b, []byte("\n")) {
		if len(line) == 0 {
			continue
		}

		t.todos = append(t.todos, string(line))
		todos.AddItem(string(line), "", rune(byte_a), nil)
		byte_a++
	}
	t.todoList = todos
	return nil
}

func (t *TodoApp) Save() error {
	var sb strings.Builder
	for _, item := range t.todos {
		sb.WriteString(item + "\n")
	}
	return os.WriteFile(t.todoFile, []byte(sb.String()), 0o644)
}

func (t *TodoApp) Refresh() {
	t.todoList.Clear()
	for _, item := range t.todos {
		t.todoList.AddItem(item, "", 0, nil)
	}
}

func (t *TodoApp) AddTodo() {
	for _, item := range t.todos {
		t.todoList.AddItem(item, "", 0, nil)
	}
}

func main() {
	todoApp := TodoApp{todoFile: "todo.txt"}
	app := tview.NewApplication()
	todoApp.app = app

	err := todoApp.PopulateTodos()
	if err != nil {
		panic(err)
	}

	todoApp.todoList.SetBorder(true).SetTitle(" Todo List ")

	help := tview.NewTextView().
		SetText("q: Quit | a: Add | d: Delete").
		SetDynamicColors(true).
		SetTextAlign(tview.AlignCenter)
	help.SetBorder(true).SetTitle(" Help ")

	box := tview.NewFlex().
		SetDirection(tview.FlexRow).
		AddItem(todoApp.todoList, 0, 1, true).
		AddItem(help, 3, 1, false)

	input := tview.NewInputField().
		SetLabel("Add todo: ")
	modal := tview.NewFlex().
		SetDirection(tview.FlexRow).
		AddItem(nil, 0, 1, false).
		AddItem(tview.NewFlex().
			AddItem(nil, 0, 1, false).
			AddItem(input, 0, 1, true).
			AddItem(nil, 0, 1, false),
			10, 1, true).
		AddItem(nil, 0, 1, false)

	input.SetDoneFunc(func(key tcell.Key) {
		switch key {
		case tcell.KeyEnter:
			todoApp.todos = append(todoApp.todos, input.GetText())
			todoApp.AddTodo()
			todoApp.Save()
			todoApp.app.SetRoot(box, true).SetFocus(todoApp.todoList)
		case tcell.KeyEsc:
			todoApp.app.SetRoot(box, true).SetFocus(todoApp.todoList)
		}
	})

	todoApp.app.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Rune() {
		case 'a':
			todoApp.app.SetRoot(modal, true).SetFocus(input)
			return nil
		case 'q':
			todoApp.app.Stop()
		}
		return event
	})

	if err := app.SetRoot(box, true).Run(); err != nil {
		panic(err)
	}
}
