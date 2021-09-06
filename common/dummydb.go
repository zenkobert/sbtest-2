package common

type DummyDB interface {
	Log(record string) error
}
