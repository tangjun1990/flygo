package main

import (
	"context"
	"errors"
	"fmt"

	"github.com/tangjun1990/flygo"
	kredis2 "github.com/tangjun1990/flygo/component/client/kredis"
	kcronlock2 "github.com/tangjun1990/flygo/component/client/kredis/kcronlock"
	kcron2 "github.com/tangjun1990/flygo/component/task/kcron"
	"github.com/tangjun1990/flygo/core/klog"
	"go.uber.org/zap"
)

var (
	redis  *kredis2.Component
	locker *kcronlock2.Component
)

func main() {
	if err := flygo.New().Invoker(initRedis).Cron(cronJob1(), cronJob2()).Run(); err != nil {
		klog.Error("start up", zap.Error(err))
	}
}

// 异常任务
func cronJob1() kcron2.Kcron {
	job := func(ctx context.Context) error {
		fmt.Println("run job1")
		return errors.New("exec job1 error")
	}

	cron := kcron2.Load("cron.testcron1").Build(kcron2.WithJob(job))
	return cron
}

// 正常任务
func cronJob2() kcron2.Kcron {
	job := func(ctx context.Context) error {
		fmt.Println("run job2")
		return nil
	}

	cron := kcron2.Load("cron.testcron2").Build(kcron2.WithJob(job))
	return cron
}

func initRedis() error {
	redis = kredis2.Load("redis.cron").Build()
	locker = kcronlock2.DefaultContainer().Build(kcronlock2.WithClient(redis))
	return nil
}

// 分布式任务
func distributeCron() kcron2.Kcron {
	cron := kcron2.Load("cron.discron").Build(
		// 设置分布式锁
		kcron2.WithLock(locker.NewLock("lock:discron:hello")),
		kcron2.WithJob(helloWorld),
	)
	return cron
}

func helloWorld(ctx context.Context) error {
	fmt.Println("run distributeCron")
	return nil
}
