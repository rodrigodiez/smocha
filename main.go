package main

import (
  "encoding/json"
  "flag"
  "fmt"
  "github.com/jinzhu/configor"
  "github.com/fatih/color"
  "io/ioutil"
  "log"
  "net/http"
  "os"
  "strings"
  "time"
  jsonSchema "github.com/lestrrat/go-jsschema"
  jsonValBuilder "github.com/lestrrat/go-jsval/builder"
)

type Test struct {
  Url string
  Should struct {
    Contain string
    HaveStatus int `yaml:"have_status"`
    MatchJsonSchema string `yaml:"match_json_schema"`

  }
}

type Testbook struct {
  Host string
  Schema string
  Rate int `default:"30"`
  Tests []Test
}


var testbook Testbook;
var green = color.New(color.FgGreen).SprintFunc();
var red = color.New(color.FgRed).SprintFunc();
var yellow = color.New(color.FgYellow).SprintFunc();

func main() {
  filename := flag.String("testbook", "testbook.yml", "testbook file");

	configor.New(&configor.Config{ENVPrefix: "SMOCHA"}).Load(&testbook, *filename)
  done := make(chan bool, len(testbook.Tests));
  throttle := time.Tick(time.Second / time.Duration(testbook.Rate));

  fmt.Printf("Starting %s tests on %s at %s requests per second\n", yellow(testbook.Schema), yellow(testbook.Host), yellow(testbook.Rate));

  for i := range testbook.Tests {
    <-throttle
    go test(testbook.Tests[i], testbook.Host, testbook.Schema, done);
  }

  for _ = range testbook.Tests {
    if ok := <-done; ok == false {
      defer os.Exit(1)
    }
  }
}

func test(test Test, host string, schema string, done chan bool) {
  res, err := http.Get(fmt.Sprintf("%s://%s%s", schema, host, test.Url));

  if err != nil {
    log.Printf("%s%s %s", red("✗"), red(test.Url), err);
    done <- false
    return
  }

  bodyBytes, err := ioutil.ReadAll(res.Body)
  defer res.Body.Close()

  if(test.Should.HaveStatus != 0) {

    if res.StatusCode != test.Should.HaveStatus {
      log.Printf("%s%s %s", red("✗"), red(test.Url), fmt.Sprintf("status code %s was expected but got %s instead", yellow(test.Should.HaveStatus), yellow(res.StatusCode)));
      done <- false
      return
    }
  }


  if(test.Should.MatchJsonSchema != "") {
    s, err := jsonSchema.ReadFile(test.Should.MatchJsonSchema);
    if err != nil {
      log.Printf("%s %s %s", red("✗"), red(test.Url), err);
      done <- false
      return
    }

    b := jsonValBuilder.New();
    v, err := b.Build(s);

    if err != nil {
      log.Printf("%s %s %s", red("✗"), red(test.Url), err);
      done <- false
      return
    }

    var input interface{};
    json.Unmarshal(bodyBytes, &input);

    if err := v.Validate(input); err != nil {
      log.Printf("%s %s %s", red("✗"), red(test.Url), err);
      done <- false
      return
    }
  }

  if(test.Should.Contain != "") {
    bodyString := string(bodyBytes)
    if false == strings.Contains(bodyString, test.Should.Contain) {
      log.Printf("%s %s %s '%s'", red("✗"), red(test.Url), "response does not contain", test.Should.Contain);
      done <- false
      return
    }
  }

  fmt.Printf("%s %s\n", green("✔"), green(test.Url));
  done <- true
}
