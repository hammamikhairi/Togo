package list

type Log struct {
	msg string
	pos Vec2i
}

type Logger struct {
	logs []*Log
}
