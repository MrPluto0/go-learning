package gee

import "strings"

type node struct {
	pattern  string  // route matched, such as /p/:lang
	part     string  // a part of route, such as :lang
	children []*node // subNotes, such as [doc, tutorial, intro]
	isWild   bool    // fuzzy matching. True if part includes : or *
}

// match the first node
func (n *node) matchChild(part string) *node {
	for _, child := range n.children {
		if child.part == part || child.isWild {
			return child
		}
	}
	return nil
}

// match all nodes
func (n *node) matchChildren(part string) []*node {
	nodes := make([]*node, 0)
	for _, child := range n.children {
		if child.part == part || child.isWild {
			nodes = append(nodes, child)
		}
	}
	return nodes
}

// height = 0 ~ len(part) - 1
func (n *node) insert(pattern string, parts []string, height int) {
	if len(parts) == height {
		n.pattern = pattern
		return
	}

	part := parts[height]
	child := n.matchChild(part)
	if child == nil {
		child = &node{part: part, isWild: part[0] == '*' || part[0] == ':'}
		n.children = append(n.children, child)
	}
	child.insert(pattern, parts, height+1)
}

func (n *node) search(parts []string, height int) *node {
	if len(parts) == height || strings.HasPrefix(n.part, "*") {
		if n.pattern == "" {
			return nil
		}
		return n
	}

	part := parts[height]
	children := n.matchChildren(part)

	for _, child := range children {
		node := child.search(parts, height+1)
		if node != nil {
			return node
		}
	}

	return nil
}
