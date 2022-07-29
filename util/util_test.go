package util

import (
	"testing"
	"time"
)

func TestToFloat(t *testing.T) {
	s := "0.3"
	if f := ToFloat(s); f != 0.3 {
		t.Errorf("Expected: 0.3 and got: %v\n", f)
	}
	es := ""
	if f := ToFloat(es); f != 0 {
		t.Errorf("Expected: 0 and got: %v\n", f)
	}
}

func TestToInt(t *testing.T) {
	s := "3"
	if f := ToFloat(s); f != 3 {
		t.Errorf("Expected: 3 and got: %v\n", f)
	}
	es := ""
	if f := ToFloat(es); f != 0 {
		t.Errorf("Expected: 0 and got: %v\n", f)
	}
}
func TestToSQ(t *testing.T) {
	if s := SQ("test"); "'test'" != s {
		t.Errorf("Expected: 'test' and got: %s\n", s)
	}
}
func TestToString(t *testing.T) {
	if s := ToString(1); "1" != s {
		t.Errorf("Expected: 1 and got: %s\n", s)
	}
}
func TestToDate(t *testing.T) {
	input := "2022-07-28 00:00:00 +0530 IST"
	if d := ToDate(input); d.String() != input {
		t.Errorf("Expected: %s and got: %s\n", input, d.String())
	}
	wrongInput := "2022-07-28"
	expected := time.Time{}
	if d := ToDate(wrongInput); d != expected {
		t.Errorf("Expected: %s and got: %s\n", input, d.String())
	}
}
func TestToJSON(t *testing.T) {
	expected := "{\"Name\":\"gopal\"}"
	input := struct {
		Name string
	}{Name: "gopal"}
	if s := ToJSONString(input); s != expected {
		t.Errorf("Expected: %s and got: %s\n", expected, s)
	}
	if s := ToJSONString(func() {}); s != string(internalServerErrorResponse) {
		t.Errorf("Expected: %s and got: %s\n", string(internalServerErrorResponse), s)
	}
}
