package sk

import (
	"fmt"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
	"os"
	"time"
)

type zapBuilder struct {
	hook lumberjack.Logger

	fileEnable bool
	fileLevel zapcore.LevelEnabler

	consoleEnable bool
	consoleLevel zapcore.LevelEnabler
}

func NewZapBuilder() *zapBuilder {
	return &zapBuilder{
		hook:lumberjack.Logger{
			Filename:"all.log",
			MaxSize:    10,
			MaxBackups: 3,
			MaxAge:     7,
			Compress:   true,
			LocalTime:false,
		},
		fileEnable:false,
		fileLevel:zapcore.WarnLevel,
		consoleEnable:false,
		consoleLevel:zapcore.InfoLevel,
	}
}

func (opts *zapBuilder) EnableFile() *zapBuilder {
	opts.fileEnable = true
	return opts
}

func (opts *zapBuilder) EnableConsole() *zapBuilder {
	opts.consoleEnable = true
	return opts
}

func (opts *zapBuilder) Filename(filename string) *zapBuilder {
	opts.hook.Filename = filename
	return opts
}

func (opts *zapBuilder) FileMegaBytes(maxSize int) *zapBuilder {
	opts.hook.MaxSize = maxSize
	return opts
}

func (opts *zapBuilder) FileMaxBackups(maxBackups int) *zapBuilder {
	opts.hook.MaxBackups = maxBackups
	return opts
}

func (opts *zapBuilder) FileMaxAge(maxAge int) *zapBuilder {
	opts.hook.MaxAge = maxAge
	return opts
}

func (opts *zapBuilder) FileCompress(compress bool) *zapBuilder {
	opts.hook.Compress = compress
	return opts
}

func (opts *zapBuilder) FileLevel(level string) *zapBuilder {
	opts.fileLevel = getLevel(level)
	return opts
}

func (opts *zapBuilder) ConsoleLevel(level string) *zapBuilder {
	opts.consoleLevel = getLevel(level)
	return opts
}

func getLevel(s string) zapcore.LevelEnabler {
	l := zapcore.Level(zapcore.DebugLevel)
	_ = l.Set(s)
	return l
}

func (opts *zapBuilder) Build() *zap.Logger {
	encoderConfig := zapcore.EncoderConfig{
		TimeKey:       "T",
		LevelKey:      "L",
		NameKey:       "K",
		CallerKey:     "N",
		MessageKey:    "M",
		StacktraceKey: "S",
		LineEnding:    zapcore.DefaultLineEnding,
		EncodeLevel:   zapcore.LowercaseLevelEncoder, // 小写编码器
		EncodeTime: func(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
			enc.AppendString(t.Format("01-02 15:04:05"))
		},
		EncodeDuration: zapcore.SecondsDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
		EncodeName:     zapcore.FullNameEncoder,
	}

	var allCore []zapcore.Core
	if opts.fileEnable {
		fmt.Printf("%+v\n", opts.hook)
		fileW := zapcore.AddSync(&opts.hook)
		core := zapcore.NewCore(zapcore.NewConsoleEncoder(encoderConfig), fileW, opts.fileLevel)
		allCore = append(allCore, core)
	}

	if opts.consoleEnable {
		consoleW := zapcore.Lock(os.Stdout)
		core := zapcore.NewCore(zapcore.NewConsoleEncoder(encoderConfig), consoleW, opts.consoleLevel)
		allCore = append(allCore, core)
	}

	core := zapcore.NewTee(allCore...)
	return zap.New(core, zap.AddCaller())
}

var (
	loggers = make(map[string]*zap.Logger)
)

func RegisterLogger(k string, logger *zap.Logger) {
	loggers[k] = logger
}

func GetLogger(k string) *zap.Logger {
	return loggers[k]
}