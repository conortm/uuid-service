package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"sync"
	"time"

	"github.com/pborman/uuid"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

const (
	uuidPath       = "uuid"
	databaseName   = "uuidService"
	collectionName = "uuids"
)

// UUID is a holder for uuids.
type UUID struct {
	ID      bson.ObjectId `json:"id" bson:"_id"`
	Key     string        `json:"key" bson:"key"`
	UUID    string        `json:"uuid" bson:"uuid"`
	Created time.Time     `json:"created" bson:"created"`
}

// UUIDController controls UUIDs.
type UUIDController struct {
	session *mgo.Session
	sync.RWMutex
}

// UUIDPost is the format of a POST to the uuid endpoint
type UUIDPost struct {
	Key string `json:"key"`
}

func newUUIDController(s *mgo.Session) *UUIDController {
	return &UUIDController{session: s}
}

func (uc *UUIDController) createUUID(key string) (*UUID, error) {
	uc.Lock()
	defer uc.Unlock()
	u := &UUID{
		ID:      bson.NewObjectId(),
		Key:     key,
		UUID:    uuid.New(),
		Created: time.Now(),
	}
	err := uc.session.DB(databaseName).C(collectionName).Insert(u)
	return u, err
}

func (uc *UUIDController) getUUID(key string) (*UUID, error) {
	uc.RLock()
	defer uc.RUnlock()
	u := &UUID{}
	err := uc.session.DB(databaseName).C(collectionName).Find(bson.M{"key": key}).One(u)
	return u, err
}

func uuidHandler(w http.ResponseWriter, r *http.Request) {
	s, err := mgo.Dial(os.Getenv("MONGO_URL"))
	defer s.Close()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	uc := newUUIDController(s)
	var u *UUID
	var key string
	switch r.Method {
	case "GET":
		key = r.URL.Path[len("/"+uuidPath+"/"):]
		// TODO: validate key?
		u, err = uc.getUUID(key)
		if err != nil {
			http.NotFound(w, r)
			return
		}
	case "POST":
		var p UUIDPost
		if err = json.NewDecoder(r.Body).Decode(&p); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		key = p.Key
		// TODO: validate key?
		// Check if uuid exists.
		u, err = uc.getUUID(key)
		if u.Key != key {
			u, err = uc.createUUID(key)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
		}
	default:
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(u); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func main() {
	http.HandleFunc("/"+uuidPath, uuidHandler)
	http.HandleFunc("/"+uuidPath+"/", uuidHandler)
	log.Fatal(http.ListenAndServe(":"+os.Getenv("PORT"), nil))
}
