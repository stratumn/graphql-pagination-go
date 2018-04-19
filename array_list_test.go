package pagination_test

import (
	"testing"

	"github.com/stratumn/graphql-pagination-go"
	"github.com/stretchr/testify/assert"
)

var arrayListTestLetters = []interface{}{
	"A", "B", "C", "D", "E",
}

func TestListFromArray_HandlesBasicSlicing_ReturnsAllElementsWithoutFilters(t *testing.T) {
	args := pagination.NewListArguments(nil)

	expected := &pagination.List{
		Items: []interface{}{"A", "B", "C", "D", "E"},
		PageInfo: pagination.PageInfo{
			StartCursor:     "YXJyYXljb25uZWN0aW9uOjA=",
			EndCursor:       "YXJyYXljb25uZWN0aW9uOjQ=",
			HasPreviousPage: false,
			HasNextPage:     false,
		},
		TotalCount: 5,
	}

	result := pagination.ListFromArray(arrayListTestLetters, args)
	assert.EqualValues(t, expected, result)
}

func TestListFromArray_HandlesBasicSlicing_RespectsASmallerFirst(t *testing.T) {
	// Create list arguments from map[string]interface{},
	// which you usually get from types.GQLParams.Args
	filter := map[string]interface{}{
		"first": 2,
	}
	args := pagination.NewListArguments(filter)

	// Alternatively, you can create list arg the following way.
	// args := pagination.NewListArguments(filter)
	// args.First = 2

	expected := &pagination.List{
		Items: []interface{}{"A", "B"},
		PageInfo: pagination.PageInfo{
			StartCursor:     "YXJyYXljb25uZWN0aW9uOjA=",
			EndCursor:       "YXJyYXljb25uZWN0aW9uOjE=",
			HasPreviousPage: false,
			HasNextPage:     true,
		},
		TotalCount: 5,
	}

	result := pagination.ListFromArray(arrayListTestLetters, args)
	assert.EqualValues(t, expected, result)
}

func TestListFromArray_HandlesBasicSlicing_RespectsAnOverlyLargeFirst(t *testing.T) {

	filter := map[string]interface{}{
		"first": 10,
	}
	args := pagination.NewListArguments(filter)

	expected := &pagination.List{
		Items: []interface{}{"A", "B", "C", "D", "E"},
		PageInfo: pagination.PageInfo{
			StartCursor:     "YXJyYXljb25uZWN0aW9uOjA=",
			EndCursor:       "YXJyYXljb25uZWN0aW9uOjQ=",
			HasPreviousPage: false,
			HasNextPage:     false,
		},
		TotalCount: 5,
	}

	result := pagination.ListFromArray(arrayListTestLetters, args)
	assert.EqualValues(t, expected, result)
}

func TestListFromArray_HandlesBasicSlicing_RespectsASmallerLast(t *testing.T) {

	filter := map[string]interface{}{
		"last": 2,
	}
	args := pagination.NewListArguments(filter)

	expected := &pagination.List{
		Items: []interface{}{"D", "E"},
		PageInfo: pagination.PageInfo{
			StartCursor:     "YXJyYXljb25uZWN0aW9uOjM=",
			EndCursor:       "YXJyYXljb25uZWN0aW9uOjQ=",
			HasPreviousPage: true,
			HasNextPage:     false,
		},
		TotalCount: 5,
	}

	result := pagination.ListFromArray(arrayListTestLetters, args)
	assert.EqualValues(t, expected, result)
}

func TestListFromArray_HandlesBasicSlicing_RespectsAnOverlyLargeLast(t *testing.T) {

	filter := map[string]interface{}{
		"last": 10,
	}
	args := pagination.NewListArguments(filter)

	expected := &pagination.List{
		Items: []interface{}{"A", "B", "C", "D", "E"},
		PageInfo: pagination.PageInfo{
			StartCursor:     "YXJyYXljb25uZWN0aW9uOjA=",
			EndCursor:       "YXJyYXljb25uZWN0aW9uOjQ=",
			HasPreviousPage: false,
			HasNextPage:     false,
		},
		TotalCount: 5,
	}

	result := pagination.ListFromArray(arrayListTestLetters, args)
	assert.EqualValues(t, expected, result)
}

