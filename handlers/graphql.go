package handlers

import (
	"encoding/json"
	"github.com/davejohnston/graphql-go-tutorial/schema"
	"fmt"
	"github.com/golang/glog"
	"github.com/gorilla/websocket"
	"github.com/graphql-go/graphql"
	"github.com/graphql-go/handler"
	"log"
	"net/http"
	"github.com/trevex/graphql-go-subscription"
    "github.com/trevex/graphql-go-subscription/examples/pubsub"
	"sync"
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
	Id string `json"id"`
	//Payload map[string]string 	 `json:"payload"`
	Payload GraphQLMessage `json:"payload"`
}

type SubscriptionIdMap struct {
	GraphqlRequestId string
	SubscriptionId subscription.SubscriptionId
}

var (
	Clients   map[*websocket.Conn]SubscriptionIdMap
	graphqlPubSub *pubsub.PubSub
	subscriptionManager *subscription.SubscriptionManager
)

func init() {
	Clients = make(map[*websocket.Conn]SubscriptionIdMap)
	graphqlPubSub = pubsub.New(4)
	subscriptionManager = subscription.NewSubscriptionManager(subscription.SubscriptionManagerConfig{
		Schema: schema.Schema,
		PubSub: graphqlPubSub,
	})

}

func GraphQLHandler() http.HandlerFunc {
	return func(response http.ResponseWriter, request *http.Request) {

		if origin := request.Header.Get("Origin"); origin != "" {
			response.Header().Set("Access-Control-Allow-Origin", origin)
			response.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
			response.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
		}
		// Stop here if its Preflighted OPTIONS request
		if request.Method == "OPTIONS" {
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
		}

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
			fmt.Printf("Publishing Message %v\n", result.Data)
			graphqlPubSub.Publish("messageAdded", result.Data)
		}

		response.WriteHeader(http.StatusOK)
		response.Header().Set("Access-Control-Allow-Origin", "*")
		response.Write(payload)
	}
}

type SubscriptionMessage struct {
	MessageAdded interface{} `json:"messageAdded"`

}

type WebsocketMessage struct {
	Type string `json:"type"`
	Id string `json:"id"`
	Payload *graphql.Result `json:"payload"`
}

func WebsocketHandler() http.HandlerFunc {
	return func(response http.ResponseWriter, request *http.Request) {
		glog.Infof("Handling Websocket....")
		var mutex = &sync.Mutex{}
		responseHeaders := http.Header{"Sec-WebSocket-Protocol": {"graphql-ws"}}
		conn, err := upgrader.Upgrade(response, request, responseHeaders)

		if err != nil {
			log.Println(err)
			return
		}

		for {
			// Read Message from the connection
			//
			mt, message, err := conn.ReadMessage()
			if err != nil {
				log.Println("read:", err)
				break
			}

			glog.Infof("RECV message: [%s] with type ID [%d]\n", message, mt)

			var payloadResponse []byte = nil

			// Convert the message to a struct
			var websocketMessage Message
			_ = json.Unmarshal(message, &websocketMessage)
			glog.Infof("received message from: [%s] - type [%s], payload [%s] with type ID [%d]\n", conn.RemoteAddr(), websocketMessage.Type, websocketMessage.Payload, mt)


			var subId subscription.SubscriptionId
			if websocketMessage.Type == "connection_init" {
				glog.Info("Received Connection Init - generating an ack")
				// Send Ack
				response := Message{
					Type: "connection_ack",
				}
				payloadResponse, _ = json.Marshal(response)
				if payloadResponse != nil {
					glog.Infof("writing: %s\n", payloadResponse)
					mutex.Lock()
					err = conn.WriteMessage(mt, payloadResponse)
					mutex.Unlock()
					if err != nil {
						log.Println("write:", err)
						break
					}
				}

			} else if websocketMessage.Type == "start" {
				glog.Warningf("Received Start - generating a subscription message for id [%s] with payload [%s]\n", websocketMessage.Id, websocketMessage.Payload)
				//query, _ := json.Marshal(websocketMessage.Payload)

				subId, err = subscriptionManager.Subscribe(subscription.SubscriptionConfig{
					Query: websocketMessage.Payload.Query,
					VariableValues: websocketMessage.Payload.Variables,
					OperationName: websocketMessage.Payload.OperationName,
					Callback: func(result *graphql.Result) error {

						fmt.Printf("Procesing Callback with Result: %v\n", result)
						if result.Errors != nil {
							log.Println("Error trying to find message: ", result.Errors)
							return nil
						}


						payload, _ := json.Marshal(result)
						// We would need to write this back to the channel
						fmt.Printf("Writing payload [%s] to websocket [%s]\n", payload, conn.RemoteAddr().String())


						websocketResponseMessage := WebsocketMessage{
							Type:    "data",
							Id:      websocketMessage.Id,
							Payload: result,
						}



						//conn.WriteMessage(1, websocketMessage)
						mutex.Lock()
						conn.WriteJSON(websocketResponseMessage)
						mutex.Unlock()
						return nil
					},
				})

				// At this point we need to define a mapping from the subscription id (provided in the request) to the
				// subId generated by the SubManager.

				if err != nil {
					fmt.Printf("Error creating subscription %s\n", err)
				}

				fmt.Printf("\t\t\t\tCreating Subscription for %d\n", subId)

				// Not completely thread safe, and assumes there will only ever be one
				// subscription per websocket.
				//
				idMap := SubscriptionIdMap {
					GraphqlRequestId: websocketMessage.Id,
					SubscriptionId: subId,

				}
				Clients[conn] = idMap

			} else if websocketMessage.Type == "stop" {
				idMap := Clients[conn]
				fmt.Printf("STOP Unsubscriving number %d......\n", idMap.SubscriptionId)
				subscriptionManager.Unsubscribe(idMap.SubscriptionId)
			} else {
				glog.Errorf("Received unkown message type [%s] with payload [%s]\n", websocketMessage.Type, websocketMessage.Payload)
			}

		}

		fmt.Printf("GOODBYE SOCKET......\n")
		idMap := Clients[conn]
		subscriptionManager.Unsubscribe(idMap.SubscriptionId)
		conn.Close()
	}
}