package pattern

import "fmt"

/*
	Реализовать паттерн «цепочка вызовов».
Объяснить применимость паттерна, его плюсы и минусы, а также реальные примеры использования данного примера на практике.
	https://en.wikipedia.org/wiki/Chain-of-responsibility_pattern


Цепочка обязанностей (англ. Chain of responsibility) — поведенческий шаблон проектирования, предназначенный для
организации в системе уровней ответственности.

Цепочка обязанностей — это поведенческий паттерн проектирования, который позволяет передавать запросы последовательно
по цепочке обработчиков. Каждый последующий обработчик решает, может ли он обработать запрос сам и стоит ли передавать
запрос дальше по цепи.

Применимость:
	-С помощью Цепочки обязанностей вы можете связать потенциальных обработчиков в одну цепь и при получении запроса
	поочерёдно спрашивать каждого из них, не хочет ли он обработать запрос.
	-Цепочка обязанностей позволяет запускать обработчиков последовательно один за другим в том порядке, в котором
	они находятся в цепочке.
	-В любой момент вы можете вмешаться в существующую цепочку и переназначить связи так, чтобы убрать или добавить новое звено.

Плюсы:
	+Уменьшает зависимость между клиентом и обработчиками.
	+Реализует принцип единственной обязанности.
	+Реализует принцип открытости/закрытости.

Минусы:
	-Запрос может остаться никем не обработанным.

В качестве примера рассмотрим госпиталь, имеющий помещения "Регистратура", "Доктор", "Комната медикаментов", "Касса"
Используя пример больницы, пациент сперва попадает в Приемное отделение. Затем, зависимо от его состояния,
Регистратура отправляет его к следующему исполнителю в цепи.
*/

// Обработчик определяет общий для всех конкретных обработчиков интерфейс. Обычно достаточно описать единственный метод обработки запросов,
// но иногда здесь может быть объявлен и метод выставления следующего обработчика
type department interface {
	execute(*patient)
	setNext(department)
}

// Базовый обработчик — опциональный класс, который позволяет избавиться от дублирования одного и того же кода во всех конкретных
// обработчиках. Обычно этот класс имеет поле для хранения ссылки на следующий обработчик в цепочке. Клиент связывает обработчики
// в цепь, подавая ссылку на следующий обработчик через конструктор или сеттер поля. Также здесь можно реализовать базовый метод
// обработки, который бы просто перенаправлял запрос следующему обработчику, проверив его наличие.

// Конкретные обработчики содержат код обработки запросов. При получении запроса каждый обработчик решает, может ли он обработать
// запрос, а также стоит ли передать его следующему объекту. В большинстве случаев обработчики могут работать сами по себе и быть
// неизменяемыми, получив все нужные детали через параметры конструктора.

// -------- Приемное отделение --------
type reception struct {
	next department
}

func (r *reception) execute(p *patient) {
	if p.registrationDone {
		fmt.Println("Patient registration already done")
		r.next.execute(p)
		return
	}
	fmt.Println("Reception registering patient")
	p.registrationDone = true
	r.next.execute(p)
}

func (r *reception) setNext(next department) {
	r.next = next
}

// ----------- Доктор -------------
type doctor struct {
	next department
}

func (d *doctor) execute(p *patient) {
	if p.doctorCheckUpDone {
		fmt.Println("Doctor checkup already done")
		d.next.execute(p)
		return
	}
	fmt.Println("Doctor checking patient")
	p.doctorCheckUpDone = true
	d.next.execute(p)
}

func (d *doctor) setNext(next department) {
	d.next = next
}

// ----------- Комната медикаментов -----------
type medical struct {
	next department
}

func (m *medical) execute(p *patient) {
	if p.medicineDone {
		fmt.Println("Medicine already given to patient")
		m.next.execute(p)
		return
	}
	fmt.Println("Medical giving medicine to patient")
	p.medicineDone = true
	m.next.execute(p)
}

func (m *medical) setNext(next department) {
	m.next = next
}

// -------- Кассир ------------
type cashier struct {
	next department
}

func (c *cashier) execute(p *patient) {
	if p.paymentDone {
		fmt.Println("Payment Done")
	}
	fmt.Println("Cashier getting money from patient")
}

func (c *cashier) setNext(next department) {
	c.next = next
}

// Клиент может либо сформировать цепочку обработчиков единожды, либо перестраивать её динамически, в зависимости от логики программы.
// Клиент может отправлять запросы любому из объектов цепочки, не обязательно первому из них.
type patient struct {
	name              string
	registrationDone  bool
	doctorCheckUpDone bool
	medicineDone      bool
	paymentDone       bool
}

func ChainPatternStart() {
	// Кассир
	cashier := &cashier{}

	// Комната медикаментов
	medical := &medical{}
	medical.setNext(cashier)

	// Доктор
	doctor := &doctor{}
	doctor.setNext(medical)

	// Регистратура
	reception := &reception{}
	reception.setNext(doctor)

	patient := &patient{name: "Alex"}
	//Patient visiting
	reception.execute(patient)
}
