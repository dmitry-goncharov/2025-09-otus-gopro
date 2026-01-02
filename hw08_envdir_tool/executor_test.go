package main

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestRunCmd(t *testing.T) {
	envs := Environment{
		"BAR":   EnvValue{Value: "bar", NeedRemove: false},
		"EMPTY": EnvValue{Value: "", NeedRemove: false},
		"FOO":   EnvValue{Value: "   foo\nwith new line", NeedRemove: false},
		"HELLO": EnvValue{Value: "\"hello\"", NeedRemove: false},
		"UNSET": EnvValue{Value: "", NeedRemove: true},
	}

	t.Run("run cmd successfully", func(t *testing.T) {
		cmd := []string{"ls"}
		resCode := RunCmd(cmd, envs)

		require.Equal(t, resCode, 0)
	})

	t.Run("run cmd failed", func(t *testing.T) {
		cmd := []string{"./testdata/echo1.sh"}
		resCode := RunCmd(cmd, envs)

		require.Equal(t, resCode, 1)
	})
}
