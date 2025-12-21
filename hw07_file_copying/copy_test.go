package main

import (
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestCopy(t *testing.T) {
	tests := []struct {
		name   string
		from   string
		to     string
		offset int64
		limit  int64
		err    error
		exp    string
	}{
		{
			name:   "negative offset",
			from:   "./testdata/input.txt",
			to:     "./out.txt",
			offset: -10,
			limit:  0,
			err:    ErrNegativeOffset,
			exp:    "./testdata/out_offset0_limit0.txt",
		},
		{
			name:   "negative limit",
			from:   "./testdata/input.txt",
			to:     "./out.txt",
			offset: 0,
			limit:  -10,
			err:    ErrNegativeLimit,
			exp:    "./testdata/out_offset0_limit0.txt",
		},
		{
			name:   "invalid source path",
			from:   "./testdata/input1.txt",
			to:     "./testdata/input.txt",
			offset: 0,
			limit:  0,
			err:    os.ErrNotExist,
			exp:    "./testdata/out_offset0_limit0.txt",
		},
		{
			name:   "invalid destination path",
			from:   "./testdata/input.txt",
			to:     "./testdata/input.txt",
			offset: 0,
			limit:  0,
			err:    ErrInvalidDestPath,
			exp:    "./testdata/out_offset0_limit0.txt",
		},
		{
			name:   "unsupport dir",
			from:   "./testdata",
			to:     "./out_offset0_limit0.txt",
			offset: 0,
			limit:  0,
			err:    ErrUnsupportedFile,
			exp:    "./testdata/out_offset0_limit0.txt",
		},
		{
			name:   "undefined file size",
			from:   "/dev/urandom",
			to:     "./out_offset0_limit0.txt",
			offset: 0,
			limit:  0,
			err:    ErrUnsupportedFile,
			exp:    "./testdata/out_offset0_limit0.txt",
		},
		{
			name:   "offset exceeds file size",
			from:   "./testdata/input.txt",
			to:     "./out_offset0_limit0.txt",
			offset: 1000000000,
			limit:  0,
			err:    ErrOffsetExceedsFileSize,
			exp:    "./testdata/out_offset0_limit0.txt",
		},
		{
			name:   "offset 0 limit 0",
			from:   "./testdata/input.txt",
			to:     "./out_offset0_limit0.txt",
			offset: 0,
			limit:  0,
			err:    nil,
			exp:    "./testdata/out_offset0_limit0.txt",
		},
		{
			name:   "offset 0 limit 10",
			from:   "./testdata/input.txt",
			to:     "./out_offset0_limit10.txt",
			offset: 0,
			limit:  10,
			err:    nil,
			exp:    "./testdata/out_offset0_limit10.txt",
		},
		{
			name:   "offset 0 limit 1000",
			from:   "./testdata/input.txt",
			to:     "./out_offset0_limit1000.txt",
			offset: 0,
			limit:  1000,
			err:    nil,
			exp:    "./testdata/out_offset0_limit1000.txt",
		},
		{
			name:   "offset 0 limit 10000",
			from:   "./testdata/input.txt",
			to:     "./out_offset0_limit10000.txt",
			offset: 0,
			limit:  10000,
			err:    nil,
			exp:    "./testdata/out_offset0_limit10000.txt",
		},
		{
			name:   "offset 100 limit 1000",
			from:   "./testdata/input.txt",
			to:     "./out_offset100_limit1000.txt",
			offset: 100,
			limit:  1000,
			err:    nil,
			exp:    "./testdata/out_offset100_limit1000.txt",
		},
		{
			name:   "offset 6000 limit 1000",
			from:   "./testdata/input.txt",
			to:     "./out_offset6000_limit1000.txt",
			offset: 6000,
			limit:  1000,
			err:    nil,
			exp:    "./testdata/out_offset6000_limit1000.txt",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			err := Copy(test.from, test.to, test.offset, test.limit)

			if test.err != nil {
				require.ErrorIs(t, err, test.err)
			} else {
				require.NoError(t, err)

				fiExp, _ := os.Stat(test.exp)
				fiAct, _ := os.Stat(test.to)
				require.Equal(t, fiExp.Size(), fiAct.Size())

				os.Remove(test.to)
			}
		})
	}
}
