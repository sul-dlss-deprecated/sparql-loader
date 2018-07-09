package sparql

import (
	"bytes"
	"io/ioutil"
	"log"
	"reflect"
	"strings"
	"testing"
)

func TestParse(t *testing.T) {
	var parseTests = []struct {
		filename string
		out      Query
	}{
		{
			filename: "fixtures/example1.txt",
			out: Query{
				Parts: []Sparql{{
					Verb: "INSERT DATA",
					Graph: []Triple{{
						Subject:   "<http://example/book3>",
						Predicate: "<http://purl.org/dc/elements/1.1/title>",
						Object:    "\"A new book\"",
					}, {
						Subject:   "<http://example/book3>",
						Predicate: "<http://purl.org/dc/elements/1.1/creator>",
						Object:    "\"A.N.Other\"",
					}},
					Body: `<http://example/book3><http://purl.org/dc/elements/1.1/title>"Anewbook";<http://purl.org/dc/elements/1.1/creator>"A.N.Other".`,
				}},
				Prefixes: map[string]string{},
			},
		},
	}

	for _, tt := range parseTests {
		query := NewQuery()

		content, err := ioutil.ReadFile(tt.filename)
		if err != nil {
			log.Fatal(err)
		}
		err = query.Parse(bytes.NewReader(content))
		if err != nil {
			t.Errorf("ERROR")
		} else {
			for i, qq := range query.Parts {
				query.Parts[i].Body = strings.Replace(strings.Replace(qq.Body, " ", "", -1), "\n", "", -1)
			}
			if !reflect.DeepEqual(query, tt.out) {
				t.Errorf("ERROR")
			}
		}
	}
}
