package main

import (
	flygo "git.4321.sh/feige/flygo"
	"git.4321.sh/feige/flygo/component/server/httpgovern"
	kin2 "git.4321.sh/feige/flygo/component/server/kin"
	"git.4321.sh/feige/flygo/core/klog"
	"git.4321.sh/feige/flygo/core/ktrace"
	"github.com/gin-gonic/gin"
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
