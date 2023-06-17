package list

import (
	"bufio"
	"log"
	"os"
	"path/filepath"
	"strings"
)

type Vec2i [2]int

func assert(condition bool, message string) {
	if condition == false {
		println()
		log.Fatal(message)
	}
}

func sortTodos(todos []*Todo) [][]*Todo {
	var left, right []*Todo
	for _, task := range todos {
		switch task.Done {
		case true:
			right = append(right, task)
		case false:
			left = append(left, task)
		}
	}
	return [][]*Todo{left, right}
}

func GetFileFullPath(path string) string {
	return filepath.Clean(path)
}

func LoadTodos(path string) []*Todo {

	var tods []*Todo
	file, err := os.Open(path)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		todoMeta := strings.Split(line, ">>")
		if len(todoMeta) != 3 {
			continue
		}
		done := true
		if todoMeta[0] == "TODO" {
			done = false
		}
		tods = append(tods, &Todo{Content: todoMeta[1], Desc: todoMeta[2], Done: done})
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
	return tods
}
