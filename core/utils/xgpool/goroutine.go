package xgpool

import (
	"context"
	"errors"
	"reflect"
	"sort"
	"time"

	"github.com/tangjun1990/flygo/core/klog"
)

var TimeoutErr = errors.New("xgpool timeout")

type response struct {
	val interface{}
	err error
}

type xgoroutine struct {
	channel chan *response
	ctx     context.Context
	val     interface{}
	err     *error
	action  string
}

type xgoroutines []*xgoroutine

func (g xgoroutines) Len() int {
	return len(g)
}

func (g xgoroutines) Less(i, j int) bool {
	iTime, _ := g[i].ctx.Deadline()
	jTime, _ := g[j].ctx.Deadline()

	return iTime.Before(jTime)
}

func (g xgoroutines) Swap(i, j int) {
	g[i], g[j] = g[j], g[i]
}

func (g *xgoroutine) setError(err error) {
	if g.err != nil {
		*g.err = err
	}
}

func (g *xgoroutine) setResp(resp *response) {
	rVal := reflect.ValueOf(g.val)
	if rVal.Kind() != reflect.Ptr {
		klog.Warn("xgpool timeout.wait, val must be ptr")
		return
	}

	if resp.err != nil {
		g.setError(resp.err)
		return
	}

	resultVal := reflect.ValueOf(resp.val)
	resultValType := resultVal.Type()

	if !rVal.IsNil() && resultValType != nil {
		rValElem := rVal.Elem()
		if rValElem.Type().String() != resultValType.String() {
			klog.Warnf("xgpool result type not equal, %s, %s", rValElem.Type().String(), resultVal.Type().String())
			return
		}
		rValElem.Set(resultVal)
	}
}

type xgpool struct {
	Ctx         context.Context
	Timeout     time.Duration
	waitChannel xgoroutines
}

func New(ctx context.Context, timeout time.Duration) *xgpool {
	return &xgpool{
		Ctx:     ctx,
		Timeout: timeout,
	}
}

const DefaultTimeout = time.Second * 5

type routineTimeoutFunc func(context.Context) (interface{}, error)

func (xg *xgpool) Submit(method routineTimeoutFunc, val interface{}, err *error, action string) {
	xg.submitJob(method, val, err, action)
}

func (xg *xgpool) submitJob(method routineTimeoutFunc, val interface{}, errVal *error, action string) {
	var ctx context.Context

	if xg.Timeout == 0 {
		xg.Timeout = DefaultTimeout
	}

	ctx, _ = context.WithTimeout(xg.Ctx, xg.Timeout)

	// 创建数据信道
	channel := make(chan *response, 1)
	xg.waitChannel = append(xg.waitChannel, &xgoroutine{
		channel: channel,
		ctx:     ctx,
		val:     val,
		err:     errVal,
		action:  action,
	})

	go xg.run(ctx, channel, method, action)
}

// 开启协程，执行方法
func (xg *xgpool) run(ctx context.Context, channel chan *response, method routineTimeoutFunc, action string) {
	defer GoRecover(action)
	defer close(channel)

	if action == "" {
		action = "unknow"
	}

	result, err := method(ctx)

	if err != nil {
		klog.Warnf("xgpool err: %s, action: %s", err.Error(), action)
	}

	xg.writeChannelStatus(channel, &response{
		val: result,
		err: err,
	})
}

func (xg *xgpool) writeChannelStatus(channel chan *response, result *response) {
	select {
	case channel <- result:
	default:
	}
}

func (xg *xgpool) Wait() {
	sort.Sort(xg.waitChannel)
	for _, routine := range xg.waitChannel {
		select {
		// 正常返回数据
		case resp := <-routine.channel:
			if resp != nil {
				routine.setResp(resp)
			}

		// context强制结束
		case <-routine.ctx.Done():
			// 还需要再读一次channel，因为select命中多个条件是无序选择的，有可能超时和返回数据同时命中
			select {
			case resp := <-routine.channel:
				if resp != nil {
					routine.setResp(resp)
				}
			default:
				routine.setError(TimeoutErr)
				klog.Warnf("xgpool timeout, action: %s", routine.action)
			}
		}
	}
}

type wait interface {
	Wait()
}

func WaitMulti(waits ...wait) {
	for _, w := range waits {
		w.Wait()
	}
}

func GoRecover(action string) {
	err := recover()
	if err != nil {
		klog.Errorf("xgpool panic: %v, action: %s", err, action)
	}
}
