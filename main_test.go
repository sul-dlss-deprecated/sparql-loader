package main

import (
	"io/ioutil"
	"log"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEncodeBody(t *testing.T) {
	var encodeTests = []struct {
		filename string
		result   bool
	}{
		{
			filename: "fixtures/decoded_query.txt",
			result:   true,
		}, {
			filename: "fixtures/encoded_query.txt",
			result:   false,
		},
	}

	for _, tt := range encodeTests {
		content, err := ioutil.ReadFile(tt.filename)
		if err != nil {
			log.Fatal(err)
		}
		assert.Equal(t, tt.result, unEscaped(string(content)))
	}
}
