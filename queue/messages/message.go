package messages

type Message struct {
	RoutingKey string `json:"routing_key"`
	Body       []byte `json:"body"`
}
