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

const uuidPath = "/uuid/"

type uuidDocument struct {
	ID      bson.ObjectId `json:"id" bson:"_id"`
	Key     string        `json:"key" bson:"key"`
	UUID    string        `json:"uuid" bson:"uuid"`
	Created time.Time     `json:"created" bson:"created"`
}

type database interface {
	insert(u *uuidDocument) error
	get(key string, u *uuidDocument) error
}

type mongoDatabase struct {
	c *mgo.Collection
	sync.RWMutex
}

func newMongoDatabase(s *mgo.Session) *mongoDatabase {
	return &mongoDatabase{c: s.DB("uuidService").C("uuids")}
}

func (mdb *mongoDatabase) insert(u *uuidDocument) error {
	mdb.Lock()
	defer mdb.Unlock()
	return mdb.c.Insert(u)
}

func (mdb *mongoDatabase) get(key string, u *uuidDocument) error {
	mdb.RLock()
	defer mdb.RUnlock()
	return mdb.c.Find(bson.M{"key": key}).One(u)
}

func uuidHandler(db database) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		key := r.URL.Path[len(uuidPath):]
		// TODO: validate key?
		httpStatus := http.StatusOK
		u := &uuidDocument{}
		err := db.get(key, u)
		// TODO: re-work GET/PUT logic.
		switch r.Method {
		case "GET":
			if err != nil {
				if err.Error() == "not found" {
					http.NotFound(w, r)
				} else {
					http.Error(w, err.Error(), http.StatusInternalServerError)
				}
				return
			}
		case "PUT":
			if u.Key != key {
				u = &uuidDocument{
					ID:      bson.NewObjectId(),
					Key:     key,
					UUID:    uuid.New(),
					Created: time.Now(),
				}
				err = db.insert(u)
				if err != nil {
					http.Error(w, err.Error(), http.StatusInternalServerError)
					return
				}
				httpStatus = http.StatusCreated
			}
		default:
			http.Error(w, "invalid request method "+r.Method, http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(httpStatus)
		if err := json.NewEncoder(w).Encode(u); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	})
}

func getDatabaseURL() string {
	databaseURL := "mongodb://mongo"
	envDatabaseURL := os.Getenv("DATABASE_URL")
	if len(envDatabaseURL) > 0 {
		databaseURL = envDatabaseURL
	}
	return databaseURL
}

func main() {
	s, err := mgo.Dial(getDatabaseURL())
	if err != nil {
		panic(err)
	}
	defer s.Close()
	db := newMongoDatabase(s)
	http.HandleFunc(uuidPath, uuidHandler(db))
	log.Fatal(http.ListenAndServe(":8080", nil))
}
