package sparql

import (
	"bytes"
	"io"
	"io/ioutil"
	"strings"
	"text/scanner"
)

// Parse parses a SPARQL query from the reader
func (query *Query) Parse(src io.Reader) error {
	b, _ := ioutil.ReadAll(src)
	s := new(scanner.Scanner).Init(bytes.NewReader(b))
	s.Mode = scanner.ScanIdents | scanner.ScanStrings

	// Positional variables
	start := 0
	uriStart := 0
	prefixStart := 0
	contentStart := 0
	level := 0

	// Content Variables
	verb := ""
	content := ""
	currentSubject := ""
	triples := []Triple{}
	currentPrefix := ""

	tok := s.Scan()
	for tok != scanner.EOF {
		switch tok {
		case -2:
			if level == 0 {
				if strings.ToUpper(verb) == "PREFIX" {
					verb = "" // Clear out the verb when it is PREFIX at the top level
				} else if len(verb) > 0 {
					verb += " "
				}
				verb += s.TokenText()
			} else {
				if len(content) > 0 {
					content += " "
				}
				content += s.TokenText()
			}

		case 123: // { - starts the main content blocks
			if level == 0 {
				start = s.Position.Offset
				contentStart = start
			}
			level++

		case 125: // } - ends the main content blocks
			level--
			if level == 0 {
				if len(triples) == 0 {
					triple := query.NewTriple(string(b[contentStart+1:s.Position.Offset]), currentSubject)
					triples = append(triples, triple)
				}
				query.Parts = append(query.Parts, Sparql{
					Verb:  verb,
					Body:  string(b[start+1 : s.Position.Offset]),
					Graph: triples,
				})
				// Reset variables
				triples = []Triple{}
				verb = ""
			}

		case 46: // .
			// If inside of a content block AND NOT inside of a URI block, construct a triple for this line
			if level == 1 {
				if uriStart == 0 {
					triple := query.NewTriple(string(b[contentStart+1:s.Position.Offset]), currentSubject)
					triples = append(triples, triple)
					contentStart = s.Position.Offset
					currentSubject = ""
				}
			}

		case 59: // ; - ends a triple line
			triple := query.NewTriple(string(b[contentStart+1:s.Position.Offset]), currentSubject)
			triples = append(triples, triple)
			contentStart = s.Position.Offset
			currentSubject = triple.Subject

		case 60: // < - starts a URI
			uriStart = s.Position.Offset

		case 62: // > - ends a URI
			// If at the TOP level and inside of a PREFIX URI add to the PREFIX map
			if level == 0 {
				if prefixStart > 0 {
					query.Prefixes[currentPrefix] = string(b[uriStart+1 : s.Position.Offset])
					// Reset the prefix variables
					prefixStart--
					currentPrefix = ""
					uriStart = 0
					verb = ""
				}
			} else {
				uriStart = 0
			}

		case 58: // :
			// If at the TOP level and building a PREFIX, get the prefix code (i.e. dc:)
			if level == 0 && prefixStart == 0 {
				currentPrefix = verb
				prefixStart++
			}
		}

		tok = s.Scan()
	}

	return nil
}