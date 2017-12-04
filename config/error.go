package config

import (
	"bytes"
	"fmt"

	"gopkg.in/go-playground/validator.v9"
)

type validateErrorWrapper struct {
	inner error
}

func (err *validateErrorWrapper) Error() string {
	buf := &bytes.Buffer{}

	if es, ok := err.inner.(validator.ValidationErrors); ok {
		fmt.Fprintln(buf, "kuma: validation failed")
		for _, e := range es {
			fmt.Fprintln(buf, e)
		}
	} else {
		fmt.Fprint(buf, err.inner.Error())
	}

	return buf.String()
}
