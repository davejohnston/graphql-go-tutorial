# graphql-go-tutorial
A go implementation of the graphql server, that is compatible with the Apollo demo: https://github.com/apollographql/graphql-tutorial

## Running the GraphQL server
The server runs on port 400.  To run execute

    $ go run  

## Testing with Curl
The server can be tested using curl.   The query and mutations are described in the sub-sections below.   In each case you can use the following curl command to  POST the requests.  You can optionally pass the response through jq to make it easier to read.

    curl -XPOST http://localhost:4000/graphql -H 'Content-Type: application/json' -d \
    '
       <pay load goes here>
    ' | jq


### Queries

This will get all channels

    {"query":"query ChannelsListQuery{channels{id name}}","operationName":"ChannelsListQuery"}

This will return channel details for the specified channel id

    {"query":"query ChannelDetailsQuery($channelId: ID!) {channel(id: $channelId) {id name messages{id text}}}","variables":{"channelId":"1"},"operationName":"ChannelDetailsQuery"}

### Mutations
   This will add a new message

    {"query":"mutation addMessage($message: MessageInput!) {addMessage(message: $message) {id text}}","variables":{"message":{"channelId":"1","text":"Yo"}},"operationName":"addMessage"}

## Testing with Apollo GraphQL Tutorial
Checkout the apollo Tutorial.

    git clone  https://github.com/apollographql/graphql-tutorial

Switch to the t7-end branch

    cd graphql-tutorial
    git checkout t7-end

With the golang server running, start the apollo Client.  

    cd client
    npm install && npm start

The browser should open on localhost:3000.  You can now use the UI to test drive the server.
