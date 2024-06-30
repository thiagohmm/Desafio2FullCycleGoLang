package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"
)

type Links struct {
	Link string
	Resp string
}

func GetCep(linkChan chan Links) {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()
	startTime := time.Now()

	linkStruct := <-linkChan
	link := linkStruct.Link
	req, err := http.NewRequestWithContext(ctx, "GET", link, nil)
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

	linkChan <- Links{link, string(body)}

}

func main() {

	//fmt.Println(GetCep("https://viacep.com.br/ws/03142001/json/"))
	//fmt.Println(GetCep("https://brasilapi.com.br/api/cep/v1/03142001"))
	c1 := make(chan Links)
	c2 := make(chan Links)

	go GetCep(c1)
	go GetCep(c2)

	link := Links{"https://viacep.com.br/ws/03142001/json/", ""}
	c1 <- link

	link = Links{"https://brasilapi.com.br/api/cep/v1/03142001", ""}
	c2 <- link

	select {

	case msg := <-c1: // rabbitmq

		fmt.Printf("Received from ViaCep: %s - %s\n", msg.Link, msg.Resp)

	case msg := <-c2: // kafka

		fmt.Printf("Received from BrasilApi: %s - %s\n", msg.Link, msg.Resp)

	case <-time.After(time.Second * 3):
		println("timeout")

	}

}
