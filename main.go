package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/fatih/color"
	"github.com/jinzhu/configor"
)

var testbook Testbook
var green = color.New(color.FgGreen).SprintFunc()
var red = color.New(color.FgRed).SprintFunc()
var yellow = color.New(color.FgYellow).SprintFunc()

func main() {
	var failedTests = 0

	filename := flag.String("testbook", "testbook.yml", "a testbook file")
	flag.Parse()

	if _, err := os.Stat(*filename); err != nil {
		defer os.Exit(1)
	}

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

func test(test Test, host string, schema string, ch chan bool) {
	res, err := http.Get(fmt.Sprintf("%s://%s%s", schema, host, test.URL))

	if err != nil {
		printErr(test, err)
		ch <- false
		return
	}
	validator := NewResponseValidator(res)

	if test.Should.HaveStatus != 0 {
		result, err := validator.HasStatus(test.Should.HaveStatus)
		if result == false {
			printErr(test, err)
			ch <- false
			return
		}
	}

	if test.Should.MatchesJSONSchema != "" {
		reader, err := os.Open(test.Should.MatchesJSONSchema)

		if err != nil {
			printErr(test, err)
			ch <- false
			return
		}
		result, err := validator.MatchesJSONSchema(reader)
		if result == false {
			printErr(test, err)
			ch <- false
			return
		}
	}

	if test.Should.Contain != "" {
		result, err := validator.Contains(test.Should.Contain)
		if result == false {
			printErr(test, err)
			ch <- false
			return
		}
	}

	if len(test.Should.HaveHeaders) > 0 {
		result, err := validator.HasHeaders(test.Should.HaveHeaders)
		if result == false {
			printErr(test, err)
			ch <- false
			return
		}
	}

	fmt.Printf("%s %s\n", green("✔"), green(test.URL))
	ch <- true
}

func printErr(test Test, err error) {
	log.Printf("%s%s %s", red("✗"), red(test.URL), err)
}
