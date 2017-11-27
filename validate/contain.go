package validate

import (
	"bytes"
	"github.com/rodrigodiez/smocha/types"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
)

func Contain(res *http.Response, test types.Test) bool {

	bodyBytes, _ := ioutil.ReadAll(res.Body)
	defer res.Body.Close()
	res.Body = ioutil.NopCloser(bytes.NewBuffer(bodyBytes))

	if false == strings.Contains(string(bodyBytes), test.Should.Contain) {
		log.Printf("%s %s %s '%s'", red("✗"), red(test.Url), "response does not contain", test.Should.Contain)
		return false
	}

	return true
}
