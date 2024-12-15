package instruction

type CodeEmitter interface {
	emitCode() []byte
}

type Statement struct {
	Label    string
	Mnemonic string
	Args     string
}

func NewStatement(label, mnemonic, args string) *Statement {
	return &Statement{label, mnemonic, args}
}
