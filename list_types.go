package relay

import (
	"fmt"
)

// ListCursor is an opaque index in a list
type ListCursor string

// PageInfo is the information about the current pagination of the list
type PageInfo struct {
	StartCursor     ListCursor `json:"startCursor"`
	EndCursor       ListCursor `json:"endCursor"`
	HasPreviousPage bool       `json:"hasPreviousPage"`
	HasNextPage     bool       `json:"hasNextPage"`
}

// List contains items with meta information about the content
type List struct {
	Items      []interface{} `json:"items"`
	PageInfo   PageInfo      `json:"pageInfo"`
	TotalCount int           `json:"totalCount"`
}

// NewList is a list constructor
func NewList() *List {
	return &List{
		Items:    []interface{}{},
		PageInfo: PageInfo{},
	}
}

// ListArguments is pagination query arguments
// Use NewListArguments() to properly initialize default values
type ListArguments struct {
	Before ListCursor `json:"before"`
	After  ListCursor `json:"after"`
	First  int        `json:"first"` // -1 for undefined, 0 would return zero results
	Last   int        `json:"last"`  //  -1 for undefined, 0 would return zero results
}

// type ListArgumentsConfig struct {
// 	Before ListCursor `json:"before"`
// 	After  ListCursor `json:"after"`

// 	// use pointers for `First` and `Last` fields
// 	// so constructor would know when to use default values
// 	First *int `json:"first"`
// 	Last  *int `json:"last"`
// }

// NewListArguments is a list arguments constructor
func NewListArguments(filters map[string]interface{}) ListArguments {
	conn := ListArguments{
		First:  -1,
		Last:   -1,
		Before: "",
		After:  "",
	}
	if filters != nil {
		if first, ok := filters["first"]; ok {
			if first, ok := first.(int); ok {
				conn.First = first
			}
		}
		if last, ok := filters["last"]; ok {
			if last, ok := last.(int); ok {
				conn.Last = last
			}
		}
		if before, ok := filters["before"]; ok {
			conn.Before = ListCursor(fmt.Sprintf("%v", before))
		}
		if after, ok := filters["after"]; ok {
			conn.After = ListCursor(fmt.Sprintf("%v", after))
		}
	}
	return conn
}
