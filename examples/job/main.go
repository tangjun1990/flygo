package main

import (
	"context"
	"fmt"

	flygo "github.com/tangjun1990/flygo"
	kjob2 "github.com/tangjun1990/flygo/component/task/kjob"
	"github.com/tangjun1990/flygo/core/klog"
	"go.uber.org/zap"
)

// go run main.go --job=demojob1,demojob2
func main() {
	if err := flygo.New().Job(NewJobRunner1(), NewJobRunner2()).Run(); err != nil {
		klog.Error("start up", zap.Error(err))
	}
}

func NewJobRunner1() *kjob2.Component {
	return kjob2.DefaultContainer().Build(
		kjob2.WithName("demojob1"),
		kjob2.WithStartFunc(job1),
	)
}

func NewJobRunner2() *kjob2.Component {
	return kjob2.DefaultContainer().Build(
		kjob2.WithName("demojob2"),
		kjob2.WithStartFunc(job2),
	)
}

func job1(ctx context.Context) error {
	fmt.Println("i am demojob1")
	klog.WithCtx(ctx).Info("i am demojob1")
	return nil
}

func job2(ctx context.Context) error {
	fmt.Println("i am demojob2")
	klog.WithCtx(ctx).Info("i am demojob2")
	return nil
}
