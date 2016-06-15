package main

import (
	"./mgodb"
	"crypto/sha1"
	"encoding/json"
	valid "github.com/asaskevich/govalidator"
	// "io/ioutil"
	"fmt"
	"log"
	"net/http"
	// "net/http/httputil"
)

// Data model: https://github.com/thinkhy/thrasher/wiki/Data-Model
// Rest API: https://github.com/thinkhy/thrasher/wiki/RESTAPI:-user-join
/**************************************************************************
Scheme
-------
{
  _id,
  name:  string,   // *required
  password: string,// *required
  type:  string,   // *required
  email: string,   // *required
  team:  string,
  tso_user: string,
  tso_password: string,
}

Example
--------
{
  _id:   111,
  name:  "Mike",
  password: "pass",
  type:  "admin",
  email: "mike@email.com",
  team:  "Unix Test",
  tso_user: "mike",
  tso_password: "mike",
}

***********************************************************/

/**********************************************************
RestAPI
=========

Path
-------
/api/join

Method
-------
post

Request Parameters
--------------------
{
   action: "join",
   name: "mike",
   email: "xx@yy.com",
   password: "xxxx"
}

Response Data
--------------
{
   ret: "0001",
   info: "ok",

   name: "user.name",
   email: "xx@yy.com",
   sha1: "user.sha1",
   type: "user.type",
   create_time: "2012-04-23T18:25:43.511Z"
   // ISO 8601 format，Refer to: http://stackoverflow.com/questions/10286204/the-right-json-date-format
}

**************************************************************/

type JoinRequestData struct {
	Action   string `json:action`
	Name     string `json:name`
	Email    string `json:email`
	Password string `json:password`
}

func handleJoin(w http.ResponseWriter, r *http.Request) *apiError {
	log.Printf("[user.handleJoin] handleJoin Entry %s\n", r.URL)

	// Read request JSON data
	var d JoinRequestData
	dec := json.NewDecoder(r.Body)
	err := dec.Decode(&d)
	if err != nil {
		str := fmt.Sprintf("failed to parse JSON %v", err)
		log.Println("[user.handleJoin] ", str)
		w.WriteHeader(http.StatusBadRequest)
		return &apiError{Error: err, Ret: -1, Info: str}
	}

	// Validate data format
	if d.Action != "join" {
		w.WriteHeader(http.StatusBadRequest)
		return &apiError{Error: err, Ret: -2, Info: `action should be "join"`}
	} else if valid.Matches(d.Name, `(?i)^[a-z\d.]{5,}$`) == false {
		info := `user name should match the pattern "/^[a-z\d.]{5,}$/i"`
		w.WriteHeader(http.StatusBadRequest)
		return &apiError{Error: err, Ret: -2, Info: info}
	} else if valid.IsEmail(d.Email) == false {
		info := `email format is invalid`
		w.WriteHeader(http.StatusBadRequest)
		return &apiError{Error: err, Ret: -2, Info: info}
	} else if valid.Matches(d.Password, `^.{6,}$`) == false {
		info := `password are case sensitive and must contain a minimum of 6 English characters`
		w.WriteHeader(http.StatusBadRequest)
		return &apiError{Error: err, Ret: -2, Info: info}
	}

	// Insert user data into MongoDB
	// Create SHA1 hash for the original password
	hash := sha1.Sum([]byte(d.Password))
	log.Printf("Hash: %v\n", hash)
	err = mgodb.InsertUser(d.Name,
		d.Email,
		hash[:],
		mgodb.UserType_Tester)
	if err != nil {
		log.Println("[user.handleJoin] ", err)
		w.WriteHeader(http.StatusBadRequest)
		// If can't insert user info into DB, we just imply there are duplicate name or email
		return &apiError{Error: err, Ret: -3, Info: "duplicate name or email exists"}
	}

	// Post-processing

	// Write back user data except for password
	p, err := mgodb.GetUserByName(d.Name)
	if err != nil {
		log.Println("[user handleJoin] Failed to get user data: ", d.Name)
		w.WriteHeader(http.StatusInternalServerError)
		return &apiError{Error: err, Ret: -4, Info: err.Error()}
	}
	p.Ret = 0
	p.Info = "OK"

	w.WriteHeader(http.StatusCreated)

	if err = json.NewEncoder(w).Encode(p); err != nil {
		log.Println("[user handleJoin] Failed to encode JSON data")
		w.WriteHeader(http.StatusInternalServerError)
		return &apiError{Error: err, Ret: -5, Info: err.Error()}
	}

	opshdlr.Write("User", d.Name, "login", "", "")
	log.Printf("[user handleJoin] %s %s joins successfully\n", d.Name, d.Email)
	return nil
}
