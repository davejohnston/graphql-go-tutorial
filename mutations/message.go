package mutations

import (
	"github.com/davejohnston/graphql-go-tutorial/types"
	"github.com/golang/glog"
	"github.com/graphql-go/graphql"
	"log"
	"strconv"
	"sync/atomic"
)

var messageID uint64 = 10

// AddMessage is a graphql query for adding new messages to a chatroom (channel)
func AddMessage() *graphql.Field {
	return &graphql.Field{
		Type:        types.MessageType, // the return type for this field
		Description: "TODO",
		Args: graphql.FieldConfigArgument{
			"message": &graphql.ArgumentConfig{Type: types.MessageInputType},
		},
		Resolve: addMessage,
	}
}

// executeCommand marshalls the graphql.go request to JSON, creates a Command struct, then assigns a UUID before
// submitting to the broker to be processed
func addMessage(params graphql.ResolveParams) (interface{}, error) {
	log.Printf("Processing GraphQL Mutation addMessage %v\n", params.Args)

	messageInput := params.Args["message"].(map[string]interface{})

	channelID := messageInput["channelId"].(string)
	text := messageInput["text"].(string)

	for _, channel := range types.ChannelList {
		if channel.ID == channelID {
			log.Printf("Found Message Channel [%s]:[%s]\n", channel.ID, channel.Name)

			// Get all the messages currently in the channel
			messages := channel.Messages
			// Generate Message ID
			atomic.AddUint64(&messageID, 1)

			glog.Infof("Creating Message with ID: %d for Channel", messageID)
			message := types.Message{
				ID:   strconv.FormatUint(messageID, 10),
				Text: text,
			}

			channel.Messages = append(messages, message)

			// Publish message for subscribers...

			return message, nil
		}
	}

	return nil, nil
}
