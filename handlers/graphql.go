package handlers

import (
	"encoding/json"
	"example.com/graphql/schema"
	"fmt"
	"github.com/golang/glog"
	"github.com/gorilla/websocket"
	"github.com/graphql-go/graphql"
	"github.com/graphql-go/handler"
	"log"
	"net/http"
	"context"
	"github.com/trevex/graphql-go-subscription"
    "github.com/trevex/graphql-go-subscription/examples/pubsub"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
} // use default options

type GraphQLMessage struct {
	Query         string                 `json:"query"`
	OperationName string                 `json:"operationName"`
	Variables     map[string]interface{} `json:"variables"`
}

//TODO this should be an enum based on the types in the graphql-ws protocol
type Message struct {
	Type string `json:"type"`
	//Payload map[string]string 	 `json:"payload"`
	Payload GraphQLMessage `json:"payload"`
}

var (
	Clients   map[*websocket.Conn]bool
	Broadcast chan *graphql.Result
    graphqlPubSub *pubsub.PubSub
	subscriptionManager *subscription.SubscriptionManager
)

func init() {
	Clients = make(map[*websocket.Conn]bool)
	Broadcast = make(chan *graphql.Result)
	graphqlPubSub = pubsub.New(4)
	subscriptionManager = subscription.NewSubscriptionManager(subscription.SubscriptionManagerConfig{
		Schema: schema.Schema,
		PubSub: graphqlPubSub,
	})

	// Go Routine to handle messages from the broker
	go handleMessages()
}

func WebsocketRegisterHandler() http.HandlerFunc {
	return func(response http.ResponseWriter, request *http.Request) {
		glog.Infof("Handling Websocket....")
		responseHeaders := http.Header{"Sec-WebSocket-Protocol": {"graphql-ws"}}
		conn, err := upgrader.Upgrade(response, request, responseHeaders)
		if err != nil {
			log.Println(err)
			return
		}

		for {
			glog.Infof("Looping round read message...")

			// Read Message from the connection
			//
			mt, message, err := conn.ReadMessage()
			if err != nil {
				log.Println("read:", err)
				break
			}

			log.Printf("received message: [%s] with type ID [%d]", message, mt)

			var payloadResponse []byte = nil

			// Convert the message to a struct
			var websocketMessage Message
			_ = json.Unmarshal(message, &websocketMessage)
			log.Printf("received message: type [%s], payload [%s] with type ID [%d]", websocketMessage.Type, websocketMessage.Payload, mt)

			if websocketMessage.Type == "connection_init" {
				glog.Info("Received Connection Init - generating an ack")
				// Send Ack
				response := Message{
					Type: "connection_ack",
				}
				payloadResponse, _ = json.Marshal(response)
				if payloadResponse != nil {
					log.Printf("writing: %s", payloadResponse)
					err = conn.WriteMessage(mt, payloadResponse)
					if err != nil {
						log.Println("write:", err)
						break
					}
				}

			} else if websocketMessage.Type == "start" {
				glog.Infof("Received Start - generating a subscription message with payload [%s]", websocketMessage.Payload)





				payloadResponse = doGraphQLStuff(response, request, websocketMessage)

				// TODO - we need to tidy this up.  But if we get to the point, where we have
				// made a subscription request, then we should break out of the loop.
				//
				Clients[conn] = true
				break

			} else {
				glog.Errorf("Received unkown message type [%s] with payload [%s]\n", websocketMessage.Type, websocketMessage.Payload)
				break
			}

		}

		glog.Infof("BYE....")
	}
}

