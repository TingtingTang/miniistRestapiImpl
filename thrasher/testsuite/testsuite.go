package main

import (
	"./mgodb"
	"encoding/json"
	"fmt"
	valid "github.com/asaskevich/govalidator"
	"log"
	"net/http"
	"regexp"
	"strconv"
	"time"
)

type TsRequestData struct {
	Action   string          `json:"action"               valid:"action,required"`
	User     string          `json:"user,omitempty"       valid:"username,required"`
	Name     string          `json:"name,omitempty"       valid:"normalname"`
	Ts       mgodb.TestSuite `json:"testsuite,omitempty"  valid:"optional"`
	Startpos int             `json:"startpos,omitempty"   valid:"optional"`
	Counter  int             `json:"counter,omitempty"    valid:"optional"`
}

type TsResponseData struct {
	*mgodb.TestSuite
	Testsuites []mgodb.TcBasicInfo `json:"testsuites,omitempty"`
	Total      int                 `json:"total,omitempty"`
	Ret        int                 `json:"ret"`
	Info       string              `json:"info"`
}

func init() {
	valid.TagMap["action"] = valid.Validator(func(action string) bool {
		if m, _ := regexp.MatchString(`^(insert|get|update|delete)_testsuites?$`, action); !m {
			return false
		} else {
			return true
		}
	})

	valid.TagMap["username"] = valid.Validator(func(username string) bool {
		if m, _ := regexp.MatchString(`^[A-Za-z\d.]{5,}$`, username); !m {
			return false
		} else {
			return true
		}
	})

	valid.TagMap["normalname"] = valid.Validator(func(normalname string) bool {
		// if m, _ := regexp.MatchString(`^([\u00c0-\u01ffa-zA-Z'\-])+$`, normalname); !m {
		if m, _ := regexp.MatchString(`^.{3,}$`, normalname); !m {
			return false
		} else {
			return true
		}
	})
}

func handleInsert(w http.ResponseWriter, r *http.Request) *apiError {
	log.Printf("[testsuite] handleInsert Entry %s\n", r.URL)

	// Read JSON data from request
	var d TsRequestData
	dec := json.NewDecoder(r.Body)
	err := dec.Decode(&d)
	if err != nil {
		str := fmt.Sprintf("failed to parse JSON %v", err)
		log.Println("[testsuite.handleInsert] ", str)
		w.WriteHeader(http.StatusBadRequest)
		return &apiError{Error: err, Ret: -1, Info: str}
	}

	// Validate data format
	if d.Action != "insert_testsuite" {
		w.WriteHeader(http.StatusBadRequest)
		return &apiError{Error: err, Ret: -2, Info: `action should be "insert_testsuite"`}
	}

	// Fill appropriate values for Creator, Modifier, CreateTime and UpdateTime
	d.Ts.Creator, d.Ts.Modifier = d.User, d.User
	d.Ts.CreateTime, d.Ts.UpdateTime = time.Now(), time.Now()

	_, err = valid.ValidateStruct(d)
	if err != nil {
		info := err.Error()
		log.Println("[testsuite.handleInsert] failed to validate input data, ", info)
		w.WriteHeader(http.StatusBadRequest)
		return &apiError{Error: err, Ret: -2, Info: info}
	}

	// Authenticate through JWT token
	if Config.TestFlag == false {
		ae := withAuthentication(w, r, d.User)
		if ae != nil {
			return ae
		}
	}

	// Insert testsuite data into MongoDB
	err = mgodb.InsertTestsuite(&d.Ts)
	if err != nil {
		log.Println("[testsuite.handleInsert] ", err)
		w.WriteHeader(http.StatusBadRequest)
		// If can't insert testsuite info into DB, we just imply there are duplicate name
		return &apiError{Error: err, Ret: -3, Info: "duplicate testsuite name exists"}
	}

	var rs TsResponseData
	rs.Ret = 0
	rs.Info = "OK"

	rs.TestSuite, err = mgodb.GetTestsuiteByName(d.Ts.Name)
	if err != nil {
		info := fmt.Sprintf("failed to get testsuite %s: %s", d.Name, err)
		log.Println("[testsuite.handleInsert] ", info)
		w.WriteHeader(http.StatusBadRequest)
		return &apiError{Error: err, Ret: -3, Info: info}
	}

	w.WriteHeader(http.StatusCreated)
	if err = json.NewEncoder(w).Encode(rs); err != nil {
		log.Println("[testsuite.handleInsert] Failed to encode JSON data")
		w.WriteHeader(http.StatusInternalServerError)
		return &apiError{Error: err, Ret: -5, Info: err.Error()}
	}

	opshdlr.Write("Testsuite", d.User, "insert", d.Ts.Name, "")
	log.Printf("[testsuite.handleInsert] %s inserted testsuite %s successfully\n", d.User, d.Ts.Name)
	return nil
}

