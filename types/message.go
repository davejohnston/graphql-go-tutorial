package types

import (
	"github.com/graphql-go/graphql"
)

type Message struct {
	Id string `json:"id"`
	Text string `json:"text"`
}

var MessageInputType = graphql.NewInputObject(graphql.InputObjectConfig{
	Name: "MessageInput",
	Fields: graphql.InputObjectConfigFieldMap{
	"channelId":       &graphql.InputObjectFieldConfig{Type: graphql.ID},
	"text": 		&graphql.InputObjectFieldConfig{Type: graphql.String},
	},
})

var MessageType = graphql.NewObject(graphql.ObjectConfig{
	Name: "Message",
	Fields: graphql.Fields{
		"id": &graphql.Field{Type: graphql.ID},
		"text": &graphql.Field{Type: graphql.String},
	},
})





