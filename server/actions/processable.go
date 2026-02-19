package actions

type Processable interface {
	Process(robot Robot) (Processable, error)
}
