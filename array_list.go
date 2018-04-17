package relay

import (
	"encoding/base64"
	"errors"
	"fmt"
	"reflect"
	"strconv"
	"strings"
)

const PREFIX = "arraylist:"

type ArraySliceMetaInfo struct {
	SliceStart  int `json:"sliceStart"`
	ArrayLength int `json:"arrayLength"`
}

/*
A simple function that accepts an array and list arguments, and returns
a list object for use in GraphQL. It uses array offsets as pagination,
so pagination will only work if the array is static.
*/
func ListFromArray(data []interface{}, args ListArguments) *List {
	return ListFromArraySlice(
		data,
		args,
		ArraySliceMetaInfo{
			SliceStart:  0,
			ArrayLength: len(data),
		},
	)
}

/*
Given a slice (subset) of an array, returns a list object for use in
GraphQL.

This function is similar to `ListFromArray`, but is intended for use
cases where you know the cardinality of the list, consider it too large
to materialize the entire array, and instead wish pass in a slice of the
total result large enough to cover the range specified in `args`.
*/
func ListFromArraySlice(
	arraySlice []interface{},
	args ListArguments,
	meta ArraySliceMetaInfo,
) *List {
	sliceEnd := meta.SliceStart + len(arraySlice)
	beforeOffset := GetOffsetWithDefault(args.Before, meta.ArrayLength)
	afterOffset := GetOffsetWithDefault(args.After, -1)

	startOffset := ternaryMax(meta.SliceStart-1, afterOffset, -1) + 1
	endOffset := ternaryMin(sliceEnd, beforeOffset, meta.ArrayLength)

	if args.First != -1 {
		endOffset = min(endOffset, startOffset+args.First)
	}

	if args.Last != -1 {
		startOffset = max(startOffset, endOffset-args.Last)
	}

	begin := max(startOffset-meta.SliceStart, 0)
	end := len(arraySlice) - (sliceEnd - endOffset)

	if begin > end {
		return NewList()
	}

	slice := arraySlice[begin:end]

	edges := []*Edge{}
	for index, value := range slice {
		edges = append(edges, &Edge{
			Cursor: OffsetToCursor(startOffset + index),
			Node:   value,
		})
	}

	var firstEdgeCursor, lastEdgeCursor ListCursor
	if len(edges) > 0 {
		firstEdgeCursor = edges[0].Cursor
		lastEdgeCursor = edges[len(edges)-1:][0].Cursor
	}

	lowerBound := 0
	if len(args.After) > 0 {
		lowerBound = afterOffset + 1
	}

	upperBound := meta.ArrayLength
	if len(args.Before) > 0 {
		upperBound = beforeOffset
	}

	hasPreviousPage := false
	if args.Last != -1 {
		hasPreviousPage = startOffset > lowerBound
	}

	hasNextPage := false
	if args.First != -1 {
		hasNextPage = endOffset < upperBound
	}

	conn := NewList()
	conn.Edges = edges
	conn.PageInfo = PageInfo{
		StartCursor:     firstEdgeCursor,
		EndCursor:       lastEdgeCursor,
		HasPreviousPage: hasPreviousPage,
		HasNextPage:     hasNextPage,
	}

	return conn
}

// Creates the cursor string from an offset
func OffsetToCursor(offset int) ListCursor {
	str := fmt.Sprintf("%v%v", PREFIX, offset)
	return ListCursor(base64.StdEncoding.EncodeToString([]byte(str)))
}

// Re-derives the offset from the cursor string.
func CursorToOffset(cursor ListCursor) (int, error) {
	str := ""
	b, err := base64.StdEncoding.DecodeString(string(cursor))
	if err == nil {
		str = string(b)
	}
	str = strings.Replace(str, PREFIX, "", -1)
	offset, err := strconv.Atoi(str)
	if err != nil {
		return 0, errors.New("Invalid cursor")
	}
	return offset, nil
}

// Return the cursor associated with an object in an array.
func CursorForObjectInList(data []interface{}, object interface{}) ListCursor {
	offset := -1
	for i, d := range data {
		// TODO: better object comparison
		if reflect.DeepEqual(d, object) {
			offset = i
			break
		}
	}
	if offset == -1 {
		return ""
	}
	return OffsetToCursor(offset)
}

func GetOffsetWithDefault(cursor ListCursor, defaultOffset int) int {
	if cursor == "" {
		return defaultOffset
	}
	offset, err := CursorToOffset(cursor)
	if err != nil {
		return defaultOffset
	}
	return offset
}

func max(a, b int) int {
	if a < b {
		return b
	}
	return a
}

func ternaryMax(a, b, c int) int {
	return max(max(a, b), c)
}

func min(a, b int) int {
	if a > b {
		return b
	}
	return a
}

func ternaryMin(a, b, c int) int {
	return min(min(a, b), c)
}
