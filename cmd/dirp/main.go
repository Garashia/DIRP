package main

import (
	"dirp/pkg/dirp"
	"errors"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"unicode/utf8"
)

// main は 3 モードを束ねる:
// 1) 単発 parse/run (-c or -f)
// 2) 一括テスト (-cases)
// 3) 実ディレクトリ作成 (-mkdir, ただし -test が優先)
func main() {
	input := flag.String("c", "", "dirp command string")
	file := flag.String("f", "", "path to a .dirp file")
	casesFile := flag.String("cases", "", "path to test cases file (one DSL pattern per line)")
	root := flag.String("root", ".", "output root directory")
	makeDirs := flag.Bool("mkdir", false, "create directories on filesystem")
	testMode := flag.Bool("test", false, "debug mode: parse and print only (no directory creation)")
	flag.Parse()

	if *casesFile != "" {
		// バッチテストは副作用なし（作成処理なし）。
		if err := runBatchCases(*casesFile); err != nil {
			fmt.Fprintln(os.Stderr, "error:", err)
			os.Exit(1)
		}
		return
	}

	src, err := readInput(*file, *input)
	if err != nil {
		fmt.Fprintln(os.Stderr, "error:", err)
		os.Exit(1)
	}
	if !utf8.ValidString(src.Text) {
		fmt.Fprintln(os.Stderr, "error: input is not valid UTF-8")
		os.Exit(1)
	}

	nodes, err := parseText(src.Text)
	if err != nil {
		// parse 層の Pos を line/col へ変換して表示を統一する。
		err = dirp.WithLineCol(err, src.Text)
		label := "inline"
		if src.File != "" {
			label = src.File
		}
		fmt.Fprintln(os.Stderr, dirp.FormatError(label, err))
		os.Exit(1)
	}

	fmt.Printf("Parsed structure for: %s\n", srcLabel(src))
	for _, n := range nodes {
		printTree(n, 0)
	}

	if *testMode && *makeDirs {
		// デバッグ安全性を優先して -test を勝たせる。
		fmt.Println("warning: both -test and -mkdir were set; running in test mode (no filesystem changes)")
	}

	if *makeDirs && !*testMode {
		if err := dirp.BuildTree(*root, nodes); err != nil {
			fmt.Printf("build error: %v\n", err)
			os.Exit(1)
		}
		fmt.Printf("Directories created under: %s\n", *root)
	} else {
		fmt.Println("Test mode: no directories were created.")
	}
}

// InputSource は「どこから読んだ入力か」を保持する。
// File が空なら inline 入力。
type InputSource struct {
	Text string
	File string
}

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
	l := dirp.NewLexer(input)
	p := dirp.NewParser(l)
	return p.Parse()
}

// runBatchCases は 1 行 1 パターンで DSL を連続検証する。
// 空行と '#...' はコメントとして無視する。
func runBatchCases(path string) error {
	b, err := os.ReadFile(path)
	if err != nil {
		return fmt.Errorf("failed to read cases file %q: %w", path, err)
	}
	if !utf8.Valid(b) {
		return errors.New("cases file is not valid UTF-8")
	}

	lines := strings.Split(string(b), "\n")
	total := 0
	passed := 0
	cleanPath := filepath.Clean(path)

	for i, line := range lines {
		pattern := strings.TrimSpace(strings.TrimRight(line, "\r"))
		if pattern == "" || strings.HasPrefix(pattern, "#") {
			continue
		}
		total++
		fmt.Printf("\n[CASE %d] %s\n", i+1, pattern)

		nodes, parseErr := parseText(pattern)
		if parseErr != nil {
			// cases ファイル内の "行番号" と、パターン内の "列番号" を合成して表示。
			parseErr = dirp.WithLineCol(parseErr, pattern)
			fmt.Printf("NG  %s\n", formatBatchCaseError(cleanPath, i+1, parseErr))
			continue
		}

		fmt.Println("OK")
		for _, n := range nodes {
			printTree(n, 0)
		}
		passed++
	}

	if total == 0 {
		return errors.New("no test cases found (empty lines and lines starting with '#' are ignored)")
	}

	fmt.Printf("\nResult: %d/%d cases passed\n", passed, total)
	if passed != total {
		return errors.New("some test cases failed")
	}
	return nil
}

// formatBatchCaseError は "casesファイルの行:列" 形式へ整形する。
func formatBatchCaseError(path string, caseLine int, err error) string {
	de, ok := dirp.AsError(err)
	if !ok {
		return fmt.Sprintf("%s:%d: parse error: %s", path, caseLine, err.Error())
	}

	col := de.Col
	if col < 1 {
		col = 1
	}
	msg := de.Msg
	if de.Cause != nil {
		msg = msg + ": " + de.Cause.Error()
	}
	return fmt.Sprintf("%s:%d:%d: %s error: %s", path, caseLine, col, de.Code.Category(), msg)
}

// printTree は parse 結果をデバッグしやすい形で可視化する。
func printTree(n *dirp.Node, indent int) {
	fmt.Printf("%s|-- %s/\n", strings.Repeat("  ", indent), n.Name)
	for _, child := range n.Children {
		printTree(child, indent+1)
	}
}
