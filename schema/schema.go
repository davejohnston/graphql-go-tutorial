package schema

import (
	"github.com/graphql-go/graphql"
	"example.com/graphql/mutations"
	"example.com/graphql/queries"
	"log"
	"example.com/graphql/subscriptions"
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
