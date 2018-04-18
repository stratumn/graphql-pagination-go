package relay

import "github.com/graphql-go/graphql"

// ListArgs returns a GraphQLFieldConfigArgumentMap appropriate to include
// on a field whose return type is a list type.
var ListArgs = graphql.FieldConfigArgument{
	"before": &graphql.ArgumentConfig{
		Type: graphql.String,
	},
	"after": &graphql.ArgumentConfig{
		Type: graphql.String,
	},
	"first": &graphql.ArgumentConfig{
		Type: graphql.Int,
	},
	"last": &graphql.ArgumentConfig{
		Type: graphql.Int,
	},
}

// NewListArgs adds pagination arguments to configMap
func NewListArgs(configMap graphql.FieldConfigArgument) graphql.FieldConfigArgument {
	for fieldName, argConfig := range ListArgs {
		configMap[fieldName] = argConfig
	}
	return configMap
}

type ListConfig struct {
	Name       string          `json:"name"`
	NodeType   *graphql.Object `json:"nodeType"`
	ListFields graphql.Fields  `json:"listFields"`
}

type EdgeType struct {
	Node   interface{} `json:"node"`
	Cursor ListCursor  `json:"cursor"`
}
type GraphQLListDefinitions struct {
	ListType *graphql.Object `json:"listType"`
}

/*
The common page info type used by all lists.
*/
var pageInfoType = graphql.NewObject(graphql.ObjectConfig{
	Name:        "PageInfo",
	Description: "Information about pagination in a list.",
	Fields: graphql.Fields{
		"hasNextPage": &graphql.Field{
			Type:        graphql.NewNonNull(graphql.Boolean),
			Description: "When paginating forwards, are there more items?",
		},
		"hasPreviousPage": &graphql.Field{
			Type:        graphql.NewNonNull(graphql.Boolean),
			Description: "When paginating backwards, are there more items?",
		},
		"startCursor": &graphql.Field{
			Type:        graphql.String,
			Description: "When paginating backwards, the cursor to continue.",
		},
		"endCursor": &graphql.Field{
			Type:        graphql.String,
			Description: "When paginating forwards, the cursor to continue.",
		},
	},
})

// ListDefinitions returns a GraphQLObjectType for a list with the given name,
// and whose nodes are of the specified type.
func ListDefinitions(config ListConfig) *GraphQLListDefinitions {

	listType := graphql.NewObject(graphql.ObjectConfig{
		Name:        config.Name + "List",
		Description: "list of items.",

		Fields: graphql.Fields{
			"items": &graphql.Field{
				Type:        graphql.NewList(config.NodeType),
				Description: "Items of the list.",
			},
			"pageInfo": &graphql.Field{
				Type:        graphql.NewNonNull(pageInfoType),
				Description: "Information to aid in pagination.",
			},
			"totalCount": &graphql.Field{
				Type:        graphql.NewNonNull(graphql.Int),
				Description: "Count of all list items.",
			},
		},
	})
	for fieldName, fieldConfig := range config.ListFields {
		listType.AddFieldConfig(fieldName, fieldConfig)
	}

	return &GraphQLListDefinitions{
		ListType: listType,
	}
}
