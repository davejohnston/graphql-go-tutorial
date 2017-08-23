#curl -H "Content-Type: application/graphql" -XPOST http://localhost:9999/graphql -d '{HelloQuery}'
#curl http://localhost:9999/graphql -d '{"query": "query { HelloQuery{} }" }'

         #"query": "query { viewer { user { name } } }"

#curl -g 'http://localhost:9999/graphql?query={HelloQuery}'
#curl http://localhost:9999/graphql -d '{ "query" : "{HelloQuery}" }' 


#curl -XPOST http://localhost:9999/graphql -H 'Content-Type: application/json' -d '{"query": "query M { HelloQuery }"}'

curl -XPOST http://localhost:9999/graphql -H 'Content-Type: application/json' -d  \
'{
  "query": "query M { HelloQuery }"
}'
