// gtools - a collection of grammar manipulation tools
// Copyright (c) 2023 Michael D Henderson. All rights reserved.

package gtools

import (
	"errors"
	"fmt"
	"io"
	"os"
)

var (
	stdin  *os.File = os.Stdin
	stdout *os.File = os.Stdout
	stderr *os.File = os.Stderr
)

const (
	EOF byte = 255
)

func fprintf(fp *os.File, format string, args ...any) {
	_, _ = fmt.Fprintf(fp, format, args...)
}

func fputs(message string, fp *os.File) {
	_, _ = fmt.Fprint(fp, message)
}

func getchar() byte {
	if stdin == nil {
		return EOF
	}
	var buffer [1]byte
	n, err := stdin.Read(buffer[:])
	if err != nil {
		if n == 0 && errors.Is(err, io.EOF) {
			return EOF
		}
		panic(err)
	}
	return buffer[0]
}

func putchar(ch byte) {
	_, _ = fmt.Fprintf(stdout, "%c", ch)
}

func SetStdin(input string) error {
	if stdin != nil {
		_ = stdin.Close()
	}

	fp, err := os.Open(input)
	if err != nil {
		return err
	}
	stdin = fp

	return nil
}

func tocstring(s string) []byte {
	return append([]byte(s), 0)
}
