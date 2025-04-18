package standard

import (
	"strconv"
	"strings"
)

type Error struct {
	Parent                   error
	File                     string
	FileLine                 int
	Function                 string
	InvokeFunctionIndex      int
	InvokeFunctionParamIndex int
	Dependence               string
	RequiredBy               []Error
}

func (a *Error) Error() string {
	var msg strings.Builder
	a.error(&msg)
	return msg.String()
}

func (a *Error) error(msg *strings.Builder) {
	if a.Parent != nil {
		msg.WriteString(a.Parent.Error())
	}
	if a.Dependence != "" {
		if msg.Len() > 0 {
			msg.WriteRune(' ')
		}
		msg.WriteString(a.Dependence)
	}
	if a.InvokeFunctionParamIndex > 0 {
		if msg.Len() > 0 {
			msg.WriteRune(' ')
		}
		msg.WriteString("(argument #")
		msg.WriteString(strconv.Itoa(a.InvokeFunctionParamIndex - 1))
		msg.WriteString(")")
	}
	switch {
	case a.Function != "":
		if msg.Len() > 0 {
			msg.WriteString(" at ")
		}
		msg.WriteString(a.Function)
		msg.WriteRune(':')
		msg.WriteString(strconv.Itoa(a.FileLine))
	case a.File != "":
		if msg.Len() > 0 {
			msg.WriteString(" at ")
		}
		msg.WriteString(a.File)
		msg.WriteRune(':')
		msg.WriteString(strconv.Itoa(a.FileLine))
	}
	if a.InvokeFunctionIndex > 0 {
		if msg.Len() > 0 {
			msg.WriteRune(' ')
		}
		msg.WriteString("(invoke #")
		msg.WriteString(strconv.Itoa(a.InvokeFunctionIndex - 1))
		msg.WriteString(")")
	}
	for _, err := range a.RequiredBy {
		if msg.Len() > 0 {
			msg.WriteString(" ← ")
		}
		msg.WriteString("required by")
		err.error(msg)
	}
}

func (a *Error) Unwrap() error {
	return a.Parent
}