func handleGet(w http.ResponseWriter, r *http.Request) *apiError {
	log.Printf("[testsuite] handleGet Entry %s\n", r.URL)

	user := r.URL.Query()["user"]
	tsname := r.URL.Query()["name"]

	// Read JSON data from request
	var d TsRequestData
	var err error

	if len(user) <= 0 || len(tsname) <= 0 {
		dec := json.NewDecoder(r.Body)
		err = dec.Decode(&d)
		if err != nil {
			str := fmt.Sprintf("failed to parse JSON %v", err)
			log.Println("[testsuite.handleGet] ", str)
			w.WriteHeader(http.StatusBadRequest)
			return &apiError{Error: err, Ret: -1, Info: str}
		}
	} else {
		d.Action = "get_testsuite"
		d.User = user[0]
		d.Name = tsname[0]
	}

	// Validate data format
	if d.Action != "get_testsuite" {
		w.WriteHeader(http.StatusBadRequest)
		return &apiError{Error: err, Ret: -2, Info: `action should be "get_testsuite"`}
	}

	_, err = valid.ValidateStruct(d)
	if err != nil {
		info := err.Error()
		log.Println("[testsuite.handleGet] failed to validate input data, ", info)
		w.WriteHeader(http.StatusBadRequest)
		return &apiError{Error: err, Ret: -2, Info: info}
	}

	// Authenticate through JWT token
	if Config.TestFlag == false {
		ae := withAuthentication(w, r, d.User)
		if ae != nil {
			return ae
		}
	}

	var rs TsResponseData
	rs.TestSuite, err = mgodb.GetTestsuiteByName(d.Name)
	if err != nil {
		info := fmt.Sprintf("failed to get testsuite %s: %s", d.Name, err)
		log.Println("[testsuite.handleGet] ", info)
		w.WriteHeader(http.StatusBadRequest)
		return &apiError{Error: err, Ret: -3, Info: info}
	}

	rs.Ret = 0
	rs.Info = "OK"
	w.WriteHeader(http.StatusOK)

	if err = json.NewEncoder(w).Encode(rs); err != nil {
		log.Println("[testsuite.handleGet] Failed to encode JSON data")
		w.WriteHeader(http.StatusInternalServerError)
		return &apiError{Error: err, Ret: -5, Info: err.Error()}
	}

	log.Printf("[testsuite.handleInsert] %s got testsuite %s successfully\n", d.User, d.Name)
	return nil
}