func TestListFromArray_HandlesPagination_RespectsFirstAndAfter(t *testing.T) {

	filter := map[string]interface{}{
		"first": 2,
		"after": "YXJyYXljb25uZWN0aW9uOjE=",
	}
	args := pagination.NewListArguments(filter)

	expected := &pagination.List{
		Items: []interface{}{"C", "D"},
		PageInfo: pagination.PageInfo{
			StartCursor:     "YXJyYXljb25uZWN0aW9uOjI=",
			EndCursor:       "YXJyYXljb25uZWN0aW9uOjM=",
			HasPreviousPage: false,
			HasNextPage:     true,
		},
		TotalCount: 5,
	}

	result := pagination.ListFromArray(arrayListTestLetters, args)
	assert.EqualValues(t, expected, result)
}

func TestListFromArray_HandlesPagination_RespectsFirstAndAfterWithLongFirst(t *testing.T) {

	filter := map[string]interface{}{
		"first": 10,
		"after": "YXJyYXljb25uZWN0aW9uOjE=",
	}
	args := pagination.NewListArguments(filter)

	expected := &pagination.List{
		Items: []interface{}{"C", "D", "E"},
		PageInfo: pagination.PageInfo{
			StartCursor:     "YXJyYXljb25uZWN0aW9uOjI=",
			EndCursor:       "YXJyYXljb25uZWN0aW9uOjQ=",
			HasPreviousPage: false,
			HasNextPage:     false,
		},
		TotalCount: 5,
	}

	result := pagination.ListFromArray(arrayListTestLetters, args)
	assert.EqualValues(t, expected, result)
}

func TestListFromArray_HandlesPagination_RespectsLastAndBefore(t *testing.T) {
	filter := map[string]interface{}{
		"last":   2,
		"before": "YXJyYXljb25uZWN0aW9uOjM=",
	}
	args := pagination.NewListArguments(filter)

	expected := &pagination.List{
		Items: []interface{}{"B", "C"},
		PageInfo: pagination.PageInfo{
			StartCursor:     "YXJyYXljb25uZWN0aW9uOjE=",
			EndCursor:       "YXJyYXljb25uZWN0aW9uOjI=",
			HasPreviousPage: true,
			HasNextPage:     false,
		},
		TotalCount: 5,
	}

	result := pagination.ListFromArray(arrayListTestLetters, args)
	assert.EqualValues(t, expected, result)
}

func TestListFromArray_HandlesPagination_RespectsLastAndBeforeWithLongLast(t *testing.T) {
	filter := map[string]interface{}{
		"last":   10,
		"before": "YXJyYXljb25uZWN0aW9uOjM=",
	}
	args := pagination.NewListArguments(filter)

	expected := &pagination.List{
		Items: []interface{}{"A", "B", "C"},
		PageInfo: pagination.PageInfo{
			StartCursor:     "YXJyYXljb25uZWN0aW9uOjA=",
			EndCursor:       "YXJyYXljb25uZWN0aW9uOjI=",
			HasPreviousPage: false,
			HasNextPage:     false,
		},
		TotalCount: 5,
	}

	result := pagination.ListFromArray(arrayListTestLetters, args)
	assert.EqualValues(t, expected, result)
}

func TestListFromArray_HandlesPagination_RespectsFirstAndAfterAndBefore_TooFew(t *testing.T) {
	filter := map[string]interface{}{
		"first":  2,
		"after":  "YXJyYXljb25uZWN0aW9uOjA=",
		"before": "YXJyYXljb25uZWN0aW9uOjQ=",
	}
	args := pagination.NewListArguments(filter)

	expected := &pagination.List{
		Items: []interface{}{"B", "C"},
		PageInfo: pagination.PageInfo{
			StartCursor:     "YXJyYXljb25uZWN0aW9uOjE=",
			EndCursor:       "YXJyYXljb25uZWN0aW9uOjI=",
			HasPreviousPage: false,
			HasNextPage:     true,
		},
		TotalCount: 5,
	}

	result := pagination.ListFromArray(arrayListTestLetters, args)
	assert.EqualValues(t, expected, result)
}

