package storage

import mgo "gopkg.in/mgo.v2"

type MongoSession struct {
	session *mgo.Session
}

func NewMongoSession(url string) (*MongoSession, error) {
	session, err := mgo.Dial(url)
	if err != nil {
		return nil, err
	}
	return &MongoSession{session}, err
}

func (s *MongoSession) Copy() *MongoSession {
	return &MongoSession{s.session.Copy()}
}

func (s *MongoSession) GetCollection(dbName string, col string) *mgo.Collection {
	return s.session.DB(dbName).C(col)
}

func (s *MongoSession) Close() {
	if s.session != nil {
		s.session.Close()
	}
}

func (s *MongoSession) DropDatabase(db string) error {
	if s.session != nil {
		return s.session.DB(db).DropDatabase()
	}
	return nil
}
