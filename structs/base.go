package structs

import (
	"strings"
	"time"
)

// ============================================================================
// BASE STRUCTS
// ============================================================================

type Method int

const (
	GET Method = iota
	POST
	PUT
	DELETE
)

type User struct {
	ID        string
	Email     string
	Name      string
	CreatedAt time.Time
	// Add any other user fields you need
}

type Common struct {
	ID       string
	Name     string
	Class    string
	Value    string
	Disabled bool
}

type Hx struct {
	Method    Method
	Params    string
	URL       string
	Target    string
	Include   string
	Trigger   string
	Swap      string
	Indicator string
	Vals      string
	Confirm   string
	PushURL   string
	Boost     bool

	OnAfterRequest  string // hx-on::after-request - fires after request completes
	OnBeforeRequest string // hx-on::before-request - fires before request
	OnAfterSettle   string // hx-on::after-settle - fires after DOM settles

	Select    string // hx-select - select subset of response to swap
	SelectOOB string // hx-select-oob - select content for out-of-band swaps
	SwapOOB   string // hx-swap-oob - out-of-band swap
	Preserve  bool   // hx-preserve - preserve element during swap
	Sync      string // hx-sync - coordinate requests
	Disabled  bool   // hx-disabled-elt - disable elements during request
	Encoding  string // hx-encoding - for file uploads (multipart/form-data)
}

// FormBehaviors holds data attributes for client-side form validation and behavior
type FormBehaviors struct {
	Constraint     []string // e.g. ["required", "email", "minLength:8"]
	DirtyWatch     bool     // Enable dirty tracking for this element
	DirtyGroup     string   // Group ID for dirty tracking
	EnableOnValid  bool     // For checkboxes: enable target when checked and form valid
	EnableTarget   string   // Selector for element to enable/disable
	ConstraintForm bool     // Mark this form for constraint validation
}

func (fb FormBehaviors) ConstraintString() string {
	if len(fb.Constraint) == 0 {
		return ""
	}
	return strings.Join(fb.Constraint, ", ")
}
