package analysis

import (
	// "fmt"

	"log"
	"strings"

	"github.com/secretVal/u-lang-lsp/cmd/lsp"
)

type State struct {
	Documents map[string]string
}

func NewState() State {
	return State{
		Documents: map[string]string{},
	}
}

func (s *State) OpenDocument(doc, text string) {
	s.Documents[doc] = text
}

func (s *State) ChangeDocument(doc, text string) {
	s.Documents[doc] = text
}

func (s *State) Hover(logger *log.Logger, id int, uri string, position lsp.Position) lsp.HoverResponse {
	line := strings.Split(s.Documents[uri], "\n")[position.Line]
	word := strings.Split(line, " ")[position.Character]
	logger.Printf("word: %s", word)
	return lsp.HoverResponse{
		Response: lsp.Response{
			RPC: "2.0",
			ID:  &id,
		},
		Result: lsp.HoverResult{
			Contents: word,
		},
	}
}

// TODO: Implement Goto Definition but only when I have vars or funcs
// func (s *State) Definition(id int, uri string, position lsp.Position) lsp.DefinitionResponse {
// 	return lsp.DefinitionResponse{
// 		Response: lsp.Response{
// 			RPC: "2.0",
// 			ID:  &id,
// 		},
// 		Result: lsp.Location{
// 			URI: uri,
// 			Range: lsp.Range{
// 				Start: lsp.Position{
// 					Line:      position.Line - 1,
// 					Character: position.Character,
// 				},
// 				End: lsp.Position{
// 					Line:      position.Line - 1,
// 					Character: position.Character,
// 				},
// 			},
// 		},
// 	}
// }
