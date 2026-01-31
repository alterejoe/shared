package structs

import (
	"strings"
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

type Common struct {
	ID       string
	Name     string
	Value    string
	Disabled bool
}

type Hx struct {
	Method    Method
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
