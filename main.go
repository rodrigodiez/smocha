package main

import (
  "flag"
  "fmt"
  "github.com/jinzhu/configor"
  "github.com/fatih/color"
  "net/http"
  "os"
  "log"
  "encoding/json"
  jsonSchema "github.com/lestrrat/go-jsschema"
  jsonValBuilder "github.com/lestrrat/go-jsval/builder"
)

type Test struct {
  Url string
  Should struct {
    HaveStatus int `yaml:"have_status"`
    MatchJsonSchema string `yaml:"match_json_schema"`
  }
}

type Testbook struct {
  Host string
  Schema string
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

  for i := range testbook.Tests {
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
    json.NewDecoder(res.Body).Decode(&input)

    if err := v.Validate(input); err != nil {
      log.Printf("%s %s %s", red("✗"), red(test.Url), err);
      done <- false
      return
    }
  }

  log.Printf("%s %s", green("✔"), green(test.Url));
  done <- true
}
