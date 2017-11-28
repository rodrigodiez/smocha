package validate

import (
	"errors"
	"fmt"
	"github.com/rodrigodiez/smocha/types"
	"net/http"
)

func Status(res *http.Response, test types.Test) (bool, error) {

	if res.StatusCode != test.Should.HaveStatus {
		return false, errors.New(fmt.Sprintf("status code %d was expected but got %d instead", test.Should.HaveStatus, res.StatusCode))
	}

	return true, nil
}