func TestListFromArray_HandlesPagination_RespectsFirstAndAfterAndBefore_TooMany(t *testing.T) {
	filter := map[string]interface{}{
		"first":  4,
		"after":  "YXJyYXljb25uZWN0aW9uOjA=",
		"before": "YXJyYXljb25uZWN0aW9uOjQ=",
	}
	args := pagination.NewListArguments(filter)

	expected := &pagination.List{
		Items: []interface{}{"B", "C", "D"},
		PageInfo: pagination.PageInfo{
			StartCursor:     "YXJyYXljb25uZWN0aW9uOjE=",
			EndCursor:       "YXJyYXljb25uZWN0aW9uOjM=",
			HasPreviousPage: false,
			HasNextPage:     false,
		},
		TotalCount: 5,
	}

	result := pagination.ListFromArray(arrayListTestLetters, args)
	assert.EqualValues(t, expected, result)
}

func TestListFromArray_HandlesPagination_RespectsFirstAndAfterAndBefore_ExactlyRight(t *testing.T) {
	filter := map[string]interface{}{
		"first":  3,
		"after":  "YXJyYXljb25uZWN0aW9uOjA=",
		"before": "YXJyYXljb25uZWN0aW9uOjQ=",
	}
	args := pagination.NewListArguments(filter)

	expected := &pagination.List{
		Items: []interface{}{"B", "C", "D"},
		PageInfo: pagination.PageInfo{
			StartCursor:     "YXJyYXljb25uZWN0aW9uOjE=",
			EndCursor:       "YXJyYXljb25uZWN0aW9uOjM=",
			HasPreviousPage: false,
			HasNextPage:     false,
		},
		TotalCount: 5,
	}

	result := pagination.ListFromArray(arrayListTestLetters, args)
	assert.EqualValues(t, expected, result)
}

func TestListFromArray_HandlesPagination_RespectsLastAndAfterAndBefore_TooFew(t *testing.T) {
	filter := map[string]interface{}{
		"last":   2,
		"after":  "YXJyYXljb25uZWN0aW9uOjA=",
		"before": "YXJyYXljb25uZWN0aW9uOjQ=",
	}
	args := pagination.NewListArguments(filter)

	expected := &pagination.List{
		Items: []interface{}{"C", "D"},
		PageInfo: pagination.PageInfo{
			StartCursor:     "YXJyYXljb25uZWN0aW9uOjI=",
			EndCursor:       "YXJyYXljb25uZWN0aW9uOjM=",
			HasPreviousPage: true,
			HasNextPage:     false,
		},
		TotalCount: 5,
	}

	result := pagination.ListFromArray(arrayListTestLetters, args)
	assert.EqualValues(t, expected, result)
}

func TestListFromArray_HandlesPagination_RespectsLasttAndAfterAndBefore_TooMany(t *testing.T) {
	filter := map[string]interface{}{
		"last":   4,
		"after":  "YXJyYXljb25uZWN0aW9uOjA=",
		"before": "YXJyYXljb25uZWN0aW9uOjQ=",
	}
	args := pagination.NewListArguments(filter)

	expected := &pagination.List{
		Items: []interface{}{"B", "C", "D"},
		PageInfo: pagination.PageInfo{
			StartCursor:     "YXJyYXljb25uZWN0aW9uOjE=",
			EndCursor:       "YXJyYXljb25uZWN0aW9uOjM=",
			HasPreviousPage: false,
			HasNextPage:     false,
		},
		TotalCount: 5,
	}

	result := pagination.ListFromArray(arrayListTestLetters, args)
	assert.EqualValues(t, expected, result)
}

