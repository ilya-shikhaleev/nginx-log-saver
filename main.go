// nginx-log-saver project main.go
package main

import (
	"encoding/json"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"io"
	"log"
	"net/http"
	"strings"
)

type RequestData struct {
	Project string      `json:"Project"`
	Date    string      `bson:"Date" json:"Date"`
	Data    interface{} `json:"Data"`
}

func parseRequest(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	var requestData RequestData
	rawValue := r.FormValue("data")
	log.Println(rawValue)

	decoder := json.NewDecoder(strings.NewReader(rawValue))
	err := decoder.Decode(&requestData)
	if err != nil {
		log.Println("Bad request ", requestData)
		log.Println("Err is ", err)
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	session, err := mgo.Dial("mongo-container")
	if err != nil {
		log.Println("Err is ", err)
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}
	defer session.Close()

	// Optional. Switch the session to a monotonic behavior.
	session.SetMode(mgo.Monotonic, true)

	c := session.DB("test").C(requestData.Project)

	var results []RequestData
	err = c.Find(bson.M{"Date": requestData.Date}).All(&results)

	if err != nil || len(results) > 0 {
		log.Println("Err is ", err)
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}
	log.Println("results", results)

	err = c.Insert(&requestData)
	if err != nil {
		log.Println("Err is ", err)
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	log.Println(requestData.Project)
	log.Println(requestData.Date)
	log.Println(requestData.Data)
	w.WriteHeader(http.StatusOK)
	io.WriteString(w, "I'm ready!")
}

func main() {	
	log.Println("Application is started.")
	defer log.Println("Application is closed.")
	http.HandleFunc("/", parseRequest)
	log.Fatal(http.ListenAndServe(":8080", nil))	
}
