#!/bin/sh

echo "Enter test type"
read type

case $type  in 
	"rest") 
		go test -v
	;;
	"db")
		cd mgodb
		go test -v
	;;
	"bench")
		cd mgodb
		go test -v -bench=. -cpu 2,4,8,16,32
	;;
	"e2e")
		e2e_test.sh
	;;
esac

