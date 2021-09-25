package grpc

import "errors"

func ErrorFromStr(s string) error {
	if s == "" {
		return nil
	}
	return errors.New(s)
}

func ErrorToStr(err error) string {
	if err == nil {
		return ""
	}
	return err.Error()
}
