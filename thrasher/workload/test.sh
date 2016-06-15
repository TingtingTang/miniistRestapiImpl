#!/bin/sh

###################################################################################
#
# Insert Workload
#
###################################################################################
invoke() {
	payload=$(echo $payload|sed -e "s/\n//g")
	echo "curl -v -X $method -d '$payload' -H "Content-Type: $apptype" $url"
	curl -v -X $method -d "$payload" -H "Content-Type: $apptype" $url
}

addr="http://127.0.0.1:8002"
method=POST
url="$addr/workload"
apptype="application/json"
read -r -d '' payload <<-'EOF'
{
	"action": "insert_workload",
	"user":   "mike32432487293",
	"workload": {
			"name": "workload_test01",
			"team":   "testteam01",
			"run_as_user":   "mega",
			"desc":   "desc01",
			"category":  "category01",
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
}
EOF

invoke

###################################################################################
#
# Get Workload
#
###################################################################################

addr="http://127.0.0.1:8002"
method=GET
url="$addr/workload?user=mike32432487293&&name=workload_test01"
apptype="application/json"
#read -r -d '' payload <<-'EOF'
#EOF

invoke

###################################################################################
#
# Update Workload
#
###################################################################################

addr="http://127.0.0.1:8002"
method=PUT
url="$addr/workload"
apptype="application/json"
read -r -d '' payload <<-'EOF'
	{ 
	"action": "update_workload",
		"user":   "mike32432487293",
		"workload": {
			"id":"573b17059f72ddeac1002ee6",
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
	}
EOF

invoke


###################################################################################
#
# Get Workloads
#
###################################################################################

addr="http://127.0.0.1:8002"
method=GET
url="$addr/workloads?user=mike32432487293&startpos=1&counter=100"
apptype="application/json"
#read -r -d '' payload <<-'EOF'
#EOF

invoke

###################################################################################
#
# Delete Workload
#
###################################################################################

addr="http://127.0.0.1:8002"
method=DELETE
url="$addr/workload"
apptype="application/json"
read -r -d '' payload <<-'EOF'
	{ 
	    "action": "delete_workload",
	    "user":   "mike32432487293",
	    "name":   "workload_test01"
	}
EOF

invoke

