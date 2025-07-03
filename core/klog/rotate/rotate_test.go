package rotate_test

import (
	"log"
	"os"
	"os/signal"
	"syscall"
	"testing"

	"github.com/stretchr/testify/assert"

	"git.4321.sh/feige/flygo/core/klog/rotate"
)

// Example of how to rotate in response to SIGHUP.
func TestLoggerRotate(t *testing.T) {
	l := &rotate.Logger{}
	log.SetOutput(l)
	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGHUP)

	go func() {
		for {
			<-c
			err := l.Rotate()
			assert.NoError(t, err)
		}
	}()

	t.Log("test done")
}
