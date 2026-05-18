package mathcmd

import (
	"bytes"
	"os"
	"os/exec"
	"strings"
	"testing"
)

// mathBin is the path to the compiled binary, set by TestMain.
var mathBin string

func TestMain(m *testing.M) {
	tmp, err := os.CreateTemp("", "math-test-*")
	if err != nil {
		panic(err)
	}
	tmp.Close()
	mathBin = tmp.Name()
	defer os.Remove(mathBin)

	if out, err := exec.Command("go", "build", "-o", mathBin, "github.com/mattmc3/mudskipper/cmd/math").CombinedOutput(); err != nil {
		panic("build failed: " + string(out))
	}

	os.Exit(m.Run())
}

type cliTest struct {
	name     string
	stdin    string
	args     []string
	wantExit int
	wantOut  []string // nil = don't check; use []string{} to assert empty
	wantErr  string   // substring match; empty = don't check
}

func runBin(t *testing.T, tc cliTest) {
	t.Helper()
	cmd := exec.Command(mathBin, tc.args...)
	cmd.Stdin = strings.NewReader(tc.stdin)
	var outBuf, errBuf bytes.Buffer
	cmd.Stdout = &outBuf
	cmd.Stderr = &errBuf

	err := cmd.Run()
	gotExit := 0
	if err != nil {
		if ee, ok := err.(*exec.ExitError); ok {
			gotExit = ee.ExitCode()
		} else {
			t.Fatalf("exec error: %v", err)
		}
	}

	if gotExit != tc.wantExit {
		t.Errorf("exit: want %d, got %d", tc.wantExit, gotExit)
	}

	if tc.wantOut != nil {
		if len(tc.wantOut) == 0 {
			if outBuf.Len() != 0 {
				t.Errorf("stdout: want empty, got %q", outBuf.String())
			}
		} else {
			got := strings.TrimRight(outBuf.String(), "\n")
			want := strings.Join(tc.wantOut, "\n")
			if got != want {
				t.Errorf("stdout: want %q, got %q", want, got)
			}
		}
	}

	if tc.wantErr != "" {
		if !strings.Contains(errBuf.String(), tc.wantErr) {
			t.Errorf("stderr: want %q in %q", tc.wantErr, errBuf.String())
		}
	}
}

// Fish parity tests — direct from https://github.com/fish-shell/fish-shell/blob/master/tests/checks/math.fish
// Differences from fish noted inline.

var fishParityMathBasic = []cliTest{
	// math.fish L6: math 3 / 2
	{name: "float_division", args: []string{"3 / 2"}, wantExit: 0, wantOut: []string{"1.5"}},
	// math.fish L8: math 10/6
	{name: "repeating_decimal", args: []string{"10/6"}, wantExit: 0, wantOut: []string{"1.666667"}},
	// math.fish L10: math -s0 10 / 6
	{name: "scale0", args: []string{"-s0", "10 / 6"}, wantExit: 0, wantOut: []string{"1"}},
	// math.fish L12: math 'floor(10 / 6)'
	{name: "floor_func", args: []string{"floor(10 / 6)"}, wantExit: 0, wantOut: []string{"1"}},
	// math.fish L14: math -s3 10/6
	{name: "scale3", args: []string{"-s3", "10/6"}, wantExit: 0, wantOut: []string{"1.667"}},
	// math.fish L16: math '10 % 6'
	{name: "modulo", args: []string{"10 % 6"}, wantExit: 0, wantOut: []string{"4"}},
	// math.fish L18: math -s0 '10 % 6'
	{name: "scale0_modulo", args: []string{"-s0", "10 % 6"}, wantExit: 0, wantOut: []string{"4"}},
	// math.fish L20: math '23 % 7'
	{name: "modulo2", args: []string{"23 % 7"}, wantExit: 0, wantOut: []string{"2"}},
	// math.fish L22: math --scale=6 '5 / 3 * 0.3'
	{name: "scale6_multiply", args: []string{"--scale=6", "5 / 3 * 0.3"}, wantExit: 0, wantOut: []string{"0.500000"}},
	// math.fish L24: math --scale=max '5 / 3'
	{name: "scale_max", args: []string{"--scale=max", "5 / 3"}, wantExit: 0, wantOut: []string{"1.666666666666667"}},
	// math.fish L26: error on non-zero scale + non-base-10
	{name: "scale_base_conflict_error", args: []string{"--scale=1", "--base=16", "2 / 3 - 1"}, wantExit: 1, wantErr: "invalid option combination"},
	// math.fish L28: math --scale=abc '5 / 3'
	{name: "scale_invalid_string", args: []string{"--scale=abc", "5 / 3"}, wantExit: 1, wantErr: "invalid scale"},
	// math.fish L30: math --scale=-1 '5 / 3'
	{name: "scale_negative_error", args: []string{"--scale=-1", "5 / 3"}, wantExit: 1, wantErr: "invalid scale"},
	// math.fish L32: math --scale=16 '5 / 3'
	{name: "scale_too_large", args: []string{"--scale=16", "5 / 3"}, wantExit: 1, wantErr: "invalid scale"},
	// math.fish L34: math "7^2"
	{name: "power", args: []string{"7^2"}, wantExit: 0, wantOut: []string{"49"}},
	// math.fish L36: math -1 + 1  (multi-arg expression)
	{name: "multi_arg_expr", args: []string{"-1", "+", "1"}, wantExit: 0, wantOut: []string{"0"}},
	// math.fish L38: math '-2 * -2'
	{name: "negative_multiply", args: []string{"-2 * -2"}, wantExit: 0, wantOut: []string{"4"}},
	// math.fish L40: math 5 \* -2
	{name: "multiply_args", args: []string{"5", "*", "-2"}, wantExit: 0, wantOut: []string{"-10"}},
	// math.fish L42: math -- -4 / 2
	{name: "double_dash_negative", args: []string{"--", "-4 / 2"}, wantExit: 0, wantOut: []string{"-2"}},
}

