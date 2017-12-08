package main

import (
	"io/ioutil"
	"net/http"
	"strings"
	"testing"
)

func TestHasStatus(t *testing.T) {
	t.Run("Right status", testHasStatusFunc(200, 200, true, false))
	t.Run("Wrong status", testHasStatusFunc(200, 404, false, true))
}

func TestMatchesJSONSchema(t *testing.T) {
	t.Run("JSON matches schema", testMatchesJSONSchemaFunc(`{"title":"Foo","type":"object","properties":{"bar":{"type":"string"}},"required":["bar"]}`, `{"bar":"baz"}`, true, false))
	t.Run("JSON does not match schema", testMatchesJSONSchemaFunc(`{"title":"Foo","type":"object","properties":{"bar":{"type":"string"}},"required":["bar"]}`, `{"qux":"baz"}`, false, true))
	t.Run("Schema is invalid", testMatchesJSONSchemaFunc(`this is not a schema`, `{"bar":"baz"}`, false, true))
}

func TestHasHeaders(t *testing.T) {
	t.Run("Header is present, right value", testHasHeadersFunc([]Header{{Name: "Foo", Value: "Bar"}}, map[string]string{"Foo": "Bar"}, true, false))
	t.Run("Headers are present, right values", testHasHeadersFunc([]Header{{Name: "Foo", Value: "Bar"}, {Name: "Baz", Value: "Qux"}}, map[string]string{"Foo": "Bar", "Baz": "Qux"}, true, false))
	t.Run("Header are present, wrong value", testHasHeadersFunc([]Header{{Name: "Foo", Value: "Bar"}}, map[string]string{"Foo": "Qux"}, false, true))
	t.Run("Headers are present, some wrong values", testHasHeadersFunc([]Header{{Name: "Foo", Value: "Bar"}, {Name: "Baz", Value: "Qux"}}, map[string]string{"Foo": "NaN", "Baz": "Qux"}, false, true))
	t.Run("Header is not present", testHasHeadersFunc([]Header{{Name: "Baz", Value: "Qux"}}, map[string]string{"Foo": "Bar"}, false, true))
}

func TestContains(t *testing.T) {
	t.Run("Body equals needle", testContainFunc("foo", "foo", true, false))
	t.Run("Body contains needle", testContainFunc("foo", "a foo mal", true, false))
	t.Run("Body does not contain needle", testContainFunc("foo", "bar", false, true))
	t.Run("Body is empty", testContainFunc("foo", "", false, true))
}

func testHasStatusFunc(expectedStatus int, responseStatus int, expectedResult bool, expectedError bool) func(*testing.T) {
	return func(t *testing.T) {
		response := &http.Response{StatusCode: responseStatus}
		validator := NewResponseValidator(response)

		result, err := validator.HasStatus(expectedStatus)

		handleResults(t, expectedResult, result, expectedError, err)
	}
}

func testMatchesJSONSchemaFunc(jsonSchema string, responseJSON string, expectedResult bool, expectedError bool) func(*testing.T) {
	return func(t *testing.T) {
		bodyReader := ioutil.NopCloser(strings.NewReader(responseJSON))
		response := &http.Response{Body: bodyReader}
		validator := NewResponseValidator(response)

		result, err := validator.MatchesJSONSchema(strings.NewReader(jsonSchema))

		handleResults(t, expectedResult, result, expectedError, err)
	}
}

func testHasHeadersFunc(expectedHeaders []Header, responseHeaders map[string]string, expectedResult bool, expectedError bool) func(*testing.T) {
	return func(t *testing.T) {
		headers := http.Header{}

		for headerName, headerValue := range responseHeaders {
			headers.Add(headerName, headerValue)
		}
		response := http.Response{Header: headers}
		validator := NewResponseValidator(&response)

		result, err := validator.HasHeaders(expectedHeaders)

		handleResults(t, expectedResult, result, expectedError, err)
	}
}

func testContainFunc(needle string, bodyContent string, expectedResult bool, expectedError bool) func(*testing.T) {
	return func(t *testing.T) {
		reader := ioutil.NopCloser(strings.NewReader(bodyContent))
		response := http.Response{Body: reader}
		validator := NewResponseValidator(&response)

		result, err := validator.Contains(needle)

		handleResults(t, expectedResult, result, expectedError, err)
	}
}

func handleResults(t *testing.T, expectedResult bool, result bool, expectedError bool, err error) {
	if result != expectedResult {
		t.Errorf("Expected result to be %t but got %t", expectedResult, result)
	}

	if (expectedError == true) && (err == nil) {
		t.Error("Expected error got", err)
	}

	if (expectedError == false) && (err != nil) {
		t.Error("Expected nil error but got", err)
	}
}
