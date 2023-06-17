package list

import "github.com/nsf/termbox-go"

const (
	PADDING_TOP  int    = 3
	PADDING_LEFT int    = 0
	PREF         string = "- [%s] "
)

func drawToTerminal(content string, x, y int, fg, bg termbox.Attribute) {
	for _, ch := range content {
		termbox.SetCell(x, y, ch, fg, bg)
		x++
	}
}

func terminalLog(msg string) {
	drawToTerminal(msg, 1, 1, termbox.ColorDefault, termbox.ColorDefault)
	termbox.Flush()
}

func clear() {
	// TODO : Draw the titles and shit
	termbox.Clear(termbox.ColorDefault, termbox.ColorDefault)
}

func drawSingleTodo(task *Todo, x, y int, fg, bg termbox.Attribute) {
	// TODO : Draw Todo with description
	// Description is limited to the bounds
	// of the Todo block
	// XXX : maybe have the Desc apear at the very bottom ?
	drawToTerminal(task.getContent(), x+PADDING_LEFT, y+PADDING_TOP, fg, bg)
}
