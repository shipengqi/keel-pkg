package recovery

import (
	"bytes"
	"errors"
	"fmt"
	"io/ioutil"
	"runtime"
)

var (
	unknown = []byte("???")
	dot     = []byte(".")
	slash   = []byte("/")
)

type RecoverFunc func(err error, errStack string)

func Recovery(rf RecoverFunc) func() {
	return func() {
		if re := recover(); re != nil {
			var err error
			switch x := re.(type) {
			case string:
				err = errors.New(x)
			case error:
				err = x
			default:
				err = errors.New("")
			}
			stackStr := stack2(3)
			rf(err, string(stackStr))
		}
	}
}

// stack returns a nicely formatted stack frame, skipping skip frames.
// Output example:
// runtime error: invalid memory address or nil pointer dereference
// /usr/local/go/src/runtime/panic.go:212 (0x435a1a)
//        panicmem: // for deferproc does not describe them. So we can't let garbage
// /usr/local/go/src/runtime/signal_unix.go:734 (0x44e492)
//        sigpanic: signalstack(&_g_.m.gsignal.stack)
// /root/gowork/src/cli/internal/action/settings.go:83 (0x1237546)
// /root/gowork/src/cli/main.go:27 (0x123bef2)
// /root/gowork/pkg/mod/github.com/spf13/cobra@v1.1.3/command.go:882 (0x5b2128)
// /root/gowork/pkg/mod/github.com/spf13/cobra@v1.1.3/command.go:818 (0x5b172e)
// /root/gowork/pkg/mod/github.com/spf13/cobra@v1.1.3/command.go:960 (0x5b25d4)
// /root/gowork/pkg/mod/github.com/spf13/cobra@v1.1.3/command.go:897 (0x123bd64)
// /root/gowork/src/cli/main.go:37 (0x123bd52)
// /usr/local/go/src/runtime/proc.go:225 (0x43a615)
//        main: exit(0)
// /usr/local/go/src/runtime/asm_amd64.s:1371 (0x46d740)
//        goexit: // gcWriteBarrier performs a heap pointer write and informs the GC.
func stack(skip int) []byte {
	buf := new(bytes.Buffer) // the returned data
	// As we loop, we open files and read them. These variables record the currently
	// loaded file.
	var lines [][]byte
	var lastFile string
	for i := skip; ; i++ { // Skip the expected number of frames
		pc, file, line, ok := runtime.Caller(i)
		if !ok {
			break
		}
		// Print this much at least.  If we can't find the source, it won't show.
		_, _ = fmt.Fprintf(buf, "%s:%d (0x%x)\n", file, line, pc)
		if file != lastFile {
			data, err := ioutil.ReadFile(file)
			if err != nil {
				continue
			}
			lines = bytes.Split(data, []byte{'\n'})
			lastFile = file
		}
		_, _ = fmt.Fprintf(buf, "\t%s: %s\n", function(pc), source(lines, line))
	}
	return buf.Bytes()
}

// stack2 returns a nicely formatted stack frame, skipping skip frames.
// Output example:
// runtime error: invalid memory address or nil pointer dereference
//        runtime.gopanic: /usr/local/go/src/runtime/panic.go:965
//        runtime.panicmem: /usr/local/go/src/runtime/panic.go:212
//        runtime.sigpanic: /usr/local/go/src/runtime/signal_unix.go:734
//        github.com/shipengqi/example.v1/cli/internal/action.(*Configuration).Init: /root/gowork/src/cli/internal/action/settings.go:83
//        main.main.func1: /root/gowork/src/cli/main.go:27
//        github.com/spf13/cobra.(*Command).preRun: /root/gowork/pkg/mod/github.com/spf13/cobra@v1.1.3/command.go:882
//        github.com/spf13/cobra.(*Command).execute: /root/gowork/pkg/mod/github.com/spf13/cobra@v1.1.3/command.go:818
//        github.com/spf13/cobra.(*Command).ExecuteC: /root/gowork/pkg/mod/github.com/spf13/cobra@v1.1.3/command.go:960
//        : /root/gowork/pkg/mod/github.com/spf13/cobra@v1.1.3/command.go:897
//        main.main: /root/gowork/src/cli/main.go:37
//        runtime.main: /usr/local/go/src/runtime/proc.go:225
func stack2(skip int) []byte {
	buf := new(bytes.Buffer) // the returned data
	st := make([]uintptr, 32)
	count := runtime.Callers(skip, st)
	callers := st[:count]
	frames := runtime.CallersFrames(callers)
	for {
		frame, ok := frames.Next()
		if !ok {
			break
		}
		_, _ = fmt.Fprintf(buf, "\t%s: %s:%d\n", frame.Func.Name(), frame.File, frame.Line)
	}
	return buf.Bytes()
}

// function returns, if possible, the name of the function containing the PC.
func function(pc uintptr) []byte {
	fn := runtime.FuncForPC(pc)
	if fn == nil {
		return []byte("???")
	}
	name := []byte(fn.Name())
	// The name includes the path name to the package, which is unnecessary
	// since the file name is already included.  Plus, it has center dots.
	// That is, we see
	//	runtime/debug.*T·ptrmethod
	// and want
	//	*T.ptrmethod
	// Also the package path might contains dot (e.g. code.google.com/...),
	// so first eliminate the path prefix
	if lastSlash := bytes.LastIndex(name, slash); lastSlash >= 0 {
		name = name[lastSlash+1:]
	}
	if period := bytes.Index(name, dot); period >= 0 {
		name = name[period+1:]
	}
	name = bytes.Replace(name, dot, dot, -1)
	return name
}

// source returns a space-trimmed slice of line.
func source(lines [][]byte, n int) []byte {
	n-- // in stack trace, lines are 1-indexed but our array is 0-indexed
	if n < 0 || n >= len(lines) {
		return unknown
	}
	return bytes.TrimSpace(lines[n])
}

// func timeFormat(t time.Time) string {
// 	var timeString = t.Format("2006/01/02 - 15:04:05")
// 	return timeString
// }
