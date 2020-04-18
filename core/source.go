package core

type source interface {
	has() bool
	next() ([][]string, error)
	close()
}
