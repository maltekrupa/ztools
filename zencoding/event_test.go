package zencoding

import "encoding/json"

type mockEventType uint8
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

func (t mockEventType) TypeName() string {
	return "mock"
}

func (t mockEventType) GetEmptyInstance() EventData {
	return new(mockEventData)
}

func (t mockEventType) MarshalJSON() ([]byte, error) {
	return json.Marshal(t.TypeName())
}

func (m *mockEventData) GetType() EventType {
	var t mockEventType
	return t
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
