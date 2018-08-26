package main

import (
	"io/ioutil"
	"log"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEncodeBody(t *testing.T) {
	var encodeTests = []struct {
		filename string
		out      *strings.Reader
	}{
		{
			filename: "fixtures/decoded_query.txt",
			out:      strings.NewReader("update=INSERT+DATA+%7B+%3Chttp%3A%2F%2Fexample.org%2Fpeople%2F1234%3E+%3Chttp%3A%2F%2Fpurl.example.org%2Fme%2FEX_0001%3E+%3Chttp%3A%2F%2Fexample.org%2Froles%2FPrincipalInvestigator_%26_Agent%3E+.+%7D"),
		}, {
			filename: "fixtures/encoded_query.txt",
			out:      strings.NewReader("query=SELECT%20%7B%20%3Fs%20%3Fp%20%3Fo%20%7D%20WHERE%20%7B%20%3Fs%20%3Fp%20%3Fo%20%7D"),
		},
	}

	for _, tt := range encodeTests {
		content, err := ioutil.ReadFile(tt.filename)
		if err != nil {
			log.Fatal(err)
		}
		assert.Equal(t, tt.out, encodeBody(string(content)))
	}
}
