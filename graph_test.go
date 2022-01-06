package main

// import (
// 	"testing"
// )

// func TestGraph(t *testing.T) {
// 	// g := Graph{}
// 	// g.insert("", state{name: "S"})
// 	// g.insert("S", state{name: "A"})
// 	// g.insert("A", state{name: "B"})
// 	// g.insert("B", state{name: "S"})
// 	// g.print()

// 	g := Graph{}
// 	g.insert("", state{name: "lvalue"})
// 	g.insert("lvalue", state{name: "\"=\""})
// 	g.insert("\"=\"", state{name: "expr"})
// 	g.insert("expr", state{name: "rvalue"})
// 	g.insert("rvalue", state{name: "\"|\""})
// 	g.insert("\"|\"", state{name: "rvalue"})
// 	g.insert("rvalue", state{name: "entity"})
// 	g.insert("rvalue", state{name: "string"})
// 	g.insert("rvalue", state{name: "E"})
// 	g.insert("rvalue", state{name: "["})
// 	g.insert("[", state{name: "expr"})
// 	g.insert("expr", state{name: "]"})
// 	g.insert("]", state{name: "E"})
// 	g.insert("string", state{name: "E"})
// 	g.insert("entity", state{name: "E"})
// 	g.print()
// }
