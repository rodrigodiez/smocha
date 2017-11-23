package main

import (
  "flag"
  "fmt"
  "github.com/jinzhu/configor"
  "github.com/fatih/color"
  "net/http"
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

	configor.Load(&testbook, *filename)
  done := make(chan bool, len(testbook.Tests));

  for i := range testbook.Tests {
    go test(testbook.Tests[i], testbook.Host, testbook.Schema, done);
  }

  for _ = range testbook.Tests {
    <-done
  }
}

func test(test Test, host string, schema string, done chan bool) {
  resp, err := http.Get(fmt.Sprintf("%s://%s%s", schema, host, test.Url));

  if err != nil {
    log.Printf("%s%s %s", red("❌"), red(test.Url), err);
    done <- true
    return
  }

  if(test.Should.HaveStatus != 0) {

    if resp.StatusCode != test.Should.HaveStatus {
      log.Printf("%s%s %s", red("❌"), red(test.Url), fmt.Sprintf("status code %s was expected but got %s instead", yellow(test.Should.HaveStatus), yellow(resp.StatusCode)));
      done <- true
      return
    }
  }


  if(test.Should.MatchJsonSchema != "") {
    s, err := jsonSchema.ReadFile(test.Should.MatchJsonSchema);
    if err != nil {
      log.Printf("%s%s %s", red("❌"), red(test.Url), err);
      done <- true
      return
    }

    b := jsonValBuilder.New();
    v, err := b.Build(s);

    if err != nil {
      log.Printf("%s%s %s", red("❌"), red(test.Url), err);
      done <- true
      return
    }

    var input interface{};
    json.Unmarshal([]byte("{\"name\":\"john\",\"age\":22,\"class\":\"mca\"}"), &input)

    if err := v.Validate(input); err != nil {
      log.Printf("%s%s %s", red("❌"), red(test.Url), err);
      done <- true
      return
    }
  }

  log.Printf("%s %s", green("✔"), green(test.Url));
  done <- true
}
