package listing

import (
	"fmt"
	"math/rand"
	"os"
	"time"
)

//----------------TASK 1----------------
func Task1() {
	a := [5]int{76, 77, 78, 79, 80}
	var b []int = a[1:4]
	fmt.Println(b) // [77, 78, 79]
}

//----------------TASK 2----------------
func Test() (x int) { // returned variable x
	//fmt.Println("!!!", x)		// x=0 - default value
	defer func() {
		x++ // x=2
	}()
	x = 1
	return // 2
}

func AnotherTest() int {
	var x int
	defer func() {
		x++
	}()
	x = 1
	return x // 1
}

//----------------TASK 3----------------
func Foo() error {
	var err *os.PathError = nil
	return err
}

func Task3() {
	err := Foo()
	fmt.Println(err)        // <nil>
	fmt.Println(err == nil) // false
}

//----------------TASK 4----------------
func Task4() {
	ch := make(chan int)
	go func() {
		for i := 0; i < 10; i++ {
			ch <- i
		}
		close(ch)
	}()

	for n := range ch {
		println(n)
	}
}

//----------------TASK 5----------------
type customError struct {
	msg string
}

func (e *customError) Error() string {
	return e.msg
}

func test() *customError {
	{
		// do something
	}
	return nil
}

func Task5() {
	var err error
	err = test()
	if err != nil {
		println("error")
		return
	}
	println("ok")
}

//----------------TASK 6----------------
func Task6() {
	var s = []string{"1", "2", "3"}
	modifySlice(s)
	fmt.Println(s)
}

func modifySlice(i []string) {
	i[0] = "3"         // s=[3,2,3] i=[3,2,3]
	i = append(i, "4") // не хватает емкости в базовом срезе. Го создаст новый слайс и скопирует в него i. i=[3,2,3,4]
	i[1] = "5"         // на предыдущем этапе потеряли ссылку на s. Теперь изменения делаем в i            i=[3,5,3,4]
	i = append(i, "6") // i=[3,5,3,4,6]
	fmt.Println("i", i)
}

//----------------TASK 7----------------
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

func Task7() {
	a := asChan(1, 3, 5, 7)
	b := asChan(2, 4, 6, 8)
	c := merge(a, b)
	for v := range c {
		fmt.Println(v)
	}
}
