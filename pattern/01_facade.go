package pattern

import (
	"fmt"
	"math/rand"
)

/*
	Реализовать паттерн «фасад».
Объяснить применимость паттерна, его плюсы и минусы,а также реальные примеры использования данного примера на практике.
	https://en.wikipedia.org/wiki/Facade_pattern


Проблема:
Минимизировать зависимость подсистем некоторой сложной системы и обмен информацией между ними.

Решение:
Фасад — простой интерфейс для работы со сложным фреймворком. Фасад не имеет всей функциональности фреймворка,
но зато скрывает его сложность от клиентов.

В примере рассматривается торговый автомат, который продает напитки в банках.
Торговый автомат состоит из кнопок выбора напитка PanelButtons, платежной системы PaymentSystem, и
драйвера электродвигателя ProductDriver, который управляет выдачей товара.

плюсы:
	+Изолирует клиентов от компонентов сложной подсистемы.

минусы:
	-Фасад рискует стать божественным объектом, привязанным ко всем классам программы. (объект, который хранит
		 в себе «слишком много» или делает «слишком много». )

Фасад - интерфейс пользователя, который просто передает управление внутренним элементам системы торгового автомата
*/

// управление электромотором
type ProductDriver struct{}

func (pd *ProductDriver) GiveProduct(ok bool, pos int) {
	if ok {
		fmt.Printf("Vending gives product by position #%d\n", pos)
	}
}

//------------------------------------
// платежная система
type PaymentSystem struct {
	profit      int
	productCost int
	position    int
}

func (p *PaymentSystem) SetProductInfo(cost int, pos int) {
	p.productCost = cost
	p.position = pos
	fmt.Printf("Cost of product is %d\n", p.productCost)
}

func (p *PaymentSystem) Pay(sum int) (bool, int) {
	if sum >= p.productCost {
		p.profit += p.productCost
		fmt.Printf("User have %d$, it is enough\n", sum)
		return true, p.position
	} else {
		fmt.Printf("User have %d$, it is not enough\n", sum)
		return false, p.position
	}
}

//------------------------------------
// панель управления для выбора товара
type PanelButtons struct {
	position int
}

// пользователь выбирает товар
func (p *PanelButtons) SelectedProductNumber() {
	p.position = rand.Intn(32)
	fmt.Printf("User select product #%d\n", p.position)
}

// получаем цену на выбранный товар
func (p *PanelButtons) GetProductInfo() (int, int) {
	if p.position%2 == 0 {
		return 10, p.position
	} else {
		return 20, p.position
	}
}

//------------------------------------
// "фасад"
type VendingMachine struct {
	panel   *PanelButtons
	payment *PaymentSystem
	driver  *ProductDriver
}

func NewVendingMachine() *VendingMachine {
	fmt.Println("---Pattern Facade---\nLoading vending machine...")
	return &VendingMachine{
		panel:   &PanelButtons{},
		payment: &PaymentSystem{},
		driver:  &ProductDriver{},
	}
}

// Реализация паттерна: вызов методов различных структур в определенном порядке.
func (v *VendingMachine) Start() {
	fmt.Println("Vending machine ready to work")
	v.panel.SelectedProductNumber()
	v.payment.SetProductInfo(v.panel.GetProductInfo())
	v.driver.GiveProduct(v.payment.Pay(rand.Intn(30)))
	fmt.Println("Vending machine work done!")
}
