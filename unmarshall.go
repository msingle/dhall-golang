package main

import (
	"reflect"

	"github.com/philandstuff/dhall-golang/ast"
)

func Unmarshal(e ast.Expr, out interface{}) {
	v := reflect.ValueOf(out)
	unmarshal(e, v.Elem())
}

func reflectValToDhallVal(val reflect.Value, typ ast.Expr) ast.Expr {
	typ = typ.Normalize()
	switch e := typ.(type) {
	case ast.Builtin:
		switch typ {
		case ast.Double:
			return ast.DoubleLit(val.Float())
		case ast.Bool:
			return ast.BoolLit(val.Bool())
		case ast.Natural:
			return ast.NaturalLit(val.Int())
		case ast.Integer:
			return ast.IntegerLit(val.Int())
			// TODO: TextLit
		case ast.List:
			panic("wrong Kind")
		default:
			panic("unknown Builtin")
		}
		// TODO: RecordLit
	case *ast.App:
		switch e.Fn {
		case ast.List:
			if val.Len() == 0 {
				return ast.EmptyList{Type: e.Arg}
			}
			slice := make([]ast.Expr, val.Len())
			for i := 0; i < val.Len(); i++ {
				slice[i] = reflectValToDhallVal(val.Index(i), e.Arg)
			}
			return ast.NonEmptyList(slice)
		default:
			panic("unknown app")
		}
	default:
		panic("Don't know what to do with val")
	}
}

func argNType(fn *ast.LambdaExpr, n int) ast.Expr {
	if n == 0 {
		return fn.Type
	}
	return argNType(fn.Body.(*ast.LambdaExpr), n-1)
}

func dhallShim(out reflect.Type, dhallFunc *ast.LambdaExpr) func([]reflect.Value) []reflect.Value {
	return func(args []reflect.Value) []reflect.Value {
		var expr ast.Expr = dhallFunc
		for i, arg := range args {
			dhallArg := reflectValToDhallVal(arg, argNType(dhallFunc, i))
			expr = &ast.App{Fn: expr, Arg: dhallArg}
		}
		ptr := reflect.New(out)
		unmarshal(expr.Normalize(), ptr.Elem())
		return []reflect.Value{ptr.Elem()}
	}
}

func unmarshal(e ast.Expr, v reflect.Value) {
	switch v.Kind() {
	case reflect.Map:
		// initialise with new (non-nil) value
		v.Set(reflect.MakeMap(v.Type()))
		if _, ok := e.(ast.EmptyList); ok {
			return
		}
		e := e.(ast.NonEmptyList)
		recordLit := e[0].(ast.RecordLit)
		if len(recordLit) != 2 {
			panic("can only unmarshal `List {mapKey : T, mapValue : U}` into go map")
		}
		for _, r := range e {
			entry := r.(ast.RecordLit)
			key := reflect.New(v.Type().Key()).Elem()
			val := reflect.New(v.Type().Elem()).Elem()
			unmarshal(entry["mapKey"], key)
			unmarshal(entry["mapValue"], val)
			v.SetMapIndex(key, val)
		}
	case reflect.Struct:
		e := e.(ast.RecordLit)
		structType := v.Type()
		for i := 0; i < structType.NumField(); i++ {
			// FIXME ignores fields in RecordLit not in Struct
			unmarshal(e[structType.Field(i).Name], v.Field(i))
		}
	case reflect.Func:
		e := e.(*ast.LambdaExpr)
		fnType := v.Type()
		returnType := fnType.Out(0)
		fn := reflect.MakeFunc(fnType, dhallShim(returnType, e))
		v.Set(fn)
	case reflect.Slice:
		if _, ok := e.(ast.EmptyList); ok {
			v.Set(reflect.MakeSlice(v.Type(), 0, 0))
			return
		}
		e := e.(ast.NonEmptyList)
		slice := reflect.MakeSlice(v.Type(), len(e), len(e))
		for i, expr := range e {
			unmarshal(expr, slice.Index(i))
		}
		v.Set(slice)
	default:
		switch e := e.(type) {
		case ast.DoubleLit:
			v.SetFloat(float64(e))
		case ast.BoolLit:
			v.SetBool(bool(e))
		case ast.NaturalLit:
			v.SetInt(int64(e))
		case ast.IntegerLit:
			v.SetInt(int64(e))
		case ast.TextLit:
			// FIXME: ensure TextLit doesn't have interpolations
			v.SetString(e.Suffix)
		}
	}
}
