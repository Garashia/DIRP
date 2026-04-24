package dirp

import (
	"os"
	"path/filepath"
)

func BuildTree(root string, nodes []*Node) error {
	for _, n := range nodes {
		if err := buildNode(root, n); err != nil {
			return err
		}
	}
	return nil
}

func buildNode(root string, n *Node) error {
	path := filepath.Join(root, n.Name)
	if err := os.MkdirAll(path, 0o755); err != nil {
		return err
	}

	for _, child := range n.Children {
		if err := buildNode(path, child); err != nil {
			return err
		}
	}
	return nil
}
