package main

// RestAPI doc: https://github.com/thinkhy/thrasher/wiki/RESTAPI-example-for-Test-Workload
// [ 2016-05-17 thinkhy ]

import (
	"./mgodb"
	"fmt"
	"github.com/ant0ine/go-json-rest/rest"
	valid "github.com/asaskevich/govalidator"
	"log"
	"net/http"
	"regexp"
	"strconv"
	"time"
)

type WkRequestData struct {
	Action   string         `json:"action"               valid:"action,required"`
	User     string         `json:"user,omitempty"       valid:"username,required"`
	Name     string         `json:"name,omitempty"       valid:"optional,normalname"`
	Workload mgodb.Workload `json:"workload,omitempty"   valid:"optional"`
	Startpos int            `json:"startpos,omitempty"   valid:"optional"`
	Counter  int            `json:"counter,omitempty"    valid:"optional"`
}

type WkResponseData struct {
	*mgodb.Workload
	Workloads []mgodb.WkBasicInfo `json:"workloads,omitempty"`
	Total     int                 `json:"total,omitempty"`
	Ret       int                 `json:"ret"`
	Info      string              `json:"info"`
}

func init() {
	valid.TagMap["action"] = valid.Validator(func(action string) bool {
		if m, _ := regexp.MatchString(`^(insert|get|update|delete)_workloads?$`, action); !m {
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

func GetWorkload(w rest.ResponseWriter, r *rest.Request) {
	module := "workload"
	function := "GetWorkload"
	log.Printf("[%s.%s] Entry", module, function)

	user := r.URL.Query()["user"]
	tsname := r.URL.Query()["name"]

	d := WkRequestData{}
	var err error
	if len(user) <= 0 || len(tsname) <= 0 {
		err = r.DecodeJsonPayload(&d)
		if err != nil {
			str := fmt.Sprintf("failed to parse JSON %v", err)
			log.Printf("[%s.%s] %s\n", module, function, str)
			WriteErrorInfo(w, &apiError{Status: http.StatusBadRequest, Ret: -1, Info: str})
			return
		}
	} else {
		d.Action = "get_workload"
		d.User = user[0]
		d.Name = tsname[0]
	}

	// Validate data format
	if d.Action != "get_workload" {
		WriteErrorInfo(w, &apiError{Status: http.StatusBadRequest, Ret: -2, Info: `action should be "get_workload"`})
		return
	}

	_, err = valid.ValidateStruct(d)
	if err != nil {
		info := err.Error()
		log.Printf("[%s.%s] failed to validate input data, %s", module, function, info)
		WriteErrorInfo(w, &apiError{Status: http.StatusBadRequest, Ret: -2, Info: info})
		return
	}

	wd := WkResponseData{}
	wd.Workload, err = mgodb.GetWorkloadByName(d.Name)
	if err != nil {
		info := fmt.Sprintf("failed to get workload %s: %v", d.Workload.Name, err)
		log.Printf("[%s.%s] %s\n", module, function, info)
		WriteErrorInfo(w, &apiError{Status: http.StatusBadRequest, Ret: -3, Info: info})
		return
	}

	wd.Ret = 0
	wd.Info = "OK"

	// have to write header prior to writting body of http response
	w.WriteHeader(http.StatusOK)

	err = w.WriteJson(&wd)
	if err != nil {
		log.Printf("[%s.%s] Failed to encode JSON data", module, function)
		WriteErrorInfo(w, &apiError{Status: http.StatusInternalServerError, Ret: -4, Info: err.Error()})
		return
	}

	log.Printf("[%s.%s] %s got workload %s successfully\n", module, function, d.User, d.Name)
	return

}

func GetWorkloads(w rest.ResponseWriter, r *rest.Request) {
	module := "workload"
	function := "GetWorkloads"
	log.Printf("[%s.%s] Entry", module, function)

	d := WkRequestData{}
	var err error

	user := r.URL.Query()["user"]
	startpos := r.URL.Query()["startpos"]
	counter := r.URL.Query()["counter"]

	if len(user) <= 0 || len(startpos) <= 0 || len(counter) <= 0 {
		err = r.DecodeJsonPayload(&d)
		if err != nil {
			str := fmt.Sprintf("failed to parse JSON %v", err)
			log.Printf("[%s.%s] %s\n", module, function, str)
			WriteErrorInfo(w, &apiError{Status: http.StatusBadRequest, Ret: -1, Info: str})
			return
		}
	} else {
		d.Action = "get_workloads"
		d.User = user[0]
		d.Startpos, _ = strconv.Atoi(startpos[0])
		d.Counter, _ = strconv.Atoi(counter[0])
	}

	// Validate data format
	if d.Action != "get_workloads" {
		WriteErrorInfo(w, &apiError{Status: http.StatusBadRequest, Ret: -2, Info: `action should be "get_workloads"`})
		return
	}

	_, err = valid.ValidateStruct(d)
	if err != nil {
		info := err.Error()
		log.Printf("[%s.%s] failed to validate input data, %s", module, function, info)
		WriteErrorInfo(w, &apiError{Status: http.StatusBadRequest, Ret: -2, Info: info})
		return
	} else if d.Startpos <= 0 || d.Counter <= 0 {
		info := "startpos and counter should be positive number"
		log.Printf("[%s.%s] failed to validate input data, %s\n", module, function, info)
		WriteErrorInfo(w, &apiError{Status: http.StatusBadRequest, Ret: -2, Info: info})
		return
	} else if d.Counter > 1000 {
		info := "counter should be less than 1000 or euqal with 1000"
		log.Printf("[%s.%s] failed to validate input data, %s\n", module, function, info)
		WriteErrorInfo(w, &apiError{Status: http.StatusBadRequest, Ret: -2, Info: info})
		return
	}

	wd := WkResponseData{}
	wd.Total, err = mgodb.GetTotalNumber()
	if err != nil {
		info := "failed to get total number of workload collection"
		log.Printf("[%s.%s] %s: %v\n", module, function, info, err)
		WriteErrorInfo(w, &apiError{Status: http.StatusBadRequest, Ret: -3, Info: info})
		return
	}

	wd.Workloads, err = mgodb.GetWorkloads(d.Startpos, d.Counter)
	if err != nil {
		info := fmt.Sprintf("failed to get workloads with startpos %v and counter %v", d.Startpos, d.Counter)
		log.Printf("[%s.%s] %s\n", module, function, info)
		WriteErrorInfo(w, &apiError{Status: http.StatusBadRequest, Ret: -3, Info: info})
		return
	}

	wd.Ret = 0
	wd.Info = "OK"

	// have to write header prior to writting body of http response
	w.WriteHeader(http.StatusOK)

	err = w.WriteJson(&wd)
	if err != nil {
		log.Printf("[%s.%s] Failed to encode JSON data", module, function)
		WriteErrorInfo(w, &apiError{Status: http.StatusInternalServerError, Ret: -4, Info: err.Error()})
		return
	}

	log.Printf("[%s.%s] %s got workloads with startpos %v and counter %v successfully\n", module, function, d.User, d.Startpos, d.Counter)
	return
}

func InsertWorkload(w rest.ResponseWriter, r *rest.Request) {
	log.Printf("[workload.InsertWorkload] Entry")
	d := WkRequestData{}
	err := r.DecodeJsonPayload(&d)
	if err != nil {
		str := fmt.Sprintf("failed to parse JSON %v", err)
		log.Println("[workload.InserWorkload] ", str)
		WriteErrorInfo(w, &apiError{Status: http.StatusBadRequest, Ret: -1, Info: str})
		return
	}

	// Validate data format
	if d.Action != "insert_workload" {
		WriteErrorInfo(w, &apiError{Status: http.StatusBadRequest, Ret: -2, Info: `action should be "insert_workload"`})
		return
	}

	// Fill appropriate values for Creator, Modifier, CreateTime and UpdateTime
	d.Workload.Creator, d.Workload.Modifier = d.User, d.User
	d.Workload.CreateTime, d.Workload.UpdateTime = time.Now(), time.Now()

	_, err = valid.ValidateStruct(d)
	if err != nil {
		info := err.Error()
		log.Println("[workload.InsertWorkload] failed to validate input data, ", info)
		WriteErrorInfo(w, &apiError{Status: http.StatusBadRequest, Ret: -2, Info: info})
		return
	}

	// Insert workload data into MongoDB
	err = mgodb.InsertWorkload(&d.Workload)
	if err != nil {
		log.Println("[workload.InsertWorkload] ", err)
		// If can't insert workload info into DB, we just imply there are duplicate name
		WriteErrorInfo(w, &apiError{Status: http.StatusBadRequest, Ret: -3, Info: "duplicate workload name exists"})
		return
	}

	wd := WkResponseData{}
	wd.Ret = 0
	wd.Info = "OK"

	wd.Workload, err = mgodb.GetWorkloadByName(d.Workload.Name)
	if err != nil {
		info := fmt.Sprintf("failed to get workload %s: %s", d.Workload.Name, err)
		log.Println("[workload.InsertWorkload] ", info)
		WriteErrorInfo(w, &apiError{Status: http.StatusBadRequest, Ret: -3, Info: info})
		return
	}

	w.WriteHeader(http.StatusCreated)
	err = w.WriteJson(&wd)
	if err != nil {
		log.Println("[workload.InsertWorkload] Failed to encode JSON data")
		WriteErrorInfo(w, &apiError{Status: http.StatusInternalServerError, Ret: -4, Info: err.Error()})
		return
	}

	opshdlr.Write("Workload", d.User, "insert", d.Workload.Name, "")
	log.Printf("[workload] %s inserted workload %s successfully\n", d.User, d.Workload.Name)
	return
}

func UpdateWorkload(w rest.ResponseWriter, r *rest.Request) {
	module := "workload"
	function := "UpdateWorkload"
	action := "update_workload"
	log.Printf("[%s.%s] Entry", module, function)

	d := WkRequestData{}
	err := r.DecodeJsonPayload(&d)
	if err != nil {
		str := fmt.Sprintf("failed to parse JSON %v", err)
		log.Printf("[%s.%s] %s\n", module, function, str)
		WriteErrorInfo(w, &apiError{Status: http.StatusBadRequest, Ret: -1, Info: str})
		return
	}

	// Validate data format
	if d.Action != action {
		WriteErrorInfo(w, &apiError{Status: http.StatusBadRequest, Ret: -2, Info: `action should be "update_workload"`})
		return
	}

	// Fill appropriate values for Modifier and UpdateTime
	d.Workload.Modifier = d.User
	d.Workload.UpdateTime = time.Now()

	_, err = valid.ValidateStruct(d)
	if err != nil {
		info := err.Error()
		log.Printf("[%s.%s] failed to validate input data, %s", module, function, info)
		WriteErrorInfo(w, &apiError{Status: http.StatusBadRequest, Ret: -2, Info: info})
		return
	}

	// clear values of Name, Creator, CreateTime which can ONLY be generated when insert workload
	// clear after validating because these values are required to be validated as tags indicate
	d.Workload.Creator = ""
	d.Workload.CreateTime = time.Time{}

	err = mgodb.UpdateWorkload(&d.Workload)
	if err != nil {
		info := fmt.Sprint("failed to get workload ", d.Workload.Name)
		log.Printf("[%s.%s] %v\n", module, function, info)
		WriteErrorInfo(w, &apiError{Status: http.StatusBadRequest, Ret: -3, Info: info})
	}

	wd := WkResponseData{}
	wd.Ret = 0
	wd.Info = "OK"

	wd.Workload, err = mgodb.GetWorkloadByName(d.Workload.Name)
	if err != nil {
		info := fmt.Sprintf("failed to get workload %s: %s", d.Workload.Name, err)
		log.Println("[workload.handleUpdate] ", info)
		WriteErrorInfo(w, &apiError{Status: http.StatusBadRequest, Ret: -3, Info: info})
		return
	}

	// have to write header prior to writting body of http response
	w.WriteHeader(http.StatusOK)

	err = w.WriteJson(&wd)
	if err != nil {
		log.Printf("[%s.%s] Failed to encode JSON data", module, function)
		WriteErrorInfo(w, &apiError{Status: http.StatusInternalServerError, Ret: -4, Info: err.Error()})
		return
	}

	opshdlr.Write("Workload", d.User, "update", d.Workload.Name, "")
	log.Printf("[%s.%s] %s updated workload %s successfully\n", module, function, d.User, wd.Workload.Name)
	return
}

func DeleteWorkload(w rest.ResponseWriter, r *rest.Request) {
	module := "workload"
	function := "DeleteWorkload"
	action := "delete_workload"
	log.Printf("[%s.%s] Entry", module, function)

	d := WkRequestData{}
	err := r.DecodeJsonPayload(&d)
	if err != nil {
		str := fmt.Sprintf("failed to parse JSON %v", err)
		log.Printf("[%s.%s] %s\n", module, function, str)
		WriteErrorInfo(w, &apiError{Status: http.StatusBadRequest, Ret: -1, Info: str})
		return
	}

	// Validate data format
	if d.Action != action {
		WriteErrorInfo(w, &apiError{Status: http.StatusBadRequest, Ret: -2, Info: `action should be "delete_workload"`})
		return
	}

	_, err = valid.ValidateStruct(d)
	if err != nil {
		info := err.Error()
		log.Printf("[%s.%s] failed to validate input data, %s", module, function, info)
		WriteErrorInfo(w, &apiError{Status: http.StatusBadRequest, Ret: -2, Info: info})
		return
	}

	g, err := mgodb.GetWorkloadByName(d.Name)
	if err != nil {
		info := fmt.Sprint("failed to get workload ", d.Name)
		log.Printf("[%s.%s] %v\n", module, function, info)
		WriteErrorInfo(w, &apiError{Status: http.StatusBadRequest, Ret: -3, Info: info})
	} else if g.Creator != d.User {
		info := fmt.Sprintf("failed to get workload %s beacause %s is not creator of workload", d.Name, d.User)
		log.Printf("[%s.%s] %v\n", module, function, info)
		WriteErrorInfo(w, &apiError{Status: http.StatusBadRequest, Ret: -3, Info: info})
	}

	err = mgodb.DeleteWorkload(d.Name)
	if err != nil {
		info := fmt.Sprintf("failed to delete workload %s: %v", d.Name, err)
		log.Printf("[%s.%s] %s\n", module, function, info)
		WriteErrorInfo(w, &apiError{Status: http.StatusBadRequest, Ret: -3, Info: info})
		return
	}

	wd := WkResponseData{}
	wd.Ret = 0
	wd.Info = "OK"

	// have to write header prior to writting body of http response
	w.WriteHeader(http.StatusOK)

	err = w.WriteJson(&wd)
	if err != nil {
		log.Printf("[%s.%s] Failed to encode JSON data", module, function)
		WriteErrorInfo(w, &apiError{Status: http.StatusInternalServerError, Ret: -4, Info: err.Error()})
		return
	}

	opshdlr.Write("Workload", d.User, "delete", d.Name, "")
	log.Printf("[%s.%s] %s deleted workload %s successfully\n", module, function, d.User, d.Name)

	return
}
