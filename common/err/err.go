package err

import "fmt"

type Err struct {
	code    int
	message string
}

func New(code int, message string) Err {
	return Err{code: code, message: message}
}

func (err Err) Error() string {
	return fmt.Sprintf("code: %v, message: %s", err.code, err.message)
}
