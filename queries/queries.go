// Package queries defines all queries available in the graphql API
package queries

import (
	"github.com/graphql-go/graphql"
)

// GetRootFields returns all the qraphql queries
func GetRootFields() graphql.Fields {
	return graphql.Fields{
		"channels": Channels(),
		"channel":  Channel(),
	}
}
