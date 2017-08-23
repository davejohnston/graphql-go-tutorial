package mutations

import (
	"example.com/graphql/types"
	"github.com/graphql-go/graphql"
	"log"
)

func AddChannel() *graphql.Field {
	return &graphql.Field{
		Type:        types.ChannelType, // the return type for this field
		Description: "TODO",
		Args: graphql.FieldConfigArgument{
			"name": &graphql.ArgumentConfig{Type: graphql.String},
		},
		Resolve: addChannel,
	}
}

func addChannel(params graphql.ResolveParams) (interface{}, error) {
	log.Printf("Processing GraphQL Mutation addChannel %v\n", params.Args)

	name := params.Args["name"].(string)
	id := len(types.ChannelList)
	id++

	channel := types.Channel{
		Id:   string(id),
		Name: name,
	}
	types.ChannelList = append(types.ChannelList, &channel)

	return channel, nil
}
