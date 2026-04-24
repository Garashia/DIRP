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
