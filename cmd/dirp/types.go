package main

import "dirp/pkg/dirp"

// InputSource は「どこから読んだ入力か」を保持する。
// File が空なら inline 入力。
type InputSource struct {
	Text string
	File string
}

type apiResponse struct {
	OK      bool           `json:"ok"`
	Mode    string         `json:"mode"`
	Source  string         `json:"source,omitempty"`
	Nodes   []jsonNode     `json:"nodes,omitempty"`
	Created bool           `json:"created,omitempty"`
	Test    bool           `json:"test,omitempty"`
	Error   *apiErrorField `json:"error,omitempty"`
	RawNodes []*dirp.Node  `json:"-"`
}

type jsonNode struct {
	Name     string     `json:"name"`
	Children []jsonNode `json:"children,omitempty"`
}

type apiErrorField struct {
	Code      string `json:"code,omitempty"`
	Category  string `json:"category,omitempty"`
	Message   string `json:"message"`
	Pos       int    `json:"pos,omitempty"`
	Line      int    `json:"line,omitempty"`
	Col       int    `json:"col,omitempty"`
	Formatted string `json:"formatted,omitempty"`
}
