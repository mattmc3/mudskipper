package mathcmd

import (
	"fmt"
	"math"
	"strings"
	"unicode"
)

// Evaluate parses and evaluates a math expression string, returning a float64.
func Evaluate(expr string) (float64, error) {
	p := &parser{input: strings.TrimSpace(expr)}
	val, err := p.parseOr()
	if err != nil {
		return 0, err
	}
	p.skipWhitespace()
	if p.pos < len(p.input) {
		return 0, fmt.Errorf("unexpected token: %q", p.input[p.pos:])
	}
	return val, nil
}

type parser struct {
	input string
	pos   int
}

func (p *parser) skipWhitespace() {
	for p.pos < len(p.input) && unicode.IsSpace(rune(p.input[p.pos])) {
		p.pos++
	}
}

func (p *parser) peek() byte {
	p.skipWhitespace()
	if p.pos >= len(p.input) {
		return 0
	}
	return p.input[p.pos]
}

func (p *parser) match(s string) bool {
	p.skipWhitespace()
	if !strings.HasPrefix(p.input[p.pos:], s) {
		return false
	}
	// Don't match prefix of longer operator
	if len(s) == 1 {
		if p.pos+1 < len(p.input) {
			next := p.input[p.pos+1]
			switch s {
			case "<":
				if next == '<' || next == '=' {
					return false
				}
			case ">":
				if next == '>' || next == '=' {
					return false
				}
			case "&":
				if next == '&' {
					return false
				}
			case "|":
				if next == '|' {
					return false
				}
			case "!":
				if next == '=' {
					return false
				}
			case "=":
				if next == '=' {
					return false
				}
			}
		}
	}
	p.pos += len(s)
	return true
}

func (p *parser) match2(s string) bool {
	p.skipWhitespace()
	if strings.HasPrefix(p.input[p.pos:], s) {
		p.pos += len(s)
		return true
	}
	return false
}

// Operator precedence (low to high): || && | & << >> + - * / % x unary ^ (right-assoc)
// Fish math: unary minus binds TIGHTER than ^, so -2^2 = (-2)^2 = 4

func (p *parser) parseOr() (float64, error) {
	left, err := p.parseAnd()
	if err != nil {
		return 0, err
	}
	for p.match2("||") {
		right, err := p.parseAnd()
		if err != nil {
			return 0, err
		}
		if left != 0 || right != 0 {
			left = 1
		} else {
			left = 0
		}
	}
	return left, nil
}

func (p *parser) parseAnd() (float64, error) {
	left, err := p.parseBitOr()
	if err != nil {
		return 0, err
	}
	for p.match2("&&") {
		right, err := p.parseBitOr()
		if err != nil {
			return 0, err
		}
		if left != 0 && right != 0 {
			left = 1
		} else {
			left = 0
		}
	}
	return left, nil
}

func (p *parser) parseBitOr() (float64, error) {
	left, err := p.parseBitAnd()
	if err != nil {
		return 0, err
	}
	for p.match("|") {
		right, err := p.parseBitAnd()
		if err != nil {
			return 0, err
		}
		left = float64(int64(left) | int64(right))
	}
	return left, nil
}

func (p *parser) parseBitAnd() (float64, error) {
	left, err := p.parseEquality()
	if err != nil {
		return 0, err
	}
	for p.match("&") {
		right, err := p.parseEquality()
		if err != nil {
			return 0, err
		}
		left = float64(int64(left) & int64(right))
	}
	return left, nil
}

func (p *parser) parseEquality() (float64, error) {
	left, err := p.parseComparison()
	if err != nil {
		return 0, err
	}
	for {
		if p.match2("==") {
			right, err := p.parseComparison()
			if err != nil {
				return 0, err
			}
			if left == right {
				left = 1
			} else {
				left = 0
			}
		} else if p.match2("!=") {
			right, err := p.parseComparison()
			if err != nil {
				return 0, err
			}
			if left != right {
				left = 1
			} else {
				left = 0
			}
		} else {
			break
		}
	}
	return left, nil
}

func (p *parser) parseComparison() (float64, error) {
	left, err := p.parseShift()
	if err != nil {
		return 0, err
	}
	p.skipWhitespace()
	if p.pos < len(p.input) {
		switch {
		case strings.HasPrefix(p.input[p.pos:], "<="),
			strings.HasPrefix(p.input[p.pos:], ">="),
			p.input[p.pos] == '<',
			p.input[p.pos] == '>':
			return 0, fmt.Errorf("logical operations are not supported, use 'test' instead")
		}
	}
	return left, nil
}

