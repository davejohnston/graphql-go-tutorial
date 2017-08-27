// Package mutations provides all available mutations in the graphql API
package mutations

import (
	"github.com/graphql-go/graphql"
)

// GetRootFields returns all the available graphql mutations.  This should be
// used when constructing a schema
func GetRootFields() graphql.Fields {
	return graphql.Fields{
		"addChannel": AddChannel(),
		"addMessage": AddMessage(),
	}
}
