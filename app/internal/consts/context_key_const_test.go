package consts

import (
	"reflect"
	"testing"
)

func TestCtxConstList(t *testing.T) {
	v := reflect.ValueOf(ContextKeys)
	tp := v.Type()

	expected := map[string]string{
		"Uuid":      "uuid",
		"Email":     "email",
		"RoomModel": "room_model",
		"IsAdmin":   "is_admin",
		"IsMember":  "is_member",
	}

	if tp.NumField() != len(expected) {
		t.Fatalf("number of fields mismatch: expected %d, got %d",
			len(expected), tp.NumField())
	}

	for i := 0; i < tp.NumField(); i++ {
		name := tp.Field(i).Name
		value := v.Field(i).String()

		expVal, ok := expected[name]
		if !ok {
			t.Errorf("unexpected field added: %s", name)
		}
		if value != expVal {
			t.Errorf("value mismatch for %s: expected %s, got %s",
				name, expVal, value)
		}
	}
}
