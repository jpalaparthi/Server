// User-Management project main.go
package main

import (
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"mgoSession"
	"models"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	_ "github.com/gorilla/securecookie"
	_ "gopkg.in/mgo.v2/bson"
)

var session *mgoSession.Session

var err error

const MONGO_CONN = "mongodb://localhost:27017"

var Tokens map[string]string

func main() {
	log.Println("server about to start on 8090 port")
	//start mongo session

	session, err = mgoSession.New(MONGO_CONN, "local")
	if err != nil {
		log.Fatal("mongodb database is not connected")
	} else {
		log.Println("mongodb session has been created")

		Tokens = make(map[string]string)

		r := mux.NewRouter().StrictSlash(false)

		r.HandleFunc("/", links)
		r.HandleFunc("/user/signin", signin)

		r.HandleFunc("/user", getAll("users"))
		r.HandleFunc("/user/{id}", getByID("users"))
		r.HandleFunc("/user/add/", user_add)
		r.HandleFunc("/user/update/", user_update)
		r.HandleFunc("/user/remove/", removeAll("users"))
		r.HandleFunc("/user/remove/{id}", removeByID("users"))

		r.HandleFunc("/user/role/", getAll("roles"))
		r.HandleFunc("/user/role/{id}", getByID("roles"))
		r.HandleFunc("/user/role/add/", role_add)
		r.HandleFunc("/user/role/update/", role_update)
		r.HandleFunc("/user/role/remove/", removeAll("roles"))
		r.HandleFunc("/user/role/remove/{id}", removeByID("roles"))

		r.HandleFunc("/user/role/scope/", getAll("scopes"))
		r.HandleFunc("/user/role/scope/{id}", getByID("scopes"))
		r.HandleFunc("/user/role/scope/add/", scope_add)
		r.HandleFunc("/user/role/scope/update/", scope_update)
		r.HandleFunc("/user/role/scope/remove/", removeAll("scopes"))
		r.HandleFunc("/user/role/scope/remove/{id}", removeByID("scopes"))

		r.HandleFunc("/log", getByID("log"))
		r.HandleFunc("/log/{id}", getByID("log"))
		r.HandleFunc("/log/add", log_add)

		http.ListenAndServe(":8090", r)
	}
}

func links(w http.ResponseWriter, r *http.Request) {

	fmt.Fprintln(w, "get all uris...")
}

//User Related URIs

func user_add(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		if !HasSignIn(r) {
			w.WriteHeader(400)
			fmt.Fprintln(w, "Message", "not athenticated user")
		} else {
			decoder := json.NewDecoder(r.Body)
			defer r.Body.Close()
			var u models.User
			err := decoder.Decode(&u)
			if err != nil {
				w.WriteHeader(400)
				fmt.Fprintln(w, err.Error())
			} else {
				m := make(map[string]interface{})
				m["username"] = u.UserName
				usr, _ := session.FindByQuery("users", m)

				if usr != nil {
					w.WriteHeader(400)
					fmt.Fprintln(w, "User already exists")
				} else {
					msg := models.ValidateUser(u)
					if msg == "" {
						fmt.Println("hitting here")
						err := session.Insert("users", u)
						if err != nil {
							w.WriteHeader(400)
							fmt.Fprintln(w, err.Error())
						} else {
							fmt.Fprintln(w, "successfully added")
						}
					} else {
						w.WriteHeader(400)
						fmt.Fprintln(w, msg)
					}
				}

			}

		}
	}
}

