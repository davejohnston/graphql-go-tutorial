package types

import (
	"github.com/graphql-go/graphql"
)

// ChannelList acts as a cheap database for us.
// When a new channel is added, its appended to the list
// When a mesasge is submitted it is added to the Channel.Messages field
var (
	ChannelList []*Channel
)

// Channel represents a chatroom chanel.  It holds a list of all messages
// submitted to that channel.
type Channel struct {
	ID       string    `json:"id"`
	Name     string    `json:"name"`
	Messages []Message `json:"messages"`
}

// ChannelType is a graphql output type.  This can be returned from
// graphql queries
var ChannelType = graphql.NewObject(graphql.ObjectConfig{
	Name: "Channel",
	Fields: graphql.Fields{
		"id":       &graphql.Field{Type: graphql.ID},
		"name":     &graphql.Field{Type: graphql.String},
		"messages": &graphql.Field{Type: graphql.NewList(MessageType)},
	},
})

// Init setups mock data
func init() {

	var soccer = Channel{
		ID:   "1",
		Name: "soccer",
		Messages: []Message{
			{
				ID:   "1",
				Text: "soccer is football",
			},
			{
				ID:   "2",
				Text: "hello soccer World Cup",
			},
		},
	}

	var baseball = Channel{
		ID:   "2",
		Name: "baseball",
		Messages: []Message{
			{
				ID:   "3",
				Text: "baseball is life",
			},
			{
				ID:   "4",
				Text: "hello baseball world series",
			},
		},
	}

	ChannelList = append(ChannelList, &soccer)
	ChannelList = append(ChannelList, &baseball)
}
