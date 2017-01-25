package models

import "gopkg.in/mgo.v2/bson"

type ServiceType struct {
	ID        bson.ObjectId `json:"id" bson:"_id,omitempty"`
	Name      string        `json:"name" bson:"name"`
	Desc      string        `json:"desc" bson:"desc"`
	Sub_Types []string      `json:"sub_types" bson:"sub_types"`
	Status    string        `json:"status" bson:"status"`
	TimeStamp string        `json:"timestamp" bson:"timestamp"`
}

func ValidateServiceType(st ServiceType) string {
	if st.Name == "" {
		return "Service Name field is empty"
	}
	if st.Desc == "" {
		return "Description field is empty"
	}

	if st.Status == "" {
		return "Status field is empty"
	}
	if st.TimeStamp == "" {
		return "Timestamp is empty"
	}
	return ""
}

func ValidateServiceTypeForUpdate(st ServiceType) string {
	if st.ID.String() == "" {
		return "id field is empty"
	}
	return ""
}