var fishParityMathFunctions = []cliTest{
	// math.fish L48: math 'max(1,2)'
	{name: "max", args: []string{"max(1,2)"}, wantExit: 0, wantOut: []string{"2"}},
	// math.fish L49: math 'min(1,2)'
	{name: "min", args: []string{"min(1,2)"}, wantExit: 0, wantOut: []string{"1"}},
	// math.fish L54: math 'round(3/2)'
	{name: "round", args: []string{"round(3/2)"}, wantExit: 0, wantOut: []string{"2"}},
	// math.fish L55: math 'floor(3/2)'
	{name: "floor", args: []string{"floor(3/2)"}, wantExit: 0, wantOut: []string{"1"}},
	// math.fish L56: math 'ceil(3/2)'
	{name: "ceil", args: []string{"ceil(3/2)"}, wantExit: 0, wantOut: []string{"2"}},
	// math.fish L60: math 'round(-3/2)'
	{name: "round_neg", args: []string{"round(-3/2)"}, wantExit: 0, wantOut: []string{"-2"}},
	// math.fish L61: math 'floor(-3/2)'
	{name: "floor_neg", args: []string{"floor(-3/2)"}, wantExit: 0, wantOut: []string{"-2"}},
	// math.fish L62: math 'ceil(-3/2)'
	{name: "ceil_neg", args: []string{"ceil(-3/2)"}, wantExit: 0, wantOut: []string{"-1"}},
	// math.fish L241: math 'log 16' → log10(16) = 1.20412  (note: fish 'log' is log base 10)
	{name: "log_base10", args: []string{"log(16)"}, wantExit: 0, wantOut: []string{"1.20412"}},
	// math.fish L283: math 'log2(8)'
	{name: "log2", args: []string{"log2(8)"}, wantExit: 0, wantOut: []string{"3"}},
	// math.fish L286: math 'log(8) / log(2)'
	{name: "log_division", args: []string{"log(8) / log(2)"}, wantExit: 0, wantOut: []string{"3"}},
	// math.fish L227: math "bitand(0xFE, 1)"
	{name: "bitand", args: []string{"bitand(0xFE, 1)"}, wantExit: 0, wantOut: []string{"0"}},
	// math.fish L229: math "bitor(0xFE, 1)"
	{name: "bitor", args: []string{"bitor(0xFE, 1)"}, wantExit: 0, wantOut: []string{"255"}},
	// math.fish L231: math "bitxor(5, 1)"
	{name: "bitxor", args: []string{"bitxor(5, 1)"}, wantExit: 0, wantOut: []string{"4"}},
	// math.fish L331: math min 2  (variadic with one arg)
	{name: "min_one_arg", args: []string{"min(2)"}, wantExit: 0, wantOut: []string{"2"}},
	// math.fish L333: math min 2, 3, 4, 5, -10, 1
	{name: "min_variadic", args: []string{"min(2,3,4,5,-10,1)"}, wantExit: 0, wantOut: []string{"-10"}},
}

