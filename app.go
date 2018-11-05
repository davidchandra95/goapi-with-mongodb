package main

import (
	"encoding/json"
	"log"
	"net/http"

	"gopkg.in/mgo.v2/bson"

	"github.com/gorilla/mux"
	. "goapi-with-mongodb/config"
	. "goapi-with-mongodb/dao"
	. "goapi-with-mongodb/models"
	utils "goapi-with-mongodb/utils"
)

var config = Config{}
var dao = MoviesDAO{}

func AllMovies(w http.ResponseWriter, r *http.Request) {
	movies, err := dao.FindAll()
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	utils.RespondWithJson(w, http.StatusOK, movies)
}

func FindMovie(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	movie, err := dao.FindById(params["id"])
	if err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "Invalid movie id")
		return
	}
	utils.RespondWithJson(w, http.StatusOK, movie)
}

func CreateMovie(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	var movie Movie
	if err := json.NewDecoder(r.Body).Decode(&movie); err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	movie.ID = bson.NewObjectId()
	if err := dao.Insert(movie); err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	utils.RespondWithJson(w, http.StatusCreated, movie)
}

func UpdateMovie(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	var movie Movie
	if err:= json.NewDecoder(r.Body).Decode(&movie); err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}
	if err := dao.Update(movie); err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	utils.RespondWithJson(w, http.StatusOK, map[string]string{"result": "success"})
}

func DeleteMovie(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	if err := dao.Delete(params["id"]); err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	utils.RespondWithJson(w, http.StatusOK, map[string]string{"result": "success"})
}

func init() {
	config.Read()

	dao.Server = config.Server
	dao.Database = config.Database
	dao.Connect()
}

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/api/movies", AllMovies).Methods("GET")
	r.HandleFunc("/api/movies", CreateMovie).Methods("POST")
	r.HandleFunc("/api/movies", UpdateMovie).Methods("PUT")
	r.HandleFunc("/api/movies/{id}", DeleteMovie).Methods("DELETE")
	r.HandleFunc("/api/movies/{id}", FindMovie).Methods("GET")

	if err := http.ListenAndServe(":3000", r); err != nil {
		log.Fatal(err)
	}
}
