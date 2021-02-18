package main

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"
)

// Configuration struct
type Configuration struct {
	Server string   `json:"server"`
	Name   string   `json:"name"`
	Token  string   `json:"token"`
	Status []string
}

// Payload struct for Rocket.Chat
type Payload struct {
	UserID string `json:"userId"`
	Data   Data   `json:"data"`
}

// Data struct is needed for the nice status text
type Data struct {
	StatusText string `json:"statusText"`
}

// ConfigurationFile is needed for global config
var ConfigurationFile Configuration

func main() {
	log.Println("That is the real shit that everyone needs!")
	jsonFile, err := os.Open(os.Getenv("HOME") + "/.config/rocketstatus/config.json")

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func() {
		for sig := range c {
			log.Println(sig)
			os.Exit(0)
		}
	}()

	if err != nil {
		initialConfig := Configuration{
			Server: "https://rocket.chat",
			Name: "abc123",
			Token: "Iamthelongtokenyoucangetfromrocketchat123111elf_",
			Status: []string{
				"TheFirstMessage",
				"TheSecondMessage",
			},
		}

		file, _ := json.MarshalIndent(initialConfig, "", "")

		_ = ioutil.WriteFile(os.Getenv("HOME") + "/.config/rocketstatus/config.json", file, 0644)
		log.Println("Created initial configuration")
	}

	byteValue, _ := ioutil.ReadAll(jsonFile)
	json.Unmarshal(byteValue, &ConfigurationFile)

	log.Println("I will do this with userID " + ConfigurationFile.Name)

	log.Println("Reset the Status")
	updateStatus("I reset myself, please wait...")

	time.Sleep(4000 * time.Millisecond)
	log.Println("Let's go 🚀")
	for {
		for _, value := range ConfigurationFile.Status {
			updateStatus(value)
			time.Sleep(3000 * time.Millisecond)
		}
	}
}

func updateStatus(shownStatus string) {
	statusBla := Data{
		StatusText: shownStatus,
	}

	data := Payload{
		UserID: ConfigurationFile.Name,
		Data:   statusBla,
	}

	payloadBytes, err := json.Marshal(data)
	if err != nil {
		log.Fatal("Excuse me, why the fuck is JSON parsing not working anymore!?")
	}
	body := bytes.NewReader(payloadBytes)

	req, err := http.NewRequest("POST", ConfigurationFile.Server + "/api/v1/users.update", body)
	if err != nil {
		log.Fatal("I think my PC is hating me to create a new object every 3 seconds lol")
	}

	req.Header.Set("X-Auth-Token", ConfigurationFile.Token)
	req.Header.Set("X-User-Id", ConfigurationFile.Name)
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Println("Eh fuck, I think I got the request limit")
		time.Sleep(10000 * time.Millisecond)
	} else {
		defer resp.Body.Close()
	}
}
