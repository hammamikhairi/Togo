package main

import (
	"flag"
	"fmt"
	. "togo/List"

	"github.com/nsf/termbox-go"
)

const TEST_PATH string = "./test"

var dirPath string

func init() {
	flag.StringVar(&dirPath, "path", "", "the path to a file or directory")
	fmt.Println(dirPath)
}

func main() {
	todos := LoadTodos(TEST_PATH)

	err := termbox.Init()
	if err != nil {
		panic(err)
	}
	defer termbox.Close()

	cursor := CursorInit(todos)

	err = cursor.DrawTodos()
	if err != nil {
		panic(err)
	}

	var act Action
	for act != ACTION_EXIT {
		cursor.LogMode()
		switch ev := termbox.PollEvent(); ev.Type {
		case termbox.EventKey:
			switch cursor.Mode {
			case NRML:
				act = cursor.HandleNRML(ev)
				cursor.DrawTodos()
			case ISRT:
				cursor.HandleISRT(ev)
			}
		case termbox.EventResize:
			cursor.DrawTodos()
		}
	}

	termbox.Close()
	cursor.Save(TEST_PATH)
}