func TestListFromArray_HandlesPagination_RespectsLastAndAfterAndBefore_ExactlyRight(t *testing.T) {
	filter := map[string]interface{}{
		"last":   3,
		"after":  "YXJyYXljb25uZWN0aW9uOjA=",
		"before": "YXJyYXljb25uZWN0aW9uOjQ=",
	}
	args := pagination.NewListArguments(filter)

	expected := &pagination.List{
		Items: []interface{}{"B", "C", "D"},
		PageInfo: pagination.PageInfo{
			StartCursor:     "YXJyYXljb25uZWN0aW9uOjE=",
			EndCursor:       "YXJyYXljb25uZWN0aW9uOjM=",
			HasPreviousPage: false,
			HasNextPage:     false,
		},
		TotalCount: 5,
	}

	result := pagination.ListFromArray(arrayListTestLetters, args)
	assert.EqualValues(t, expected, result)
}

func TestListFromArray_HandlesCursorItemCases_ReturnsNoElementsIfFirstIsZero(t *testing.T) {
	filter := map[string]interface{}{
		"first": 0,
	}
	args := pagination.NewListArguments(filter)

	expected := &pagination.List{
		Items: []interface{}{},
		PageInfo: pagination.PageInfo{
			HasPreviousPage: false,
			HasNextPage:     true,
		},
		TotalCount: 5,
	}

	result := pagination.ListFromArray(arrayListTestLetters, args)
	assert.EqualValues(t, expected, result)
}

func TestListFromArray_HandlesCursorItemCases_ReturnsAllElementsIfCursorsAreInvalid(t *testing.T) {
	filter := map[string]interface{}{
		"before": "invalid",
		"after":  "invalid",
	}
	args := pagination.NewListArguments(filter)

	expected := &pagination.List{
		Items: []interface{}{"A", "B", "C", "D", "E"},
		PageInfo: pagination.PageInfo{
			StartCursor:     "YXJyYXljb25uZWN0aW9uOjA=",
			EndCursor:       "YXJyYXljb25uZWN0aW9uOjQ=",
			HasPreviousPage: false,
			HasNextPage:     false,
		},
		TotalCount: 5,
	}

	result := pagination.ListFromArray(arrayListTestLetters, args)
	assert.EqualValues(t, expected, result)
}

func TestListFromArray_HandlesCursorItemCases_ReturnsAllElementsIfCursorsAreOnTheOutside(t *testing.T) {
	filter := map[string]interface{}{
		"before": "YXJyYXljb25uZWN0aW9uOjYK",     // ==> offset: int(6)
		"after":  "YXJyYXljb25uZWN0aW9uOi0xCg==", // ==> offset: int(-1)
	}
	args := pagination.NewListArguments(filter)

	expected := &pagination.List{
		Items: []interface{}{"A", "B", "C", "D", "E"},
		PageInfo: pagination.PageInfo{
			StartCursor:     "YXJyYXljb25uZWN0aW9uOjA=",
			EndCursor:       "YXJyYXljb25uZWN0aW9uOjQ=",
			HasPreviousPage: false,
			HasNextPage:     false,
		},
		TotalCount: 5,
	}

	result := pagination.ListFromArray(arrayListTestLetters, args)
	assert.EqualValues(t, expected, result)
}

func TestListFromArray_HandlesCursorItemCases_ReturnsNullIfCursorsIsConsecutive(t *testing.T) {
	filter := map[string]interface{}{
		"before": "YXJyYXljb25uZWN0aW9uOjM=", // ==> offset: int(3)
		"after":  "YXJyYXljb25uZWN0aW9uOjI=", // ==> offset: int(2)
	}
	args := pagination.NewListArguments(filter)

	expected := &pagination.List{
		Items:      []interface{}{},
		PageInfo:   pagination.PageInfo{},
		TotalCount: 5,
	}

	result := pagination.ListFromArray(arrayListTestLetters, args)
	assert.EqualValues(t, expected, result)
}

