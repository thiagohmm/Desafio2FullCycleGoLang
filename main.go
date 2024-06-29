package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"
)

func GetCep(cep string) string {
	ctx, cancel := context.WithTimeout(context.Background(), 500*time.Millisecond)
	defer cancel()
	startTime := time.Now()
	// Create a new client
	req, err := http.NewRequestWithContext(ctx, "GET", "https://viacep.com.br/ws/"+cep+"/json/", nil)
	if err != nil {

		if ctx.Err() == context.DeadlineExceeded {
			log.Fatal("Timeout")
		}
		log.Fatal(err)
	}
	response, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	defer response.Body.Close()

	duration := time.Since(startTime)

	fmt.Println("Response received")
	fmt.Println("Status Code: ", response.StatusCode)
	fmt.Println("Request Duration: ", duration)

	body, err := io.ReadAll(response.Body)
	if err != nil {
		log.Fatal(err)
	}
	return string(body)

}

func main() {

	fmt.Println(GetCep("03142001"))

}