func handleUpdate(w http.ResponseWriter, r *http.Request) *apiError {
	log.Printf("[testsuite] handleUpdate Entry %s\n", r.URL)

	// Read JSON data from request
	var d TsRequestData
	dec := json.NewDecoder(r.Body)
	err := dec.Decode(&d)
	if err != nil {
		str := fmt.Sprintf("failed to parse JSON %v", err)
		log.Println("[testsuite.handleUpdate] ", str)
		w.WriteHeader(http.StatusBadRequest)
		return &apiError{Error: err, Ret: -1, Info: str}
	}

	// Validate data format
	if d.Action != "update_testsuite" {
		w.WriteHeader(http.StatusBadRequest)
		return &apiError{Error: err, Ret: -2, Info: `action should be "update_testsuite"`}
	}

	// Fill appropriate values for Modifier and UpdateTime
	d.Ts.Modifier = d.User
	d.Ts.UpdateTime = time.Now()

	_, err = valid.ValidateStruct(d)
	if err != nil {
		info := err.Error()
		log.Println("[testsuite.handleUpdate] failed to validate input data, ", info)
		w.WriteHeader(http.StatusBadRequest)
		return &apiError{Error: err, Ret: -2, Info: info}
	}

	// Authenticate through JWT token
	if Config.TestFlag == false {
		ae := withAuthentication(w, r, d.User)
		if ae != nil {
			return ae
		}
	}

	// clear values of Name, Creator, CreateTime which can ONLY be generated when insert testsuite
	// clear after validating because these values are required to be validated as tags indicate
	// d.Ts.Name = ""
	d.Ts.Creator = ""
	d.Ts.CreateTime = time.Time{}

	// Insert testsuite data into MongoDB
	err = mgodb.UpdateTestsuite(d.Ts)
	if err != nil {
		log.Println("[testsuite.handleUpdate] ", err)
		w.WriteHeader(http.StatusBadRequest)
		return &apiError{Error: err, Ret: -3, Info: fmt.Sprint("failed to update testsuite ", d.Ts.Name)}
	}

	var rs TsResponseData
	rs.Ret = 0
	rs.Info = "OK"

	rs.TestSuite, err = mgodb.GetTestsuiteByName(d.Ts.Name)
	if err != nil {
		info := fmt.Sprintf("failed to get testsuite %s: %s", d.Name, err)
		log.Println("[testsuite.handleUpdate] ", info)
		w.WriteHeader(http.StatusBadRequest)
		return &apiError{Error: err, Ret: -3, Info: info}
	}

	w.WriteHeader(http.StatusOK)
	if err = json.NewEncoder(w).Encode(rs); err != nil {
		log.Println("[testsuite.handleUpdate] Failed to encode JSON data")
		w.WriteHeader(http.StatusInternalServerError)
		return &apiError{Error: err, Ret: -5, Info: err.Error()}
	}

	opshdlr.Write("Testsuite", d.User, "update", d.Ts.Name, "")
	log.Printf("[testsuite.handleUpdate] %s updated testsuite %s successfully\n", d.User, d.Ts.Name)
	return nil
}

func handleDelete(w http.ResponseWriter, r *http.Request) *apiError {
	log.Printf("[testsuite] handleDelete Entry %s\n", r.URL)

	// Read JSON data from request
	var d TsRequestData
	dec := json.NewDecoder(r.Body)
	err := dec.Decode(&d)
	if err != nil {
		str := fmt.Sprintf("failed to parse JSON %v", err)
		log.Println("[testsuite.handleDelete] ", str)
		w.WriteHeader(http.StatusBadRequest)
		return &apiError{Error: err, Ret: -1, Info: str}
	}

	// Validate data format
	if d.Action != "delete_testsuite" {
		w.WriteHeader(http.StatusBadRequest)
		return &apiError{Error: err, Ret: -2, Info: `action should be "delete_testsuite"`}
	}

	_, err = valid.ValidateStruct(d)
	if err != nil {
		info := err.Error()
		log.Println("[testsuite.handleDelete] failed to validate input data, ", info)
		w.WriteHeader(http.StatusBadRequest)
		return &apiError{Error: err, Ret: -2, Info: info}
	}

	if Config.TestFlag == false {
		ae := withAuthentication(w, r, d.User)
		if ae != nil {
			return ae
		}
	}

	g, err := mgodb.GetTestsuiteByName(d.Name)
	if err != nil {
		info := fmt.Sprint("failed to get testsuite ", d.Name)
		log.Println("[testsuite.handleDelete] ", info)
		w.WriteHeader(http.StatusBadRequest)
		return &apiError{Error: err, Ret: -3, Info: info}
	} else if g.Creator != d.User {
		info := fmt.Sprintf("failed to get testsuite %s beacause %s is not creator of testsuite", d.Name, d.User)
		log.Println("[testsuite.handleDelete] ", info)
		w.WriteHeader(http.StatusBadRequest)
		return &apiError{Error: err, Ret: -3, Info: info}
	}

	// Insert testsuite data into MongoDB
	err = mgodb.DeleteTestsuite(d.Name)
	if err != nil {
		log.Println("[testsuite.handleDelete] ", err)
		w.WriteHeader(http.StatusBadRequest)
		return &apiError{Error: err, Ret: -3, Info: fmt.Sprintf("failed to delete testsuit %s", d.Name)}
	}

	var rs TsResponseData
	rs.Ret = 0
	rs.Info = "OK"

	w.WriteHeader(http.StatusOK)

	if err = json.NewEncoder(w).Encode(rs); err != nil {
		log.Println("[testsuite.handleDelete] Failed to encode JSON data")
		w.WriteHeader(http.StatusInternalServerError)
		return &apiError{Error: err, Ret: -5, Info: err.Error()}
	}

	opshdlr.Write("Testsuite", d.User, "delete", d.Name, "")
	log.Printf("[testsuite.handleDelete] %s deleted testsuite %s successfully\n", d.User, d.Name)
	return nil
}

