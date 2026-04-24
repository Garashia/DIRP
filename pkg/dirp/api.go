package dirp

import (
	"fmt"
	"os"
	"path/filepath"
	"unicode/utf8"
)

// Parse は外部利用向けの最小エントリポイント。
// 文字列 DSL を受け取り、ディレクトリ木(Node配列)を返す。
func Parse(input string) ([]*Node, error) {
	if !utf8.ValidString(input) {
		return nil, NewError(ErrInvalidArgument, 0, "input is not valid UTF-8")
	}
	l := NewLexer(input)
	p := NewParser(l)
	return p.Parse()
}

// ParseFile は DSL ファイルを読み込んで Parse するヘルパー。
// 呼び出し側で path を保持しておくと、FormatError と組み合わせて表示しやすい。
func ParseFile(path string) ([]*Node, string, error) {
	clean := filepath.Clean(path)
	b, err := os.ReadFile(clean)
	if err != nil {
		return nil, "", fmt.Errorf("failed to read %q: %w", clean, err)
	}
	text := string(b)
	nodes, parseErr := Parse(text)
	if parseErr != nil {
		return nil, text, parseErr
	}
	return nodes, text, nil
}

// Build は外部利用向けの作成 API。
// Parse 結果を受け取り、実際にディレクトリを作成する。
func Build(root string, nodes []*Node) error {
	return BuildTree(root, nodes)
}
