package main

import (
	"io/ioutil"
	"log"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCorrectlyURIEncoded(t *testing.T) {
	content, err := ioutil.ReadFile("fixtures/encoded_query.txt")
	if err != nil {
		log.Fatal(err)
	}
	assert.True(t, correctlyURIEncoded(string(content)))
}

func TestNotCorrectlyURIEncoded(t *testing.T) {
	content, err := ioutil.ReadFile("fixtures/decoded_query.txt")
	if err != nil {
		log.Fatal(err)
	}
	assert.False(t, correctlyURIEncoded(string(content)))
}
