package main

import (
	"errors"
	"fmt"
	"io"
	"os"

	"github.com/cheggaaa/pb/v3"
)

var (
	ErrNegativeOffset        = errors.New("negative offset")
	ErrNegativeLimit         = errors.New("negative limit")
	ErrInvalidDestPath       = errors.New("invalid destination path")
	ErrUnsupportedFile       = errors.New("unsupported file")
	ErrOffsetExceedsFileSize = errors.New("offset exceeds file size")
)

func Copy(fromPath, toPath string, offset, limit int64) error {
	if offset < 0 {
		return ErrNegativeOffset
	}
	if limit < 0 {
		return ErrNegativeLimit
	}
	if fromPath == toPath {
		return ErrInvalidDestPath
	}

	sourceFile, err := os.Open(fromPath)
	if err != nil {
		return err
	}
	defer func() {
		err := sourceFile.Close()
		if err != nil {
			fmt.Println("Error on close source file", err)
		}
	}()

	fi, err := sourceFile.Stat()
	if err != nil {
		return err
	}

	if fi.IsDir() || !fi.Mode().IsRegular() {
		return ErrUnsupportedFile
	}

	if offset > fi.Size() {
		return ErrOffsetExceedsFileSize
	}

	destFile, err := os.Create(toPath)
	if err != nil {
		return err
	}
	defer func() {
		err := destFile.Close()
		if err != nil {
			fmt.Println("Error on close destination file", err)
		}
	}()

	if offset > 0 {
		_, err = sourceFile.Seek(offset, io.SeekStart)
		if err != nil {
			return err
		}
	}

	var total int64
	if limit == 0 || limit > fi.Size()-offset {
		total = fi.Size() - offset
	} else {
		total = limit
	}

	bar := pb.Full.Start64(total)

	barReader := bar.NewProxyReader(sourceFile)

	_, err = io.CopyN(destFile, barReader, total)

	bar.Finish()

	if err != nil {
		return err
	}
	return nil
}
