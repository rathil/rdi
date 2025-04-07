package standard

type wrapError struct {
	msg string
	err error
}

func (a *wrapError) Error() string {
	msg := a.msg
	if a.err != nil {
		msg += ", err: " + a.err.Error()
	}
	return msg
}

func (a *wrapError) Unwrap() error {
	return a.err
}
