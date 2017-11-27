package validate

import (
	"errors"
	"fmt"
	"github.com/rodrigodiez/smocha/types"
	"net/http"
)

func HaveHeaders(res *http.Response, test types.Test) (bool, error) {


	for i := range test.Should.HaveHeaders {
		expectedHeader := test.Should.HaveHeaders[i]
		headerValue := res.Header.Get(expectedHeader.Name)

		if headerValue == "" {
			return false, errors.New(fmt.Sprintf("response header '%s' was not found", expectedHeader.Name))
		}

		if headerValue != expectedHeader.Value {
			return false, errors.New(fmt.Sprintf("response header '%s' expected to contain '%s' but '%s' was found instead", expectedHeader.Name, expectedHeader.Value, headerValue));
		}
	}

	return true, nil
}
