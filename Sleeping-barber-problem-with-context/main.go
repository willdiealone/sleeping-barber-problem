package main

import (
	"context"
	"fmt"
	"github.com/fatih/color"
	"math/rand"
	"sync"
	"time"
)

var (
	cutDuration     = 1000 * time.Millisecond
	timeOpen        = 10
	arrivalRate     = 100
	barbersDoneChan = make(chan bool)
	seatedContains  = 10
	clientChan      = make(chan string, seatedContains)
	Closing         = make(chan bool)
	ctx, cancel     = context.WithCancel(context.Background())
	wg              sync.WaitGroup
)

func main() {
	color.Yellow("Барбершоп открыт")
	color.Yellow("----------------")
	rand.Seed(time.Now().UnixNano())
	wg.Add(5)
	shop := NewBarberShop(0, cutDuration, true, barbersDoneChan, clientChan)
	shop.AddNewBarbers("Andrey")
	shop.AddNewBarbers("Yuliya")
	shop.AddNewBarbers("Dennis")
	i := 0
	go func() {
		defer wg.Done()
		for {
			randomMilliseconds := rand.Int() % (2 * arrivalRate)
			select {
			case <-time.After(time.Duration(randomMilliseconds) * time.Millisecond):
				shop.AddClient(fmt.Sprintf("#%d", i))
				i++
			case <-Closing:
				cancel()
				return
			}
		}
	}()

	go func() {
		defer wg.Done()
		<-time.After(time.Duration(timeOpen) * time.Second)
		Closing <- true
		shop.closingShop()
	}()

	wg.Wait()
}
