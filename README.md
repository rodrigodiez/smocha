# Smocha
Smocha is a simple smoke tests framework built in Go.

It uses `yaml` files called _testbooks_ that contain your tests as configuration instead of code.

# Features
At the moment the features of Smocha are rather simple. It allows:

- Concurrent, throttled tests using _Goroutines_ (default rate is 30 rps although you can override in your testbook)
- Throttled tests
- HTTP / HTTPS GET requests
- Check the status code of the response
- Check that the response contains a string
- Check that the response contains some headers
- Match the JSON of the response against a _Json Schema_ file
- Sensible stdout/stderr logging
- Non-zero exit code on failed tests
- Allow overrides of YAML values using environment variables

# Testbook example
```yaml
host: myhost.com
schema: https
rate: 20
tests:
- url: /
  should:
    have_status: 200
    contain: success
    have_headers:
      - { name: 'Content-Type', value: 'application/javascript'}
      - { name: 'X-Custom-Header', value: 'foo'}

- url: /json-endpoint
  should:
    have_status: 200
    match_json_schema: MyOwnSchema.schema.json
```

# Installation
## Go get
```
go get -u github.com/rodrigodiez/smocha
```
## Binary download
Binaries for Linux, macOS and Windows are available in the [Releases](https://github.com/rodrigodiez/smocha/releases) section of this repository

## Brew
Brew installation is not available yet but it will be ready soon. Check [#6](https://github.com/rodrigodiez/smocha/issues/6) for updated details

# Running Smocha
```bash
# testbook file defaults to testbook.yml if not specified
smocha
# you can also specify your own file
smocha my-testbook-file.yml
# you can also override some values in the testbook by using env variables
env SMOCHA_HOST=google.com smocha
```
