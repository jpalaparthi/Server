// service-type project main.go
//
package main

import (
	"encoding/json"
	"fmt"
	"log"
	"mgoSession"
	"models"
	"net/http"
	_ "reflect"
	"time"

	"github.com/gorilla/mux"
	_ "gopkg.in/mgo.v2/bson"
)

var session *mgoSession.Session

var err error

const MONGO_CONN = "mongodb://localhost:27017"

func main() {
	log.Println("server about to start on 8090 port")
	//start mongo session

	session, err = mgoSession.New(MONGO_CONN, "local")
	if err != nil {
		log.Fatal("mongodb database is not connected")
	} else {
		log.Println("mongodb session has been created")
		r := mux.NewRouter()
		srv := &http.Server{
			Handler: r,
			Addr:    "127.0.0.1:8090",
			// Good practice: enforce timeouts for servers you create!
			WriteTimeout: 15 * time.Second,
			ReadTimeout:  15 * time.Second,
		}

		r.HandleFunc("/", all)
		r.HandleFunc("/service-type/add", Add)
		r.HandleFunc("/service-type/update", Update)
		r.HandleFunc("/service-type/delete/{id}", DeleteByID)
		r.HandleFunc("/service-type/del/all/", DeleteAll)
		r.HandleFunc("/service-type/list/{id}", ListByID)
		r.HandleFunc("/service-type/list/all/", List)
		log.Fatal(srv.ListenAndServe())
	}
}

func all(w http.ResponseWriter, r *http.Request) {

	fmt.Fprintln(w, "get all uris...")
}

func Add(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		if IsAthenticatedUser(r) {
			decoder := json.NewDecoder(r.Body)
			defer r.Body.Close()
			var st models.Service_Type
			err := decoder.Decode(&st)
			if err != nil {
				w.WriteHeader(400)
				fmt.Fprintln(w, err.Error())
			} else {
				msg := models.ValidateService_Type(st)
				if msg == "" {
					msg2 := models.ValidateService_TypeforUpdate(st)
					if msg2 == "" {
						err := session.Insert("service_type", st)
						if err != nil {
							w.WriteHeader(400)
							fmt.Fprintln(w, err.Error())

						} else {
							fmt.Fprintln(w, "successfully added")

						}
					} else {
						w.WriteHeader(400)
						fmt.Fprintln(w, msg2)

					}

				} else {
					w.WriteHeader(400)
					fmt.Fprintln(w, msg)

				}
			}
		} else {
			w.WriteHeader(400)
			fmt.Fprintln(w, "Message", "not athenticated user")
		}
	}

}

func Update(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		if IsAthenticatedUser(r) {
			decoder := json.NewDecoder(r.Body)
			defer r.Body.Close()
			var st models.Service_Type
			err := decoder.Decode(&st)
			if err != nil {
				w.WriteHeader(400)
				fmt.Fprintln(w, err.Error())
			} else {
				msg := models.ValidateService_Type(st)
				if msg == "" {
					msg2 := models.ValidateService_TypeforUpdate(st)
					if msg2 == "" {
						err := session.UpdateByID("service_type", st.ID, st)
						if err != nil {
							w.WriteHeader(400)
							fmt.Fprintln(w, err.Error())

						} else {
							fmt.Fprintln(w, "successfully updated")

						}
					} else {
						w.WriteHeader(400)
						fmt.Fprintln(w, msg2)
					}
				} else {
					w.WriteHeader(400)
					fmt.Fprintln(w, msg)
				}
			}
		} else {
			w.WriteHeader(400)
			fmt.Fprintln(w, "Message", "not athenticated user")
		}
	}

}

func List(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		result, err := session.ListAll("service_type")
		if err != nil {
			w.WriteHeader(400)
			fmt.Fprintln(w, err.Error())
		} else {

			jData, err := json.Marshal(result)
			if err != nil {
				w.WriteHeader(400)
				fmt.Fprintln(w, err.Error())
			}
			w.Header().Set("Content-Type", "application/json")
			w.Write(jData)

		}
	}
}

func ListByID(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		vars := mux.Vars(r)
		_id := vars["id"]
		fmt.Println(_id)
		result, err := session.ListByID("service_type", _id)
		if err != nil {
			w.WriteHeader(400)
			fmt.Fprintln(w, err.Error())
		} else {

			jData, err := json.Marshal(result)
			if err != nil {
				w.WriteHeader(400)
				fmt.Fprintln(w, err.Error())
			}
			w.Header().Set("Content-Type", "application/json")
			w.Write(jData)

		}
	}
}

func DeleteAll(w http.ResponseWriter, r *http.Request) {
	if r.Method == "DELETE" {
		err := session.DeleteAll("service_type")
		if err != nil {
			w.WriteHeader(400)
			fmt.Fprintln(w, err.Error())
		} else {
			fmt.Fprintln(w, "successfully deleted")
		}
	}
}

func DeleteByID(w http.ResponseWriter, r *http.Request) {
	if r.Method == "DELETE" {
		vars := mux.Vars(r)
		_id := vars["id"]
		err := session.DeleteByID("service_type", _id)
		if err != nil {
			w.WriteHeader(400)
			fmt.Fprintln(w, err.Error())
		} else {
			fmt.Fprintln(w, "successfully deleted")
		}
	}
}

// helper functions
func IsAthenticatedUser(r *http.Request) bool {
	username := r.Header.Get("username")
	password := r.Header.Get("password")
	role := r.Header.Get("role")
	if username == "admin" && password == "admin" && role == "admin" {
		return true
	} else {
		return false
	}

}
