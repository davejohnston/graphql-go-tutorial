package queries

import (
	"example.com/graphql/types"
	"github.com/graphql-go/graphql"
	"log"
)



func Channels() *graphql.Field {
	return &graphql.Field{
		Type:        graphql.NewList(types.ChannelType),
		Description: "Return all channels",
		Resolve:     channels,
	}
}

func Channel() *graphql.Field {
	return &graphql.Field{
		Type:        types.ChannelType, // the return type for this field
		Description: "Execute a new command - this creates a new command record",
		Args: graphql.FieldConfigArgument{
			"id": &graphql.ArgumentConfig{Type: graphql.ID},
		},
		Resolve: channel,
	}
}

func channels(params graphql.ResolveParams) (interface{}, error) {
	log.Printf("Processing GraphQL Query Channels\n")
	return types.ChannelList, nil
}

func channel(params graphql.ResolveParams) (interface{}, error) {
	log.Printf("Processing GraphQL Query Channel [%v]\n", params.Args)

	channelId := params.Args["id"].(string)

	for index := range types.ChannelList {
		if types.ChannelList[index].Id == channelId {
			log.Printf("Found Channel [%v]", types.ChannelList[index])
			return types.ChannelList[index], nil
		}
	}
	log.Printf("Failed to find Channel with ID [%s]", channelId)
	return nil, nil
}

