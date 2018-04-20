package pagination

import (
	"github.com/graphql-go/graphql"
)

// ResolveSingleInputFn is ...
type ResolveSingleInputFn func(input interface{}) interface{}

// PluralIdentifyingRootFieldConfig is ...
type PluralIdentifyingRootFieldConfig struct {
	ArgName            string               `json:"argName"`
	InputType          graphql.Input        `json:"inputType"`
	OutputType         graphql.Output       `json:"outputType"`
	ResolveSingleInput ResolveSingleInputFn `json:"resolveSingleInput"`
	Description        string               `json:"description"`
}

// PluralIdentifyingRootField is ...
func PluralIdentifyingRootField(config PluralIdentifyingRootFieldConfig) *graphql.Field {
	inputArgs := graphql.FieldConfigArgument{}
	if config.ArgName != "" {
		inputArgs[config.ArgName] = &graphql.ArgumentConfig{
			Type: graphql.NewNonNull(graphql.NewList(graphql.NewNonNull(config.InputType))),
		}
	}

	return &graphql.Field{
		Description: config.Description,
		Type:        graphql.NewList(config.OutputType),
		Args:        inputArgs,
		Resolve: func(p graphql.ResolveParams) (interface{}, error) {
			inputs, ok := p.Args[config.ArgName]
			if !ok {
				return nil, nil
			}

			if config.ResolveSingleInput == nil {
				return nil, nil
			}
			switch inputs := inputs.(type) {
			case []interface{}:
				res := []interface{}{}
				for _, input := range inputs {
					r := config.ResolveSingleInput(input)
					res = append(res, r)
				}
				return res, nil
			}
			return nil, nil
		},
	}
}
