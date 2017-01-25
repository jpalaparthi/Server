package models

import "gopkg.in/mgo.v2/bson"

type Log struct {
	ID          bson.ObjectId `json:"id" bson:"_id,omitempty"`
	User        string        `json:"user" bson:"user"`
	URI         string        `json:"uri" bson:"uri"`
	IP          string        `json:"ip" bson:"ip"`
	Input       string        `json:"input" bson:"input"`
	Output      string        `json:"output" bson:"output"`
	Result_Type string        `json:"result_type" bson:"result_type"`
	Device      string        `json:"device" bson:"device"`
	Device_Name string        `json:"device_name" bson:"device_name"`
	Status      string        `json:"status" bson:"status"`
	TimeStamp   string        `json:"timestamp" bson:"timestamp"`
}
