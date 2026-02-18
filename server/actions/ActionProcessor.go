package actions

func Process(p Processable) (R any, err error) {
	return p.Process()
}
