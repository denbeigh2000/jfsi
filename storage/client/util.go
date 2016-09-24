package client

import (
	"fmt"
	"io"
	"io/ioutil"
)

func handleUnknownError(code int, r io.Reader) error {
	body, err := ioutil.ReadAll(r)
	if err != nil {
		return fmt.Errorf("Code %v received, but error reading body: %v",
			code, err)
	}

	return fmt.Errorf("%v: %v", code, string(body))
}
