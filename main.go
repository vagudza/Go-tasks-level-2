package main

func main() {
	// Patterns

	// 1. Facade
	//vm := pattern.NewVendingMachine()
	//vm.Start()

	// 2. Builder
	//pattern.BuilderPatternStart()

	// 3. Visitor
	//pattern.VisitorPatternStart()

	// 4. Command
	//pattern.CommandPatternStart()

	// 5. Chain of responsibility
	//pattern.ChainPatternStart()

	// 6. Factory
	//pattern.FactoryMethodPatternStart()

	// 7. Strategy
	//pattern.StrategyPatternStart()

	// 8. State
	//pattern.StatePatternStart()

	// Listing task 1
	//listing.Task1()

	// Listing task 2
	//fmt.Println(listing.Test())
	//fmt.Println(listing.AnotherTest())

	// Listing task 3
	//listing.Task3()

	// Listing task 4
	//listing.Task4()

	// Listing task 5
	//listing.Task5()

	// Listing task 6
	//listing.Task6()

	// Listing task 7
	//listing.Task7()
}

/*
Список вопросов:
	1. Если в реализации паттерна Фасад имеется хоть малейшая логика (например, используется оператор if),
	то это по-прежнему паттерн Фасад? см. 01_facade.go, функция Start()

	2. Паттерн Builder решает проблему инициализации сложного объекта. Я привел пример простой настройки параметров.
	Насколько я понимаю, для такого примитивной настройки можно просто создать один конструктор, который принимает
	путь к файлу с конфигурацией и просто считывает все параметры. Правильно ли я понимаю, что Builder практически
	применим для инициализации объекта, при которой должна выполняться некоторая логика, а не просто сохранение параметров?
	Например, в конкретном строителе использовать условные операторы, циклы, и другую логику.

	3. Почему в https://refactoring.guru/ru/design-patterns/visitor/go/example написано, что
	"мы не можем оставить только один метод visit(shape) в интерфейсе посетителя? Это невозможно из-за того,
	что язык Go не поддерживает перегрузку методов, поэтому вы не можете иметь методы с одинаковыми именами,
	но разными параметрами". Я реализовал такой метод в паттерне visitor (90 строка)

	4. при работе с реальными случаями возникают серьезные ограничения, все из-за ограничений языка:
	Реализация паттерна "посетитель" работает только в том случае, если и интерфейс посетителя, и структуры находятся в одном пакете,
	чтобы избежать циклов импорта (они оба нуждаются друг в друге). https://gist.github.com/francoishill/f0624e7760aacdc96b42

	5. В посетителе поля посещаемых структур должны быть не приватными. Иначе при посещении из другого пакета мы не сможем
	получить доступ к неэкспортируемым полям

*/
