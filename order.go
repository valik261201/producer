package main

import (
	"math/rand"
	"time"
)

type Order struct {
	Id         int   `json:"id"`
	Items      []int `json:"items"`
	Priority   int   `json:"priority"`
	MaxWait    int   `json:"max-wait"`
	PickUpTime int   `json:"pick-up-time"`
}

func genRandomNum(min, max int) int {
	// Intn generates random number between [0,n)
	return min + rand.Intn(max-min)
}

func genItems() []int {
	n := genRandomNum(1, 10)
	var items = make([]int, n)

	for i := 0; i < n; i++ {
		items[i] = genRandomNum(1, 10)
	}

	return items
}

func genMaxWait(menu Foods) int {
	maxTime := 0

	for i := 0; i < len(menu.Foods); i++ {
		if menu.Foods[i].Preparation_time > maxTime {
			maxTime = menu.Foods[i].Preparation_time
		}
	}
	maxWait := float32(maxTime) * 1.3
	return int(maxWait)
}

func getUnixTimestamp() int {
	now := time.Now()
	sec := now.Unix()

	return int(sec)
}

func genOrder() Order {

	var menu Foods

	return Order{
		Id:         int(genRandomNum(1, 100)),
		Items:      genItems(),
		Priority:   int(genRandomNum(1, 5)),
		MaxWait:    genMaxWait(menu),
		PickUpTime: getUnixTimestamp(),
	}

}
