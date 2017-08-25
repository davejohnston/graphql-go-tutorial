package mutations

import (
	"github.com/graphql-go/graphql"
	"github.com/davejohnston/graphql-go-tutorial/types"
	"log"
	"github.com/golang/glog"
	"strconv"
	"sync/atomic"
)

var messageId uint64 = 10

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

	channelId := messageInput["channelId"].(string)
	text := messageInput["text"].(string)

	for _, channel := range types.ChannelList {
		if channel.Id == channelId {
			log.Printf("Found Message Channel [%s]:[%s]\n", channel.Id, channel.Name)

			// Get all the messages currently in the channel
			messages := channel.Messages
			// Generate Message ID
			atomic.AddUint64(&messageId, 1)

			glog.Infof("Creating Message with ID: %d for Channel", messageId, )
			message := types.Message {
				Id: strconv.FormatUint(messageId, 10),
				Text: text,
			}

			channel.Messages = append(messages, message)

			// Publish message for subscribers...

			return message, nil
		}
	}

	return nil, nil
}


