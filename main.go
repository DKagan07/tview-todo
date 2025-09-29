package main

import (
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"

	"tviewdo/components"
	model "tviewdo/todo"
)

func main() {
	todoApp := model.NewTodoApp("todo.txt")
	app := tview.NewApplication()
	todoApp.App = app

	err := todoApp.PopulateTodos()
	if err != nil {
		panic(err)
	}
	todoApp.TodoList.SetBorder(true)
	todoApp.TodoList.SetTitle(" Todo List ")

	help := components.NewHelpText()

	container := tview.NewFlex().
		SetDirection(tview.FlexRow).
		AddItem(todoApp.TodoList, 0, 1, true).
		AddItem(help, 3, 1, false)

	todoApp.SetKeybindings(container)

	todoApp.App.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		// Close the app with q and ctrl+c
		if event.Key() == tcell.KeyCtrlC {
			todoApp.App.Stop()
		}
		return event
	})

	// Start the app
	if err := app.SetRoot(container, true).Run(); err != nil {
		panic(err)
	}
}
