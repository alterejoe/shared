package structs

import "fmt"

type User struct {
	ID    string
	Name  string
	Email string
}

type Common struct {
	ID    string
	Name  string
	Value string
}

type Method string

const (
	NONE   Method = ""
	GET    Method = "get"
	POST   Method = "post"
	PUT    Method = "put"
	DELETE Method = "delete"
)

type Hx struct {
	Method    Method
	Target    string
	Include   string
	Trigger   string
	Swap      string
	Indicator string
	Vals      string

	Confirm string
	URL     string
}

func (h Hx) Attr() string {

	switch h.Method {
	case GET:
		return fmt.Sprintf("hx-get=%s", h.URL)

	case POST:
		return fmt.Sprintf("hx-post=%s", h.URL)
	case PUT:
		return fmt.Sprintf("hx-put=%s", h.URL)
	case DELETE:
		return fmt.Sprintf("hx-delete=%s", h.URL)
	default:
		return ""
	}
}

type Link struct {
	Href string
}
