package runtime

import (
	"io/ioutil"
	"testing"

	"github.com/matryer/is"
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