func (p *parser) parseShift() (float64, error) {
	left, err := p.parseAddSub()
	if err != nil {
		return 0, err
	}
	for {
		if p.match2("<<") {
			right, err := p.parseAddSub()
			if err != nil {
				return 0, err
			}
			left = float64(int64(left) << uint(right))
		} else if p.match2(">>") {
			right, err := p.parseAddSub()
			if err != nil {
				return 0, err
			}
			left = float64(int64(left) >> uint(right))
		} else {
			break
		}
	}
	return left, nil
}

func (p *parser) parseAddSub() (float64, error) {
	left, err := p.parseMulDiv()
	if err != nil {
		return 0, err
	}
	for {
		if p.match("+") {
			right, err := p.parseMulDiv()
			if err != nil {
				return 0, err
			}
			left += right
		} else if p.match("-") {
			right, err := p.parseMulDiv()
			if err != nil {
				return 0, err
			}
			left -= right
		} else {
			break
		}
	}
	return left, nil
}

func (p *parser) parseMulDiv() (float64, error) {
	left, err := p.parsePower()
	if err != nil {
		return 0, err
	}
	for {
		if p.match("*") {
			right, err := p.parsePower()
			if err != nil {
				return 0, err
			}
			left *= right
		} else if p.match("/") {
			right, err := p.parsePower()
			if err != nil {
				return 0, err
			}
			if right == 0 {
				return 0, fmt.Errorf("division by zero")
			}
			left /= right
		} else if p.match("%") {
			right, err := p.parsePower()
			if err != nil {
				return 0, err
			}
			if right == 0 {
				return 0, fmt.Errorf("division by zero")
			}
			left = math.Mod(left, right)
		} else if p.matchX() {
			// 'x' as multiply (fish compat: "5 x 4" = 20)
			right, err := p.parsePower()
			if err != nil {
				return 0, err
			}
			left *= right
		} else {
			break
		}
	}
	return left, nil
}

// matchX matches 'x' as a multiplication operator (surrounded by non-alphanumeric).
func (p *parser) matchX() bool {
	p.skipWhitespace()
	if p.pos >= len(p.input) || p.input[p.pos] != 'x' {
		return false
	}
	// Must not be start of identifier like 'xor'
	if p.pos+1 < len(p.input) && (unicode.IsLetter(rune(p.input[p.pos+1])) || unicode.IsDigit(rune(p.input[p.pos+1]))) {
		return false
	}
	p.pos++
	return true
}

// parsePower is right-associative. Fish: unary binds tighter than ^, so -2^2 = (-2)^2 = 4.
// Achieved by having parsePower call parseUnary for the base.
func (p *parser) parsePower() (float64, error) {
	base, err := p.parseUnary()
	if err != nil {
		return 0, err
	}
	if p.match("^") {
		// Right-associative: exponent is another parsePower call
		exp, err := p.parsePower()
		if err != nil {
			return 0, err
		}
		return math.Pow(base, exp), nil
	}
	return base, nil
}

func (p *parser) parseUnary() (float64, error) {
	p.skipWhitespace()
	if p.match("-") {
		val, err := p.parseUnary()
		return -val, err
	}
	if p.match("+") {
		return p.parseUnary()
	}
	if p.match("!") {
		val, err := p.parseUnary()
		if err != nil {
			return 0, err
		}
		if val == 0 {
			return 1, nil
		}
		return 0, nil
	}
	if p.match("~") {
		val, err := p.parseUnary()
		if err != nil {
			return 0, err
		}
		return float64(^int64(val)), nil
	}
	return p.parsePrimary()
}

