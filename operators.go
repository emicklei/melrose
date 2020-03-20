package melrose

type Sequenceable interface {
	Storable
	S() Sequence
}

type Storable interface {
	Storex() string
}
