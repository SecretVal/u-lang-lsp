package main

import (
	"bufio"
	"encoding/json"
	"log"
	"os"

	"github.com/secretVal/u-lang-lsp/cmd/analysis"
	"github.com/secretVal/u-lang-lsp/cmd/lsp"
	"github.com/secretVal/u-lang-lsp/cmd/rpc"
)

func main() {
	logger := getLogger("/home/lukas/Dokumente/coding/u-lang/u-lang-lsp/log.txt")
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Split(rpc.Split)

	state := analysis.NewState()

	for scanner.Scan() {
		msg := scanner.Bytes()
		method, contets, err := rpc.DecodeMessage(msg)
		if err != nil {
			logger.Printf("Got an Error: %s", err)
			continue
		}
		handleMessage(logger, state, method, contets)
	}
}

func handleMessage(logger *log.Logger, state analysis.State, method string, contents []byte) {
	logger.Printf("Recieved msg with method: %s", method)
	switch method {
	case "initialize":
		var req lsp.InitializeRequest
		if err := json.Unmarshal(contents, &req); err != nil {
			logger.Printf("Couldn't parse this: %s", err)
		}
		logger.Printf("Conected to %s %s", req.Params.ClientInfo.Name, req.Params.ClientInfo.Version)
		writeResponse(lsp.NewInitializeResponse(req.ID))

	case "textDocument/didOpen":
		var notification lsp.DidOpenTextDocumentNotification
		if err := json.Unmarshal(contents, &notification); err != nil {
			logger.Printf("Couldn't parse this: %s", err)
		}
		logger.Printf("Opened %s", notification.Params.TextDocument.URI)
		state.OpenDocument(notification.Params.TextDocument.URI, notification.Params.TextDocument.Text)

	case "textDocument/didChange":
		var notification lsp.DidChangeTextDocumentNotification
		if err := json.Unmarshal(contents, &notification); err != nil {
			logger.Printf("Couldn't parse this: %s", err)
		}
		logger.Printf("Changed %s", notification.Params.TextDocument.URI)
		for _, change := range notification.Params.ContentChanges {
			state.ChangeDocument(notification.Params.TextDocument.URI, change.Text)
		}
	case "textDocument/hover":
		var req lsp.HoverRequest
		if err := json.Unmarshal(contents, &req); err != nil {
			logger.Printf("Couldn't parse this: %s", err)
		}
		res := state.Hover(logger, req.ID, req.Params.TextDocument.URI, req.Params.Position)
		log.Print(res)
		writeResponse(res)
		// TODO: Implement Goto Definition but only when I have vars or funcs
		// case "textDocument/definition":
		// 	var req lsp.DefinitionRequest
		// 	if err := json.Unmarshal(contents, &req); err != nil {
		// 		logger.Printf("Couldn't parse this: %s", err)
		// 	}
		// 	res := state.Definition(req.ID, req.Params.TextDocument.URI, req.Params.Position)
		// 	writeResponse(res)
	}
}

func getLogger(filename string) *log.Logger {
	logfile, err := os.OpenFile(filename, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0666)
	if err != nil {
		panic("not a good logfile")
	}
	return log.New(logfile, "[u-lang-lsp]", log.Ldate|log.Ltime|log.Lshortfile)
}

func writeResponse(msg any) {
	reply := rpc.EncodeMessage(msg)
	os.Stdout.Write([]byte(reply))
}
