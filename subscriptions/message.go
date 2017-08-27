// Package subscriptions provides all available subscription queries in the graphql API
package subscriptions

import (
	"github.com/davejohnston/graphql-go-tutorial/types"
	"github.com/golang/glog"
	"github.com/graphql-go/graphql"
	"github.com/graphql-go/graphql/gqlerrors"
)

// MessageAdded defines the graphql subcription, that should handle
// messageAdded events.
func MessageAdded() *graphql.Field {
	return &graphql.Field{
		Type:        types.MessageType, // the return type for this field
		Description: "TODO",
		Args: graphql.FieldConfigArgument{
			"channelId": &graphql.ArgumentConfig{Type: graphql.ID},
		},
		Resolve: messageAdded,
	}
}

// messageAdded uses the provided RootValue to lookup the message from
// the database.
// When the subscription is first created, we subscribe to that queue.
// (this is done in the handler)
// We setup a callback, so that everytime a message is published to the
// queue, this subscription query is trigger, followed by the callback which
// is responsible for writing the response back to the client
//
func messageAdded(params graphql.ResolveParams) (interface{}, error) {
	glog.Infof("Processing GraphQL Subscription messageAdded %v\n", params.Args)

	channelID := params.Args["channelId"]

	test := params.Info.RootValue.(map[string]interface{})
	payload := test["addMessage"].(map[string]interface{})

	channels := types.ChannelList
	for _, channel := range channels {
		if channel.ID == channelID {
			messages := channel.Messages
			for _, message := range messages {
				if message.ID == payload["id"].(string) {
					return message, nil
				}
			}
		}

	}
	return nil, gqlerrors.NewFormattedError("message for channel not found")
}
