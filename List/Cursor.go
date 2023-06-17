package list

import (
	"fmt"
	"os"

	"github.com/nsf/termbox-go"
)

type Mode uint8

const (
	ISRT Mode = iota
	NRML
)

type Action uint8

const (
	ACTION_NONE Action = iota
	ACTION_EXIT
)

type Cursor struct {
	Pos    Vec2i
	Todos  [][]*Todo
	Fg, Bg termbox.Attribute
	*Logger
	Mode
}

func CursorInit(todos []*Todo) *Cursor {

	return &Cursor{
		Pos:   [2]int{0, 0},
		Fg:    termbox.ColorDarkGray,
		Bg:    termbox.ColorWhite,
		Todos: sortTodos(todos),
		Mode:  NRML,
	}
}

func (cr *Cursor) getCurrent() *Todo {
	if len(cr.Todos[cr.Pos[0]]) == 0 {
		return nil
	}
	return cr.Todos[cr.Pos[0]][cr.Pos[1]]
}

func (cr *Cursor) DrawTodos() error {
	var (
		x, w, colsWidth int
		fg, bg          termbox.Attribute
	)

	w, _ = termbox.Size()
	colsWidth = w / 2

	clear()

	for col, colTodos := range cr.Todos {
		for y, todo := range colTodos {
			x = col * colsWidth
			fg, bg = termbox.ColorDefault, termbox.ColorDefault
			if cr.Pos[0] == col && cr.Pos[1] == y {
				fg, bg = cr.Fg, cr.Bg
				if cr.Mode == ISRT {
					bg = termbox.ColorLightYellow
				}
			}
			drawSingleTodo(todo, x, y, fg, bg)
		}
	}

	if cr.Pos[0] == 0 {
		drawToTerminal("<TODO>", 0, PADDING_TOP-1, cr.Fg, cr.Bg)
		drawToTerminal("<DONE>", colsWidth, PADDING_TOP-1, cr.Fg, termbox.ColorDefault)
	} else {
		drawToTerminal("<TODO>", 0, PADDING_TOP-1, cr.Fg, termbox.ColorDefault)
		drawToTerminal("<DONE>", colsWidth, PADDING_TOP-1, cr.Fg, cr.Bg)
	}

	// _, h := termbox.Size()
	// if curr := cr.getCurrent(); curr != nil {
	// 	drawToTerminal(curr.Desc, 0, h-1, termbox.ColorDefault, termbox.ColorDefault)
	// }

	return termbox.Flush()
}
func (cr *Cursor) Move() {

	if len(cr.Todos[cr.Pos[0]]) == 0 {
		return
	}

	curr := cr.getCurrent()
	switch cr.Pos[0] {
	case 0:
		curr.Done = true
	case 1:
		curr.Done = false
	}
	cr.Todos[(cr.Pos[0]+1)%2] = append(cr.Todos[(cr.Pos[0]+1)%2], cr.Todos[cr.Pos[0]][cr.Pos[1]])
	cr.Todos[cr.Pos[0]] = append(cr.Todos[cr.Pos[0]][:cr.Pos[1]], cr.Todos[cr.Pos[0]][cr.Pos[1]+1:]...)

	if cr.Pos[1] != 0 {
		cr.Pos[1] -= 1
	}
}

func (cr *Cursor) MoveTodoDown() {
	if cr.Pos[1]-1 < 0 {
		return
	}

	cr.Todos[cr.Pos[0]][cr.Pos[1]], cr.Todos[cr.Pos[0]][cr.Pos[1]-1] = cr.Todos[cr.Pos[0]][cr.Pos[1]-1], cr.Todos[cr.Pos[0]][cr.Pos[1]]
	cr.Pos[1]--
}

func (cr *Cursor) MoveTodoUp() {
	if cr.Pos[1]+1 >= len(cr.Todos[cr.Pos[0]]) {
		return
	}
	cr.Todos[cr.Pos[0]][cr.Pos[1]], cr.Todos[cr.Pos[0]][cr.Pos[1]+1] = cr.Todos[cr.Pos[0]][cr.Pos[1]+1], cr.Todos[cr.Pos[0]][cr.Pos[1]]
	cr.Pos[1]++
}

func (cr *Cursor) goUp() {
	if len(cr.Todos[cr.Pos[0]]) == 0 {
		return
	}
	cr.Pos[1] = (cr.Pos[1] + len(cr.Todos[cr.Pos[0]]) - 1) % len(cr.Todos[cr.Pos[0]])
}
func (cr *Cursor) goDown() {
	if len(cr.Todos[cr.Pos[0]]) == 0 {
		return
	}
	cr.Pos[1] = (cr.Pos[1] + 1) % len(cr.Todos[cr.Pos[0]])
}

