package sparql

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewGraph(t *testing.T) {
	var graphTests = []struct {
		in  string
		out Graph
	}{
		{
			in: "GRAPH <http://sul.stanford.edu/rialto/sources/test2>",
			out: Graph{
				URI: "<http://sul.stanford.edu/rialto/sources/test2>",
			},
		},
	}

	query := Query{
		Prefixes: map[string]string{
			"dc": "http://purl.org/dc/elements/1.1/",
		},
	}

	for _, tt := range graphTests {
		actual := query.NewGraph(tt.in)
		assert.Equal(t, tt.out, actual)
	}
}
