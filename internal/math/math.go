package mathcmd

import (
	"bufio"
	"fmt"
	"io"
	"math"
	"strconv"
	"strings"
)

const scaleMax = 15

func Run(args []string, stdin io.Reader, stdout, stderr io.Writer) int {
	// Parse flags manually — math args look like expressions, not typical getopt flags.
	scale := -1 // -1 = auto
	base := 10
	var exprArgs []string

	i := 0
	for i < len(args) {
		arg := args[i]
		if arg == "--" {
			exprArgs = append(exprArgs, args[i+1:]...)
			break
		}
		if arg == "--help" || arg == "-h" {
			fmt.Fprint(stdout, `Usage: math [-s SCALE] [-b BASE] EXPRESSION ...

  Evaluate arithmetic expressions.

Options:
  -s, --scale=N    Decimal places (0-15, or 'max' for full precision)
  -b, --base=N     Output base (2-16, or 'hex'/'octal', default 10)
  -h, --help       Show this help message

Constants: pi, e, tau, inf
Functions: sin cos tan asin acos atan atan2 sinh cosh tanh
           sqrt exp ln log log2 log10 abs floor ceil round pow max min
           bitand bitor bitxor
`)
			return 0
		}

		// --scale=N or -sN or -s N
		if arg == "--scale" || arg == "-s" {
			i++
			if i >= len(args) {
				fmt.Fprintln(stderr, "math: missing scale value")
				return 1
			}
			s, ok := parseScale(args[i], stderr)
			if !ok {
				return 1
			}
			scale = s
			i++
			continue
		}
		if strings.HasPrefix(arg, "--scale=") {
			s, ok := parseScale(arg[8:], stderr)
			if !ok {
				return 1
			}
			scale = s
			i++
			continue
		}
		if strings.HasPrefix(arg, "-s") && len(arg) > 2 {
			s, ok := parseScale(arg[2:], stderr)
			if !ok {
				return 1
			}
			scale = s
			i++
			continue
		}

		// --base=N or -bN or -b N or --base N
		if arg == "--base" || arg == "-b" {
			i++
			if i >= len(args) {
				fmt.Fprintln(stderr, "math: missing base value")
				return 1
			}
			b, ok := parseBase(args[i], stderr)
			if !ok {
				return 1
			}
			base = b
			i++
			continue
		}
		if strings.HasPrefix(arg, "--base=") {
			b, ok := parseBase(arg[7:], stderr)
			if !ok {
				return 1
			}
			base = b
			i++
			continue
		}
		if strings.HasPrefix(arg, "-b") && len(arg) > 2 {
			b, ok := parseBase(arg[2:], stderr)
			if !ok {
				return 1
			}
			base = b
			i++
			continue
		}

		exprArgs = append(exprArgs, arg)
		i++
	}

	if base != 10 && scale > 0 {
		fmt.Fprintln(stderr, "math: invalid option combination, non-zero scale value only valid for base 10")
		return 1
	}

	// Read from stdin if available and no args given (fish behavior: stdin takes precedence)
	if len(exprArgs) == 0 {
		scanner := bufio.NewScanner(stdin)
		if scanner.Scan() {
			line := strings.TrimSpace(scanner.Text())
			if line != "" {
				exprArgs = []string{line}
			}
		}
	}

	if len(exprArgs) == 0 {
		fmt.Fprintln(stderr, "math: expected >= 1 arguments; got 0")
		return 1
	}

	expr := strings.Join(exprArgs, " ")
	val, err := Evaluate(expr)
	if err != nil {
		fmt.Fprintf(stderr, "math: Error: %v\n", err)
		return 1
	}

	fmt.Fprintln(stdout, formatResult(val, scale, base))
	return 0
}

func parseScale(s string, stderr io.Writer) (int, bool) {
	if s == "max" {
		return scaleMax, true
	}
	n, err := strconv.Atoi(s)
	if err != nil || n < 0 || n > 15 {
		if err == nil {
			fmt.Fprintf(stderr, "math: %d: invalid scale\n", n)
		} else {
			fmt.Fprintf(stderr, "math: %s: invalid scale\n", s)
		}
		return 0, false
	}
	return n, true
}

func parseBase(s string, stderr io.Writer) (int, bool) {
	switch strings.ToLower(s) {
	case "hex":
		return 16, true
	case "octal":
		return 8, true
	case "binary":
		return 2, true
	}
	n, err := strconv.Atoi(s)
	if err != nil || n < 2 || n > 16 {
		fmt.Fprintf(stderr, "math: %s: invalid base value\n", s)
		return 0, false
	}
	return n, true
}

func formatResult(val float64, scale, base int) string {
	if math.IsInf(val, 1) {
		return "inf"
	}
	if math.IsInf(val, -1) {
		return "-inf"
	}
	if math.IsNaN(val) {
		return "nan"
	}

	if base != 10 {
		n := int64(val)
		neg := n < 0
		if neg {
			n = -n
		}
		digits := strconv.FormatInt(n, base)
		var prefix string
		switch base {
		case 16:
			prefix = "0x"
		case 8:
			if digits != "0" {
				prefix = "0"
			}
		}
		if neg {
			return "-" + prefix + digits
		}
		return prefix + digits
	}

	if scale == 0 {
		return strconv.FormatInt(int64(val), 10)
	}

	if scale == scaleMax {
		return strconv.FormatFloat(val, 'f', 15, 64)
	}

	if scale >= 0 {
		return strconv.FormatFloat(val, 'f', scale, 64)
	}

	// Auto: whole numbers as integers, floats at 6 decimal places (trimmed)
	if val == math.Trunc(val) && !math.IsInf(val, 0) && val >= -1e15 && val <= 1e15 {
		return strconv.FormatInt(int64(val), 10)
	}
	s := strconv.FormatFloat(val, 'f', 6, 64)
	s = strings.TrimRight(s, "0")
	s = strings.TrimRight(s, ".")
	return s
}