func TestListFromArray_HandlesCursorItemCases_ReturnsNoElementsIfCursorsCross(t *testing.T) {
	filter := map[string]interface{}{
		"before": "YXJyYXljb25uZWN0aW9uOjI=", // ==> offset: int(2)
		"after":  "YXJyYXljb25uZWN0aW9uOjQ=", // ==> offset: int(4)
	}
	args := pagination.NewListArguments(filter)

	expected := &pagination.List{
		Items:      []interface{}{},
		PageInfo:   pagination.PageInfo{},
		TotalCount: 0, // better to raise an error if possible
	}

	result := pagination.ListFromArray(arrayListTestLetters, args)
	assert.EqualValues(t, expected, result)
}

func TestListFromArray_CursorForObjectInList_ReturnsAnItemCursor_GivenAnArrayAndAMemberObject(t *testing.T) {
	letterBCursor := pagination.CursorForObjectInList(arrayListTestLetters, "B")
	expected := pagination.ListCursor("YXJyYXljb25uZWN0aW9uOjE=")
	assert.EqualValues(t, expected, letterBCursor)
}

func TestListFromArray_CursorForObjectInList_ReturnsEmptyCursor_GivenAnArrayAndANonMemberObject(t *testing.T) {
	letterFCursor := pagination.CursorForObjectInList(arrayListTestLetters, "F")
	assert.EqualValues(t, "", letterFCursor)
}

func TestListFromArraySlice_JustRightArraySlice(t *testing.T) {
	filter := map[string]interface{}{
		"first": 2,
		"after": "YXJyYXljb25uZWN0aW9uOjA=",
	}
	args := pagination.NewListArguments(filter)

	expected := &pagination.List{
		Items: []interface{}{"B", "C"},
		PageInfo: pagination.PageInfo{
			StartCursor:     "YXJyYXljb25uZWN0aW9uOjE=",
			EndCursor:       "YXJyYXljb25uZWN0aW9uOjI=",
			HasPreviousPage: false,
			HasNextPage:     true,
		},
		TotalCount: 2,
	}

	result := pagination.ListFromArraySlice(
		arrayListTestLetters[1:3],
		args,
		pagination.ArraySliceMetaInfo{
			SliceStart:  1,
			ArrayLength: 5,
		},
	)
	assert.EqualValues(t, expected, result)
}

func TestListFromArraySlice_OversizedSliceLeft(t *testing.T) {
	filter := map[string]interface{}{
		"first": 2,
		"after": "YXJyYXljb25uZWN0aW9uOjA=",
	}
	args := pagination.NewListArguments(filter)

	expected := &pagination.List{
		Items: []interface{}{"B", "C"},
		PageInfo: pagination.PageInfo{
			StartCursor:     "YXJyYXljb25uZWN0aW9uOjE=",
			EndCursor:       "YXJyYXljb25uZWN0aW9uOjI=",
			HasPreviousPage: false,
			HasNextPage:     true,
		},
		TotalCount: 3,
	}

	result := pagination.ListFromArraySlice(
		arrayListTestLetters[0:3],
		args,
		pagination.ArraySliceMetaInfo{
			SliceStart:  0,
			ArrayLength: 5,
		},
	)
	assert.EqualValues(t, expected, result)
}

func TestListFromArraySlice_OversizedSliceRight(t *testing.T) {
	filter := map[string]interface{}{
		"first": 1,
		"after": "YXJyYXljb25uZWN0aW9uOjE=",
	}
	args := pagination.NewListArguments(filter)

	expected := &pagination.List{
		Items: []interface{}{"C"},
		PageInfo: pagination.PageInfo{
			StartCursor:     "YXJyYXljb25uZWN0aW9uOjI=",
			EndCursor:       "YXJyYXljb25uZWN0aW9uOjI=",
			HasPreviousPage: false,
			HasNextPage:     true,
		},
		TotalCount: 2,
	}

	result := pagination.ListFromArraySlice(
		arrayListTestLetters[2:4],
		args,
		pagination.ArraySliceMetaInfo{
			SliceStart:  2,
			ArrayLength: 5,
		},
	)
	assert.EqualValues(t, expected, result)
}

