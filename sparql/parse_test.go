package sparql

import (
	"io/ioutil"
	"log"
	"testing"

	"github.com/stretchr/testify/assert"
)

/*
 // NOTE: Leaving this old test in place until/if we remove/replace the Parse method
func TestParse(t *testing.T) {
	// var testGraph = "http://sul.stanford.edu/rialto/sources/test2"
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
				NamedGraph: nil,
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
				NamedGraph: &testGraph,
			},
		}, {
			filename: "fixtures/etl_orgs.txt",
			out: Query{
				Prefixes:   map[string]string{},
				NamedGraph: nil,
			},
		}, {
			filename: "fixtures/test.sparql",
			out: Query{
				Prefixes:   map[string]string{},
				NamedGraph: nil,
			},
		},
	}

	for _, tt := range parseTests {
		query := NewQuery()

		content, err := ioutil.ReadFile(tt.filename)
		if err != nil {
			log.Fatal(err)
		}
		err = query.Parse(string(content))
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
*/

func TestExtractEntities(t *testing.T) {
	var entitiesTests = []struct {
		filename string
		out      []string
	}{
		{
			filename: "fixtures/example1.txt",
			out:      []string{"http://example/book3", ""},
		}, {
			filename: "fixtures/example_with_graph_2.txt",
			out: []string{
				"http://sul.stanford.edu/rialto/context/positions/capFaculty_Bio-ABC_3784",
				"http://sul.stanford.edu/rialto/agents/orgs/Child_Health_Research_Institute",
				"http://sul.stanford.edu/rialto/context/positions/capFaculty_Stanford_Neurosciences_Institute_3784",
				"http://sul.stanford.edu/rialto/agents/people/3784",
				"http://sul.stanford.edu/rialto/context/names/75872"},
		}, {
			filename: "fixtures/etl_orgs.txt",
			out: []string{
				"http://sul.stanford.edu/rialto/agents/orgs/vice-provost-for-student-affairs/dean-of-educational-resources/office-of-accessible-education/oae-operations",
				"http://sul.stanford.edu/rialto/agents/orgs/vice-provost-for-student-affairs/dean-of-educational-resources/womens-community-center",
			},
		}, {
			filename: "fixtures/test.sparql",
			out: []string{
				"http://sul.stanford.edu/rialto/context/addresses/2b3d13e9-4cea-49a5-ac73-34c603461ab1_WOS:000347715900024",
				"http://sul.stanford.edu/rialto/agents/people/2b3d13e9-4cea-49a5-ac73-34c603461ab1",
				"http://sul.stanford.edu/rialto/publications/693df4b21c16dee82a08b966bb7070db",
				"http://sul.stanford.edu/rialto/context/relationships/WOS:000347715900024_e4143964-9548-4959-b04d-e434fa4219d2",
				"http://sul.stanford.edu/rialto/context/addresses/e4143964-9548-4959-b04d-e434fa4219d2_WOS:000347715900024",
				"http://sul.stanford.edu/rialto/agents/people/e4143964-9548-4959-b04d-e434fa4219d2",
				"http://sul.stanford.edu/rialto/context/relationships/WOS:000347715900024_6e99a7c0-eed4-4c4d-8f8e-ec4505c5a3a4",
				"http://sul.stanford.edu/rialto/context/addresses/6e99a7c0-eed4-4c4d-8f8e-ec4505c5a3a4_WOS:000347715900024",
				"http://sul.stanford.edu/rialto/agents/people/6e99a7c0-eed4-4c4d-8f8e-ec4505c5a3a4",
				"http://sul.stanford.edu/rialto/context/relationships/WOS:000347715900024_494ed70c-ecad-42a1-8e14-7c947eaec4ab",
				"http://sul.stanford.edu/rialto/context/addresses/494ed70c-ecad-42a1-8e14-7c947eaec4ab_WOS:000347715900024",
				"http://sul.stanford.edu/rialto/agents/people/494ed70c-ecad-42a1-8e14-7c947eaec4ab",
				"http://sul.stanford.edu/rialto/context/relationships/WOS:000347715900024_96418d53-d54d-4068-a33c-c2c875a41958",
				"http://sul.stanford.edu/rialto/context/addresses/96418d53-d54d-4068-a33c-c2c875a41958_WOS:000347715900024",
				"http://sul.stanford.edu/rialto/agents/people/96418d53-d54d-4068-a33c-c2c875a41958",
				"http://sul.stanford.edu/rialto/context/relationships/WOS:000347715900024_3974dc70-83d8-49cb-8e5e-835a94822044",
				"http://sul.stanford.edu/rialto/context/addresses/3974dc70-83d8-49cb-8e5e-835a94822044_WOS:000347715900024",
				"http://sul.stanford.edu/rialto/agents/people/3974dc70-83d8-49cb-8e5e-835a94822044",
				"http://sul.stanford.edu/rialto/context/relationships/WOS:000347715900024_7e577f2a-e2fd-4770-b961-c805ffcc5ad3",
				"http://sul.stanford.edu/rialto/context/addresses/7e577f2a-e2fd-4770-b961-c805ffcc5ad3_WOS:000347715900024",
				"http://sul.stanford.edu/rialto/agents/people/7e577f2a-e2fd-4770-b961-c805ffcc5ad3",
				"http://sul.stanford.edu/rialto/context/relationships/WOS:000347715900024_49664e8c-d613-41db-b58e-b1492ed5c7bd",
				"http://sul.stanford.edu/rialto/context/addresses/49664e8c-d613-41db-b58e-b1492ed5c7bd_WOS:000347715900024",
				"http://sul.stanford.edu/rialto/agents/people/49664e8c-d613-41db-b58e-b1492ed5c7bd",
				"http://sul.stanford.edu/rialto/context/relationships/WOS:000347715900024_35ae8118-8f72-4540-a595-caf8bc40d6ff",
				"http://sul.stanford.edu/rialto/context/addresses/35ae8118-8f72-4540-a595-caf8bc40d6ff_WOS:000347715900024",
				"http://sul.stanford.edu/rialto/agents/people/35ae8118-8f72-4540-a595-caf8bc40d6ff",
				"http://sul.stanford.edu/rialto/context/relationships/WOS:000347715900024_3854c4ab-d91a-43b1-9e66-f2cf5e756a4c",
				"http://sul.stanford.edu/rialto/context/addresses/3854c4ab-d91a-43b1-9e66-f2cf5e756a4c_WOS:000347715900024",
				"http://sul.stanford.edu/rialto/agents/people/3854c4ab-d91a-43b1-9e66-f2cf5e756a4c",
				"http://sul.stanford.edu/rialto/context/relationships/WOS:000347715900024_e4d7d217-59d3-4703-8ae9-cf525065a8fe",
				"http://sul.stanford.edu/rialto/context/addresses/e4d7d217-59d3-4703-8ae9-cf525065a8fe_WOS:000347715900024",
				"http://sul.stanford.edu/rialto/agents/people/e4d7d217-59d3-4703-8ae9-cf525065a8fe",
				"http://sul.stanford.edu/rialto/context/relationships/WOS:000347715900024_eeb64a2d-dfe8-4a43-9cc2-8c5fa1698dfd",
				"http://sul.stanford.edu/rialto/context/addresses/eeb64a2d-dfe8-4a43-9cc2-8c5fa1698dfd_WOS:000347715900024",
				"http://sul.stanford.edu/rialto/agents/people/eeb64a2d-dfe8-4a43-9cc2-8c5fa1698dfd",
				"http://sul.stanford.edu/rialto/context/relationships/WOS:000347715900024_2b319ae2-078d-468d-ae35-870ece493fbf",
				"http://sul.stanford.edu/rialto/context/addresses/2b319ae2-078d-468d-ae35-870ece493fbf_WOS:000347715900024",
				"http://sul.stanford.edu/rialto/agents/people/2b319ae2-078d-468d-ae35-870ece493fbf",
				"http://sul.stanford.edu/rialto/context/relationships/WOS:000347715900024_3e50e3c8-e55d-49b6-814d-69cc27d1e475",
				"http://sul.stanford.edu/rialto/context/addresses/3e50e3c8-e55d-49b6-814d-69cc27d1e475_WOS:000347715900024",
				"http://sul.stanford.edu/rialto/agents/people/3e50e3c8-e55d-49b6-814d-69cc27d1e475",
				"http://sul.stanford.edu/rialto/context/relationships/WOS:000347715900024_6dd9fa77-bc2f-4cb6-bd93-9e3ef02d7e99",
				"http://sul.stanford.edu/rialto/context/addresses/6dd9fa77-bc2f-4cb6-bd93-9e3ef02d7e99_WOS:000347715900024",
				"http://sul.stanford.edu/rialto/agents/people/6dd9fa77-bc2f-4cb6-bd93-9e3ef02d7e99",
			},
		},
	}

	for _, tt := range entitiesTests {
		query := NewQuery()

		content, err := ioutil.ReadFile(tt.filename)
		if err != nil {
			log.Fatal(err)
		}
		actual, err := query.ExtractEntities(string(content))
		if err != nil {
			t.Errorf("ERROR")
		} else {
			assert.Equal(t, tt.out, actual)
		}
	}
}
