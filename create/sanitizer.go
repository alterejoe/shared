package create

import "github.com/microcosm-cc/bluemonday"

func CreateSanitizer() *bluemonday.Policy {
	return bluemonday.StrictPolicy()
}
