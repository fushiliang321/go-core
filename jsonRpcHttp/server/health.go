package server

type Health struct {
}

type params struct {
	Name string `json:"name,omitempty"`
}

var result = "success"

func (s *Health) Check(params *params) (*string, error) {
	return &result, nil
}
