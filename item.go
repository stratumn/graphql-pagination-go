package pagination

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"strings"

	"context"

	"github.com/graphql-go/graphql"
)

// ItemDefinitions is the GraphQL object type for an item list
type ItemDefinitions struct {
	ItemInterface *graphql.Interface
	ItemField     *graphql.Field
}

// ItemDefinitionsConfig is the configuration object for item list
type ItemDefinitionsConfig struct {
	IDFetcher   IDFetcherFn
	TypeResolve graphql.ResolveTypeFn
}

// IDFetcherFn returns the the object from an id
type IDFetcherFn func(id string, info graphql.ResolveInfo, ctx context.Context) (interface{}, error)

// GlobalIDFetcherFn returns the id of an object
type GlobalIDFetcherFn func(obj interface{}, info graphql.ResolveInfo, ctx context.Context) (string, error)

// NewItemDefinitions constructs a `Item` interface that objects can implement,
// and a field config for a `Item` root field, given a function to map from an
// ID to an underlying object, and a function to map from an underlying object
// to the concrete GraphQLObjectType it corresponds to.
// If the typeResolver is omitted, object resolution on the interface will be
// handled with the `isTypeOf` method on object types, as with any GraphQL
// interface without a provided `resolveType` method.
func NewItemDefinitions(config ItemDefinitionsConfig) *ItemDefinitions {
	ItemInterface := graphql.NewInterface(graphql.InterfaceConfig{
		Name:        "Item",
		Description: "An object with an ID",
		Fields: graphql.Fields{
			"id": &graphql.Field{
				Type:        graphql.NewNonNull(graphql.ID),
				Description: "The id of the object",
			},
		},
		ResolveType: config.TypeResolve,
	})

	ItemField := &graphql.Field{
		Name:        "Item",
		Description: "Fetches an object given its ID",
		Type:        ItemInterface,
		Args: graphql.FieldConfigArgument{
			"id": &graphql.ArgumentConfig{
				Type:        graphql.NewNonNull(graphql.ID),
				Description: "The ID of an object",
			},
		},
		Resolve: func(p graphql.ResolveParams) (interface{}, error) {
			if config.IDFetcher == nil {
				return nil, nil
			}
			id := ""
			if iid, ok := p.Args["id"]; ok {
				id = fmt.Sprintf("%v", iid)
			}
			return config.IDFetcher(id, p.Info, p.Context)
		},
	}
	return &ItemDefinitions{
		ItemInterface: ItemInterface,
		ItemField:     ItemField,
	}
}

// ResolvedGlobalID is the type and id of an object
type ResolvedGlobalID struct {
	Type string `json:"type"`
	ID   string `json:"id"`
}

// ToGlobalID takes a type name and an ID specific to that type name, and returns a
// "global ID" that is unique among all types.
func ToGlobalID(ttype string, id string) string {
	str := ttype + ":" + id
	encStr := base64.StdEncoding.EncodeToString([]byte(str))
	return encStr
}

// FromGlobalID takes the "global ID" created by toGlobalID, and returns the type name and ID
// used to create it.
func FromGlobalID(globalID string) *ResolvedGlobalID {
	strID := ""
	b, err := base64.StdEncoding.DecodeString(globalID)
	if err == nil {
		strID = string(b)
	}
	tokens := strings.Split(strID, ":")
	if len(tokens) < 2 {
		return nil
	}
	return &ResolvedGlobalID{
		Type: tokens[0],
		ID:   tokens[1],
	}
}

// GlobalIDField creates the configuration for an id field on a Item, using `toGlobalId` to
// construct the ID from the provided typename. The type-specific ID is fetcher
// by calling idFetcher on the object, or if not provided, by accessing the `id`
// property on the object.
func GlobalIDField(typeName string, idFetcher GlobalIDFetcherFn) *graphql.Field {
	return &graphql.Field{
		Name:        "id",
		Description: "The ID of an object",
		Type:        graphql.NewNonNull(graphql.ID),
		Resolve: func(p graphql.ResolveParams) (interface{}, error) {
			id := ""
			if idFetcher != nil {
				fetched, err := idFetcher(p.Source, p.Info, p.Context)
				id = fmt.Sprintf("%v", fetched)
				if err != nil {
					return id, err
				}
			} else {
				// try to get from p.Source (data)
				var objMap interface{}
				b, _ := json.Marshal(p.Source)
				_ = json.Unmarshal(b, &objMap)
				switch obj := objMap.(type) {
				case map[string]interface{}:
					if iid, ok := obj["id"]; ok {
						id = fmt.Sprintf("%v", iid)
					}
				}
			}
			globalID := ToGlobalID(typeName, id)
			return globalID, nil
		},
	}
}
