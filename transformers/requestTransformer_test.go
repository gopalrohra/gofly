package transformers

import (
	"net/http"
	"testing"
)

func TestParseParameters(t *testing.T) {
	r, err := http.NewRequest("GET", "/orders/1", http.NoBody)
	if err != nil {
		t.Error("Failed to initialize mocked http.request")
	}
	transformer := NewRequestTransformer(r, "/orders/:id")
	if transformer.pathParameters["id"] != "1" {
		t.Errorf("Expected: 1 and got: %s\n", transformer.pathParameters["id"])
	}
}
