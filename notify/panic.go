package notify

var Panic = RuntimePanic

func RuntimePanic(err error) error {
	panic(err)
}
