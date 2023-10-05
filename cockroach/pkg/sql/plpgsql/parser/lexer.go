// Copyright 2023 The Cockroach Authors.
//
// Use of this software is governed by the Business Source License
// included in the file licenses/BSL.txt.
//
// As of the Change Date specified in that file, in accordance with
// the Business Source License, use of this software will be governed
// by the Apache License, Version 2.0, included in the file
// licenses/APL.txt.

package parser

import (
	"strings"

	"github.com/cockroachdb/cockroach/pkg/sql/parser"
	"github.com/cockroachdb/cockroach/pkg/sql/pgwire/pgcode"
	"github.com/cockroachdb/cockroach/pkg/sql/pgwire/pgerror"
	"github.com/cockroachdb/cockroach/pkg/sql/sem/plpgsqltree"
	"github.com/cockroachdb/cockroach/pkg/sql/sem/tree"
	"github.com/cockroachdb/cockroach/pkg/sql/types"
	unimp "github.com/cockroachdb/cockroach/pkg/util/errorutil/unimplemented"
	"github.com/cockroachdb/errors"
)

type lexer struct {
	in string
	// tokens contains tokens generated by the scanner.
	tokens []plpgsqlSymType

	// The type that should be used when an INT or SERIAL is encountered.
	nakedIntType *types.T

	// lastPos is the position into the tokens slice of the last
	// token returned by Lex().
	lastPos int

	stmt *plpgsqltree.PLpgSQLStmtBlock

	// numPlaceholders is 1 + the highest placeholder index encountered.
	numPlaceholders int
	numAnnotations  tree.AnnotationIdx

	lastError error

	parser plpgsqlParser
}

func (l *lexer) init(sql string, tokens []plpgsqlSymType, nakedIntType *types.T, p plpgsqlParser) {
	l.in = sql
	l.tokens = tokens
	l.lastPos = -1
	l.stmt = nil
	l.numPlaceholders = 0
	l.numAnnotations = 0
	l.lastError = nil
	l.nakedIntType = nakedIntType
	l.parser = p
}

// cleanup is used to avoid holding on to memory unnecessarily (for the cases
// where we reuse a scanner).
func (l *lexer) cleanup() {
	l.tokens = nil
	l.stmt = nil
	l.lastError = nil
}

// Lex lexes a token from input.
// to push the tokens back (move l.pos back).
func (l *lexer) Lex(lval *plpgsqlSymType) int {
	l.lastPos++
	// The core lexing takes place in the scanner. Here we do a small bit of post
	// processing of the lexical tokens so that the grammar only requires
	// one-token lookahead despite SQL requiring multi-token lookahead in some
	// cases. These special cases are handled below and the returned tokens are
	// adjusted to reflect the lookahead (LA) that occurred.
	if l.lastPos >= len(l.tokens) {
		lval.id = 0
		lval.pos = int32(len(l.in))
		lval.str = "EOF"
		return 0
	}
	*lval = l.tokens[l.lastPos]

	switch lval.id {
	case RETURN:
		nextToken := plpgsqlSymType{}
		if l.lastPos+1 < len(l.tokens) {
			nextToken = l.tokens[l.lastPos+1]
		}
		switch nextToken.id {
		case NEXT:
			lval.id = RETURN_NEXT
		case QUERY:
			lval.id = RETURN_QUERY
		}
	case END:
		nextToken := plpgsqlSymType{}
		if l.lastPos+1 < len(l.tokens) {
			nextToken = l.tokens[l.lastPos+1]
		}
		switch nextToken.id {
		case IF:
			lval.id = END_IF
		case CASE:
			lval.id = END_CASE
		}
	case NO:
		nextToken := plpgsqlSymType{}
		if l.lastPos+1 < len(l.tokens) {
			nextToken = l.tokens[l.lastPos+1]
		}
		switch nextToken.id {
		case SCROLL:
			lval.id = NO_SCROLL
		}
	}

	return int(lval.id)
}

// MakeExecSqlStmt makes a PLpgSQLStmtExecSql from current token position.
// TODO(chengxiong): we need to fill in variables as well.
func (l *lexer) MakeExecSqlStmt(startTokenID int) *plpgsqltree.PLpgSQLStmtExecSql {
	sqlToks := make([]string, 0)
	if startTokenID == 0 || startTokenID == ';' {
		l.setErr(errors.AssertionFailedf("plpgsql_execsql: invalid start token"))
	}
	if int(l.lastToken().id) != startTokenID {
		l.setErr(errors.AssertionFailedf("plpgsql_execsql: given start token does not match current pos of lexer"))
	}

	var hasInto bool
	var hasStrict bool
	var preTok plpgsqlSymType
	tok := l.lastToken()
	for {
		if !hasInto {
			sqlToks = append(sqlToks, tok.Str())
		}
		preTok = tok
		l.Lex(&tok)
		if tok.id == ';' {
			break
		}
		if tok.id == 0 {
			l.setErr(errors.AssertionFailedf("unexpected end of function definition"))
		}
		if hasInto && tok.id == STRICT {
			hasStrict = true
			continue
		}
		if tok.id == INTO {
			if preTok.id == INSERT {
				continue
			}
			if preTok.id == MERGE {
				continue
			}
			if startTokenID == IMPORT {
				continue
			}
			if hasInto {
				l.setErr(errors.AssertionFailedf("plpgsql_execsql: INTO specified more than once"))
			}
			hasInto = true
		}
	}
	return &plpgsqltree.PLpgSQLStmtExecSql{
		SqlStmt: strings.Join(sqlToks, " "),
		Into:    hasInto,
		Strict:  hasStrict,
	}
}

