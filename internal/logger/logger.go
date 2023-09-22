package logger

import "go.uber.org/zap"

var _ Logger = (*logger)(nil)
var _ Logger = (*noLog)(nil)

type Logger interface {
	Debug(args ...interface{})
	Debugf(template string, args ...interface{})
	Info(args ...interface{})
	Infof(template string, args ...interface{})
	Error(args ...interface{})
	Errorf(template string, args ...interface{})
	Fatal(args ...interface{})
	Fatalf(template string, args ...interface{})
	NoError(err error, msg ...string)
}

func NewLoggerFromZap(zapLogger *zap.SugaredLogger) Logger {
	return &logger{zapLogger}
}

var NoLog = noLog{}

type noLog struct{}

func (*noLog) Debug(args ...interface{})                   {}
func (*noLog) Debugf(template string, args ...interface{}) {}
func (*noLog) Info(args ...interface{})                    {}
func (*noLog) Infof(template string, args ...interface{})  {}
func (*noLog) Error(args ...interface{})                   {}
func (*noLog) Errorf(template string, args ...interface{}) {}
func (*noLog) Fatal(args ...interface{})                   {}
func (*noLog) Fatalf(template string, args ...interface{}) {}
func (*noLog) NoError(err error, msg ...string) {
	if err != nil {
		panic(err)
	}
}

type logger struct {
	*zap.SugaredLogger
}

func (l *logger) NoError(err error, msg ...string) {
	if err != nil {
		if len(msg) > 0 {
			l.Fatalf("\n\nERROR (%s): %v\n\n", err, msg)
		}
		l.Fatalf("\n\nERROR: %v\n\n", err)
	}
}