func TestListFromArraySlice_OversizedSliceBoth(t *testing.T) {
	filter := map[string]interface{}{
		"first": 1,
		"after": "YXJyYXljb25uZWN0aW9uOjE=",
	}
	args := pagination.NewListArguments(filter)

	expected := &pagination.List{
		Items: []interface{}{"C"},
		PageInfo: pagination.PageInfo{
			StartCursor:     "YXJyYXljb25uZWN0aW9uOjI=",
			EndCursor:       "YXJyYXljb25uZWN0aW9uOjI=",
			HasPreviousPage: false,
			HasNextPage:     true,
		},
		TotalCount: 3,
	}

	result := pagination.ListFromArraySlice(
		arrayListTestLetters[1:4],
		args,
		pagination.ArraySliceMetaInfo{
			SliceStart:  1,
			ArrayLength: 5,
		},
	)
	assert.EqualValues(t, expected, result)
}

func TestListFromArraySlice_UndersizedSliceLeft(t *testing.T) {
	filter := map[string]interface{}{
		"first": 3,
		"after": "YXJyYXljb25uZWN0aW9uOjE=",
	}
	args := pagination.NewListArguments(filter)

	expected := &pagination.List{
		Items: []interface{}{"D", "E"},
		PageInfo: pagination.PageInfo{
			StartCursor:     "YXJyYXljb25uZWN0aW9uOjM=",
			EndCursor:       "YXJyYXljb25uZWN0aW9uOjQ=",
			HasPreviousPage: false,
			HasNextPage:     false,
		},
		TotalCount: 2,
	}

	result := pagination.ListFromArraySlice(
		arrayListTestLetters[3:5],
		args,
		pagination.ArraySliceMetaInfo{
			SliceStart:  3,
			ArrayLength: 5,
		},
	)
	assert.EqualValues(t, expected, result)
}

func TestListFromArraySlice_UndersizedSliceRight(t *testing.T) {
	filter := map[string]interface{}{
		"first": 3,
		"after": "YXJyYXljb25uZWN0aW9uOjE=",
	}
	args := pagination.NewListArguments(filter)

	expected := &pagination.List{
		Items: []interface{}{"C", "D"},
		PageInfo: pagination.PageInfo{
			StartCursor:     "YXJyYXljb25uZWN0aW9uOjI=",
			EndCursor:       "YXJyYXljb25uZWN0aW9uOjM=",
			HasPreviousPage: false,
			HasNextPage:     true,
		},
		TotalCount: 2,
	}

	result := pagination.ListFromArraySlice(
		arrayListTestLetters[2:4],
		args,
		pagination.ArraySliceMetaInfo{
			SliceStart:  2,
			ArrayLength: 5,
		},
	)
	assert.EqualValues(t, expected, result)
}

func TestListFromArraySlice_UndersizedSliceBoth(t *testing.T) {
	filter := map[string]interface{}{
		"first": 3,
		"after": "YXJyYXljb25uZWN0aW9uOjE=",
	}
	args := pagination.NewListArguments(filter)

	expected := &pagination.List{
		Items: []interface{}{"D"},
		PageInfo: pagination.PageInfo{
			StartCursor:     "YXJyYXljb25uZWN0aW9uOjM=",
			EndCursor:       "YXJyYXljb25uZWN0aW9uOjM=",
			HasPreviousPage: false,
			HasNextPage:     true,
		},
		TotalCount: 1,
	}

	result := pagination.ListFromArraySlice(
		arrayListTestLetters[3:4],
		args,
		pagination.ArraySliceMetaInfo{
			SliceStart:  3,
			ArrayLength: 5,
		},
	)
	assert.EqualValues(t, expected, result)
}
