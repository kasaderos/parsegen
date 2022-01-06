package main

// import "fmt"

// type Node struct {
// 	state
// 	// todo do as linked list
// 	nexts []*Node
// }

// func (nd *Node) Marked() bool {
// 	return isTitled(nd.state.name)
// }

// type Graph struct {
// 	Head *Node
// }

// // key is a state name where data is appended to the nexts array
// func (g *Graph) insert(key string, data state) {
// 	if g.Head == nil {
// 		g.Head = &Node{state: data}
// 		return
// 	}
// 	stack := Stack{}
// 	states := make(map[string]bool, 0)
// 	stack.Push(Frame{node: g.Head})
// 	states[g.Head.name] = true
// 	for !stack.Empty() {
// 		front := stack.Top()
// 		stack.Pop()
// 		// push back
// 		if front.node.name == key {
// 			front.node.nexts = append(front.node.nexts, &Node{state: data})
// 			return
// 		}
// 		// visited
// 		states[front.node.name] = true

// 		for _, next := range front.node.nexts {
// 			// not visited nodes
// 			if _, ok := states[next.name]; !ok {
// 				stack.Push(Frame{node: next})
// 			}
// 		}
// 	}
// }

// func (g *Graph) print() {
// 	stack := Stack{}
// 	states := make(map[string]bool, 0)
// 	stack.Push(Frame{node: g.Head})
// 	states[g.Head.name] = true

// 	for !stack.Empty() {
// 		front := stack.Top()
// 		stack.Pop()
// 		// visited
// 		states[front.node.name] = true

// 		fmt.Printf("%s -> ", front.node.name)
// 		for _, next := range front.node.nexts {
// 			// not visited nodes
// 			if _, ok := states[next.name]; !ok {
// 				stack.Push(Frame{node: next})
// 			}
// 			fmt.Printf("%s, ", next.name)
// 		}
// 		fmt.Printf("\n")
// 	}
// }
