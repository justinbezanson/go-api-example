package main

import (
	"net/http"
	"encoding/json"
	"sync"
	"io/ioutil"
	"fmt"
	"time"
	"strings"
	"os"
	"math/rand"
)

type Coaster struct {
	Name string `json:"name"`
	Manufacturer string `json:"manufacturer"`
	ID string `json:"id"`
	InPark string `json:"in_park"`
	Height int `json:"height"`
}

type AdminPortal struct {
	password string
}

type coastersController struct {
	sync.Mutex
	data map[string]Coaster
}

func newAdminPortal() *AdminPortal {
	password := os.Getenv("ADMIN_PASSWORD")
	if password == "" {
		panic("ADMIN_PASSWORD environment variable is not set")
	}

	return &AdminPortal{
		password: password,
	}
}

func (a *AdminPortal) handler(w http.ResponseWriter, r *http.Request) {
	user, pass, ok := r.BasicAuth()
	if !ok || user != "admin" || pass != a.password {
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte("401 - Unauthorized\n"))
		return
	}

	w.Write([]byte("<html><body><h1>Welcome to the Admin Portal</h1></body></html>"))
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

func (h *coastersController) randomCosaster(w http.ResponseWriter, r *http.Request) {
	ids := make([]string, len(h.data))
	h.Lock()
	i := 0
	for id := range h.data {
		ids[i] = id
		i++
	}
	defer h.Unlock()

	var target string
	if len(ids) == 0 {
		w.WriteHeader(http.StatusNotFound)
		return
	} else if len(ids) == 1 {
		target = ids[0]
	} else {
		rand.Seed(time.Now().UnixNano())
		target = ids[rand.Intn(len(ids))]
	}

	w.Header().Set("Location", fmt.Sprintf("/coasters/%s", target))
	w.WriteHeader(http.StatusFound)
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

	if(id == "random") {
		h.randomCosaster(w, r)
		return
	}

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
	adminPortal := newAdminPortal()
	http.HandleFunc("/admin", adminPortal.handler)

	coasterController := newCoasterController();
	http.HandleFunc("/coasters", coasterController.coasters)
	http.HandleFunc("/coasters/", coasterController.show)

	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		panic(err)
	}
}