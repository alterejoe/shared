package create

import "github.com/microcosm-cc/bluemonday"

func GetSanitizer() *bluemonday.Policy {
	return bluemonday.StrictPolicy()
}
