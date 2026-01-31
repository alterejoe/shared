package interfaces

// validator is kept separate as it will be integral for testing
// and it is one of the few dependencies that will not have a pointer
// reference. It will be created each time it is needed
type Validator interface {
	Valid() bool
	AddFieldError(field, message string)
	CheckField(ok bool, key, message string)
	GetError() error
}
