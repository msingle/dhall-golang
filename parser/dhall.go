
package parser

import (
"bytes"
"errors"
"fmt"
"io"
"io/ioutil"
"math"
"net"
"net/url"
"os"
"path"
"strconv"
"strings"
"unicode"
"unicode/utf8"
)
import . "github.com/philandstuff/dhall-golang/ast"

// Helper function for parsing all the operator parsing blocks
// see OrExpression for an example of how this is used
func ParseOperator(opcode int, first, rest interface{}) Expr {
    out := first.(Expr)
    if rest == nil { return out }
    for _, b := range rest.([]interface{}) {
        nextExpr := b.([]interface{})[3].(Expr)
        out = Operator{OpCode: opcode, L: out, R: nextExpr}
    }
    return out
}


var g = &grammar {
	rules: []*rule{
{
	name: "DhallFile",
	pos: position{line: 36, col: 1, offset: 621},
	expr: &actionExpr{
	pos: position{line: 36, col: 13, offset: 635},
	run: (*parser).callonDhallFile1,
	expr: &seqExpr{
	pos: position{line: 36, col: 13, offset: 635},
	exprs: []interface{}{
&labeledExpr{
	pos: position{line: 36, col: 13, offset: 635},
	label: "e",
	expr: &ruleRefExpr{
	pos: position{line: 36, col: 15, offset: 637},
	name: "CompleteExpression",
},
},
&ruleRefExpr{
	pos: position{line: 36, col: 34, offset: 656},
	name: "EOF",
},
	},
},
},
},
{
	name: "CompleteExpression",
	pos: position{line: 38, col: 1, offset: 679},
	expr: &actionExpr{
	pos: position{line: 38, col: 22, offset: 702},
	run: (*parser).callonCompleteExpression1,
	expr: &seqExpr{
	pos: position{line: 38, col: 22, offset: 702},
	exprs: []interface{}{
&ruleRefExpr{
	pos: position{line: 38, col: 22, offset: 702},
	name: "_",
},
&labeledExpr{
	pos: position{line: 38, col: 24, offset: 704},
	label: "e",
	expr: &ruleRefExpr{
	pos: position{line: 38, col: 26, offset: 706},
	name: "Expression",
},
},
&ruleRefExpr{
	pos: position{line: 38, col: 37, offset: 717},
	name: "_",
},
	},
},
},
},
{
	name: "EOL",
	pos: position{line: 40, col: 1, offset: 738},
	expr: &choiceExpr{
	pos: position{line: 40, col: 7, offset: 746},
	alternatives: []interface{}{
&litMatcher{
	pos: position{line: 40, col: 7, offset: 746},
	val: "\n",
	ignoreCase: false,
},
&actionExpr{
	pos: position{line: 40, col: 14, offset: 753},
	run: (*parser).callonEOL3,
	expr: &litMatcher{
	pos: position{line: 40, col: 14, offset: 753},
	val: "\r\n",
	ignoreCase: false,
},
},
	},
},
},
{
	name: "BlockComment",
	pos: position{line: 42, col: 1, offset: 790},
	expr: &seqExpr{
	pos: position{line: 42, col: 16, offset: 807},
	exprs: []interface{}{
&litMatcher{
	pos: position{line: 42, col: 16, offset: 807},
	val: "{-",
	ignoreCase: false,
},
&ruleRefExpr{
	pos: position{line: 42, col: 21, offset: 812},
	name: "BlockCommentContinue",
},
	},
},
},
{
	name: "BlockCommentChunk",
	pos: position{line: 44, col: 1, offset: 834},
	expr: &choiceExpr{
	pos: position{line: 45, col: 5, offset: 860},
	alternatives: []interface{}{
&ruleRefExpr{
	pos: position{line: 45, col: 5, offset: 860},
	name: "BlockComment",
},
&charClassMatcher{
	pos: position{line: 46, col: 5, offset: 877},
	val: "[\\t\\u0020-\\U0010ffff]",
	chars: []rune{'\t',},
	ranges: []rune{' ','\U0010ffff',},
	ignoreCase: false,
	inverted: false,
},
&ruleRefExpr{
	pos: position{line: 47, col: 5, offset: 903},
	name: "EOL",
},
	},
},
},
{
	name: "BlockCommentContinue",
	pos: position{line: 49, col: 1, offset: 908},
	expr: &choiceExpr{
	pos: position{line: 49, col: 24, offset: 933},
	alternatives: []interface{}{
&litMatcher{
	pos: position{line: 49, col: 24, offset: 933},
	val: "-}",
	ignoreCase: false,
},
&seqExpr{
	pos: position{line: 49, col: 31, offset: 940},
	exprs: []interface{}{
&ruleRefExpr{
	pos: position{line: 49, col: 31, offset: 940},
	name: "BlockCommentChunk",
},
&ruleRefExpr{
	pos: position{line: 49, col: 49, offset: 958},
	name: "BlockCommentContinue",
},
	},
},
	},
},
},
{
	name: "NotEOL",
	pos: position{line: 51, col: 1, offset: 980},
	expr: &charClassMatcher{
	pos: position{line: 51, col: 10, offset: 991},
	val: "[\\t\\u0020-\\U0010ffff]",
	chars: []rune{'\t',},
	ranges: []rune{' ','\U0010ffff',},
	ignoreCase: false,
	inverted: false,
},
},
{
	name: "LineComment",
	pos: position{line: 53, col: 1, offset: 1014},
	expr: &actionExpr{
	pos: position{line: 53, col: 15, offset: 1030},
	run: (*parser).callonLineComment1,
	expr: &seqExpr{
	pos: position{line: 53, col: 15, offset: 1030},
	exprs: []interface{}{
&litMatcher{
	pos: position{line: 53, col: 15, offset: 1030},
	val: "--",
	ignoreCase: false,
},
&labeledExpr{
	pos: position{line: 53, col: 20, offset: 1035},
	label: "content",
	expr: &actionExpr{
	pos: position{line: 53, col: 29, offset: 1044},
	run: (*parser).callonLineComment5,
	expr: &zeroOrMoreExpr{
	pos: position{line: 53, col: 29, offset: 1044},
	expr: &ruleRefExpr{
	pos: position{line: 53, col: 29, offset: 1044},
	name: "NotEOL",
},
},
},
},
&ruleRefExpr{
	pos: position{line: 53, col: 68, offset: 1083},
	name: "EOL",
},
	},
},
},
},
{
	name: "WhitespaceChunk",
	pos: position{line: 55, col: 1, offset: 1112},
	expr: &choiceExpr{
	pos: position{line: 55, col: 19, offset: 1132},
	alternatives: []interface{}{
&litMatcher{
	pos: position{line: 55, col: 19, offset: 1132},
	val: " ",
	ignoreCase: false,
},
&litMatcher{
	pos: position{line: 55, col: 25, offset: 1138},
	val: "\t",
	ignoreCase: false,
},
&ruleRefExpr{
	pos: position{line: 55, col: 32, offset: 1145},
	name: "EOL",
},
&ruleRefExpr{
	pos: position{line: 55, col: 38, offset: 1151},
	name: "LineComment",
},
&ruleRefExpr{
	pos: position{line: 55, col: 52, offset: 1165},
	name: "BlockComment",
},
	},
},
},
{
	name: "_",
	pos: position{line: 57, col: 1, offset: 1179},
	expr: &zeroOrMoreExpr{
	pos: position{line: 57, col: 5, offset: 1185},
	expr: &ruleRefExpr{
	pos: position{line: 57, col: 5, offset: 1185},
	name: "WhitespaceChunk",
},
},
},
{
	name: "_1",
	pos: position{line: 59, col: 1, offset: 1203},
	expr: &oneOrMoreExpr{
	pos: position{line: 59, col: 6, offset: 1210},
	expr: &ruleRefExpr{
	pos: position{line: 59, col: 6, offset: 1210},
	name: "WhitespaceChunk",
},
},
},
{
	name: "Digit",
	pos: position{line: 61, col: 1, offset: 1228},
	expr: &charClassMatcher{
	pos: position{line: 61, col: 9, offset: 1238},
	val: "[0-9]",
	ranges: []rune{'0','9',},
	ignoreCase: false,
	inverted: false,
},
},
{
	name: "HexDig",
	pos: position{line: 63, col: 1, offset: 1245},
	expr: &choiceExpr{
	pos: position{line: 63, col: 10, offset: 1256},
	alternatives: []interface{}{
&ruleRefExpr{
	pos: position{line: 63, col: 10, offset: 1256},
	name: "Digit",
},
&charClassMatcher{
	pos: position{line: 63, col: 18, offset: 1264},
	val: "[a-f]i",
	ranges: []rune{'a','f',},
	ignoreCase: true,
	inverted: false,
},
	},
},
},
{
	name: "SimpleLabelFirstChar",
	pos: position{line: 65, col: 1, offset: 1272},
	expr: &charClassMatcher{
	pos: position{line: 65, col: 24, offset: 1297},
	val: "[A-Za-z_]",
	chars: []rune{'_',},
	ranges: []rune{'A','Z','a','z',},
	ignoreCase: false,
	inverted: false,
},
},
{
	name: "SimpleLabelNextChar",
	pos: position{line: 66, col: 1, offset: 1307},
	expr: &charClassMatcher{
	pos: position{line: 66, col: 23, offset: 1331},
	val: "[A-Za-z0-9_/-]",
	chars: []rune{'_','/','-',},
	ranges: []rune{'A','Z','a','z','0','9',},
	ignoreCase: false,
	inverted: false,
},
},
{
	name: "SimpleLabel",
	pos: position{line: 67, col: 1, offset: 1346},
	expr: &choiceExpr{
	pos: position{line: 67, col: 15, offset: 1362},
	alternatives: []interface{}{
&actionExpr{
	pos: position{line: 67, col: 15, offset: 1362},
	run: (*parser).callonSimpleLabel2,
	expr: &seqExpr{
	pos: position{line: 67, col: 15, offset: 1362},
	exprs: []interface{}{
&ruleRefExpr{
	pos: position{line: 67, col: 15, offset: 1362},
	name: "Keyword",
},
&oneOrMoreExpr{
	pos: position{line: 67, col: 23, offset: 1370},
	expr: &ruleRefExpr{
	pos: position{line: 67, col: 23, offset: 1370},
	name: "SimpleLabelNextChar",
},
},
	},
},
},
&actionExpr{
	pos: position{line: 68, col: 13, offset: 1434},
	run: (*parser).callonSimpleLabel7,
	expr: &seqExpr{
	pos: position{line: 68, col: 13, offset: 1434},
	exprs: []interface{}{
&notExpr{
	pos: position{line: 68, col: 13, offset: 1434},
	expr: &ruleRefExpr{
	pos: position{line: 68, col: 14, offset: 1435},
	name: "Keyword",
},
},
&ruleRefExpr{
	pos: position{line: 68, col: 22, offset: 1443},
	name: "SimpleLabelFirstChar",
},
&zeroOrMoreExpr{
	pos: position{line: 68, col: 43, offset: 1464},
	expr: &ruleRefExpr{
	pos: position{line: 68, col: 43, offset: 1464},
	name: "SimpleLabelNextChar",
},
},
	},
},
},
	},
},
},
{
	name: "Label",
	pos: position{line: 75, col: 1, offset: 1565},
	expr: &actionExpr{
	pos: position{line: 75, col: 9, offset: 1575},
	run: (*parser).callonLabel1,
	expr: &labeledExpr{
	pos: position{line: 75, col: 9, offset: 1575},
	label: "label",
	expr: &ruleRefExpr{
	pos: position{line: 75, col: 15, offset: 1581},
	name: "SimpleLabel",
},
},
},
},
{
	name: "NonreservedLabel",
	pos: position{line: 77, col: 1, offset: 1616},
	expr: &choiceExpr{
	pos: position{line: 77, col: 20, offset: 1637},
	alternatives: []interface{}{
&actionExpr{
	pos: position{line: 77, col: 20, offset: 1637},
	run: (*parser).callonNonreservedLabel2,
	expr: &seqExpr{
	pos: position{line: 77, col: 20, offset: 1637},
	exprs: []interface{}{
&andExpr{
	pos: position{line: 77, col: 20, offset: 1637},
	expr: &seqExpr{
	pos: position{line: 77, col: 22, offset: 1639},
	exprs: []interface{}{
&ruleRefExpr{
	pos: position{line: 77, col: 22, offset: 1639},
	name: "Reserved",
},
&ruleRefExpr{
	pos: position{line: 77, col: 31, offset: 1648},
	name: "SimpleLabelNextChar",
},
	},
},
},
&labeledExpr{
	pos: position{line: 77, col: 52, offset: 1669},
	label: "label",
	expr: &ruleRefExpr{
	pos: position{line: 77, col: 58, offset: 1675},
	name: "Label",
},
},
	},
},
},
&actionExpr{
	pos: position{line: 78, col: 19, offset: 1721},
	run: (*parser).callonNonreservedLabel10,
	expr: &seqExpr{
	pos: position{line: 78, col: 19, offset: 1721},
	exprs: []interface{}{
&notExpr{
	pos: position{line: 78, col: 19, offset: 1721},
	expr: &ruleRefExpr{
	pos: position{line: 78, col: 20, offset: 1722},
	name: "Reserved",
},
},
&labeledExpr{
	pos: position{line: 78, col: 29, offset: 1731},
	label: "label",
	expr: &ruleRefExpr{
	pos: position{line: 78, col: 35, offset: 1737},
	name: "Label",
},
},
	},
},
},
	},
},
},
{
	name: "AnyLabel",
	pos: position{line: 80, col: 1, offset: 1766},
	expr: &ruleRefExpr{
	pos: position{line: 80, col: 12, offset: 1779},
	name: "Label",
},
},
{
	name: "DoubleQuoteChunk",
	pos: position{line: 83, col: 1, offset: 1787},
	expr: &choiceExpr{
	pos: position{line: 84, col: 6, offset: 1813},
	alternatives: []interface{}{
&ruleRefExpr{
	pos: position{line: 84, col: 6, offset: 1813},
	name: "Interpolation",
},
&actionExpr{
	pos: position{line: 85, col: 6, offset: 1832},
	run: (*parser).callonDoubleQuoteChunk3,
	expr: &seqExpr{
	pos: position{line: 85, col: 6, offset: 1832},
	exprs: []interface{}{
&litMatcher{
	pos: position{line: 85, col: 6, offset: 1832},
	val: "\\",
	ignoreCase: false,
},
&labeledExpr{
	pos: position{line: 85, col: 11, offset: 1837},
	label: "e",
	expr: &ruleRefExpr{
	pos: position{line: 85, col: 13, offset: 1839},
	name: "DoubleQuoteEscaped",
},
},
	},
},
},
&ruleRefExpr{
	pos: position{line: 86, col: 6, offset: 1881},
	name: "DoubleQuoteChar",
},
	},
},
},
{
	name: "DoubleQuoteEscaped",
	pos: position{line: 88, col: 1, offset: 1898},
	expr: &choiceExpr{
	pos: position{line: 89, col: 8, offset: 1928},
	alternatives: []interface{}{
&litMatcher{
	pos: position{line: 89, col: 8, offset: 1928},
	val: "\"",
	ignoreCase: false,
},
&litMatcher{
	pos: position{line: 90, col: 8, offset: 1939},
	val: "$",
	ignoreCase: false,
},
&litMatcher{
	pos: position{line: 91, col: 8, offset: 1950},
	val: "\\",
	ignoreCase: false,
},
&litMatcher{
	pos: position{line: 92, col: 8, offset: 1962},
	val: "/",
	ignoreCase: false,
},
&actionExpr{
	pos: position{line: 93, col: 8, offset: 1973},
	run: (*parser).callonDoubleQuoteEscaped6,
	expr: &litMatcher{
	pos: position{line: 93, col: 8, offset: 1973},
	val: "b",
	ignoreCase: false,
},
},
&actionExpr{
	pos: position{line: 94, col: 8, offset: 2013},
	run: (*parser).callonDoubleQuoteEscaped8,
	expr: &litMatcher{
	pos: position{line: 94, col: 8, offset: 2013},
	val: "f",
	ignoreCase: false,
},
},
&actionExpr{
	pos: position{line: 95, col: 8, offset: 2053},
	run: (*parser).callonDoubleQuoteEscaped10,
	expr: &litMatcher{
	pos: position{line: 95, col: 8, offset: 2053},
	val: "n",
	ignoreCase: false,
},
},
&actionExpr{
	pos: position{line: 96, col: 8, offset: 2093},
	run: (*parser).callonDoubleQuoteEscaped12,
	expr: &litMatcher{
	pos: position{line: 96, col: 8, offset: 2093},
	val: "r",
	ignoreCase: false,
},
},
&actionExpr{
	pos: position{line: 97, col: 8, offset: 2133},
	run: (*parser).callonDoubleQuoteEscaped14,
	expr: &litMatcher{
	pos: position{line: 97, col: 8, offset: 2133},
	val: "t",
	ignoreCase: false,
},
},
&actionExpr{
	pos: position{line: 98, col: 8, offset: 2173},
	run: (*parser).callonDoubleQuoteEscaped16,
	expr: &seqExpr{
	pos: position{line: 98, col: 8, offset: 2173},
	exprs: []interface{}{
&litMatcher{
	pos: position{line: 98, col: 8, offset: 2173},
	val: "u",
	ignoreCase: false,
},
&ruleRefExpr{
	pos: position{line: 98, col: 12, offset: 2177},
	name: "HexDig",
},
&ruleRefExpr{
	pos: position{line: 98, col: 19, offset: 2184},
	name: "HexDig",
},
&ruleRefExpr{
	pos: position{line: 98, col: 26, offset: 2191},
	name: "HexDig",
},
&ruleRefExpr{
	pos: position{line: 98, col: 33, offset: 2198},
	name: "HexDig",
},
	},
},
},
	},
},
},
{
	name: "DoubleQuoteChar",
	pos: position{line: 103, col: 1, offset: 2330},
	expr: &choiceExpr{
	pos: position{line: 104, col: 6, offset: 2355},
	alternatives: []interface{}{
&charClassMatcher{
	pos: position{line: 104, col: 6, offset: 2355},
	val: "[\\x20-\\x21]",
	ranges: []rune{' ','!',},
	ignoreCase: false,
	inverted: false,
},
&charClassMatcher{
	pos: position{line: 105, col: 6, offset: 2372},
	val: "[\\x23-\\x5b]",
	ranges: []rune{'#','[',},
	ignoreCase: false,
	inverted: false,
},
&charClassMatcher{
	pos: position{line: 106, col: 6, offset: 2389},
	val: "[\\x5d-\\U0010ffff]",
	ranges: []rune{']','\U0010ffff',},
	ignoreCase: false,
	inverted: false,
},
	},
},
},
{
	name: "DoubleQuoteLiteral",
	pos: position{line: 108, col: 1, offset: 2408},
	expr: &actionExpr{
	pos: position{line: 108, col: 22, offset: 2431},
	run: (*parser).callonDoubleQuoteLiteral1,
	expr: &seqExpr{
	pos: position{line: 108, col: 22, offset: 2431},
	exprs: []interface{}{
&litMatcher{
	pos: position{line: 108, col: 22, offset: 2431},
	val: "\"",
	ignoreCase: false,
},
&labeledExpr{
	pos: position{line: 108, col: 26, offset: 2435},
	label: "chunks",
	expr: &zeroOrMoreExpr{
	pos: position{line: 108, col: 33, offset: 2442},
	expr: &ruleRefExpr{
	pos: position{line: 108, col: 33, offset: 2442},
	name: "DoubleQuoteChunk",
},
},
},
&litMatcher{
	pos: position{line: 108, col: 51, offset: 2460},
	val: "\"",
	ignoreCase: false,
},
	},
},
},
},
{
	name: "SingleQuoteContinue",
	pos: position{line: 125, col: 1, offset: 2928},
	expr: &choiceExpr{
	pos: position{line: 126, col: 7, offset: 2958},
	alternatives: []interface{}{
&seqExpr{
	pos: position{line: 126, col: 7, offset: 2958},
	exprs: []interface{}{
&ruleRefExpr{
	pos: position{line: 126, col: 7, offset: 2958},
	name: "Interpolation",
},
&ruleRefExpr{
	pos: position{line: 126, col: 21, offset: 2972},
	name: "SingleQuoteContinue",
},
	},
},
&seqExpr{
	pos: position{line: 127, col: 7, offset: 2998},
	exprs: []interface{}{
&ruleRefExpr{
	pos: position{line: 127, col: 7, offset: 2998},
	name: "EscapedQuotePair",
},
&ruleRefExpr{
	pos: position{line: 127, col: 24, offset: 3015},
	name: "SingleQuoteContinue",
},
	},
},
&seqExpr{
	pos: position{line: 128, col: 7, offset: 3041},
	exprs: []interface{}{
&ruleRefExpr{
	pos: position{line: 128, col: 7, offset: 3041},
	name: "EscapedInterpolation",
},
&ruleRefExpr{
	pos: position{line: 128, col: 28, offset: 3062},
	name: "SingleQuoteContinue",
},
	},
},
&seqExpr{
	pos: position{line: 129, col: 7, offset: 3088},
	exprs: []interface{}{
&ruleRefExpr{
	pos: position{line: 129, col: 7, offset: 3088},
	name: "SingleQuoteChar",
},
&ruleRefExpr{
	pos: position{line: 129, col: 23, offset: 3104},
	name: "SingleQuoteContinue",
},
	},
},
&litMatcher{
	pos: position{line: 130, col: 7, offset: 3130},
	val: "''",
	ignoreCase: false,
},
	},
},
},
{
	name: "EscapedQuotePair",
	pos: position{line: 132, col: 1, offset: 3136},
	expr: &actionExpr{
	pos: position{line: 132, col: 20, offset: 3157},
	run: (*parser).callonEscapedQuotePair1,
	expr: &litMatcher{
	pos: position{line: 132, col: 20, offset: 3157},
	val: "'''",
	ignoreCase: false,
},
},
},
{
	name: "EscapedInterpolation",
	pos: position{line: 136, col: 1, offset: 3292},
	expr: &actionExpr{
	pos: position{line: 136, col: 24, offset: 3317},
	run: (*parser).callonEscapedInterpolation1,
	expr: &litMatcher{
	pos: position{line: 136, col: 24, offset: 3317},
	val: "''${",
	ignoreCase: false,
},
},
},
{
	name: "SingleQuoteChar",
	pos: position{line: 138, col: 1, offset: 3359},
	expr: &choiceExpr{
	pos: position{line: 139, col: 6, offset: 3384},
	alternatives: []interface{}{
&charClassMatcher{
	pos: position{line: 139, col: 6, offset: 3384},
	val: "[\\x20-\\U0010ffff]",
	ranges: []rune{' ','\U0010ffff',},
	ignoreCase: false,
	inverted: false,
},
&litMatcher{
	pos: position{line: 140, col: 6, offset: 3407},
	val: "\t",
	ignoreCase: false,
},
&ruleRefExpr{
	pos: position{line: 141, col: 6, offset: 3417},
	name: "EOL",
},
	},
},
},
{
	name: "SingleQuoteLiteral",
	pos: position{line: 143, col: 1, offset: 3422},
	expr: &actionExpr{
	pos: position{line: 143, col: 22, offset: 3445},
	run: (*parser).callonSingleQuoteLiteral1,
	expr: &seqExpr{
	pos: position{line: 143, col: 22, offset: 3445},
	exprs: []interface{}{
&litMatcher{
	pos: position{line: 143, col: 22, offset: 3445},
	val: "''",
	ignoreCase: false,
},
&ruleRefExpr{
	pos: position{line: 143, col: 27, offset: 3450},
	name: "EOL",
},
&labeledExpr{
	pos: position{line: 143, col: 31, offset: 3454},
	label: "content",
	expr: &ruleRefExpr{
	pos: position{line: 143, col: 39, offset: 3462},
	name: "SingleQuoteContinue",
},
},
	},
},
},
},
{
	name: "Interpolation",
	pos: position{line: 161, col: 1, offset: 4012},
	expr: &actionExpr{
	pos: position{line: 161, col: 17, offset: 4030},
	run: (*parser).callonInterpolation1,
	expr: &seqExpr{
	pos: position{line: 161, col: 17, offset: 4030},
	exprs: []interface{}{
&litMatcher{
	pos: position{line: 161, col: 17, offset: 4030},
	val: "${",
	ignoreCase: false,
},
&labeledExpr{
	pos: position{line: 161, col: 22, offset: 4035},
	label: "e",
	expr: &ruleRefExpr{
	pos: position{line: 161, col: 24, offset: 4037},
	name: "CompleteExpression",
},
},
&litMatcher{
	pos: position{line: 161, col: 43, offset: 4056},
	val: "}",
	ignoreCase: false,
},
	},
},
},
},
{
	name: "TextLiteral",
	pos: position{line: 163, col: 1, offset: 4079},
	expr: &choiceExpr{
	pos: position{line: 163, col: 15, offset: 4095},
	alternatives: []interface{}{
&ruleRefExpr{
	pos: position{line: 163, col: 15, offset: 4095},
	name: "DoubleQuoteLiteral",
},
&ruleRefExpr{
	pos: position{line: 163, col: 36, offset: 4116},
	name: "SingleQuoteLiteral",
},
	},
},
},
{
	name: "Reserved",
	pos: position{line: 166, col: 1, offset: 4221},
	expr: &choiceExpr{
	pos: position{line: 167, col: 5, offset: 4238},
	alternatives: []interface{}{
&actionExpr{
	pos: position{line: 167, col: 5, offset: 4238},
	run: (*parser).callonReserved2,
	expr: &litMatcher{
	pos: position{line: 167, col: 5, offset: 4238},
	val: "Natural/build",
	ignoreCase: false,
},
},
&actionExpr{
	pos: position{line: 168, col: 5, offset: 4287},
	run: (*parser).callonReserved4,
	expr: &litMatcher{
	pos: position{line: 168, col: 5, offset: 4287},
	val: "Natural/fold",
	ignoreCase: false,
},
},
&actionExpr{
	pos: position{line: 169, col: 5, offset: 4334},
	run: (*parser).callonReserved6,
	expr: &litMatcher{
	pos: position{line: 169, col: 5, offset: 4334},
	val: "Natural/isZero",
	ignoreCase: false,
},
},
&actionExpr{
	pos: position{line: 170, col: 5, offset: 4385},
	run: (*parser).callonReserved8,
	expr: &litMatcher{
	pos: position{line: 170, col: 5, offset: 4385},
	val: "Natural/even",
	ignoreCase: false,
},
},
&actionExpr{
	pos: position{line: 171, col: 5, offset: 4432},
	run: (*parser).callonReserved10,
	expr: &litMatcher{
	pos: position{line: 171, col: 5, offset: 4432},
	val: "Natural/odd",
	ignoreCase: false,
},
},
&actionExpr{
	pos: position{line: 172, col: 5, offset: 4477},
	run: (*parser).callonReserved12,
	expr: &litMatcher{
	pos: position{line: 172, col: 5, offset: 4477},
	val: "Natural/toInteger",
	ignoreCase: false,
},
},
&actionExpr{
	pos: position{line: 173, col: 5, offset: 4534},
	run: (*parser).callonReserved14,
	expr: &litMatcher{
	pos: position{line: 173, col: 5, offset: 4534},
	val: "Natural/show",
	ignoreCase: false,
},
},
&actionExpr{
	pos: position{line: 174, col: 5, offset: 4581},
	run: (*parser).callonReserved16,
	expr: &litMatcher{
	pos: position{line: 174, col: 5, offset: 4581},
	val: "Integer/toDouble",
	ignoreCase: false,
},
},
&actionExpr{
	pos: position{line: 175, col: 5, offset: 4636},
	run: (*parser).callonReserved18,
	expr: &litMatcher{
	pos: position{line: 175, col: 5, offset: 4636},
	val: "Integer/show",
	ignoreCase: false,
},
},
&actionExpr{
	pos: position{line: 176, col: 5, offset: 4683},
	run: (*parser).callonReserved20,
	expr: &litMatcher{
	pos: position{line: 176, col: 5, offset: 4683},
	val: "Double/show",
	ignoreCase: false,
},
},
&actionExpr{
	pos: position{line: 177, col: 5, offset: 4728},
	run: (*parser).callonReserved22,
	expr: &litMatcher{
	pos: position{line: 177, col: 5, offset: 4728},
	val: "List/build",
	ignoreCase: false,
},
},
&actionExpr{
	pos: position{line: 178, col: 5, offset: 4771},
	run: (*parser).callonReserved24,
	expr: &litMatcher{
	pos: position{line: 178, col: 5, offset: 4771},
	val: "List/fold",
	ignoreCase: false,
},
},
&actionExpr{
	pos: position{line: 179, col: 5, offset: 4812},
	run: (*parser).callonReserved26,
	expr: &litMatcher{
	pos: position{line: 179, col: 5, offset: 4812},
	val: "List/length",
	ignoreCase: false,
},
},
&actionExpr{
	pos: position{line: 180, col: 5, offset: 4857},
	run: (*parser).callonReserved28,
	expr: &litMatcher{
	pos: position{line: 180, col: 5, offset: 4857},
	val: "List/head",
	ignoreCase: false,
},
},
&actionExpr{
	pos: position{line: 181, col: 5, offset: 4898},
	run: (*parser).callonReserved30,
	expr: &litMatcher{
	pos: position{line: 181, col: 5, offset: 4898},
	val: "List/last",
	ignoreCase: false,
},
},
&actionExpr{
	pos: position{line: 182, col: 5, offset: 4939},
	run: (*parser).callonReserved32,
	expr: &litMatcher{
	pos: position{line: 182, col: 5, offset: 4939},
	val: "List/indexed",
	ignoreCase: false,
},
},
&actionExpr{
	pos: position{line: 183, col: 5, offset: 4986},
	run: (*parser).callonReserved34,
	expr: &litMatcher{
	pos: position{line: 183, col: 5, offset: 4986},
	val: "List/reverse",
	ignoreCase: false,
},
},
&actionExpr{
	pos: position{line: 184, col: 5, offset: 5033},
	run: (*parser).callonReserved36,
	expr: &litMatcher{
	pos: position{line: 184, col: 5, offset: 5033},
	val: "Optional/build",
	ignoreCase: false,
},
},
&actionExpr{
	pos: position{line: 185, col: 5, offset: 5084},
	run: (*parser).callonReserved38,
	expr: &litMatcher{
	pos: position{line: 185, col: 5, offset: 5084},
	val: "Optional/fold",
	ignoreCase: false,
},
},
&actionExpr{
	pos: position{line: 186, col: 5, offset: 5133},
	run: (*parser).callonReserved40,
	expr: &litMatcher{
	pos: position{line: 186, col: 5, offset: 5133},
	val: "Text/show",
	ignoreCase: false,
},
},
&actionExpr{
	pos: position{line: 187, col: 5, offset: 5174},
	run: (*parser).callonReserved42,
	expr: &litMatcher{
	pos: position{line: 187, col: 5, offset: 5174},
	val: "Bool",
	ignoreCase: false,
},
},
&actionExpr{
	pos: position{line: 188, col: 5, offset: 5206},
	run: (*parser).callonReserved44,
	expr: &litMatcher{
	pos: position{line: 188, col: 5, offset: 5206},
	val: "True",
	ignoreCase: false,
},
},
&actionExpr{
	pos: position{line: 189, col: 5, offset: 5238},
	run: (*parser).callonReserved46,
	expr: &litMatcher{
	pos: position{line: 189, col: 5, offset: 5238},
	val: "False",
	ignoreCase: false,
},
},
&actionExpr{
	pos: position{line: 190, col: 5, offset: 5272},
	run: (*parser).callonReserved48,
	expr: &litMatcher{
	pos: position{line: 190, col: 5, offset: 5272},
	val: "Optional",
	ignoreCase: false,
},
},
&actionExpr{
	pos: position{line: 191, col: 5, offset: 5312},
	run: (*parser).callonReserved50,
	expr: &litMatcher{
	pos: position{line: 191, col: 5, offset: 5312},
	val: "Natural",
	ignoreCase: false,
},
},
&actionExpr{
	pos: position{line: 192, col: 5, offset: 5350},
	run: (*parser).callonReserved52,
	expr: &litMatcher{
	pos: position{line: 192, col: 5, offset: 5350},
	val: "Integer",
	ignoreCase: false,
},
},
&actionExpr{
	pos: position{line: 193, col: 5, offset: 5388},
	run: (*parser).callonReserved54,
	expr: &litMatcher{
	pos: position{line: 193, col: 5, offset: 5388},
	val: "Double",
	ignoreCase: false,
},
},
&actionExpr{
	pos: position{line: 194, col: 5, offset: 5424},
	run: (*parser).callonReserved56,
	expr: &litMatcher{
	pos: position{line: 194, col: 5, offset: 5424},
	val: "Text",
	ignoreCase: false,
},
},
&actionExpr{
	pos: position{line: 195, col: 5, offset: 5456},
	run: (*parser).callonReserved58,
	expr: &litMatcher{
	pos: position{line: 195, col: 5, offset: 5456},
	val: "List",
	ignoreCase: false,
},
},
&actionExpr{
	pos: position{line: 196, col: 5, offset: 5488},
	run: (*parser).callonReserved60,
	expr: &litMatcher{
	pos: position{line: 196, col: 5, offset: 5488},
	val: "None",
	ignoreCase: false,
},
},
&actionExpr{
	pos: position{line: 197, col: 5, offset: 5520},
	run: (*parser).callonReserved62,
	expr: &litMatcher{
	pos: position{line: 197, col: 5, offset: 5520},
	val: "Type",
	ignoreCase: false,
},
},
&actionExpr{
	pos: position{line: 198, col: 5, offset: 5552},
	run: (*parser).callonReserved64,
	expr: &litMatcher{
	pos: position{line: 198, col: 5, offset: 5552},
	val: "Kind",
	ignoreCase: false,
},
},
&actionExpr{
	pos: position{line: 199, col: 5, offset: 5584},
	run: (*parser).callonReserved66,
	expr: &litMatcher{
	pos: position{line: 199, col: 5, offset: 5584},
	val: "Sort",
	ignoreCase: false,
},
},
	},
},
},
{
	name: "If",
	pos: position{line: 201, col: 1, offset: 5613},
	expr: &litMatcher{
	pos: position{line: 201, col: 6, offset: 5620},
	val: "if",
	ignoreCase: false,
},
},
{
	name: "Then",
	pos: position{line: 202, col: 1, offset: 5625},
	expr: &litMatcher{
	pos: position{line: 202, col: 8, offset: 5634},
	val: "then",
	ignoreCase: false,
},
},
{
	name: "Else",
	pos: position{line: 203, col: 1, offset: 5641},
	expr: &litMatcher{
	pos: position{line: 203, col: 8, offset: 5650},
	val: "else",
	ignoreCase: false,
},
},
{
	name: "Let",
	pos: position{line: 204, col: 1, offset: 5657},
	expr: &litMatcher{
	pos: position{line: 204, col: 7, offset: 5665},
	val: "let",
	ignoreCase: false,
},
},
{
	name: "In",
	pos: position{line: 205, col: 1, offset: 5671},
	expr: &litMatcher{
	pos: position{line: 205, col: 6, offset: 5678},
	val: "in",
	ignoreCase: false,
},
},
{
	name: "As",
	pos: position{line: 206, col: 1, offset: 5683},
	expr: &litMatcher{
	pos: position{line: 206, col: 6, offset: 5690},
	val: "as",
	ignoreCase: false,
},
},
{
	name: "Using",
	pos: position{line: 207, col: 1, offset: 5695},
	expr: &litMatcher{
	pos: position{line: 207, col: 9, offset: 5705},
	val: "using",
	ignoreCase: false,
},
},
{
	name: "Merge",
	pos: position{line: 208, col: 1, offset: 5713},
	expr: &litMatcher{
	pos: position{line: 208, col: 9, offset: 5723},
	val: "merge",
	ignoreCase: false,
},
},
{
	name: "Missing",
	pos: position{line: 209, col: 1, offset: 5731},
	expr: &actionExpr{
	pos: position{line: 209, col: 11, offset: 5743},
	run: (*parser).callonMissing1,
	expr: &litMatcher{
	pos: position{line: 209, col: 11, offset: 5743},
	val: "missing",
	ignoreCase: false,
},
},
},
{
	name: "True",
	pos: position{line: 210, col: 1, offset: 5779},
	expr: &litMatcher{
	pos: position{line: 210, col: 8, offset: 5788},
	val: "True",
	ignoreCase: false,
},
},
{
	name: "False",
	pos: position{line: 211, col: 1, offset: 5795},
	expr: &litMatcher{
	pos: position{line: 211, col: 9, offset: 5805},
	val: "False",
	ignoreCase: false,
},
},
{
	name: "Infinity",
	pos: position{line: 212, col: 1, offset: 5813},
	expr: &litMatcher{
	pos: position{line: 212, col: 12, offset: 5826},
	val: "Infinity",
	ignoreCase: false,
},
},
{
	name: "NaN",
	pos: position{line: 213, col: 1, offset: 5837},
	expr: &litMatcher{
	pos: position{line: 213, col: 7, offset: 5845},
	val: "NaN",
	ignoreCase: false,
},
},
{
	name: "Some",
	pos: position{line: 214, col: 1, offset: 5851},
	expr: &litMatcher{
	pos: position{line: 214, col: 8, offset: 5860},
	val: "Some",
	ignoreCase: false,
},
},
{
	name: "Keyword",
	pos: position{line: 216, col: 1, offset: 5868},
	expr: &choiceExpr{
	pos: position{line: 217, col: 5, offset: 5884},
	alternatives: []interface{}{
&ruleRefExpr{
	pos: position{line: 217, col: 5, offset: 5884},
	name: "If",
},
&ruleRefExpr{
	pos: position{line: 217, col: 10, offset: 5889},
	name: "Then",
},
&ruleRefExpr{
	pos: position{line: 217, col: 17, offset: 5896},
	name: "Else",
},
&ruleRefExpr{
	pos: position{line: 218, col: 5, offset: 5905},
	name: "Let",
},
&ruleRefExpr{
	pos: position{line: 218, col: 11, offset: 5911},
	name: "In",
},
&ruleRefExpr{
	pos: position{line: 219, col: 5, offset: 5918},
	name: "Using",
},
&ruleRefExpr{
	pos: position{line: 219, col: 13, offset: 5926},
	name: "Missing",
},
&ruleRefExpr{
	pos: position{line: 219, col: 23, offset: 5936},
	name: "As",
},
&ruleRefExpr{
	pos: position{line: 220, col: 5, offset: 5943},
	name: "True",
},
&ruleRefExpr{
	pos: position{line: 220, col: 12, offset: 5950},
	name: "False",
},
&ruleRefExpr{
	pos: position{line: 221, col: 5, offset: 5960},
	name: "Infinity",
},
&ruleRefExpr{
	pos: position{line: 221, col: 16, offset: 5971},
	name: "NaN",
},
&ruleRefExpr{
	pos: position{line: 222, col: 5, offset: 5979},
	name: "Merge",
},
&ruleRefExpr{
	pos: position{line: 222, col: 13, offset: 5987},
	name: "Some",
},
	},
},
},
{
	name: "Optional",
	pos: position{line: 224, col: 1, offset: 5993},
	expr: &litMatcher{
	pos: position{line: 224, col: 12, offset: 6006},
	val: "Optional",
	ignoreCase: false,
},
},
{
	name: "Text",
	pos: position{line: 225, col: 1, offset: 6017},
	expr: &litMatcher{
	pos: position{line: 225, col: 8, offset: 6026},
	val: "Text",
	ignoreCase: false,
},
},
{
	name: "List",
	pos: position{line: 226, col: 1, offset: 6033},
	expr: &litMatcher{
	pos: position{line: 226, col: 8, offset: 6042},
	val: "List",
	ignoreCase: false,
},
},
{
	name: "Combine",
	pos: position{line: 228, col: 1, offset: 6050},
	expr: &choiceExpr{
	pos: position{line: 228, col: 11, offset: 6062},
	alternatives: []interface{}{
&litMatcher{
	pos: position{line: 228, col: 11, offset: 6062},
	val: "/\\",
	ignoreCase: false,
},
&litMatcher{
	pos: position{line: 228, col: 19, offset: 6070},
	val: "∧",
	ignoreCase: false,
},
	},
},
},
{
	name: "CombineTypes",
	pos: position{line: 229, col: 1, offset: 6076},
	expr: &choiceExpr{
	pos: position{line: 229, col: 16, offset: 6093},
	alternatives: []interface{}{
&litMatcher{
	pos: position{line: 229, col: 16, offset: 6093},
	val: "//\\\\",
	ignoreCase: false,
},
&litMatcher{
	pos: position{line: 229, col: 27, offset: 6104},
	val: "⩓",
	ignoreCase: false,
},
	},
},
},
{
	name: "Prefer",
	pos: position{line: 230, col: 1, offset: 6110},
	expr: &choiceExpr{
	pos: position{line: 230, col: 10, offset: 6121},
	alternatives: []interface{}{
&litMatcher{
	pos: position{line: 230, col: 10, offset: 6121},
	val: "//",
	ignoreCase: false,
},
&litMatcher{
	pos: position{line: 230, col: 17, offset: 6128},
	val: "⫽",
	ignoreCase: false,
},
	},
},
},
{
	name: "Lambda",
	pos: position{line: 231, col: 1, offset: 6134},
	expr: &choiceExpr{
	pos: position{line: 231, col: 10, offset: 6145},
	alternatives: []interface{}{
&litMatcher{
	pos: position{line: 231, col: 10, offset: 6145},
	val: "\\",
	ignoreCase: false,
},
&litMatcher{
	pos: position{line: 231, col: 17, offset: 6152},
	val: "λ",
	ignoreCase: false,
},
	},
},
},
{
	name: "Forall",
	pos: position{line: 232, col: 1, offset: 6157},
	expr: &choiceExpr{
	pos: position{line: 232, col: 10, offset: 6168},
	alternatives: []interface{}{
&litMatcher{
	pos: position{line: 232, col: 10, offset: 6168},
	val: "forall",
	ignoreCase: false,
},
&litMatcher{
	pos: position{line: 232, col: 21, offset: 6179},
	val: "∀",
	ignoreCase: false,
},
	},
},
},
{
	name: "Arrow",
	pos: position{line: 233, col: 1, offset: 6185},
	expr: &choiceExpr{
	pos: position{line: 233, col: 9, offset: 6195},
	alternatives: []interface{}{
&litMatcher{
	pos: position{line: 233, col: 9, offset: 6195},
	val: "->",
	ignoreCase: false,
},
&litMatcher{
	pos: position{line: 233, col: 16, offset: 6202},
	val: "→",
	ignoreCase: false,
},
	},
},
},
{
	name: "Exponent",
	pos: position{line: 235, col: 1, offset: 6209},
	expr: &seqExpr{
	pos: position{line: 235, col: 12, offset: 6222},
	exprs: []interface{}{
&litMatcher{
	pos: position{line: 235, col: 12, offset: 6222},
	val: "e",
	ignoreCase: true,
},
&zeroOrOneExpr{
	pos: position{line: 235, col: 17, offset: 6227},
	expr: &charClassMatcher{
	pos: position{line: 235, col: 17, offset: 6227},
	val: "[+-]",
	chars: []rune{'+','-',},
	ignoreCase: false,
	inverted: false,
},
},
&oneOrMoreExpr{
	pos: position{line: 235, col: 23, offset: 6233},
	expr: &ruleRefExpr{
	pos: position{line: 235, col: 23, offset: 6233},
	name: "Digit",
},
},
	},
},
},
{
	name: "NumericDoubleLiteral",
	pos: position{line: 237, col: 1, offset: 6241},
	expr: &actionExpr{
	pos: position{line: 237, col: 24, offset: 6266},
	run: (*parser).callonNumericDoubleLiteral1,
	expr: &seqExpr{
	pos: position{line: 237, col: 24, offset: 6266},
	exprs: []interface{}{
&zeroOrOneExpr{
	pos: position{line: 237, col: 24, offset: 6266},
	expr: &charClassMatcher{
	pos: position{line: 237, col: 24, offset: 6266},
	val: "[+-]",
	chars: []rune{'+','-',},
	ignoreCase: false,
	inverted: false,
},
},
&oneOrMoreExpr{
	pos: position{line: 237, col: 30, offset: 6272},
	expr: &ruleRefExpr{
	pos: position{line: 237, col: 30, offset: 6272},
	name: "Digit",
},
},
&choiceExpr{
	pos: position{line: 237, col: 39, offset: 6281},
	alternatives: []interface{}{
&seqExpr{
	pos: position{line: 237, col: 39, offset: 6281},
	exprs: []interface{}{
&litMatcher{
	pos: position{line: 237, col: 39, offset: 6281},
	val: ".",
	ignoreCase: false,
},
&oneOrMoreExpr{
	pos: position{line: 237, col: 43, offset: 6285},
	expr: &ruleRefExpr{
	pos: position{line: 237, col: 43, offset: 6285},
	name: "Digit",
},
},
&zeroOrOneExpr{
	pos: position{line: 237, col: 50, offset: 6292},
	expr: &ruleRefExpr{
	pos: position{line: 237, col: 50, offset: 6292},
	name: "Exponent",
},
},
	},
},
&ruleRefExpr{
	pos: position{line: 237, col: 62, offset: 6304},
	name: "Exponent",
},
	},
},
	},
},
},
},
{
	name: "DoubleLiteral",
	pos: position{line: 245, col: 1, offset: 6460},
	expr: &choiceExpr{
	pos: position{line: 245, col: 17, offset: 6478},
	alternatives: []interface{}{
&labeledExpr{
	pos: position{line: 245, col: 17, offset: 6478},
	label: "d",
	expr: &ruleRefExpr{
	pos: position{line: 245, col: 19, offset: 6480},
	name: "NumericDoubleLiteral",
},
},
&actionExpr{
	pos: position{line: 246, col: 5, offset: 6505},
	run: (*parser).callonDoubleLiteral4,
	expr: &ruleRefExpr{
	pos: position{line: 246, col: 5, offset: 6505},
	name: "Infinity",
},
},
&actionExpr{
	pos: position{line: 247, col: 5, offset: 6557},
	run: (*parser).callonDoubleLiteral6,
	expr: &seqExpr{
	pos: position{line: 247, col: 5, offset: 6557},
	exprs: []interface{}{
&litMatcher{
	pos: position{line: 247, col: 5, offset: 6557},
	val: "-",
	ignoreCase: false,
},
&ruleRefExpr{
	pos: position{line: 247, col: 9, offset: 6561},
	name: "Infinity",
},
	},
},
},
&actionExpr{
	pos: position{line: 248, col: 5, offset: 6614},
	run: (*parser).callonDoubleLiteral10,
	expr: &ruleRefExpr{
	pos: position{line: 248, col: 5, offset: 6614},
	name: "NaN",
},
},
	},
},
},
{
	name: "NaturalLiteral",
	pos: position{line: 250, col: 1, offset: 6657},
	expr: &actionExpr{
	pos: position{line: 250, col: 18, offset: 6676},
	run: (*parser).callonNaturalLiteral1,
	expr: &oneOrMoreExpr{
	pos: position{line: 250, col: 18, offset: 6676},
	expr: &ruleRefExpr{
	pos: position{line: 250, col: 18, offset: 6676},
	name: "Digit",
},
},
},
},
{
	name: "IntegerLiteral",
	pos: position{line: 255, col: 1, offset: 6765},
	expr: &actionExpr{
	pos: position{line: 255, col: 18, offset: 6784},
	run: (*parser).callonIntegerLiteral1,
	expr: &seqExpr{
	pos: position{line: 255, col: 18, offset: 6784},
	exprs: []interface{}{
&charClassMatcher{
	pos: position{line: 255, col: 18, offset: 6784},
	val: "[+-]",
	chars: []rune{'+','-',},
	ignoreCase: false,
	inverted: false,
},
&ruleRefExpr{
	pos: position{line: 255, col: 22, offset: 6788},
	name: "NaturalLiteral",
},
	},
},
},
},
{
	name: "DeBruijn",
	pos: position{line: 263, col: 1, offset: 6940},
	expr: &actionExpr{
	pos: position{line: 263, col: 12, offset: 6953},
	run: (*parser).callonDeBruijn1,
	expr: &seqExpr{
	pos: position{line: 263, col: 12, offset: 6953},
	exprs: []interface{}{
&ruleRefExpr{
	pos: position{line: 263, col: 12, offset: 6953},
	name: "_",
},
&litMatcher{
	pos: position{line: 263, col: 14, offset: 6955},
	val: "@",
	ignoreCase: false,
},
&ruleRefExpr{
	pos: position{line: 263, col: 18, offset: 6959},
	name: "_",
},
&labeledExpr{
	pos: position{line: 263, col: 20, offset: 6961},
	label: "index",
	expr: &ruleRefExpr{
	pos: position{line: 263, col: 26, offset: 6967},
	name: "NaturalLiteral",
},
},
	},
},
},
},
{
	name: "Variable",
	pos: position{line: 265, col: 1, offset: 7023},
	expr: &actionExpr{
	pos: position{line: 265, col: 12, offset: 7036},
	run: (*parser).callonVariable1,
	expr: &seqExpr{
	pos: position{line: 265, col: 12, offset: 7036},
	exprs: []interface{}{
&labeledExpr{
	pos: position{line: 265, col: 12, offset: 7036},
	label: "name",
	expr: &ruleRefExpr{
	pos: position{line: 265, col: 17, offset: 7041},
	name: "NonreservedLabel",
},
},
&labeledExpr{
	pos: position{line: 265, col: 34, offset: 7058},
	label: "index",
	expr: &zeroOrOneExpr{
	pos: position{line: 265, col: 40, offset: 7064},
	expr: &ruleRefExpr{
	pos: position{line: 265, col: 40, offset: 7064},
	name: "DeBruijn",
},
},
},
	},
},
},
},
{
	name: "Identifier",
	pos: position{line: 273, col: 1, offset: 7227},
	expr: &choiceExpr{
	pos: position{line: 273, col: 14, offset: 7242},
	alternatives: []interface{}{
&ruleRefExpr{
	pos: position{line: 273, col: 14, offset: 7242},
	name: "Variable",
},
&ruleRefExpr{
	pos: position{line: 273, col: 25, offset: 7253},
	name: "Reserved",
},
	},
},
},
{
	name: "PathCharacter",
	pos: position{line: 275, col: 1, offset: 7263},
	expr: &choiceExpr{
	pos: position{line: 276, col: 6, offset: 7286},
	alternatives: []interface{}{
&litMatcher{
	pos: position{line: 276, col: 6, offset: 7286},
	val: "!",
	ignoreCase: false,
},
&charClassMatcher{
	pos: position{line: 277, col: 6, offset: 7298},
	val: "[\\x24-\\x27]",
	ranges: []rune{'$','\'',},
	ignoreCase: false,
	inverted: false,
},
&charClassMatcher{
	pos: position{line: 278, col: 6, offset: 7315},
	val: "[\\x2a-\\x2b]",
	ranges: []rune{'*','+',},
	ignoreCase: false,
	inverted: false,
},
&charClassMatcher{
	pos: position{line: 279, col: 6, offset: 7332},
	val: "[\\x2d-\\x2e]",
	ranges: []rune{'-','.',},
	ignoreCase: false,
	inverted: false,
},
&charClassMatcher{
	pos: position{line: 280, col: 6, offset: 7349},
	val: "[\\x30-\\x3b]",
	ranges: []rune{'0',';',},
	ignoreCase: false,
	inverted: false,
},
&litMatcher{
	pos: position{line: 281, col: 6, offset: 7366},
	val: "=",
	ignoreCase: false,
},
&charClassMatcher{
	pos: position{line: 282, col: 6, offset: 7378},
	val: "[\\x40-\\x5a]",
	ranges: []rune{'@','Z',},
	ignoreCase: false,
	inverted: false,
},
&charClassMatcher{
	pos: position{line: 283, col: 6, offset: 7395},
	val: "[\\x5e-\\x7a]",
	ranges: []rune{'^','z',},
	ignoreCase: false,
	inverted: false,
},
&litMatcher{
	pos: position{line: 284, col: 6, offset: 7412},
	val: "|",
	ignoreCase: false,
},
&litMatcher{
	pos: position{line: 285, col: 6, offset: 7424},
	val: "~",
	ignoreCase: false,
},
	},
},
},
{
	name: "UnquotedPathComponent",
	pos: position{line: 287, col: 1, offset: 7432},
	expr: &actionExpr{
	pos: position{line: 287, col: 25, offset: 7458},
	run: (*parser).callonUnquotedPathComponent1,
	expr: &oneOrMoreExpr{
	pos: position{line: 287, col: 25, offset: 7458},
	expr: &ruleRefExpr{
	pos: position{line: 287, col: 25, offset: 7458},
	name: "PathCharacter",
},
},
},
},
{
	name: "PathComponent",
	pos: position{line: 289, col: 1, offset: 7505},
	expr: &actionExpr{
	pos: position{line: 289, col: 17, offset: 7523},
	run: (*parser).callonPathComponent1,
	expr: &seqExpr{
	pos: position{line: 289, col: 17, offset: 7523},
	exprs: []interface{}{
&litMatcher{
	pos: position{line: 289, col: 17, offset: 7523},
	val: "/",
	ignoreCase: false,
},
&labeledExpr{
	pos: position{line: 289, col: 21, offset: 7527},
	label: "u",
	expr: &ruleRefExpr{
	pos: position{line: 289, col: 23, offset: 7529},
	name: "UnquotedPathComponent",
},
},
	},
},
},
},
{
	name: "Path",
	pos: position{line: 291, col: 1, offset: 7570},
	expr: &actionExpr{
	pos: position{line: 291, col: 8, offset: 7579},
	run: (*parser).callonPath1,
	expr: &labeledExpr{
	pos: position{line: 291, col: 8, offset: 7579},
	label: "cs",
	expr: &oneOrMoreExpr{
	pos: position{line: 291, col: 11, offset: 7582},
	expr: &ruleRefExpr{
	pos: position{line: 291, col: 11, offset: 7582},
	name: "PathComponent",
},
},
},
},
},
{
	name: "Local",
	pos: position{line: 300, col: 1, offset: 7856},
	expr: &choiceExpr{
	pos: position{line: 300, col: 9, offset: 7866},
	alternatives: []interface{}{
&ruleRefExpr{
	pos: position{line: 300, col: 9, offset: 7866},
	name: "ParentPath",
},
&ruleRefExpr{
	pos: position{line: 300, col: 22, offset: 7879},
	name: "HerePath",
},
&ruleRefExpr{
	pos: position{line: 300, col: 33, offset: 7890},
	name: "HomePath",
},
&ruleRefExpr{
	pos: position{line: 300, col: 44, offset: 7901},
	name: "AbsolutePath",
},
	},
},
},
{
	name: "ParentPath",
	pos: position{line: 302, col: 1, offset: 7915},
	expr: &actionExpr{
	pos: position{line: 302, col: 14, offset: 7930},
	run: (*parser).callonParentPath1,
	expr: &seqExpr{
	pos: position{line: 302, col: 14, offset: 7930},
	exprs: []interface{}{
&litMatcher{
	pos: position{line: 302, col: 14, offset: 7930},
	val: "..",
	ignoreCase: false,
},
&labeledExpr{
	pos: position{line: 302, col: 19, offset: 7935},
	label: "p",
	expr: &ruleRefExpr{
	pos: position{line: 302, col: 21, offset: 7937},
	name: "Path",
},
},
	},
},
},
},
{
	name: "HerePath",
	pos: position{line: 303, col: 1, offset: 7993},
	expr: &actionExpr{
	pos: position{line: 303, col: 12, offset: 8006},
	run: (*parser).callonHerePath1,
	expr: &seqExpr{
	pos: position{line: 303, col: 12, offset: 8006},
	exprs: []interface{}{
&litMatcher{
	pos: position{line: 303, col: 12, offset: 8006},
	val: ".",
	ignoreCase: false,
},
&labeledExpr{
	pos: position{line: 303, col: 16, offset: 8010},
	label: "p",
	expr: &ruleRefExpr{
	pos: position{line: 303, col: 18, offset: 8012},
	name: "Path",
},
},
	},
},
},
},
{
	name: "HomePath",
	pos: position{line: 304, col: 1, offset: 8051},
	expr: &actionExpr{
	pos: position{line: 304, col: 12, offset: 8064},
	run: (*parser).callonHomePath1,
	expr: &seqExpr{
	pos: position{line: 304, col: 12, offset: 8064},
	exprs: []interface{}{
&litMatcher{
	pos: position{line: 304, col: 12, offset: 8064},
	val: "~",
	ignoreCase: false,
},
&labeledExpr{
	pos: position{line: 304, col: 16, offset: 8068},
	label: "p",
	expr: &ruleRefExpr{
	pos: position{line: 304, col: 18, offset: 8070},
	name: "Path",
},
},
	},
},
},
},
{
	name: "AbsolutePath",
	pos: position{line: 305, col: 1, offset: 8125},
	expr: &actionExpr{
	pos: position{line: 305, col: 16, offset: 8142},
	run: (*parser).callonAbsolutePath1,
	expr: &labeledExpr{
	pos: position{line: 305, col: 16, offset: 8142},
	label: "p",
	expr: &ruleRefExpr{
	pos: position{line: 305, col: 18, offset: 8144},
	name: "Path",
},
},
},
},
{
	name: "Scheme",
	pos: position{line: 307, col: 1, offset: 8200},
	expr: &seqExpr{
	pos: position{line: 307, col: 10, offset: 8211},
	exprs: []interface{}{
&litMatcher{
	pos: position{line: 307, col: 10, offset: 8211},
	val: "http",
	ignoreCase: false,
},
&zeroOrOneExpr{
	pos: position{line: 307, col: 17, offset: 8218},
	expr: &litMatcher{
	pos: position{line: 307, col: 17, offset: 8218},
	val: "s",
	ignoreCase: false,
},
},
	},
},
},
{
	name: "HttpRaw",
	pos: position{line: 309, col: 1, offset: 8224},
	expr: &actionExpr{
	pos: position{line: 309, col: 11, offset: 8236},
	run: (*parser).callonHttpRaw1,
	expr: &seqExpr{
	pos: position{line: 309, col: 11, offset: 8236},
	exprs: []interface{}{
&ruleRefExpr{
	pos: position{line: 309, col: 11, offset: 8236},
	name: "Scheme",
},
&litMatcher{
	pos: position{line: 309, col: 18, offset: 8243},
	val: "://",
	ignoreCase: false,
},
&ruleRefExpr{
	pos: position{line: 309, col: 24, offset: 8249},
	name: "Authority",
},
&ruleRefExpr{
	pos: position{line: 309, col: 34, offset: 8259},
	name: "Path",
},
&zeroOrOneExpr{
	pos: position{line: 309, col: 39, offset: 8264},
	expr: &seqExpr{
	pos: position{line: 309, col: 41, offset: 8266},
	exprs: []interface{}{
&litMatcher{
	pos: position{line: 309, col: 41, offset: 8266},
	val: "?",
	ignoreCase: false,
},
&ruleRefExpr{
	pos: position{line: 309, col: 45, offset: 8270},
	name: "Query",
},
	},
},
},
	},
},
},
},
{
	name: "Authority",
	pos: position{line: 311, col: 1, offset: 8327},
	expr: &seqExpr{
	pos: position{line: 311, col: 13, offset: 8341},
	exprs: []interface{}{
&zeroOrOneExpr{
	pos: position{line: 311, col: 13, offset: 8341},
	expr: &seqExpr{
	pos: position{line: 311, col: 14, offset: 8342},
	exprs: []interface{}{
&ruleRefExpr{
	pos: position{line: 311, col: 14, offset: 8342},
	name: "Userinfo",
},
&litMatcher{
	pos: position{line: 311, col: 23, offset: 8351},
	val: "@",
	ignoreCase: false,
},
	},
},
},
&ruleRefExpr{
	pos: position{line: 311, col: 29, offset: 8357},
	name: "Host",
},
&zeroOrOneExpr{
	pos: position{line: 311, col: 34, offset: 8362},
	expr: &seqExpr{
	pos: position{line: 311, col: 35, offset: 8363},
	exprs: []interface{}{
&litMatcher{
	pos: position{line: 311, col: 35, offset: 8363},
	val: ":",
	ignoreCase: false,
},
&ruleRefExpr{
	pos: position{line: 311, col: 39, offset: 8367},
	name: "Port",
},
	},
},
},
	},
},
},
{
	name: "Userinfo",
	pos: position{line: 313, col: 1, offset: 8375},
	expr: &zeroOrMoreExpr{
	pos: position{line: 313, col: 12, offset: 8388},
	expr: &choiceExpr{
	pos: position{line: 313, col: 14, offset: 8390},
	alternatives: []interface{}{
&ruleRefExpr{
	pos: position{line: 313, col: 14, offset: 8390},
	name: "Unreserved",
},
&ruleRefExpr{
	pos: position{line: 313, col: 27, offset: 8403},
	name: "PctEncoded",
},
&ruleRefExpr{
	pos: position{line: 313, col: 40, offset: 8416},
	name: "SubDelims",
},
&litMatcher{
	pos: position{line: 313, col: 52, offset: 8428},
	val: ":",
	ignoreCase: false,
},
	},
},
},
},
{
	name: "Host",
	pos: position{line: 315, col: 1, offset: 8436},
	expr: &choiceExpr{
	pos: position{line: 315, col: 8, offset: 8445},
	alternatives: []interface{}{
&ruleRefExpr{
	pos: position{line: 315, col: 8, offset: 8445},
	name: "IPLiteral",
},
&ruleRefExpr{
	pos: position{line: 315, col: 20, offset: 8457},
	name: "RegName",
},
	},
},
},
{
	name: "Port",
	pos: position{line: 317, col: 1, offset: 8466},
	expr: &zeroOrMoreExpr{
	pos: position{line: 317, col: 8, offset: 8475},
	expr: &ruleRefExpr{
	pos: position{line: 317, col: 8, offset: 8475},
	name: "Digit",
},
},
},
{
	name: "IPLiteral",
	pos: position{line: 319, col: 1, offset: 8483},
	expr: &seqExpr{
	pos: position{line: 319, col: 13, offset: 8497},
	exprs: []interface{}{
&litMatcher{
	pos: position{line: 319, col: 13, offset: 8497},
	val: "[",
	ignoreCase: false,
},
&ruleRefExpr{
	pos: position{line: 319, col: 17, offset: 8501},
	name: "IPv6address",
},
&litMatcher{
	pos: position{line: 319, col: 29, offset: 8513},
	val: "]",
	ignoreCase: false,
},
	},
},
},
{
	name: "IPv6address",
	pos: position{line: 321, col: 1, offset: 8518},
	expr: &actionExpr{
	pos: position{line: 321, col: 15, offset: 8534},
	run: (*parser).callonIPv6address1,
	expr: &seqExpr{
	pos: position{line: 321, col: 15, offset: 8534},
	exprs: []interface{}{
&zeroOrMoreExpr{
	pos: position{line: 321, col: 15, offset: 8534},
	expr: &ruleRefExpr{
	pos: position{line: 321, col: 16, offset: 8535},
	name: "HexDig",
},
},
&litMatcher{
	pos: position{line: 321, col: 25, offset: 8544},
	val: ":",
	ignoreCase: false,
},
&zeroOrMoreExpr{
	pos: position{line: 321, col: 29, offset: 8548},
	expr: &choiceExpr{
	pos: position{line: 321, col: 30, offset: 8549},
	alternatives: []interface{}{
&ruleRefExpr{
	pos: position{line: 321, col: 30, offset: 8549},
	name: "HexDig",
},
&litMatcher{
	pos: position{line: 321, col: 39, offset: 8558},
	val: ":",
	ignoreCase: false,
},
&litMatcher{
	pos: position{line: 321, col: 45, offset: 8564},
	val: ".",
	ignoreCase: false,
},
	},
},
},
	},
},
},
},
{
	name: "RegName",
	pos: position{line: 327, col: 1, offset: 8718},
	expr: &zeroOrMoreExpr{
	pos: position{line: 327, col: 11, offset: 8730},
	expr: &choiceExpr{
	pos: position{line: 327, col: 12, offset: 8731},
	alternatives: []interface{}{
&ruleRefExpr{
	pos: position{line: 327, col: 12, offset: 8731},
	name: "Unreserved",
},
&ruleRefExpr{
	pos: position{line: 327, col: 25, offset: 8744},
	name: "PctEncoded",
},
&ruleRefExpr{
	pos: position{line: 327, col: 38, offset: 8757},
	name: "SubDelims",
},
	},
},
},
},
{
	name: "PChar",
	pos: position{line: 329, col: 1, offset: 8770},
	expr: &choiceExpr{
	pos: position{line: 329, col: 9, offset: 8780},
	alternatives: []interface{}{
&ruleRefExpr{
	pos: position{line: 329, col: 9, offset: 8780},
	name: "Unreserved",
},
&ruleRefExpr{
	pos: position{line: 329, col: 22, offset: 8793},
	name: "PctEncoded",
},
&ruleRefExpr{
	pos: position{line: 329, col: 35, offset: 8806},
	name: "SubDelims",
},
&charClassMatcher{
	pos: position{line: 329, col: 47, offset: 8818},
	val: "[:@]",
	chars: []rune{':','@',},
	ignoreCase: false,
	inverted: false,
},
	},
},
},
{
	name: "Query",
	pos: position{line: 331, col: 1, offset: 8824},
	expr: &zeroOrMoreExpr{
	pos: position{line: 331, col: 9, offset: 8834},
	expr: &choiceExpr{
	pos: position{line: 331, col: 10, offset: 8835},
	alternatives: []interface{}{
&ruleRefExpr{
	pos: position{line: 331, col: 10, offset: 8835},
	name: "PChar",
},
&charClassMatcher{
	pos: position{line: 331, col: 18, offset: 8843},
	val: "[/?]",
	chars: []rune{'/','?',},
	ignoreCase: false,
	inverted: false,
},
	},
},
},
},
{
	name: "PctEncoded",
	pos: position{line: 333, col: 1, offset: 8851},
	expr: &seqExpr{
	pos: position{line: 333, col: 14, offset: 8866},
	exprs: []interface{}{
&litMatcher{
	pos: position{line: 333, col: 14, offset: 8866},
	val: "%",
	ignoreCase: false,
},
&ruleRefExpr{
	pos: position{line: 333, col: 18, offset: 8870},
	name: "HexDig",
},
&ruleRefExpr{
	pos: position{line: 333, col: 25, offset: 8877},
	name: "HexDig",
},
	},
},
},
{
	name: "Unreserved",
	pos: position{line: 335, col: 1, offset: 8885},
	expr: &charClassMatcher{
	pos: position{line: 335, col: 14, offset: 8900},
	val: "[A-Za-z0-9._~-]",
	chars: []rune{'.','_','~','-',},
	ranges: []rune{'A','Z','a','z','0','9',},
	ignoreCase: false,
	inverted: false,
},
},
{
	name: "SubDelims",
	pos: position{line: 337, col: 1, offset: 8917},
	expr: &choiceExpr{
	pos: position{line: 337, col: 13, offset: 8931},
	alternatives: []interface{}{
&litMatcher{
	pos: position{line: 337, col: 13, offset: 8931},
	val: "!",
	ignoreCase: false,
},
&litMatcher{
	pos: position{line: 337, col: 19, offset: 8937},
	val: "$",
	ignoreCase: false,
},
&litMatcher{
	pos: position{line: 337, col: 25, offset: 8943},
	val: "&",
	ignoreCase: false,
},
&litMatcher{
	pos: position{line: 337, col: 31, offset: 8949},
	val: "'",
	ignoreCase: false,
},
&litMatcher{
	pos: position{line: 337, col: 37, offset: 8955},
	val: "(",
	ignoreCase: false,
},
&litMatcher{
	pos: position{line: 337, col: 43, offset: 8961},
	val: ")",
	ignoreCase: false,
},
&litMatcher{
	pos: position{line: 337, col: 49, offset: 8967},
	val: "*",
	ignoreCase: false,
},
&litMatcher{
	pos: position{line: 337, col: 55, offset: 8973},
	val: "+",
	ignoreCase: false,
},
&litMatcher{
	pos: position{line: 337, col: 61, offset: 8979},
	val: ",",
	ignoreCase: false,
},
&litMatcher{
	pos: position{line: 337, col: 67, offset: 8985},
	val: ";",
	ignoreCase: false,
},
&litMatcher{
	pos: position{line: 337, col: 73, offset: 8991},
	val: "=",
	ignoreCase: false,
},
	},
},
},
{
	name: "Http",
	pos: position{line: 339, col: 1, offset: 8996},
	expr: &actionExpr{
	pos: position{line: 339, col: 8, offset: 9005},
	run: (*parser).callonHttp1,
	expr: &labeledExpr{
	pos: position{line: 339, col: 8, offset: 9005},
	label: "u",
	expr: &ruleRefExpr{
	pos: position{line: 339, col: 10, offset: 9007},
	name: "HttpRaw",
},
},
},
},
{
	name: "Env",
	pos: position{line: 341, col: 1, offset: 9052},
	expr: &actionExpr{
	pos: position{line: 341, col: 7, offset: 9060},
	run: (*parser).callonEnv1,
	expr: &seqExpr{
	pos: position{line: 341, col: 7, offset: 9060},
	exprs: []interface{}{
&litMatcher{
	pos: position{line: 341, col: 7, offset: 9060},
	val: "env:",
	ignoreCase: false,
},
&labeledExpr{
	pos: position{line: 341, col: 14, offset: 9067},
	label: "v",
	expr: &choiceExpr{
	pos: position{line: 341, col: 17, offset: 9070},
	alternatives: []interface{}{
&ruleRefExpr{
	pos: position{line: 341, col: 17, offset: 9070},
	name: "BashEnvironmentVariable",
},
&ruleRefExpr{
	pos: position{line: 341, col: 43, offset: 9096},
	name: "PosixEnvironmentVariable",
},
	},
},
},
	},
},
},
},
{
	name: "BashEnvironmentVariable",
	pos: position{line: 343, col: 1, offset: 9141},
	expr: &actionExpr{
	pos: position{line: 343, col: 27, offset: 9169},
	run: (*parser).callonBashEnvironmentVariable1,
	expr: &seqExpr{
	pos: position{line: 343, col: 27, offset: 9169},
	exprs: []interface{}{
&charClassMatcher{
	pos: position{line: 343, col: 27, offset: 9169},
	val: "[A-Za-z_]",
	chars: []rune{'_',},
	ranges: []rune{'A','Z','a','z',},
	ignoreCase: false,
	inverted: false,
},
&zeroOrMoreExpr{
	pos: position{line: 343, col: 36, offset: 9178},
	expr: &charClassMatcher{
	pos: position{line: 343, col: 36, offset: 9178},
	val: "[A-Za-z0-9_]",
	chars: []rune{'_',},
	ranges: []rune{'A','Z','a','z','0','9',},
	ignoreCase: false,
	inverted: false,
},
},
	},
},
},
},
{
	name: "PosixEnvironmentVariable",
	pos: position{line: 347, col: 1, offset: 9234},
	expr: &actionExpr{
	pos: position{line: 347, col: 28, offset: 9263},
	run: (*parser).callonPosixEnvironmentVariable1,
	expr: &seqExpr{
	pos: position{line: 347, col: 28, offset: 9263},
	exprs: []interface{}{
&litMatcher{
	pos: position{line: 347, col: 28, offset: 9263},
	val: "\"",
	ignoreCase: false,
},
&labeledExpr{
	pos: position{line: 347, col: 32, offset: 9267},
	label: "v",
	expr: &ruleRefExpr{
	pos: position{line: 347, col: 34, offset: 9269},
	name: "PosixEnvironmentVariableContent",
},
},
&litMatcher{
	pos: position{line: 347, col: 66, offset: 9301},
	val: "\"",
	ignoreCase: false,
},
	},
},
},
},
{
	name: "PosixEnvironmentVariableContent",
	pos: position{line: 351, col: 1, offset: 9326},
	expr: &actionExpr{
	pos: position{line: 351, col: 35, offset: 9362},
	run: (*parser).callonPosixEnvironmentVariableContent1,
	expr: &labeledExpr{
	pos: position{line: 351, col: 35, offset: 9362},
	label: "v",
	expr: &oneOrMoreExpr{
	pos: position{line: 351, col: 37, offset: 9364},
	expr: &ruleRefExpr{
	pos: position{line: 351, col: 37, offset: 9364},
	name: "PosixEnvironmentVariableCharacter",
},
},
},
},
},
{
	name: "PosixEnvironmentVariableCharacter",
	pos: position{line: 360, col: 1, offset: 9577},
	expr: &choiceExpr{
	pos: position{line: 361, col: 7, offset: 9621},
	alternatives: []interface{}{
&actionExpr{
	pos: position{line: 361, col: 7, offset: 9621},
	run: (*parser).callonPosixEnvironmentVariableCharacter2,
	expr: &litMatcher{
	pos: position{line: 361, col: 7, offset: 9621},
	val: "\\\"",
	ignoreCase: false,
},
},
&actionExpr{
	pos: position{line: 362, col: 7, offset: 9661},
	run: (*parser).callonPosixEnvironmentVariableCharacter4,
	expr: &litMatcher{
	pos: position{line: 362, col: 7, offset: 9661},
	val: "\\\\",
	ignoreCase: false,
},
},
&actionExpr{
	pos: position{line: 363, col: 7, offset: 9701},
	run: (*parser).callonPosixEnvironmentVariableCharacter6,
	expr: &litMatcher{
	pos: position{line: 363, col: 7, offset: 9701},
	val: "\\a",
	ignoreCase: false,
},
},
&actionExpr{
	pos: position{line: 364, col: 7, offset: 9741},
	run: (*parser).callonPosixEnvironmentVariableCharacter8,
	expr: &litMatcher{
	pos: position{line: 364, col: 7, offset: 9741},
	val: "\\b",
	ignoreCase: false,
},
},
&actionExpr{
	pos: position{line: 365, col: 7, offset: 9781},
	run: (*parser).callonPosixEnvironmentVariableCharacter10,
	expr: &litMatcher{
	pos: position{line: 365, col: 7, offset: 9781},
	val: "\\f",
	ignoreCase: false,
},
},
&actionExpr{
	pos: position{line: 366, col: 7, offset: 9821},
	run: (*parser).callonPosixEnvironmentVariableCharacter12,
	expr: &litMatcher{
	pos: position{line: 366, col: 7, offset: 9821},
	val: "\\n",
	ignoreCase: false,
},
},
&actionExpr{
	pos: position{line: 367, col: 7, offset: 9861},
	run: (*parser).callonPosixEnvironmentVariableCharacter14,
	expr: &litMatcher{
	pos: position{line: 367, col: 7, offset: 9861},
	val: "\\r",
	ignoreCase: false,
},
},
&actionExpr{
	pos: position{line: 368, col: 7, offset: 9901},
	run: (*parser).callonPosixEnvironmentVariableCharacter16,
	expr: &litMatcher{
	pos: position{line: 368, col: 7, offset: 9901},
	val: "\\t",
	ignoreCase: false,
},
},
&actionExpr{
	pos: position{line: 369, col: 7, offset: 9941},
	run: (*parser).callonPosixEnvironmentVariableCharacter18,
	expr: &litMatcher{
	pos: position{line: 369, col: 7, offset: 9941},
	val: "\\v",
	ignoreCase: false,
},
},
&charClassMatcher{
	pos: position{line: 370, col: 7, offset: 9981},
	val: "[\\x20-\\x21]",
	ranges: []rune{' ','!',},
	ignoreCase: false,
	inverted: false,
},
&charClassMatcher{
	pos: position{line: 371, col: 7, offset: 9999},
	val: "[\\x23-\\x3c]",
	ranges: []rune{'#','<',},
	ignoreCase: false,
	inverted: false,
},
&charClassMatcher{
	pos: position{line: 372, col: 7, offset: 10017},
	val: "[\\x3e-\\x5b]",
	ranges: []rune{'>','[',},
	ignoreCase: false,
	inverted: false,
},
&charClassMatcher{
	pos: position{line: 373, col: 7, offset: 10035},
	val: "[\\x5d-\\x7e]",
	ranges: []rune{']','~',},
	ignoreCase: false,
	inverted: false,
},
	},
},
},
{
	name: "ImportType",
	pos: position{line: 375, col: 1, offset: 10048},
	expr: &choiceExpr{
	pos: position{line: 375, col: 14, offset: 10063},
	alternatives: []interface{}{
&ruleRefExpr{
	pos: position{line: 375, col: 14, offset: 10063},
	name: "Missing",
},
&ruleRefExpr{
	pos: position{line: 375, col: 24, offset: 10073},
	name: "Local",
},
&ruleRefExpr{
	pos: position{line: 375, col: 32, offset: 10081},
	name: "Http",
},
&ruleRefExpr{
	pos: position{line: 375, col: 39, offset: 10088},
	name: "Env",
},
	},
},
},
{
	name: "ImportHashed",
	pos: position{line: 377, col: 1, offset: 10093},
	expr: &actionExpr{
	pos: position{line: 377, col: 16, offset: 10110},
	run: (*parser).callonImportHashed1,
	expr: &labeledExpr{
	pos: position{line: 377, col: 16, offset: 10110},
	label: "i",
	expr: &ruleRefExpr{
	pos: position{line: 377, col: 18, offset: 10112},
	name: "ImportType",
},
},
},
},
{
	name: "Import",
	pos: position{line: 379, col: 1, offset: 10179},
	expr: &choiceExpr{
	pos: position{line: 379, col: 10, offset: 10190},
	alternatives: []interface{}{
&actionExpr{
	pos: position{line: 379, col: 10, offset: 10190},
	run: (*parser).callonImport2,
	expr: &seqExpr{
	pos: position{line: 379, col: 10, offset: 10190},
	exprs: []interface{}{
&labeledExpr{
	pos: position{line: 379, col: 10, offset: 10190},
	label: "i",
	expr: &ruleRefExpr{
	pos: position{line: 379, col: 12, offset: 10192},
	name: "ImportHashed",
},
},
&ruleRefExpr{
	pos: position{line: 379, col: 25, offset: 10205},
	name: "_",
},
&ruleRefExpr{
	pos: position{line: 379, col: 27, offset: 10207},
	name: "As",
},
&ruleRefExpr{
	pos: position{line: 379, col: 30, offset: 10210},
	name: "_1",
},
&ruleRefExpr{
	pos: position{line: 379, col: 33, offset: 10213},
	name: "Text",
},
	},
},
},
&actionExpr{
	pos: position{line: 380, col: 10, offset: 10310},
	run: (*parser).callonImport10,
	expr: &labeledExpr{
	pos: position{line: 380, col: 10, offset: 10310},
	label: "i",
	expr: &ruleRefExpr{
	pos: position{line: 380, col: 12, offset: 10312},
	name: "ImportHashed",
},
},
},
	},
},
},
{
	name: "LetBinding",
	pos: position{line: 383, col: 1, offset: 10407},
	expr: &actionExpr{
	pos: position{line: 383, col: 14, offset: 10422},
	run: (*parser).callonLetBinding1,
	expr: &seqExpr{
	pos: position{line: 383, col: 14, offset: 10422},
	exprs: []interface{}{
&ruleRefExpr{
	pos: position{line: 383, col: 14, offset: 10422},
	name: "Let",
},
&ruleRefExpr{
	pos: position{line: 383, col: 18, offset: 10426},
	name: "_1",
},
&labeledExpr{
	pos: position{line: 383, col: 21, offset: 10429},
	label: "label",
	expr: &ruleRefExpr{
	pos: position{line: 383, col: 27, offset: 10435},
	name: "NonreservedLabel",
},
},
&ruleRefExpr{
	pos: position{line: 383, col: 44, offset: 10452},
	name: "_",
},
&labeledExpr{
	pos: position{line: 383, col: 46, offset: 10454},
	label: "a",
	expr: &zeroOrOneExpr{
	pos: position{line: 383, col: 48, offset: 10456},
	expr: &seqExpr{
	pos: position{line: 383, col: 49, offset: 10457},
	exprs: []interface{}{
&ruleRefExpr{
	pos: position{line: 383, col: 49, offset: 10457},
	name: "Annotation",
},
&ruleRefExpr{
	pos: position{line: 383, col: 60, offset: 10468},
	name: "_",
},
	},
},
},
},
&litMatcher{
	pos: position{line: 384, col: 13, offset: 10484},
	val: "=",
	ignoreCase: false,
},
&ruleRefExpr{
	pos: position{line: 384, col: 17, offset: 10488},
	name: "_",
},
&labeledExpr{
	pos: position{line: 384, col: 19, offset: 10490},
	label: "v",
	expr: &ruleRefExpr{
	pos: position{line: 384, col: 21, offset: 10492},
	name: "Expression",
},
},
&ruleRefExpr{
	pos: position{line: 384, col: 32, offset: 10503},
	name: "_",
},
	},
},
},
},
{
	name: "Expression",
	pos: position{line: 399, col: 1, offset: 10812},
	expr: &choiceExpr{
	pos: position{line: 400, col: 7, offset: 10833},
	alternatives: []interface{}{
&actionExpr{
	pos: position{line: 400, col: 7, offset: 10833},
	run: (*parser).callonExpression2,
	expr: &seqExpr{
	pos: position{line: 400, col: 7, offset: 10833},
	exprs: []interface{}{
&ruleRefExpr{
	pos: position{line: 400, col: 7, offset: 10833},
	name: "Lambda",
},
&ruleRefExpr{
	pos: position{line: 400, col: 14, offset: 10840},
	name: "_",
},
&litMatcher{
	pos: position{line: 400, col: 16, offset: 10842},
	val: "(",
	ignoreCase: false,
},
&ruleRefExpr{
	pos: position{line: 400, col: 20, offset: 10846},
	name: "_",
},
&labeledExpr{
	pos: position{line: 400, col: 22, offset: 10848},
	label: "label",
	expr: &ruleRefExpr{
	pos: position{line: 400, col: 28, offset: 10854},
	name: "NonreservedLabel",
},
},
&ruleRefExpr{
	pos: position{line: 400, col: 45, offset: 10871},
	name: "_",
},
&litMatcher{
	pos: position{line: 400, col: 47, offset: 10873},
	val: ":",
	ignoreCase: false,
},
&ruleRefExpr{
	pos: position{line: 400, col: 51, offset: 10877},
	name: "_1",
},
&labeledExpr{
	pos: position{line: 400, col: 54, offset: 10880},
	label: "t",
	expr: &ruleRefExpr{
	pos: position{line: 400, col: 56, offset: 10882},
	name: "Expression",
},
},
&ruleRefExpr{
	pos: position{line: 400, col: 67, offset: 10893},
	name: "_",
},
&litMatcher{
	pos: position{line: 400, col: 69, offset: 10895},
	val: ")",
	ignoreCase: false,
},
&ruleRefExpr{
	pos: position{line: 400, col: 73, offset: 10899},
	name: "_",
},
&ruleRefExpr{
	pos: position{line: 400, col: 75, offset: 10901},
	name: "Arrow",
},
&ruleRefExpr{
	pos: position{line: 400, col: 81, offset: 10907},
	name: "_",
},
&labeledExpr{
	pos: position{line: 400, col: 83, offset: 10909},
	label: "body",
	expr: &ruleRefExpr{
	pos: position{line: 400, col: 88, offset: 10914},
	name: "Expression",
},
},
	},
},
},
&actionExpr{
	pos: position{line: 403, col: 7, offset: 11031},
	run: (*parser).callonExpression22,
	expr: &seqExpr{
	pos: position{line: 403, col: 7, offset: 11031},
	exprs: []interface{}{
&ruleRefExpr{
	pos: position{line: 403, col: 7, offset: 11031},
	name: "If",
},
&ruleRefExpr{
	pos: position{line: 403, col: 10, offset: 11034},
	name: "_1",
},
&labeledExpr{
	pos: position{line: 403, col: 13, offset: 11037},
	label: "cond",
	expr: &ruleRefExpr{
	pos: position{line: 403, col: 18, offset: 11042},
	name: "Expression",
},
},
&ruleRefExpr{
	pos: position{line: 403, col: 29, offset: 11053},
	name: "_",
},
&ruleRefExpr{
	pos: position{line: 403, col: 31, offset: 11055},
	name: "Then",
},
&ruleRefExpr{
	pos: position{line: 403, col: 36, offset: 11060},
	name: "_1",
},
&labeledExpr{
	pos: position{line: 403, col: 39, offset: 11063},
	label: "t",
	expr: &ruleRefExpr{
	pos: position{line: 403, col: 41, offset: 11065},
	name: "Expression",
},
},
&ruleRefExpr{
	pos: position{line: 403, col: 52, offset: 11076},
	name: "_",
},
&ruleRefExpr{
	pos: position{line: 403, col: 54, offset: 11078},
	name: "Else",
},
&ruleRefExpr{
	pos: position{line: 403, col: 59, offset: 11083},
	name: "_1",
},
&labeledExpr{
	pos: position{line: 403, col: 62, offset: 11086},
	label: "f",
	expr: &ruleRefExpr{
	pos: position{line: 403, col: 64, offset: 11088},
	name: "Expression",
},
},
	},
},
},
&actionExpr{
	pos: position{line: 406, col: 7, offset: 11174},
	run: (*parser).callonExpression38,
	expr: &seqExpr{
	pos: position{line: 406, col: 7, offset: 11174},
	exprs: []interface{}{
&labeledExpr{
	pos: position{line: 406, col: 7, offset: 11174},
	label: "bindings",
	expr: &oneOrMoreExpr{
	pos: position{line: 406, col: 16, offset: 11183},
	expr: &ruleRefExpr{
	pos: position{line: 406, col: 16, offset: 11183},
	name: "LetBinding",
},
},
},
&ruleRefExpr{
	pos: position{line: 406, col: 28, offset: 11195},
	name: "In",
},
&ruleRefExpr{
	pos: position{line: 406, col: 31, offset: 11198},
	name: "_1",
},
&labeledExpr{
	pos: position{line: 406, col: 34, offset: 11201},
	label: "b",
	expr: &ruleRefExpr{
	pos: position{line: 406, col: 36, offset: 11203},
	name: "Expression",
},
},
	},
},
},
&actionExpr{
	pos: position{line: 413, col: 7, offset: 11443},
	run: (*parser).callonExpression47,
	expr: &seqExpr{
	pos: position{line: 413, col: 7, offset: 11443},
	exprs: []interface{}{
&ruleRefExpr{
	pos: position{line: 413, col: 7, offset: 11443},
	name: "Forall",
},
&ruleRefExpr{
	pos: position{line: 413, col: 14, offset: 11450},
	name: "_",
},
&litMatcher{
	pos: position{line: 413, col: 16, offset: 11452},
	val: "(",
	ignoreCase: false,
},
&ruleRefExpr{
	pos: position{line: 413, col: 20, offset: 11456},
	name: "_",
},
&labeledExpr{
	pos: position{line: 413, col: 22, offset: 11458},
	label: "label",
	expr: &ruleRefExpr{
	pos: position{line: 413, col: 28, offset: 11464},
	name: "NonreservedLabel",
},
},
&ruleRefExpr{
	pos: position{line: 413, col: 45, offset: 11481},
	name: "_",
},
&litMatcher{
	pos: position{line: 413, col: 47, offset: 11483},
	val: ":",
	ignoreCase: false,
},
&ruleRefExpr{
	pos: position{line: 413, col: 51, offset: 11487},
	name: "_1",
},
&labeledExpr{
	pos: position{line: 413, col: 54, offset: 11490},
	label: "t",
	expr: &ruleRefExpr{
	pos: position{line: 413, col: 56, offset: 11492},
	name: "Expression",
},
},
&ruleRefExpr{
	pos: position{line: 413, col: 67, offset: 11503},
	name: "_",
},
&litMatcher{
	pos: position{line: 413, col: 69, offset: 11505},
	val: ")",
	ignoreCase: false,
},
&ruleRefExpr{
	pos: position{line: 413, col: 73, offset: 11509},
	name: "_",
},
&ruleRefExpr{
	pos: position{line: 413, col: 75, offset: 11511},
	name: "Arrow",
},
&ruleRefExpr{
	pos: position{line: 413, col: 81, offset: 11517},
	name: "_",
},
&labeledExpr{
	pos: position{line: 413, col: 83, offset: 11519},
	label: "body",
	expr: &ruleRefExpr{
	pos: position{line: 413, col: 88, offset: 11524},
	name: "Expression",
},
},
	},
},
},
&actionExpr{
	pos: position{line: 416, col: 7, offset: 11633},
	run: (*parser).callonExpression67,
	expr: &seqExpr{
	pos: position{line: 416, col: 7, offset: 11633},
	exprs: []interface{}{
&labeledExpr{
	pos: position{line: 416, col: 7, offset: 11633},
	label: "o",
	expr: &ruleRefExpr{
	pos: position{line: 416, col: 9, offset: 11635},
	name: "OperatorExpression",
},
},
&ruleRefExpr{
	pos: position{line: 416, col: 28, offset: 11654},
	name: "_",
},
&ruleRefExpr{
	pos: position{line: 416, col: 30, offset: 11656},
	name: "Arrow",
},
&ruleRefExpr{
	pos: position{line: 416, col: 36, offset: 11662},
	name: "_",
},
&labeledExpr{
	pos: position{line: 416, col: 38, offset: 11664},
	label: "e",
	expr: &ruleRefExpr{
	pos: position{line: 416, col: 40, offset: 11666},
	name: "Expression",
},
},
	},
},
},
&actionExpr{
	pos: position{line: 417, col: 7, offset: 11725},
	run: (*parser).callonExpression76,
	expr: &seqExpr{
	pos: position{line: 417, col: 7, offset: 11725},
	exprs: []interface{}{
&ruleRefExpr{
	pos: position{line: 417, col: 7, offset: 11725},
	name: "Merge",
},
&ruleRefExpr{
	pos: position{line: 417, col: 13, offset: 11731},
	name: "_1",
},
&labeledExpr{
	pos: position{line: 417, col: 16, offset: 11734},
	label: "h",
	expr: &ruleRefExpr{
	pos: position{line: 417, col: 18, offset: 11736},
	name: "ImportExpression",
},
},
&ruleRefExpr{
	pos: position{line: 417, col: 35, offset: 11753},
	name: "_1",
},
&labeledExpr{
	pos: position{line: 417, col: 38, offset: 11756},
	label: "u",
	expr: &ruleRefExpr{
	pos: position{line: 417, col: 40, offset: 11758},
	name: "ImportExpression",
},
},
&ruleRefExpr{
	pos: position{line: 417, col: 57, offset: 11775},
	name: "_",
},
&litMatcher{
	pos: position{line: 417, col: 59, offset: 11777},
	val: ":",
	ignoreCase: false,
},
&ruleRefExpr{
	pos: position{line: 417, col: 63, offset: 11781},
	name: "_1",
},
&labeledExpr{
	pos: position{line: 417, col: 66, offset: 11784},
	label: "a",
	expr: &ruleRefExpr{
	pos: position{line: 417, col: 68, offset: 11786},
	name: "ApplicationExpression",
},
},
	},
},
},
&ruleRefExpr{
	pos: position{line: 420, col: 7, offset: 11907},
	name: "EmptyList",
},
&ruleRefExpr{
	pos: position{line: 421, col: 7, offset: 11923},
	name: "AnnotatedExpression",
},
	},
},
},
{
	name: "Annotation",
	pos: position{line: 423, col: 1, offset: 11944},
	expr: &actionExpr{
	pos: position{line: 423, col: 14, offset: 11959},
	run: (*parser).callonAnnotation1,
	expr: &seqExpr{
	pos: position{line: 423, col: 14, offset: 11959},
	exprs: []interface{}{
&litMatcher{
	pos: position{line: 423, col: 14, offset: 11959},
	val: ":",
	ignoreCase: false,
},
&ruleRefExpr{
	pos: position{line: 423, col: 18, offset: 11963},
	name: "_1",
},
&labeledExpr{
	pos: position{line: 423, col: 21, offset: 11966},
	label: "a",
	expr: &ruleRefExpr{
	pos: position{line: 423, col: 23, offset: 11968},
	name: "Expression",
},
},
	},
},
},
},
{
	name: "AnnotatedExpression",
	pos: position{line: 425, col: 1, offset: 11998},
	expr: &actionExpr{
	pos: position{line: 426, col: 1, offset: 12022},
	run: (*parser).callonAnnotatedExpression1,
	expr: &seqExpr{
	pos: position{line: 426, col: 1, offset: 12022},
	exprs: []interface{}{
&labeledExpr{
	pos: position{line: 426, col: 1, offset: 12022},
	label: "e",
	expr: &ruleRefExpr{
	pos: position{line: 426, col: 3, offset: 12024},
	name: "OperatorExpression",
},
},
&labeledExpr{
	pos: position{line: 426, col: 22, offset: 12043},
	label: "a",
	expr: &zeroOrOneExpr{
	pos: position{line: 426, col: 24, offset: 12045},
	expr: &seqExpr{
	pos: position{line: 426, col: 25, offset: 12046},
	exprs: []interface{}{
&ruleRefExpr{
	pos: position{line: 426, col: 25, offset: 12046},
	name: "_",
},
&ruleRefExpr{
	pos: position{line: 426, col: 27, offset: 12048},
	name: "Annotation",
},
	},
},
},
},
	},
},
},
},
{
	name: "EmptyList",
	pos: position{line: 431, col: 1, offset: 12173},
	expr: &actionExpr{
	pos: position{line: 431, col: 13, offset: 12187},
	run: (*parser).callonEmptyList1,
	expr: &seqExpr{
	pos: position{line: 431, col: 13, offset: 12187},
	exprs: []interface{}{
&litMatcher{
	pos: position{line: 431, col: 13, offset: 12187},
	val: "[",
	ignoreCase: false,
},
&ruleRefExpr{
	pos: position{line: 431, col: 17, offset: 12191},
	name: "_",
},
&litMatcher{
	pos: position{line: 431, col: 19, offset: 12193},
	val: "]",
	ignoreCase: false,
},
&ruleRefExpr{
	pos: position{line: 431, col: 23, offset: 12197},
	name: "_",
},
&litMatcher{
	pos: position{line: 431, col: 25, offset: 12199},
	val: ":",
	ignoreCase: false,
},
&ruleRefExpr{
	pos: position{line: 431, col: 29, offset: 12203},
	name: "_1",
},
&ruleRefExpr{
	pos: position{line: 431, col: 32, offset: 12206},
	name: "List",
},
&ruleRefExpr{
	pos: position{line: 431, col: 37, offset: 12211},
	name: "_",
},
&labeledExpr{
	pos: position{line: 431, col: 39, offset: 12213},
	label: "a",
	expr: &ruleRefExpr{
	pos: position{line: 431, col: 41, offset: 12215},
	name: "ImportExpression",
},
},
	},
},
},
},
{
	name: "OperatorExpression",
	pos: position{line: 435, col: 1, offset: 12278},
	expr: &ruleRefExpr{
	pos: position{line: 435, col: 22, offset: 12301},
	name: "ImportAltExpression",
},
},
{
	name: "ImportAltExpression",
	pos: position{line: 437, col: 1, offset: 12322},
	expr: &ruleRefExpr{
	pos: position{line: 437, col: 24, offset: 12347},
	name: "OrExpression",
},
},
{
	name: "OrExpression",
	pos: position{line: 439, col: 1, offset: 12361},
	expr: &actionExpr{
	pos: position{line: 439, col: 26, offset: 12388},
	run: (*parser).callonOrExpression1,
	expr: &seqExpr{
	pos: position{line: 439, col: 26, offset: 12388},
	exprs: []interface{}{
&labeledExpr{
	pos: position{line: 439, col: 26, offset: 12388},
	label: "first",
	expr: &ruleRefExpr{
	pos: position{line: 439, col: 32, offset: 12394},
	name: "PlusExpression",
},
},
&labeledExpr{
	pos: position{line: 439, col: 55, offset: 12417},
	label: "rest",
	expr: &zeroOrMoreExpr{
	pos: position{line: 439, col: 60, offset: 12422},
	expr: &seqExpr{
	pos: position{line: 439, col: 61, offset: 12423},
	exprs: []interface{}{
&ruleRefExpr{
	pos: position{line: 439, col: 61, offset: 12423},
	name: "_",
},
&litMatcher{
	pos: position{line: 439, col: 63, offset: 12425},
	val: "||",
	ignoreCase: false,
},
&ruleRefExpr{
	pos: position{line: 439, col: 68, offset: 12430},
	name: "_",
},
&ruleRefExpr{
	pos: position{line: 439, col: 70, offset: 12432},
	name: "PlusExpression",
},
	},
},
},
},
	},
},
},
},
{
	name: "PlusExpression",
	pos: position{line: 441, col: 1, offset: 12498},
	expr: &actionExpr{
	pos: position{line: 441, col: 26, offset: 12525},
	run: (*parser).callonPlusExpression1,
	expr: &seqExpr{
	pos: position{line: 441, col: 26, offset: 12525},
	exprs: []interface{}{
&labeledExpr{
	pos: position{line: 441, col: 26, offset: 12525},
	label: "first",
	expr: &ruleRefExpr{
	pos: position{line: 441, col: 32, offset: 12531},
	name: "TextAppendExpression",
},
},
&labeledExpr{
	pos: position{line: 441, col: 55, offset: 12554},
	label: "rest",
	expr: &zeroOrMoreExpr{
	pos: position{line: 441, col: 60, offset: 12559},
	expr: &seqExpr{
	pos: position{line: 441, col: 61, offset: 12560},
	exprs: []interface{}{
&ruleRefExpr{
	pos: position{line: 441, col: 61, offset: 12560},
	name: "_",
},
&litMatcher{
	pos: position{line: 441, col: 63, offset: 12562},
	val: "+",
	ignoreCase: false,
},
&ruleRefExpr{
	pos: position{line: 441, col: 67, offset: 12566},
	name: "_1",
},
&labeledExpr{
	pos: position{line: 441, col: 70, offset: 12569},
	label: "e",
	expr: &ruleRefExpr{
	pos: position{line: 441, col: 72, offset: 12571},
	name: "TextAppendExpression",
},
},
	},
},
},
},
	},
},
},
},
{
	name: "TextAppendExpression",
	pos: position{line: 443, col: 1, offset: 12645},
	expr: &actionExpr{
	pos: position{line: 443, col: 26, offset: 12672},
	run: (*parser).callonTextAppendExpression1,
	expr: &seqExpr{
	pos: position{line: 443, col: 26, offset: 12672},
	exprs: []interface{}{
&labeledExpr{
	pos: position{line: 443, col: 26, offset: 12672},
	label: "first",
	expr: &ruleRefExpr{
	pos: position{line: 443, col: 32, offset: 12678},
	name: "ListAppendExpression",
},
},
&labeledExpr{
	pos: position{line: 443, col: 55, offset: 12701},
	label: "rest",
	expr: &zeroOrMoreExpr{
	pos: position{line: 443, col: 60, offset: 12706},
	expr: &seqExpr{
	pos: position{line: 443, col: 61, offset: 12707},
	exprs: []interface{}{
&ruleRefExpr{
	pos: position{line: 443, col: 61, offset: 12707},
	name: "_",
},
&litMatcher{
	pos: position{line: 443, col: 63, offset: 12709},
	val: "++",
	ignoreCase: false,
},
&ruleRefExpr{
	pos: position{line: 443, col: 68, offset: 12714},
	name: "_",
},
&labeledExpr{
	pos: position{line: 443, col: 70, offset: 12716},
	label: "e",
	expr: &ruleRefExpr{
	pos: position{line: 443, col: 72, offset: 12718},
	name: "ListAppendExpression",
},
},
	},
},
},
},
	},
},
},
},
{
	name: "ListAppendExpression",
	pos: position{line: 445, col: 1, offset: 12798},
	expr: &actionExpr{
	pos: position{line: 445, col: 26, offset: 12825},
	run: (*parser).callonListAppendExpression1,
	expr: &seqExpr{
	pos: position{line: 445, col: 26, offset: 12825},
	exprs: []interface{}{
&labeledExpr{
	pos: position{line: 445, col: 26, offset: 12825},
	label: "first",
	expr: &ruleRefExpr{
	pos: position{line: 445, col: 32, offset: 12831},
	name: "AndExpression",
},
},
&labeledExpr{
	pos: position{line: 445, col: 55, offset: 12854},
	label: "rest",
	expr: &zeroOrMoreExpr{
	pos: position{line: 445, col: 60, offset: 12859},
	expr: &seqExpr{
	pos: position{line: 445, col: 61, offset: 12860},
	exprs: []interface{}{
&ruleRefExpr{
	pos: position{line: 445, col: 61, offset: 12860},
	name: "_",
},
&litMatcher{
	pos: position{line: 445, col: 63, offset: 12862},
	val: "#",
	ignoreCase: false,
},
&ruleRefExpr{
	pos: position{line: 445, col: 67, offset: 12866},
	name: "_",
},
&labeledExpr{
	pos: position{line: 445, col: 69, offset: 12868},
	label: "e",
	expr: &ruleRefExpr{
	pos: position{line: 445, col: 71, offset: 12870},
	name: "AndExpression",
},
},
	},
},
},
},
	},
},
},
},
{
	name: "AndExpression",
	pos: position{line: 447, col: 1, offset: 12943},
	expr: &actionExpr{
	pos: position{line: 447, col: 26, offset: 12970},
	run: (*parser).callonAndExpression1,
	expr: &seqExpr{
	pos: position{line: 447, col: 26, offset: 12970},
	exprs: []interface{}{
&labeledExpr{
	pos: position{line: 447, col: 26, offset: 12970},
	label: "first",
	expr: &ruleRefExpr{
	pos: position{line: 447, col: 32, offset: 12976},
	name: "CombineExpression",
},
},
&labeledExpr{
	pos: position{line: 447, col: 55, offset: 12999},
	label: "rest",
	expr: &zeroOrMoreExpr{
	pos: position{line: 447, col: 60, offset: 13004},
	expr: &seqExpr{
	pos: position{line: 447, col: 61, offset: 13005},
	exprs: []interface{}{
&ruleRefExpr{
	pos: position{line: 447, col: 61, offset: 13005},
	name: "_",
},
&litMatcher{
	pos: position{line: 447, col: 63, offset: 13007},
	val: "&&",
	ignoreCase: false,
},
&ruleRefExpr{
	pos: position{line: 447, col: 68, offset: 13012},
	name: "_",
},
&labeledExpr{
	pos: position{line: 447, col: 70, offset: 13014},
	label: "e",
	expr: &ruleRefExpr{
	pos: position{line: 447, col: 72, offset: 13016},
	name: "CombineExpression",
},
},
	},
},
},
},
	},
},
},
},
{
	name: "CombineExpression",
	pos: position{line: 449, col: 1, offset: 13086},
	expr: &actionExpr{
	pos: position{line: 449, col: 26, offset: 13113},
	run: (*parser).callonCombineExpression1,
	expr: &seqExpr{
	pos: position{line: 449, col: 26, offset: 13113},
	exprs: []interface{}{
&labeledExpr{
	pos: position{line: 449, col: 26, offset: 13113},
	label: "first",
	expr: &ruleRefExpr{
	pos: position{line: 449, col: 32, offset: 13119},
	name: "PreferExpression",
},
},
&labeledExpr{
	pos: position{line: 449, col: 55, offset: 13142},
	label: "rest",
	expr: &zeroOrMoreExpr{
	pos: position{line: 449, col: 60, offset: 13147},
	expr: &seqExpr{
	pos: position{line: 449, col: 61, offset: 13148},
	exprs: []interface{}{
&ruleRefExpr{
	pos: position{line: 449, col: 61, offset: 13148},
	name: "_",
},
&ruleRefExpr{
	pos: position{line: 449, col: 63, offset: 13150},
	name: "Combine",
},
&ruleRefExpr{
	pos: position{line: 449, col: 71, offset: 13158},
	name: "_",
},
&labeledExpr{
	pos: position{line: 449, col: 73, offset: 13160},
	label: "e",
	expr: &ruleRefExpr{
	pos: position{line: 449, col: 75, offset: 13162},
	name: "PreferExpression",
},
},
	},
},
},
},
	},
},
},
},
{
	name: "PreferExpression",
	pos: position{line: 451, col: 1, offset: 13239},
	expr: &actionExpr{
	pos: position{line: 451, col: 26, offset: 13266},
	run: (*parser).callonPreferExpression1,
	expr: &seqExpr{
	pos: position{line: 451, col: 26, offset: 13266},
	exprs: []interface{}{
&labeledExpr{
	pos: position{line: 451, col: 26, offset: 13266},
	label: "first",
	expr: &ruleRefExpr{
	pos: position{line: 451, col: 32, offset: 13272},
	name: "CombineTypesExpression",
},
},
&labeledExpr{
	pos: position{line: 451, col: 55, offset: 13295},
	label: "rest",
	expr: &zeroOrMoreExpr{
	pos: position{line: 451, col: 60, offset: 13300},
	expr: &seqExpr{
	pos: position{line: 451, col: 61, offset: 13301},
	exprs: []interface{}{
&ruleRefExpr{
	pos: position{line: 451, col: 61, offset: 13301},
	name: "_",
},
&ruleRefExpr{
	pos: position{line: 451, col: 63, offset: 13303},
	name: "Prefer",
},
&ruleRefExpr{
	pos: position{line: 451, col: 70, offset: 13310},
	name: "_",
},
&labeledExpr{
	pos: position{line: 451, col: 72, offset: 13312},
	label: "e",
	expr: &ruleRefExpr{
	pos: position{line: 451, col: 74, offset: 13314},
	name: "CombineTypesExpression",
},
},
	},
},
},
},
	},
},
},
},
{
	name: "CombineTypesExpression",
	pos: position{line: 453, col: 1, offset: 13408},
	expr: &actionExpr{
	pos: position{line: 453, col: 26, offset: 13435},
	run: (*parser).callonCombineTypesExpression1,
	expr: &seqExpr{
	pos: position{line: 453, col: 26, offset: 13435},
	exprs: []interface{}{
&labeledExpr{
	pos: position{line: 453, col: 26, offset: 13435},
	label: "first",
	expr: &ruleRefExpr{
	pos: position{line: 453, col: 32, offset: 13441},
	name: "TimesExpression",
},
},
&labeledExpr{
	pos: position{line: 453, col: 55, offset: 13464},
	label: "rest",
	expr: &zeroOrMoreExpr{
	pos: position{line: 453, col: 60, offset: 13469},
	expr: &seqExpr{
	pos: position{line: 453, col: 61, offset: 13470},
	exprs: []interface{}{
&ruleRefExpr{
	pos: position{line: 453, col: 61, offset: 13470},
	name: "_",
},
&ruleRefExpr{
	pos: position{line: 453, col: 63, offset: 13472},
	name: "CombineTypes",
},
&ruleRefExpr{
	pos: position{line: 453, col: 76, offset: 13485},
	name: "_",
},
&labeledExpr{
	pos: position{line: 453, col: 78, offset: 13487},
	label: "e",
	expr: &ruleRefExpr{
	pos: position{line: 453, col: 80, offset: 13489},
	name: "TimesExpression",
},
},
	},
},
},
},
	},
},
},
},
{
	name: "TimesExpression",
	pos: position{line: 455, col: 1, offset: 13569},
	expr: &actionExpr{
	pos: position{line: 455, col: 26, offset: 13596},
	run: (*parser).callonTimesExpression1,
	expr: &seqExpr{
	pos: position{line: 455, col: 26, offset: 13596},
	exprs: []interface{}{
&labeledExpr{
	pos: position{line: 455, col: 26, offset: 13596},
	label: "first",
	expr: &ruleRefExpr{
	pos: position{line: 455, col: 32, offset: 13602},
	name: "EqualExpression",
},
},
&labeledExpr{
	pos: position{line: 455, col: 55, offset: 13625},
	label: "rest",
	expr: &zeroOrMoreExpr{
	pos: position{line: 455, col: 60, offset: 13630},
	expr: &seqExpr{
	pos: position{line: 455, col: 61, offset: 13631},
	exprs: []interface{}{
&ruleRefExpr{
	pos: position{line: 455, col: 61, offset: 13631},
	name: "_",
},
&litMatcher{
	pos: position{line: 455, col: 63, offset: 13633},
	val: "*",
	ignoreCase: false,
},
&ruleRefExpr{
	pos: position{line: 455, col: 67, offset: 13637},
	name: "_",
},
&labeledExpr{
	pos: position{line: 455, col: 69, offset: 13639},
	label: "e",
	expr: &ruleRefExpr{
	pos: position{line: 455, col: 71, offset: 13641},
	name: "EqualExpression",
},
},
	},
},
},
},
	},
},
},
},
{
	name: "EqualExpression",
	pos: position{line: 457, col: 1, offset: 13711},
	expr: &actionExpr{
	pos: position{line: 457, col: 26, offset: 13738},
	run: (*parser).callonEqualExpression1,
	expr: &seqExpr{
	pos: position{line: 457, col: 26, offset: 13738},
	exprs: []interface{}{
&labeledExpr{
	pos: position{line: 457, col: 26, offset: 13738},
	label: "first",
	expr: &ruleRefExpr{
	pos: position{line: 457, col: 32, offset: 13744},
	name: "NotEqualExpression",
},
},
&labeledExpr{
	pos: position{line: 457, col: 55, offset: 13767},
	label: "rest",
	expr: &zeroOrMoreExpr{
	pos: position{line: 457, col: 60, offset: 13772},
	expr: &seqExpr{
	pos: position{line: 457, col: 61, offset: 13773},
	exprs: []interface{}{
&ruleRefExpr{
	pos: position{line: 457, col: 61, offset: 13773},
	name: "_",
},
&litMatcher{
	pos: position{line: 457, col: 63, offset: 13775},
	val: "==",
	ignoreCase: false,
},
&ruleRefExpr{
	pos: position{line: 457, col: 68, offset: 13780},
	name: "_",
},
&labeledExpr{
	pos: position{line: 457, col: 70, offset: 13782},
	label: "e",
	expr: &ruleRefExpr{
	pos: position{line: 457, col: 72, offset: 13784},
	name: "NotEqualExpression",
},
},
	},
},
},
},
	},
},
},
},
{
	name: "NotEqualExpression",
	pos: position{line: 459, col: 1, offset: 13854},
	expr: &actionExpr{
	pos: position{line: 459, col: 26, offset: 13881},
	run: (*parser).callonNotEqualExpression1,
	expr: &seqExpr{
	pos: position{line: 459, col: 26, offset: 13881},
	exprs: []interface{}{
&labeledExpr{
	pos: position{line: 459, col: 26, offset: 13881},
	label: "first",
	expr: &ruleRefExpr{
	pos: position{line: 459, col: 32, offset: 13887},
	name: "ApplicationExpression",
},
},
&labeledExpr{
	pos: position{line: 459, col: 55, offset: 13910},
	label: "rest",
	expr: &zeroOrMoreExpr{
	pos: position{line: 459, col: 60, offset: 13915},
	expr: &seqExpr{
	pos: position{line: 459, col: 61, offset: 13916},
	exprs: []interface{}{
&ruleRefExpr{
	pos: position{line: 459, col: 61, offset: 13916},
	name: "_",
},
&litMatcher{
	pos: position{line: 459, col: 63, offset: 13918},
	val: "!=",
	ignoreCase: false,
},
&ruleRefExpr{
	pos: position{line: 459, col: 68, offset: 13923},
	name: "_",
},
&labeledExpr{
	pos: position{line: 459, col: 70, offset: 13925},
	label: "e",
	expr: &ruleRefExpr{
	pos: position{line: 459, col: 72, offset: 13927},
	name: "ApplicationExpression",
},
},
	},
},
},
},
	},
},
},
},
{
	name: "ApplicationExpression",
	pos: position{line: 462, col: 1, offset: 14001},
	expr: &actionExpr{
	pos: position{line: 462, col: 25, offset: 14027},
	run: (*parser).callonApplicationExpression1,
	expr: &seqExpr{
	pos: position{line: 462, col: 25, offset: 14027},
	exprs: []interface{}{
&labeledExpr{
	pos: position{line: 462, col: 25, offset: 14027},
	label: "f",
	expr: &ruleRefExpr{
	pos: position{line: 462, col: 27, offset: 14029},
	name: "FirstApplicationExpression",
},
},
&labeledExpr{
	pos: position{line: 462, col: 54, offset: 14056},
	label: "rest",
	expr: &zeroOrMoreExpr{
	pos: position{line: 462, col: 59, offset: 14061},
	expr: &seqExpr{
	pos: position{line: 462, col: 60, offset: 14062},
	exprs: []interface{}{
&ruleRefExpr{
	pos: position{line: 462, col: 60, offset: 14062},
	name: "_1",
},
&ruleRefExpr{
	pos: position{line: 462, col: 63, offset: 14065},
	name: "ImportExpression",
},
	},
},
},
},
	},
},
},
},
{
	name: "FirstApplicationExpression",
	pos: position{line: 471, col: 1, offset: 14308},
	expr: &choiceExpr{
	pos: position{line: 472, col: 8, offset: 14346},
	alternatives: []interface{}{
&actionExpr{
	pos: position{line: 472, col: 8, offset: 14346},
	run: (*parser).callonFirstApplicationExpression2,
	expr: &seqExpr{
	pos: position{line: 472, col: 8, offset: 14346},
	exprs: []interface{}{
&ruleRefExpr{
	pos: position{line: 472, col: 8, offset: 14346},
	name: "Merge",
},
&ruleRefExpr{
	pos: position{line: 472, col: 14, offset: 14352},
	name: "_1",
},
&labeledExpr{
	pos: position{line: 472, col: 17, offset: 14355},
	label: "h",
	expr: &ruleRefExpr{
	pos: position{line: 472, col: 19, offset: 14357},
	name: "ImportExpression",
},
},
&ruleRefExpr{
	pos: position{line: 472, col: 36, offset: 14374},
	name: "_1",
},
&labeledExpr{
	pos: position{line: 472, col: 39, offset: 14377},
	label: "u",
	expr: &ruleRefExpr{
	pos: position{line: 472, col: 41, offset: 14379},
	name: "ImportExpression",
},
},
	},
},
},
&actionExpr{
	pos: position{line: 475, col: 8, offset: 14482},
	run: (*parser).callonFirstApplicationExpression11,
	expr: &seqExpr{
	pos: position{line: 475, col: 8, offset: 14482},
	exprs: []interface{}{
&ruleRefExpr{
	pos: position{line: 475, col: 8, offset: 14482},
	name: "Some",
},
&ruleRefExpr{
	pos: position{line: 475, col: 13, offset: 14487},
	name: "_1",
},
&labeledExpr{
	pos: position{line: 475, col: 16, offset: 14490},
	label: "e",
	expr: &ruleRefExpr{
	pos: position{line: 475, col: 18, offset: 14492},
	name: "ImportExpression",
},
},
	},
},
},
&ruleRefExpr{
	pos: position{line: 476, col: 8, offset: 14547},
	name: "ImportExpression",
},
	},
},
},
{
	name: "ImportExpression",
	pos: position{line: 478, col: 1, offset: 14565},
	expr: &choiceExpr{
	pos: position{line: 478, col: 20, offset: 14586},
	alternatives: []interface{}{
&ruleRefExpr{
	pos: position{line: 478, col: 20, offset: 14586},
	name: "Import",
},
&ruleRefExpr{
	pos: position{line: 478, col: 29, offset: 14595},
	name: "SelectorExpression",
},
	},
},
},
{
	name: "SelectorExpression",
	pos: position{line: 480, col: 1, offset: 14615},
	expr: &actionExpr{
	pos: position{line: 480, col: 22, offset: 14638},
	run: (*parser).callonSelectorExpression1,
	expr: &seqExpr{
	pos: position{line: 480, col: 22, offset: 14638},
	exprs: []interface{}{
&labeledExpr{
	pos: position{line: 480, col: 22, offset: 14638},
	label: "e",
	expr: &ruleRefExpr{
	pos: position{line: 480, col: 24, offset: 14640},
	name: "PrimitiveExpression",
},
},
&labeledExpr{
	pos: position{line: 480, col: 44, offset: 14660},
	label: "ls",
	expr: &zeroOrMoreExpr{
	pos: position{line: 480, col: 47, offset: 14663},
	expr: &seqExpr{
	pos: position{line: 480, col: 48, offset: 14664},
	exprs: []interface{}{
&ruleRefExpr{
	pos: position{line: 480, col: 48, offset: 14664},
	name: "_",
},
&litMatcher{
	pos: position{line: 480, col: 50, offset: 14666},
	val: ".",
	ignoreCase: false,
},
&ruleRefExpr{
	pos: position{line: 480, col: 54, offset: 14670},
	name: "_",
},
&ruleRefExpr{
	pos: position{line: 480, col: 56, offset: 14672},
	name: "AnyLabel",
},
	},
},
},
},
	},
},
},
},
{
	name: "PrimitiveExpression",
	pos: position{line: 490, col: 1, offset: 14905},
	expr: &choiceExpr{
	pos: position{line: 491, col: 7, offset: 14935},
	alternatives: []interface{}{
&ruleRefExpr{
	pos: position{line: 491, col: 7, offset: 14935},
	name: "DoubleLiteral",
},
&ruleRefExpr{
	pos: position{line: 492, col: 7, offset: 14955},
	name: "NaturalLiteral",
},
&ruleRefExpr{
	pos: position{line: 493, col: 7, offset: 14976},
	name: "IntegerLiteral",
},
&ruleRefExpr{
	pos: position{line: 494, col: 7, offset: 14997},
	name: "TextLiteral",
},
&actionExpr{
	pos: position{line: 495, col: 7, offset: 15015},
	run: (*parser).callonPrimitiveExpression6,
	expr: &seqExpr{
	pos: position{line: 495, col: 7, offset: 15015},
	exprs: []interface{}{
&litMatcher{
	pos: position{line: 495, col: 7, offset: 15015},
	val: "{",
	ignoreCase: false,
},
&ruleRefExpr{
	pos: position{line: 495, col: 11, offset: 15019},
	name: "_",
},
&labeledExpr{
	pos: position{line: 495, col: 13, offset: 15021},
	label: "r",
	expr: &ruleRefExpr{
	pos: position{line: 495, col: 15, offset: 15023},
	name: "RecordTypeOrLiteral",
},
},
&ruleRefExpr{
	pos: position{line: 495, col: 35, offset: 15043},
	name: "_",
},
&litMatcher{
	pos: position{line: 495, col: 37, offset: 15045},
	val: "}",
	ignoreCase: false,
},
	},
},
},
&actionExpr{
	pos: position{line: 496, col: 7, offset: 15073},
	run: (*parser).callonPrimitiveExpression14,
	expr: &seqExpr{
	pos: position{line: 496, col: 7, offset: 15073},
	exprs: []interface{}{
&litMatcher{
	pos: position{line: 496, col: 7, offset: 15073},
	val: "<",
	ignoreCase: false,
},
&ruleRefExpr{
	pos: position{line: 496, col: 11, offset: 15077},
	name: "_",
},
&labeledExpr{
	pos: position{line: 496, col: 13, offset: 15079},
	label: "u",
	expr: &ruleRefExpr{
	pos: position{line: 496, col: 15, offset: 15081},
	name: "UnionType",
},
},
&ruleRefExpr{
	pos: position{line: 496, col: 25, offset: 15091},
	name: "_",
},
&litMatcher{
	pos: position{line: 496, col: 27, offset: 15093},
	val: ">",
	ignoreCase: false,
},
	},
},
},
&ruleRefExpr{
	pos: position{line: 497, col: 7, offset: 15121},
	name: "NonEmptyListLiteral",
},
&ruleRefExpr{
	pos: position{line: 498, col: 7, offset: 15147},
	name: "Identifier",
},
&actionExpr{
	pos: position{line: 499, col: 7, offset: 15164},
	run: (*parser).callonPrimitiveExpression24,
	expr: &seqExpr{
	pos: position{line: 499, col: 7, offset: 15164},
	exprs: []interface{}{
&litMatcher{
	pos: position{line: 499, col: 7, offset: 15164},
	val: "(",
	ignoreCase: false,
},
&ruleRefExpr{
	pos: position{line: 499, col: 11, offset: 15168},
	name: "_",
},
&labeledExpr{
	pos: position{line: 499, col: 14, offset: 15171},
	label: "e",
	expr: &ruleRefExpr{
	pos: position{line: 499, col: 16, offset: 15173},
	name: "Expression",
},
},
&ruleRefExpr{
	pos: position{line: 499, col: 27, offset: 15184},
	name: "_",
},
&litMatcher{
	pos: position{line: 499, col: 29, offset: 15186},
	val: ")",
	ignoreCase: false,
},
	},
},
},
	},
},
},
{
	name: "RecordTypeOrLiteral",
	pos: position{line: 501, col: 1, offset: 15209},
	expr: &choiceExpr{
	pos: position{line: 502, col: 7, offset: 15239},
	alternatives: []interface{}{
&actionExpr{
	pos: position{line: 502, col: 7, offset: 15239},
	run: (*parser).callonRecordTypeOrLiteral2,
	expr: &litMatcher{
	pos: position{line: 502, col: 7, offset: 15239},
	val: "=",
	ignoreCase: false,
},
},
&ruleRefExpr{
	pos: position{line: 503, col: 7, offset: 15294},
	name: "NonEmptyRecordType",
},
&ruleRefExpr{
	pos: position{line: 504, col: 7, offset: 15319},
	name: "NonEmptyRecordLiteral",
},
&actionExpr{
	pos: position{line: 505, col: 7, offset: 15347},
	run: (*parser).callonRecordTypeOrLiteral6,
	expr: &litMatcher{
	pos: position{line: 505, col: 7, offset: 15347},
	val: "",
	ignoreCase: false,
},
},
	},
},
},
{
	name: "RecordTypeField",
	pos: position{line: 507, col: 1, offset: 15393},
	expr: &actionExpr{
	pos: position{line: 507, col: 19, offset: 15413},
	run: (*parser).callonRecordTypeField1,
	expr: &seqExpr{
	pos: position{line: 507, col: 19, offset: 15413},
	exprs: []interface{}{
&labeledExpr{
	pos: position{line: 507, col: 19, offset: 15413},
	label: "name",
	expr: &ruleRefExpr{
	pos: position{line: 507, col: 24, offset: 15418},
	name: "AnyLabel",
},
},
&ruleRefExpr{
	pos: position{line: 507, col: 33, offset: 15427},
	name: "_",
},
&litMatcher{
	pos: position{line: 507, col: 35, offset: 15429},
	val: ":",
	ignoreCase: false,
},
&ruleRefExpr{
	pos: position{line: 507, col: 39, offset: 15433},
	name: "_1",
},
&labeledExpr{
	pos: position{line: 507, col: 42, offset: 15436},
	label: "expr",
	expr: &ruleRefExpr{
	pos: position{line: 507, col: 47, offset: 15441},
	name: "Expression",
},
},
	},
},
},
},
{
	name: "MoreRecordType",
	pos: position{line: 510, col: 1, offset: 15498},
	expr: &actionExpr{
	pos: position{line: 510, col: 18, offset: 15517},
	run: (*parser).callonMoreRecordType1,
	expr: &seqExpr{
	pos: position{line: 510, col: 18, offset: 15517},
	exprs: []interface{}{
&ruleRefExpr{
	pos: position{line: 510, col: 18, offset: 15517},
	name: "_",
},
&litMatcher{
	pos: position{line: 510, col: 20, offset: 15519},
	val: ",",
	ignoreCase: false,
},
&ruleRefExpr{
	pos: position{line: 510, col: 24, offset: 15523},
	name: "_",
},
&labeledExpr{
	pos: position{line: 510, col: 26, offset: 15525},
	label: "f",
	expr: &ruleRefExpr{
	pos: position{line: 510, col: 28, offset: 15527},
	name: "RecordTypeField",
},
},
	},
},
},
},
{
	name: "NonEmptyRecordType",
	pos: position{line: 511, col: 1, offset: 15559},
	expr: &actionExpr{
	pos: position{line: 512, col: 7, offset: 15588},
	run: (*parser).callonNonEmptyRecordType1,
	expr: &seqExpr{
	pos: position{line: 512, col: 7, offset: 15588},
	exprs: []interface{}{
&labeledExpr{
	pos: position{line: 512, col: 7, offset: 15588},
	label: "first",
	expr: &ruleRefExpr{
	pos: position{line: 512, col: 13, offset: 15594},
	name: "RecordTypeField",
},
},
&labeledExpr{
	pos: position{line: 512, col: 29, offset: 15610},
	label: "rest",
	expr: &zeroOrMoreExpr{
	pos: position{line: 512, col: 34, offset: 15615},
	expr: &ruleRefExpr{
	pos: position{line: 512, col: 34, offset: 15615},
	name: "MoreRecordType",
},
},
},
	},
},
},
},
{
	name: "RecordLiteralField",
	pos: position{line: 526, col: 1, offset: 16199},
	expr: &actionExpr{
	pos: position{line: 526, col: 22, offset: 16222},
	run: (*parser).callonRecordLiteralField1,
	expr: &seqExpr{
	pos: position{line: 526, col: 22, offset: 16222},
	exprs: []interface{}{
&labeledExpr{
	pos: position{line: 526, col: 22, offset: 16222},
	label: "name",
	expr: &ruleRefExpr{
	pos: position{line: 526, col: 27, offset: 16227},
	name: "AnyLabel",
},
},
&ruleRefExpr{
	pos: position{line: 526, col: 36, offset: 16236},
	name: "_",
},
&litMatcher{
	pos: position{line: 526, col: 38, offset: 16238},
	val: "=",
	ignoreCase: false,
},
&ruleRefExpr{
	pos: position{line: 526, col: 42, offset: 16242},
	name: "_",
},
&labeledExpr{
	pos: position{line: 526, col: 44, offset: 16244},
	label: "expr",
	expr: &ruleRefExpr{
	pos: position{line: 526, col: 49, offset: 16249},
	name: "Expression",
},
},
	},
},
},
},
{
	name: "MoreRecordLiteral",
	pos: position{line: 529, col: 1, offset: 16306},
	expr: &actionExpr{
	pos: position{line: 529, col: 21, offset: 16328},
	run: (*parser).callonMoreRecordLiteral1,
	expr: &seqExpr{
	pos: position{line: 529, col: 21, offset: 16328},
	exprs: []interface{}{
&ruleRefExpr{
	pos: position{line: 529, col: 21, offset: 16328},
	name: "_",
},
&litMatcher{
	pos: position{line: 529, col: 23, offset: 16330},
	val: ",",
	ignoreCase: false,
},
&ruleRefExpr{
	pos: position{line: 529, col: 27, offset: 16334},
	name: "_",
},
&labeledExpr{
	pos: position{line: 529, col: 29, offset: 16336},
	label: "f",
	expr: &ruleRefExpr{
	pos: position{line: 529, col: 31, offset: 16338},
	name: "RecordLiteralField",
},
},
	},
},
},
},
{
	name: "NonEmptyRecordLiteral",
	pos: position{line: 530, col: 1, offset: 16373},
	expr: &actionExpr{
	pos: position{line: 531, col: 7, offset: 16405},
	run: (*parser).callonNonEmptyRecordLiteral1,
	expr: &seqExpr{
	pos: position{line: 531, col: 7, offset: 16405},
	exprs: []interface{}{
&labeledExpr{
	pos: position{line: 531, col: 7, offset: 16405},
	label: "first",
	expr: &ruleRefExpr{
	pos: position{line: 531, col: 13, offset: 16411},
	name: "RecordLiteralField",
},
},
&labeledExpr{
	pos: position{line: 531, col: 32, offset: 16430},
	label: "rest",
	expr: &zeroOrMoreExpr{
	pos: position{line: 531, col: 37, offset: 16435},
	expr: &ruleRefExpr{
	pos: position{line: 531, col: 37, offset: 16435},
	name: "MoreRecordLiteral",
},
},
},
	},
},
},
},
{
	name: "UnionType",
	pos: position{line: 545, col: 1, offset: 17025},
	expr: &choiceExpr{
	pos: position{line: 545, col: 13, offset: 17039},
	alternatives: []interface{}{
&ruleRefExpr{
	pos: position{line: 545, col: 13, offset: 17039},
	name: "NonEmptyUnionType",
},
&ruleRefExpr{
	pos: position{line: 545, col: 33, offset: 17059},
	name: "EmptyUnionType",
},
	},
},
},
{
	name: "EmptyUnionType",
	pos: position{line: 547, col: 1, offset: 17075},
	expr: &actionExpr{
	pos: position{line: 547, col: 18, offset: 17094},
	run: (*parser).callonEmptyUnionType1,
	expr: &litMatcher{
	pos: position{line: 547, col: 18, offset: 17094},
	val: "",
	ignoreCase: false,
},
},
},
{
	name: "NonEmptyUnionType",
	pos: position{line: 549, col: 1, offset: 17126},
	expr: &actionExpr{
	pos: position{line: 549, col: 21, offset: 17148},
	run: (*parser).callonNonEmptyUnionType1,
	expr: &seqExpr{
	pos: position{line: 549, col: 21, offset: 17148},
	exprs: []interface{}{
&labeledExpr{
	pos: position{line: 549, col: 21, offset: 17148},
	label: "first",
	expr: &ruleRefExpr{
	pos: position{line: 549, col: 27, offset: 17154},
	name: "UnionVariant",
},
},
&labeledExpr{
	pos: position{line: 549, col: 40, offset: 17167},
	label: "rest",
	expr: &zeroOrMoreExpr{
	pos: position{line: 549, col: 45, offset: 17172},
	expr: &seqExpr{
	pos: position{line: 549, col: 46, offset: 17173},
	exprs: []interface{}{
&ruleRefExpr{
	pos: position{line: 549, col: 46, offset: 17173},
	name: "_",
},
&litMatcher{
	pos: position{line: 549, col: 48, offset: 17175},
	val: "|",
	ignoreCase: false,
},
&ruleRefExpr{
	pos: position{line: 549, col: 52, offset: 17179},
	name: "_",
},
&ruleRefExpr{
	pos: position{line: 549, col: 54, offset: 17181},
	name: "UnionVariant",
},
	},
},
},
},
	},
},
},
},
{
	name: "UnionVariant",
	pos: position{line: 569, col: 1, offset: 17903},
	expr: &seqExpr{
	pos: position{line: 569, col: 16, offset: 17920},
	exprs: []interface{}{
&ruleRefExpr{
	pos: position{line: 569, col: 16, offset: 17920},
	name: "AnyLabel",
},
&zeroOrOneExpr{
	pos: position{line: 569, col: 25, offset: 17929},
	expr: &seqExpr{
	pos: position{line: 569, col: 26, offset: 17930},
	exprs: []interface{}{
&ruleRefExpr{
	pos: position{line: 569, col: 26, offset: 17930},
	name: "_",
},
&litMatcher{
	pos: position{line: 569, col: 28, offset: 17932},
	val: ":",
	ignoreCase: false,
},
&ruleRefExpr{
	pos: position{line: 569, col: 32, offset: 17936},
	name: "_1",
},
&ruleRefExpr{
	pos: position{line: 569, col: 35, offset: 17939},
	name: "Expression",
},
	},
},
},
	},
},
},
{
	name: "MoreList",
	pos: position{line: 571, col: 1, offset: 17953},
	expr: &actionExpr{
	pos: position{line: 571, col: 12, offset: 17966},
	run: (*parser).callonMoreList1,
	expr: &seqExpr{
	pos: position{line: 571, col: 12, offset: 17966},
	exprs: []interface{}{
&litMatcher{
	pos: position{line: 571, col: 12, offset: 17966},
	val: ",",
	ignoreCase: false,
},
&ruleRefExpr{
	pos: position{line: 571, col: 16, offset: 17970},
	name: "_",
},
&labeledExpr{
	pos: position{line: 571, col: 18, offset: 17972},
	label: "e",
	expr: &ruleRefExpr{
	pos: position{line: 571, col: 20, offset: 17974},
	name: "Expression",
},
},
&ruleRefExpr{
	pos: position{line: 571, col: 31, offset: 17985},
	name: "_",
},
	},
},
},
},
{
	name: "NonEmptyListLiteral",
	pos: position{line: 573, col: 1, offset: 18004},
	expr: &actionExpr{
	pos: position{line: 574, col: 7, offset: 18034},
	run: (*parser).callonNonEmptyListLiteral1,
	expr: &seqExpr{
	pos: position{line: 574, col: 7, offset: 18034},
	exprs: []interface{}{
&litMatcher{
	pos: position{line: 574, col: 7, offset: 18034},
	val: "[",
	ignoreCase: false,
},
&ruleRefExpr{
	pos: position{line: 574, col: 11, offset: 18038},
	name: "_",
},
&labeledExpr{
	pos: position{line: 574, col: 13, offset: 18040},
	label: "first",
	expr: &ruleRefExpr{
	pos: position{line: 574, col: 19, offset: 18046},
	name: "Expression",
},
},
&ruleRefExpr{
	pos: position{line: 574, col: 30, offset: 18057},
	name: "_",
},
&labeledExpr{
	pos: position{line: 574, col: 32, offset: 18059},
	label: "rest",
	expr: &zeroOrMoreExpr{
	pos: position{line: 574, col: 37, offset: 18064},
	expr: &ruleRefExpr{
	pos: position{line: 574, col: 37, offset: 18064},
	name: "MoreList",
},
},
},
&litMatcher{
	pos: position{line: 574, col: 47, offset: 18074},
	val: "]",
	ignoreCase: false,
},
	},
},
},
},
{
	name: "EOF",
	pos: position{line: 584, col: 1, offset: 18350},
	expr: &notExpr{
	pos: position{line: 584, col: 7, offset: 18358},
	expr: &anyMatcher{
	line: 584, col: 8, offset: 18359,
},
},
},
	},
}
func (c *current) onDhallFile1(e interface{}) (interface{}, error) {
 return e, nil 
}

