package types

import (
	"github.com/graphql-go/graphql"
)

// Message defines a structure used to represent chatroom messages
type Message struct {
	ID   string `json:"id"`
	Text string `json:"text"`
}

// MessageInputType is a graphql input type.  This can be used in
// graphql queries, mutations and subscriptions to provide a message
// body.
var MessageInputType = graphql.NewInputObject(graphql.InputObjectConfig{
	Name: "MessageInput",
	Fields: graphql.InputObjectConfigFieldMap{
		"channelId": &graphql.InputObjectFieldConfig{Type: graphql.ID},
		"text":      &graphql.InputObjectFieldConfig{Type: graphql.String},
	},
})

// MessageType is a graphql output type.  This can be returned by graphql
// queries, mutations and subscriptions.
var MessageType = graphql.NewObject(graphql.ObjectConfig{
	Name: "Message",
	Fields: graphql.Fields{
		"id":   &graphql.Field{Type: graphql.ID},
		"text": &graphql.Field{Type: graphql.String},
	},
})
