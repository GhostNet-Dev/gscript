package object

type Environment struct {
	store        map[string]Object
	typeStore    map[string]*Environment
	outer        *Environment
	ProgramParam interface{}
}

func NewEnclosedEnvironment(outer *Environment) *Environment {
	env := NewEnvironment(outer.ProgramParam)
	env.outer = outer
	return env
}

func NewEnvironment(programParam interface{}) *Environment {
	s := make(map[string]Object)
	t := make(map[string]*Environment)
	return &Environment{store: s, typeStore: t, outer: nil, ProgramParam: programParam}
}

func (e *Environment) Get(name string) (Object, bool) {
	obj, ok := e.store[name]
	if !ok && e.outer != nil {
		obj, ok = e.outer.Get(name)
	}
	return obj, ok
}

func (e *Environment) Set(name string, val Object) Object {
	e.store[name] = val
	return val
}

func (e *Environment) TypeDefine(name string) *Environment {
	newEnv := NewEnvironment(e.ProgramParam)
	e.typeStore[name] = newEnv
	return newEnv
}

func (e *Environment) GetType(name string) (*Environment, bool) {
	obj, ok := e.typeStore[name]
	if !ok && e.outer != nil {
		obj, ok = e.outer.GetType(name)
	}
	return obj, ok
}