func (p *parser) callonDhallFile1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onDhallFile1(stack["e"])
}

func (c *current) onCompleteExpression1(e interface{}) (interface{}, error) {
 return e, nil 
}

func (p *parser) callonCompleteExpression1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onCompleteExpression1(stack["e"])
}

func (c *current) onEOL3() (interface{}, error) {
 return []byte{'\n'}, nil 
}

func (p *parser) callonEOL3() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onEOL3()
}

func (c *current) onLineComment5() (interface{}, error) {
 return string(c.text), nil
}

func (p *parser) callonLineComment5() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onLineComment5()
}

func (c *current) onLineComment1(content interface{}) (interface{}, error) {
 return content, nil 
}

func (p *parser) callonLineComment1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onLineComment1(stack["content"])
}

func (c *current) onSimpleLabel2() (interface{}, error) {
 return string(c.text), nil 
}

func (p *parser) callonSimpleLabel2() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onSimpleLabel2()
}

func (c *current) onSimpleLabel7() (interface{}, error) {
            return string(c.text), nil
          
}

func (p *parser) callonSimpleLabel7() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onSimpleLabel7()
}

func (c *current) onLabel1(label interface{}) (interface{}, error) {
 return label, nil 
}

func (p *parser) callonLabel1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onLabel1(stack["label"])
}

