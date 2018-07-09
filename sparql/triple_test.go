package sparql

import (
	"testing"
)

func TestNewTriple(t *testing.T) {
	var graphTests = []struct {
		in  string
		out Triple
	}{
		{
			in: "<http://example/book3> <http://purl.org/dc/elements/1.1/title> \"A new book\"",
			out: Triple{
				Subject:   "<http://example/book3>",
				Predicate: "<http://purl.org/dc/elements/1.1/title>",
				Object:    "\"A new book\"",
			},
		},
		{
			in: "<http://example/book3> dc:title \"A new book\"",
			out: Triple{
				Subject:   "<http://example/book3>",
				Predicate: "http://purl.org/dc/elements/1.1/title",
				Object:    "\"A new book\"",
			},
		},
	}

	query := Query{
		Prefixes: map[string]string{
			"dc": "http://purl.org/dc/elements/1.1/",
		},
	}

	for _, tt := range graphTests {
		actual := query.NewTriple(tt.in, "")
		if actual != tt.out {
			t.Errorf("ERROR")
		}
	}
}
