package main

func any(p *Parser, end byte) func() bool {
	return func() bool {
		for p.cc() != end {
			p.gc()
		}
		p.gc()
		return p.eof
	}
}

func str(p *Parser, s []byte) func() bool {
	return func() bool {
		for _, b := range s {
			if b != p.cc() {
				return true
			}
			p.gc()
		}
		return p.eof
	}
}

/*

	rules for BNF -> genFSM -> BFSM

	generate parser:
	input: BNF data
	BNF -> BFSM -> Rules -> FSM
	output:
	FSM(data)


*/
