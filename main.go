// Это простая демонстрация того, как решить дилемму спящего парикмахера, классическую задачу информатики.
// который иллюстрирует сложности, возникающие при наличии нескольких процессов операционной системы. Здесь у нас есть
// конечное число парикмахеров, конечное число мест в зале ожидания, фиксированный период времени, в течение которого парикмахерская работает
// и клиенты приходят (примерно) через регулярные промежутки времени. Когда парикмахеру нечего делать, он проверяет
// комнату ожидания для новых клиентов, и если есть один или несколько, происходит стрижка. В противном случае парикмахер пойдет
// спать, пока не прибудет новый клиент. Итак, правила следующие:
//
// 		- если клиентов нет, парикмахер засыпает в кресле
// 		- клиент должен разбудить парикмахера, если он спит
//		- если клиент приходит, пока парикмахер работает, клиент уходит, если все стулья заняты и
//			 садимся на пустой стул, если он доступен
// 		- когда парикмахер заканчивает стрижку, он осматривает зал ожидания, есть ли ожидающие клиенты
// 			 и засыпаем, если их нет
// 		- магазин может перестать принимать новых клиентов во время закрытия, но парикмахеры не могут уйти, пока зал ожидания не будет пустой
// 		- после закрытия магазина и отсутствия клиентов в зоне ожидания парикмахер
//		     идет домой
//
// Спящий парикмахер был первоначально предложен в 1965 году Эдсгером Дейкстрой.
//
// Суть этой проблемы и ее решения заключалась в том, чтобы прояснить, что во многих случаях использование
// семафоров (мьютексов) не нужно.

package main

import (
	"fmt"
	"github.com/fatih/color"
	"math/rand"
	"time"
)

// переменные
var (
	seatingCapaсity = 10
	arrivalRate     = 100
	cutDuration     = 1000 * time.Millisecond
	timeOpen        = 10 * time.Second
)

func main() {

	//заполняем наш генератор случайных чисел
	rand.Seed(time.Now().UnixNano())

	// печатаем приветственное сообщение
	color.Yellow("The Sleeping Barber Problem")
	color.Yellow("---------------------------")

	// создаем каналы, если они нам нужны
	clientChan := make(chan string, seatingCapaсity)
	doneChan := make(chan bool)

	// создаем парикмахерскую
	shop := &BarberShop{
		ShopCapacity:    seatingCapaсity,
		HairCutDuration: cutDuration,
		NumberOfBarbers: 0,
		ClientsChan:     clientChan,
		BarbersDoneChan: doneChan,
		Open:            true,
	}

	color.Green("Барбершоп открыт!")

	// добавляем парикмахеров
	shop.AddBarber("Frank")
	shop.AddBarber("Gerrard")
	shop.AddBarber("Alex")
	shop.AddBarber("Milton")
	shop.AddBarber("Fred")
	shop.AddBarber("Kelly")

	// запускаем парикмахерскую как горутину
	shopClosing := make(chan bool)
	defer close(shopClosing)
	closed := make(chan bool)
	defer close(closed)

	go func() {
		<-time.After(timeOpen)
		shopClosing <- true
		shop.CloseShopForThisDay()
		closed <- true
	}()

	// добавляем клиентов
	i := 1
	go func() {
		for {
			randomMilliseconds := rand.Int() % (2 * arrivalRate)
			select {
			case <-shopClosing:
				return
			case <-time.After(time.Millisecond * time.Duration(randomMilliseconds)):
				shop.AddClient(fmt.Sprintf("Клиент №%d", i))
				i++
			}
		}
	}()
	// блокируем, пока парикмахерская не закроется
	<-closed
}