func (c *current) onNonreservedLabel2(label interface{}) (interface{}, error) {
 return label, nil 
}

func (p *parser) callonNonreservedLabel2() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onNonreservedLabel2(stack["label"])
}

func (c *current) onNonreservedLabel10(label interface{}) (interface{}, error) {
 return label, nil 
}

func (p *parser) callonNonreservedLabel10() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onNonreservedLabel10(stack["label"])
}

func (c *current) onDoubleQuoteChunk3(e interface{}) (interface{}, error) {
 return e, nil 
}

func (p *parser) callonDoubleQuoteChunk3() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onDoubleQuoteChunk3(stack["e"])
}

func (c *current) onDoubleQuoteEscaped6() (interface{}, error) {
 return []byte("\b"), nil 
}

func (p *parser) callonDoubleQuoteEscaped6() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onDoubleQuoteEscaped6()
}

func (c *current) onDoubleQuoteEscaped8() (interface{}, error) {
 return []byte("\f"), nil 
}

func (p *parser) callonDoubleQuoteEscaped8() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onDoubleQuoteEscaped8()
}

func (c *current) onDoubleQuoteEscaped10() (interface{}, error) {
 return []byte("\n"), nil 
}

func (p *parser) callonDoubleQuoteEscaped10() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onDoubleQuoteEscaped10()
}

