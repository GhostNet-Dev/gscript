package evaluator

import (
	"github.com/GhostNet-Dev/glambda/object"
)

var builtins = map[string]*object.Builtin{
	"len":    object.GetBuiltinByName("len"),
	"first":  object.GetBuiltinByName("first"),
	"last":   object.GetBuiltinByName("last"),
	"rest":   object.GetBuiltinByName("rest"),
	"push":   object.GetBuiltinByName("push"),
	"puts":   object.GetBuiltinByName("puts"),
	"int":    object.GetBuiltinByName("int"),
	"string": object.GetBuiltinByName("string"),
}

func AddBuiltIn(name string, builtin *object.Builtin) {
	builtins[name] = builtin
}
