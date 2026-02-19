package actions

type Processable interface {
	Process(robot Robot) (R any, err error)
	Debug()
}
