package main

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestReadDir(t *testing.T) {
	expEnvs := Environment{
		"BAR":   EnvValue{Value: "bar", NeedRemove: false},
		"EMPTY": EnvValue{Value: "", NeedRemove: false},
		"FOO":   EnvValue{Value: "   foo\nwith new line", NeedRemove: false},
		"HELLO": EnvValue{Value: "\"hello\"", NeedRemove: false},
		"UNSET": EnvValue{Value: "", NeedRemove: true},
	}

	t.Run("read envs", func(t *testing.T) {
		envs, err := ReadDir("./testdata/env")

		require.Equal(t, expEnvs, envs)
		require.NoError(t, err)
	})

	t.Run("read envs from not existing dir", func(t *testing.T) {
		envs, err := ReadDir("./testdata/env1")

		require.Nil(t, envs)
		require.Error(t, err)
	})
}