func GraphQLHandler() http.HandlerFunc {
	return func(response http.ResponseWriter, request *http.Request) {

		if origin := request.Header.Get("Origin"); origin != "" {
			glog.Info("Looks like we got an Origin")
			response.Header().Set("Access-Control-Allow-Origin", origin)
			fmt.Println("Origin: " + origin)
			response.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
			response.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
		}
		// Stop here if its Preflighted OPTIONS request
		if request.Method == "OPTIONS" {
			glog.Info("Looks like we got an Options")
			return
		}

		// parse http.Request into handler.RequestOptions
		opts := handler.NewRequestOptions(request)

		rootValue := map[string]interface{}{
			"response": response,
			"request":  request,
			"viewer":   "john_doe", // TODO extract identifier from the ctx for subscription messages

		}

		// Construct the graphql.go Query from the request.
		// N.B we pass the ctx in here, so that the graphql.go mutation
		// and send the payload to the broker.
		params := graphql.Params{
			Schema:         schema.Schema,
			RequestString:  opts.Query,
			VariableValues: opts.Variables,
			OperationName:  opts.OperationName,
			RootObject:     rootValue,
			Context:        context.WithValue(context.Background(), "broadcast", Broadcast),
		}

		glog.Infof("XXX Root: %s", rootValue)
		glog.Infof("XXX Operation Name: %s", opts.OperationName)
		glog.Infof("XXX Variables: %s", opts.Variables)
		glog.Infof("XXX Query: %s", opts.Query)

		// If there was an error, it should be
		// included in the result, so we send it back to the client
		result := graphql.Do(params)
		payload, err := json.Marshal(result)
		if err != nil {
			log.Println("[CommandHandler] Unable to marshal JSON for publishing: ", err)
			http.Error(response, err.Error(), http.StatusInternalServerError)
			return
		}

		if opts.OperationName == "addMessage" {

			// Write Message to message queue

			//websocketMessage := WebsocketMessage{
			//	Type:    "data",
			//	Id:      "1",
			//	Payload: result,
			//}

			//{"type":"data","id":"1","payload":{"data":{"messageAdded":{"id":"5","text":"hello, world","__typename":"Message"}}}}

			//{"type":"data","id":"1","payload":{"data":{"addMessage":{"__typename":"Message","id":"3","text":"Hello, World"}}}}

			// {"data":{"messageAdded":{"__typename":"Message","id":null,"text":null}}}

			//subPayload, err := json.Marshal(websocketMessage)
			//if err != nil {
			//	log.Println("[CommandHandler] Unable to marshal JSON for publishing: ", err)
			//	return
			//}

			Broadcast <- result
			graphqlPubSub.Publish("messageAdded", result)

		}

		response.WriteHeader(http.StatusOK)
		//response.Header().Set("Content-Type", "application/json")
		//response.Header().Set("Content-Type", "application/json")
		response.Header().Set("Access-Control-Allow-Origin", "*")
		response.Write(payload)
	}
}

func doGraphQLStuff(response http.ResponseWriter, request *http.Request, message Message) []byte {
	// Send GraphQL PAyload....

	glog.Infof("Operation Name %s", message.Payload.OperationName)
	glog.Infof("Query %s", message.Payload.Query)
	glog.Infof("Variables %s", message.Payload.Variables)

	/*
			XXX Operation Name: addMessage
		    XXX Variables: map[message:map[channelId:1 text:Hello, World]]
		    XXX Query: mutation addMessage($message: MessageInput!) {
		  addMessage(message: $message) {
		    id
		    text
		    __typename
		  }
	*/



	// Make graphql request
	subId, err := subscriptionManager.Subscribe(subscription.SubscriptionConfig{
		Query: message.Payload.Query,
		VariableValues: message.Payload.Variables,
		OperationName:  message.Payload.OperationName,
		Callback: func(result *graphql.Result) error {
			str, _ := json.Marshal(result)
			fmt.Printf("XXXXXX %s XXXXXX", str)
			return nil
		},
	})

	glog.Infof("Successfully Suscribed User: [%s]\n", subId)

	if err != nil {
		glog.Errorf("Error from subscription [%s]\n", err)
	}



	return nil
}


type WebsocketMessage struct {
	Type string `json:"type"`
    Id string `json:"id"`
	Payload *graphql.Result `json:"payload"`
}

func handleMessages() {
	glog.Info("HandleMessages")

	for {
		glog.Info("Handle Message Loop")
		// Grab the next message from the broadcast channel
		msg := <-Broadcast
		payload, err := json.Marshal(msg)
		if err != nil {
			log.Println("[CommandHandler] Unable to marshal JSON for publishing: ", err)
		}
		glog.Infof("Payload: %s", string(payload))

		// Send it out to every client that is currently connected
		for client := range Clients {
			glog.Infof("Writing Message [%s]\n", payload)
			err := client.WriteMessage(1, payload)
			if err != nil {
				log.Printf("error: %v", err)
				client.Close()
				delete(Clients, client)
			}
		}
	}
}
