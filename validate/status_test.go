package validate

import (
	"net/http"
  "testing"
  "github.com/rodrigodiez/smocha/types"
)

func TestReturnsTrueOnRightStatus(t *testing.T) {
  response := &http.Response{StatusCode: 200}
  test := &types.Test{Should:types.Should{HaveStatus: 200}}

  result, _ := Status(response, *test)

  if result!=true {
    t.Error("Expected true but got", result);
  }
}

func TestReturnsNilErrorOnRightStatus(t *testing.T) {
  response := &http.Response{StatusCode: 200}
  test := &types.Test{Should:types.Should{HaveStatus: 200}}

  _, err := Status(response, *test)

  if err!=nil {
    t.Error("Expected nil error but got", err);
  }
}

func TestReturnsFalseOnWrongStatus(t *testing.T) {
  response := &http.Response{StatusCode: 404}
  test := &types.Test{Should:types.Should{HaveStatus: 200}}

  result, _ := Status(response, *test)

  if result!=false {
    t.Error("Expected false but got", result);
  }
}

func TestReturnsErrorOnWrongStatus(t *testing.T) {
  response := &http.Response{StatusCode: 404}
  test := &types.Test{Should:types.Should{HaveStatus: 200}}

  _, err := Status(response, *test)

  if err==nil {
    t.Error("Expected an error but got", err);
  }
}
