//: ----------------------------------------------------------------------------
//: Copyright (C) 2019 Helmut Wahrmann.
//:
//: file:    api-server.go
//: details: REST API-Server. Provides a basic REST interface
//: author:  Helmut Wahrmann
//: date:    03/09/2019
//:
//: Licensed under the Apache License, Version 2.0 (the "License");
//: you may not use this file except in compliance with the License.
//: You may obtain a copy of the License at
//:
//:     http://www.apache.org/licenses/LICENSE-2.0
//:
//: Unless required by applicable law or agreed to in writing, software
//: distributed under the License is distributed on an "AS IS" BASIS,
//: WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
//: See the License for the specific language governing permissions and
//: limitations under the License.
//: ----------------------------------------------------------------------------

package main

import (
	"context"
	"encoding/json"
	"errors"
	"log"
	"net"
	"net/http"
	"strconv"
	"strings"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	mongoclient *mongo.Client
	collection  *mongo.Collection
)

// RestAPIServer represents a RestApiServer
type RestAPIServer struct{}

// NewAPIServer returns a new RestApiServer instance
func NewAPIServer() *RestAPIServer {
	return &RestAPIServer{}
}

// Starts the REST Api Serverand connects to MongoDB on the EndPointServer
func (r *RestAPIServer) run() error {
	var err error

	// Connect to MongoDB
	clientOptions := options.Client().ApplyURI("mongodb://" + opts.EndPointServer + ":27017")
	clientOptions.SetAuth(options.Credential{AuthMechanism: "SCRAM-SHA-1", AuthSource: "admin", Username: opts.APIUser, Password: opts.APIUserPwd})

	mongoclient, err = mongo.Connect(context.TODO(), clientOptions)

	if err != nil {
		opts.Logger.Fatal(err)
		return errors.New("Error connecting to EndPoint Server")
	}

	// Check the connection
	err = mongoclient.Ping(context.TODO(), nil)

	if err != nil {
		log.Fatal(err)
		return errors.New("Error: couldn't reach EndPoint Server")
	}

	// Connect to machinedetail database
	collection = mongoclient.Database("endpoint-server").Collection("machinedetail")
	if err != nil {
		log.Fatal(err)
		return errors.New("Error: couldn't connect to machinedetail database")
	}

	opts.Logger.Infof("Connected to EndPointServer on %s", opts.EndPointServer)

	// Start HTTP Server ro serve REST requests
	mux := http.NewServeMux()
	mux.HandleFunc("/machine/", getMachine())
	addr := net.JoinHostPort("0.0.0.0", strconv.Itoa(opts.MonitorPort))

	opts.Logger.Infof("Starting REST API Server on port %s", strconv.Itoa(opts.MonitorPort))
	err = http.ListenAndServe(addr, mux)
	if err != nil {
		opts.Logger.Fatal(err)
		return err
	}

	return nil
}

// Shut down the connection to MongoDB
func (r *RestAPIServer) shutdown() {
	err := mongoclient.Disconnect(context.TODO())

	if err != nil {
		log.Fatal(err)
	}
	opts.Logger.Info("Connection to EndPointServer closed.")
}

// getMachine returns the result of the mongodb query to machinedetail
func getMachine() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		type Data struct {
			Machine string
			IP      string
			Score   string
			Status  string
		}

		data := Data{}

		machine := strings.Replace(r.URL.Path, "/machine/", "", 1)

		var result bson.Raw

		//filter := bson.D{{"machine.machineName", bson.D{{"$eq", machine}}}}
		filter := bson.D{{"machine.machineName", bson.M{"$regex": "^" + machine + "$", "$options": "i"}}}
		err := collection.FindOne(context.TODO(), filter).Decode(&result)
		if err != nil {
			opts.Logger.Info(err)
			data.Machine = machine
			if strings.Contains(err.Error(), "mongo: no documents in result") {
				data.Status = "404"
			} else {
				data.Status = "500"
			}
		} else {
			var myjson map[string]interface{}
			json.Unmarshal([]byte(result.String()), &myjson)

			data.Machine = myjson["machine"].(map[string]interface{})["machineName"].(string)
			data.Score = myjson["score"].(map[string]interface{})["$numberInt"].(string)
			data.IP = myjson["machine"].(map[string]interface{})["networkInterfaces"].([]interface{})[0].(map[string]interface{})["ipv4"].([]interface{})[0].(string)
			data.Status = "200"
		}

		j, err := json.Marshal(data)
		if err != nil {
			opts.Logger.Info(err)
		}

		if _, err = w.Write(j); err != nil {
			opts.Logger.Info(err)
		}
	}
}
