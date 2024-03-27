package cronjob

import (
	"time"
)

// run in background with goroutine
// using memory storage
type AppWorker struct {
	// function to run
	task Task

	// delay before run
	delay int

	// repeat period
	repeat int
}

type Task = func()

func NewAppWorker() *AppWorker {
	return &AppWorker{}
}

func (aw *AppWorker) SetTask(fn Task) {
	aw.task = fn
}

func (aw *AppWorker) SetDelay(seconds int) {
	aw.delay = seconds
}

func (aw *AppWorker) SetRepeatPeriod(seconds int) {
	aw.repeat = seconds
}

func (aw *AppWorker) Run() {
	go func() {
		defer recoverFunc()
		aw.run()
	}()
}

func (aw *AppWorker) run() {
	time.Sleep(time.Duration(aw.delay) * time.Second)

	aw.task()

	ticker := time.NewTicker(time.Duration(aw.repeat) * time.Second)
	aw.scheduleWorkerTask(ticker)
}

func (aw *AppWorker) scheduleWorkerTask(ticker *time.Ticker) {
	for range ticker.C {
		aw.task()
	}
}
