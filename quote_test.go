package quotes

import (
	"reflect"
	"testing"
)

func TestQuote_SerializeDeserialize(t *testing.T) {
	tests := []struct {
		name  string
		quote Quote
	}{
		{"Test01", Quote{ID: 42, Author: "Test", Text: "This is a test", Source: "unknown"}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			serialized, err := tt.quote.Serialize()
			if err != nil {
				t.Errorf("Quote.Serialize() error = %v", err)
				return
			}
			restored := Quote{}
			err = restored.Deserialize(serialized)
			if err != nil {
				t.Errorf("Quote.Deserialize() error = %v", err)
				return
			}
			if !reflect.DeepEqual(tt.quote, restored) {
				t.Errorf("Quote.Serialize() -> Quote.Deserialize() = %#v, want %#v", restored, tt.quote)
			}
		})
	}
}
