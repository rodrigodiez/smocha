package types

type Test struct {
	Url    string
	Should Should
}

type Should struct {
	Contain         string
	HaveStatus      int      `yaml:"have_status"`
	MatchJsonSchema string   `yaml:"match_json_schema"`
	HaveHeaders     []Header `yaml:"have_headers"`
}

type Header struct {
	Name  string
	Value string
}
