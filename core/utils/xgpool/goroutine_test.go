package xgpool

import (
	"context"
	"errors"
	"fmt"
	"testing"
	"time"
)

type testStruct struct {
	Num int
}

func Test_Attach(t *testing.T) {
	g := New(context.Background(), 2*time.Second)

	var data1 int
	g.Submit(func(ctx context.Context) (interface{}, error) {
		time.Sleep(time.Second * 1)
		return 3, nil
	}, &data1, nil, "normal_xg")

	data2 := new(testStruct)
	g.Submit(func(ctx context.Context) (interface{}, error) {
		time.Sleep(time.Second * 2)
		return &testStruct{
			Num: 333,
		}, nil
	}, &data2, nil, "timeout_xg")

	var data3 int
	g.Submit(func(ctx context.Context) (interface{}, error) {
		time.Sleep(time.Second * 3)
		panic("amsaklsa")
		return 4, nil
	}, &data3, nil, "timeout_xg_panic")

	var data4 int
	var err4 error
	g.Submit(func(ctx context.Context) (interface{}, error) {
		time.Sleep(time.Second * 3)
		return 0, errors.New("invalid data error")
	}, &data4, &err4, "timeout_err")

	var data5 int
	var err5 error
	g.Submit(func(ctx context.Context) (interface{}, error) {
		time.Sleep(time.Second * 1)
		return 0, errors.New("invalid data error")
	}, &data5, &err5, "normal_err")

	g.Wait()

	fmt.Println(data1)
	fmt.Println(data2)
	fmt.Println(data3)
	fmt.Println(data4)
	fmt.Println(err4)
	fmt.Println(errors.Is(err4, TimeoutErr))

	fmt.Println(data5)
	fmt.Println(err5)
	fmt.Println(errors.Is(err5, TimeoutErr))
}
