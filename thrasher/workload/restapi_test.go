package main

import (
	//assert "github.com/stretchr/testify/assert"
	"./mgodb"
	"encoding/json"
	"fmt"
	. "github.com/stretchr/testify/require"
	"gopkg.in/mgo.v2/bson"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"
	// "regexp"
	"strings"
	"time"
)

var (
	server       *httptest.Server
	create_time  time.Time
	workloadId   bson.ObjectId
	workloadName string
	// creator string
	// reader   io.Reader
	// testsuiteUrl string
	// jwtToken string
	// userID   string
)

func init() {
	loadConfig()
	setupInfluxDB()
	mgodb.Setup(Config.DbAddress)
	api := MakeApi()
	server = httptest.NewServer(api.MakeHandler())
	Config.WorkloadUrl = fmt.Sprintf("%s%s", server.URL, Config.WorkloadUrl)
	Config.WorkloadsUrl = fmt.Sprintf("%s%s", server.URL, Config.WorkloadsUrl)
	Config.TestFlag = true
}

func TestInsert(t *testing.T) {
	var err error
	tsJson := `{ 
                  "action": "insert_workload",
                  "user":   "mike32432487293",
	              "workload": {
				  "name": "workload_test01",
				  "team":   "testteam01",
				  "run_as_user":   "mega",
				  "desc":   "desc01",
				  "category": "category01",
				  "tag":    ["tagA", "tagB"],
				  "setup_script":"setupScript",
				  "clean_script":"cleanScript",
				  "exe_script":  "exeScript",
				  "env_script":  "envScript",
				  "machine": ["sy01", "sy02", "sy03"],
				  "testsuites": [
					{
					 "id": "560d34ce0699616af8b86843",
					 "name": "ts22", 
					 "weight": 100 
					},
					{
					 "id": "560d34ce0699616af8b86842",
					 "name": "ts23", 
					 "weight": 200 
					}
				  ]	
		      }
		  }`
	log.Println("+++ User mike32432487293 creates a new workload workload_test01")
	reader := strings.NewReader(tsJson)
	create_time = time.Now()
	request, err := http.NewRequest("POST", Config.WorkloadUrl, reader)
	Nil(t, err, "Failed to invoke http.NewRequest")
	request.Header.Set("Content-Type", "application/json")
	log.Printf("Post a request at URL %s\n", Config.WorkloadUrl)
	res, err := http.DefaultClient.Do(request)
	Nil(t, err, "Failed to invoke http.DefaultClient.Do")
	Equal(t, http.StatusCreated, res.StatusCode, "")

	var d WkResponseData
	dec := json.NewDecoder(res.Body)
	err = dec.Decode(&d)
	Nil(t, err, "Failed to parse response JSON")
	Equal(t, 0, d.Ret, "Ret should be correct")
	Equal(t, "OK", d.Info, "Info should be correct")

}

