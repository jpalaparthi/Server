// service-provider project main.go
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

func main() {
	fmt.Println("Hello World!")
}
