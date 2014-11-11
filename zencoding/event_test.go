package zencoding

import "encoding/json"

var (
	mockEventType = EventType{
		TypeName:         "mock",
		GetEmptyInstance: newMockEvent,
	}
)

func newMockEvent() EventData {
	return new(mockEventData)
}

type mockEventData struct {
	A string
	B int
	C *string
}

type encodedMockEvent struct {
	A string
	B int
	C *string
}

func (m *mockEventData) GetType() EventType {
	return mockEventType
}

func (m *mockEventData) MarshalJSON() ([]byte, error) {
	e := encodedMockEvent{
		A: m.A,
		B: m.B,
		C: m.C,
	}
	return json.Marshal(&e)
}

func (m *mockEventData) UnmarshalJSON(b []byte) error {
	e := encodedMockEvent{}
	if err := json.Unmarshal(b, &e); err != nil {
		return err
	}
	m.A = e.A
	m.B = e.B
	m.C = e.C
	return nil
}

func (m *mockEventData) saneDefaults() *mockEventData {
	m.A = "a"
	m.B = 123
	m.C = nil
	return m
}
