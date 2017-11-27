package validate

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/rodrigodiez/smocha/types"
	"io/ioutil"
	"net/http"
	"strings"
)

func Contain(res *http.Response, test types.Test) (bool, error) {

	bodyBytes, _ := ioutil.ReadAll(res.Body)
	defer res.Body.Close()
	res.Body = ioutil.NopCloser(bytes.NewBuffer(bodyBytes))

	if false == strings.Contains(string(bodyBytes), test.Should.Contain) {
		return false, errors.New(fmt.Sprintf("response does not contain '%s'", test.Should.Contain))
	}

	return true, nil
}
