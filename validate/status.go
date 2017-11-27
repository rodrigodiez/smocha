package validate

import (
	"fmt"
	"github.com/fatih/color"
	"github.com/rodrigodiez/smocha/types"
	"log"
	"net/http"
)

var green = color.New(color.FgGreen).SprintFunc()
var red = color.New(color.FgRed).SprintFunc()
var yellow = color.New(color.FgYellow).SprintFunc()

func Status(res *http.Response, test types.Test) bool {

	if res.StatusCode != test.Should.HaveStatus {
		log.Printf("%s%s %s", red("✗"), red(test.Url), fmt.Sprintf("status code %s was expected but got %s instead", yellow(test.Should.HaveStatus), yellow(res.StatusCode)))
		return false
	}

	return true
}
