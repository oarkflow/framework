package log

import (
	"context"
	"errors"
	log2 "github.com/sujit-baniya/framework/contracts/log"

	"github.com/sujit-baniya/framework/facades"
	"github.com/sujit-baniya/framework/log/logger"

	"github.com/gookit/color"
	"github.com/sirupsen/logrus"
)

type Logrus struct {
	instance *logrus.Logger
	log2.Writer
}

func NewLogrus(logger *logrus.Logger, writer log2.Writer) log2.Log {
	return &Logrus{
		instance: logger,
		Writer:   writer,
	}
}

func logrusInstance() *logrus.Logger {
	instance := logrus.New()
	instance.SetLevel(logrus.DebugLevel)

	if facades.Config != nil {
		logging := facades.Config.GetString("logging.default")
		if logging != "" {
			if err := registerHook(instance, logging); err != nil {
				color.Redln("Init facades.Log error: " + err.Error())

				return nil
			}
		}
	}

	return instance
}

func (r *Logrus) WithContext(ctx context.Context) log2.Writer {
	switch r.Writer.(type) {
	case *Writer:
		return NewWriter(r.instance.WithContext(ctx))
	default:
		return r.Writer
	}
}

type Writer struct {
	instance *logrus.Entry
}

func NewWriter(instance *logrus.Entry) log2.Writer {
	return &Writer{instance: instance}
}

func (r *Writer) Debug(args ...interface{}) {
	r.instance.Debug(args...)
}

func (r *Writer) Debugf(format string, args ...interface{}) {
	r.instance.Debugf(format, args...)
}

func (r *Writer) Info(args ...interface{}) {
	r.instance.Info(args...)
}

func (r *Writer) Infof(format string, args ...interface{}) {
	r.instance.Infof(format, args...)
}

func (r *Writer) Warning(args ...interface{}) {
	r.instance.Warning(args...)
}

func (r *Writer) Warningf(format string, args ...interface{}) {
	r.instance.Warningf(format, args...)
}

func (r *Writer) Error(args ...interface{}) {
	r.instance.Error(args...)
}

func (r *Writer) Errorf(format string, args ...interface{}) {
	r.instance.Errorf(format, args...)
}

func (r *Writer) Fatal(args ...interface{}) {
	r.instance.Fatal(args...)
}

func (r *Writer) Fatalf(format string, args ...interface{}) {
	r.instance.Fatalf(format, args...)
}

func (r *Writer) Panic(args ...interface{}) {
	r.instance.Panic(args...)
}

func (r *Writer) Panicf(format string, args ...interface{}) {
	r.instance.Panicf(format, args...)
}

func registerHook(instance *logrus.Logger, channel string) error {
	channelPath := "logging.channels." + channel
	driver := facades.Config.GetString(channelPath + ".driver")

	var hook logrus.Hook
	var err error
	switch driver {
	case log2.StackDriver:
		for _, stackChannel := range facades.Config.Get(channelPath + ".channels").([]string) {
			if stackChannel == channel {
				return errors.New("stack drive can't include self channel")
			}

			if err := registerHook(instance, stackChannel); err != nil {
				return err
			}
		}

		return nil
	case log2.SingleDriver:
		logLogger := &logger.Single{}
		hook, err = logLogger.Handle(channelPath)
		if err != nil {
			return err
		}
	case log2.DailyDriver:
		logLogger := &logger.Daily{}
		hook, err = logLogger.Handle(channelPath)
		if err != nil {
			return err
		}
	case log2.CustomDriver:
		logLogger := facades.Config.Get(channelPath + ".via").(log2.Logger)
		logHook, err := logLogger.Handle(channelPath)
		if err != nil {
			return err
		}

		hook = &Hook{logHook}
	default:
		return errors.New("Error logging channel: " + channel)
	}

	instance.AddHook(hook)

	return nil
}

type Hook struct {
	instance log2.Hook
}

func (h *Hook) Levels() []logrus.Level {
	levels := h.instance.Levels()
	var logrusLevels []logrus.Level
	for _, item := range levels {
		logrusLevels = append(logrusLevels, logrus.Level(item))
	}

	return logrusLevels
}

func (h *Hook) Fire(entry *logrus.Entry) error {
	return h.instance.Fire(&Entry{
		ctx:     entry.Context,
		level:   log2.Level(entry.Level),
		time:    entry.Time,
		message: entry.Message,
	})
}
