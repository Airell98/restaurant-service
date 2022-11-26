package main

import (
	"math/rand"
	"restaurant-service/handler/rest"
	"time"
)



func main() {
	rand.Seed(time.Now().Unix())
	rest.StartApp()
}