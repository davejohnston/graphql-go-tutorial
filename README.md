# graphql-go-tutorial
A go implementation of the graphql server, that is compatible with the Apollo demo: https://github.com/apollographql/graphql-tutorial

## Testing with Curl
The server can be tested using curl.   The query and mutations are described in the sub-sections below.   In each case you can use the following curl command to  POST the requests.  You can optionally pass the response through jq to make it easier to read.

    curl -XPOST http://locahost:4000 -H 'Content-Type: application/json' -d \
    '
       <pay load goes here>
    ' | jq


### Queries

This will get all channels
    ...
### Mutations
   This will add a new message
    {"query":"mutation addMessage($message: MessageInput!) {\n  addMessage(message: $message) {\n    id\n    text\n    __typename\n  }\n}\n","variables":{"message":{"channelId":"1","text":"Yo"}},"operationName":"addMessage"}
