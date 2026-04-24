package main

import "dirp/pkg/dirp"

// runSingle は単発入力の Parse/Build を実行し、結果を apiResponse で返す。
// JSON 出力とCLI表示のどちらでも再利用できるよう、副作用の判断をここへ集約する。
func runSingle(input string, file string, root string, makeDirs bool, testMode bool) (*apiResponse, error) {
	resp := &apiResponse{OK: false, Mode: "single"}

	src, err := readInput(file, input)
	if err != nil {
		resp.Error = buildJSONError(err, "", "input")
		return resp, err
	}

	label := "inline"
	if src.File != "" {
		label = src.File
	}
	resp.Source = label

	nodes, err := parseText(src.Text)
	if err != nil {
		err = dirp.WithLineCol(err, src.Text)
		resp.Error = buildJSONError(err, src.Text, label)
		return resp, err
	}

	resp.RawNodes = nodes
	resp.Nodes = toJSONNodes(nodes)
	resp.Test = testMode

	if makeDirs && !testMode {
		if err := dirp.Build(root, nodes); err != nil {
			resp.Error = buildJSONError(err, "", "build")
			return resp, err
		}
		resp.Created = true
	}

	resp.OK = true
	return resp, nil
}
