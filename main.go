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
	uuidPath       = "/uuid/"
	databaseName   = "uuidService"
	collectionName = "uuids"
)

var databaseURL = "mongodb://mongo"

type document struct {
	ID      bson.ObjectId `json:"id" bson:"_id"`
	Key     string        `json:"key" bson:"key"`
	UUID    string        `json:"uuid" bson:"uuid"`
	Created time.Time     `json:"created" bson:"created"`
}

type documentController struct {
	session *mgo.Session
	sync.RWMutex
}

func newDocumentController(s *mgo.Session) *documentController {
	return &documentController{session: s}
}

func (dc *documentController) createDocument(key string) (*document, error) {
	dc.Lock()
	defer dc.Unlock()
	d := &document{
		ID:      bson.NewObjectId(),
		Key:     key,
		UUID:    uuid.New(),
		Created: time.Now(),
	}
	err := dc.session.DB(databaseName).C(collectionName).Insert(d)
	return d, err
}

func (dc *documentController) getDocument(key string) (*document, error) {
	dc.RLock()
	defer dc.RUnlock()
	d := &document{}
	err := dc.session.DB(databaseName).C(collectionName).Find(bson.M{"key": key}).One(d)
	return d, err
}

func uuidHandler(w http.ResponseWriter, r *http.Request) {
	s, err := mgo.Dial(databaseURL)
	defer s.Close()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	dc := newDocumentController(s)
	var d *document
	key := r.URL.Path[len(uuidPath):]
	// TODO: validate key?
	httpStatus := http.StatusOK
	d, err = dc.getDocument(key)
	// TODO: re-work GET/PUT logic.
	switch r.Method {
	case "GET":
		if err != nil {
			http.NotFound(w, r)
			return
		}
	case "PUT":
		if d.Key != key {
			d, err = dc.createDocument(key)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			httpStatus = http.StatusCreated
		}
	default:
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(httpStatus)
	if err := json.NewEncoder(w).Encode(d); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func main() {
	envDatabaseURL := os.Getenv("DATABASE_URL")
	if len(envDatabaseURL) > 0 {
		databaseURL = envDatabaseURL
	}
	http.HandleFunc(uuidPath, uuidHandler)
	log.Fatal(http.ListenAndServe(":8080", nil))
}
