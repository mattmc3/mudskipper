package completecmd

import (
	"fmt"
	"io"
	"strings"
	"text/template"
)

var zshAppendTmpl = template.Must(template.New("zsh-append").Delims("<%", "%>").Parse(
	`<% if .Condition %><%- .CondVar %>+=(<% .Condition %> <% .OptSpec %>)
<% else %><% .OptsVar %>+=(<% .OptSpec %>)
<% end %>`))

var zshRegisterTmpl = template.Must(template.New("zsh-register").Delims("<%", "%>").Parse(
	`typeset -f <%.Fn%> > /dev/null 2>&1 || {
  <%.Fn%>() {
    local -a _opts=("${<%.OptsVar%>[@]}")
    local _i
    for (( _i=1; _i<=${#<%.CondVar%>[@]}; _i+=2 )); do
      eval "${<%.CondVar%>[_i]}" 2>/dev/null && _opts+=("${<%.CondVar%>[_i+1]}")
    done
    _arguments -s "${_opts[@]}"
  }
  compdef <%.Fn%> <%.Target%>
}
`))

var zshEraseTmpl = template.Must(template.New("zsh-erase").Delims("<%", "%>").Parse(
	`unset <%.OptsVar%> <%.CondVar%> 2>/dev/null
unfunction <%.Fn%> 2>/dev/null
compdef -d <%.Target%>
`))

var zshWrapsTmpl = template.Must(template.New("zsh-wraps").Delims("<%", "%>").Parse(
	`typeset -f <%.WrapFn%> > /dev/null 2>&1 && compdef <%.WrapFn%> <%.Target%>
`))

type zshData struct {
	Fn        string
	OptsVar   string
	CondVar   string
	Target    string
	OptSpec   string
	Condition string
	WrapFn    string
}

func emitZsh(sp spec, stdout, stderr io.Writer) int {
	cmd := sp.effectiveCommand()
	if cmd == "" {
		fmt.Fprintln(stderr, "complete: -c/--command or -p/--path required")
		return 1
	}

	safeName := cmdToIdent(cmd)
	d := zshData{
		Fn:        "_complete_" + safeName,
		OptsVar:   "_complete_" + safeName + "_opts",
		CondVar:   "_complete_" + safeName + "_cond_opts",
		Target:    sp.compdefTarget(),
		OptSpec:   zshOptSpec(sp),
		Condition: zshSingleQuote(sp.condition),
		WrapFn:    "_complete_" + cmdToIdent(sp.wraps),
	}

	switch {
	case sp.erase:
		return execTmpl(zshEraseTmpl, d, stdout, stderr)
	case sp.wraps != "":
		return execTmpl(zshWrapsTmpl, d, stdout, stderr)
	}
	if d.OptSpec != "" {
		if err := zshAppendTmpl.Execute(stdout, d); err != nil {
			fmt.Fprintf(stderr, "complete: %v\n", err)
			return 1
		}
	}
	return execTmpl(zshRegisterTmpl, d, stdout, stderr)
}

func execTmpl(t *template.Template, data any, stdout, stderr io.Writer) int {
	if err := t.Execute(stdout, data); err != nil {
		fmt.Fprintf(stderr, "complete: %v\n", err)
		return 1
	}
	return 0
}

func zshOptSpec(sp spec) string {
	hasFlag := sp.shortOption != "" || sp.longOption != "" || sp.oldOption != ""
	if !hasFlag {
		if sp.arguments != "" {
			return fmt.Sprintf("'*: :(%s)'", sp.arguments)
		}
		return ""
	}

	desc := zshEscDesc(sp.description)
	var specs []string
	for _, flag := range zshFlagVariants(sp) {
		var entry string
		switch {
		case sp.arguments != "":
			entry = fmt.Sprintf("'%s=[%s]: :(%s)'", flag, desc, sp.arguments)
		case sp.requireParam && sp.noFiles:
			entry = fmt.Sprintf("'%s=[%s]: :( )'", flag, desc)
		case sp.requireParam:
			entry = fmt.Sprintf("'%s=[%s]: :_files'", flag, desc)
		default:
			entry = fmt.Sprintf("'%s[%s]'", flag, desc)
		}
		specs = append(specs, entry)
	}
	return strings.Join(specs, " ")
}

func zshFlagVariants(sp spec) []string {
	var variants []string
	if sp.shortOption != "" && sp.longOption != "" {
		variants = append(variants, fmt.Sprintf("{-%s,--%s}", sp.shortOption, sp.longOption))
	} else if sp.shortOption != "" {
		variants = append(variants, "-"+sp.shortOption)
	} else if sp.longOption != "" {
		variants = append(variants, "--"+sp.longOption)
	}
	if sp.oldOption != "" {
		variants = append(variants, "-"+sp.oldOption)
	}
	return variants
}

func zshEscDesc(s string) string {
	s = strings.ReplaceAll(s, `\`, `\\`)
	s = strings.ReplaceAll(s, "'", `'\''`)
	s = strings.ReplaceAll(s, "[", `\[`)
	s = strings.ReplaceAll(s, "]", `\]`)
	return s
}

func zshSingleQuote(s string) string {
	if s == "" {
		return ""
	}
	return "'" + strings.ReplaceAll(s, "'", `'\''`) + "'"
}
