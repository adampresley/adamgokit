package cron

import (
	"context"

	cronv3 "github.com/robfig/cron/v3"
)

var internalCron *cronv3.Cron

func init() {
	internalCron = cronv3.New()
}

func Add(schedule string, cronFunc func()) {
	internalCron.AddFunc(schedule, cronFunc)
}

func Start() {
	internalCron.Start()
}

func Stop() context.Context {
	return internalCron.Stop()
}
