package handlers

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"workmate/pkg/ioboundtask"
)

type requestTask struct {
	ID   int    `json:"ID"`
	Work string `json:"Work"`
}

func CreateTask(rw http.ResponseWriter, r *http.Request) {
	data, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}

	work := struct {
		Work string `json:"work"`
	}{}

	err = json.Unmarshal(data, &work)
	if err != nil {
		log.Println(err)
		http.Error(rw, err.Error(), http.StatusBadRequest)
		return
	}

	go ioboundtask.CreateTask(work.Work)

	rw.WriteHeader(http.StatusCreated)
}

func DeleteTask(rw http.ResponseWriter, r *http.Request) {
	data, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}

	id := struct {
		ID int `json:"id"`
	}{}

	err = json.Unmarshal(data, &id)
	if err != nil {
		log.Println(err)
		http.Error(rw, err.Error(), http.StatusBadRequest)
		return
	}

	go ioboundtask.DeleteTask(id.ID)

	rw.WriteHeader(http.StatusOK)
}

func GetTask(rw http.ResponseWriter, r *http.Request) {
	data, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}

	id := struct {
		ID int `json:"id"`
	}{}

	err = json.Unmarshal(data, &id)
	if err != nil {
		log.Println(err)
		http.Error(rw, err.Error(), http.StatusBadRequest)
		return
	}

	err = ioboundtask.GetTask(rw, id.ID)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}
}

func GetAllTask(rw http.ResponseWriter, r *http.Request) {
	err := ioboundtask.GetAllTask(rw)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}
}
