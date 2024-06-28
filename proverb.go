package main

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"
)

type Proverb struct {
	Content   string    `json:"content"`
	CreatedOn time.Time `json:"createdon"`
}

var proverbStore = make(map[string]Proverb)
var id int = 0

func PostProverbHandler(w http.ResponseWriter, r *http.Request) {
	var proverb Proverb
	err := json.NewDecoder(r.Body).Decode(&proverb)
	if err != nil {
		panic(err)
	}
	proverb.CreatedOn = time.Now()
	id++
	k := strconv.Itoa(id)
	proverbStore[k] = proverb
	j, err := json.Marshal(proverb)
	if err != nil {
		panic(err)
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	w.Write(j)
}

func GetProverbHandler(w http.ResponseWriter, r *http.Request) {
	var proverb []Proverb
	for _, v := range proverbStore {
		proverb = append(proverb, v)
	}
	w.Header().Set("Content-Type", "application/json")
	j, err := json.Marshal(proverb)
	if err != nil {
		panic(err)
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(j)
}

func GetParticularProverbHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	k := vars["id"]
	if proverb, ok := proverbStore[k]; ok {
		j, err := json.Marshal(proverb)
		if err != nil {
			panic(err)
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write(j)
	} else {
		w.WriteHeader(http.StatusNotFound)
	}
}

func PutProverbHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	k := vars["id"]
	var newproverb Proverb
	err := json.NewDecoder(r.Body).Decode(&newproverb)
	if err != nil {
		panic(err)
	}
	if proverb, ok := proverbStore[k]; ok {
		newproverb.CreatedOn = proverb.CreatedOn
		delete(proverbStore, k)

		proverbStore[k] = newproverb
		j, err := json.Marshal(newproverb)
		if err != nil {
			panic(err)
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write(j)
	} else {
		w.WriteHeader(http.StatusNotFound)
	}
}

func DeleteProverbHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	k := vars["id"]
	if _, ok := proverbStore[k]; ok {
		delete(proverbStore, k)
		w.WriteHeader(http.StatusNoContent)
	} else {
		w.WriteHeader(http.StatusNotFound)
	}
}

func main() {
	r := mux.NewRouter().StrictSlash(false)
	r.HandleFunc("/api/proverbs/", PostProverbHandler).Methods("POST")
	r.HandleFunc("/api/proverbs/", GetProverbHandler).Methods("GET")
	r.HandleFunc("/api/proverbs/{id}", GetParticularProverbHandler).Methods("GET")
	r.HandleFunc("/api/proverbs/{id}", PutProverbHandler).Methods("PUT")
	r.HandleFunc("/api/proverbs/{id}", DeleteProverbHandler).Methods("DELETE")
	server := &http.Server{
		Addr:    ":8080",
		Handler: r,
	}
	server.ListenAndServe()
}
