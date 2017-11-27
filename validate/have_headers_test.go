package validate

import (
	"net/http"
  "testing"
  "github.com/rodrigodiez/smocha/types"
)

func TestReturnsTrueWhenTheHeaderIsFoundWithTheRightValue(t *testing.T) {
  headers := &http.Header{}
  headers.Add("Foo-Header", "Bar-Value")
  response := &http.Response{Header:*headers}

  test := &types.Test{Should:types.Should{HaveHeaders: []types.Header{{Name: "Foo-Header", Value: "Bar-Value"}}}}

  result, _ := HaveHeaders(response, *test)

  if result!=true {
    t.Error("Expected true but got", result);
  }
}

func TestReturnsFalseWhenTheHeaderIsNotFound(t *testing.T) {
  headers := &http.Header{}
  headers.Add("Foo-Header", "Bar-Value")
  response := &http.Response{Header:*headers}

  test := &types.Test{Should:types.Should{HaveHeaders: []types.Header{{Name: "Baz-Header", Value: "Bar-Value"}}}}

  result, _ := HaveHeaders(response, *test)

  if result!=false {
    t.Error("Expected false but got", result);
  }
}

func TestReturnsFalseWhenTheHeaderIsFoundButDoesNotContainTheRightValue(t *testing.T) {
  headers := &http.Header{}
  headers.Add("Foo-Header", "Bar-Value")
  response := &http.Response{Header:*headers}

  test := &types.Test{Should:types.Should{HaveHeaders: []types.Header{{Name: "Foo-Header", Value: "Qux-Value"}}}}

  result, _ := HaveHeaders(response, *test)

  if result!=false {
    t.Error("Expected false but got", result);
  }
}

func TestReturnsTrueWhenAllHeadersAreFoundWithTheRightValue(t *testing.T) {
  headers := &http.Header{}
  headers.Add("Foo-Header", "Bar-Value")
  headers.Add("Baz-Header", "Qux-Value")
  response := &http.Response{Header:*headers}

  test := &types.Test{Should:types.Should{HaveHeaders: []types.Header{{Name: "Foo-Header", Value: "Bar-Value"}, {Name: "Baz-Header", Value: "Qux-Value"}}}}

  result, _ := HaveHeaders(response, *test)

  if result!=true {
    t.Error("Expected true but got", result);
  }
}

func TestReturnsFalseWhenNotAllHeadersAreFoundWithTheRightValue(t *testing.T) {
  headers := &http.Header{}
  headers.Add("Foo-Header", "Bar-Value")
  headers.Add("Baz-Header", "Qux-Value")
  response := &http.Response{Header:*headers}

  test := &types.Test{Should:types.Should{HaveHeaders: []types.Header{{Name: "Foo-Header", Value: "Bar-Value"}, {Name: "Baz-Header", Value: "LOL-Value"}}}}

  result, _ := HaveHeaders(response, *test)

  if result!=false {
    t.Error("Expected false but got", result);
  }
}
