package runner

import (
	"context"
	"os"
	"sync"
)

type Runner interface {
	Ctx() context.Context
	AddProcess(int)
	DoneProcess(string)
	AwaitStop()
}

type app struct {
	ctx  context.Context
	stop context.CancelFunc
}

type processes struct {
	wg   sync.WaitGroup
	errs chan error
}

type stop struct {
	exist  bool
	signal chan os.Signal
}

type run struct {
	app       app
	processes processes
	stop      stop
}

const (
	colorBlack = iota + 30
	colorRed
	colorGreen
	colorYellow
	colorBlue
	colorMagenta
	colorCyan
	colorWhite

	colorBold     = 1
	colorDarkGray = 90
)