func (c *current) onDoubleQuoteEscaped12() (interface{}, error) {
 return []byte("\r"), nil 
}

func (p *parser) callonDoubleQuoteEscaped12() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onDoubleQuoteEscaped12()
}

func (c *current) onDoubleQuoteEscaped14() (interface{}, error) {
 return []byte("\t"), nil 
}

func (p *parser) callonDoubleQuoteEscaped14() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onDoubleQuoteEscaped14()
}

func (c *current) onDoubleQuoteEscaped16() (interface{}, error) {
        i, err := strconv.ParseInt(string(c.text[1:]), 16, 32)
        return []byte(string([]rune{rune(i)})), err
     
}

func (p *parser) callonDoubleQuoteEscaped16() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onDoubleQuoteEscaped16()
}

func (c *current) onDoubleQuoteLiteral1(chunks interface{}) (interface{}, error) {
    var str strings.Builder
    var outChunks Chunks
    for _, chunk := range chunks.([]interface{}) {
        switch e := chunk.(type) {
        case []byte:
                str.Write(e)
        case Expr:
                outChunks = append(outChunks, Chunk{str.String(), e})
                str.Reset()
        default:
                return nil, errors.New("can't happen")
        }
    }
    return TextLit{Chunks: outChunks, Suffix: str.String()}, nil
}

