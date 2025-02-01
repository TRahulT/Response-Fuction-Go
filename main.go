package main

import (
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

type Movie struct {
	ID       string    `json:"id"`
	Isbn     string    `json:"isbn"`

	Director *Director `json:"director"`
}

type Director struct {
	Firstname string `json:"firstname"`
	Lastname  string `json:"lastname"`
}

// Response is the structure for all API responses
type Response struct {
	Successful bool        `json:"successful"`
	Message    string      `json:"message"`
	Data       interface{} `json:"data,omitempty"`
}

var movies []Movie

func respondWithJSON(w http.ResponseWriter, status int, response Response) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(response)
}

func getMovies(w http.ResponseWriter, r *http.Request) {
	response := Response{
		Successful: true,
		Message:    "Movies retrieved successfully",
		Data:       movies,
	}
	respondWithJSON(w, http.StatusOK, response)
}

func deleteMovie(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id := params["id"]

	for index, item := range movies {
		if item.ID == id {
			movies = append(movies[:index], movies[index+1:]...)
			response := Response{
				Successful: true,
				Message:    "Movie deleted successfully",
			}
			respondWithJSON(w, http.StatusOK, response)
			return
		}
	}

	response := Response{
		Successful: false,
		Message:    "Movie not found",
	}
	respondWithJSON(w, http.StatusNotFound, response)
}

func getMovie(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id := params["id"]

	for _, item := range movies {
		if item.ID == id {
			response := Response{
				Successful: true,
				Message:    "Movie retrieved successfully",
				Data:       item,
			}
			respondWithJSON(w, http.StatusOK, response)
			return
		}
	}

	response := Response{
		Successful: false,
		Message:    "Movie not found",
	}
	respondWithJSON(w, http.StatusNotFound, response)
}

func createMovie(w http.ResponseWriter, r *http.Request) {
	var movie Movie
	if err := json.NewDecoder(r.Body).Decode(&movie); err != nil {
		response := Response{
			Successful: false,
			Message:    "Invalid input",
		}
		respondWithJSON(w, http.StatusBadRequest, response)
		return
	}

	if  movie.Isbn == "" || movie.Director == nil {
		response := Response{
			Successful: false,
			Message:    "Missing required fields",
		}
		respondWithJSON(w, http.StatusBadRequest, response)
		return
	}

	movie.ID = strconv.Itoa(rand.Intn(100000000))
	movies = append(movies, movie)
	response := Response{
		Successful: true,
		Message:    "Movie created successfully",
		Data:       movie,
	}
	respondWithJSON(w, http.StatusCreated, response)
}

func updateMovie(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id := params["id"]

	for index, item := range movies {
		if item.ID == id {
			movies = append(movies[:index], movies[index+1:]...)

			var movie Movie
			if err := json.NewDecoder(r.Body).Decode(&movie); err != nil {
				response := Response{
					Successful: false,
					Message:    "Invalid input",
				}
				respondWithJSON(w, http.StatusBadRequest, response)
				return
			}

			if  movie.Isbn == "" || movie.Director == nil {
				response := Response{
					Successful: false,
					Message:    "Missing required fields",
				}
				respondWithJSON(w, http.StatusBadRequest, response)
				return
			}

			movie.ID = id
			movies = append(movies, movie)
			response := Response{
				Successful: true,
				Message:    "Movie updated successfully",
				Data:       movie,
			}
			respondWithJSON(w, http.StatusOK, response)
			return
		}
	}

	response := Response{
		Successful: false,
		Message:    "Movie not found",
	}
	respondWithJSON(w, http.StatusNotFound, response)
}

func main() {
	r := mux.NewRouter()

	movies = append(movies, Movie{ID: "1", Isbn: "438227",  Director: &Director{Firstname: "John", Lastname: "Carter"}})
	movies = append(movies, Movie{ID: "2", Isbn: "875483",  Director: &Director{Firstname: "Rahul", Lastname: "Gujjar"}})
	movies = append(movies, Movie{ID: "3", Isbn: "432553",  Director: &Director{Firstname: "Raju", Lastname: "Photographer"}})

	r.HandleFunc("/movies", getMovies).Methods("GET")
	r.HandleFunc("/movies/{id}", getMovie).Methods("GET")
	r.HandleFunc("/movies", createMovie).Methods("POST")
	r.HandleFunc("/movies/{id}", updateMovie).Methods("PUT")
	r.HandleFunc("/movies/{id}", deleteMovie).Methods("DELETE")

	fmt.Println("Starting server at port 8000")
	log.Fatal(http.ListenAndServe(":8000", r))
}
