// Package subscriptions provides all available subscription queries in the graphql API
package subscriptions

import (
	"github.com/graphql-go/graphql"
)

// GetRootFields returns all the subscription queries and should be
// used to build your query
func GetRootFields() graphql.Fields {
	return graphql.Fields{
		"messageAdded": MessageAdded(),
	}
}