func (p *parser) callonDoubleQuoteLiteral1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onDoubleQuoteLiteral1(stack["chunks"])
}

func (c *current) onEscapedQuotePair1() (interface{}, error) {
 return []byte("''"), nil 
}

func (p *parser) callonEscapedQuotePair1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onEscapedQuotePair1()
}

func (c *current) onEscapedInterpolation1() (interface{}, error) {
 return []byte("$\u007b"), nil 
}

func (p *parser) callonEscapedInterpolation1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onEscapedInterpolation1()
}

func (c *current) onSingleQuoteLiteral1(content interface{}) (interface{}, error) {
    var str strings.Builder
    var outChunks Chunks
    chunk, ok := content.([]interface{})
    for ; ok; chunk, ok = chunk[1].([]interface{}) {
        switch e := chunk[0].(type) {
        case []byte:
            str.Write(e)
        case Expr:
                outChunks = append(outChunks, Chunk{str.String(), e})
                str.Reset()
        default:
            return nil, errors.New("unimplemented")
        }
    }
    return RemoveLeadingCommonIndent(TextLit{Chunks: outChunks, Suffix: str.String()}), nil
}

func (p *parser) callonSingleQuoteLiteral1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onSingleQuoteLiteral1(stack["content"])
}

func (c *current) onInterpolation1(e interface{}) (interface{}, error) {
 return e, nil 
}

