package main

// import (
// 	"fmt"
// 	"os"
// )

// var tab = "    "
// var packageCode = `
// package bnfparser
// `
// var parserCode = `
// type Parser struct {
// 	msg *Message
// 	data []byte
// 	ind int
// 	EOF bool
// }

// func NewParser(data []byte, msg *Message) *Parser {
// 	return &Parser{msg: msg, data: data}
// }

// func (p *Parser) Parse() {
// 	p.Text()
// }

// func (p *Parser) Message() *Message {
// 	return p.msg
// }

// func (p *Parser) cc() byte {
// 	return p.data[p.ind]
// }

// func (p *Parser) gp() int {
// 	return p.ind
// }

// func (p *Parser) gc() {
// 	if p.ind >= len(p.data) {
// 		p.EOF = true
// 		return
// 	}
// 	p.ind++
// }
// func (p *Parser) slice(start, end int) {
// 	return p.data[start:end]
// }

// func isAlpha(b byte) bool {
// 	return (b >= 'a' && b <= 'z') || (b >= 'A' && b <= 'Z')
// }

// func (p *Parser) word() bool {
// 	for !p.EOF {
// 		char := p.cc()
// 		if !isAlpha(char) {
// 			return false
// 		}
// 		gc()
// 	}
// }

// func (p *Parser) checkString(str []byte) bool {
// 	for _, b := range str {
// 		if p.EOF || b != p.cc() {
// 			return false
// 		}
// 		p.gc()
// 	}
// 	return true
// }
// `

// // genFunc generates entity function. Example
// // lvalue := rvalue1 rvalue2 .. rvalueN
// // func (p *Parser) lvalue() {
// //      rvalue1()
// //      rvalue2()
// //          ..
// //      rvalueN()
// //}
// func genFunc(r *Rule) string {
// 	// declare 'start', 'end' indexes of rvalue
// 	body := tab + "var start, end int\n"
// 	body += tab + "start = p.gp()\n"
// 	for _, rvalue := range r.rvalue {
// 		body += tab
// 		switch rvalue.typ {
// 		case String:
// 			// generate checking strings
// 			body += fmt.Sprintf("if !p.checkStr(%s) {\n return \n}\n", rvalue)
// 		case Entity:
// 			// generate entity func as entity()
// 			body += fmt.Sprintf("if !p.%s() {\n return \n}\n", rvalue)
// 		case Word:
// 			// check is entity a word, like word(entity)
// 			body += fmt.Sprintf("if !p.word(%s) {\n return \n}\n", rvalue)
// 		}
// 	}
// 	body += tab + "end = p.gp()\n"

// 	// write to Message field
// 	if r.Marked {
// 		body += tab + fmt.Sprintf("p.msg.%s = p.slice(start, end)\n", r.lvalue.value)
// 	}

// 	return fmt.Sprintf("func (p *Parser) %s(){\n%s}\n", r.lvalue.value, body)
// }

// func appendField(fields *string, lvalue *entity) {
// 	*fields += tab + fmt.Sprintf("%s []byte `json:\"%s\"`\n", lvalue.value, lvalue.value)
// }

// func genCode(rules []*Rule) error {
// 	f, err := os.OpenFile("", os.O_CREATE|os.O_TRUNC|os.O_RDWR, 0666)
// 	if err != nil {
// 		return err
// 	}

// 	fields := ""
// 	funcsCode := make([]string, 0)
// 	for _, rule := range rules {
// 		funcCode := genFunc(rule)
// 		if rule.Marked {
// 			appendField(&fields, &rule.lvalue)
// 		}
// 		funcsCode = append(funcsCode, funcCode)
// 	}

// 	messageCode := fmt.Sprintf("type Message struct{\n%s}\n", fields)

// 	f.Write([]byte(packageCode))
// 	f.Write([]byte(messageCode))
// 	f.Write([]byte(parserCode))
// 	for _, funcCode := range funcsCode {
// 		f.Write([]byte(funcCode))
// 		f.Write([]byte("\n"))
// 	}

// 	return f.Close()
// }

// func genParser(bnf []byte) error {
// 	rules, err := getRules(bnf)
// 	if err != nil {
// 		return err
// 	}

// 	return genCode(rules)
// }
