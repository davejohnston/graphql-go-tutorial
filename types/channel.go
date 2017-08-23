package types

import (
	"github.com/graphql-go/graphql"
)

var (
	ChannelList []*Channel
)

type Channel struct {
	Id string `json:"id"`
	Name string `json:"name"`
	Messages []Message `json:"messages"`
}

var ChannelType = graphql.NewObject(graphql.ObjectConfig{
	Name: "Channel",
	Fields: graphql.Fields{
		"id": &graphql.Field{Type: graphql.ID},
		"name": &graphql.Field{Type: graphql.String},
		"messages": &graphql.Field{Type: graphql.NewList(MessageType)},
	},
})

// Init setups mock data
func init() {

	var soccer = Channel{
		Id:   "1",
		Name: "soccer",
		Messages: []Message{
			{
				Id:   "1",
				Text: "soccer is football",
			},
			{
				Id:   "2",
				Text: "hello soccer World Cup",
			},
		},
	}

	var baseball = Channel{
		Id:   "2",
		Name: "baseball",
		Messages: []Message{
			{
				Id:   "3",
				Text: "baseball is life",
			},
			{
				Id:   "4",
				Text: "hello baseball world series",
			},
		},
	}

	ChannelList = append(ChannelList, &soccer)
	ChannelList = append(ChannelList, &baseball)
}
