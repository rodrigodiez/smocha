package types

type Test struct {
	Url    string
	Should Should
}

type Should struct {
	Contain         string
	HaveStatus      int    `yaml:"have_status"`
	MatchJsonSchema string `yaml:"match_json_schema"`
}
