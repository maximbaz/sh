// Copyright (c) 2016, Daniel Martí <mvdan@mvdan.cc>
// See LICENSE for licensing information

package sh

import (
	"bytes"
	"fmt"
	"io"
)

type Prog struct {
	Stmts []Node
}

func (p Prog) String() string {
	return nodeJoin(p.Stmts, "; ")
}

type Node interface {
	fmt.Stringer
}

func nodeJoin(ns []Node, sep string) string {
	var b bytes.Buffer
	for i, n := range ns {
		if i > 0 {
			io.WriteString(&b, sep)
		}
		io.WriteString(&b, n.String())
	}
	return b.String()
}

func stmtJoin(ns []Node) string {
	return nodeJoin(ns, "; ")
}

func wordJoin(ns []Node) string {
	return nodeJoin(ns, " ")
}

type Command struct {
	Args []Node

	Background bool
}

func (c Command) String() string {
	suffix := ""
	if c.Background {
		suffix += " &"
	}
	return wordJoin(c.Args) + suffix
}

type Redirect struct {
	Op  Token
	Obj Node
}

func (r Redirect) String() string {
	return r.Op.String() + r.Obj.String()
}

type Subshell struct {
	Stmts []Node
}

func (s Subshell) String() string {
	return "(" + stmtJoin(s.Stmts) + ")"
}

type Block struct {
	Stmts []Node
}

func (b Block) String() string {
	return "{ " + stmtJoin(b.Stmts) + "; }"
}

type IfStmt struct {
	Cond      Node
	ThenStmts []Node
	Elifs     []Node
	ElseStmts []Node
}

func (s IfStmt) String() string {
	var b bytes.Buffer
	io.WriteString(&b, "if ")
	io.WriteString(&b, s.Cond.String())
	io.WriteString(&b, "; then ")
	io.WriteString(&b, stmtJoin(s.ThenStmts))
	for _, n := range s.Elifs {
		e := n.(Elif)
		io.WriteString(&b, "; ")
		io.WriteString(&b, e.String())
	}
	if len(s.ElseStmts) > 0 {
		io.WriteString(&b, "; else ")
		io.WriteString(&b, stmtJoin(s.ElseStmts))
	}
	io.WriteString(&b, "; fi")
	return b.String()
}

type Elif struct {
	Cond      Node
	ThenStmts []Node
}

func (e Elif) String() string {
	return "elif " + e.Cond.String() + "; then " + stmtJoin(e.ThenStmts)
}

type WhileStmt struct {
	Cond    Node
	DoStmts []Node
}

func (w WhileStmt) String() string {
	return "while " + w.Cond.String() + "; do " + stmtJoin(w.DoStmts) + "; done"
}

type ForStmt struct {
	Name     Node
	WordList []Node
	DoStmts  []Node
}

func (f ForStmt) String() string {
	var b bytes.Buffer
	io.WriteString(&b, "for ")
	io.WriteString(&b, f.Name.String())
	io.WriteString(&b, " in ")
	io.WriteString(&b, wordJoin(f.WordList))
	io.WriteString(&b, "; do ")
	io.WriteString(&b, stmtJoin(f.DoStmts))
	io.WriteString(&b, "; done")
	return b.String()
}

type BinaryExpr struct {
	X, Y Node
	Op   Token
}

func (b BinaryExpr) String() string {
	return fmt.Sprintf("%s %s %s", b.X, b.Op, b.Y)
}

type FuncDecl struct {
	Name Lit
	Body Node
}

func (f FuncDecl) String() string {
	return fmt.Sprintf("%s() %s", f.Name, f.Body)
}

type Word struct {
	Parts []Node
}

func (w Word) String() string {
	return nodeJoin(w.Parts, "")
}

type Lit struct {
	Val string
}

func (l Lit) String() string {
	return l.Val
}

type DblQuoted struct {
	Parts []Node
}

func (q DblQuoted) String() string {
	return `"` + nodeJoin(q.Parts, "") + `"`
}

type CmdSubst struct {
	Stmts []Node
}

func (c CmdSubst) String() string {
	return "$(" + stmtJoin(c.Stmts) + ")"
}

type ParamExp struct {
	Short bool
	Text  string
}

func (p ParamExp) String() string {
	if p.Short {
		return "$" + p.Text
	}
	return "${" + p.Text + "}"
}

type ArithmExp struct {
	Text string
}

func (a ArithmExp) String() string {
	return "$((" + a.Text + "))"
}

type CaseStmt struct {
	Name     Node
	Patterns []Node
}

func (c CaseStmt) String() string {
	return "case " + c.Name.String() + " in " + nodeJoin(c.Patterns, ";; ") + "; esac"
}

type CasePattern struct {
	Parts []Node
	Stmts []Node
}

func (c CasePattern) String() string {
	return nodeJoin(c.Parts, " | ") + ") " + stmtJoin(c.Stmts)
}
