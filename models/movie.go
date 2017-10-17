package models

import "gopkg.in/mgo.v2/bson"

// service type model
type Movie1 struct {
	ID        bson.ObjectId `json:"id" bson:"_id,omitempty"`
	Title     string        `json:"title" bson:"title"`
	Tags      string        `json:"tags" bson:"tags"`
	Lang      string        `json:"lang" bson:"lang"`
	Banner    string        `json:"banner" bson:"banner"`
	ReleaseOn string        `json:"releaseon" bson:"releaseon"`
	Pics      []string      `json:"pics" bson:"pics"`
	Status    string        `json:"status" bson:"status"`
	Timestamp string        `json:"timestamp" bson:"timestamp"`
}

type Movie struct {
	ID      bson.ObjectId `json:"id" bson:"_id,omitempty"`
	BatchNo string        `json:"batchno" bson:"batchno"`
	Title   string        `json:"title" bson:"title"`
	MovieId string        `json:"movieid" bson:"movieid"`
	Wiki    string        `json:"wiki" bson:"wiki"`
}

type Pic struct {
	ID        bson.ObjectId `json:"id" bson:"_id,omitempty"`
	PicId     string        `json:"picid" bson:"picid"`
	MovieId   string        `json:"movieid" bson:"movieid"`
	Title     string        `json:"title" bson:"title"`
	FileType  string        `json:"filetype" bson:"filetype"`
	PicPath   string        `json:"picpath" bson:"picpath"`
	PicStatus string        `json:"picstatus" bson:"picstatus"`
	Relevance string        `json:"relevence" bson:"relevence"`
	Option    string        `json:"option" bson:"option"`
	Status    string        `json:"status" bson:"status"`
	Timestamp string        `json:"timestamp" bson:"timestamp"`
}

func ValidatePic(P Pic) string {
	if P.MovieId == "" {
		return "Movie Id field is empty"
	}
	if P.Title == "" {
		return "Title field is empty"
	}
	if P.Status == "" {
		return "Status field is empty"
	}

	if P.Timestamp == "" {
		return "Timestamp is empty"
	}
	return ""
}

// Validating each and every field of Service_Type object.
// Any additional validations , can be developed here..
func ValidateMovie(M Movie) string {
	if M.Title == "" {
		return "Title field is empty"
	}
	if M.BatchNo == "" {
		return "BatchNo field is empty"
	}
	/*if M.Lang == "" {
		return "Language field is empty"
	}
	if M.ReleaseOn == "" {
		return "Release On field is empty"
	}

	if M.Tags == "" {
		return "Tags field is empty"
	}

	if M.Status == "" {
		return "Status field is empty"
	}

	if M.Timestamp == "" {
		return "Timestamp is empty"
	}*/
	return ""
}