func (p *parser) callonInterpolation1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onInterpolation1(stack["e"])
}

func (c *current) onReserved2() (interface{}, error) {
 return NaturalBuild, nil 
}

func (p *parser) callonReserved2() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onReserved2()
}

func (c *current) onReserved4() (interface{}, error) {
 return NaturalFold, nil 
}

func (p *parser) callonReserved4() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onReserved4()
}

func (c *current) onReserved6() (interface{}, error) {
 return NaturalIsZero, nil 
}

func (p *parser) callonReserved6() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onReserved6()
}

func (c *current) onReserved8() (interface{}, error) {
 return NaturalEven, nil 
}

func (p *parser) callonReserved8() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onReserved8()
}

func (c *current) onReserved10() (interface{}, error) {
 return NaturalOdd, nil 
}

func (p *parser) callonReserved10() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onReserved10()
}

func (c *current) onReserved12() (interface{}, error) {
 return NaturalToInteger, nil 
}

func (p *parser) callonReserved12() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onReserved12()
}

func (c *current) onReserved14() (interface{}, error) {
 return NaturalShow, nil 
}

func (p *parser) callonReserved14() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onReserved14()
}

func (c *current) onReserved16() (interface{}, error) {
 return IntegerToDouble, nil 
}

func (p *parser) callonReserved16() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onReserved16()
}

