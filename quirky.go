package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"strings"

	qrcode "github.com/skip2/go-qrcode"
)

var doubleSize = flag.Bool("d", false, "Make QR code double size")
var inverted = flag.Bool("i", false, "Invert colors")

func normal() {
	if *inverted {
		fmt.Print("\x1b[0m")
	} else {
		fmt.Print("\x1b[7m")
	}
}

func invert() {
	if *inverted {
		fmt.Print("\x1b[7m")
	} else {
		fmt.Print("\x1b[0m")
	}
}

func printCode(bitmap [][]bool) {
	width := len(bitmap)
	height := len(bitmap[0])

	for y := 0; y < width; y += 2 {
		lastInverted := false
		if !*inverted {
			lastInverted = true
			normal()
		}

		for x := 0; x < height; x++ {
			upper := bitmap[y][x]

			var lower bool
			if y+1 < width {
				lower = bitmap[y+1][x]
			} else {
				lower = false
			}

			if upper == lower {
				if upper && !lastInverted {
					invert()
					lastInverted = true
				} else if !upper && lastInverted {
					normal()
					lastInverted = false
				}
				fmt.Print(" ")
			} else {
				if upper && !lastInverted {
					invert()
					lastInverted = true
				} else if !upper && lastInverted {
					normal()
					lastInverted = false
				}
				fmt.Print("▄")
			}
		}
		fmt.Println("\x1b[0m")
	}
}

func printCodeDouble(bitmap [][]bool) {
	width := len(bitmap)
	height := len(bitmap[0])

	for y := 0; y < width; y++ {
		lastInverted := false
		if !*inverted {
			lastInverted = true
			normal()
		}

		for x := 0; x < height; x++ {
			if bitmap[y][x] && !lastInverted {
				invert()
				lastInverted = true
			} else if !bitmap[y][x] && lastInverted {
				normal()
				lastInverted = false
			}
			fmt.Print("  ")
		}
		fmt.Println("\x1b[0m")
	}
}

func main() {
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: %s [options] [text | -]\n\n", os.Args[0])
		fmt.Fprintln(os.Stderr, "  Specify \"-\" or just nothing to read from stdin.")
		fmt.Fprintln(os.Stderr)
		fmt.Fprintln(os.Stderr, "Options:")
		flag.PrintDefaults()
	}
	flag.Parse()

	data := flag.Args()
	var dataString string
	if len(data) > 1 {
		flag.Usage()
		os.Exit(1)
	} else if len(data) == 0 || data[0] == "-" {
		input := bufio.NewReader(os.Stdin)
		line, _, _ := input.ReadLine()
		dataString = strings.TrimSpace(string(line[:]))
	} else {
		dataString = strings.TrimSpace(data[0])
	}

	if len(dataString) == 0 {
		fmt.Println("No data was provided")
		os.Exit(1)
	}

	qr, err := qrcode.New(dataString, qrcode.Low)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	if *doubleSize {
		printCodeDouble(qr.Bitmap())
	} else {
		printCode(qr.Bitmap())
	}
}
