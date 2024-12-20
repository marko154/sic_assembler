package symtable

const (
	UNKNOWN = -1
)

type SymTable struct {
	// maps lablel to address
	table map[string]int
}

func NewSymTable() *SymTable {
	return &SymTable{}
}

func (s *SymTable) Set(label string, address int) {
	s.table[label] = address
}

func (s *SymTable) Get(label string) (int, bool) {
	value, ok := s.table[label]
	return value, ok
}

func (s *SymTable) Has(label string) bool {
	_, ok := s.table[label]
	return ok
}
