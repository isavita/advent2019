package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

// Node represents an object in space.
type Node struct {
	name     string
	children []*Node
	parent   *Node
}

// findOrCreateNode finds or creates a node with the given name.
func findOrCreateNode(name string, nodes map[string]*Node) *Node {
	if node, ok := nodes[name]; ok {
		return node
	}
	node := &Node{name: name}
	nodes[name] = node
	return node
}

// buildOrbitMap builds the orbit map from the input data.
func buildOrbitMap(scanner *bufio.Scanner) map[string]*Node {
	nodes := make(map[string]*Node)
	for scanner.Scan() {
		parts := strings.Split(scanner.Text(), ")")
		center := findOrCreateNode(parts[0], nodes)
		orbiter := findOrCreateNode(parts[1], nodes)
		center.children = append(center.children, orbiter)
		orbiter.parent = center
	}
	return nodes
}

// pathToRoot creates a slice of nodes from a given node to the root.
func pathToRoot(node *Node) []*Node {
	var path []*Node
	for node != nil {
		path = append(path, node)
		node = node.parent
	}
	return path
}

// findCommonAncestor finds the common ancestor of two nodes and returns the paths to it.
func findCommonAncestor(node1, node2 *Node) (int, int) {
	path1 := pathToRoot(node1)
	path2 := pathToRoot(node2)

	i := len(path1) - 1
	j := len(path2) - 1

	// Move upwards until the paths diverge.
	for i >= 0 && j >= 0 && path1[i] == path2[j] {
		i--
		j--
	}
	return i + 1, j + 1
}

func main() {
	file, err := os.Open("day6/input.txt")
	if err != nil {
		fmt.Println("Error opening file:", err)
		return
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	orbitMap := buildOrbitMap(scanner)

	// Calculate the number of orbital transfers.
	transfersYOU, transfersSAN := findCommonAncestor(orbitMap["YOU"].parent, orbitMap["SAN"].parent)
	fmt.Println(transfersYOU + transfersSAN)
}
