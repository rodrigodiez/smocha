package main

import (
	"flag"
	"fmt"
	"github.com/fatih/color"
	"github.com/jinzhu/configor"
	"github.com/rodrigodiez/smocha/types"
	"github.com/rodrigodiez/smocha/validate"
	"log"
	"net/http"
	"os"
	"time"
)

var testbook types.Testbook
var green = color.New(color.FgGreen).SprintFunc()
var red = color.New(color.FgRed).SprintFunc()
var yellow = color.New(color.FgYellow).SprintFunc()

func main() {
	var failedTests = 0

	filename := flag.String("testbook", "testbook.yml", "a testbook file")
	flag.Parse()

	configor.New(&configor.Config{ENVPrefix: "SMOCHA"}).Load(&testbook, *filename)
	ch := make(chan bool, len(testbook.Tests))
	throttle := time.Tick(time.Second / time.Duration(testbook.Rate))

	fmt.Printf("Starting %s tests on %s at %s requests per second\n", yellow(testbook.Schema), yellow(testbook.Host), yellow(testbook.Rate))

	for i := range testbook.Tests {
		<-throttle
		go test(testbook.Tests[i], testbook.Host, testbook.Schema, ch)
	}

	for _ = range testbook.Tests {
		if ok := <-ch; ok == false {
			failedTests++
			defer os.Exit(1)
		}
	}

	fmt.Printf("%d tests (%s, %s)\n", len(testbook.Tests), green(fmt.Sprintf("%d passed", len(testbook.Tests)-failedTests)), red(fmt.Sprintf("%d failed", failedTests)))
}

func test(test types.Test, host string, schema string, ch chan bool) {
	res, err := http.Get(fmt.Sprintf("%s://%s%s", schema, host, test.Url))

	if err != nil {
		printErr(test, err)
		ch <- false
		return
	}

	if test.Should.HaveStatus != 0 {
		result, err := validate.Status(res, test)
		if result == false {
			printErr(test, err)
			ch <- false
			return
		}
	}

	if test.Should.MatchJsonSchema != "" {
		result, err := validate.MatchJsonSchema(res, test)
		if result == false {
			printErr(test, err)
			ch <- false
			return
		}
	}

	if test.Should.Contain != "" {
		result, err := validate.Contain(res, test)
		if result == false {
			printErr(test, err)
			ch <- false
			return
		}
	}

	fmt.Printf("%s %s\n", green("✔"), green(test.Url))
	ch <- true
}

func printErr(test types.Test, err error) {
	log.Printf("%s%s %s", red("✗"), red(test.Url), err)
}
