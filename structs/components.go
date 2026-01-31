package structs

import "strconv"

// ============================================================================
// COMPONENT STRUCTS
// ============================================================================

type Notice struct {
	Common
	AlignText string
}

func (n Notice) Align() bool {
	switch n.AlignText {
	case "left":
		return false
	default:
		return true
	}
}

type Radio struct {
	Common
	Hx
	FormBehaviors
	Checked  bool
	Required bool
}

type Checkbox struct {
	Common
	Hx
	FormBehaviors
	Checked  bool
	Required bool
}

type Input struct {
	Common
	Hx
	FormBehaviors
	Placeholder string
	InputType   string
	Disabled    bool
	Autofocus   bool
	Required    bool
}

func (i Input) Type() string {
	switch i.InputType {
	case "password":
		return "password"
	case "email":
		return "email"
	case "number":
		return "number"
	case "date":
		return "date"
	case "time":
		return "time"
	case "datetime-local":
		return "datetime-local"
	default:
		return "text"
	}
}

type Textarea struct {
	Common
	Hx
	FormBehaviors
	Placeholder string
	Rows        int
	Disabled    bool
	Autofocus   bool
	Required    bool
}

func (t Textarea) RowsString() string {
	if t.Rows <= 0 {
		return "3"
	}
	return strconv.Itoa(t.Rows)
}

type Select struct {
	Common
	Hx
	FormBehaviors
	Disabled bool
	Required bool
}

type Link struct {
	Common
	Hx
	Href string
}

type Form struct {
	Common
	Hx
	FormBehaviors
	DirtyScope bool // Enable dirty tracking scope for this form
}
