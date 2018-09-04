package main


import (
	"encoding/json"
	"log"
	"net/http"
	"gopkg.in/mgo.v2/bson"
	"fmt"
	"strings"
	"io/ioutil"
	"bytes"

	"github.com/gorilla/mux"
	. "github.com/79sphelps/go-angular-mongodb-crud-app/app_api/config"
	. "github.com/79sphelps/go-angular-mongodb-crud-app/app_api/dao"
	. "github.com/79sphelps/go-angular-mongodb-crud-app/app_api/models"
)


var config = Config{}
var dao = UsersDAO{}

// GET list of users
func AllUsersEndPoint(w http.ResponseWriter, r *http.Request) {
	users, err := dao.FindAll()
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	respondWithJSON(w, http.StatusOK, users)
}

// GET a user by its ID
func FindUserEndpoint(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	user, err := dao.FindById(params["id"])
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid User ID")
		return
	}
	respondWithJSON(w, http.StatusOK, user)
}

// POST a new user
func CreateUserEndPoint(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	var user User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}
	user.ID = bson.NewObjectId()
	if err := dao.Insert(user); err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	respondWithJSON(w, http.StatusCreated, user)
}

// PUT update an existing user
func UpdateUserEndPoint(w http.ResponseWriter, r *http.Request) {
	/*
	buf, bodyErr := ioutil.ReadAll(r.Body)
	if bodyErr != nil {
		log.Print("bodyErr ", bodyErr.Error())
		http.Error(w, bodyErr.Error(), http.StatusInternalServerError)
		return
	}
	
	rdr1 := ioutil.NopCloser(bytes.NewBuffer(buf))
	rdr2 := ioutil.NopCloser(bytes.NewBuffer(buf))
	log.Printf("BODY: %q", rdr1)
	r.Body = rdr2
	*/

	defer r.Body.Close()
	var user User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}
	if err := dao.Update(user); err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	respondWithJSON(w, http.StatusOK, map[string]string{"result": "success"})
}

// DELETE an existing user
func DeleteUserEndPoint(w http.ResponseWriter, r *http.Request) {
	fmt.Printf("Req: %s %s\n", r.Host, r.URL.Path) 
	tmp := strings.Split(r.URL.Path, "/")
	id := tmp[3]
	fmt.Printf("ID: %s\n", id)

	defer r.Body.Close()

	if err := dao.Delete(id); err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	respondWithJSON(w, http.StatusOK, map[string]string{"result": "success"})
}

// DELETE an existing user
func DeleteUserEndPoint_BACKUP(w http.ResponseWriter, r *http.Request) {
	fmt.Printf("Req: %s %s\n", r.Host, r.URL.Path) 
	tmp := strings.Split(r.URL.Path, "/")
	id := tmp[3]
	fmt.Printf("ID: %s\n", id)

	//-------------------------------------------------
	buf, bodyErr := ioutil.ReadAll(r.Body)
	if bodyErr != nil {
		log.Print("bodyErr ", bodyErr.Error())
		http.Error(w, bodyErr.Error(), http.StatusInternalServerError)
		return
	}
	
	rdr1 := ioutil.NopCloser(bytes.NewBuffer(buf))
	rdr2 := ioutil.NopCloser(bytes.NewBuffer(buf))
	log.Printf("BODY: %q", rdr1)
	r.Body = rdr2

	//-------------------------------------------------

	defer r.Body.Close()
	/*
	var user User
	
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}
	*/

	//user := new(User)
	//var user User
	//user.ID = bson.NewObjectId()
	//user.ID = bson.ObjectIdHex(id)

	//if err := dao.Delete(user); err != nil {
	if err := dao.Delete(id); err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	respondWithJSON(w, http.StatusOK, map[string]string{"result": "success"})
}




// formatRequest generates ascii representation of a request
func formatRequest(r *http.Request) string {
 // Create return string
 var request []string
 // Add the request string
 url := fmt.Sprintf("%v %v %v", r.Method, r.URL, r.Proto)
 request = append(request, url)
 // Add the host
 request = append(request, fmt.Sprintf("Host: %v", r.Host))
 // Loop through headers
 for name, headers := range r.Header {
   name = strings.ToLower(name)
   for _, h := range headers {
     request = append(request, fmt.Sprintf("%v: %v", name, h))
   }
 }
 
 // If this is a POST, add post data
 if r.Method == "POST" {
    r.ParseForm()
    request = append(request, "\n")
    request = append(request, r.Form.Encode())
 } 
  // Return the request as a string
  return strings.Join(request, "\n")
}



func respondWithError(w http.ResponseWriter, code int, msg string) {
	respondWithJSON(w, code, map[string]string{"error": msg})
}

func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	response, _ := json.Marshal(payload)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(response)
}

// Parse the configuration file 'config.toml', and establish a connection to DB
func init() {
	config.Read()

	dao.Server = config.Server
	dao.Database = config.Database
	dao.Connect()
}

// Define HTTP request routes
func main() {
	r := mux.NewRouter()
	r.HandleFunc("/api/users", AllUsersEndPoint).Methods("GET")
	r.HandleFunc("/api/users", CreateUserEndPoint).Methods("POST")
	r.HandleFunc("/api/users/{id}", UpdateUserEndPoint).Methods("PUT")
	r.HandleFunc("/api/users/{id}", DeleteUserEndPoint).Methods("DELETE")
	r.HandleFunc("/api/users/{id}", FindUserEndpoint).Methods("GET")
	if err := http.ListenAndServe(":3000", r); err != nil {
		log.Fatal(err)
	}
}
