package main

import (
	mongoapp "api/mongo_managment"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

func put(w http.ResponseWriter, r *http.Request) {
	var tList mongoapp.TodoList
	reqBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Println(err)
	}

	json.Unmarshal(reqBody, &tList)

	_, err = mongoapp.CreateItem(tList)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		log.Println(err)
		return
	}

	w.WriteHeader(http.StatusCreated)
	log.Printf("Item {%s:%v} was added to Mongo DB\n", tList.Message, tList.Done)
}

func getAll(w http.ResponseWriter, r *http.Request) {
	lists, err := mongoapp.GetItems()
	if err != nil {
		log.Println(err)
		json.NewEncoder(w).Encode("{}")
		log.Println(err)
		return
	}
	json.NewEncoder(w).Encode(lists)
}

func update(w http.ResponseWriter, r *http.Request) {
	var tList mongoapp.TodoList
	reqBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Println(err)
		return
	}
	json.Unmarshal(reqBody, &tList)

	err = mongoapp.UpdateItem(tList.Id, tList.Done, tList.Message)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		log.Println(err)
		return
	}
	json.NewEncoder(w).Encode(http.StatusOK)
	w.WriteHeader(http.StatusOK)
	log.Printf("Item {%s:%v} was updated in Mongo DB\n", tList.Message, tList.Done)
}

func delete(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	res, err := strconv.Atoi(id)
	if err != nil {
		log.Println(err)
		return
	}
	err = mongoapp.DeleteItem(res)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		log.Println(err)
		return
	}
	w.WriteHeader(http.StatusOK)
	log.Printf("Item %s was deleted from Mongo DB\n", id)
}

func getItem(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	res, err := strconv.Atoi(id)
	if err != nil {
		log.Println(err)
		return
	}

	item, err := mongoapp.GetItem(res)
	if err != nil {
		log.Println(err)
		return
	}

	json.NewEncoder(w).Encode(item)
}

func main() {
	mongoapp.Setup()

	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/put", put).Methods("POST")
	router.HandleFunc("/getAll", getAll).Methods("GET")
	router.HandleFunc("/update", update).Methods("POST")
	router.HandleFunc("/delete/{id}", delete).Methods("DELETE")
	router.HandleFunc("/get/{id}", getItem).Methods("GET")
	log.Println("API is listening on port 3000")
	http.ListenAndServe("0.0.0.0:3000", router)
}