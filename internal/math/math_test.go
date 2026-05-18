package mathcmd

import (
	"strings"
	"testing"
)

func runMath(args ...string) (int, string, string) {
	var out, err strings.Builder
	code := Run(args, strings.NewReader(""), &out, &err)
	return code, out.String(), err.String()
}

func assertOut(t *testing.T, want, got string) {
	t.Helper()
	got = strings.TrimRight(got, "\n")
	if got != want {
		t.Errorf("want %q, got %q", want, got)
	}
}

// math.fish L6: math 3 / 2
func TestMath_float_division(t *testing.T) {
	_, out, _ := runMath("3 / 2")
	assertOut(t, "1.5", out)
}

// math.fish L8: math 10/6
func TestMath_repeating_decimal(t *testing.T) {
	_, out, _ := runMath("10/6")
	assertOut(t, "1.666667", out)
}

// math.fish L10: math -s0 10 / 6
func TestMath_scale0(t *testing.T) {
	_, out, _ := runMath("-s0", "10 / 6")
	assertOut(t, "1", out)
}

// math.fish L12: math 'floor(10 / 6)'
func TestMath_floor_func(t *testing.T) {
	_, out, _ := runMath("floor(10 / 6)")
	assertOut(t, "1", out)
}

// math.fish L14: math -s3 10/6
func TestMath_scale3(t *testing.T) {
	_, out, _ := runMath("-s3", "10/6")
	assertOut(t, "1.667", out)
}

// math.fish L16: math '10 % 6'
func TestMath_modulo(t *testing.T) {
	_, out, _ := runMath("10 % 6")
	assertOut(t, "4", out)
}

// math.fish L22: math --scale=max '5 / 3'
func TestMath_scale_max(t *testing.T) {
	_, out, _ := runMath("--scale=max", "5 / 3")
	assertOut(t, "1.666666666666667", out)
}

// math.fish L24: math --scale=1 --base=16 "2 / 3 - 1" → error
func TestMath_scale_and_base_error(t *testing.T) {
	code, _, stderr := runMath("--scale=1", "--base=16", "2 / 3 - 1")
	if code == 0 {
		t.Error("want error exit")
	}
	if !strings.Contains(stderr, "invalid option combination") {
		t.Errorf("stderr: want 'invalid option combination', got %q", stderr)
	}
}

// math.fish L29: math "7^2"
func TestMath_power(t *testing.T) {
	_, out, _ := runMath("7^2")
	assertOut(t, "49", out)
}

// math.fish L31: math -1 + 1
func TestMath_multi_arg_expression(t *testing.T) {
	_, out, _ := runMath("-1", "+", "1")
	assertOut(t, "0", out)
}

// math.fish L33: math '-2 * -2'
func TestMath_negative_multiply(t *testing.T) {
	_, out, _ := runMath("-2 * -2")
	assertOut(t, "4", out)
}

// math.fish L39: math 'max(1,2)' / 'min(1,2)'
func TestMath_max_min(t *testing.T) {
	_, out, _ := runMath("max(1,2)")
	assertOut(t, "2", out)
	_, out, _ = runMath("min(1,2)")
	assertOut(t, "1", out)
}

// math.fish L44-L52: round/floor/ceil
func TestMath_round_floor_ceil(t *testing.T) {
	tests := []struct{ expr, want string }{
		{"round(3/2)", "2"},
		{"floor(3/2)", "1"},
		{"ceil(3/2)", "2"},
		{"round(-3/2)", "-2"},
		{"floor(-3/2)", "-2"},
		{"ceil(-3/2)", "-1"},
	}
	for _, tt := range tests {
		_, out, _ := runMath(tt.expr)
		assertOut(t, tt.want, out)
	}
}

// math.fish L57-L60: large integers
func TestMath_large_integers(t *testing.T) {
	tests := []struct{ expr, want string }{
		{"1", "1"}, {"10", "10"}, {"100", "100"}, {"1000", "1000"},
		{"10^15", "1000000000000000"},
	}
	for _, tt := range tests {
		_, out, _ := runMath(tt.expr)
		assertOut(t, tt.want, out)
	}
}

// math.fish: constants
func TestMath_constants(t *testing.T) {
	code, out, _ := runMath("pi")
	if code != 0 {
		t.Error("want exit 0")
	}
	if !strings.HasPrefix(strings.TrimRight(out, "\n"), "3.14159") {
		t.Errorf("pi: got %q", out)
	}
}

// base output
func TestMath_base16(t *testing.T) {
	// fish outputs 0x prefix for hex base
	_, out, _ := runMath("--base=16", "255")
	assertOut(t, "0xff", out)
}

func TestMath_base2(t *testing.T) {
	_, out, _ := runMath("-b2", "8")
	assertOut(t, "1000", out)
}

// Trig
func TestMath_sin_cos(t *testing.T) {
	_, out, _ := runMath("-s0", "sin(0)")
	assertOut(t, "0", out)
	_, out, _ = runMath("-s0", "cos(0)")
	assertOut(t, "1", out)
}

// Comparison operators error — use `test` instead
func TestMath_comparison_errors(t *testing.T) {
	tests := []string{"1 < 2", "2 > 1", "1 <= 1", "1 >= 1"}
	for _, expr := range tests {
		code, _, stderr := runMath(expr)
		if code == 0 {
			t.Errorf("%q: want error exit, got 0", expr)
		}
		if !strings.Contains(stderr, "logical operations") {
			t.Errorf("%q: stderr: want 'logical operations', got %q", expr, stderr)
		}
	}
}

// Bitwise
func TestMath_bitwise(t *testing.T) {
	_, out, _ := runMath("3 & 5")
	assertOut(t, "1", out)
	_, out, _ = runMath("3 | 5")
	assertOut(t, "7", out)
	_, out, _ = runMath("1 << 3")
	assertOut(t, "8", out)
}

// Invalid scale errors
func TestMath_gcd(t *testing.T) {
	_, out, _ := runMath("gcd(12,8)")
	assertOut(t, "4", out)
	_, out, _ = runMath("gcd(100,75)")
	assertOut(t, "25", out)
	_, out, _ = runMath("gcd(7,13)")
	assertOut(t, "1", out)
}

func TestMath_roundf(t *testing.T) {
	_, out, _ := runMath("roundf(1.666667, 2)")
	assertOut(t, "1.67", out)
	_, out, _ = runMath("roundf(-1.5, 0)")
	assertOut(t, "-2", out)
}

func TestMath_floorf(t *testing.T) {
	_, out, _ := runMath("floorf(1.666, 2)")
	assertOut(t, "1.66", out)
	_, out, _ = runMath("floorf(-1.666, 2)")
	assertOut(t, "-1.67", out)
}

func TestMath_ceilf(t *testing.T) {
	_, out, _ := runMath("ceilf(1.661, 2)")
	assertOut(t, "1.67", out)
	_, out, _ = runMath("ceilf(-1.666, 2)")
	assertOut(t, "-1.66", out)
}

func TestMath_invalid_scale(t *testing.T) {
	code, _, stderr := runMath("--scale=abc", "5/3")
	if code == 0 {
		t.Error("want error")
	}
	if !strings.Contains(stderr, "invalid scale") {
		t.Errorf("stderr: %q", stderr)
	}
}
