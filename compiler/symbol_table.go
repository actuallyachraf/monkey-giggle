package compiler

// SymbolScope represents the scope of a symbol
type SymbolScope string

const (
	// GlobalScope marks symbols that are global (accessible from anywhere)
	GlobalScope SymbolScope = "GLOBAL"
	// LocalScope marks symbols that are local to the current frame
	LocalScope SymbolScope = "LOCAL"
	// BuiltinScope marks symbols (functions) that are part of language
	BuiltinScope SymbolScope = "BUILTIN"
	// FreeScope marks free variables
	FreeScope SymbolScope = "FREE"
)

// Symbol represents an identifier, its scope and index in the table
type Symbol struct {
	Name  string
	Scope SymbolScope
	Index int
}

// SymbolTable is an associative map that maps identifiers to symbols
type SymbolTable struct {
	Outer          *SymbolTable
	FreeSyms       []Symbol
	store          map[string]Symbol
	numDefinitions int
}

// NewSymbolTable creates a new symbol table instance.
func NewSymbolTable() *SymbolTable {
	return &SymbolTable{
		FreeSyms:       []Symbol{},
		store:          make(map[string]Symbol),
		numDefinitions: 0,
	}
}

// NewEnclosedSymbolTable creates a new symbol table instance with outer
func NewEnclosedSymbolTable(outer *SymbolTable) *SymbolTable {
	s := NewSymbolTable()
	s.Outer = outer
	return s
}

// Define a new symbole for a given identifier
func (s *SymbolTable) Define(name string) Symbol {
	sym := Symbol{Name: name, Index: s.numDefinitions}

	if s.Outer == nil {
		sym.Scope = GlobalScope
	} else {
		sym.Scope = LocalScope
	}

	s.store[name] = sym
	s.numDefinitions++

	return sym
}

// Resolve a symbol by it's name
func (s *SymbolTable) Resolve(name string) (Symbol, bool) {
	obj, ok := s.store[name]

	if !ok && s.Outer != nil {
		obj, ok := s.Outer.Resolve(name)
		if !ok {
			return obj, ok
		}
		if obj.Scope == GlobalScope || obj.Scope == BuiltinScope {
			return obj, ok
		}
		free := s.defineFree(obj)
		return free, true
	}

	return obj, ok
}

// DefineBuiltIn loads the builtin function into the symbol table
func (s *SymbolTable) DefineBuiltIn(index int, name string) Symbol {
	sym := Symbol{Name: name, Index: index, Scope: BuiltinScope}
	s.store[name] = sym
	return sym
}

// defineFree adds a symbol to the free var scope
func (s *SymbolTable) defineFree(orig Symbol) Symbol {
	s.FreeSyms = append(s.FreeSyms, orig)
	sym := Symbol{Name: orig.Name, Index: len(s.FreeSyms) - 1, Scope: FreeScope}
	s.store[orig.Name] = sym

	return sym
}
