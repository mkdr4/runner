package runner

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
)

func Setup() Runner {
	appCtx, appStop := context.WithCancel(context.Background())

	return &run{
		app: app{
			ctx:  appCtx,
			stop: appStop,
		},
		processes: processes{
			wg:   sync.WaitGroup{},
			errs: make(chan error),
		},
		stop: stop{
			exist:  false,
			signal: make(chan os.Signal, 1),
		},
	}
}

func (r *run) Ctx() context.Context { return r.app.ctx }

func (r *run) AddProcess(i int) { r.processes.wg.Add(i) }

func (r *run) DoneProcess() {
	if err := recover(); err != nil {
		r.stop.signal <- os.Interrupt
		r.processes.wg.Add(1)
		r.processes.errs <- fmt.Errorf("%v", err)
	}

	r.processes.wg.Done()
}

func (r *run) AwaitStop() {
	<-r.stopSignal()
	fmt.Printf("%s %s %s\n",
		colorize(time.Now().Format(time.TimeOnly), colorDarkGray, true),
		colorize("INF", colorGreen, true),
		"runner stoped",
	)

	logCtx, logStop := context.WithCancel(context.Background())
	defer logStop()
	go r.logErrors(logCtx)

	r.stopApplication()
}

func (r *run) stopSignal() chan os.Signal {
	signal.Notify(r.stop.signal, os.Interrupt, syscall.SIGTERM)
	return r.stop.signal
}

func (r *run) stopApplication() {
	r.app.stop()
	r.processes.wg.Wait()
}

func (r *run) logErrors(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			return
		case v := <-r.processes.errs:
			fmt.Printf("%s %s %s\n",
				colorize(time.Now().Format(time.TimeOnly), colorDarkGray, true),
				colorize("ERR", colorRed, true),
				colorize(v, colorRed, true),
			)
			r.processes.wg.Done()
		}
	}
}

func colorize(s interface{}, c int, enable bool) string {
	if !enable {
		return fmt.Sprintf("%s", s)
	}
	return fmt.Sprintf("\x1b[%dm%v\x1b[0m", c, s)
}
