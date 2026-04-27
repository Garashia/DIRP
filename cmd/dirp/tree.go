package main

import (
	"dirp/pkg/dirp"
	"fmt"
	"io"
	"os"
	"strings"
)

// printTree は parse 結果をデバッグしやすい形で可視化する。
func printTree(n *dirp.Node, indent int) {
	printTreeTo(os.Stdout, n, indent)
}

func printTreeTo(w io.Writer, n *dirp.Node, indent int) {
	fmt.Fprintf(w, "%s|-- %s/\n", strings.Repeat("  ", indent), n.Name)
	for _, child := range n.Children {
		printTreeTo(w, child, indent+1)
	}
}

// printGraphNodes は複数ルートノードを ASCII の構成図として表示する。
func printGraphNodes(nodes []*dirp.Node) {
	printGraphNodesTo(os.Stdout, nodes)
}

func printGraphNodesTo(w io.Writer, nodes []*dirp.Node) {
	for i, n := range nodes {
		isLast := i == len(nodes)-1
		printGraphNodeTo(w, n, "", isLast)
	}
}

func printGraphNodeTo(w io.Writer, n *dirp.Node, prefix string, isLast bool) {
	connector := "|-- "
	nextPrefix := prefix + "|   "
	if isLast {
		connector = "`-- "
		nextPrefix = prefix + "    "
	}
	fmt.Fprintf(w, "%s%s%s/\n", prefix, connector, n.Name)
	for i, child := range n.Children {
		printGraphNodeTo(w, child, nextPrefix, i == len(n.Children)-1)
	}
}
