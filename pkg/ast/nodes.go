package ast

import sitter "github.com/smacker/go-tree-sitter"

// GetClosestNodeContainsRange returns the closest node that contains the given range.
func GetClosestNodeContainsRange(node *sitter.Node, startPos uint32, endPos uint32) *sitter.Node {
	if node.StartByte() <= startPos && node.EndByte() >= endPos {
		for i := uint32(0); i < node.ChildCount(); i++ {
			child := node.Child(int(i))
			containsNode := GetClosestNodeContainsRange(child, startPos, endPos)
			if containsNode != nil {
				return containsNode
			}
		}
		return node
	}
	return nil
}
