package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/gorilla/mux"
	"golang.org/x/sync/semaphore"
)

func PrettyStruct(data interface{}) (string, error) {
	val, err := json.MarshalIndent(data, "", "    ")
	if err != nil {
		return "", err
	}
	return string(val), nil
}

func postOrder(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var order Order
	_ = json.NewDecoder(r.Body).Decode(&order)

	ordersAggregator.Enqueue(order)

	json.NewEncoder(w).Encode(&order)

	time.Sleep(time.Second * 3)

	ord, err := PrettyStruct(order)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Print("\nProducer recieved the data:\n", ord)
}

func makeOrder() {
	var wg sync.WaitGroup
	// nr of items we can send before the aggregator will block
	// nr of concurrent items (10)
	sem := semaphore.NewWeighted(10)
	for i := 0; i < 20; i++ {
		wg.Add(1)
		time.Sleep(time.Millisecond * 500)

		//aquire the semaphore which will block
		if err := sem.Acquire(context.Background(), 1); err != nil {
			log.Fatal(err)
		}

		go performPostRequest(&wg)

		// release the semaphore
		defer sem.Release(1)
	}
	wg.Wait()
}

func performPostRequest(wg *sync.WaitGroup) {
	const myUrl = "http://localhost:8080/aggregator"

	order := genOrder()

	requestBody, _ := json.Marshal(order)

	response, err := http.Post(myUrl, "application/json", bytes.NewBuffer(requestBody))

	if err != nil {
		panic(err)
	}

	defer response.Body.Close()
	fmt.Printf("\nData was sent to Aggregator\n")
}

func main() {
	router := mux.NewRouter()

	//URL path and the function to handle
	router.HandleFunc("/producer", postOrder).Methods("POST")

	go makeOrder()

	http.ListenAndServe(":3030", router)
}
