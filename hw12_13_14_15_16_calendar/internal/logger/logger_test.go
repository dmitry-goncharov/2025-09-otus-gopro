package logger

import (
	"fmt"
	"io"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestLogger(t *testing.T) {
	type test struct {
		logLevel    string
		msgLevel    string
		shouldPrint bool
	}
	tests := []test{
		{
			logLevel:    LevelDebug,
			msgLevel:    LevelDebug,
			shouldPrint: true,
		},
		{
			logLevel:    LevelDebug,
			msgLevel:    LevelInfo,
			shouldPrint: true,
		},
		{
			logLevel:    LevelDebug,
			msgLevel:    LevelWarn,
			shouldPrint: true,
		},
		{
			logLevel:    LevelDebug,
			msgLevel:    LevelError,
			shouldPrint: true,
		},
		{
			logLevel:    LevelInfo,
			msgLevel:    LevelDebug,
			shouldPrint: false,
		},
		{
			logLevel:    LevelInfo,
			msgLevel:    LevelInfo,
			shouldPrint: true,
		},
		{
			logLevel:    LevelInfo,
			msgLevel:    LevelWarn,
			shouldPrint: true,
		},
		{
			logLevel:    LevelInfo,
			msgLevel:    LevelError,
			shouldPrint: true,
		},
		{
			logLevel:    LevelWarn,
			msgLevel:    LevelDebug,
			shouldPrint: false,
		},
		{
			logLevel:    LevelWarn,
			msgLevel:    LevelInfo,
			shouldPrint: false,
		},
		{
			logLevel:    LevelWarn,
			msgLevel:    LevelWarn,
			shouldPrint: true,
		},
		{
			logLevel:    LevelWarn,
			msgLevel:    LevelError,
			shouldPrint: true,
		},
		{
			logLevel:    LevelError,
			msgLevel:    LevelDebug,
			shouldPrint: false,
		},
		{
			logLevel:    LevelError,
			msgLevel:    LevelInfo,
			shouldPrint: false,
		},
		{
			logLevel:    LevelError,
			msgLevel:    LevelWarn,
			shouldPrint: false,
		},
		{
			logLevel:    LevelError,
			msgLevel:    LevelError,
			shouldPrint: true,
		},
	}
	for _, tc := range tests {
		var msg string
		if tc.shouldPrint {
			msg = fmt.Sprintf("%s logger should print %s", tc.logLevel, tc.msgLevel)
		} else {
			msg = fmt.Sprintf("%s logger should not print %s", tc.logLevel, tc.msgLevel)
		}
		t.Run(msg, func(t *testing.T) {
			rescueStdout := os.Stdout
			r, w, _ := os.Pipe()
			os.Stdout = w

			log, err := New(tc.logLevel, false)

			require.NoError(t, err)

			switch tc.msgLevel {
			case LevelDebug:
				log.Debug(msg)
			case LevelInfo:
				log.Info(msg)
			case LevelWarn:
				log.Warn(msg)
			case LevelError:
				log.Error(msg)
			}

			w.Close()
			out, _ := io.ReadAll(r)
			os.Stdout = rescueStdout

			if tc.shouldPrint {
				require.Contains(t, string(out), msg)
			} else {
				require.NotContains(t, string(out), msg)
			}
		})
	}

	t.Run("invalid log level", func(t *testing.T) {
		log, err := New("some", false)

		require.Error(t, err)
		require.Nil(t, log)
	})
}
