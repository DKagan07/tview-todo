package model

import (
	"bytes"
	"fmt"
	"os"
	"strings"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

type TodoApp struct {
	App      *tview.Application
	TodoFile string
	TodoList *tview.List
	Todos    []string
}

func (t *TodoApp) PopulateTodos() error {
	b, err := os.ReadFile("todo.txt")
	if err != nil {
		return err
	}

	t.TodoList = tview.NewList()
	for line := range bytes.SplitSeq(b, []byte("\n")) {
		if len(line) == 0 {
			continue
		}

		t.Todos = append(t.Todos, string(line))
		t.AddTodo(string(line))
	}
	return nil
}

func (t *TodoApp) Save() error {
	var sb strings.Builder
	for _, item := range t.Todos {
		sb.WriteString(item + "\n")
	}
	return os.WriteFile(t.TodoFile, []byte(sb.String()), 0o644)
}

func (t *TodoApp) Refresh() {
	t.TodoList.Clear()
	for _, item := range t.Todos {
		t.AddTodo(item)
	}
}

func (t *TodoApp) AddTodo(item string) {
	t.TodoList.AddItem(item, "", '-', nil)
}

func (t *TodoApp) DeleteTodo(idx int) {
	modal := tview.NewModal().
		SetText("Are you sure you want to delete this todo?").
		AddButtons([]string{"Yes", "No"})

	modal.SetBorder(true)
	modal.SetBorderStyle(tcell.StyleDefault.Background(tcell.ColorBlack))
	modal.SetBackgroundColor(tcell.ColorBlack)
	modal.SetBorderColor(tcell.ColorWhite)
	modal.SetTextColor(tcell.ColorWhite)

	modal.SetDoneFunc(func(_ int, buttonLabel string) {
		if buttonLabel == "Yes" {
			t.TodoList = t.TodoList.RemoveItem(idx)
			t.Todos = append(t.Todos[:idx], t.Todos[idx+1:]...)
			if err := t.Save(); err != nil {
				panic(err)
			}
			t.Refresh()
		}
		t.App.SetRoot(t.TodoList, true).SetFocus(t.TodoList)
	})

	frame := tview.NewFrame(modal).SetBorders(1, 1, 1, 1, 1, 1).Clear()
	t.App.SetRoot(modal, true).SetFocus(frame)
}

func (t *TodoApp) SetKeybindings(container *tview.Flex) {
	// Creating components for the add to todo modal

	// Add the keybinds
	t.TodoList.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Rune() {
		case 'a':
			input := tview.NewInputField().
				SetLabel("Add todo: ").
				SetFieldBackgroundColor(tcell.ColorBlack)
			inputFrame := tview.NewFrame(input).
				SetBorders(1, 1, 2, 2, 1, 1)
			input.SetDoneFunc(func(key tcell.Key) {
				switch key {
				case tcell.KeyEnter:
					t.Todos = append(t.Todos, input.GetText())
					t.AddTodo(input.GetText())
					if err := t.Save(); err != nil {
						panic(err)
					}
					t.Refresh()
					t.App.SetRoot(container, true).SetFocus(t.TodoList)
				case tcell.KeyEsc:
					t.App.SetRoot(container, true).SetFocus(t.TodoList)
				}
			})
			modal := tview.NewFlex().
				SetDirection(tview.FlexRow).
				AddItem(nil, 0, 1, false).
				AddItem(tview.NewFlex().
					AddItem(nil, 0, 1, false).
					AddItem(inputFrame, 0, 1, true).
					AddItem(nil, 0, 1, false),
					10, 1, true).
				AddItem(nil, 0, 1, false)

			t.App.SetRoot(modal, true).SetFocus(input)
			return nil
		case 'd':
			currentItemIdx := t.TodoList.GetCurrentItem()
			t.DeleteTodo(currentItemIdx)
		case 'u':
			currentItemIdx := t.TodoList.GetCurrentItem()
			currentItem := t.Todos[currentItemIdx]

			text := tview.NewTextView().
				SetText(fmt.Sprintf("Current todo: '%s'", currentItem)).
				SetDynamicColors(true).
				SetTextAlign(tview.AlignCenter)
			text.SetBorder(true)
			text.SetTitle(" Update ")

			input := tview.NewInputField().
				SetLabel("Update todo: ").
				SetFieldBackgroundColor(tcell.ColorBlack)
			input.SetBorder(true)
			input.SetTitle(" Update To ")
			inputFrame := tview.NewFrame(input).
				SetBorders(2, 2, 2, 2, 2, 2)
			input.SetDoneFunc(func(key tcell.Key) {
				switch key {
				case tcell.KeyEnter:
					t.Todos[currentItemIdx] = input.GetText()
					if err := t.Save(); err != nil {
						panic(err)
					}
					t.Refresh()
					t.App.SetRoot(container, true).SetFocus(t.TodoList)
				case tcell.KeyEsc:
					t.App.SetRoot(container, true).SetFocus(t.TodoList)
				}
			})

			modal := tview.NewFlex().
				SetDirection(tview.FlexRow).
				AddItem(nil, 0, 1, false).
				AddItem(tview.NewFlex().
					AddItem(nil, 0, 1, false).
					AddItem(text, 0, 1, false).
					AddItem(nil, 0, 1, false),
					3, 3, false).
				AddItem(tview.NewFlex().
					AddItem(nil, 0, 1, false).
					AddItem(inputFrame, 0, 1, true).
					AddItem(nil, 0, 1, false),
					10, 1, true).
				AddItem(nil, 0, 1, false)

			t.App.SetRoot(modal, true).SetFocus(input)
		case 'q':
			t.App.Stop()
		}

		return event
	})
}

func NewTodoApp(todoFile string) *TodoApp {
	return &TodoApp{TodoFile: todoFile}
}
