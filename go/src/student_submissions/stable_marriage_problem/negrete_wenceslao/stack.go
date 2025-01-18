// Package main implements a stable marriage problem solver
package main

import (
	"fmt"
	"strings"
)

// node represents a single element in the stack
type node struct {
	name string
	next *node
}

// stack implements a LIFO (Last In First Out) data structure
type stack struct {
	top    *node
	length int
}

// NewStack creates and returns a new empty stack
func NewStack() *stack {
	return &stack{
		top:    nil,
		length: 0,
	}
}

// Clone creates a deep copy of the stack
func (s *stack) Clone() *stack {
	newStack := NewStack()
	if s.top == nil {
		return newStack
	}

	// Create a map of old nodes to new nodes
	nodeMap := make(map[*node]*node)
	
	// First pass: create all new nodes
	current := s.top
	for current != nil {
		nodeMap[current] = &node{name: current.name}
		current = current.next
	}
	
	// Second pass: connect all the nodes
	current = s.top
	for current != nil {
		if current.next != nil {
			nodeMap[current].next = nodeMap[current.next]
		}
		current = current.next
	}
	
	newStack.top = nodeMap[s.top]
	newStack.length = s.length
	return newStack
}

// Push adds a new element to the top of the stack
func (s *stack) push(name string) {
	n := &node{
		name: name,
		next: s.top,
	}
	s.top = n
	s.length++
}

// Pop removes and returns the top element from the stack
func (s *stack) pop() (string, error) {
	if s.isEmpty() {
		return "", fmt.Errorf("stack is empty")
	}
	n := s.top
	s.top = n.next
	s.length--
	return n.name, nil
}

// Peek returns the top element without removing it
func (s *stack) peek() (string, error) {
	if s.isEmpty() {
		return "", fmt.Errorf("stack is empty")
	}
	return s.top.name, nil
}

// IsEmpty returns true if the stack has no elements
func (s *stack) isEmpty() bool {
	return s.top == nil
}

// Size returns the current number of elements in the stack
func (s *stack) size() int {
	return s.length
}

// Print returns a string representation of the stack
func (s *stack) print() string {
	if s.isEmpty() {
		return "empty"
	}
	
	var sb strings.Builder
	current := s.top
	for current != nil {
		sb.WriteString(current.name)
		if current.next != nil {
			sb.WriteString(" -> ")
		}
		current = current.next
	}
	return sb.String()
}
