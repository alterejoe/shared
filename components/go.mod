module github.com/alterejoe/shared/components

go 1.24.7

require (
	github.com/a-h/templ v0.3.1001
	github.com/alterejoe/shared/structs v0.0.0-20260131031456-308870a41570
)

require github.com/google/go-cmp v0.7.0 // indirect

replace github.com/alterejoe/shared/structs => ../structs/
