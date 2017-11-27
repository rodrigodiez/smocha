package types

type Testbook struct {
	Host   string
	Schema string
	Rate   int `default:"30"`
	Tests  []Test
}