func (c *current) onReserved18() (interface{}, error) {
 return IntegerShow, nil 
}

func (p *parser) callonReserved18() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onReserved18()
}

func (c *current) onReserved20() (interface{}, error) {
 return DoubleShow, nil 
}

func (p *parser) callonReserved20() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onReserved20()
}

func (c *current) onReserved22() (interface{}, error) {
 return ListBuild, nil 
}

func (p *parser) callonReserved22() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onReserved22()
}

func (c *current) onReserved24() (interface{}, error) {
 return ListFold, nil 
}

func (p *parser) callonReserved24() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onReserved24()
}

func (c *current) onReserved26() (interface{}, error) {
 return ListLength, nil 
}

func (p *parser) callonReserved26() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onReserved26()
}

func (c *current) onReserved28() (interface{}, error) {
 return ListHead, nil 
}

func (p *parser) callonReserved28() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onReserved28()
}

func (c *current) onReserved30() (interface{}, error) {
 return ListLast, nil 
}

func (p *parser) callonReserved30() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onReserved30()
}

func (c *current) onReserved32() (interface{}, error) {
 return ListIndexed, nil 
}

func (p *parser) callonReserved32() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onReserved32()
}

func (c *current) onReserved34() (interface{}, error) {
 return ListReverse, nil 
}

func (p *parser) callonReserved34() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onReserved34()
}

func (c *current) onReserved36() (interface{}, error) {
 return OptionalBuild, nil 
}

func (p *parser) callonReserved36() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onReserved36()
}

func (c *current) onReserved38() (interface{}, error) {
 return OptionalFold, nil 
}

func (p *parser) callonReserved38() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onReserved38()
}

func (c *current) onReserved40() (interface{}, error) {
 return TextShow, nil 
}

func (p *parser) callonReserved40() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onReserved40()
}

func (c *current) onReserved42() (interface{}, error) {
 return Bool, nil 
}

func (p *parser) callonReserved42() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onReserved42()
}

func (c *current) onReserved44() (interface{}, error) {
 return True, nil 
}

func (p *parser) callonReserved44() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onReserved44()
}

func (c *current) onReserved46() (interface{}, error) {
 return False, nil 
}

func (p *parser) callonReserved46() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onReserved46()
}

func (c *current) onReserved48() (interface{}, error) {
 return Optional, nil 
}

func (p *parser) callonReserved48() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onReserved48()
}

func (c *current) onReserved50() (interface{}, error) {
 return Natural, nil 
}

func (p *parser) callonReserved50() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onReserved50()
}

func (c *current) onReserved52() (interface{}, error) {
 return Integer, nil 
}

func (p *parser) callonReserved52() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onReserved52()
}

func (c *current) onReserved54() (interface{}, error) {
 return Double, nil 
}

func (p *parser) callonReserved54() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onReserved54()
}

func (c *current) onReserved56() (interface{}, error) {
 return Text, nil 
}

func (p *parser) callonReserved56() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onReserved56()
}

func (c *current) onReserved58() (interface{}, error) {
 return List, nil 
}

func (p *parser) callonReserved58() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onReserved58()
}

func (c *current) onReserved60() (interface{}, error) {
 return None, nil 
}

func (p *parser) callonReserved60() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onReserved60()
}

func (c *current) onReserved62() (interface{}, error) {
 return Type, nil 
}

func (p *parser) callonReserved62() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onReserved62()
}

func (c *current) onReserved64() (interface{}, error) {
 return Kind, nil 
}

func (p *parser) callonReserved64() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onReserved64()
}

func (c *current) onReserved66() (interface{}, error) {
 return Sort, nil 
}

func (p *parser) callonReserved66() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onReserved66()
}

func (c *current) onMissing1() (interface{}, error) {
 return Missing{}, nil 
}

func (p *parser) callonMissing1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onMissing1()
}

func (c *current) onNumericDoubleLiteral1() (interface{}, error) {
      d, err := strconv.ParseFloat(string(c.text), 64)
      if err != nil {
         return nil, err
      }
      return DoubleLit(d), nil
}

func (p *parser) callonNumericDoubleLiteral1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onNumericDoubleLiteral1()
}

func (c *current) onDoubleLiteral4() (interface{}, error) {
 return DoubleLit(math.Inf(1)), nil 
}

func (p *parser) callonDoubleLiteral4() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onDoubleLiteral4()
}

func (c *current) onDoubleLiteral6() (interface{}, error) {
 return DoubleLit(math.Inf(-1)), nil 
}

func (p *parser) callonDoubleLiteral6() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onDoubleLiteral6()
}

func (c *current) onDoubleLiteral10() (interface{}, error) {
 return DoubleLit(math.NaN()), nil 
}

func (p *parser) callonDoubleLiteral10() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onDoubleLiteral10()
}

func (c *current) onNaturalLiteral1() (interface{}, error) {
      i, err := strconv.Atoi(string(c.text))
      return NaturalLit(i), err
}

func (p *parser) callonNaturalLiteral1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onNaturalLiteral1()
}

func (c *current) onIntegerLiteral1() (interface{}, error) {
      i, err := strconv.Atoi(string(c.text))
      if err != nil {
         return nil, err
      }
      return IntegerLit(i), nil
}

func (p *parser) callonIntegerLiteral1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onIntegerLiteral1()
}

func (c *current) onDeBruijn1(index interface{}) (interface{}, error) {
 return int(index.(NaturalLit)), nil 
}

func (p *parser) callonDeBruijn1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onDeBruijn1(stack["index"])
}

func (c *current) onVariable1(name, index interface{}) (interface{}, error) {
    if index != nil {
        return Var{Name:name.(string), Index:index.(int)}, nil
    } else {
        return Var{Name:name.(string)}, nil
    }
}

func (p *parser) callonVariable1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onVariable1(stack["name"], stack["index"])
}

func (c *current) onUnquotedPathComponent1() (interface{}, error) {
 return string(c.text), nil 
}

func (p *parser) callonUnquotedPathComponent1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onUnquotedPathComponent1()
}

func (c *current) onPathComponent1(u interface{}) (interface{}, error) {
 return u, nil 
}

func (p *parser) callonPathComponent1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onPathComponent1(stack["u"])
}

func (c *current) onPath1(cs interface{}) (interface{}, error) {
    // urgh, have to convert []interface{} to []string
    components := make([]string, len(cs.([]interface{})))
    for i, component := range cs.([]interface{}) {
        components[i] = component.(string)
    }
    return path.Join(components...), nil
}

func (p *parser) callonPath1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onPath1(stack["cs"])
}

func (c *current) onParentPath1(p interface{}) (interface{}, error) {
 return Local(path.Join("..", p.(string))), nil 
}

func (p *parser) callonParentPath1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onParentPath1(stack["p"])
}

func (c *current) onHerePath1(p interface{}) (interface{}, error) {
 return Local(p.(string)), nil 
}

func (p *parser) callonHerePath1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onHerePath1(stack["p"])
}

func (c *current) onHomePath1(p interface{}) (interface{}, error) {
 return Local(path.Join("~", p.(string))), nil 
}

func (p *parser) callonHomePath1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onHomePath1(stack["p"])
}

func (c *current) onAbsolutePath1(p interface{}) (interface{}, error) {
 return Local(path.Join("/", p.(string))), nil 
}

func (p *parser) callonAbsolutePath1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onAbsolutePath1(stack["p"])
}

func (c *current) onHttpRaw1() (interface{}, error) {
 return url.ParseRequestURI(string(c.text)) 
}

func (p *parser) callonHttpRaw1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onHttpRaw1()
}

func (c *current) onIPv6address1() (interface{}, error) {
    addr := net.ParseIP(string(c.text))
    if addr == nil { return nil, errors.New("Malformed IPv6 address") }
    return string(c.text), nil
}

func (p *parser) callonIPv6address1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onIPv6address1()
}

func (c *current) onHttp1(u interface{}) (interface{}, error) {
 return MakeRemote(u.(*url.URL)) 
}

func (p *parser) callonHttp1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onHttp1(stack["u"])
}

func (c *current) onEnv1(v interface{}) (interface{}, error) {
 return v, nil 
}

func (p *parser) callonEnv1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onEnv1(stack["v"])
}

func (c *current) onBashEnvironmentVariable1() (interface{}, error) {
  return EnvVar(string(c.text)), nil
}

func (p *parser) callonBashEnvironmentVariable1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onBashEnvironmentVariable1()
}

func (c *current) onPosixEnvironmentVariable1(v interface{}) (interface{}, error) {
  return v, nil
}

func (p *parser) callonPosixEnvironmentVariable1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onPosixEnvironmentVariable1(stack["v"])
}

func (c *current) onPosixEnvironmentVariableContent1(v interface{}) (interface{}, error) {
  var b strings.Builder
  for _, c := range v.([]interface{}) {
    _, err := b.Write(c.([]byte))
    if err != nil { return nil, err }
  }
  return EnvVar(b.String()), nil
}

func (p *parser) callonPosixEnvironmentVariableContent1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onPosixEnvironmentVariableContent1(stack["v"])
}

func (c *current) onPosixEnvironmentVariableCharacter2() (interface{}, error) {
 return []byte{0x22}, nil 
}

func (p *parser) callonPosixEnvironmentVariableCharacter2() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onPosixEnvironmentVariableCharacter2()
}

func (c *current) onPosixEnvironmentVariableCharacter4() (interface{}, error) {
 return []byte{0x5c}, nil 
}

func (p *parser) callonPosixEnvironmentVariableCharacter4() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onPosixEnvironmentVariableCharacter4()
}

func (c *current) onPosixEnvironmentVariableCharacter6() (interface{}, error) {
 return []byte{0x07}, nil 
}

func (p *parser) callonPosixEnvironmentVariableCharacter6() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onPosixEnvironmentVariableCharacter6()
}

func (c *current) onPosixEnvironmentVariableCharacter8() (interface{}, error) {
 return []byte{0x08}, nil 
}

func (p *parser) callonPosixEnvironmentVariableCharacter8() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onPosixEnvironmentVariableCharacter8()
}

func (c *current) onPosixEnvironmentVariableCharacter10() (interface{}, error) {
 return []byte{0x0c}, nil 
}

func (p *parser) callonPosixEnvironmentVariableCharacter10() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onPosixEnvironmentVariableCharacter10()
}

func (c *current) onPosixEnvironmentVariableCharacter12() (interface{}, error) {
 return []byte{0x0a}, nil 
}

func (p *parser) callonPosixEnvironmentVariableCharacter12() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onPosixEnvironmentVariableCharacter12()
}

func (c *current) onPosixEnvironmentVariableCharacter14() (interface{}, error) {
 return []byte{0x0d}, nil 
}

func (p *parser) callonPosixEnvironmentVariableCharacter14() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onPosixEnvironmentVariableCharacter14()
}

func (c *current) onPosixEnvironmentVariableCharacter16() (interface{}, error) {
 return []byte{0x09}, nil 
}

func (p *parser) callonPosixEnvironmentVariableCharacter16() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onPosixEnvironmentVariableCharacter16()
}

func (c *current) onPosixEnvironmentVariableCharacter18() (interface{}, error) {
 return []byte{0x0b}, nil 
}

func (p *parser) callonPosixEnvironmentVariableCharacter18() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onPosixEnvironmentVariableCharacter18()
}

func (c *current) onImportHashed1(i interface{}) (interface{}, error) {
 return ImportHashed{Fetchable: i.(Fetchable)}, nil 
}

func (p *parser) callonImportHashed1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onImportHashed1(stack["i"])
}

func (c *current) onImport2(i interface{}) (interface{}, error) {
 return Embed(Import{ImportHashed: i.(ImportHashed), ImportMode: RawText}), nil 
}

func (p *parser) callonImport2() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onImport2(stack["i"])
}

func (c *current) onImport10(i interface{}) (interface{}, error) {
 return Embed(Import{ImportHashed: i.(ImportHashed), ImportMode: Code}), nil 
}

func (p *parser) callonImport10() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onImport10(stack["i"])
}

func (c *current) onLetBinding1(label, a, v interface{}) (interface{}, error) {
    if a != nil {
        return Binding{
            Variable: label.(string),
            Annotation: a.([]interface{})[0].(Expr),
            Value: v.(Expr),
        }, nil
    } else {
        return Binding{
            Variable: label.(string),
            Value: v.(Expr),
        }, nil
    }
}

func (p *parser) callonLetBinding1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onLetBinding1(stack["label"], stack["a"], stack["v"])
}

func (c *current) onExpression2(label, t, body interface{}) (interface{}, error) {
          return &LambdaExpr{Label:label.(string), Type:t.(Expr), Body: body.(Expr)}, nil
      
}

func (p *parser) callonExpression2() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onExpression2(stack["label"], stack["t"], stack["body"])
}

func (c *current) onExpression22(cond, t, f interface{}) (interface{}, error) {
          return BoolIf{cond.(Expr),t.(Expr),f.(Expr)},nil
      
}

func (p *parser) callonExpression22() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onExpression22(stack["cond"], stack["t"], stack["f"])
}

func (c *current) onExpression38(bindings, b interface{}) (interface{}, error) {
        bs := make([]Binding, len(bindings.([]interface{})))
        for i, binding := range bindings.([]interface{}) {
            bs[i] = binding.(Binding)
        }
        return MakeLet(b.(Expr), bs...), nil
      
}

func (p *parser) callonExpression38() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onExpression38(stack["bindings"], stack["b"])
}

func (c *current) onExpression47(label, t, body interface{}) (interface{}, error) {
          return &Pi{Label:label.(string), Type:t.(Expr), Body: body.(Expr)}, nil
      
}

func (p *parser) callonExpression47() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onExpression47(stack["label"], stack["t"], stack["body"])
}

func (c *current) onExpression67(o, e interface{}) (interface{}, error) {
 return FnType(o.(Expr),e.(Expr)), nil 
}

func (p *parser) callonExpression67() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onExpression67(stack["o"], stack["e"])
}

func (c *current) onExpression76(h, u, a interface{}) (interface{}, error) {
          return Merge{Handler:h.(Expr), Union:u.(Expr), Annotation:a.(Expr)}, nil
      
}

func (p *parser) callonExpression76() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onExpression76(stack["h"], stack["u"], stack["a"])
}

func (c *current) onAnnotation1(a interface{}) (interface{}, error) {
 return a, nil 
}

func (p *parser) callonAnnotation1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onAnnotation1(stack["a"])
}

func (c *current) onAnnotatedExpression1(e, a interface{}) (interface{}, error) {
        if a == nil { return e, nil }
        return Annot{e.(Expr), a.([]interface{})[1].(Expr)}, nil
    
}

func (p *parser) callonAnnotatedExpression1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onAnnotatedExpression1(stack["e"], stack["a"])
}

func (c *current) onEmptyList1(a interface{}) (interface{}, error) {
          return EmptyList{a.(Expr)},nil
}

func (p *parser) callonEmptyList1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onEmptyList1(stack["a"])
}

func (c *current) onOrExpression1(first, rest interface{}) (interface{}, error) {
return ParseOperator(OrOp, first, rest), nil
}

func (p *parser) callonOrExpression1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onOrExpression1(stack["first"], stack["rest"])
}

func (c *current) onPlusExpression1(first, rest interface{}) (interface{}, error) {
return ParseOperator(PlusOp, first, rest), nil
}

func (p *parser) callonPlusExpression1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onPlusExpression1(stack["first"], stack["rest"])
}

func (c *current) onTextAppendExpression1(first, rest interface{}) (interface{}, error) {
return ParseOperator(TextAppendOp, first, rest), nil
}

func (p *parser) callonTextAppendExpression1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onTextAppendExpression1(stack["first"], stack["rest"])
}

func (c *current) onListAppendExpression1(first, rest interface{}) (interface{}, error) {
return ParseOperator(ListAppendOp, first, rest), nil
}

func (p *parser) callonListAppendExpression1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onListAppendExpression1(stack["first"], stack["rest"])
}

func (c *current) onAndExpression1(first, rest interface{}) (interface{}, error) {
return ParseOperator(AndOp, first, rest), nil
}

func (p *parser) callonAndExpression1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onAndExpression1(stack["first"], stack["rest"])
}

func (c *current) onCombineExpression1(first, rest interface{}) (interface{}, error) {
return ParseOperator(RecordMergeOp, first, rest), nil
}

func (p *parser) callonCombineExpression1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onCombineExpression1(stack["first"], stack["rest"])
}

func (c *current) onPreferExpression1(first, rest interface{}) (interface{}, error) {
return ParseOperator(RightBiasedRecordMergeOp, first, rest), nil
}

func (p *parser) callonPreferExpression1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onPreferExpression1(stack["first"], stack["rest"])
}

func (c *current) onCombineTypesExpression1(first, rest interface{}) (interface{}, error) {
return ParseOperator(RecordTypeMergeOp, first, rest), nil
}

func (p *parser) callonCombineTypesExpression1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onCombineTypesExpression1(stack["first"], stack["rest"])
}

func (c *current) onTimesExpression1(first, rest interface{}) (interface{}, error) {
return ParseOperator(TimesOp, first, rest), nil
}

func (p *parser) callonTimesExpression1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onTimesExpression1(stack["first"], stack["rest"])
}

func (c *current) onEqualExpression1(first, rest interface{}) (interface{}, error) {
return ParseOperator(EqOp, first, rest), nil
}

func (p *parser) callonEqualExpression1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onEqualExpression1(stack["first"], stack["rest"])
}

func (c *current) onNotEqualExpression1(first, rest interface{}) (interface{}, error) {
return ParseOperator(NeOp, first, rest), nil
}

func (p *parser) callonNotEqualExpression1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onNotEqualExpression1(stack["first"], stack["rest"])
}

func (c *current) onApplicationExpression1(f, rest interface{}) (interface{}, error) {
          e := f.(Expr)
          if rest == nil { return e, nil }
          for _, arg := range rest.([]interface{}) {
              e = Apply(e, arg.([]interface{})[1].(Expr))
          }
          return e,nil
      
}

func (p *parser) callonApplicationExpression1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onApplicationExpression1(stack["f"], stack["rest"])
}

func (c *current) onFirstApplicationExpression2(h, u interface{}) (interface{}, error) {
             return Merge{Handler:h.(Expr), Union:u.(Expr)}, nil
          
}

func (p *parser) callonFirstApplicationExpression2() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onFirstApplicationExpression2(stack["h"], stack["u"])
}

func (c *current) onFirstApplicationExpression11(e interface{}) (interface{}, error) {
 return Some{e.(Expr)}, nil 
}

func (p *parser) callonFirstApplicationExpression11() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onFirstApplicationExpression11(stack["e"])
}

func (c *current) onSelectorExpression1(e, ls interface{}) (interface{}, error) {
    expr := e.(Expr)
    labels := ls.([]interface{})
    for _, labelSelector := range labels {
        label := labelSelector.([]interface{})[3]
        expr = Field{expr, label.(string)}
    }
    return expr, nil
}

func (p *parser) callonSelectorExpression1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onSelectorExpression1(stack["e"], stack["ls"])
}

func (c *current) onPrimitiveExpression6(r interface{}) (interface{}, error) {
 return r, nil 
}

func (p *parser) callonPrimitiveExpression6() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onPrimitiveExpression6(stack["r"])
}

func (c *current) onPrimitiveExpression14(u interface{}) (interface{}, error) {
 return u, nil 
}

func (p *parser) callonPrimitiveExpression14() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onPrimitiveExpression14(stack["u"])
}

func (c *current) onPrimitiveExpression24(e interface{}) (interface{}, error) {
 return e, nil 
}

func (p *parser) callonPrimitiveExpression24() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onPrimitiveExpression24(stack["e"])
}

func (c *current) onRecordTypeOrLiteral2() (interface{}, error) {
 return RecordLit(map[string]Expr{}), nil 
}

func (p *parser) callonRecordTypeOrLiteral2() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onRecordTypeOrLiteral2()
}

func (c *current) onRecordTypeOrLiteral6() (interface{}, error) {
 return Record(map[string]Expr{}), nil 
}

func (p *parser) callonRecordTypeOrLiteral6() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onRecordTypeOrLiteral6()
}

func (c *current) onRecordTypeField1(name, expr interface{}) (interface{}, error) {
    return []interface{}{name, expr}, nil
}

func (p *parser) callonRecordTypeField1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onRecordTypeField1(stack["name"], stack["expr"])
}

func (c *current) onMoreRecordType1(f interface{}) (interface{}, error) {
return f, nil
}

func (p *parser) callonMoreRecordType1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onMoreRecordType1(stack["f"])
}

func (c *current) onNonEmptyRecordType1(first, rest interface{}) (interface{}, error) {
          fields := rest.([]interface{})
          content := make(map[string]Expr, len(fields)+1)
          content[first.([]interface{})[0].(string)] = first.([]interface{})[1].(Expr)
          for _, field := range(fields) {
              fieldName := field.([]interface{})[0].(string)
              if _, ok := content[fieldName]; ok {
                  return nil, fmt.Errorf("Duplicate field %s in record", fieldName)
              }
              content[fieldName] = field.([]interface{})[1].(Expr)
          }
          return Record(content), nil
      
}

func (p *parser) callonNonEmptyRecordType1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onNonEmptyRecordType1(stack["first"], stack["rest"])
}

func (c *current) onRecordLiteralField1(name, expr interface{}) (interface{}, error) {
    return []interface{}{name, expr}, nil
}

func (p *parser) callonRecordLiteralField1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onRecordLiteralField1(stack["name"], stack["expr"])
}

