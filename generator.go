package main

var (
	InitialState = state{name: "S"}
	FinalState   = state{name: "E"}
)

type Rule struct {
	lvalue string
	// groups of linked list
	groups [][]state
}

func isTitled(s string) bool {
	return len(s) > 0 && s[0] >= 'A' && s[0] <= 'Z'
}

type Nonterminals struct {
	data []string
}

// func genRules(raw []Nonterminals) []*Rule {
// 	result := make([]*Rule, 0)
// 	for _, rule := range rules {
// 		rule
// 	}
// 	return nil
// }

// genFSM generate FSM from rules. The rules must have initial and final state.
// Example:
// Rules:
//   Rule:
//     lvalue: "S"
func genFSM(rules []*Rule) *FSM {
	g := &Graph{}
	g.insert("", InitialState)
	for _, rule := range rules {
		for _, group := range rule.groups {
			if len(group) == 0 {
				continue
			}
			g.insert(rule.lvalue, group[0])
			for i := 1; i < len(group); i++ {
				g.insert(group[i-1].name, group[i])
			}
		}
	}
	return &FSM{
		graph:   g,
		current: g.Head,
	}
}
