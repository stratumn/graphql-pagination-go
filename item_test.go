package pagination_test

import (
	"errors"
	"fmt"
	"reflect"
	"testing"

	"context"

	"github.com/graphql-go/graphql"
	"github.com/graphql-go/graphql/gqlerrors"
	"github.com/graphql-go/graphql/language/location"
	"github.com/graphql-go/graphql/testutil"
	pagination "github.com/stratumn/graphql-pagination-go"
)

type user struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}
type photo struct {
	ID    int `json:"id"`
	Width int `json:"width"`
}

var ItemTestUserData = map[string]*user{
	"1": &user{1, "John Doe"},
	"2": &user{2, "Jane Smith"},
}
var ItemTestPhotoData = map[string]*photo{
	"3": &photo{3, 300},
	"4": &photo{4, 400},
}

// declare types first, define later in init()
// because they all depend on ItemTestDef
var ItemTestUserType *graphql.Object
var ItemTestPhotoType *graphql.Object

var ItemTestDef = pagination.NewItemDefinitions(pagination.ItemDefinitionsConfig{
	IDFetcher: func(id string, info graphql.ResolveInfo, ctx context.Context) (interface{}, error) {
		if user, ok := ItemTestUserData[id]; ok {
			return user, nil
		}
		if photo, ok := ItemTestPhotoData[id]; ok {
			return photo, nil
		}
		return nil, errors.New("Unknown Item")
	},
	TypeResolve: func(p graphql.ResolveTypeParams) *graphql.Object {
		switch p.Value.(type) {
		case *user:
			return ItemTestUserType
		case *photo:
			return ItemTestPhotoType
		default:
			panic(fmt.Sprintf("Unknown object type `%v`", p.Value))
		}
	},
})
var ItemTestQueryType = graphql.NewObject(graphql.ObjectConfig{
	Name: "Query",
	Fields: graphql.Fields{
		"Item": ItemTestDef.ItemField,
	},
})

// be careful not to define schema here, since ItemTestUserType and ItemTestPhotoType wouldn't be defined till init()
var ItemTestSchema graphql.Schema