func (p *parser) parsePrimary() (float64, error) {
	p.skipWhitespace()
	if p.pos >= len(p.input) {
		return 0, fmt.Errorf("unexpected end of expression")
	}

	// Parenthesized expression
	if p.input[p.pos] == '(' {
		p.pos++
		val, err := p.parseOr()
		if err != nil {
			return 0, err
		}
		p.skipWhitespace()
		if p.pos >= len(p.input) || p.input[p.pos] != ')' {
			return 0, fmt.Errorf("missing closing parenthesis")
		}
		p.pos++
		return val, nil
	}

	// Hex literal
	if strings.HasPrefix(p.input[p.pos:], "0x") || strings.HasPrefix(p.input[p.pos:], "0X") {
		p.pos += 2
		start := p.pos
		for p.pos < len(p.input) && isHexDigit(p.input[p.pos]) {
			p.pos++
		}
		if p.pos == start {
			return 0, fmt.Errorf("invalid hex literal")
		}
		var n int64
		fmt.Sscanf(p.input[start:p.pos], "%x", &n)
		return float64(n), nil
	}

	// Number
	if unicode.IsDigit(rune(p.input[p.pos])) || p.input[p.pos] == '.' {
		return p.parseNumber()
	}

	// Identifier (constant or function)
	if unicode.IsLetter(rune(p.input[p.pos])) || p.input[p.pos] == '_' {
		return p.parseIdent()
	}

	return 0, fmt.Errorf("unexpected character: %q", string(p.input[p.pos]))
}

func (p *parser) parseNumber() (float64, error) {
	start := p.pos
	for p.pos < len(p.input) && (unicode.IsDigit(rune(p.input[p.pos])) || p.input[p.pos] == '.') {
		p.pos++
	}
	// Scientific notation
	if p.pos < len(p.input) && (p.input[p.pos] == 'e' || p.input[p.pos] == 'E') {
		p.pos++
		if p.pos < len(p.input) && (p.input[p.pos] == '+' || p.input[p.pos] == '-') {
			p.pos++
		}
		for p.pos < len(p.input) && unicode.IsDigit(rune(p.input[p.pos])) {
			p.pos++
		}
	}
	var f float64
	_, err := fmt.Sscanf(p.input[start:p.pos], "%f", &f)
	return f, err
}

func (p *parser) parseIdent() (float64, error) {
	start := p.pos
	for p.pos < len(p.input) && (unicode.IsLetter(rune(p.input[p.pos])) || unicode.IsDigit(rune(p.input[p.pos])) || p.input[p.pos] == '_') {
		p.pos++
	}
	name := p.input[start:p.pos]

	// Check for function call
	p.skipWhitespace()
	if p.pos < len(p.input) && p.input[p.pos] == '(' {
		p.pos++ // consume '('
		var args []float64
		p.skipWhitespace()
		if p.pos < len(p.input) && p.input[p.pos] != ')' {
			arg, err := p.parseOr()
			if err != nil {
				return 0, err
			}
			args = append(args, arg)
			p.skipWhitespace()
			for p.pos < len(p.input) && p.input[p.pos] == ',' {
				p.pos++
				arg, err := p.parseOr()
				if err != nil {
					return 0, err
				}
				args = append(args, arg)
				p.skipWhitespace()
			}
		}
		if p.pos >= len(p.input) || p.input[p.pos] != ')' {
			return 0, fmt.Errorf("missing closing parenthesis")
		}
		p.pos++
		return callFunc(name, args)
	}

	// Constant
	switch name {
	case "pi":
		return math.Pi, nil
	case "e":
		return math.E, nil
	case "tau":
		return 2 * math.Pi, nil
	case "inf":
		return math.Inf(1), nil
	case "true":
		return 1, nil
	case "false":
		return 0, nil
	}
	return 0, fmt.Errorf("unknown variable: %q", name)
}