func TestGet(t *testing.T) {
	var err error
	tsJson := `{ 
                  "action": "get_workload",
                  "user":   "mike32432487293",
	          "name":   "workload_test01"
		  }`
	log.Println("+++ User mike32432487293 gets a workload workload_test01")
	reader := strings.NewReader(tsJson)
	request, err := http.NewRequest("GET", Config.WorkloadUrl, reader)
	Nil(t, err, "Failed to invoke http.NewRequest")
	request.Header.Set("Content-Type", "application/json")
	log.Printf("Get a request at URL %s\n", Config.WorkloadUrl)
	res, err := http.DefaultClient.Do(request)
	Nil(t, err, "Failed to invoke http.DefaultClient.Do")
	Equal(t, http.StatusOK, res.StatusCode, "")

	var d WkResponseData
	dec := json.NewDecoder(res.Body)
	err = dec.Decode(&d)
	Nil(t, err, "Failed to parse response JSON")
	Equal(t, 0, d.Ret, "Ret should be correct")
	Equal(t, "OK", d.Info, "Info should be correct")
	workloadId = d.Id
	Equal(t, "workload_test01", d.Workload.Name, "Workload.name should be correct")
	Equal(t, "testteam01", d.Workload.Team, "Workload.team should be correct")
	Equal(t, "mega", d.Workload.Run_as_user, "Workload.run_as_user should be correct")
	Equal(t, "desc01", d.Workload.Desc, "Workload.desc should be correct")
	Equal(t, "category01", d.Workload.Category, "Workload.category should be correct")
	Equal(t, "workload_test01", d.Workload.Name, "Workload.name should be correct")
	Equal(t, "workload_test01", d.Workload.Name, "Workload.name should be correct")
	Equal(t, "workload_test01", d.Workload.Name, "Workload.name should be correct")
	tag := []string{"tagA", "tagB"}
	Equal(t, tag, d.Workload.Tag, "Workload.tag should be correct")
	Condition(t,
		func() (success bool) {
			return d.CreateTime == d.UpdateTime && d.CreateTime.Sub(create_time) < 1*time.Second
		},
		"Createtime shuold be within 2 seconds")
	Equal(t, "setupScript", d.SetupScript, "Setup script should be correct")
	Equal(t, "cleanScript", d.CleanScript, "Clean script should be correct")
	Equal(t, "exeScript", d.ExeScript, "Execution script should be correct")
	Equal(t, "envScript", d.EnvScript, "Env script should be correct")
	machines := []string{"sy01", "sy02", "sy03"}
	Equal(t, machines, d.Machines, "Machines should be correct")
	testsuites := []mgodb.Testsuite{
		{Id: bson.ObjectIdHex("560d34ce0699616af8b86843"),
			Name:   "ts22",
			Weight: 100,
		},
		{Id: bson.ObjectIdHex("560d34ce0699616af8b86842"),
			Name:   "ts23",
			Weight: 200,
		},
	}
	Equal(t, testsuites, d.Testsuites, "Testsuites should be correct")

}

func TestUpdate(t *testing.T) {
	var err error
	tsJson := fmt.Sprintf(`{ 
                  "action": "update_workload",
                  "user":   "mike32432487293",
	              "workload": {
			      "id":"%v",
				  "name": "workload_test02",
				  "team":   "testteam02",
				  "run_as_user":   "wellie",
				  "desc":   "desc02",
				  "category":   "category02",
				  "tag":    ["tagAA", "tagBB"],
				  "setup_script":"setupScript01",
				  "clean_script":"cleanScript01",
				  "exe_script":  "exeScript01",
				  "env_script":  "envScript01",
				  "machine": ["sy01A", "sy02A", "sy03A"],
				  "testsuites": [
					{
					 "id": "560d34ce0699616af8b86843",
					 "name": "ts22", 
					 "weight": 200 
					},
					{
					 "id": "560d34ce0699616af8b86842",
					 "name": "ts23", 
					 "weight": 300 
					}
				  ]	
		      }
		  }`, workloadId.Hex())
	method := "PUT"
	url := Config.WorkloadUrl

	log.Println("+++ User mike32432487293 updates a workload workload_test02")
	reader := strings.NewReader(tsJson)
	create_time = time.Now()
	request, err := http.NewRequest(method, url, reader)
	Nil(t, err, "Failed to invoke http.NewRequest")
	request.Header.Set("Content-Type", "application/json")
	log.Printf("%s a request at URL %s\n", method, url)
	res, err := http.DefaultClient.Do(request)
	Nil(t, err, "Failed to invoke http.DefaultClient.Do")
	Equal(t, http.StatusOK, res.StatusCode, "")

	var d WkResponseData
	dec := json.NewDecoder(res.Body)
	err = dec.Decode(&d)
	Nil(t, err, "Failed to parse response JSON")
	Equal(t, 0, d.Ret, "Ret should be correct")
	Equal(t, "OK", d.Info, "Info should be correct")

	Equal(t, "workload_test02", d.Workload.Name, "Workload.name should be correct")
	Equal(t, "testteam02", d.Workload.Team, "Workload.team should be correct")
	Equal(t, "wellie", d.Workload.Run_as_user, "Workload.run_as_user should be correct")
	Equal(t, "desc02", d.Workload.Desc, "Workload.desc should be correct")
	Equal(t, "category02", d.Workload.Category, "Workload.category should be correct")
	tag := []string{"tagAA", "tagBB"}
	Equal(t, tag, d.Workload.Tag, "Workload.tag should be correct")
	Condition(t,
		func() (success bool) {
			return time.Since(d.UpdateTime) < 1*time.Second
		},
		"Createtime shuold be within 2 seconds")
	Equal(t, "setupScript01", d.SetupScript, "Setup script should be correct")
	Equal(t, "cleanScript01", d.CleanScript, "Clean script should be correct")
	Equal(t, "exeScript01", d.ExeScript, "Execution script should be correct")
	Equal(t, "envScript01", d.EnvScript, "Env script should be correct")
	machines := []string{"sy01A", "sy02A", "sy03A"}
	Equal(t, machines, d.Machines, "Machines should be correct")
	testsuites := []mgodb.Testsuite{
		{Id: bson.ObjectIdHex("560d34ce0699616af8b86843"),
			Name:   "ts22",
			Weight: 200,
		},
		{Id: bson.ObjectIdHex("560d34ce0699616af8b86842"),
			Name:   "ts23",
			Weight: 300,
		},
	}
	Equal(t, testsuites, d.Testsuites, "Testsuites should be correct")
}

