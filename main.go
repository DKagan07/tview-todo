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

	t.todoList = tview.NewList()
	for line := range bytes.SplitSeq(b, []byte("\n")) {
		if len(line) == 0 {
			continue
		}

		t.todos = append(t.todos, string(line))
		t.AddTodo(string(line))
	}
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
		t.AddTodo(item)
	}
}

func (t *TodoApp) AddTodo(item string) {
	t.todoList.AddItem(item, "", '-', nil)
}

func (t *TodoApp) DeleteTodo(idx int) {
	t.todoList = t.todoList.RemoveItem(idx)
	t.todos = append(t.todos[:idx], t.todos[idx+1:]...)
	t.Save()
	t.Refresh()
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

	container := tview.NewFlex().
		SetDirection(tview.FlexRow).
		AddItem(todoApp.todoList, 0, 1, true).
		AddItem(help, 3, 1, false)

	input := tview.NewInputField().
		SetLabel("Add todo: ").
		SetFieldBackgroundColor(tcell.ColorBlack)
	inputFrame := tview.NewFrame(input).
		SetBorders(1, 1, 2, 2, 1, 1)

	modal := tview.NewFlex().
		SetDirection(tview.FlexRow).
		AddItem(nil, 0, 1, false).
		AddItem(tview.NewFlex().
			AddItem(nil, 0, 1, false).
			AddItem(inputFrame, 0, 1, true).
			AddItem(nil, 0, 1, false),
			10, 1, true).
		AddItem(nil, 0, 1, false)

	input.SetDoneFunc(func(key tcell.Key) {
		switch key {
		case tcell.KeyEnter:
			todoApp.todos = append(todoApp.todos, input.GetText())
			todoApp.AddTodo(input.GetText())
			todoApp.Save()
			todoApp.Refresh()
			todoApp.app.SetRoot(container, true).SetFocus(todoApp.todoList)
		case tcell.KeyEsc:
			todoApp.app.SetRoot(container, true).SetFocus(todoApp.todoList)
		}
	})

	todoApp.todoList.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Rune() {
		case 'a':
			todoApp.app.SetRoot(modal, true).SetFocus(input)
			return nil
		case 'd':
			currentItemIdx := todoApp.todoList.GetCurrentItem()
			todoApp.DeleteTodo(currentItemIdx)
		}
		return event
	})

	todoApp.app.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Rune() {
		case 'q':
			todoApp.app.Stop()
		}

		if event.Key() == tcell.KeyCtrlC {
			todoApp.app.Stop()
		}
		return event
	})

	if err := app.SetRoot(container, true).Run(); err != nil {
		panic(err)
	}
}