func (c *current) onMoreRecordLiteral1(f interface{}) (interface{}, error) {
return f, nil
}

func (p *parser) callonMoreRecordLiteral1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onMoreRecordLiteral1(stack["f"])
}

func (c *current) onNonEmptyRecordLiteral1(first, rest interface{}) (interface{}, error) {
          fields := rest.([]interface{})
          content := make(map[string]Expr, len(fields)+1)
          content[first.([]interface{})[0].(string)] = first.([]interface{})[1].(Expr)
          for _, field := range(fields) {
              fieldName := field.([]interface{})[0].(string)
              if _, ok := content[fieldName]; ok {
                  return nil, fmt.Errorf("Duplicate field %s in record", fieldName)
              }
              content[fieldName] = field.([]interface{})[1].(Expr)
          }
          return RecordLit(content), nil
      
}

func (p *parser) callonNonEmptyRecordLiteral1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onNonEmptyRecordLiteral1(stack["first"], stack["rest"])
}

func (c *current) onEmptyUnionType1() (interface{}, error) {
 return UnionType{}, nil 
}

func (p *parser) callonEmptyUnionType1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onEmptyUnionType1()
}

func (c *current) onNonEmptyUnionType1(first, rest interface{}) (interface{}, error) {
    alternatives := make(map[string]Expr)
    first2 := first.([]interface{})
    if first2[1] == nil {
        alternatives[first2[0].(string)] = nil
    } else {
        alternatives[first2[0].(string)] = first2[1].([]interface{})[3].(Expr)
    }
    if rest == nil { return UnionType(alternatives), nil }
    for _, alternativeSyntax := range rest.([]interface{}) {
        alternative := alternativeSyntax.([]interface{})[3].([]interface{})
        if alternative[1] == nil {
            alternatives[alternative[0].(string)] = nil
        } else {
            alternatives[alternative[0].(string)] = alternative[1].([]interface{})[3].(Expr)
        }
    }
    return UnionType(alternatives), nil
}

func (p *parser) callonNonEmptyUnionType1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onNonEmptyUnionType1(stack["first"], stack["rest"])
}

func (c *current) onMoreList1(e interface{}) (interface{}, error) {
return e, nil
}

func (p *parser) callonMoreList1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onMoreList1(stack["e"])
}

func (c *current) onNonEmptyListLiteral1(first, rest interface{}) (interface{}, error) {
          exprs := rest.([]interface{})
          content := make([]Expr, len(exprs)+1)
          content[0] = first.(Expr)
          for i, expr := range(exprs) {
              content[i+1] = expr.(Expr)
          }
          return NonEmptyList(content), nil
      
}

func (p *parser) callonNonEmptyListLiteral1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onNonEmptyListLiteral1(stack["first"], stack["rest"])
}


var (
	// errNoRule is returned when the grammar to parse has no rule.
	errNoRule          = errors.New("grammar has no rule")

	// errInvalidEncoding is returned when the source is not properly
	// utf8-encoded.
	errInvalidEncoding = errors.New("invalid encoding")

	// errNoMatch is returned if no match could be found.
	errNoMatch         = errors.New("no match found")
)

// Option is a function that can set an option on the parser. It returns
// the previous setting as an Option.
type Option func(*parser) Option

// Debug creates an Option to set the debug flag to b. When set to true,
// debugging information is printed to stdout while parsing.
//
// The default is false.
func Debug(b bool) Option {
	return func(p *parser) Option {
		old := p.debug
		p.debug = b
		return Debug(old)
	}
}

// Memoize creates an Option to set the memoize flag to b. When set to true,
// the parser will cache all results so each expression is evaluated only
// once. This guarantees linear parsing time even for pathological cases,
// at the expense of more memory and slower times for typical cases.
//
// The default is false.
func Memoize(b bool) Option {
	return func(p *parser) Option {
		old := p.memoize
		p.memoize = b
		return Memoize(old)
	}
}

// Recover creates an Option to set the recover flag to b. When set to
// true, this causes the parser to recover from panics and convert it
// to an error. Setting it to false can be useful while debugging to
// access the full stack trace.
//
// The default is true.
func Recover(b bool) Option {
	return func(p *parser) Option {
		old := p.recover
		p.recover = b
		return Recover(old)
	}
}

// ParseFile parses the file identified by filename.
func ParseFile(filename string, opts ...Option) (interface{}, error) {
	f, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	return ParseReader(filename, f, opts...)
}

// ParseReader parses the data from r using filename as information in the
// error messages.
func ParseReader(filename string, r io.Reader, opts ...Option) (interface{}, error) {
	b, err := ioutil.ReadAll(r)
	if err != nil {
		return nil, err
	}

	return Parse(filename, b, opts...)
}

// Parse parses the data from b using filename as information in the
// error messages.
func Parse(filename string, b []byte, opts ...Option) (interface{}, error) {
	return newParser(filename, b, opts...).parse(g)
}

// position records a position in the text.
type position struct {
	line, col, offset int
}

func (p position) String() string {
	return fmt.Sprintf("%d:%d [%d]", p.line, p.col, p.offset)
}

// savepoint stores all state required to go back to this point in the
// parser.
type savepoint struct {
	position
	rn rune
	w  int
}

type current struct {
	pos  position // start position of the match
	text []byte   // raw text of the match
}

// the AST types...

type grammar struct {
	pos   position
	rules []*rule
}

type rule struct {
	pos         position
	name        string
	displayName string
	expr        interface{}
}

type choiceExpr struct {
	pos          position
	alternatives []interface{}
}

type actionExpr struct {
	pos    position
	expr   interface{}
	run    func(*parser) (interface{}, error)
}

type seqExpr struct {
	pos   position
	exprs []interface{}
}

type labeledExpr struct {
	pos   position
	label string
	expr  interface{}
}

type expr struct {
	pos  position
	expr interface{}
}

type andExpr expr
type notExpr expr
type zeroOrOneExpr expr
type zeroOrMoreExpr expr
type oneOrMoreExpr expr

type ruleRefExpr struct {
	pos  position
	name string
}

type andCodeExpr struct {
	pos position
	run func(*parser) (bool, error)
}

type notCodeExpr struct {
	pos position
	run func(*parser) (bool, error)
}

type litMatcher struct {
	pos        position
	val        string
	ignoreCase bool
}

type charClassMatcher struct {
	pos        position
	val        string
	chars      []rune
	ranges     []rune
	classes    []*unicode.RangeTable
	ignoreCase bool
	inverted   bool
}

type anyMatcher position

// errList cumulates the errors found by the parser.
type errList []error

func (e *errList) add(err error) {
	*e = append(*e, err)
}

func (e errList) err() error {
	if len(e) == 0 {
		return nil
	}
	e.dedupe()
	return e
}

func (e *errList) dedupe() {
	var cleaned []error
	set := make(map[string]bool)
	for _, err := range *e {
		if msg := err.Error(); !set[msg] {
			set[msg] = true
			cleaned = append(cleaned, err)
		}
	}
	*e = cleaned
}

func (e errList) Error() string {
	switch len(e) {
	case 0:
		return ""
	case 1:
		return e[0].Error()
	default:
		var buf bytes.Buffer

		for i, err := range e {
			if i > 0 {
				buf.WriteRune('\n')
			}
			buf.WriteString(err.Error())
		}
		return buf.String()
	}
}

// parserError wraps an error with a prefix indicating the rule in which
// the error occurred. The original error is stored in the Inner field.
type parserError struct {
	Inner  error
	pos    position
	prefix string
}

// Error returns the error message.
func (p *parserError) Error() string {
	return p.prefix + ": " + p.Inner.Error()
}

// newParser creates a parser with the specified input source and options.
func newParser(filename string, b []byte, opts ...Option) *parser {
	p := &parser{
		filename: filename,
		errs: new(errList),
		data: b,
		pt: savepoint{position: position{line: 1}},
		recover: true,
	}
	p.setOptions(opts)
	return p
}

// setOptions applies the options to the parser.
func (p *parser) setOptions(opts []Option) {
	for _, opt := range opts {
		opt(p)
	}
}

type resultTuple struct {
	v interface{}
	b bool
	end savepoint
}

type parser struct {
	filename string
	pt       savepoint
	cur      current

	data []byte
	errs *errList

	recover bool
	debug bool
	depth  int

	memoize bool
	// memoization table for the packrat algorithm:
	// map[offset in source] map[expression or rule] {value, match}
	memo map[int]map[interface{}]resultTuple

	// rules table, maps the rule identifier to the rule node
	rules  map[string]*rule
	// variables stack, map of label to value
	vstack []map[string]interface{}
	// rule stack, allows identification of the current rule in errors
	rstack []*rule

	// stats
	exprCnt int
}

// push a variable set on the vstack.
func (p *parser) pushV() {
	if cap(p.vstack) == len(p.vstack) {
		// create new empty slot in the stack
		p.vstack = append(p.vstack, nil)
	} else {
		// slice to 1 more
		p.vstack = p.vstack[:len(p.vstack)+1]
	}

	// get the last args set
	m := p.vstack[len(p.vstack)-1]
	if m != nil && len(m) == 0 {
		// empty map, all good
		return
	}

	m = make(map[string]interface{})
	p.vstack[len(p.vstack)-1] = m
}

// pop a variable set from the vstack.
func (p *parser) popV() {
	// if the map is not empty, clear it
	m := p.vstack[len(p.vstack)-1]
	if len(m) > 0 {
		// GC that map
		p.vstack[len(p.vstack)-1] = nil
	}
	p.vstack = p.vstack[:len(p.vstack)-1]
}

func (p *parser) print(prefix, s string) string {
	if !p.debug {
		return s
	}

	fmt.Printf("%s %d:%d:%d: %s [%#U]\n",
		prefix, p.pt.line, p.pt.col, p.pt.offset, s, p.pt.rn)
	return s
}

func (p *parser) in(s string) string {
	p.depth++
	return p.print(strings.Repeat(" ", p.depth) + ">", s)
}

func (p *parser) out(s string) string {
	p.depth--
	return p.print(strings.Repeat(" ", p.depth) + "<", s)
}

func (p *parser) addErr(err error) {
	p.addErrAt(err, p.pt.position)
}

func (p *parser) addErrAt(err error, pos position) {
	var buf bytes.Buffer
	if p.filename != "" {
		buf.WriteString(p.filename)
	}
	if buf.Len() > 0 {
		buf.WriteString(":")
	}
	buf.WriteString(fmt.Sprintf("%d:%d (%d)", pos.line, pos.col, pos.offset))
	if len(p.rstack) > 0 {
		if buf.Len() > 0 {
			buf.WriteString(": ")
		}
		rule := p.rstack[len(p.rstack)-1]
		if rule.displayName != "" {
			buf.WriteString("rule " + rule.displayName)
		} else {
			buf.WriteString("rule " + rule.name)
		}
	}
	pe := &parserError{Inner: err, pos: pos, prefix: buf.String()}
	p.errs.add(pe)
}

// read advances the parser to the next rune.
func (p *parser) read() {
	p.pt.offset += p.pt.w
	rn, n := utf8.DecodeRune(p.data[p.pt.offset:])
	p.pt.rn = rn
	p.pt.w = n
	p.pt.col++
	if rn == '\n' {
		p.pt.line++
		p.pt.col = 0
	}

	if rn == utf8.RuneError {
		if n == 1 {
			p.addErr(errInvalidEncoding)
		}
	}
}

// restore parser position to the savepoint pt.
func (p *parser) restore(pt savepoint) {
	if p.debug {
		defer p.out(p.in("restore"))
	}
	if pt.offset == p.pt.offset {
		return
	}
	p.pt = pt
}

// get the slice of bytes from the savepoint start to the current position.
func (p *parser) sliceFrom(start savepoint) []byte {
	return p.data[start.position.offset:p.pt.position.offset]
}

func (p *parser) getMemoized(node interface{}) (resultTuple, bool) {
	if len(p.memo) == 0 {
		return resultTuple{}, false
	}
	m := p.memo[p.pt.offset]
	if len(m) == 0 {
		return resultTuple{}, false
	}
	res, ok := m[node]
	return res, ok
}

func (p *parser) setMemoized(pt savepoint, node interface{}, tuple resultTuple) {
	if p.memo == nil {
		p.memo = make(map[int]map[interface{}]resultTuple)
	}
	m := p.memo[pt.offset]
	if m == nil {
		m = make(map[interface{}]resultTuple)
		p.memo[pt.offset] = m
	}
	m[node] = tuple
}

func (p *parser) buildRulesTable(g *grammar) {
	p.rules = make(map[string]*rule, len(g.rules))
	for _, r := range g.rules {
		p.rules[r.name] = r
	}
}

func (p *parser) parse(g *grammar) (val interface{}, err error) {
	if len(g.rules) == 0 {
		p.addErr(errNoRule)
		return nil, p.errs.err()
	}

	// TODO : not super critical but this could be generated
	p.buildRulesTable(g)

	if p.recover {
		// panic can be used in action code to stop parsing immediately
		// and return the panic as an error.
		defer func() {
			if e := recover(); e != nil {
				if p.debug {
					defer p.out(p.in("panic handler"))
				}
				val = nil
				switch e := e.(type) {
				case error:
					p.addErr(e)
				default:
					p.addErr(fmt.Errorf("%v", e))
				}
				err = p.errs.err()
			}
		}()
	}

	// start rule is rule [0]
	p.read() // advance to first rune
	val, ok := p.parseRule(g.rules[0])
	if !ok {
		if len(*p.errs) == 0 {
			// make sure this doesn't go out silently
			p.addErr(errNoMatch)
		}
		return nil, p.errs.err()
	}
	return val, p.errs.err()
}

func (p *parser) parseRule(rule *rule) (interface{}, bool) {
	if p.debug {
		defer p.out(p.in("parseRule " + rule.name))
	}

	if p.memoize {
		res, ok := p.getMemoized(rule)
		if ok {
			p.restore(res.end)
			return res.v, res.b
		}
	}

	start := p.pt
	p.rstack = append(p.rstack, rule)
	p.pushV()
	val, ok := p.parseExpr(rule.expr)
	p.popV()
	p.rstack = p.rstack[:len(p.rstack)-1]
	if ok && p.debug {
		p.print(strings.Repeat(" ", p.depth) + "MATCH", string(p.sliceFrom(start)))
	}

	if p.memoize {
		p.setMemoized(start, rule, resultTuple{val, ok, p.pt})
	}
	return val, ok
}

func (p *parser) parseExpr(expr interface{}) (interface{}, bool) {
	var pt savepoint
	var ok bool

	if p.memoize {
		res, ok := p.getMemoized(expr)
		if ok {
			p.restore(res.end)
			return res.v, res.b
		}
		pt = p.pt
	}

	p.exprCnt++
	var val interface{}
	switch expr := expr.(type) {
	case *actionExpr:
		val, ok = p.parseActionExpr(expr)
	case *andCodeExpr:
		val, ok = p.parseAndCodeExpr(expr)
	case *andExpr:
		val, ok = p.parseAndExpr(expr)
	case *anyMatcher:
		val, ok = p.parseAnyMatcher(expr)
	case *charClassMatcher:
		val, ok = p.parseCharClassMatcher(expr)
	case *choiceExpr:
		val, ok = p.parseChoiceExpr(expr)
	case *labeledExpr:
		val, ok = p.parseLabeledExpr(expr)
	case *litMatcher:
		val, ok = p.parseLitMatcher(expr)
	case *notCodeExpr:
		val, ok = p.parseNotCodeExpr(expr)
	case *notExpr:
		val, ok = p.parseNotExpr(expr)
	case *oneOrMoreExpr:
		val, ok = p.parseOneOrMoreExpr(expr)
	case *ruleRefExpr:
		val, ok = p.parseRuleRefExpr(expr)
	case *seqExpr:
		val, ok = p.parseSeqExpr(expr)
	case *zeroOrMoreExpr:
		val, ok = p.parseZeroOrMoreExpr(expr)
	case *zeroOrOneExpr:
		val, ok = p.parseZeroOrOneExpr(expr)
	default:
		panic(fmt.Sprintf("unknown expression type %T", expr))
	}
	if p.memoize {
		p.setMemoized(pt, expr, resultTuple{val, ok, p.pt})
	}
	return val, ok
}

func (p *parser) parseActionExpr(act *actionExpr) (interface{}, bool) {
	if p.debug {
		defer p.out(p.in("parseActionExpr"))
	}

	start := p.pt
	val, ok := p.parseExpr(act.expr)
	if ok {
		p.cur.pos = start.position
		p.cur.text = p.sliceFrom(start)
		actVal, err := act.run(p)
		if err != nil {
			p.addErrAt(err, start.position)
		}
		val = actVal
	}
	if ok && p.debug {
		p.print(strings.Repeat(" ", p.depth) + "MATCH", string(p.sliceFrom(start)))
	}
	return val, ok
}

func (p *parser) parseAndCodeExpr(and *andCodeExpr) (interface{}, bool) {
	if p.debug {
		defer p.out(p.in("parseAndCodeExpr"))
	}

	ok, err := and.run(p)
	if err != nil {
		p.addErr(err)
	}
	return nil, ok
}

func (p *parser) parseAndExpr(and *andExpr) (interface{}, bool) {
	if p.debug {
		defer p.out(p.in("parseAndExpr"))
	}

	pt := p.pt
	p.pushV()
	_, ok := p.parseExpr(and.expr)
	p.popV()
	p.restore(pt)
	return nil, ok
}

func (p *parser) parseAnyMatcher(any *anyMatcher) (interface{}, bool) {
	if p.debug {
		defer p.out(p.in("parseAnyMatcher"))
	}

	if p.pt.rn != utf8.RuneError {
		start := p.pt
		p.read()
		return p.sliceFrom(start), true
	}
	return nil, false
}

func (p *parser) parseCharClassMatcher(chr *charClassMatcher) (interface{}, bool) {
	if p.debug {
		defer p.out(p.in("parseCharClassMatcher"))
	}

	cur := p.pt.rn
	// can't match EOF
	if cur == utf8.RuneError {
		return nil, false
	}
	start := p.pt
	if chr.ignoreCase {
		cur = unicode.ToLower(cur)
	}

	// try to match in the list of available chars
	for _, rn := range chr.chars {
		if rn == cur {
			if chr.inverted {
				return nil, false
			}
			p.read()
			return p.sliceFrom(start), true
		}
	}

	// try to match in the list of ranges
	for i := 0; i < len(chr.ranges); i += 2 {
		if cur >= chr.ranges[i] && cur <= chr.ranges[i+1] {
			if chr.inverted {
				return nil, false
			}
			p.read()
			return p.sliceFrom(start), true
		}
	}

	// try to match in the list of Unicode classes
	for _, cl := range chr.classes {
		if unicode.Is(cl, cur) {
			if chr.inverted {
				return nil, false
			}
			p.read()
			return p.sliceFrom(start), true
		}
	}

	if chr.inverted {
		p.read()
		return p.sliceFrom(start), true
	}
	return nil, false
}

func (p *parser) parseChoiceExpr(ch *choiceExpr) (interface{}, bool) {
	if p.debug {
		defer p.out(p.in("parseChoiceExpr"))
	}

	for _, alt := range ch.alternatives {
		p.pushV()
		val, ok := p.parseExpr(alt)
		p.popV()
		if ok {
			return val, ok
		}
	}
	return nil, false
}

func (p *parser) parseLabeledExpr(lab *labeledExpr) (interface{}, bool) {
	if p.debug {
		defer p.out(p.in("parseLabeledExpr"))
	}

	p.pushV()
	val, ok := p.parseExpr(lab.expr)
	p.popV()
	if ok && lab.label != "" {
		m := p.vstack[len(p.vstack)-1]
		m[lab.label] = val
	}
	return val, ok
}

func (p *parser) parseLitMatcher(lit *litMatcher) (interface{}, bool) {
	if p.debug {
		defer p.out(p.in("parseLitMatcher"))
	}

	start := p.pt
	for _, want := range lit.val {
		cur := p.pt.rn
		if lit.ignoreCase {
			cur = unicode.ToLower(cur)
		}
		if cur != want {
			p.restore(start)
			return nil, false
		}
		p.read()
	}
	return p.sliceFrom(start), true
}

func (p *parser) parseNotCodeExpr(not *notCodeExpr) (interface{}, bool) {
	if p.debug {
		defer p.out(p.in("parseNotCodeExpr"))
	}

	ok, err := not.run(p)
	if err != nil {
		p.addErr(err)
	}
	return nil, !ok
}

func (p *parser) parseNotExpr(not *notExpr) (interface{}, bool) {
	if p.debug {
		defer p.out(p.in("parseNotExpr"))
	}

	pt := p.pt
	p.pushV()
	_, ok := p.parseExpr(not.expr)
	p.popV()
	p.restore(pt)
	return nil, !ok
}

func (p *parser) parseOneOrMoreExpr(expr *oneOrMoreExpr) (interface{}, bool) {
	if p.debug {
		defer p.out(p.in("parseOneOrMoreExpr"))
	}

	var vals []interface{}

	for {
		p.pushV()
		val, ok := p.parseExpr(expr.expr)
		p.popV()
		if !ok {
			if len(vals) == 0 {
				// did not match once, no match
				return nil, false
			}
			return vals, true
		}
		vals = append(vals, val)
	}
}

func (p *parser) parseRuleRefExpr(ref *ruleRefExpr) (interface{}, bool) {
	if p.debug {
		defer p.out(p.in("parseRuleRefExpr " + ref.name))
	}

	if ref.name == "" {
		panic(fmt.Sprintf("%s: invalid rule: missing name", ref.pos))
	}

	rule := p.rules[ref.name]
	if rule == nil {
		p.addErr(fmt.Errorf("undefined rule: %s", ref.name))
		return nil, false
	}
	return p.parseRule(rule)
}

func (p *parser) parseSeqExpr(seq *seqExpr) (interface{}, bool) {
	if p.debug {
		defer p.out(p.in("parseSeqExpr"))
	}

	var vals []interface{}

	pt := p.pt
	for _, expr := range seq.exprs {
		val, ok := p.parseExpr(expr)
		if !ok {
			p.restore(pt)
			return nil, false
		}
		vals = append(vals, val)
	}
	return vals, true
}

func (p *parser) parseZeroOrMoreExpr(expr *zeroOrMoreExpr) (interface{}, bool) {
	if p.debug {
		defer p.out(p.in("parseZeroOrMoreExpr"))
	}

	var vals []interface{}

	for {
		p.pushV()
		val, ok := p.parseExpr(expr.expr)
		p.popV()
		if !ok {
			return vals, true
		}
		vals = append(vals, val)
	}
}

func (p *parser) parseZeroOrOneExpr(expr *zeroOrOneExpr) (interface{}, bool) {
	if p.debug {
		defer p.out(p.in("parseZeroOrOneExpr"))
	}

	p.pushV()
	val, _ := p.parseExpr(expr.expr)
	p.popV()
	// whether it matched or not, consider it a match
	return val, true
}

func rangeTable(class string) *unicode.RangeTable {
	if rt, ok := unicode.Categories[class]; ok {
		return rt
	}
	if rt, ok := unicode.Properties[class]; ok {
		return rt
	}
	if rt, ok := unicode.Scripts[class]; ok {
		return rt
	}

	// cannot happen
	panic(fmt.Sprintf("invalid Unicode class: %s", class))
}