func TestGetWorkloads(t *testing.T) {
	tsJson := `{ 
                  "action":"get_workloads",
                  "user":"mike32432487293",
				  "startpos":1,
				  "counter":100
                }`
	method := "GET"
	url := Config.WorkloadsUrl

	log.Println("+++ User mike32432487293 get workloads")
	reader := strings.NewReader(tsJson)
	create_time = time.Now()
	request, err := http.NewRequest(method, url, reader)
	Nil(t, err, "Failed to invoke http.NewRequest")
	request.Header.Set("Content-Type", "application/json")
	log.Printf("%s a request at URL %s\n", method, url)
	res, err := http.DefaultClient.Do(request)
	Nil(t, err, "Failed to invoke http.DefaultClient.Do")

	var d WkResponseData
	dec := json.NewDecoder(res.Body)
	err = dec.Decode(&d)
	Nil(t, err, "Failed to parse response JSON")
	Equal(t, "OK", d.Info, "Info should be correct")
	Equal(t, 0, d.Ret, "Ret should be correct")
	Equal(t, http.StatusOK, res.StatusCode, "")
	True(t, (d.Total >= 1), "Total number should >= 1")

	for _, wk := range d.Workloads {
		log.Println(wk)
		True(t, wk.Id.Valid(), "Workload ID should be valid")
		Equal(t, "workload_test02", wk.Name, "Name should be correct")
		Equal(t, "testteam02", wk.Team, "Team should be correct")
		Equal(t, "desc02", wk.Desc, "Desc should be correct")
		Equal(t, "category02", wk.Category, "Category should be correct")
		Equal(t, "mike32432487293", wk.Creator, "Creator should be correct")
		Equal(t, "mike32432487293", wk.Modifier, "Modifier should be correct")
	}

}

func TestDelete(t *testing.T) {
	var err error
	tsJson := `{ 
		    "action": "delete_workload",
		    "user":   "mike32432487293",
		    "name":   "workload_test02"
		  }`
	log.Println("+++ User mike32432487293 creates a new workload workload_test02")
	reader := strings.NewReader(tsJson)
	create_time = time.Now()
	request, err := http.NewRequest("DELETE", Config.WorkloadUrl, reader)
	Nil(t, err, "Failed to invoke http.NewRequest")
	request.Header.Set("Content-Type", "application/json")
	log.Printf("DELETE a request at URL %s\n", Config.WorkloadUrl)
	res, err := http.DefaultClient.Do(request)
	Nil(t, err, "Failed to invoke http.DefaultClient.Do")
	Equal(t, http.StatusOK, res.StatusCode, "")

	var d WkResponseData
	dec := json.NewDecoder(res.Body)
	err = dec.Decode(&d)
	Nil(t, err, "Failed to parse response JSON")
	Equal(t, 0, d.Ret, "Ret should be correct")
	Equal(t, "OK", d.Info, "Info should be correct")
}
