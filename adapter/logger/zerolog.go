package logger

import (
	"context"
	"fmt"
	"os"
	"reflect"
	"strings"
	"time"

	"github.com/logrusorgru/aurora"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/pkgerrors"

	"narou/infrastructure/log"
)

type levelLogger struct {
	logger zerolog.Logger
}

func NewLogger(ctx context.Context) log.Logger {
	output := zerolog.ConsoleWriter{
		Out:             os.Stdout,
		NoColor:         false,
		TimeFormat:      time.RFC3339,
		PartsOrder:      nil,
		FormatTimestamp: nil,
		FormatLevel: func(i interface{}) string {
			switch i.(string) {
			case "error":
				return aurora.Red(strings.ToUpper(fmt.Sprintf("| %-6s|", i))).String()
			case "debug":
				return aurora.Yellow(strings.ToUpper(fmt.Sprintf("| %-6s|", i))).String()
			case "info":
				return aurora.Green(strings.ToUpper(fmt.Sprintf("| %-6s|", i))).String()
			default:
				return aurora.Green(strings.ToUpper(fmt.Sprintf("| %-6s|", i))).String()
			}
		},
		FormatMessage: func(i interface{}) string {
			return fmt.Sprintf("*** %s ****", i)
		},
		FormatFieldValue: func(i interface{}) string {
			return strings.ToUpper(fmt.Sprintf("%s", i))
		},
		FormatErrFieldName: func(i interface{}) string {
			return fmt.Sprintf("%s:", i)
		},
		FormatErrFieldValue: nil,
	}

	lg := &levelLogger{
		logger: zerolog.New(output).With().Logger(),
	}
	lg.logger.WithContext(ctx)
	return lg
}

func (l *levelLogger) Debug(msg string, keyvals ...interface{}) {
	l.print(l.logger.Debug, msg, keyvals...)
}

func (l *levelLogger) Info(msg string, keyvals ...interface{}) {
	l.print(l.logger.Info, msg, keyvals...)
}

func (l *levelLogger) Warn(msg string, keyvals ...interface{}) {
	l.print(l.logger.Warn, msg, keyvals...)
}

func (l *levelLogger) Error(msg string, keyvals ...interface{}) {
	l.print(l.logger.Error, msg, keyvals...)
}

func (l *levelLogger) Fatal(msg string, keyvals ...interface{}) {
	l.print(l.logger.Fatal, msg, keyvals...)
}

func (l *levelLogger) Panic(msg string, keyvals ...interface{}) {
	l.print(l.logger.Panic, msg, keyvals...)
}

func (l *levelLogger) Log(keyvals ...interface{}) {
	if len(keyvals) == 1 {
		l.print(l.logger.Log, fmt.Sprint(keyvals[0]))
	} else {
		l.print(l.logger.Log, "", keyvals...)
	}
}

func (l *levelLogger) print(lvl func() *zerolog.Event, msg string, keyvals ...interface{}) {
	// when log.Info("test"), keyValues is [[]]
	if !(len(keyvals) == 1 && reflect.ValueOf(keyvals[0]).IsNil()) &&
		len(keyvals)%2 != 0 {
		panic("illegal format")
	}

	if len(keyvals) == 1 && reflect.ValueOf(keyvals[0]).IsNil() {
		keyvals = nil
	}

	event := lvl()
	event = event.Timestamp()

	for i := 0; i < len(keyvals); i += 2 {
		k := convertKey(keyvals[i])
		v := keyvals[i+1]
		event = apply(event, k, v)
	}

	if msg == "" {
		event.Send()
	} else {
		event.Msg(msg)
	}
}

func convertKey(k interface{}) string {
	switch x := k.(type) {
	case string:
		return x
	case fmt.Stringer:
		return x.String()
	default:
		return fmt.Sprint(k)
	}
}

func apply(e *zerolog.Event, k string, v interface{}) *zerolog.Event {
	switch x := v.(type) {
	case int:
		e = e.Int(k, x)
	case []int:
		e = e.Ints(k, x)
	case int64:
		e = e.Int64(k, x)
	case string:
		e = e.Str(k, x)
	case []string:
		e = e.Strs(k, x)
	case time.Time:
		e = e.Time(k, x)
	case []time.Time:
		e = e.Times(k, x)
	case time.Duration:
		e = e.Dur(k, x)
	case []time.Duration:
		e = e.Durs(k, x)
	case float32:
		e = e.Float32(k, x)
	case []float32:
		e = e.Floats32(k, x)
	case float64:
		e = e.Float64(k, x)
	case []float64:
		e = e.Floats64(k, x)
	case error:
		e = e.AnErr(k, x)
		if trace := pkgerrors.MarshalStack(x); trace != nil {
			e = e.Interface("stack trace", trace)
		}
	case []error:
		e = e.Errs(k, x)
	}

	return e
}
