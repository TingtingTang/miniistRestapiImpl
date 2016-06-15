#!/bin/bash
# user daemon
# chkconfig: 345 20 80
# description: user daemon
# processname: user

DAEMON_PATH="."

DAEMON=go 
DAEMONOPTS="run main.go join.go login.go respond.go user.go helper.go"

NAME=testhub_user_daemon
DESC="User daemon for testhub"
RUNDIR=/tmp/run
LOGDIR=/tmp/log
PIDFILE=/tmp/run/$NAME.pid
SCRIPTNAME=/etc/init.d/$NAME

mkdir -p $RUNDIR
mkdir -p $LOGDIR

case "$1" in
start)
	printf "%-50s" "Starting $NAME..."
	cd $DAEMON_PATH
        export GOPATH="$HOME/go"
        echo "GOPATH: $GOPATH "
        CMD="$DAEMON $DAEMONOPTS > $LOGDIR/${PID}.log 2>&1 && echo $! "
	echo "Run command: $CMD"
	PID=`($DAEMON $DAEMONOPTS > $LOGDIR/${PID}.log 2>&1)  && echo $! ` 
	echo "AFTER it"
	# echo "Saving PID" $PID " to " $PIDFILE
        if [ $? -ne 0 -o -z "$PID" ]; then
	    echo "PID: $PID "
	    echo "GOPATH: $GOPATH "
            printf "%s\n" "Fail"
        else
            echo $PID > $PIDFILE
            printf "%s\n" "Ok"
        fi
;;
status)
        printf "%-50s" "Checking $NAME..."
        if [ -f $PIDFILE ]; then
            PID=`cat $PIDFILE`
            if [ -z "`ps axf | grep ${PID} | grep -v grep`" ]; then
                printf "%s\n" "Process dead but pidfile exists"
            else
                echo "Running"
            fi
        else
            printf "%s\n" "Service not running"
        fi
;;
stop)
        printf "%-50s" "Stopping $NAME"
            PID=`cat $PIDFILE`
            cd $DAEMON_PATH
        if [ -f $PIDFILE ]; then
            kill -HUP $PID
            printf "%s\n" "Ok"
            rm -f $PIDFILE
        else
            printf "%s\n" "pidfile not found"
        fi
;;

restart)
  	$0 stop
  	$0 start
;;

*)
        echo "Usage: $0 {status|start|stop|restart}"
        exit 1
esac
