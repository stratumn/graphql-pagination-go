package starwars_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/graphql-go/graphql"
	"github.com/stratumn/graphql-pagination-go/examples/starwars"
)

func TestList_TestFetching_CorrectlyFetchesTheFirstShipOfTheRebels(t *testing.T) {
	query := `
        query RebelsShipsQuery {
          rebels {
            name,
            ships(first: 1) {
							items {
								name
							}
            }
          }
        }
      `
	expected := &graphql.Result{
		Data: map[string]interface{}{
			"rebels": map[string]interface{}{
				"name": "Alliance to Restore the Republic",
				"ships": map[string]interface{}{
					"items": []interface{}{
						map[string]interface{}{
							"name": "X-Wing",
						},
					},
				},
			},
		},
	}
	result := graphql.Do(graphql.Params{
		Schema:        starwars.Schema,
		RequestString: query,
	})
	assert.EqualValues(t, expected, result)
}
func TestList_TestFetching_CorrectlyFetchesTheFirstTwoShipsOfTheRebels(t *testing.T) {
	query := `
        query MoreRebelShipsQuery {
          rebels {
            name,
            ships(first: 2) {
              items {
                name
              }
            }
          }
        }
      `
	expected := &graphql.Result{
		Data: map[string]interface{}{
			"rebels": map[string]interface{}{
				"name": "Alliance to Restore the Republic",
				"ships": map[string]interface{}{
					"items": []interface{}{
						map[string]interface{}{
							"name": "X-Wing",
						},
						map[string]interface{}{
							"name": "Y-Wing",
						},
					},
				},
			},
		},
	}
	result := graphql.Do(graphql.Params{
		Schema:        starwars.Schema,
		RequestString: query,
	})
	assert.EqualValues(t, expected, result)
}
func TestList_TestFetching_CorrectlyFetchesTheNextThreeShipsOfTheRebelsWithACursor(t *testing.T) {
	query := `
        query EndOfRebelShipsQuery {
          rebels {
            name,
            ships(first: 3 after: "YXJyYXljb25uZWN0aW9uOjE=") {
              items {
                name
              }
            }
          }
        }
      `
	expected := &graphql.Result{
		Data: map[string]interface{}{
			"rebels": map[string]interface{}{
				"name": "Alliance to Restore the Republic",
				"ships": map[string]interface{}{
					"items": []interface{}{
						map[string]interface{}{
							"name": "A-Wing",
						},
						map[string]interface{}{
							"name": "Millenium Falcon",
						},
						map[string]interface{}{
							"name": "Home One",
						},
					},
				},
			},
		},
	}
	result := graphql.Do(graphql.Params{
		Schema:        starwars.Schema,
		RequestString: query,
	})
	assert.EqualValues(t, expected, result)
}
func TestList_TestFetching_CorrectlyFetchesNoShipsOfTheRebelsAtTheEndOfTheList(t *testing.T) {
	query := `
        query RebelsQuery {
          rebels {
            name,
            ships(first: 3 after: "YXJyYXljb25uZWN0aW9uOjQ=") {
              items {
                name
              }
            }
          }
        }
      `
	expected := &graphql.Result{
		Data: map[string]interface{}{
			"rebels": map[string]interface{}{
				"name": "Alliance to Restore the Republic",
				"ships": map[string]interface{}{
					"items": []interface{}{},
				},
			},
		},
	}
	result := graphql.Do(graphql.Params{
		Schema:        starwars.Schema,
		RequestString: query,
	})
	assert.EqualValues(t, expected, result)
}
func TestList_TestFetching_CorrectlyIdentifiesTheEndOfTheList(t *testing.T) {
	query := `
        query EndOfRebelShipsQuery {
          rebels {
            name,
            originalShips: ships(first: 2) {
              items {
                name
              }
              pageInfo {
                hasNextPage
              }
            }
            moreShips: ships(first: 3 after: "YXJyYXljb25uZWN0aW9uOjE=") {
              items {
                name
              }
              pageInfo {
                hasNextPage
              }
            }
          }
        }
      `
	expected := &graphql.Result{
		Data: map[string]interface{}{
			"rebels": map[string]interface{}{
				"name": "Alliance to Restore the Republic",
				"originalShips": map[string]interface{}{
					"items": []interface{}{
						map[string]interface{}{
							"name": "X-Wing",
						},
						map[string]interface{}{
							"name": "Y-Wing",
						},
					},
					"pageInfo": map[string]interface{}{
						"hasNextPage": true,
					},
				},
				"moreShips": map[string]interface{}{
					"items": []interface{}{
						map[string]interface{}{
							"name": "A-Wing",
						},
						map[string]interface{}{
							"name": "Millenium Falcon",
						},
						map[string]interface{}{
							"name": "Home One",
						},
					},
					"pageInfo": map[string]interface{}{
						"hasNextPage": false,
					},
				},
			},
		},
	}
	result := graphql.Do(graphql.Params{
		Schema:        starwars.Schema,
		RequestString: query,
	})
	assert.EqualValues(t, expected, result)
}
