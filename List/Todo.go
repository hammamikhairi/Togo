package list

import "fmt"

type Todo struct {
	Content, Desc string
	Done          bool
	// Pos        Origin
}

func (td *Todo) getContent() string {
	var pref string
	if td.Done {
		pref = fmt.Sprintf(PREF, "x")
	} else {
		pref = fmt.Sprintf(PREF, " ")
	}

	return pref + td.Content
}

func (td *Todo) getSave() string {
	var pref string
	if td.Done {
		pref = "DONE>>"
	} else {
		pref = "TODO>>"
	}

	return pref + td.Content + ">>" + td.Desc
}