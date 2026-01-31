package structs

type Radio struct {
	Checked bool
	Color   string
	Hx
	Common
}

type Button struct {
	Common
	Hx
	Color    string
	Size     string
	Radius   string
	Border   string
	Variant  string
	Shadow   string
	Disabled bool
}

type Checkbox struct {
	Common
	Hx
	Checked bool
	Color   string
	Size    string
	Radius  string
	Border  string
	Shadow  string
}