func (l *lexer) MakeDynamicExecuteStmt() *plpgsqltree.PLpgSQLStmtDynamicExecute {
	cmdStr, _ := l.ReadSqlConstruct(INTO, USING, ';')
	ret := &plpgsqltree.PLpgSQLStmtDynamicExecute{
		Query: cmdStr,
	}

	var lval plpgsqlSymType
	l.Lex(&lval)
	for {
		if lval.id == INTO {
			if ret.Into {
				l.setErr(errors.AssertionFailedf("seen multiple INTO"))
			}
			ret.Into = true
			nextTok := l.Peek()
			if nextTok.id == int32(STRICT) {
				l.Lex(&lval)
				ret.Strict = true
			}
			// TODO we need to read each "INTO" variable name instead of just a
			// string.
			l.ReadSqlExpressionStr2(USING, ';')
			l.Lex(&lval)
		} else if lval.id == USING {
			if ret.Params != nil {
				l.setErr(errors.AssertionFailedf("seen multiple USINGs"))
			}
			ret.Params = make([]plpgsqltree.PLpgSQLExpr, 0)
			for {
				l.ReadSqlConstruct(',', ';', INTO)
				ret.Params = append(ret.Params, nil)
				l.Lex(&lval)
				if lval.id == ';' {
					break
				}
			}
		} else if lval.id == ';' {
			break
		} else {
			l.setErr(errors.AssertionFailedf("syntax error"))
		}
	}

	return ret
}

func (l *lexer) ProcessForOpenCursor(nullCursorExplicitExpr bool) *plpgsqltree.PLpgSQLStmtOpen {
	openStmt := &plpgsqltree.PLpgSQLStmtOpen{}
	openStmt.CursorOptions = plpgsqltree.PLpgSQLCursorOptFastPlan.Mask()

	if nullCursorExplicitExpr {
		if l.Peek().id == NO {
			l.lastPos++
			if l.Peek().id == SCROLL {
				openStmt.CursorOptions |= plpgsqltree.PLpgSQLCursorOptNoScroll.Mask()
				l.lastPos++
			}
		} else if l.Peek().id == SCROLL {
			openStmt.CursorOptions |= plpgsqltree.PLpgSQLCursorOptScroll.Mask()
			l.lastPos++
		}

		if l.Peek().id != FOR {
			l.setErr(pgerror.New(pgcode.Syntax, "syntax error, expected \"FOR\""))
			return nil
		}

		l.lastPos++
		if l.Peek().id == EXECUTE {
			l.lastPos++
			dynamicQuery, endToken := l.ReadSqlExpressionStr2(USING, ';')
			openStmt.DynamicQuery = dynamicQuery
			l.lastPos++
			if endToken == USING {
				// Continue reading for params for the sql expression till the ending
				// token is not a comma.
				openStmt.Params = make([]string, 0)
				for {
					param, endToken := l.ReadSqlExpressionStr2(',', ';')
					openStmt.Params = append(openStmt.Params, param)
					if endToken != ',' {
						break
					}
					l.lastPos++
				}
			}
		} else {
			openStmt.Query = l.ReadSqlExpressionStr(';')
		}
	} else {
		// read_cursor_args()
		openStmt.ArgQuery = "hello"
	}
	return openStmt
}

// ReadSqlExpressionStr returns the string from the l.lastPos till it sees
// the terminator for the first time. The returned string is made by tokens
// between the starting index (included) to the terminator (not included).
// TODO(plpgsql-team): pass the output to the sql parser
// (i.e. sqlParserImpl.Parse()).
func (l *lexer) ReadSqlExpressionStr(terminator int) (sqlStr string) {
	sqlStr, _ = l.ReadSqlConstruct(terminator, 0, 0)
	return sqlStr
}

func (l *lexer) ReadSqlExpressionStr2(
	terminator1 int, terminator2 int,
) (sqlStr string, terminatorMet int) {
	return l.ReadSqlConstruct(terminator1, terminator2, 0)
}

