package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"strings"

	jsonSchema "github.com/lestrrat/go-jsschema"
	jsonValBuilder "github.com/lestrrat/go-jsval/builder"
)

// ResponseValidator validates its Response
type ResponseValidator struct {
	response *http.Response
}

// NewResponseValidator creates a new ResponseValidator from a Response
func NewResponseValidator(res *http.Response) *ResponseValidator {
	return &ResponseValidator{response: res}
}

// HasStatus checks whether the Response has a given status code
func (r *ResponseValidator) HasStatus(status int) (bool, error) {

	if r.response.StatusCode != status {
		return false, errors.New(fmt.Sprintf("status code %d was expected but got %d instead", status, r.response.StatusCode))
	}

	return true, nil
}

// MatchesJSONSchema checks whether the Response body matches the Json schema
// defined in a given file
func (r *ResponseValidator) MatchesJSONSchema(in io.Reader) (bool, error) {
	schema, err := jsonSchema.Read(in)
	if err != nil {
		return false, err
	}

	b := jsonValBuilder.New()
	jsonValidator, err := b.Build(schema)

	if err != nil {
		return false, err
	}

	var input interface{}
	bodyBytes, _ := ioutil.ReadAll(r.response.Body)
	defer r.response.Body.Close()
	r.response.Body = ioutil.NopCloser(bytes.NewBuffer(bodyBytes))

	json.Unmarshal(bodyBytes, &input)

	if err := jsonValidator.Validate(input); err != nil {
		return false, err
	}

	return true, nil
}

// HasHeaders checks whether the Response has all the provided headers with
// their values
func (r *ResponseValidator) HasHeaders(headers []Header) (bool, error) {

	for i := range headers {
		expectedHeader := headers[i]
		headerValue := r.response.Header.Get(expectedHeader.Name)

		if headerValue == "" {
			return false, errors.New(fmt.Sprintf("response header '%s' was not found", expectedHeader.Name))
		}

		if headerValue != expectedHeader.Value {
			return false, errors.New(fmt.Sprintf("response header '%s' expected to contain '%s' but '%s' was found instead", expectedHeader.Name, expectedHeader.Value, headerValue))
		}
	}

	return true, nil
}

// Contains checks whether the response body contains a given string
func (r *ResponseValidator) Contains(needle string) (bool, error) {

	bodyBytes, _ := ioutil.ReadAll(r.response.Body)
	defer r.response.Body.Close()
	r.response.Body = ioutil.NopCloser(bytes.NewBuffer(bodyBytes))

	if false == strings.Contains(string(bodyBytes), needle) {
		return false, errors.New(fmt.Sprintf("response does not contain '%s'", needle))
	}

	return true, nil
}
