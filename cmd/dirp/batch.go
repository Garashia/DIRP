package main

import (
	"dirp/pkg/dirp"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
)

// runBatchCases は 1 行 1 パターンで DSL を連続検証する。
// 空行と '#...' はコメントとして無視する。
func runBatchCases(path string) error {
	return runBatchCasesTo(path, os.Stdout)
}

// runBatchCasesTo は runBatchCases の出力先注入版。
// テスト時に bytes.Buffer を渡せるようにし、stdout 依存を減らす。
func runBatchCasesTo(path string, out io.Writer) error {
	b, err := os.ReadFile(path)
	if err != nil {
		return fmt.Errorf("failed to read cases file %q: %w", path, err)
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
		fmt.Fprintf(out, "\n[CASE %d] %s\n", i+1, pattern)

		nodes, parseErr := parseText(pattern)
		if parseErr != nil {
			// cases ファイル内の "行番号" と、パターン内の "列番号" を合成して表示。
			parseErr = dirp.WithLineCol(parseErr, pattern)
			fmt.Fprintf(out, "NG  %s\n", formatBatchCaseError(cleanPath, i+1, parseErr))
			continue
		}

		fmt.Fprintln(out, "OK")
		for _, n := range nodes {
			printTreeTo(out, n, 0)
		}
		passed++
	}

	if total == 0 {
		return errors.New("no test cases found (empty lines and lines starting with '#' are ignored)")
	}

	fmt.Fprintf(out, "\nResult: %d/%d cases passed\n", passed, total)
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
