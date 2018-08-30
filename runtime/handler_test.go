package runtime

import (
	"io/ioutil"
	"testing"

	"github.com/matryer/is"
	"github.com/sul-dlss-labs/sparql-loader/sparql"
)

func TestCorrectlyURIEncoded(t *testing.T) {
	is := is.New(t)
	content, err := ioutil.ReadFile("../fixtures/encoded_query.txt")
	is.NoErr(err)
	is.True(correctlyURIEncoded(string(content)))
}

func TestNotCorrectlyURIEncoded(t *testing.T) {
	is := is.New(t)
	content, err := ioutil.ReadFile("../fixtures/decoded_query.txt")
	is.NoErr(err)
	is.Equal(false, correctlyURIEncoded(string(content)))
}

func TestUniqueSubjects(t *testing.T) {
	is := is.New(t)
	testTriples := []sparql.Triple{
		sparql.Triple{Subject: "<http://example.com/test1>"},
		sparql.Triple{Subject: "<http://example.com/test2>"},
		sparql.Triple{Subject: "<http://example.com/test3>"},
		sparql.Triple{Subject: "<http://example.com/test1>"},
		sparql.Triple{Subject: "<http://example.com/test4>"},
	}
	actualSubjects := uniqueSubjects(testTriples)
	is.Equal(4, len(actualSubjects))
	is.Equal("http://example.com/test1", actualSubjects[0])
	is.Equal("http://example.com/test2", actualSubjects[1])
	is.Equal("http://example.com/test3", actualSubjects[2])
	is.Equal("http://example.com/test4", actualSubjects[3])
}
