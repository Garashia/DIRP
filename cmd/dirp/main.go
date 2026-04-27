package main

import (
	"flag"
	"fmt"
	"os"
)

// main は 3 モードを束ねる:
// 1) 単発 parse/run (-c or -f)
// 2) 一括テスト (-cases)
// 3) 実ディレクトリ作成 (-mkdir, ただし -test が優先)
func main() {
	input := flag.String("c", "", "dirp command string")
	file := flag.String("f", "", "path to a .dirp file")
	casesFile := flag.String("cases", "", "path to test cases file (one DSL pattern per line)")
	jsonOut := flag.Bool("json", false, "output AST / errors as JSON (for API integration)")
	graphMode := flag.Bool("Graph", false, "print directory graph view in ASCII")
	root := flag.String("root", ".", "output root directory")
	makeDirs := flag.Bool("mkdir", false, "create directories on filesystem")
	testMode := flag.Bool("test", false, "debug mode: parse and print only (no directory creation)")
	flag.Parse()

	if *casesFile != "" {
		// バッチテストは副作用なし（作成処理なし）。
		if err := runBatchCases(*casesFile); err != nil {
			if *jsonOut {
				emitJSON(apiResponse{
					OK:    false,
					Mode:  "cases",
					Error: buildJSONError(err, "", "cases"),
				})
				os.Exit(1)
			}
			fmt.Fprintln(os.Stderr, "error:", err)
			os.Exit(1)
		}
		if *jsonOut {
			emitJSON(apiResponse{OK: true, Mode: "cases"})
		}
		return
	}

	resp, err := runSingle(*input, *file, *root, *makeDirs, *testMode)
	if *jsonOut {
		emitJSON(resp)
		if err != nil {
			os.Exit(1)
		}
		return
	}

	if err != nil {
		if resp != nil && resp.Error != nil && resp.Error.Formatted != "" {
			fmt.Fprintln(os.Stderr, resp.Error.Formatted)
		} else {
			fmt.Fprintln(os.Stderr, "error:", err)
		}
		os.Exit(1)
	}

	if *graphMode {
		fmt.Printf("Graph for: %s\n", resp.Source)
		printGraphNodes(resp.RawNodes)
	} else {
		fmt.Printf("Parsed structure for: %s\n", resp.Source)
		for _, n := range resp.RawNodes {
			printTree(n, 0)
		}
	}
	if *testMode && *makeDirs {
		// デバッグ安全性を優先して -test を勝たせる。
		fmt.Println("warning: both -test and -mkdir were set; running in test mode (no filesystem changes)")
	}
	if resp.Created {
		fmt.Printf("Directories created under: %s\n", *root)
	} else {
		fmt.Println("Test mode: no directories were created.")
	}
}