func callFunc(name string, args []float64) (float64, error) {
	req := func(n int) error {
		if len(args) != n {
			return fmt.Errorf("%s() requires %d argument(s), got %d", name, n, len(args))
		}
		return nil
	}
	switch name {
	case "abs", "fabs":
		if err := req(1); err != nil {
			return 0, err
		}
		return math.Abs(args[0]), nil
	case "sin":
		if err := req(1); err != nil {
			return 0, err
		}
		return math.Sin(args[0]), nil
	case "cos":
		if err := req(1); err != nil {
			return 0, err
		}
		return math.Cos(args[0]), nil
	case "tan":
		if err := req(1); err != nil {
			return 0, err
		}
		return math.Tan(args[0]), nil
	case "asin":
		if err := req(1); err != nil {
			return 0, err
		}
		return math.Asin(args[0]), nil
	case "acos":
		if err := req(1); err != nil {
			return 0, err
		}
		return math.Acos(args[0]), nil
	case "atan":
		if err := req(1); err != nil {
			return 0, err
		}
		return math.Atan(args[0]), nil
	case "atan2":
		if err := req(2); err != nil {
			return 0, err
		}
		return math.Atan2(args[0], args[1]), nil
	case "sinh":
		if err := req(1); err != nil {
			return 0, err
		}
		return math.Sinh(args[0]), nil
	case "cosh":
		if err := req(1); err != nil {
			return 0, err
		}
		return math.Cosh(args[0]), nil
	case "tanh":
		if err := req(1); err != nil {
			return 0, err
		}
		return math.Tanh(args[0]), nil
	case "sqrt":
		if err := req(1); err != nil {
			return 0, err
		}
		return math.Sqrt(args[0]), nil
	case "ln":
		if err := req(1); err != nil {
			return 0, err
		}
		return math.Log(args[0]), nil
	case "log":
		// fish: log is log base 10
		if err := req(1); err != nil {
			return 0, err
		}
		return math.Log10(args[0]), nil
	case "log2":
		if err := req(1); err != nil {
			return 0, err
		}
		return math.Log2(args[0]), nil
	case "log10":
		if err := req(1); err != nil {
			return 0, err
		}
		return math.Log10(args[0]), nil
	case "exp":
		if err := req(1); err != nil {
			return 0, err
		}
		return math.Exp(args[0]), nil
	case "floor":
		if err := req(1); err != nil {
			return 0, err
		}
		return math.Floor(args[0]), nil
	case "ceil":
		if err := req(1); err != nil {
			return 0, err
		}
		return math.Ceil(args[0]), nil
	case "round":
		if err := req(1); err != nil {
			return 0, err
		}
		return math.Round(args[0]), nil
	case "pow":
		if err := req(2); err != nil {
			return 0, err
		}
		return math.Pow(args[0], args[1]), nil
	case "max":
		if len(args) == 0 {
			return 0, fmt.Errorf("max() requires at least 1 argument")
		}
		m := args[0]
		for _, a := range args[1:] {
			if a > m {
				m = a
			}
		}
		return m, nil
	case "min":
		if len(args) == 0 {
			return 0, fmt.Errorf("min() requires at least 1 argument")
		}
		m := args[0]
		for _, a := range args[1:] {
			if a < m {
				m = a
			}
		}
		return m, nil
	// Bitwise functions (fish uses these instead of operators for clarity)
	case "fac":
		if err := req(1); err != nil {
			return 0, err
		}
		return factorial(int64(args[0])), nil
	case "ncr":
		if err := req(2); err != nil {
			return 0, err
		}
		n, r := int64(args[0]), int64(args[1])
		if r < 0 || r > n {
			return 0, nil
		}
		return factorial(n) / (factorial(r) * factorial(n-r)), nil
	case "npr":
		if err := req(2); err != nil {
			return 0, err
		}
		n, r := int64(args[0]), int64(args[1])
		if r < 0 || r > n {
			return 0, nil
		}
		return factorial(n) / factorial(n-r), nil
	case "bitand":
		if err := req(2); err != nil {
			return 0, err
		}
		return float64(int64(args[0]) & int64(args[1])), nil
	case "bitor":
		if err := req(2); err != nil {
			return 0, err
		}
		return float64(int64(args[0]) | int64(args[1])), nil
	case "bitxor":
		if err := req(2); err != nil {
			return 0, err
		}
		return float64(int64(args[0]) ^ int64(args[1])), nil
	case "gcd":
		if err := req(2); err != nil {
			return 0, err
		}
		return float64(gcd(int64(math.Abs(args[0])), int64(math.Abs(args[1])))), nil
	case "roundf":
		if err := req(2); err != nil {
			return 0, err
		}
		m := math.Pow(10, args[1])
		return math.Round(args[0]*m) / m, nil
	case "floorf":
		if err := req(2); err != nil {
			return 0, err
		}
		m := math.Pow(10, args[1])
		return math.Floor(args[0]*m) / m, nil
	case "ceilf":
		if err := req(2); err != nil {
			return 0, err
		}
		m := math.Pow(10, args[1])
		return math.Ceil(args[0]*m) / m, nil
	default:
		return 0, fmt.Errorf("unknown function: %q", name)
	}
}

func gcd(a, b int64) int64 {
	for b != 0 {
		a, b = b, a%b
	}
	return a
}

func factorial(n int64) float64 {
	if n < 0 {
		return math.NaN()
	}
	result := 1.0
	for i := int64(2); i <= n; i++ {
		result *= float64(i)
	}
	return result
}

func isHexDigit(b byte) bool {
	return (b >= '0' && b <= '9') || (b >= 'a' && b <= 'f') || (b >= 'A' && b <= 'F')
}
