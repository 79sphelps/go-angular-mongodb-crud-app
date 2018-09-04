package dao

import (
	"log"

	. "github.com/79sphelps/movies-restapi/models"
	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)


type UsersDAO struct {
	Server   string
	Database string
}

var db *mgo.Database

const (
	COLLECTION = "users"
)

// Establish a connection to database
func (m *UsersDAO) Connect() {
	session, err := mgo.Dial(m.Server)
	if err != nil {
		log.Fatal(err)
	}
	db = session.DB(m.Database)
}

// Find list of Users
func (m *UsersDAO) FindAll() ([]User, error) {
	var users []User
	err := db.C(COLLECTION).Find(bson.M{}).All(&users)
	
	return users, err
}

// Find a user by its id
func (m *UsersDAO) FindById(id string) (User, error) {
	var user User
	err := db.C(COLLECTION).FindId(bson.ObjectIdHex(id)).One(&user)
	return user, err
}

// Insert a user into database
func (m *UsersDAO) Insert(user User) error {
	err := db.C(COLLECTION).Insert(&user)
	return err
}

// Delete an existing user
//func (m *UsersDAO) Delete(user User) error {
func (m *UsersDAO) Delete(id string) error {
	var user User
	db.C(COLLECTION).FindId(bson.ObjectIdHex(id)).One(&user)
	err := db.C(COLLECTION).Remove(&user)
	
	//err := db.C(COLLECTION).Remove(&user)
	return err
}

// Update an existing user
func (m *UsersDAO) Update(user User) error {
	err := db.C(COLLECTION).UpdateId(user.ID, &user)
	return err
}
