package main

import (
	"fmt"
	"net/http"
	"time"

	goreplica "github.com/JonCSykes/GoReplica"
)

func main() {

	tr := &http.Transport{
		MaxIdleConns:       10,
		IdleConnTimeout:    30 * time.Second,
		DisableCompression: true,
	}

	httpClient := &http.Client{Transport: tr}

	replicaClient := goreplica.Client{ServiceEndpoint: "https://api.replicastudios.com", ClientID: "xxx", ClientSecret: "xxx", HTTPClient: httpClient}

	replicaClient.Auth()

	fmt.Println("Access Token : " + replicaClient.AccessToken)

	voices, err := replicaClient.GetVoices()
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println(voices)

	urls, speechErr := replicaClient.GetSpeech("This is just a test.", "d6ad9af8-6361-4c44-a574-9c2b24e73dc2", 128, 44100, goreplica.MP3)
	if speechErr != nil {
		fmt.Println(err)
	}

	fmt.Println(urls)
}