var fishParityMathIntegers = []cliTest{
	// math.fish L68-71: integer formatting
	{name: "int_1", args: []string{"1"}, wantExit: 0, wantOut: []string{"1"}},
	{name: "int_1000", args: []string{"1000"}, wantExit: 0, wantOut: []string{"1000"}},
	// math.fish L76: math '10^15'
	{name: "power_15", args: []string{"10^15"}, wantExit: 0, wantOut: []string{"1000000000000000"}},
	// math.fish L77: math '-10^14'  (fish: unary binds tighter → (-10)^14 = 10^14)
	{name: "neg_power_even", args: []string{"-10^14"}, wantExit: 0, wantOut: []string{"100000000000000"}},
	// math.fish L78: math '-10^15'  ((-10)^15 = -10^15)
	{name: "neg_power_odd", args: []string{"-10^15"}, wantExit: 0, wantOut: []string{"-1000000000000000"}},
	// math.fish L92: math -2^2
	{name: "neg_base_power", args: []string{"-2^2"}, wantExit: 0, wantOut: []string{"4"}},
	// math.fish L95: math -s0 '1.0 / 2.0'
	{name: "scale0_float_div", args: []string{"-s0", "1.0 / 2.0"}, wantExit: 0, wantOut: []string{"0"}},
}

var fishParityMathBase = []cliTest{
	// math.fish L254: math --base=16 255 / 15
	{name: "hex_output", args: []string{"--base=16", "255 / 15"}, wantExit: 0, wantOut: []string{"0x11"}},
	// math.fish L256: math -bhex 16 x 2
	{name: "hex_alias_x_multiply", args: []string{"-bhex", "16 x 2"}, wantExit: 0, wantOut: []string{"0x20"}},
	// math.fish L258: math --base hex 12 + 0x50
	{name: "hex_alias_space", args: []string{"--base", "hex", "12 + 0x50"}, wantExit: 0, wantOut: []string{"0x5c"}},
	// math.fish L260: math --base hex 0
	{name: "hex_zero", args: []string{"--base", "hex", "0"}, wantExit: 0, wantOut: []string{"0x0"}},
	// math.fish L262: math --base hex -1
	{name: "hex_negative", args: []string{"--base", "hex", "-1"}, wantExit: 0, wantOut: []string{"-0x1"}},
	// math.fish L264: math --base hex -15
	{name: "hex_negative_15", args: []string{"--base", "hex", "-15"}, wantExit: 0, wantOut: []string{"-0xf"}},
	// math.fish L268: math --base octal 0
	{name: "octal_zero", args: []string{"--base", "octal", "0"}, wantExit: 0, wantOut: []string{"0"}},
	// math.fish L270: math --base octal -1
	{name: "octal_negative", args: []string{"--base", "octal", "-1"}, wantExit: 0, wantOut: []string{"-01"}},
	// math.fish L278: math --base notabase
	{name: "invalid_base", args: []string{"--base", "notabase"}, wantExit: 1, wantErr: "invalid base"},
}

var fishParityMathStdin = []cliTest{
	// math.fish L375: echo 5 + 6 | math
	{name: "stdin_expr", stdin: "5 + 6\n", args: []string{}, wantExit: 0, wantOut: []string{"11"}},
}

func TestFishParityMath(t *testing.T) {
	run := func(group string, tests []cliTest) {
		t.Run(group, func(t *testing.T) {
			for _, tc := range tests {
				t.Run(tc.name, func(t *testing.T) {
					runBin(t, tc)
				})
			}
		})
	}
	run("basic", fishParityMathBasic)
	run("functions", fishParityMathFunctions)
	run("integers", fishParityMathIntegers)
	run("base", fishParityMathBase)
	run("stdin", fishParityMathStdin)
}
