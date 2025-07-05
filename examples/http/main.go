package main

import (
	"github.com/gin-gonic/gin"
	flygo "github.com/tangjun1990/flygo"
	"github.com/tangjun1990/flygo/component/server/httpgovern"
	kin2 "github.com/tangjun1990/flygo/component/server/kin"
	"github.com/tangjun1990/flygo/core/klog"
	"github.com/tangjun1990/flygo/core/ktrace"
)

// go run main.go --config=config.toml, then visit http://127.0.0.1:9009/api/v1/hello
func main() {
	if err := flygo.New().
		Serve(httpgovern.Load("server.httpgovern").Build(),
			NewServer(),
		).Run(); err != nil {
		klog.Panic("startup", klog.FieldErr(err))
	}
}

func NewServer() *kin2.Component {
	server := kin2.Load("server.http").Build()
	apiV1 := server.Group("/api/v1")

	apiV1.Use()
	{
		apiV1.GET("/hello", HandleHello)
	}

	return server
}

func HandleHello(ctx *gin.Context) {

	ctx.JSON(200, gin.H{
		"trace_id": ktrace.ExtractTraceID(ctx.Request.Context()),
	})
}
