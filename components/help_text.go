package components

import "github.com/rivo/tview"

func NewHelpText() *tview.TextView {
	help := tview.NewTextView().
		SetText("q: Quit | a: Add | d: Delete | u: Update").
		SetDynamicColors(true).
		SetTextAlign(tview.AlignCenter)
	help.SetBorder(true).SetTitle(" Help ")

	return help
}
