package queries

import (
	"github.com/davejohnston/graphql-go-tutorial/types"
	"github.com/golang/glog"
	"github.com/graphql-go/graphql"
	"log"
)

// Channels query returns a list of all channels
func Channels() *graphql.Field {
	return &graphql.Field{
		Type:        graphql.NewList(types.ChannelType),
		Description: "Return all channels",
		Resolve:     channels,
	}
}

// Channel given a specific channel id returns a channel
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

	channelID := params.Args["id"].(string)

	for _, channel := range types.ChannelList {
		if channel.ID == channelID {
			glog.Infof("Found Channel [%s] ID:%s", channel.Name, channel.ID)
			return channel, nil
		}
	}
	log.Printf("Failed to find Channel with ID [%s]", channelID)
	return nil, nil
}