// func handleDelete(w http.ResponseWriter, r *http.Request) *apiError {
func handleGetTestsuites(w http.ResponseWriter, r *http.Request) *apiError {
	function := "handleGetTestsuites"
	log.Printf("[testsuite] %s Entry %s\n", function, r.URL)

	// Read JSON data from request
	var d TsRequestData
	var err error

	user := r.URL.Query()["user"]
	startpos := r.URL.Query()["startpos"]
	counter := r.URL.Query()["counter"]

	if len(user) <= 0 || len(startpos) <= 0 || len(counter) <= 0 {
		dec := json.NewDecoder(r.Body)
		err = dec.Decode(&d)
		if err != nil {
			str := fmt.Sprintf("failed to parse JSON %v", err)
			log.Printf("[testsuite.%s] %s\n", function, str)
			w.WriteHeader(http.StatusBadRequest)
			return &apiError{Error: err, Ret: -1, Info: str}
		}
	} else {
		d.Action = "get_testsuites"
		d.User = user[0]
		d.Startpos, _ = strconv.Atoi(startpos[0])
		d.Counter, _ = strconv.Atoi(counter[0])
	}

	// Validate data format
	if d.Action != "get_testsuites" {
		w.WriteHeader(http.StatusBadRequest)
		return &apiError{Error: err, Ret: -2, Info: `action should be "get_testsuites"`}
	}

	_, err = valid.ValidateStruct(d)
	if err != nil {
		info := err.Error()
		log.Printf("[testsuite.%s] failed to validate input data, %s\n", function, info)
		w.WriteHeader(http.StatusBadRequest)
		return &apiError{Error: err, Ret: -2, Info: info}
	} else if d.Startpos <= 0 || d.Counter <= 0 {
		info := "startpos and counter should be positive number"
		log.Printf("[testsuite.%s] failed to validate input data, %s\n", function, info)
		w.WriteHeader(http.StatusBadRequest)
		return &apiError{Error: err, Ret: -2, Info: info}
	} else if d.Counter > 1000 {
		info := "counter should be less than 1000 or euqal with 1000"
		log.Printf("[testsuite.%s] failed to validate input data, %s\n", function, info)
		w.WriteHeader(http.StatusBadRequest)
		return &apiError{Error: err, Ret: -2, Info: info}
	}

	if Config.TestFlag == false {
		ae := withAuthentication(w, r, d.User)
		if ae != nil {
			return ae
		}
	}

	var rs TsResponseData
	rs.Total, err = mgodb.GetTotalNumber()
	if err != nil {
		info := "failed to get total number of testsuite collection"
		log.Printf("[testsuite.handleGetTestsuites] %s: %v\n", info, err)
		w.WriteHeader(http.StatusBadRequest)
		return &apiError{Error: err, Ret: -3, Info: info}
	}

	rs.Testsuites, err = mgodb.GetTestsuites(d.Startpos, d.Counter)
	if err != nil {
		info := fmt.Sprintf("failed to get testsuite with startpos %v and counter %v", d.Startpos, d.Counter)
		log.Printf("[testsuite.%s] ", function, info)
		w.WriteHeader(http.StatusBadRequest)
		return &apiError{Error: err, Ret: -3, Info: info}
	}

	rs.Ret = 0
	rs.Info = "OK"
	w.WriteHeader(http.StatusOK)

	if err = json.NewEncoder(w).Encode(rs); err != nil {
		log.Printf("[testsuite.%s] Failed to encode JSON data", function)
		w.WriteHeader(http.StatusInternalServerError)
		return &apiError{Error: err, Ret: -5, Info: err.Error()}
	}

	log.Printf("[testsuite.%s] %v got testsuites successfully(total: %v, startpos: %v, counter: %v)\n",
		function, d.User, rs.Total, d.Startpos, d.Counter)
	return nil
}
