package pattern

import "fmt"

/*
	Реализовать паттерн «строитель».
Объяснить применимость паттерна, его плюсы и минусы, а также реальные примеры использования данного примера на практике.
	https://en.wikipedia.org/wiki/Builder_pattern


Строитель (англ. Builder) — порождающий шаблон проектирования предоставляет способ создания составного объекта.
Отделяет конструирование сложного объекта от его представления так, что в результате одного и того
же процесса конструирования могут получаться разные представления.

Проблема:
    Инициализация очень сложного, большого объекта со множеством параметров инициализации. Использовать один конструктор
    с множеством параметров - плохо (телескопический конструктор - анти-паттерн)

Решение:
    Паттерн Строитель предлагает вынести конструирование объекта за пределы его собственного класса, поручив это дело
    отдельным объектам, называемым строителями.

плюсы:
    +Позволяет создавать продукты пошагово.
    +Позволяет использовать один и тот же код для создания различных продуктов.
    +Изолирует сложный код сборки продукта от его основной бизнес-логики.

минусы:
    -Усложняет код программы из-за введения дополнительных классов.
    -Клиент будет привязан к конкретным классам строителей, так как в интерфейсе директора может не быть метода получения результата.

    В качесте примера реализации создадим услугу "Бессерверные вычисления". Пользователь сам выбирает настройки облака, и Builder
    подготавливает для него объект с настроенными параметрами. 2 услуги: облачные вычисления и облачное хранилище
*/

// Интерфейс строителя
type iBuilder interface {
	setCPU()
	setRAM()
	setStorageType()
	setCore()
	getService() *Service
}

func getBuilder(builderType string) iBuilder {
	if builderType == "cloud computing" {
		return &CloudComputeBuilder{}
	}

	if builderType == "object storage" {
		return &ObjectStorageBuilder{}
	}
	return nil
}

// конкретный строитель 1
type CloudComputeBuilder struct {
	CPUsNum     int
	RAM         int
	StorageSize int
	StorageType string
	Core        string
}

func (cc *CloudComputeBuilder) setCPU() {
	cc.CPUsNum = 8
}

func (cc *CloudComputeBuilder) setRAM() {
	cc.RAM = 64
}

func (cc *CloudComputeBuilder) setStorageType() {
	cc.StorageType = "SSD"
	cc.StorageSize = 256
}

func (cc *CloudComputeBuilder) setCore() {
	cc.Core = "Ubuntu LTS 18.01"
}

func (cc *CloudComputeBuilder) getService() *Service {
	return &Service{
		CPUsNum:     cc.CPUsNum,
		RAM:         cc.RAM,
		StorageSize: cc.StorageSize,
		StorageType: cc.StorageType,
		Core:        cc.Core,
	}
}

// конкретный строитель 2
type ObjectStorageBuilder struct {
	CPUsNum     int
	RAM         int
	StorageSize int
	StorageType string
	Core        string
}

func (objs *ObjectStorageBuilder) setCPU() {
	objs.CPUsNum = 2
}

func (objs *ObjectStorageBuilder) setRAM() {
	objs.RAM = 16
}

func (objs *ObjectStorageBuilder) setStorageType() {
	objs.StorageType = "HDD"
	objs.StorageSize = 32768
}

func (objs *ObjectStorageBuilder) setCore() {
	objs.Core = "Debian"
}

func (objs *ObjectStorageBuilder) getService() *Service {
	return &Service{
		CPUsNum:     objs.CPUsNum,
		RAM:         objs.RAM,
		StorageSize: objs.StorageSize,
		StorageType: objs.StorageType,
		Core:        objs.Core,
	}
}

// продукт строительства
type Service struct {
	CPUsNum     int
	RAM         int
	StorageSize int
	StorageType string
	Core        string
}

func (s *Service) printSettings() {
	fmt.Printf("Cloud Computing Service\nCPUs:%d\nRAM:%d\nStorage:%s\n", s.CPUsNum, s.RAM, s.StorageType)
	fmt.Printf("Storage size:%d\nCore:%s\n", s.StorageSize, s.Core)
}

// Директор
type Director struct {
	builder iBuilder
}

func newDirector(b iBuilder) *Director {
	return &Director{
		builder: b,
	}
}

func (d *Director) setBuilder(b iBuilder) {
	d.builder = b
}

// процесс построения сервиса конкретным строителем
func (d *Director) buildService() *Service {
	d.builder.setCPU()
	d.builder.setRAM()
	d.builder.setStorageType()
	d.builder.setCore()
	return d.builder.getService()
}

// клиентское применение:
func BuilderPatternStart() {
	// создаем работников:
	cloudComputingBuilder := getBuilder("cloud computing")
	objectStorageBuilder := getBuilder("object storage")

	fmt.Printf("\n\n---Pattern Builder---\nCreating new cloud service for cloud computing:\n")
	director := newDirector(cloudComputingBuilder)
	cloudComputingService := director.buildService()
	cloudComputingService.printSettings()

	fmt.Printf("\nCreating new cloud service for storage data:\n")
	director.setBuilder(objectStorageBuilder)
	objectStorageService := director.buildService()
	objectStorageService.printSettings()
}
