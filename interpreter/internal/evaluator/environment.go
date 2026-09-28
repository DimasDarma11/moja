package evaluator

// Environment handles variable storage and lexical scoping.
type Environment struct {
	store     map[string]Object
	iterators map[string]bool
	outer     *Environment
}

// NewEnvironment creates a new root (global) environment.
func NewEnvironment() *Environment {
	return &Environment{
		store:     make(map[string]Object),
		iterators: make(map[string]bool),
		outer:     nil,
	}
}

// NewEnclosedEnvironment creates a new child environment.
func NewEnclosedEnvironment(outer *Environment) *Environment {
	env := NewEnvironment()
	env.outer = outer
	return env
}

// Get retrieves a variable by walking up the environment chain.
func (e *Environment) Get(name string) (Object, bool) {
	obj, ok := e.store[name]
	if !ok && e.outer != nil {
		return e.outer.Get(name)
	}
	return obj, ok
}

// Set sets or updates a variable.
// Follows the "walk up scope chain" rule:
// If variable exists in current or any outer scope, update it there.
// If it does not exist anywhere, create it in the CURRENT scope.
// If the variable is an iterator in any scope, returns a RuntimeError.
func (e *Environment) Set(name string, val Object, line int, col int) Object {
	// Check if this variable is marked as iterator anywhere in scope chain
	if e.isIterator(name) {
		return &RuntimeError{
			Line:    line,
			Column:  col,
			Message: "Variabel iterator '" + name + "' tidak dapat diubah di dalam loop",
		}
	}

	// Walk up chain to find where it is defined
	target := e.findEnvDefining(name)
	if target != nil {
		target.store[name] = val
	} else {
		// Define in current scope
		e.store[name] = val
	}
	return val
}

// DefineIterator defines an iterator variable in the current scope.
func (e *Environment) DefineIterator(name string, val Object) {
	e.store[name] = val
	e.iterators[name] = true
}

func (e *Environment) isIterator(name string) bool {
	if e.iterators[name] {
		return true
	}
	if e.outer != nil {
		return e.outer.isIterator(name)
	}
	return false
}

func (e *Environment) findEnvDefining(name string) *Environment {
	if _, ok := e.store[name]; ok {
		return e
	}
	if e.outer != nil {
		return e.outer.findEnvDefining(name)
	}
	return nil
}
