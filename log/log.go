package log

import (
	"fmt"
	"io"
	"log/slog"
	"os"
	"runtime"
	"time"

	sf "github.com/samber/slog-formatter"
	"gopkg.in/natefinch/lumberjack.v2"
)

var logger = slog.New(DefaultHandler(slog.LevelInfo))

func SetLogger(l *slog.Logger) {
	logger = l
}

func Logger() *slog.Logger {
	return logger
}

func Debug(msg string, args ...any) {
	logger.Debug(msg, args...)
}

func Info(msg string, args ...any) {
	logger.Info(msg, args...)
}

func Warn(msg string, args ...any) {
	logger.Warn(msg, args...)
}

func Error(msg string, args ...any) {
	logger.Error(msg, args...)
}

func PanicF(msg string, args ...any) {
	// 使用 defer 捕获 panic
	defer func() {
		if r := recover(); r != nil {
			logger.Error("Recovered from panic", "msg", msg, "args", args, "error", r)
			// 记录调用栈
			buf := make([]byte, 1<<16) // 64KB
			stackSize := runtime.Stack(buf, true)
			logger.Error("Stack trace", "stack", string(buf[:stackSize]))
		}
	}()

	// 触发 panic
	panic(fmt.Sprintf(msg, args...))
}

func DefaultHandler(level slog.Level) slog.Handler {
	cn, err := time.LoadLocation("Asia/Shanghai")
	if err != nil {
		// 处理加载时区的错误
		fmt.Printf("Error loading time location, error: %+v", err)
		cn = time.UTC // 或者使用默认的 UTC
	}
	return sf.NewFormatterHandler(sf.TimeFormatter(time.DateTime, time.UTC))(
		slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
			Level: level,
			ReplaceAttr: func(_ []string, a slog.Attr) slog.Attr {
				if a.Key != "time" {
					return a
				}
				newTime := a.Value.Time().In(cn)
				return slog.Time(a.Key, newTime)
			},
		}),
	)
}

func JSONHandler(w io.Writer, level slog.Level) slog.Handler {
	cn, _ := time.LoadLocation("Asia/Shanghai")
	return sf.NewFormatterHandler(sf.TimeFormatter(time.DateTime, time.UTC))(
		slog.NewJSONHandler(w, &slog.HandlerOptions{
			Level: level,
			ReplaceAttr: func(_ []string, a slog.Attr) slog.Attr {
				if a.Key != "time" {
					return a
				}
				newTime := a.Value.Time().In(cn)
				return slog.Time(a.Key, newTime)
			},
		}),
	)
}

func TextHandler(w io.Writer, level slog.Level) slog.Handler {
	cn, _ := time.LoadLocation("Asia/Shanghai")
	return sf.NewFormatterHandler(sf.TimeFormatter(time.DateTime, time.UTC))(
		slog.NewTextHandler(w, &slog.HandlerOptions{
			Level: level,
			ReplaceAttr: func(_ []string, a slog.Attr) slog.Attr {
				if a.Key != "time" {
					return a
				}
				newTime := a.Value.Time().In(cn)
				return slog.Time(a.Key, newTime)
			},
		}),
	)
}

func FileWriter(outputFile string) io.WriteCloser {
	return &lumberjack.Logger{
		Filename:   outputFile,
		MaxSize:    10,
		MaxAge:     1,
		MaxBackups: 1,
		LocalTime:  true,
	}
}

func Init(outputFile, level string) {
	l := slog.LevelInfo
	switch level {
	case "debug":
		l = slog.LevelDebug
	case "warn":
		l = slog.LevelWarn
	case "error":
		l = slog.LevelError
	}

	if outputFile != "" {
		SetLogger(slog.New(TextHandler(FileWriter(outputFile), l)))
	} else {
		SetLogger(slog.New(DefaultHandler(l)))
	}
}
