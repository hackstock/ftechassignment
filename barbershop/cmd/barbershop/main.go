package main

import (
	"barbershop/pkg/shop"
	"fmt"
	"log"
	"math/rand"
	"time"

	"github.com/kelseyhightower/envconfig"
	"go.uber.org/zap"
)

var env = struct {
	SeatingCapacity  int    `envconfig:"SEAT_CAPACITY" default:"10"`
	ArrivalRate      int    `envconfig:"ARRIVAL_RATE_IN_SEC" default:"100"`
	DurationPerCut   int    `envconfig:"CUT_DURATION_IN_MSEC" default:"1000"`
	ShopOpenDuration int    `envconfig:"OPENED_UNTIL_IN_SEC" default:"10"`
	Environment      string `envconfig:"ENVIRONMENT" default:"development"`
}{}

func init() {
	err := envconfig.Process("", &env)
	if err != nil {
		log.Fatalf("failed loading env vars : %v", err)
	}
}

func initLogger(environment string) (*zap.Logger, error) {
	if environment == "development" {
		return zap.NewDevelopment()
	}

	return zap.NewProduction()
}

func main() {
	rand.Seed(time.Now().UnixNano())
	logger, err := initLogger(env.Environment)
	if err != nil {
		log.Fatalf("failed initializing logger : %v", err)
	}

	clients := make(chan string, env.SeatingCapacity)
	done := make(chan bool)

	cfg := &shop.ConfigOptions{
		SeatingCapacity: env.SeatingCapacity,
		DurationPerCut:  env.DurationPerCut,
		Clients:         clients,
		Done:            done,
	}

	shop := shop.NewShop(cfg, logger)

	shop.AddBarber("Kofi")
	shop.AddBarber("Edward")
	shop.AddBarber("Rosina")

	shopClosing := make(chan bool)
	closed := make(chan bool)

	go func() {
		<-time.After(time.Duration(env.ShopOpenDuration) * time.Second)
		shopClosing <- true
		shop.Close()
		closed <- true
	}()

	customerId := 1
	go func() {
		for {
			waitPeriod := rand.Int() % (2 * env.ArrivalRate)
			select {
			case <-shopClosing:
				return
			case <-time.After(time.Millisecond * time.Duration(waitPeriod)):
				shop.AddCustomer(fmt.Sprintf("Customer %d", customerId))
				customerId++
			}
		}
	}()

	<-closed
}
