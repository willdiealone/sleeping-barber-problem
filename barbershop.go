package main

import (
	"github.com/fatih/color"
	"time"
)

type BarberShop struct {
	ShopCapacity    int
	HairCutDuration time.Duration
	NumberOfBarbers int
	BarbersDoneChan chan bool
	ClientsChan     chan string
	Open            bool
}

// AddBarber ... метод добавляет барбера, который сразу начинает работать
func (shop *BarberShop) AddBarber(barberName string) {
	// добавляем барбера
	shop.NumberOfBarbers++
	go func() {
		for {
			isSleeping := false
			color.Yellow("%s пошел в комнату ожидания проверить есть ли посетители...", barberName)
			// если клиентов нет идем спать.
			if len(shop.ClientsChan) == 0 {
				color.Yellow("%s пошел спать, так как клиентов нет.", barberName)
				isSleeping = true
			}
			// слушаем клиентов из канала
			if client, shopOpen := <-shop.ClientsChan; shopOpen {
				// если барбер спит
				if isSleeping {
					color.Yellow("%s разбудил %s", client, barberName)
					isSleeping = false
				}
				// стрижка
				shop.cutHair(barberName, client)
			} else {
				if isSleeping {
					color.Yellow("%s просыпается так как пора идти домой", barberName)
					isSleeping = false
				}
				// если магазин закрыт то барбер идет домой
				shop.sendBarberToHome(barberName)
				return
			}
		}
	}()
}

// метод имплементирует стрижку клиента
func (shop *BarberShop) cutHair(barberName, client string) {
	color.Green("%s стрижет волосы клиента с имененм %s", barberName, client)
	time.Sleep(shop.HairCutDuration)
	color.Green("%s закончил стричь волосы клиента с имененм %s", barberName, client)
}

// метод отправляет барбера домой
func (shop *BarberShop) sendBarberToHome(barberName string) {
	color.Cyan("%s пошел домой", barberName)
	shop.BarbersDoneChan <- true
}

// CloseShopForThisDay ... метод закрывает барбершоп
func (shop *BarberShop) CloseShopForThisDay() {
	color.Cyan("Барбершоп на сегодня закрывается.")
	// закрываем канал с клиентами так как мы больше не можем их принимать
	close(shop.ClientsChan)
	// магазин закрылся
	shop.Open = false
	for i := 1; i <= shop.NumberOfBarbers; i++ {
		// поток блокируется до тех про пока не закончат свою работу все барберы
		<-shop.BarbersDoneChan
	}
	close(shop.BarbersDoneChan)
	color.Green("-------------------------------------------------")
	color.Green("Барбершоп закрыт, все барберы разошлись по домам.")
}

// AddClient ... метод добавялет клиента
func (shop *BarberShop) AddClient(clientName string) {
	color.Green("*** Пришел клиент %s ***", clientName)
	if shop.Open {
		select {
		case shop.ClientsChan <- clientName:
			color.Yellow("%s занял место в зале ожидания.", clientName)
		default:
			color.Red("В зале ожидания нет места, %s уходит!!!", clientName)
		}
	} else {
		color.Red("Барбершоп уже закрыт, %s уходит!", clientName)
	}
}
