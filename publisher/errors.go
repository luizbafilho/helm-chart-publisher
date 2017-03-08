package publisher

import "fmt"

type ResourceNotFoundErr string

func (e ResourceNotFoundErr) Error() string {
	return string(e)
}

type StorageErr struct {
	Err     error
	Message string
}

func (e StorageErr) Error() string {
	return fmt.Sprintf("%s: %s", e.Err, e.Message)
}

type HelmErr struct {
	Err     error
	Message string
}

func (e HelmErr) Error() string {
	return fmt.Sprintf("%s: %s", e.Err, e.Message)
}
