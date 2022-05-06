module github.com/ariary/notionion

go 1.17

replace github.com/jomei/notionapi => ../notionapi

require (
	github.com/elazarl/goproxy v0.0.0-20220417044921-416226498f94
	github.com/jomei/notionapi v0.0.0-00010101000000-000000000000
)

require github.com/pkg/errors v0.9.1 // indirect
