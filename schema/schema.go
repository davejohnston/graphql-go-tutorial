package schema

import (
	"github.com/graphql-go/graphql"
	"github.com/davejohnston/graphql-go-tutorial/mutations"
	"github.com/davejohnston/graphql-go-tutorial/queries"
	"log"
	"github.com/davejohnston/graphql-go-tutorial/subscriptions"
)

var (
	// DiscoverySchema should be used when handling graphql.go requests
	Schema graphql.Schema
)

func buildSchema() graphql.Schema {

	schemaConfig := graphql.SchemaConfig{
		Query: graphql.NewObject(graphql.ObjectConfig{
			Name:   "RootQuery",
			Fields: queries.GetRootFields(),
		}),
		Mutation: graphql.NewObject(graphql.ObjectConfig{
			Name:   "RootMutation",
			Fields: mutations.GetRootFields(),
		}),
		Subscription: graphql.NewObject(graphql.ObjectConfig{
			Name: 	"RootSubscription",
			Fields: subscriptions.GetRootFields(),
		}),
	}

	schema, err := graphql.NewSchema(schemaConfig)
	if err != nil {
		log.Fatalf("Failed to create new  discovery schema, error: %v", err)
	}

	return schema
}

func init() {
	Schema = buildSchema()
}
