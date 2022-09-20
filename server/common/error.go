package common

type MyError struct {
	code int
	msg  string
}

func New(code int, msg string) error {
	return &MyError{
		code: code,
		msg:  msg,
	}
}

func (m *MyError) Error() string {
	return m.msg
}

func (m *MyError) Code() int {
	return m.code
}
