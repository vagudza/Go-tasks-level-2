package main

import (
	"fmt"
	"sync"
	"time"
)

/*
=== Or channel ===

Реализовать функцию, которая будет объединять один или более done каналов в single канал если один из его составляющих каналов закроется.
Одним из вариантов было бы очевидно написать выражение при помощи select, которое бы реализовывало эту связь,
однако иногда неизестно общее число done каналов, с которыми вы работаете в рантайме.
В этом случае удобнее использовать вызов единственной функции, которая, приняв на вход один или более or каналов, реализовывала весь функционал.

Определение функции:
var or func(channels ...<- chan interface{}) <- chan interface{}

Пример использования функции:
sig := func(after time.Duration) <- chan interface{} {
	c := make(chan interface{})
	go func() {
		defer close(c)
		time.Sleep(after)
}()
return c
}

start := time.Now()
<-or (
	sig(2*time.Hour),
	sig(5*time.Minute),
	sig(1*time.Second),
	sig(1*time.Hour),
	sig(1*time.Minute),
)

fmt.Printf(“fone after %v”, time.Since(start))


OK go vet -c=10 task.go
OK golint task.go

Запуск:
go run task.go
*/

func main() {

	or := func(channels ...<-chan interface{}) <-chan interface{} {
		out := make(chan interface{})
		wg := sync.WaitGroup{}
		wg.Add(len(channels))

		for _, channel := range channels {

			// передаем канал как значение в горутину.  Просто передайте переменную в качестве параметра, так она будет продублирована,
			// и каждая горутина будет иметь разные переменные для себя.
			// Иначе начнется состояние гонки за разделяемую переменную channel, итерируемую в for:
			// Проверка программы на состояние гонки: go run -race main.go
			//
			// go func() {
			// 	for value := range channel {
			// 		out <- value
			// 	}
			//
			go func(ch <-chan interface{}) {
				for value := range ch {
					out <- value
				}
				fmt.Printf("Channel closed right now\n")
				wg.Done()
			}(channel)
		}

		go func() {
			wg.Wait()
			close(out)
		}()

		return out
	}

	sig := func(after time.Duration) <-chan interface{} {
		// создаем канал, в который можно писать и читать
		c := make(chan interface{})

		go func() {
			defer close(c)
			time.Sleep(after)
		}()

		// почему тогда нет ошибки типов, ведь возвращаем канал, из которого можно только читать?
		return c
	}

	start := time.Now()
	<-or(
		sig(2*time.Second),
		sig(5*time.Second),
		sig(1*time.Second),
		sig(3*time.Second),
		sig(4*time.Second),
	)

	fmt.Printf("fone after %v", time.Since(start))
}
