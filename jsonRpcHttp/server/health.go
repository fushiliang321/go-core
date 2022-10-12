package server

type Health struct {
}

type params struct {
	Name string `json:"name"`
}

func (s *Health) Check(params *params, result *string) error {
	*result = "success"
	return nil
}
