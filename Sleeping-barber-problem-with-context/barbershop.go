package main

import (
	"github.com/fatih/color"
	"time"
)

type BarberShop struct {
	NumberOfBarbers int
	HairCutDuration time.Duration
	Open            bool
	BarbersDoneChan chan bool
	ClientChan      chan string
}

func NewBarberShop(numberOfBarbers int, hairCutDuration time.Duration, open bool, barbersDoneChan chan bool, clientChan chan string) *BarberShop {
	return &BarberShop{
		NumberOfBarbers: numberOfBarbers,
		HairCutDuration: hairCutDuration,
		Open:            open,
		BarbersDoneChan: barbersDoneChan,
		ClientChan:      clientChan,
	}
}

func (shop *BarberShop) AddNewBarbers(barberName string) {
	defer wg.Done()
	shop.NumberOfBarbers++
	var isSleeping bool
	color.Green("%s идет проверять комнату ожидания...", barberName)
	go func() {
		for {
			if len(shop.ClientChan) == 0 {
				color.Green("Клиентов нет, по этому %s пошел спать", barberName)
				isSleeping = true
			}
			select {
			case client, barberShopOpen := <-shop.ClientChan:
				if barberShopOpen {
					if isSleeping {
						color.Blue("Клиент %s разбудил %s", client, barberName)
						isSleeping = false
					}
					shop.cutHair(client, barberName)
					color.Green("%s идет проверять комнату ожидания...", barberName)
				}
			case <-ctx.Done():
				if len(shop.ClientChan) == 0 {
					if isSleeping {
						color.Yellow("%s просыпается так как пора идти домой", barberName)
						isSleeping = false
						shop.barberGoingToHome(barberName)
					}
				} else {
					continue
				}
				return
			}
		}
	}()
}

func (shop *BarberShop) cutHair(clientName, barberName string) {
	color.Green("%s стрижет клиента %s", barberName, clientName)
	time.Sleep(shop.HairCutDuration)
	color.Green("%s закончил стричь клиента %s", barberName, clientName)
}

func (shop *BarberShop) barberGoingToHome(barberName string) {
	shop.BarbersDoneChan <- true
	color.Yellow("%s закончил и пошел домой!", barberName)
}

func (shop *BarberShop) closingShop() {
	shop.Open = false
	color.Yellow("Барбершоп на сегодня, закрывается")
	close(shop.ClientChan)
	for i := 1; i <= shop.NumberOfBarbers; i++ {
		<-shop.BarbersDoneChan
	}
	close(barbersDoneChan)
	color.Yellow("Барбершоп закрывает, все барберы пошли домой!")
}

func (shop *BarberShop) AddClient(clientName string) {
	if shop.Open {
		select {
		case shop.ClientChan <- clientName:
			color.Blue("Пришел клиент %s, зашел в зал ожидания!", clientName)
		default:
			color.Red("В зале ожидания нет места, %s уходит!!!", clientName)
		}
	} else {
		color.Red("Пришел клиент %s, но барбрешоп уже закрыт!!!", clientName)
	}
}
