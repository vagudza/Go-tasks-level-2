package pattern

import (
	"fmt"
)

/*
	Реализовать паттерн «фабричный метод».
Объяснить применимость паттерна, его плюсы и минусы, а также реальные примеры использования данного примера на практике.
	https://en.wikipedia.org/wiki/Factory_method_pattern



	Фабричный метод (Factory method) также известный как Виртуальный конструктор (Virtual Constructor) - пораждающий шаблон проектирования,
	определяющий общий интерфейс создания объектов в родительском классе и позволяющий изменять создаваемые объекты в дочерних классах.

	Шаблон позволяет классу делегировать создание объектов подклассам. Используется, когда:
    	-Классу заранее неизвестно, объекты каких подклассов ему нужно создать.
    	-Обязанности делегируются подклассу, а знания о том, какой подкласс принимает эти обязанности, локализованы.
    	-Создаваемые объекты родительского класса специализируются подклассами.


В примере реализации показано, как обеспечить хранилище данных с различными бэкэндами, такими как хранилище в памяти, дисковое хранилище.
*/

type Store interface {
	Save(string) error
}

// Различные реализации
type StorageType int

const (
	mongoStorage StorageType = 1 << iota
	memoryStorage
)

func NewStore(t StorageType) Store {
	switch t {
	case memoryStorage:
		return newMemoryStorage()
	case mongoStorage:
		return newMongoStorage()
	default:
		fmt.Println("unknown Storage type")
		return nil
	}
}

type MemoryStorage struct{}

func (ms *MemoryStorage) Save(s string) error {
	// сохранение записи в ОЗУ
	fmt.Printf("Запись '%s' успешно сохранена в Монго\n", s)
	return nil
}

func newMemoryStorage() *MemoryStorage {
	return &MemoryStorage{}
}

type MongoStorage struct{}

func (ms *MongoStorage) Save(s string) error {
	// сохранение записи в БД
	fmt.Printf("Запись '%s' успешно сохранена в Монго\n", s)
	return nil
}

func newMongoStorage() *MongoStorage {
	return &MongoStorage{}
}

// Использование
// С фабричным методом пользователь может определить тип хранилища, который он хочет.
func FactoryMethodPatternStart() {
	memStorage := NewStore(memoryStorage)
	memStorage.Save("memory")

	monStorage := NewStore(mongoStorage)
	monStorage.Save("Mongo")
}
