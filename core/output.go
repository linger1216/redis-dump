package core

type output interface {
	save(commands []string) error
	close()
}
