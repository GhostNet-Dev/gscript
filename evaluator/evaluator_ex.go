package evaluator

import (
	"github.com/GhostNet-Dev/glambda/ast"
	"github.com/GhostNet-Dev/glambda/object"
)

func evalTypeExpression(node *ast.TypeStatement, env *object.Environment) object.Object {
	name := node.Name.Value

	switch node.Type.Value {
	case "struct":
		structEnv := env.TypeDefine(name)
		return &object.Struct{Name: name, Value: Eval(node.Body, structEnv)}
	}

	return nil
}

func evalTypeIdentifier(node *ast.TypeIdentifier, env *object.Environment) object.Object {
	val := Eval(node.Variable, env)
	env.Set(node.Value, val)

	return &object.Struct{}
}

func evalObjectBlockStatement(block *ast.ObjectBlockStatement, env *object.Environment) object.Object {
	var result object.Object

	for _, statement := range block.Statements {
		result = Eval(statement, env)
	}
	return result
}

func evalForExpression(ie *ast.ForExpression, env *object.Environment) object.Object {
	Eval(ie.Init, env)
	condition := Eval(ie.Condition, env)
	if isError(condition) {
		return condition
	}

	for isTruthy(condition) {
		Eval(ie.Consequence, env)
		Eval(ie.Increment, env)
		condition = Eval(ie.Condition, env)
	}
	return NULL
}
