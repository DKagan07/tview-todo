package components

import "github.com/rivo/tview"

func NewHelpText() *tview.TextView {
	help := tview.NewTextView().
		SetText("q: Quit | a: Add | d: Delete").
		SetDynamicColors(true).
		SetTextAlign(tview.AlignCenter)
	help.SetBorder(true).SetTitle(" Help ")

	return help
}
