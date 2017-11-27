package validate

import (
	"bytes"
	"encoding/json"
	jsonSchema "github.com/lestrrat/go-jsschema"
	jsonValBuilder "github.com/lestrrat/go-jsval/builder"
	"github.com/rodrigodiez/smocha/types"
	"io/ioutil"
	"net/http"
)

func MatchJsonSchema(res *http.Response, test types.Test) (bool, error) {
	s, err := jsonSchema.ReadFile(test.Should.MatchJsonSchema)
	if err != nil {
		return false, err
	}

	b := jsonValBuilder.New()
	v, err := b.Build(s)

	if err != nil {
		return false, err
	}

	var input interface{}
	bodyBytes, _ := ioutil.ReadAll(res.Body)
	defer res.Body.Close()
	res.Body = ioutil.NopCloser(bytes.NewBuffer(bodyBytes))

	json.Unmarshal(bodyBytes, &input)

	if err := v.Validate(input); err != nil {
		return false, err
	}

	return true, nil
}
