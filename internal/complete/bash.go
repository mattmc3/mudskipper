package completecmd

import (
	"fmt"
	"io"
	"strings"
	"text/template"
)

var bashAppendTmpl = template.Must(template.New("bash-append").Delims("<%", "%>").Parse(
	`<% if .Condition %><% .CondVar %>+=(<% .Condition %> <% .Word %>)
<% else %><% .WordsVar %>+=(<% .Word %>)
<% end %>`))

var bashRegisterTmpl = template.Must(template.New("bash-register").Delims("<%", "%>").Parse(
	`declare -f <%.Fn%> > /dev/null 2>&1 || {
  <%.Fn%>() {
    local cur="${COMP_WORDS[COMP_CWORD]}"
    local -a _words=("${<%.WordsVar%>[@]}")
    local _i
    for (( _i=0; _i<${#<%.CondVar%>[@]}; _i+=2 )); do
      eval "${<%.CondVar%>[_i]}" 2>/dev/null && _words+=("${<%.CondVar%>[_i+1]}")
    done
<% if .KeepOrder %>    local _w
    for _w in "${_words[@]}"; do
      [[ "$_w" == "$cur"* ]] && COMPREPLY+=("$_w")
    done
<% else %>    COMPREPLY=($(compgen -W "${_words[*]}" -- "$cur"))
<% end %><% if .ForceFiles %>    COMPREPLY+=($(compgen -f -- "$cur"))
<% end %>  }
  builtin complete -F <%.Fn%> <%.Target%>
}
`))

var bashEraseTmpl = template.Must(template.New("bash-erase").Delims("<%", "%>").Parse(
	`unset <%.WordsVar%> <%.CondVar%> 2>/dev/null
unset -f <%.Fn%> 2>/dev/null
builtin complete -r <%.Target%> 2>/dev/null
`))

var bashWrapsTmpl = template.Must(template.New("bash-wraps").Delims("<%", "%>").Parse(
	`declare -f <%.WrapFn%> > /dev/null 2>&1 && builtin complete -F <%.WrapFn%> <%.Target%>
`))

type bashData struct {
	Fn         string
	WordsVar   string
	CondVar    string
	Target     string
	Word       string
	Condition  string
	WrapFn     string
	KeepOrder  bool
	ForceFiles bool
}

func emitBash(sp spec, stdout, stderr io.Writer) int {
	cmd := sp.effectiveCommand()
	if cmd == "" {
		fmt.Fprintln(stderr, "complete: -c/--command or -p/--path required")
		return 1
	}

	safeName := cmdToIdent(cmd)
	d := bashData{
		Fn:         "_complete_" + safeName,
		WordsVar:   "_complete_" + safeName + "_words",
		CondVar:    "_complete_" + safeName + "_cond_words",
		Target:     sp.compdefTarget(),
		Condition:  bashQuote(sp.condition),
		WrapFn:     "_complete_" + cmdToIdent(sp.wraps),
		KeepOrder:  sp.keepOrder,
		ForceFiles: sp.forceFiles,
	}

	switch {
	case sp.erase:
		return execTmpl(bashEraseTmpl, d, stdout, stderr)
	case sp.wraps != "":
		return execTmpl(bashWrapsTmpl, d, stdout, stderr)
	}

	for _, w := range bashSpecWords(sp) {
		d.Word = bashQuote(w)
		if err := bashAppendTmpl.Execute(stdout, d); err != nil {
			fmt.Fprintf(stderr, "complete: %v\n", err)
			return 1
		}
	}
	return execTmpl(bashRegisterTmpl, d, stdout, stderr)
}

func bashSpecWords(sp spec) []string {
	var words []string
	if sp.longOption != "" {
		words = append(words, "--"+sp.longOption)
	}
	if sp.shortOption != "" {
		words = append(words, "-"+sp.shortOption)
	}
	if sp.oldOption != "" {
		words = append(words, "-"+sp.oldOption)
	}
	if sp.longOption == "" && sp.shortOption == "" && sp.oldOption == "" && sp.arguments != "" {
		words = append(words, strings.Fields(sp.arguments)...)
	}
	return words
}

func bashQuote(s string) string {
	if s == "" {
		return ""
	}
	if !strings.ContainsAny(s, " \t\n\"'\\$`") {
		return s
	}
	return "'" + strings.ReplaceAll(s, "'", `'\''`) + "'"
}
