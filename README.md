# Smocha
Smocha is a simple smoke tests framework built in Go.

It uses `yaml` files called _testbooks_ that contain your tests as configuration instead of code.

# Features
At the moment the features of Smocha are rather simple. It allows:

- Concurrent tests using _goroutines_
- HTTP / HTTPS GET requests
- Check the status code of the response
- Check that the response contains a string
- Match the JSON of the response against a _Json Schema_ file
- stdout / stderr logging
- error exit code on failed tests
- allow overrides using environment variables

# Testbook example
```yaml
host: myhost.com
schema: https
tests:
- url: /
  should:
    have_status: 200
    contain: success!

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
# testbook file defaults to testbook.yml if not specified
smocha
# you can also specify your own file
smocha my-testbook-file.yml
# you can also override some values in the testbook by using env variables
env SMOCHA_HOST=google.com smocha
```
