package main

import (
	"dirp/pkg/dirp"
	"errors"
	"fmt"
	"os"
	"path/filepath"
)

// readInput は入力源の優先順位と相互排他を管理する。
func readInput(path string, inline string) (*InputSource, error) {
	if path != "" && inline != "" {
		return nil, errors.New("use either -f or -c, not both")
	}
	if path != "" {
		b, err := os.ReadFile(path)
		if err != nil {
			return nil, fmt.Errorf("failed to read %q: %w", path, err)
		}
		return &InputSource{Text: string(b), File: filepath.Clean(path)}, nil
	}
	if inline != "" {
		return &InputSource{Text: inline}, nil
	}
	return nil, errors.New("no input provided; use -c, -f, or -cases")
}

// srcLabel はログ出力用の入力識別子を返す。
func srcLabel(src *InputSource) string {
	if src.File != "" {
		return src.File
	}
	return src.Text
}

// parseText は lexer+parser の配線を隠蔽する小さなユーティリティ。
func parseText(input string) ([]*dirp.Node, error) {
	return dirp.Parse(input)
}
