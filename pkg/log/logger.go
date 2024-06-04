package log

import (
	"github.com/rs/zerolog"
	"os"
	"runtime/debug"
	"time"
)

var Logger zerolog.Logger

func init() {
	buildInfo, _ := debug.ReadBuildInfo()

	Logger = zerolog.New(zerolog.ConsoleWriter{Out: os.Stderr, TimeFormat: time.DateTime}).
		Level(zerolog.TraceLevel).
		With().
		Timestamp().
		Caller().
		Int("pid", os.Getpid()).
		Str("go_version", buildInfo.GoVersion).
		Logger()

}
