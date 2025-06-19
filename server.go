package main

import (
	"net/http"
	"encoding/json"
	"sync"
	"io/ioutil"
	"fmt"
	"time"
	"strings"
)

type Coaster struct {
	Name string `json:"name"`
	Manufacturer string `json:"manufacturer"`
	ID string `json:"id"`
	InPark string `json:"in_park"`
	Height int `json:"height"`
}

type coastersController struct {
	sync.Mutex
	data map[string]Coaster
}

func newCoasterController() *coastersController {
	return &coastersController{
		data: map[string]Coaster{
			"1": {
				Name: "Roller Coaster 1", 
				Manufacturer: "Coaster Inc.", 
				ID: "1", 
				InPark: "Park A", 
				Height: 50,
			},
		},
	}
}

func (h *coastersController) coasters(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		h.index(w, r)
		return
	case http.MethodPost:
		h.store(w, r)
		return
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
		w.Write([]byte("Method not allowed\n"))
		return
	}
}

/**
 * GET /coasters/:id
 */
func (h *coastersController) show(w http.ResponseWriter, r *http.Request) {
	parts := strings.Split(r.URL.Path, "/")
	if len(parts) != 3 {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	
	id := parts[2]

	h.Lock()
	coaster, ok := h.data[id]
	if !ok {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("Coaster not found\n"))
		h.Unlock()
		return
	}
	h.Unlock()

	jsonBtyes, err := json.Marshal(coaster)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(jsonBtyes)
}

/**
 * POST /coasters
 */
func (h *coastersController) store(w http.ResponseWriter, r *http.Request) {
	bodyBytes, err := ioutil.ReadAll(r.Body)
	defer r.Body.Close()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}

	ct := r.Header.Get("Content-Type")
	if ct != "application/json" {
		w.WriteHeader(http.StatusUnsupportedMediaType)
		w.Write([]byte("Content-Type must be application/json\n"))
		return
	}
	
	var coaster Coaster
	err = json.Unmarshal(bodyBytes, &coaster)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Invalid JSON format\n"))
		return
	}

	coaster.ID = fmt.Sprintf("%d", time.Now().UnixNano())
	
	h.Lock()
	h.data[coaster.ID] = coaster
	w.WriteHeader(http.StatusCreated)
	w.Write([]byte("Coaster created successfully\n"))
	defer h.Unlock()
}

/**
 * GET /coasters
 */
func (h *coastersController) index(w http.ResponseWriter, r *http.Request) {
	coasters := make([]Coaster, len(h.data))

	h.Lock()
	i := 0
	for _, coaster := range h.data {
		coasters[i] = coaster
		i++
	}
	h.Unlock()

	jsonBtyes, err := json.Marshal(coasters)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(jsonBtyes)
}

func main() {
	coasterController := newCoasterController();

	http.HandleFunc("/coasters", coasterController.coasters)
	http.HandleFunc("/coasters/", coasterController.show)

	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		panic(err)
	}
}