func init() {
	ItemTestUserType = graphql.NewObject(graphql.ObjectConfig{
		Name: "User",
		Fields: graphql.Fields{
			"id": &graphql.Field{
				Type: graphql.NewNonNull(graphql.ID),
			},
			"name": &graphql.Field{
				Type: graphql.String,
			},
		},
		Interfaces: []*graphql.Interface{ItemTestDef.ItemInterface},
	})
	ItemTestPhotoType = graphql.NewObject(graphql.ObjectConfig{
		Name: "Photo",
		Fields: graphql.Fields{
			"id": &graphql.Field{
				Type: graphql.NewNonNull(graphql.ID),
			},
			"width": &graphql.Field{
				Type: graphql.Int,
			},
		},
		Interfaces: []*graphql.Interface{ItemTestDef.ItemInterface},
	})

	ItemTestSchema, _ = graphql.NewSchema(graphql.SchemaConfig{
		Query: ItemTestQueryType,
		Types: []graphql.Type{ItemTestUserType, ItemTestPhotoType},
	})
}
func TestItemInterfaceAndFields_AllowsRefetching_GetsTheCorrectIDForUsers(t *testing.T) {
	query := `{
        Item(id: "1") {
          id
        }
      }`
	expected := &graphql.Result{
		Data: map[string]interface{}{
			"Item": map[string]interface{}{
				"id": "1",
			},
		},
	}
	result := graphql.Do(graphql.Params{
		Schema:        ItemTestSchema,
		RequestString: query,
	})
	if !reflect.DeepEqual(result, expected) {
		t.Fatalf("wrong result, graphql result diff: %v", testutil.Diff(expected, result))
	}
}
func TestItemInterfaceAndFields_AllowsRefetching_GetsTheCorrectIDForPhotos(t *testing.T) {
	query := `{
        Item(id: "4") {
          id
        }
      }`
	expected := &graphql.Result{
		Data: map[string]interface{}{
			"Item": map[string]interface{}{
				"id": "4",
			},
		},
	}
	result := graphql.Do(graphql.Params{
		Schema:        ItemTestSchema,
		RequestString: query,
	})
	if !reflect.DeepEqual(result, expected) {
		t.Fatalf("wrong result, graphql result diff: %v", testutil.Diff(expected, result))
	}
}
func TestItemInterfaceAndFields_AllowsRefetching_GetsTheCorrectNameForUsers(t *testing.T) {
	query := `{
        Item(id: "1") {
          id
          ... on User {
            name
          }
        }
      }`
	expected := &graphql.Result{
		Data: map[string]interface{}{
			"Item": map[string]interface{}{
				"id":   "1",
				"name": "John Doe",
			},
		},
	}
	result := graphql.Do(graphql.Params{
		Schema:        ItemTestSchema,
		RequestString: query,
	})
	if !reflect.DeepEqual(result, expected) {
		t.Fatalf("wrong result, graphql result diff: %v", testutil.Diff(expected, result))
	}
}
func TestItemInterfaceAndFields_AllowsRefetching_GetsTheCorrectWidthForPhotos(t *testing.T) {
	query := `{
        Item(id: "4") {
          id
          ... on Photo {
            width
          }
        }
      }`
	expected := &graphql.Result{
		Data: map[string]interface{}{
			"Item": map[string]interface{}{
				"id":    "4",
				"width": 400,
			},
		},
	}
	result := graphql.Do(graphql.Params{
		Schema:        ItemTestSchema,
		RequestString: query,
	})
	if !reflect.DeepEqual(result, expected) {
		t.Fatalf("wrong result, graphql result diff: %v", testutil.Diff(expected, result))
	}
}
func TestItemInterfaceAndFields_AllowsRefetching_GetsTheCorrectTypeNameForUsers(t *testing.T) {
	query := `{
        Item(id: "1") {
          id
          __typename
        }
      }`
	expected := &graphql.Result{
		Data: map[string]interface{}{
			"Item": map[string]interface{}{
				"id":         "1",
				"__typename": "User",
			},
		},
	}
	result := graphql.Do(graphql.Params{
		Schema:        ItemTestSchema,
		RequestString: query,
	})
	if !reflect.DeepEqual(result, expected) {
		t.Fatalf("wrong result, graphql result diff: %v", testutil.Diff(expected, result))
	}
}
func TestItemInterfaceAndFields_AllowsRefetching_GetsTheCorrectTypeNameForPhotos(t *testing.T) {
	query := `{
        Item(id: "4") {
          id
          __typename
        }
      }`
	expected := &graphql.Result{
		Data: map[string]interface{}{
			"Item": map[string]interface{}{
				"id":         "4",
				"__typename": "Photo",
			},
		},
	}
	result := graphql.Do(graphql.Params{
		Schema:        ItemTestSchema,
		RequestString: query,
	})
	if !reflect.DeepEqual(result, expected) {
		t.Fatalf("wrong result, graphql result diff: %v", testutil.Diff(expected, result))
	}
}
func TestItemInterfaceAndFields_AllowsRefetching_IgnoresPhotoFragmentsOnUser(t *testing.T) {
	query := `{
        Item(id: "1") {
          id
          ... on Photo {
            width
          }
        }
      }`
	expected := &graphql.Result{
		Data: map[string]interface{}{
			"Item": map[string]interface{}{
				"id": "1",
			},
		},
	}
	result := graphql.Do(graphql.Params{
		Schema:        ItemTestSchema,
		RequestString: query,
	})
	if !reflect.DeepEqual(result, expected) {
		t.Fatalf("wrong result, graphql result diff: %v", testutil.Diff(expected, result))
	}
}
func TestItemInterfaceAndFields_AllowsRefetching_ReturnsNullForBadIDs(t *testing.T) {
	query := `{
        Item(id: "5") {
          id
        }
      }`
	expected := &graphql.Result{
		Data: map[string]interface{}{
			"Item": nil,
		},
		Errors: []gqlerrors.FormattedError{
			{
				Message:   "Unknown Item",
				Locations: []location.SourceLocation{},
			},
		},
	}
	result := graphql.Do(graphql.Params{
		Schema:        ItemTestSchema,
		RequestString: query,
	})

	if !reflect.DeepEqual(result, expected) {
		t.Fatalf("wrong result, graphql result diff: %v", testutil.Diff(expected, result))
	}
}
func TestItemInterfaceAndFields_CorrectlyIntrospects_HasCorrectItemInterface(t *testing.T) {
	query := `{
        __type(name: "Item") {
          name
          kind
          fields {
            name
            type {
              kind
              ofType {
                name
                kind
              }
            }
          }
        }
      }`
	expected := &graphql.Result{
		Data: map[string]interface{}{
			"__type": map[string]interface{}{
				"name": "Item",
				"kind": "INTERFACE",
				"fields": []interface{}{
					map[string]interface{}{
						"name": "id",
						"type": map[string]interface{}{
							"kind": "NON_NULL",
							"ofType": map[string]interface{}{
								"name": "ID",
								"kind": "SCALAR",
							},
						},
					},
				},
			},
		},
	}
	result := graphql.Do(graphql.Params{
		Schema:        ItemTestSchema,
		RequestString: query,
	})
	if !reflect.DeepEqual(result, expected) {
		t.Fatalf("wrong result, graphql result diff: %v", testutil.Diff(expected, result))
	}
}
func TestItemInterfaceAndFields_CorrectlyIntrospects_HasCorrectItemRootField(t *testing.T) {
	query := `{
        __schema {
          queryType {
            fields {
              name
              type {
                name
                kind
              }
              args {
                name
                type {
                  kind
                  ofType {
                    name
                    kind
                  }
                }
              }
            }
          }
        }
      }`
	expected := &graphql.Result{
		Data: map[string]interface{}{
			"__schema": map[string]interface{}{
				"queryType": map[string]interface{}{
					"fields": []interface{}{
						map[string]interface{}{
							"name": "Item",
							"type": map[string]interface{}{
								"name": "Item",
								"kind": "INTERFACE",
							},
							"args": []interface{}{
								map[string]interface{}{
									"name": "id",
									"type": map[string]interface{}{
										"kind": "NON_NULL",
										"ofType": map[string]interface{}{
											"name": "ID",
											"kind": "SCALAR",
										},
									},
								},
							},
						},
					},
				},
			},
		},
	}
	result := graphql.Do(graphql.Params{
		Schema:        ItemTestSchema,
		RequestString: query,
	})
	if !reflect.DeepEqual(result, expected) {
		t.Fatalf("wrong result, graphql result diff: %v", testutil.Diff(expected, result))
	}
}
