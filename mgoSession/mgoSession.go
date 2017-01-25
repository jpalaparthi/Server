// mgoSession project mgoSession.go
package mgoSession

import (
	"log"
	_ "time"

	_ "models"

	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

type Session struct {
	Con    string
	Sess   *mgo.Session
	DBName string
}

func New(conn string, db string) (dbs *Session, e error) {
	dbs = &Session{}
	log.Println("session created")
	dbs.Sess, e = mgo.Dial(conn)

	if e != nil {
		return nil, e
	} else {
		dbs.DBName = db
		//defer dbs.Sess.Close()
	}
	return dbs, nil
}

//start mongo session

func (s *Session) Insert(col string, data interface{}) error {
	err := s.Sess.DB(s.DBName).C(col).Insert(data)
	if err != nil {
		return err
	} else {
		return nil
	}

}

func (s *Session) UpdateByID(col string, id interface{}, data interface{}) error {
	err := s.Sess.DB(s.DBName).C(col).UpdateId(id, data)
	if err != nil {
		return err
	} else {
		return nil
	}
}

func (s *Session) DeleteByID(col string, id string) error {

	_id := bson.ObjectIdHex(id)
	err := s.Sess.DB(s.DBName).C(col).RemoveId(_id)
	if err != nil {
		return err
	} else {
		return nil
	}
}

func (s *Session) DeleteAll(col string) error {
	_, err := s.Sess.DB(s.DBName).C(col).RemoveAll(nil)
	if err != nil {
		return err
	} else {
		return nil
	}

}

func (s *Session) ListByID(col string, id string) (map[string]interface{}, error) {
	_id := bson.ObjectIdHex(id)
	var result map[string]interface{}
	err := s.Sess.DB(s.DBName).C(col).FindId(_id).One(&result)
	if err != nil {
		return result, err
	} else {
		return result, nil
	}
	return result, nil
}

func (s *Session) FindByQuery(col string, query map[string]interface{}) (map[string]interface{}, error) {
	var result map[string]interface{}
	s.Sess.DB(s.DBName).C((col))
	err := s.Sess.DB(s.DBName).C(col).Find(bson.M(query)).One(&result)
	if err != nil {
		return result, err
	} else {
		return result, nil
	}
	return result, nil
}

func (s *Session) ListAll(col string) ([]interface{}, error) {
	var result []interface{}
	err := s.Sess.DB(s.DBName).C(col).Find(nil).All(&result)
	if err != nil {
		log.Println(err)
		return result, err
	} else {
		return result, nil
	}
	return result, nil
}
