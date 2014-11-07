package zencoding

type encodedGrab struct {
	Host   string            `json:"host"`
	Domain *string           `json:"domain"`
	Time   string            `json:"time"`
	Log    []ConnectionEvent `json:"log"`
}

type encodedConnectionEvent struct {
	Type  EventType `json:"type"`
	Data  EventData `json:"data"`
	Error *string   `json:"error"`
}
