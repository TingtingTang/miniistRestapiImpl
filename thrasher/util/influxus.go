package main

import (
	"fmt"
	"github.com/Sirupsen/logrus"
	client "github.com/influxdata/influxdb/client/v2"
	"github.com/vlad-doru/influxus"
	"time"
)

func main() {
	// Create the InfluxDB client.
	influxClient, err := client.NewHTTPClient(client.HTTPConfig{
		Addr:     "http://45.55.21.6:8086",
		Username: "root",
		Password: "amwuck00",
	})
	if err != nil {
		logrus.Fatalf("Error while creating the client: %v", err)
	} else {
		fmt.Println("Seems OK now")
	}
	// Create and add the hook.
	hook, err := influxus.NewHook(
		&influxus.Config{
			Client:             influxClient,
			Database:           "testhub", // DATABASE MUST BE CREATED
			DefaultMeasurement: "logrus",
			BatchSize:          100, // default is 100
			BatchInterval:      5,   // default is 5 seconds
		})
	if err != nil {
		logrus.Fatalf("Error while creating the hook: %v", err)
	}
	// Add the hook to the standard logger.
	// logrus.StandardLogger().Hooks.Add(hook)
	logrus.AddHook(hook)

	logrus.WithFields(logrus.Fields{
		"animal": "walrus",
	}).Info("A walrus appears")

	logrus.Debug("Useful debugging information.")
	logrus.Info("Something noteworthy happened!")
	//logrus.Warn("You should probably take a look at this.")
	//logrus.Error("Something failed but I'm not quitting.")
	time.Sleep(10 * time.Second)
	fmt.Println("OK so far")
	// Calls os.Exit(1) after logging
	logrus.Fatal("Bye.")
	// Calls panic() after logging
	logrus.Panic("I'm bailing.")
}
