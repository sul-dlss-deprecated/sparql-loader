package sparql

import (
	"strings"
)

type Graph struct {
	URI string
}

func (query *Query) NewGraph(line string) Graph {
	parts := strings.Fields(line)
	if len(parts) != 2 {
		return Graph{}
	}
	return Graph{URI: parts[1]}
}
