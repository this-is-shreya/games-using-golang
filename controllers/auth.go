package controllers

import (
	"encoding/json"
	"net/http"

	"example.com/games/database"
	"example.com/games/middleware"
)

type user struct {
	Email      string `json:"email"`
	First_name string `json:"first_name"`
	Last_name  string `json:"last_name"`
}
type token struct {
	TokenString string
	Success     bool
	Message     string
	Error       error
}

func Signup(w http.ResponseWriter, r *http.Request) {
	var resp = &token{}
	w.Header().Set("Content-Type", "application/json")
	Db := database.Database
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
	var n int32
	err := Db.QueryRow(`SELECT COUNT(*) FROM user WHERE email=?`, u.Email).Scan(&n)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		resp.Success = false
		resp.Message = "Internal server error2"
		json.NewEncoder(w).Encode(resp)
		return

	}
	var roomNo int32
	error := Db.QueryRow(`SELECT COUNT(*) FROM user`).Scan(&roomNo)
	if error != nil {
		w.WriteHeader(http.StatusInternalServerError)
		resp.Success = false
		resp.Message = "Internal server error"
		json.NewEncoder(w).Encode(resp)
		return

	}
	if n == 0 {

		_, err := Db.Exec(`INSERT INTO user VALUES(?,?,?,?)`, u.Email, u.First_name, u.Last_name, (roomNo + 1))
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			resp.Success = false
			resp.Message = "Internal server error"
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
			resp.Message = "Internal server error"
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

	w.Header().Set("Content-Type", "application/json")
	var resp = &token{}

	Db := database.Database
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
	var n int32
	err := Db.QueryRow(`SELECT COUNT(*) FROM user WHERE email=?`, u.Email).Scan(&n)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		resp.Success = false
		resp.Error = err
		json.NewEncoder(w).Encode(resp)

		return

	}

	if n == 0 {
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
