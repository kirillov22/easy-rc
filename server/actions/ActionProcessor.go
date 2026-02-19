package actions

func Process(p Processable, robot Robot) (R any, err error) {
	return p.Process(robot)
}
