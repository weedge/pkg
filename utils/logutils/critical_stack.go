package logutils

import (
	"fmt"
	"runtime/debug"
	"strings"

	"github.com/cloudwego/kitex/pkg/klog"
)

func Criticalf(format string, args ...interface{}) {
	stack := string(debug.Stack())
	stack = strings.ReplaceAll(stack, "\n\t", "]<-")
	stack = strings.ReplaceAll(stack, "\n", "  [")
	klog.Errorf("[CRITICAL] stack %s | %s", stack, fmt.Sprintf(format, args...))
}

func CriticalError(err error) {
	Criticalf(err.Error())
}

func CriticalIfError(err error) {
	if err != nil {
		CriticalError(err)
	}
}
