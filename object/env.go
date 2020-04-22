package object

// Environment represents binding maps for let statements.
type Environment struct {
	store map[string]Object
	outer *Environment
}

// NewEnv creates a new environment instance
func NewEnv() *Environment {
	return &Environment{
		store: make(map[string]Object),
		outer: nil,
	}
}

// NewEnclosedEnvironment creates an environment that extends an outer one.
func NewEnclosedEnvironment(outerEnv *Environment) *Environment {
	return &Environment{
		store: make(map[string]Object),
		outer: outerEnv,
	}
}

// Get an object by it's binding identifier
func (env *Environment) Get(name string) (Object, bool) {
	obj, ok := env.store[name]
	if !ok && env.outer != nil {
		obj, ok = env.outer.Get(name)
	}
	return obj, ok
}

// Set an object by it's binding identifier
func (env *Environment) Set(name string, obj Object) Object {
	env.store[name] = obj
	return obj
}
