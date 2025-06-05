package main

import "net/http"
import "encoding/json"

type Coaster struct {
	Name string `json:"name"`
	Manufacturer string `json:"manufacturer"`
	ID string `json:"id"`
	InPark string `json:"in_park"`
	Height int `json:"height"`
}

type coastersHandler struct {
	store map[string]Coaster
}

func newCoasterHandler() *coastersHandler {
	return &coastersHandler{
		store: map[string]Coaster{
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

func (h *coastersHandler) get(w http.ResponseWriter, r *http.Request) {
	coasters := make([]Coaster, len(h.store))
	i := 0
	for _, coaster := range h.store {
		coasters[i] = coaster
		i++
	}

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
	coastersHandler := newCoasterHandler()

	http.HandleFunc("/coasters", coastersHandler.get)

	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		panic(err)
	}
}