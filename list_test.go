package pagination_test

import (
	"reflect"
	"testing"

	"github.com/graphql-go/graphql"
	"github.com/graphql-go/graphql/testutil"
	pagination "github.com/stratumn/graphql-pagination-go"
)

var listTestAllUsers = []interface{}{
	&user{Name: "Dan"},
	&user{Name: "Nick"},
	&user{Name: "Lee"},
	&user{Name: "Joe"},
	&user{Name: "Tim"},
}
var listTestUserType *graphql.Object
var listTestQueryType *graphql.Object
var listTestSchema graphql.Schema
var listTestListDef *pagination.GraphQLListDefinitions

func init() {
	listTestUserType = graphql.NewObject(graphql.ObjectConfig{
		Name: "User",
		Fields: graphql.Fields{
			"name": &graphql.Field{
				Type: graphql.String,
			},
			// re-define `friends` field later because `listTestUserType` has `listTestListDef` has `listTestUserType` (cyclic-reference)
			"friends": &graphql.Field{},
		},
	})

	listTestListDef = pagination.ListDefinitions(pagination.ListConfig{
		Name:     "Friend",
		NodeType: listTestUserType,
		EdgeFields: graphql.Fields{
			"friendshipTime": &graphql.Field{
				Type: graphql.String,
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					return "Yesterday", nil
				},
			},
		},
		ListFields: graphql.Fields{
			"totalCount": &graphql.Field{
				Type: graphql.Int,
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					return len(listTestAllUsers), nil
				},
			},
		},
	})

	// define `friends` field here after getting list definition
	listTestUserType.AddFieldConfig("friends", &graphql.Field{
		Type: listTestListDef.ListType,
		Args: pagination.ListArgs,
		Resolve: func(p graphql.ResolveParams) (interface{}, error) {
			arg := pagination.NewListArguments(p.Args)
			res := pagination.ListFromArray(listTestAllUsers, arg)
			return res, nil
		},
	})

	listTestQueryType = graphql.NewObject(graphql.ObjectConfig{
		Name: "Query",
		Fields: graphql.Fields{
			"user": &graphql.Field{
				Type: listTestUserType,
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					return listTestAllUsers[0], nil
				},
			},
		},
	})
	var err error
	listTestSchema, err = graphql.NewSchema(graphql.SchemaConfig{
		Query: listTestQueryType,
	})
	if err != nil {
		panic(err)
	}

}

func TestListDefinition_IncludesListAndEdgeFields(t *testing.T) {
	query := `
      query FriendsQuery {
        user {
          friends(first: 2) {
            totalCount
            edges {
              friendshipTime
              node {
                name
              }
            }
          }
        }
      }
    `
	expected := &graphql.Result{
		Data: map[string]interface{}{
			"user": map[string]interface{}{
				"friends": map[string]interface{}{
					"totalCount": 5,
					"edges": []interface{}{
						map[string]interface{}{
							"friendshipTime": "Yesterday",
							"node": map[string]interface{}{
								"name": "Dan",
							},
						},
						map[string]interface{}{
							"friendshipTime": "Yesterday",
							"node": map[string]interface{}{
								"name": "Nick",
							},
						},
					},
				},
			},
		},
	}
	result := graphql.Do(graphql.Params{
		Schema:        listTestSchema,
		RequestString: query,
	})
	if !reflect.DeepEqual(result, expected) {
		t.Fatalf("wrong result, graphql result diff: %v", testutil.Diff(expected, result))
	}
}
