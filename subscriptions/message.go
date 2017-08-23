package subscriptions

import (
	"example.com/graphql/types"
	"github.com/golang/glog"
	"github.com/graphql-go/graphql"
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

	glog.Infof("YYYYY  Subscription messageAdded %v\n", params)

	return params.Info.RootValue, nil

}
