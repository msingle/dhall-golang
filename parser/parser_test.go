package parser_test

import (
	"math"

	. "github.com/philandstuff/dhall-golang/ast"
	"github.com/philandstuff/dhall-golang/parser"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/ginkgo/extensions/table"
	. "github.com/onsi/gomega"
)

func ParseAndCompare(input string, expected interface{}) {
	root, err := parser.Parse("test", []byte(input))
	Expect(err).ToNot(HaveOccurred())
	Expect(root).To(Equal(expected))
}

func ParseAndFail(input string) {
	_, err := parser.Parse("test", []byte(input))
	Expect(err).To(HaveOccurred())
}

var _ = Describe("Expression", func() {
	DescribeTable("simple expressions", ParseAndCompare,
		Entry("Type", `Type`, Type),
		Entry("Kind", `Kind`, Kind),
		Entry("Sort", `Sort`, Sort),
		Entry("Double", `Double`, Double),
		Entry("Text", `Text`, Text),
		Entry("DoubleLit", `3.0`, DoubleLit(3.0)),
		Entry("DoubleLit with exponent", (`3E5`), DoubleLit(3e5)),
		Entry("DoubleLit with sign", (`+3.0`), DoubleLit(3.0)),
		Entry("DoubleLit with everything", (`-5.0e1`), DoubleLit(-50.0)),
		Entry("Infinity", `Infinity`, DoubleLit(math.Inf(1))),
		Entry("-Infinity", `-Infinity`, DoubleLit(math.Inf(-1))),
		Entry("Integer", `Integer`, Integer),
		Entry("IntegerLit", `+1234`, IntegerLit(1234)),
		Entry("IntegerLit", `-3`, IntegerLit(-3)),
		Entry("Annotated expression", `3 : Natural`, Annot{NaturalLit(3), Natural}),
	)
	DescribeTable("bools", ParseAndCompare,
		Entry("Bool", `Bool`, Bool),
		Entry("True", `True`, BoolLit(true)),
		Entry("False", `False`, BoolLit(false)),
		Entry("if True then x else y", `if True then x else y`, BoolIf{True, MkVar("x"), MkVar("y")}),
	)
	DescribeTable("naturals", ParseAndCompare,
		Entry("Natural", `Natural`, Natural),
		Entry("NaturalLit", `1234`, NaturalLit(1234)),
		Entry("NaturalLit", `3`, NaturalLit(3)),
		Entry("NaturalPlus", `3 + 5`, NaturalPlus(NaturalLit(3), NaturalLit(5))),
		Entry("NaturalTimes", `3 * 5`, NaturalTimes(NaturalLit(3), NaturalLit(5))),
		Entry("Natural op order #1", `3 * 4 + 5`, NaturalPlus(NaturalTimes(NaturalLit(3), NaturalLit(4)), NaturalLit(5))),
		Entry("Natural op order #2", `3 + 4 * 5`, NaturalPlus(NaturalLit(3), NaturalTimes(NaturalLit(4), NaturalLit(5)))),
		// Check that if we skip whitespace, it parses
		// correctly as function application, not natural
		// addition
		Entry("Plus without whitespace", `3 +5`, Apply(NaturalLit(3), IntegerLit(5))),
	)
	DescribeTable("double-quoted text literals", ParseAndCompare,
		Entry("Empty TextLit", `""`, TextLit{}),
		Entry("Simple TextLit", `"foo"`, TextLit{Suffix: "foo"}),
		Entry(`TextLit escape "`, `"\""`, TextLit{Suffix: `"`}),
		Entry(`TextLit escape $`, `"\$"`, TextLit{Suffix: `$`}),
		Entry(`TextLit escape \`, `"\\"`, TextLit{Suffix: `\`}),
		Entry(`TextLit escape /`, `"\/"`, TextLit{Suffix: `/`}),
		Entry(`TextLit escape \b`, `"\b"`, TextLit{Suffix: "\b"}),
		Entry(`TextLit escape \f`, `"\f"`, TextLit{Suffix: "\f"}),
		Entry(`TextLit escape \n`, `"\n"`, TextLit{Suffix: "\n"}),
		Entry(`TextLit escape \r`, `"\r"`, TextLit{Suffix: "\r"}),
		Entry(`TextLit escape \t`, `"\t"`, TextLit{Suffix: "\t"}),
		Entry(`TextLit escape \u2200`, `"\u2200"`, TextLit{Suffix: "∀"}),
		Entry(`TextLit escape \u03bb`, `"\u03bb"`, TextLit{Suffix: "λ"}),
		Entry(`TextLit escape \u03BB`, `"\u03BB"`, TextLit{Suffix: "λ"}),
		Entry("Interpolated TextLit", `"foo ${"bar"} baz"`,
			TextLit{Chunks{Chunk{"foo ", TextLit{Suffix: "bar"}}},
				" baz"},
		),
	)
	DescribeTable("single-quoted text literals", ParseAndCompare,
		Entry("Empty TextLit", `''
''`, TextLit{}),
		Entry("Simple TextLit with no newlines", `''
foo''`, TextLit{Suffix: "foo"}),
		Entry("Simple TextLit with newlines", `''
foo
''`, TextLit{Suffix: "foo\n"}),
		Entry("TextLit with space indent", `''
  foo
  bar
  ''`, TextLit{Suffix: "foo\nbar\n"}),
		Entry("TextLit with tab indent", `''
		foo
		bar
		''`, TextLit{Suffix: "foo\nbar\n"}),
		Entry("TextLit with mixed tab/space indent", `''
	  foo
	  bar
	  ''`, TextLit{Suffix: "foo\nbar\n"}),
		Entry("TextLit with weird indenting", `''
	    foo
	  	bar
	  ''`, TextLit{Suffix: "  foo\n\tbar\n"}),
		Entry(`Escape ''`, `''
'''
''`, TextLit{Suffix: "''\n"}),
		Entry(`Escape ${`, `''
''${
''`, TextLit{Suffix: "${\n"}),
		Entry("Interpolation", `''
foo ${"bar"}
baz
''`,
			TextLit{Chunks{Chunk{"foo ", TextLit{Suffix: "bar"}}},
				"\nbaz\n"},
		),
	)
	DescribeTable("simple expressions", ParseAndCompare,
		Entry("Identifier", `x`, MkVar("x")),
		Entry("Identifier with index", `x@1`, Var{"x", 1}),
		Entry("Identifier with reserved prefix", `Listicle`, MkVar("Listicle")),
		Entry("Identifier with reserved prefix and index", `Listicle@3`, Var{"Listicle", 3}),
	)
	DescribeTable("lists", ParseAndCompare,
		Entry("List Natural", `List Natural`, Apply(List, Natural)),
		Entry("[3]", `[3]`, MakeList(NaturalLit(3))),
		Entry("[3,4]", `[3,4]`, MakeList(NaturalLit(3), NaturalLit(4))),
		Entry("[] : List Natural", `[] : List Natural`, EmptyList{Apply(List, Natural)}),
		Entry("[3] : List Natural", `[3] : List Natural`, Annot{MakeList(NaturalLit(3)), Apply(List, Natural)}),
		Entry("a # b", `a # b`, ListAppend(MkVar("a"), MkVar("b"))),
	)
	DescribeTable("optionals", ParseAndCompare,
		Entry("Optional Natural", `Optional Natural`, Apply(Optional, Natural)),
		Entry("Some 3", `Some 3`, Some{NaturalLit(3)}),
		Entry("None Natural", `None Natural`, Apply(None, Natural)),
	)
	DescribeTable("records", ParseAndCompare,
		Entry("{}", `{}`, Record{}),
		Entry("{=}", `{=}`, RecordLit{}),
		Entry("{foo : Natural}", `{foo : Natural}`, Record{"foo": Natural}),
		Entry("{foo = 3}", `{foo = 3}`, RecordLit{"foo": NaturalLit(3)}),
		Entry("{foo : Natural, bar : Integer}", `{foo : Natural, bar: Integer}`, Record{"foo": Natural, "bar": Integer}),
		Entry("{foo = 3 , bar = +3}", `{foo = 3 , bar = +3}`, RecordLit{"foo": NaturalLit(3), "bar": IntegerLit(3)}),
		Entry("t.x", `t.x`, Field{MkVar("t"), "x"}),
		Entry("t.x.y", `t.x.y`, Field{Field{MkVar("t"), "x"}, "y"}),
	)
	DescribeTable("imports", ParseAndCompare,
		Entry("bash envvar text import", `env:FOO as Text`, Embed(MakeEnvVarImport("FOO", RawText))),
		Entry("posix envvar text import", `env:"FOO" as Text`, Embed(MakeEnvVarImport("FOO", RawText))),
		Entry("posix envvar text import", `env:"foo\nbar\a!" as Text`, Embed(MakeEnvVarImport("foo\nbar\a!", RawText))),
		Entry("bash envvar code import", `env:FOO`, Embed(MakeEnvVarImport("FOO", Code))),
		Entry("posix envvar code import", `env:"FOO"`, Embed(MakeEnvVarImport("FOO", Code))),
		Entry("posix envvar code import", `env:"foo\nbar\a!"`, Embed(MakeEnvVarImport("foo\nbar\a!", Code))),
		Entry("missing", `missing`, Embed(MakeImport(Missing(struct{}{}), Code))),
		Entry("local here-path import", `./local`, Embed(MakeLocalImport("local", Code))),
		Entry("local parent-path import", `../local`, Embed(MakeLocalImport("../local", Code))),
		Entry("local home import", `~/in/home`, Embed(MakeLocalImport("~/in/home", Code))),
		Entry("local absolute import", `/local`, Embed(MakeLocalImport("/local", Code))),
		Entry("simple remote", `https://example.com/foo`, Embed(MakeRemoteImport("https://example.com/foo", Code))),
		Entry("http remote", `http://example.com/foo`, Embed(MakeRemoteImport("http://example.com/foo", Code))),
		Entry("remote with query string", `https://example.com/foo?bar=baz&fred=jim`, Embed(MakeRemoteImport("https://example.com/foo?bar=baz&fred=jim", Code))),
		Entry("remote with port", `https://example.com:8080/foo`, Embed(MakeRemoteImport("https://example.com:8080/foo", Code))),
		Entry("remote with userinfo", `https://foo:bar@example.com/foo`, Embed(MakeRemoteImport("https://foo:bar@example.com/foo", Code))),
		Entry("remote with IPv4 address", `https://127.0.0.1/foo`, Embed(MakeRemoteImport("https://127.0.0.1/foo", Code))),
		Entry("remote with IPv6 address", `https://[cafe:d00d::1234]/foo`, Embed(MakeRemoteImport("https://[cafe:d00d::1234]/foo", Code))),
		// unimplemented yet. don't care too much about these features
		PEntry("remote with headers", ``, nil),
	)
	// can't test NaN using ParseAndCompare because NaN ≠ NaN
	It("handles NaN correctly", func() {
		root, err := parser.Parse("test", []byte(`NaN`))
		Expect(err).ToNot(HaveOccurred())
		f := float64(root.(DoubleLit))
		Expect(math.IsNaN(f)).To(BeTrue())
	})
	DescribeTable("lambda expressions", ParseAndCompare,
		Entry("simple λ",
			`λ(foo : bar) → baz`,
			&LambdaExpr{
				"foo", MkVar("bar"), MkVar("baz")}),
		Entry(`simple \`,
			`\(foo : bar) → baz`,
			&LambdaExpr{
				"foo", MkVar("bar"), MkVar("baz")}),
		Entry("with line comment",
			"λ(foo : bar) --asdf\n → baz",
			&LambdaExpr{
				"foo", MkVar("bar"), MkVar("baz")}),
		Entry("with block comment",
			"λ(foo : bar) {-asdf\n-} → baz",
			&LambdaExpr{
				"foo", MkVar("bar"), MkVar("baz")}),
		Entry("simple ∀",
			`∀(foo : bar) → baz`,
			&Pi{"foo", MkVar("bar"), MkVar("baz")}),
		Entry("arrow type has implicit _ var",
			`foo → bar`,
			FnType(MkVar("foo"), MkVar("bar"))),
		Entry(`simple forall`,
			`forall(foo : bar) → baz`,
			&Pi{"foo", MkVar("bar"), MkVar("baz")}),
		Entry("with line comment",
			"∀(foo : bar) --asdf\n → baz",
			&Pi{"foo", MkVar("bar"), MkVar("baz")}),
	)
	DescribeTable("applications", ParseAndCompare,
		Entry("identifier application",
			`foo bar`,
			Apply(
				MkVar("foo"),
				MkVar("bar"),
			)),
		Entry("lambda application",
			`(λ(foo : bar) → baz) quux`,
			Apply(
				&LambdaExpr{
					"foo", MkVar("bar"), MkVar("baz")},
				MkVar("quux"))),
	)
	DescribeTable("lets", ParseAndCompare,
		Entry("simple let",
			`let x = y in z`,
			MakeLet(MkVar("z"), Binding{
				Variable: "x", Value: MkVar("y"),
			})),
		Entry("lambda application",
			`(λ(foo : bar) → baz) quux`,
			Apply(
				&LambdaExpr{
					"foo", MkVar("bar"), MkVar("baz")},
				MkVar("quux"))),
	)
	Describe("Expected failures", func() {
		// these keywords should fail to parse unless they're part of
		// a larger expression
		DescribeTable("keywords", ParseAndFail,
			Entry("if", `if`),
			Entry("then", `then`),
			Entry("else", `else`),
			Entry("let", `let`),
			Entry("in", `in`),
			Entry("using", `using`),
			Entry("as", `as`),
			Entry("merge", `merge`),
			Entry("Some", `Some`),
		)
		DescribeTable("bad URLs", ParseAndFail,
			Entry("bad IPv6", `https://[11111::22222]/abc`),
		)
		DescribeTable("other expected failures", ParseAndFail,
			Entry("annotation without required space", `3 :Natural`),
			Entry("unannotated list", `[]`),
		)
	})
})
