Что выведет программа? Объяснить вывод программы.

```go
package main

import (
	"fmt"
	"math/rand"
	"time"
)

func asChan(vs ...int) <-chan int {
	c := make(chan int)

	go func() {
		for _, v := range vs {
			c <- v
			time.Sleep(time.Duration(rand.Intn(1000)) * time.Millisecond)
		}

		close(c)
	}()
	return c
}

func merge(a, b <-chan int) <-chan int {
	c := make(chan int)
	go func() {
		for {
			select {
			case v := <-a:
				c <- v
			case v := <-b:
				c <- v
			}
		}
	}()
	return c
}

func main() {

	a := asChan(1, 3, 5, 7)
	b := asChan(2, 4 ,6, 8)
	c := merge(a, b )
	for v := range c {
		fmt.Println(v)
	}
}
```

Ответ:
```
Программа выведет нечетные числа из a и четные числа из b в случайном порядке, но порядок сохранится отдельно для четных и 
отдельно нечетных чисел, а после - нули. После того, как горутина внутри asChan() закроет канал, функция merge() будет читать 
из закрытых каналов. Поскольку тип каналов - int, то в бесконечном цикле будут извлекаться нули из закрытых каналов. Вывод
значений из канала, возвращаемого merge() - в main() в бесконечном цикле с range.

```
