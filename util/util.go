package util

import (
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	"github.com/gopalrohra/flyapi/log"
)

var internalServerErrorResponse, _ = json.Marshal(map[string]string{"status": "Error", "message": "Oops, internal server error."})

// ToJSONString converts any value to json string.
func ToJSONString(o interface{}) string {
	r, err := json.Marshal(o)
	if err != nil {
		return string(internalServerErrorResponse)
	}
	return string(r)
}

// ToFloat converts string to float32 & handles error
func ToFloat(v string) float32 {
	r, err := strconv.ParseFloat(v, 32)
	if err != nil {
		log.Error(err.Error())
		return 0
	}
	return float32(r)
}

// ToInt converts string to int and handles error.
func ToInt(v string) int {
	r, err := strconv.ParseInt(v, 10, 32)
	if err != nil {
		log.Error(err.Error())
		return 0
	}
	return int(r)
}

// ToDate converts string to date and handles error
func ToDate(v string) time.Time {
	layout := "2006-01-02 15:04:05 -0700 MST"
	t, err := time.Parse(layout, v)
	if err != nil {
		log.Error(err.Error())
		return time.Time{}
	}
	return t
}

// SQ takes a string and surrounds it with ' for db operation
func SQ(s string) string {
	return fmt.Sprintf("'%s'", s)
}
func ToString(v int) string {
	return fmt.Sprint(v)
}