func (l *lexer) ReadSqlConstruct(
	terminator1 int, terminators ...int,
) (sqlStr string, terminatorMet int) {
	if l.parser.Lookahead() != -1 {
		// Push back the lookahead token so that it can be included in the string.
		l.PushBack(1)
	}
	parenLevel := 0
	startPos := l.lastPos + 1
	for l.lastPos < len(l.tokens) {
		tok := l.Peek()
		if int(tok.id) == terminator1 && parenLevel == 0 {
			terminatorMet = terminator1
			break
		}
		for _, term := range terminators {
			if int(tok.id) == term && parenLevel == 0 {
				terminatorMet = term
			}
		}
		if terminatorMet != 0 {
			break
		}
		if tok.id == '(' || tok.id == '[' {
			parenLevel++
		} else if tok.id == ')' || tok.id == ']' {
			parenLevel--
			if parenLevel < 0 {
				panic(errors.AssertionFailedf("wrongly nested parentheses"))
			}
		}
		l.lastPos++
	}
	if parenLevel != 0 {
		panic(errors.AssertionFailedf("parentheses is badly nested"))
	}
	if startPos > l.lastPos {
		//TODO(jane): show the terminator in the panic message.
		l.setErr(errors.New("missing SQL expression"))
	}
	end := len(l.in)
	if l.lastPos+1 < len(l.tokens) {
		end = int(l.tokens[l.lastPos+1].Pos())
	}
	start := int(l.tokens[startPos].Pos())
	return l.in[start:end], terminatorMet
}

func (l *lexer) ProcessQueryForCursorWithoutExplicitExpr(openStmt *plpgsqltree.PLpgSQLStmtOpen) {
	l.lastPos++
	if int(l.Peek().id) == EXECUTE {
		dynamicQuery, endToken := l.ReadSqlExpressionStr2(USING, ';')
		openStmt.DynamicQuery = dynamicQuery
		if endToken == USING {
			var expr string
			for {
				expr, endToken = l.ReadSqlExpressionStr2(',', ';')
				openStmt.Params = append(openStmt.Params, expr)
				if endToken != ',' {
					break
				}
			}
		}
	} else {
		openStmt.Query = l.ReadSqlExpressionStr(';')
	}
}

// Peek peeks
func (l *lexer) Peek() plpgsqlSymType {
	if l.lastPos+1 < len(l.tokens) {
		return l.tokens[l.lastPos+1]
	}
	return plpgsqlSymType{}
}

// PushBack rewinds the lexer by n tokens.
func (l *lexer) PushBack(n int) {
	if n < 0 {
		panic(errors.AssertionFailedf("negative n provided to PushBack"))
	}
	l.lastPos -= n
	if l.lastPos < -1 {
		// Return to the initialized state.
		l.lastPos = -1
	}
	if n >= 1 {
		// Invalidate the parser lookahead token.
		l.parser.(*plpgsqlParserImpl).char = -1
	}
}

func (l *lexer) lastToken() plpgsqlSymType {
	if l.lastPos < 0 {
		return plpgsqlSymType{}
	}

	if l.lastPos >= len(l.tokens) {
		return plpgsqlSymType{
			id:  0,
			pos: int32(len(l.in)),
			str: "EOF",
		}
	}
	return l.tokens[l.lastPos]
}

// SetStmt is called from the parser when the statement is constructed.
func (l *lexer) SetStmt(stmt plpgsqltree.PLpgSQLStatement) {
	l.stmt = stmt.(*plpgsqltree.PLpgSQLStmtBlock)
}

// setErr is called from parsing action rules to register an error observed
// while running the action. That error becomes the actual "cause" of the
// syntax error.
func (l *lexer) setErr(err error) {
	err = pgerror.WithCandidateCode(err, pgcode.Syntax)
	l.lastError = err
	lastTok := l.lastToken()
	l.lastError = parser.PopulateErrorDetails(lastTok.id, lastTok.str, lastTok.pos, l.lastError, l.in)
}

func (l *lexer) Error(e string) {
	e = strings.TrimPrefix(e, "syntax error: ") // we'll add it again below.
	err := pgerror.WithCandidateCode(errors.Newf("%s", e), pgcode.Syntax)
	lastTok := l.lastToken()
	l.lastError = parser.PopulateErrorDetails(lastTok.id, lastTok.str, lastTok.pos, err, l.in)
}

// Unimplemented wraps Error, setting lastUnimplementedError.
func (l *lexer) Unimplemented(feature string) {
	l.lastError = unimp.New(feature, "this syntax")
	lastTok := l.lastToken()
	l.lastError = parser.PopulateErrorDetails(lastTok.id, lastTok.str, lastTok.pos, l.lastError, l.in)
	l.lastError = &tree.UnsupportedError{
		Err:         l.lastError,
		FeatureName: feature,
	}
}

func (l *lexer) GetTypeFromValidSQLSyntax(sqlStr string) (tree.ResolvableTypeReference, error) {
	return parser.GetTypeFromValidSQLSyntax(sqlStr)
}

func (l *lexer) ParseExpr(sqlStr string) (plpgsqltree.PLpgSQLExpr, error) {
	return parser.ParseExpr(sqlStr)
}