package sglog

import (
	"fmt"
	"io"
	"log"
	"os"
	"runtime"
)

var (
	Trace   Logger
	Debug   Logger
	Info    Logger
	Warning Logger
	Error   Logger
)

func init() {
	SetTrace(nil)
	SetDebug(nil)
	SetInfo(nil)
	SetWarning(os.Stderr)
	SetError(os.Stderr)
}

func setLogger(logger *Logger, handle io.Writer, prefix string) {
	if handle == nil {
		logger.Enabled = false
		logger.delegate = nil
	} else {
		logger.Enabled = true
		logger.delegate = log.New(
			handle,
			fmt.Sprintf("[%s] ", prefix),
			log.Ldate|log.Ltime|log.Lmicroseconds|log.Lshortfile)
	}
}

func SetTrace(handle io.Writer) {
	setLogger(&Trace, handle, "TRC")
}

func SetDebug(handle io.Writer) {
	setLogger(&Debug, handle, "DBG")
}

func SetInfo(handle io.Writer) {
	setLogger(&Info, handle, "INF")
}

func SetWarning(handle io.Writer) {
	setLogger(&Warning, handle, "WRN")
}

func SetError(handle io.Writer) {
	setLogger(&Error, handle, "ERR")
}

//----------------------------------------------------------------------------
// Logger
//----------------------------------------------------------------------------

type Logger struct {
	Enabled  bool
	delegate *log.Logger
}

func getCallerName() string {
	pc, _, _, _ := runtime.Caller(2)
	f := runtime.FuncForPC(pc)
	return f.Name()
}

func (l *Logger) PrintEnter() {
	if l.Enabled {
		l.delegate.Printf("=> %s", getCallerName())
	}
}

func (l *Logger) PrintEnterAnon(desc string) {
	if l.Enabled {
		l.delegate.Printf("=> %s (%s)", getCallerName(), desc)
	}
}

func (l *Logger) PrintLeave() {
	if l.Enabled {
		l.delegate.Printf("<= %s", getCallerName())
	}
}

func (l *Logger) PrintLeaveAnon(desc string) {
	if l.Enabled {
		l.delegate.Printf("<= %s (%s)", getCallerName(), desc)
	}
}

func (l *Logger) Println(args ...interface{}) {
	if l.Enabled {
		l.delegate.Println(args...)
	}
}

func (l *Logger) Printf(format string, args ...interface{}) {
	if l.Enabled {
		l.delegate.Printf(format, args...)
	}
}

func (l *Logger) PrintStack(all bool) {
	if l.Enabled {
		buf := make([]byte, 512)
		for {
			n := runtime.Stack(buf, true)
			if n < len(buf) {
				break
			}
			buf = make([]byte, len(buf)*2)
		}
		l.Printf("Stack (all goroutines: %v)\n%v", all, string(buf))
	}
}
