package uuider

type UUIDer interface {
	Generate() (*string, error)
}