func (cr *Cursor) Delete() {
	if len(cr.Todos[cr.Pos[0]]) == 0 {
		return
	}
	cr.Todos[cr.Pos[0]] = append(cr.Todos[cr.Pos[0]][:cr.Pos[1]], cr.Todos[cr.Pos[0]][cr.Pos[1]+1:]...)
	if cr.Pos[1] != 0 {

		cr.Pos[1]--
	}
}

func (cr *Cursor) LogMode() {
	// Log Mode ISRT NRML
	if cr.Mode == ISRT {
		drawToTerminal("ISRT", 0, 0, termbox.ColorDefault, termbox.ColorBlue)
	}

	if cr.Mode == NRML {
		drawToTerminal("NRML", 0, 0, termbox.ColorDefault, termbox.ColorGreen)
	}
	termbox.Flush()
}

func (cr *Cursor) Switch() {
	left, right := len(cr.Todos[0]), len(cr.Todos[1])

	switch cr.Pos[0] {
	case 0:
		// left -> right
		cr.Pos[0] = 1
		if cr.Pos[1] >= right {
			if right == 0 {
				cr.Pos[1] = right
			} else {
				cr.Pos[1] = right - 1
			}
		}
	case 1:
		// right -> left
		cr.Pos[0] = 0
		if cr.Pos[1] >= left {
			if left == 0 {
				cr.Pos[1] = left
			} else {
				cr.Pos[1] = left - 1
			}
		}
	}
}

func (cr *Cursor) SetCursor() {
	x, y := cr.getTerminalCursorPos()
	termbox.SetCursor(x, y)
	termbox.Flush()
}

func (cr *Cursor) getTerminalCursorPos() (int, int) {
	var (
		y int = cr.Pos[1] + PADDING_TOP
		x int
	)

	x = len(cr.getCurrent().getContent())
	if cr.Pos[0] == 1 {
		w, _ := termbox.Size()
		x += w / 2
	}

	return x, y
}

func (cr *Cursor) SetCurrentCusrorCell(ch rune, fg, bg termbox.Attribute) {
	x, y := cr.getTerminalCursorPos()
	termbox.SetCell(x-1, y, ch, fg, bg)
	termbox.Flush()
}

func (cr *Cursor) AddTODO() {
	todo := Todo{
		Content: "",
		Desc:    "",
		Done:    false,
	}

	if cr.Pos[0] == 1 {
		todo.Done = true
	}

	cr.Todos[cr.Pos[0]] = append(cr.Todos[cr.Pos[0]], &todo)
	cr.Mode = ISRT
	cr.Pos[1] = len(cr.Todos[cr.Pos[0]]) - 1
	cr.DrawTodos()
}

func (cr *Cursor) HandleNRML(ev termbox.Event) Action {
	if ev.Key == termbox.KeyEsc || ev.Ch == 'q' || ev.Ch == 'Q' {
		return ACTION_EXIT
	}
	if ev.Ch == 'r' {
		cr.Mode = ISRT
		cr.SetCursor()
		return ACTION_NONE
	}
	switch ev.Key {
	case termbox.KeyArrowUp:
		cr.goUp()
	case termbox.KeyArrowDown:
		cr.goDown()
	case termbox.KeyTab:
		cr.Switch()
	case termbox.KeyEnter:
		cr.Move()
	default:
		switch ev.Ch {
		case 'j':
			cr.MoveTodoUp()
		case 'k':
			cr.MoveTodoDown()
		case 'd':
			cr.Delete()
		case 'a':
			cr.AddTODO()
		}
	}
	return ACTION_NONE
}

func (cr *Cursor) HandleISRT(ev termbox.Event) {
	if ev.Key == termbox.KeyEsc || ev.Key == termbox.KeyEnter {
		cr.Mode = NRML
		termbox.HideCursor()
		cr.DrawTodos()
		return
	}

	if ev.Key == 127 {
		if len(cr.getCurrent().Content) > 0 {
			cr.SetCurrentCusrorCell(' ', termbox.ColorDefault, termbox.ColorDefault)
			cr.getCurrent().Content = cr.getCurrent().Content[:len(cr.getCurrent().Content)-1]
		}
	} else {
		cr.getCurrent().Content += string(ev.Ch)
		cr.SetCurrentCusrorCell(ev.Ch, cr.Fg, termbox.ColorLightYellow)
	}
	cr.SetCursor()
}

func (cr *Cursor) Save(path string) error {
	file, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		return err
	}
	defer file.Close()

	for _, colTodos := range cr.Todos {
		for _, todo := range colTodos {
			_, err = file.WriteString(todo.getSave() + "\n")
			if err != nil {
				return err
			}
		}
	}
	fmt.Println("saved to ./test")
	return nil
}