package validate

import (
	"net/http"
  "testing"
  "github.com/rodrigodiez/smocha/types"
	"strings"
 "io/ioutil"
)

func TestReturnsTrueWhenStringIsFound(t *testing.T) {
	reader := ioutil.NopCloser(strings.NewReader("foo"))
  response := &http.Response{Body: reader}
  test := &types.Test{Should:types.Should{Contain: "foo"}}

  result, _ := Contain(response, *test)

  if result!=true {
    t.Error("Expected true but got", result);
  }
}

func TestReturnsNilErrorWhenStringIsFound(t *testing.T) {
	reader := ioutil.NopCloser(strings.NewReader("foo"))
  response := &http.Response{Body: reader}
  test := &types.Test{Should:types.Should{Contain: "foo"}}

  _, err := Contain(response, *test)

  if err!=nil {
    t.Error("Expected nil error but got", err);
  }
}

func TestReturnsFalseWhenStringIsNotFound(t *testing.T) {
	reader := ioutil.NopCloser(strings.NewReader("foo"))
  response := &http.Response{Body: reader}
  test := &types.Test{Should:types.Should{Contain: "bar"}}

  result, _ := Contain(response, *test)

  if result!=false {
    t.Error("Expected false but got", result);
  }
}

func TestReturnsErrorWhenStringIsNotFound(t *testing.T) {
	reader := ioutil.NopCloser(strings.NewReader("foo"))
  response := &http.Response{Body: reader}
  test := &types.Test{Should:types.Should{Contain: "bar"}}

  _, err := Contain(response, *test)

  if err==nil {
    t.Error("Expected error but got", err);
  }
}
