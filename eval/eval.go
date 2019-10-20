package eval

import (
	"fmt"

	. "github.com/philandstuff/dhall-golang/core"
)

type Env map[string][]Value

// Eval normalizes Term to a Value.
func Eval(t Term) Value {
	return evalWith(t, Env{}, false)
}

// AlphaBetaEval alpha-beta-normalizes Term to a Value.
func AlphaBetaEval(t Term) Value {
	return evalWith(t, Env{}, true)
}

func evalWith(t Term, e Env, shouldAlphaNormalize bool) Value {
	switch t := t.(type) {
	case Universe:
		return t
	case Builtin:
		return t
	case BoundVar:
		if t.Index >= len(e[t.Name]) {
			panic(fmt.Sprintf("Eval: unbound variable %s", t))
		}
		return e[t.Name][t.Index]
	case LocalVar:
		return t
	case FreeVar:
		return t
	case LambdaTerm:
		v := LambdaValue{
			Label:  t.Label,
			Domain: evalWith(t.Type, e, shouldAlphaNormalize),
			Fn: func(x Value) Value {
				newEnv := Env{}
				for k, v := range e {
					newEnv[k] = v
				}
				newEnv[t.Label] = append([]Value{x}, newEnv[t.Label]...)
				return evalWith(t.Body, newEnv, shouldAlphaNormalize)
			}}
		if shouldAlphaNormalize {
			v.Label = "_"
		}
		return v
	case PiTerm:
		v := PiValue{
			Label:  t.Label,
			Domain: evalWith(t.Type, e, shouldAlphaNormalize),
			Range: func(x Value) Value {
				newEnv := Env{}
				for k, v := range e {
					newEnv[k] = v
				}
				newEnv[t.Label] = append([]Value{x}, newEnv[t.Label]...)
				return evalWith(t.Body, newEnv, shouldAlphaNormalize)
			}}
		if shouldAlphaNormalize {
			v.Label = "_"
		}
		return v
	case AppTerm:
		fn := evalWith(t.Fn, e, shouldAlphaNormalize)
		arg := evalWith(t.Arg, e, shouldAlphaNormalize)
		if f, ok := fn.(LambdaValue); ok {
			return f.Fn(arg)
		}
		return AppValue{
			Fn:  fn,
			Arg: arg,
		}
	case Let:
		return TextLitVal{Suffix: "let unimplemented"}
	case Annot:
		return evalWith(t.Expr, e, shouldAlphaNormalize)
	case DoubleLit:
		return t
	case TextLitTerm:
		return TextLitVal{Suffix: "TextLit unimplemented but here's the suffix: " + t.Suffix}
	case BoolLit:
		return t
	case IfTerm:
		return TextLitVal{Suffix: "If unimplemented"}
	case NaturalLit:
		return t
	case IntegerLit:
		return t
	case OpTerm:
		return TextLitVal{Suffix: "OpTerm unimplemented"}
	case EmptyList:
		return EmptyListVal{Type: evalWith(t.Type, e, shouldAlphaNormalize)}
	case NonEmptyList:
		return TextLitVal{Suffix: "NonEmptyList unimplemented"}
	case Some:
		return SomeVal{evalWith(t.Val, e, shouldAlphaNormalize)}
	case RecordType:
		return TextLitVal{Suffix: "RecordType unimplemented"}
	case RecordLit:
		return TextLitVal{Suffix: "RecordLit unimplemented"}
	case ToMap:
		return TextLitVal{Suffix: "ToMap unimplemented"}
	case Field:
		return TextLitVal{Suffix: "Field unimplemented"}
	case Project:
		return TextLitVal{Suffix: "Project unimplemented"}
	case ProjectType:
		return TextLitVal{Suffix: "ProjectType unimplemented"}
	case UnionType:
		return TextLitVal{Suffix: "UnionType unimplemented"}
	case Merge:
		return TextLitVal{Suffix: "Merge unimplemented"}
	case Assert:
		return AssertVal{Annotation: evalWith(t.Annotation, e, shouldAlphaNormalize)}
	default:
		panic(fmt.Sprint("unknown term type", t))
	}
}