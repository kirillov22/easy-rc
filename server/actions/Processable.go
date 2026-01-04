package actions

type Processable interface {
	Process() (R any, err error)
	Debug()
}
