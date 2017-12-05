package main

// Test defines the structure of a test
type Test struct {
	URL    string `yaml:"url"`
	Should Should
}

// Should defines the list of validators enabled for this test
type Should struct {
	Contain           string
	HaveStatus        int      `yaml:"have_status"`
	MatchesJSONSchema string   `yaml:"match_json_schema"`
	HaveHeaders       []Header `yaml:"have_headers"`
}

// Header defines a header
type Header struct {
	Name  string
	Value string
}

// Testbook defines the structure of a testbook
type Testbook struct {
	Host   string
	Schema string
	Rate   int `default:"30"`
	Tests  []Test
}
