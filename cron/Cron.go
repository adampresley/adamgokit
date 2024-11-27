package cron

import (
	"context"

	cronv3 "github.com/robfig/cron/v3"
)

var internalCron *cronv3.Cron

func init() {
	internalCron = cronv3.New()
}

/*
Add adds a function to run on a schedule. The schedule is defined by the
standard cron syntax. See https://crontab.guru
*/
func Add(schedule string, cronFunc func()) {
	internalCron.AddFunc(schedule, cronFunc)
}

/*
Start starts the cron scheduler and runner. This is asyncronous.
*/
func Start() {
	internalCron.Start()
}

/*
Stop attempts to stop all running cron processes.
*/
func Stop() context.Context {
	return internalCron.Stop()
}
