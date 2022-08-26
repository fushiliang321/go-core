package types

type Producer struct {
	Exchange   string
	RoutingKey string
	Data       any
}
