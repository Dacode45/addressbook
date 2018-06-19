package storage

import mgo "gopkg.in/mgo.v2"

// MongoSession wraps the hasle around dialing and database in an easy object
type MongoSession struct {
	session *mgo.Session
}

// NewMongoSession creates a new session. Will fail if the database is unreachable
func NewMongoSession(url string) (*MongoSession, error) {
	session, err := mgo.Dial(url)
	if err != nil {
		return nil, err
	}
	return &MongoSession{session}, err
}

// Copy copies a session. Same as the Session interface.
func (s *MongoSession) Copy() *MongoSession {
	return &MongoSession{s.session.Copy()}
}

// GetCollection retrieves the mongodby collection
func (s *MongoSession) GetCollection(dbName string, col string) *mgo.Collection {
	return s.session.DB(dbName).C(col)
}

// Close closes a session
func (s *MongoSession) Close() {
	if s.session != nil {
		s.session.Close()
	}
}

// DropDatabase drops the database
func (s *MongoSession) DropDatabase(db string) error {
	if s.session != nil {
		return s.session.DB(db).DropDatabase()
	}
	return nil
}
