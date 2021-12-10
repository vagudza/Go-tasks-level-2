package pattern

import "fmt"

/*
	Реализовать паттерн «стратегия».
Объяснить применимость паттерна, его плюсы и минусы, а также реальные примеры использования данного примера на практике.
	https://en.wikipedia.org/wiki/Strategy_pattern


	Стратегия — это поведенческий паттерн, позволяющий выбор поведения алгоритма в ходе исполнения.
	Этот паттерн определяет алгоритмы, инкапсулирует их и использует их взаимозаменяемо.

	Паттерн Стратегия позволяет вам изменять внутренности объекта. Паттерн Декоратор позволяет вам изменять оболочку объекта.

	Другие объекты содержат ссылку на объект-стратегию и делегируют ей работу. Программа может подменить этот объект
	другим, если требуется иной способ решения задачи.


	Пример: Реализация взаимозаменяемого объекта оператора, который оперирует целыми числами.
*/

type Operator interface {
	Apply(int, int) int
}

type Operation struct {
	Operator Operator
}

func (o *Operation) Operate(leftValue, rightValue int) int {
	return o.Operator.Apply(leftValue, rightValue)
}

type Addition struct{}

func (Addition) Apply(lval, rval int) int {
	return lval + rval
}

type Multiplication struct{}

func (Multiplication) Apply(lval, rval int) int {
	return lval * rval
}

func StrategyPatternStart() {
	mult := Operation{Multiplication{}}
	fmt.Println(mult.Operate(3, 5)) // 15

	add := Operation{Addition{}}
	fmt.Println(add.Operate(3, 5)) // 8
}
