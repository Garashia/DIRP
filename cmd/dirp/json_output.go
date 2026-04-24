package main

import (
	"dirp/pkg/dirp"
	"encoding/json"
	"os"
)

func buildJSONError(err error, input string, source string) *apiErrorField {
	out := &apiErrorField{Message: err.Error()}
	de, ok := dirp.AsError(err)
	if !ok {
		return out
	}

	out.Code = de.Code.String()
	out.Category = de.Code.Category()
	out.Message = de.Msg
	out.Pos = de.Pos
	out.Line = de.Line
	out.Col = de.Col

	if source != "" {
		withLC := err
		if input != "" && (de.Line == 0 || de.Col == 0) {
			withLC = dirp.WithLineCol(err, input)
		}
		out.Formatted = dirp.FormatError(source, withLC)
	}
	return out
}

func emitJSON(v any) {
	enc := json.NewEncoder(os.Stdout)
	enc.SetEscapeHTML(false)
	enc.SetIndent("", "  ")
	_ = enc.Encode(v)
}

func toJSONNodes(nodes []*dirp.Node) []jsonNode {
	out := make([]jsonNode, 0, len(nodes))
	for _, n := range nodes {
		out = append(out, toJSONNode(n))
	}
	return out
}

func toJSONNode(n *dirp.Node) jsonNode {
	out := jsonNode{Name: n.Name}
	if len(n.Children) > 0 {
		out.Children = toJSONNodes(n.Children)
	}
	return out
}
