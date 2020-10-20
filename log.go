package main

import "github.com/gookit/color"

type log struct {
	color   *color.Color
	address string
}

func newLog(address string, idx int) *log {
	colors := []color.Color{color.Green, color.Blue, color.Magenta, color.Cyan, color.Yellow, color.Red, color.Gray}
	return &log{address: address, color: &colors[idx%len(colors)]}
}

func (l *log) printLog(color color.Color, data ...interface{}) {
	l.color.Printf("%s | %s\n", l.address, color.Render(data...))
}

func (l *log) debug(data ...interface{}) {
	l.printLog(color.Gray, data...)
}

func (l *log) info(data ...interface{}) {
	l.printLog(color.White, data...)
}

func (l *log) warn(data ...interface{}) {
	l.printLog(color.Yellow, data...)
}

func (l *log) error(data ...interface{}) {
	l.printLog(color.Red, data...)
}
