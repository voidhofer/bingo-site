package model

import (
	// local imports
	"github.com/voidhofer/bingo_site/app/shared/database"
	// external imports
	"gopkg.in/mgo.v2/bson"
)

// Client struct is used for storing users (non administrative user, client only)
type Client struct {
	// client infos
	ObjectID  bson.ObjectId `json:"_id" bson:"_id"`
	FirstName string        `json:"fname" bson:"fname"`
	LastName  string        `json:"lname" bson:"lname"`
	Email     string        `json:"email" bson:"email"`
	Password  string        `json:"password" bson:"password"`
	Status    int           `json:"status" bson:"status"`
}

// ClientID returns the client's ID
func (u *Client) ClientID() string {
	r := ""
	r = u.ObjectID.Hex()
	return r
}

// ClientByEmail returns client with given email address
func ClientByEmail(email string) (Client, error) {
	var err error
	result := Client{}
	if database.CheckConnection() {
		dbsess := database.Mongo.Copy()
		defer dbsess.Close()
		c := dbsess.DB(database.ReadConfig().MongoDB.Database).C("user")
		err = c.Find(bson.M{"email": email}).One(&result)
	} else {
		err = ErrUnavailable
	}
	return result, standardizeError(err)
}

// ClientCreate creates new user record in DB
// temporarly disabled: ', phone, group, preflang, company, taxnum, country, state, city, district, postcode, address, pdtype, house, floor, door, dcompany, dcountry, dstate, dcity, ddistrict, dpostcode, daddress, dpdtype, dhouse, dfloor, ddoor string, promotions bool'
func ClientCreate(firstName, lastName, email, password string) error {
	var err error
	if database.CheckConnection() {
		session := database.Mongo.Copy()
		defer session.Close()
		c := session.DB(database.ReadConfig().MongoDB.Database).C("user")
		client := &Client{
			ObjectID:  bson.NewObjectId(),
			FirstName: firstName,
			LastName:  lastName,
			Email:     email,
			Password:  password,
			Status:    1}
		err = c.Insert(client)
	} else {
		err = ErrUnavailable
	}
	return standardizeError(err)
}
