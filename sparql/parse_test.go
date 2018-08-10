package sparql

import (
	"bytes"
	"io/ioutil"
	"log"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
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
				Prefixes:   map[string]string{},
				NamedGraph: "",
			},
		}, {
			filename: "fixtures/example_with_graph_2.txt",
			out: Query{
				Parts: []Sparql{{
					Verb: "INSERT",
					Graph: []Triple{{
						Subject:   "<http://sul.stanford.edu/rialto/context/positions/capFaculty_Bio-ABC_3784>",
						Predicate: "<http://vivoweb.org/ontology/core#hrJobTitle>",
						Object:    "\"George A. Zimmermann Professor and Professor of Pediatrics\"",
					}, {
						Subject:   "<http://sul.stanford.edu/rialto/agents/orgs/Child_Health_Research_Institute>",
						Predicate: "<http://vivoweb.org/ontology/core#relatedBy>",
						Object:    "<http://sul.stanford.edu/rialto/context/positions/capFaculty_Child_Health_Research_Institute_3784>",
					}, {
						Subject:   "<http://sul.stanford.edu/rialto/context/positions/capFaculty_Stanford_Neurosciences_Institute_3784>",
						Predicate: "<http://vivoweb.org/ontology/core#relates>",
						Object:    "<http://sul.stanford.edu/rialto/agents/people/3784>",
					}, {
						Subject:   "<http://sul.stanford.edu/rialto/agents/people/3784>",
						Predicate: "<http://vivoweb.org/ontology/core#relatedBy>",
						Object:    "<http://sul.stanford.edu/rialto/context/relationships/75872_3784>",
					}, {
						Subject:   "<http://sul.stanford.edu/rialto/context/names/75872>",
						Predicate: "<http://www.w3.org/2006/vcard/ns#given-name>",
						Object:    "\"Noga\"",
					}},
					Body: `<http://sul.stanford.edu/rialto/context/positions/capFaculty_Bio-ABC_3784><http://vivoweb.org/ontology/core#hrJobTitle>"GeorgeA.ZimmermannProfessorandProfessorofPediatrics".<http://sul.stanford.edu/rialto/agents/orgs/Child_Health_Research_Institute><http://vivoweb.org/ontology/core#relatedBy><http://sul.stanford.edu/rialto/context/positions/capFaculty_Child_Health_Research_Institute_3784>.<http://sul.stanford.edu/rialto/context/positions/capFaculty_Stanford_Neurosciences_Institute_3784><http://vivoweb.org/ontology/core#relates><http://sul.stanford.edu/rialto/agents/people/3784>.<http://sul.stanford.edu/rialto/agents/people/3784><http://vivoweb.org/ontology/core#relatedBy><http://sul.stanford.edu/rialto/context/relationships/75872_3784>.<http://sul.stanford.edu/rialto/context/names/75872><http://www.w3.org/2006/vcard/ns#given-name>"Noga".`,
				}},
				Prefixes:   map[string]string{},
				NamedGraph: "http://sul.stanford.edu/rialto/sources/test2",
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
			assert.Equal(t, tt.out, query)
		}
	}
}
