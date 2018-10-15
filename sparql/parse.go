package sparql

import (
	"strings"
)

// ExtractEntities parses a SPARQL query from the reader and returns an array of entities (subjects)
func (query *Query) ExtractEntities(src string) ([]string, error) {
	entities := []string{}
	blocks := strings.Split(src, "}")
	for _, block := range blocks {
		blockParts := strings.Split(block, "{") // Get the VERB of the query and any sub sections (i.e. GRAPHS)
		if len(blockParts) > 0 {
			tripleBlock := strings.TrimSpace(blockParts[len(blockParts)-1]) // By default we grab the last part of the block.
			if len(tripleBlock) > 1 {                                       // We have something more than ".", ";", or blank
				// Every triple line should end with either "." or ";"
				startPos := 0
				pos := strings.Index(tripleBlock, " ;")
				for pos > startPos {
					aTriple := query.NewTriple(tripleBlock[startPos:pos], "")
					entities = appendEntity(entities, aTriple.Subject)
					startPos = pos + 2
					pos = strings.Index(tripleBlock[startPos:], " ;") + startPos
				}

				pos = strings.Index(string(tripleBlock[startPos:]), " .") + startPos // We need to add our currert start position as index starts from 0
				// This check must happen in order to determine if there is literal block
				// that may match our end-of-statement marker.
				if strings.Count(string(tripleBlock[startPos:pos]), "\"")%2 != 0 {
					pos = strings.Index(string(tripleBlock[pos+2]), " .")
				}
				for pos > startPos {
					aTriple := query.NewTriple(tripleBlock[startPos:pos], "")
					entities = appendEntity(entities, aTriple.Subject)
					startPos = pos + 2
					pos = strings.Index(tripleBlock[startPos:], " .") + startPos // We need to add our currert start position as index starts from 0
				}

			}
		}
	}
	return entities, nil
}

func appendEntity(entities []string, newEntity string) []string {
	// Strip URI wrapper from entity
	newEntity = strings.Replace(newEntity, "<", "", -1)
	newEntity = strings.Replace(newEntity, ">", "", -1)

	for i := range entities {
		if entities[i] == newEntity {
			return entities
		}
	}
	return append(entities, newEntity)
}
