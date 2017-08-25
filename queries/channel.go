package queries

import (
	"github.com/davejohnston/graphql-go-tutorial/types"
	"github.com/golang/glog"
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
	glog.Infof("Processing GraphQL Query Channels\n")
	return types.ChannelList, nil
}

func channel(params graphql.ResolveParams) (interface{}, error) {
	glog.Infof("Processing GraphQL Query Channel [%v]\n", params.Args)

	channelId := params.Args["id"].(string)

	for _, channel := range types.ChannelList {
		if channel.Id == channelId {
			glog.Infof("Found Channel [%s] ID:%s", channel.Name, channel.Id)
			return channel, nil
		}
	}
	log.Printf("Failed to find Channel with ID [%s]", channelId)
	return nil, nil
}
