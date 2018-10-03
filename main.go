package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"os"
	"unicode/utf8"
)

var quoteStylePtr = flag.Uint("quote", 0, `quotation style:
	0 : ＂＂ (default)
	1 : 「」
	2 : 『』`)
var quoteMap map[uint]map[bool]rune = map[uint]map[bool]rune{
	0: map[bool]rune{false: 0xfee0, true: 0xfee0},
	1: map[bool]rune{false: 0x2fea, true: 0x2feb},
	2: map[bool]rune{false: 0x2fec, true: 0x2fed},
}

func convert(r io.Reader, w io.Writer) error {
	scanner := bufio.NewScanner(r)
	scanner.Split(bufio.ScanRunes)

	var seenQuote bool
	for scanner.Scan() {
		r, size := utf8.DecodeRune(scanner.Bytes())

		if size != 0 && r != utf8.RuneError && !(r < 0x0020 || r > 0x007e) {
			switch r {
			case 0x0020:
				r = 0x3000
			case 0x0022:
				r += quoteMap[*quoteStylePtr][seenQuote]
				seenQuote = !seenQuote
			default:
				r += 0xfee0
			}
		}

		fmt.Fprint(w, string(r))
	}

	if err := scanner.Err(); err != nil {
		return err
	}

	return nil
}

func main() {
	flag.Parse()
	if err := convert(os.Stdin, os.Stdout); err != nil {
		fmt.Fprintln(os.Stderr, err)
	}
}
