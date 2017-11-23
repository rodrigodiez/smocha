# Smocha
Smocha is a simple smoke tests framework built in Go.

It uses `yaml` files called _testbooks_ that contain your tests as configuration instead of code.

# Features
At the moment the features of Smocha are rather simple. It allows:

- Concurrent tests using _goroutines_
- HTTP / HTTPS GET requests
- Validate the status code of the response
- Match the JSON of the response against a _Json Schema_ file
- stdout / stderr logging
- error exit code on failed tests

# Testbook example
```yaml
host: myhost.com
schema: https
tests:
- url: /
  should:
    have_status: 200

- url: /json-endpoint
  should:
    have_status: 200
    match_json_schema: MyOwnSchema.schema.json
```

# Installation
```
go get -u github.com/rodrigodiez/smocha
```

# Running Smocha
```bash
#file defaults to testbook.yml if not specified
smocha my-testbook-file.yml
```
