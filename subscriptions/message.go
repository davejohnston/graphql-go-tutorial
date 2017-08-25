package subscriptions

import (
	"fmt"
	"github.com/davejohnston/graphql-go-tutorial/types"
	"github.com/golang/glog"
	"github.com/graphql-go/graphql"
	"github.com/graphql-go/graphql/gqlerrors"
)

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

// executeCommand marshalls the graphql.go request to JSON, creates a Command struct, then assigns a UUID before
// submitting to the broker to be processed
func messageAdded(params graphql.ResolveParams) (interface{}, error) {
	glog.Infof("Processing GraphQL Subscription messageAdded %v\n", params.Args)

	channelId := params.Args["channelId"]

	test := params.Info.RootValue.(map[string]interface{})
	payload := test["addMessage"].(map[string]interface{})

	channels := types.ChannelList
	for _, channel := range channels {
		fmt.Printf("Comparing Channel id [%s] with requested channel id [%s]\n", channel.Id, channelId)
		if channel.Id == channelId {
			fmt.Printf("\tGot channel match\n")
			messages := channel.Messages
			for _, message := range messages {
				fmt.Printf("\t\tComparing Message id [%s] with requested message 	 id [%s]\n", message.Id, payload["id"].(string))
				if message.Id == payload["id"].(string) {
					fmt.Printf("\t\t\tGot Message match\n")
					fmt.Println("Message contains: ", message.Text)

					return message, nil
				}
			}
		}

	}

	return nil, gqlerrors.NewFormattedError("message for channel not found")
}
