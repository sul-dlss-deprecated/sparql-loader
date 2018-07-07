package sparql

import (
	"bytes"
	"io/ioutil"
	"strings"
	"text/scanner"
)

// Triple is exported in order to be used when parsing
type Triple struct {
	Subject   string
	Predicate string
	Object    string
}

// NewTriple takes a string and default subject
func (query *Query) NewTriple(line string, defaultSubject string) Triple {
	byt, _ := ioutil.ReadAll(strings.NewReader(line))
	scan := new(scanner.Scanner).Init(bytes.NewReader(byt))
	scan.Mode = scanner.ScanIdents | scanner.ScanStrings

	objectLiteral := ""

	token := scan.Scan()
	for token != scanner.EOF {
		switch token {

		case -6:
			objectLiteral = scan.TokenText()

		default:
			if len(objectLiteral) > 0 {
				objectLiteral += scan.TokenText()
			}
		}
		token = scan.Scan()
	}

	line = strings.Replace(line, objectLiteral, "", -1)
	tripleArray := strings.Fields(line)
	if len(objectLiteral) > 0 {
		tripleArray = append(tripleArray, objectLiteral)
	}
	if len(tripleArray) != 3 {
		return Triple{
			Subject:   defaultSubject,
			Predicate: query.replacePrefix(tripleArray[0]),
			Object:    query.replacePrefix(tripleArray[1])}
	}
	return Triple{
		Subject:   query.replacePrefix(tripleArray[0]),
		Predicate: query.replacePrefix(tripleArray[1]),
		Object:    query.replacePrefix(tripleArray[2])}

}

func (query *Query) replacePrefix(input string) string {
	i := strings.Index(input, ":")
	if i > -1 && len(query.Prefixes[input[:i]]) > 0 {
		return query.Prefixes[input[:i]] + input[i+1:]
	}
	return input

}
