package pattern

import "fmt"

/*
	Реализовать паттерн «состояние».
Объяснить применимость паттерна, его плюсы и минусы, а также реальные примеры использования данного примера на практике.
	https://en.wikipedia.org/wiki/State_pattern


Состояние — это поведенческий паттерн, позволяющий динамически изменять поведение объекта при смене его состояния.
Поведения, зависящие от состояния, переезжают в отдельные классы. Первоначальный класс хранит ссылку на один из таких
объектов-состояний и делегирует ему работу.

Паттерн должен применяться:

    когда поведение объекта зависит от его состояния
    поведение объекта должно изменяться во время выполнения программы
    состояний достаточно много и использовать для этого условные операторы, разбросанные по коду, достаточно затруднительно

Пример:
*/

// MobileAlertStater - общий интерфейс для различных состояний
type MobileAlertStater interface {
	Alert() string
}

// MobileAlert реализует оповещение в зависимости от его состояния.
type MobileAlert struct {
	state MobileAlertStater
}

func (a *MobileAlert) Alert() string {
	return a.state.Alert()
}

// SetState изменяет состояния
func (a *MobileAlert) SetState(state MobileAlertStater) {
	a.state = state
}

func NewMobileAlert() *MobileAlert {
	return &MobileAlert{state: &MobileAlertVibration{}}
}

// MobileAlertVibration реализует состояние телефона в беззвучном режиме
type MobileAlertVibration struct {
}

func (a *MobileAlertVibration) Alert() string {
	return "Vrrr... Brrr... Vrrr..."
}

// MobileAlertSong реализует состояние телефона в обычном режиме
type MobileAlertSong struct {
}

// Alert returns a alert string
func (a *MobileAlertSong) Alert() string {
	return "Белые розы, Белые розы. Беззащитны шипы..."
}

func StatePatternStart() {
	mobile := NewMobileAlert()
	fmt.Println(mobile.Alert())

	mobile.SetState(&MobileAlertSong{})
	fmt.Println(mobile.Alert())
}
