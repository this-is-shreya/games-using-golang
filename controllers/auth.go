package controllers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"example.com/games/database"
	"example.com/games/middleware"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type user struct {
	Email      string `json:"email"`
	First_name string `json:"first_name"`
	Last_name  string `json:"last_name"`
}
type mongoResponse struct {
	Email      string `bson:"email"`
	First_name string `bson:"first_name"`
	Last_name  string `bson:"last_name"`
	Room       int    `bson:"room"`
	Id         string `bson:"_id"`
}
type token struct {
	TokenString string
	Success     bool
	Message     string
	Error       error
}

var Coll *mongo.Collection

func Signup(w http.ResponseWriter, r *http.Request) {

	conn := database.IsConnected

	if conn {
		Coll = database.Client.Database("golang").Collection("user")

	} else {
		fmt.Println(conn)
		return
	}
	var resp = &token{}
	w.Header().Set("Content-Type", "application/json")
	// Db := database.Database
	var u user
	decoder := json.NewDecoder(r.Body)
	dErr := decoder.Decode(&u)
	if dErr != nil {
		w.WriteHeader(http.StatusInternalServerError)
		resp.Success = false
		resp.Message = "Internal server error1"
		json.NewEncoder(w).Encode(resp)
		return
	}
	var n = 0
	filter := bson.D{
		{Key: "email", Value: bson.D{{Key: "$eq", Value: u.Email}}},
	}

	var userObtained mongoResponse
	opts := options.FindOne()

	err := Coll.FindOne(database.Ctx, filter, opts).Decode(&userObtained)
	// fmt.Println(err)
	if err != mongo.ErrNoDocuments && err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		resp.Success = false
		resp.Message = "Internal server error2"
		resp.Error = err
		json.NewEncoder(w).Encode(resp)
		return

	}

	if userObtained == (mongoResponse{}) {
		//find the room no, and insert
		opts = options.FindOne().SetSort(bson.D{{Key: "_id", Value: -1}}).SetProjection(bson.D{{Key: "room", Value: 1}})
		err = Coll.FindOne(database.Ctx, bson.D{}, opts).Decode(&n)
		fmt.Println(err)
		if err != mongo.ErrNoDocuments && err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			resp.Success = false
			resp.Message = "Internal server error3"
			json.NewEncoder(w).Encode(resp)
			return

		}
		_, err = Coll.InsertOne(database.Ctx, bson.D{
			{Key: "email", Value: u.Email},
			{Key: "first_name", Value: u.First_name},
			{Key: "last_name", Value: u.Last_name},
			{Key: "room", Value: n + 1}})
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			resp.Success = false
			resp.Message = "Internal server error4"
			json.NewEncoder(w).Encode(resp)
			return

		}
		tokenString := middleware.CreateToken(u.Email, w)
		if tokenString != "" {
			w.WriteHeader(http.StatusCreated)
			var resp = token{
				TokenString: tokenString,
				Success:     true,
				Message:     "signed up successfully",
			}
			json.NewEncoder(w).Encode(resp)

		} else {
			w.WriteHeader(http.StatusInternalServerError)
			resp.Success = false
			resp.Message = "Internal server error5"
			json.NewEncoder(w).Encode(resp)
			return

		}
	} else {
		w.WriteHeader(http.StatusBadRequest)
		resp.Success = false
		resp.Message = "email already exists. Log in to continue"
		json.NewEncoder(w).Encode(resp)
		return

	}
}
func Login(w http.ResponseWriter, r *http.Request) {

	conn := database.IsConnected

	if conn {
		Coll = database.Client.Database("golang").Collection("user")

	} else {
		fmt.Println(conn)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	var resp = &token{}

	var u user
	decoder := json.NewDecoder(r.Body)
	dErr := decoder.Decode(&u)
	if dErr != nil {
		w.WriteHeader(http.StatusInternalServerError)
		resp.Success = false
		resp.Message = "Internal server error"
		json.NewEncoder(w).Encode(resp)

		return

	}

	err := Coll.FindOne(database.Ctx, bson.D{{Key: "email", Value: bson.D{{Key: "$eq", Value: u.Email}}}}, options.FindOne()).Decode(&u)
	fmt.Println(err)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		resp.Success = false
		resp.Error = err
		resp.Message = "email id doesn't exist"
		json.NewEncoder(w).Encode(resp)

		return

	}

	if u == (user{}) {
		w.WriteHeader(http.StatusNotFound)
		resp.Success = false
		resp.Message = "Email doesn't exist"
		json.NewEncoder(w).Encode(resp)

		return

	} else {
		tokenString := middleware.CreateToken(u.Email, w)
		if tokenString == "" {
			w.WriteHeader(http.StatusInternalServerError)
			resp.Success = false
			resp.Message = "Internal server error"
			json.NewEncoder(w).Encode(resp)
			return
		}
		w.WriteHeader(http.StatusCreated)
		resp.TokenString = tokenString
		resp.Success = true
		resp.Message = "Logged in successfully"

		json.NewEncoder(w).Encode(resp)

		return

	}
}
