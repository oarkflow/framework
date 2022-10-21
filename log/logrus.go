package log

import (
	"errors"
	"fmt"

	"github.com/sujit-baniya/framework/contracts/log"
	"github.com/sujit-baniya/framework/facades"
	"github.com/sujit-baniya/framework/log/logger"

	"github.com/gookit/color"
	"github.com/sirupsen/logrus"
)

type Logrus struct {
	Instance *logrus.Logger
	Test     bool
}

func NewLogrus() log.Log {
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

	return &Logrus{instance, false}
}

func (r *Logrus) Testing(is bool) log.Log {
	r.Test = is

	return r
}

func (r *Logrus) Debug(args ...any) {
	if r.Test {
		fmt.Print("Debug: ")
		fmt.Println(args...)
		return
	}

	r.Instance.Debug(args...)
}

func (r *Logrus) Debugf(format string, args ...any) {
	if r.Test {
		fmt.Print("Debugf: ")
		fmt.Printf(format+"\n", args...)
		return
	}

	r.Instance.Debugf(format, args...)
}

func (r *Logrus) Info(args ...any) {
	if r.Test {
		fmt.Print("Info: ")
		fmt.Println(args...)
		return
	}

	r.Instance.Info(args...)
}

func (r *Logrus) Infof(format string, args ...any) {
	if r.Test {
		fmt.Print("Infof: ")
		fmt.Printf(format+"\n", args...)
		return
	}

	r.Instance.Infof(format, args...)
}

func (r *Logrus) Warning(args ...any) {
	if r.Test {
		fmt.Print("Warningf: ")
		fmt.Println(args...)
		return
	}

	r.Instance.Warning(args...)
}

func (r *Logrus) Warningf(format string, args ...any) {
	if r.Test {
		fmt.Print("Warningf: ")
		fmt.Printf(format+"\n", args...)
		return
	}

	r.Instance.Warningf(format, args...)
}

func (r *Logrus) Error(args ...any) {
	if r.Test {
		fmt.Print("Error: ")
		fmt.Println(args...)
		return
	}

	r.Instance.Error(args...)
}

func (r *Logrus) Errorf(format string, args ...any) {
	if r.Test {
		fmt.Print("Errorf: ")
		fmt.Printf(format+"\n", args...)
		return
	}

	r.Instance.Errorf(format, args...)
}

func (r *Logrus) Fatal(args ...any) {
	if r.Test {
		fmt.Print("Error: ")
		fmt.Println(args...)
		return
	}

	r.Instance.Fatal(args...)
}

func (r *Logrus) Fatalf(format string, args ...any) {
	if r.Test {
		fmt.Print("Errorf: ")
		fmt.Printf(format+"\n", args...)
		return
	}

	r.Instance.Fatalf(format, args...)
}

func (r *Logrus) Panic(args ...any) {
	if r.Test {
		fmt.Print("Panic: ")
		fmt.Println(args...)
		return
	}

	r.Instance.Panic(args...)
}

func (r *Logrus) Panicf(format string, args ...any) {
	if r.Test {
		fmt.Print("Panicf: ")
		fmt.Printf(format+"\n", args...)
		return
	}

	r.Instance.Panicf(format, args...)
}

func registerHook(instance *logrus.Logger, channel string) error {
	driver := facades.Config.GetString("logging.channels." + channel + ".driver")
	channelPath := "logging.channels." + channel

	var hook logrus.Hook
	var err error
	switch driver {
	case log.StackDriver:
		for _, stackChannel := range facades.Config.Get("logging.channels." + channel + ".channels").([]string) {
			if stackChannel == channel {
				return errors.New("stack drive can't include self channel")
			}

			if err := registerHook(instance, stackChannel); err != nil {
				return err
			}
		}

		return nil
	case log.SingleDriver:
		logLogger := &logger.Single{}
		hook, err = logLogger.Handle(channelPath)
		if err != nil {
			return err
		}
	case log.DailyDriver:
		logLogger := &logger.Daily{}
		hook, err = logLogger.Handle(channelPath)
		if err != nil {
			return err
		}
	case log.CustomDriver:
		logLogger := facades.Config.Get("logging.channels." + channel + ".via").(log.Logger)
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
	instance log.Hook
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
		level:   log.Level(entry.Level),
		time:    entry.Time,
		message: entry.Message,
	})
}
