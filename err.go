package nolang

import "fmt"

const (
	ErrCodeEnd   = "CodeEnd"
	ErrWrongType = "WrongType"
)

type NoLangErr struct {
	Line   int
	Reason string
}

func (e *NoLangErr) Error() string {
	return e.String()
}

func (e *NoLangErr) String() string {
	return fmt.Sprintf("Err: '%s' Line: %d", e.Reason, e.Line)
}

func newError(reason string, line int) error {
	return &NoLangErr{Line: line, Reason: reason}
}

func isErrEq(e error, reason string) bool {
	if nlErr, ok := e.(*NoLangErr); ok {
		return nlErr.Reason == reason
	}
	return false
}
