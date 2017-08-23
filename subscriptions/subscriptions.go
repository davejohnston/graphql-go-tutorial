// Package mutations provides all available mutations in the graphql API
package subscriptions


import (
	"github.com/graphql-go/graphql"
)

func GetRootFields() graphql.Fields {
	return graphql.Fields{
		"messageAdded": MessageAdded(),
	}
}