func user_update(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		if HasSignIn(r) {
			decoder := json.NewDecoder(r.Body)
			defer r.Body.Close()
			var u models.User
			err := decoder.Decode(&u)
			if err != nil {
				w.WriteHeader(400)
				fmt.Fprintln(w, err.Error())
			} else {
				msg := models.ValidateUser(u)
				if msg == "" {
					msg2 := models.ValidateUserForUpdate(u)
					if msg2 == "" {
						err := session.UpdateByID("users", u.ID, u)
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

// Role related URI

func role_add(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		if !HasSignIn(r) {
			w.WriteHeader(400)
			fmt.Fprintln(w, "Message", "not athenticated user")
		} else {
			decoder := json.NewDecoder(r.Body)
			defer r.Body.Close()
			var r models.Role
			err := decoder.Decode(&r)
			if err != nil {
				w.WriteHeader(400)
				fmt.Fprintln(w, err.Error())
			} else {
				m := make(map[string]interface{})
				m["name"] = r.Name
				usr, _ := session.FindByQuery("roles", m)

				if usr != nil {
					w.WriteHeader(400)
					fmt.Fprintln(w, "Role already exists")
				} else {
					msg := models.ValidateRole(r)
					if msg == "" {
						fmt.Println("hitting here")
						err := session.Insert("roles", r)
						if err != nil {
							w.WriteHeader(400)
							fmt.Fprintln(w, err.Error())
						} else {
							fmt.Fprintln(w, "successfully added")
						}
					} else {
						w.WriteHeader(400)
						fmt.Fprintln(w, msg)
					}
				}

			}

		}
	}
}

func role_update(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		if HasSignIn(r) {
			decoder := json.NewDecoder(r.Body)
			defer r.Body.Close()
			var r models.Role
			err := decoder.Decode(&r)
			if err != nil {
				w.WriteHeader(400)
				fmt.Fprintln(w, err.Error())
			} else {
				msg := models.ValidateRole(r)
				if msg == "" {
					msg2 := models.ValidateRoleForUpdate(r)
					if msg2 == "" {
						err := session.UpdateByID("roles", r.ID, r)
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

//scope

// Role related URI

func scope_add(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		if !HasSignIn(r) {
			w.WriteHeader(400)
			fmt.Fprintln(w, "Message", "not athenticated user")
		} else {
			decoder := json.NewDecoder(r.Body)
			defer r.Body.Close()
			var s models.Scope
			err := decoder.Decode(&s)
			if err != nil {
				w.WriteHeader(400)
				fmt.Fprintln(w, err.Error())
			} else {
				m := make(map[string]interface{})
				m["name"] = s.Name
				usr, _ := session.FindByQuery("roles", m)

				if usr == nil {
					w.WriteHeader(400)
					fmt.Fprintln(w, "role does not exists")
				} else {
					msg := models.ValidateScope(s)
					if msg == "" {
						err := session.Insert("scopes", s)
						if err != nil {
							w.WriteHeader(400)
							fmt.Fprintln(w, err.Error())
						} else {
							fmt.Fprintln(w, "successfully added")
						}
					} else {
						w.WriteHeader(400)
						fmt.Fprintln(w, msg)
					}
				}

			}

		}
	}
}

func scope_update(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		if HasSignIn(r) {
			decoder := json.NewDecoder(r.Body)
			defer r.Body.Close()
			var s models.Scope
			err := decoder.Decode(&s)
			if err != nil {
				w.WriteHeader(400)
				fmt.Fprintln(w, err.Error())
			} else {
				msg := models.ValidateScope(s)
				if msg == "" {
					msg2 := models.ValidateScopeForUpdate(s)
					if msg2 == "" {
						err := session.UpdateByID("scopes", s.ID, s)
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

//Service Type related HandleFuncs

func serviceType_add(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		if !HasSignIn(r) {
			w.WriteHeader(400)
			fmt.Fprintln(w, "Message", "not athenticated user")
		} else {
			decoder := json.NewDecoder(r.Body)
			defer r.Body.Close()
			var s models.Scope
			err := decoder.Decode(&s)
			if err != nil {
				w.WriteHeader(400)
				fmt.Fprintln(w, err.Error())
			} else {
				m := make(map[string]interface{})
				m["name"] = s.Name
				usr, _ := session.FindByQuery("roles", m)

				if usr == nil {
					w.WriteHeader(400)
					fmt.Fprintln(w, "role does not exists")
				} else {
					msg := models.ValidateScope(s)
					if msg == "" {
						err := session.Insert("scopes", s)
						if err != nil {
							w.WriteHeader(400)
							fmt.Fprintln(w, err.Error())
						} else {
							fmt.Fprintln(w, "successfully added")
						}
					} else {
						w.WriteHeader(400)
						fmt.Fprintln(w, msg)
					}
				}

			}

		}
	}
}

func serviceType_update(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		if HasSignIn(r) {
			decoder := json.NewDecoder(r.Body)
			defer r.Body.Close()
			var s models.Scope
			err := decoder.Decode(&s)
			if err != nil {
				w.WriteHeader(400)
				fmt.Fprintln(w, err.Error())
			} else {
				msg := models.ValidateScope(s)
				if msg == "" {
					msg2 := models.ValidateScopeForUpdate(s)
					if msg2 == "" {
						err := session.UpdateByID("scopes", s.ID, s)
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

func getAll(collection string) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "GET" {
			if !HasSignIn(r) {
				w.WriteHeader(400)
				fmt.Fprintln(w, "Message", "not athenticated user")
			} else {
				result, err := session.ListAll(collection)
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
	}
}

func getByID(collection string) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "GET" {
			if !HasSignIn(r) {
				w.WriteHeader(400)
				fmt.Fprintln(w, "Message", "not athenticated user")
			} else {
				vars := mux.Vars(r)
				_id := vars["id"]

				result, err := session.ListByID(collection, _id)
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
	}
}

func removeByID(collection string) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "DELETE" {
			if !HasSignIn(r) {
				w.WriteHeader(400)
				fmt.Fprintln(w, "Message", "not athenticated user")
			} else {
				vars := mux.Vars(r)
				_id := vars["id"]
				err := session.DeleteByID(collection, _id)
				if err != nil {
					w.WriteHeader(400)
					fmt.Fprintln(w, err.Error())
				} else {
					fmt.Fprintln(w, "successfully removed")
				}
			}
		}

	}

}

func removeAll(collection string) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "DELETE" {
			if !HasSignIn(r) {
				w.WriteHeader(400)
				fmt.Fprintln(w, "Message", "not athenticated user")
			} else {
				err := session.DeleteAll(collection)
				if err != nil {
					w.WriteHeader(400)
					fmt.Fprintln(w, err.Error())
				} else {
					fmt.Fprintln(w, "successfully removed")
				}
			}
		}

	}

}

//Log related URI
func logs(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		result, err := session.ListAll("logs")
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

func log_get(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		vars := mux.Vars(r)
		_id := vars["id"]

		result, err := session.ListByID("logs", _id)
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

func log_add(w http.ResponseWriter, r *http.Request) {}

// user sign in stuff...
func signin(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		username := r.FormValue("username")
		password := r.FormValue("password")
		//role:=r.FormValue("role")
		m := make(map[string]interface{})
		m["username"] = username
		m["password"] = password
		usr, _ := session.FindByQuery("users", m)
		if usr != nil {
			token := GetSessionID()
			for v, k := range Tokens {
				if v == username {
					delete(Tokens, k)
					break
				}
			}
			Tokens[token] = username

			w.Header().Add("Set-Cookie", "token="+token+"; httponly; Path=/")
			fmt.Fprintln(w, "successfully signed in")
		} else {
			w.WriteHeader(400)
			fmt.Fprintln(w, "user does not exists")
		}
	}
}

func signout(w http.ResponseWriter, r *http.Request) {
	if r.Method == "Delete" {
		token := r.Header.Get("token")
		_, k := Tokens[token]
		if k {
			delete(Tokens, token)
			fmt.Fprintln(w, "successfully signed out")
		} else {
			fmt.Fprintln(w, "invalid token")
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

func HasSignIn(r *http.Request) bool {
	token := r.Header.Get("token")
	if token != "" {

		if _, k := Tokens[token]; k == true {
			return true
		}

	}
	return false
}

//To get a GUID based session ID
func GetSessionID() string {
	// generate 32 bits timestamp
	unix32bits := uint32(time.Now().UTC().Unix())

	buff := make([]byte, 12)

	rand.Read(buff)
	return fmt.Sprintf("%x-%x-%x-%x-%x-%x", unix32bits, buff[0:2], buff[2:4], buff[4:6], buff[6:8], buff[8:])
}
