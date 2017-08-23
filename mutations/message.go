package mutations

import (
	"github.com/graphql-go/graphql"
	"example.com/graphql/types"
	"log"
	"github.com/golang/glog"
	"strconv"
)



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

	message := params.Args["message"].(map[string]interface{})

	channelId := message["channelId"].(string)
	glog.Infof("Channel ID: %s", channelId)
	text := message["text"].(string)
	glog.Infof("Message: %s", text)

	for index := range types.ChannelList {
		if types.ChannelList[index].Id == channelId {
			log.Printf("Found Channel [%v]", types.ChannelList[index])

			channel := types.ChannelList[index]
			messages := channel.Messages

			messageId := len(messages)
			glog.Infof("Message Size: %d", messageId)
			messageId++

			glog.Infof("Creating Message with ID: %d", messageId)
			message := types.Message {
				Id: strconv.Itoa(messageId),
				Text: text,
			}

			channel.Messages = append(messages, message)

			// Publish message for subscribers...

			return message, nil
		}
	}

	return nil, nil
}


