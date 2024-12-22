//go:generate go install github.com/dmarkham/enumer@latest
//go:generate enumer -type=TaskStatus
package main

import "strings"

type TaskStatus int

const (
	Pending TaskStatus = iota
	Doing
	Completed
	Cancelled
	Unknown
)

func (s TaskStatus) LowerString() string {
	return strings.ToLower(s.String())
}
