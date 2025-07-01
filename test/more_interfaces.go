package test_data

import (
	"fmt"
	"io"
)

// ===============================
// СЦЕНАРИИ ДЛЯ more_interfaces.go
// ===============================

// Кейс 23: Вложенные и расширяющие интерфейсы
type BaseIO interface {
	Read(p []byte) (n int, err error) // используется (через io.Reader)
	Close() error                     // используется
}

type ExtendedIO interface {
	BaseIO
	Seek(offset int64, whence int) (int64, error) // не используется
}

// Кейс 24: Интерфейс с методом, перекрывающим встроенный тип
type Stringer interface {
	String() string // используется (стандартный интерфейс)
}

type CustomStringer interface {
	String(format string) string // не используется
}

// Кейс 25: Цепочка интерфейсов
type First interface {
	FirstMethod() // используется
}

type Second interface {
	First
	SecondMethod() // не используется
}

type Third interface {
	Second
	ThirdMethod() // используется
}

// Кейс 26: Интерфейсы, используемые как анонимные поля в структурах
type Greeter interface {
	Greet() string // используется
}

type Speaker interface {
	Speak() // не используется
}

// Кейс 27: Использование метода через присваивание другому интерфейсу
type Source interface {
	ReadSource() string // используется
}

type Destination interface {
	ReadSource() string // используется (тот же метод)
}

// ===============================
// СТРУКТУРЫ И РЕАЛИЗАЦИИ
// ===============================

// Для Кейса 23
type File struct{}

func (f *File) Read(p []byte) (n int, err error) { return 0, nil }
func (f *File) Close() error                     { return nil }
func (f *File) Seek(offset int64, whence int) (int64, error) {
	return 0, nil
}

// Для Кейса 24
type DataItem struct{}

func (d DataItem) String() string                   { return "data" }
func (d DataItem) StringWithFormat(f string) string { return "formatted_data" }

// Для Кейса 25
type ChainImpl struct{}

func (c *ChainImpl) FirstMethod()  {}
func (c *ChainImpl) SecondMethod() {}
func (c *ChainImpl) ThirdMethod()  {}

// Для Кейса 26
type GreeterImpl struct{}

func (g *GreeterImpl) Greet() string { return "hello" }
func (g *GreeterImpl) Speak()        {}

type Robot struct {
	Greeter // используется
	Speaker // не используется
}

// Для Кейса 27
type DataSource struct{}

func (ds *DataSource) ReadSource() string { return "source data" }

// ===============================
// ФУНКЦИИ ИСПОЛЬЗОВАНИЯ
// ===============================

func UseMoreInterfaces() {
	// Кейс 23: Вложенные интерфейсы
	var closer io.Closer = &File{}
	closer.Close() // Использует BaseIO.Close

	var reader io.Reader = &File{}
	reader.Read(nil) // Использует BaseIO.Read

	// Кейс 24: Неявный вызов fmt.Stringer
	var s Stringer = DataItem{}
	fmt.Println(s) // Использует Stringer.String

	// Кейс 25: Цепочка интерфейсов
	var t Third = &ChainImpl{}
	t.FirstMethod()
	t.ThirdMethod()

	// Кейс 26: Анонимные поля
	robot := Robot{Greeter: &GreeterImpl{}}
	robot.Greet() // Использует Greeter.Greet

	// Кейс 27: Присваивание интерфейсов
	var src Source = &DataSource{}
	var dst Destination = src // Присваивание, которое делает метод используемым для обоих
	dst.ReadSource()
}
