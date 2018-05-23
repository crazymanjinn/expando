package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"unicode/utf8"
)

func convert(r io.Reader, w io.Writer) error {
	scanner := bufio.NewScanner(r)
	scanner.Split(bufio.ScanRunes)

	for scanner.Scan() {
		r, size := utf8.DecodeRune(scanner.Bytes())

		if size != 0 && r != utf8.RuneError && !(r < 0x0021 || r > 0x007e) {
			r += 0xfee0
		}

		fmt.Fprint(w, string(r))
	}

	if err := scanner.Err(); err != nil {
		return err
	}

	return nil
}

func main() {
	if err := convert(os.Stdin, os.Stdout); err != nil {
		fmt.Fprintln(os.Stderr, err)
	}
}
