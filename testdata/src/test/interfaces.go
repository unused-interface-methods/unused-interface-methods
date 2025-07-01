package test_data

import (
	"context"
	"io"
)

// ===============================
// БАЗОВЫЕ КЕЙСЫ 1-16
// ===============================

// Кейс 1: Встроенные интерфейсы (embedded interfaces)
type Reader interface {
	io.Reader
	CustomRead() error // используется через type assertion
}

type Writer interface {
	io.Writer
	CustomWrite() error // используется
}

// Кейс 2: Интерфейс с контекстом
type ServiceWithContext interface {
	ProcessWithContext(ctx context.Context, data string) error // используется
	CancelOperation(ctx context.Context) error                 // используется в type switch
}

// Кейс 3: Интерфейс с вариативными параметрами
type Logger interface {
	Log(level string, args ...interface{}) error // используется
	Debug(args ...string)                        // используется как значение функции
}

// Кейс 4: Интерфейс с функциональными типами
type EventHandler interface {
	OnEvent(callback func(string) error) error           // используется
	OnError(handler func(error) bool) error              // want "method \"OnError\" of interface \"EventHandler\" is declared but not used"
	Subscribe(filter func(string) bool, cb func()) error // want "method \"Subscribe\" of interface \"EventHandler\" is declared but not used"
}

// Кейс 5: Интерфейс с каналами
type ChannelProcessor interface {
	SendData(ch chan<- string) error    // используется
	ReceiveData(ch <-chan string) error // want "method \"ReceiveData\" of interface \"ChannelProcessor\" is declared but not used"
	ProcessStream(ch chan string) error // want "method \"ProcessStream\" of interface \"ChannelProcessor\" is declared but not used"
}

// Кейс 6: Интерфейс с мапами и слайсами
type DataProcessor interface {
	ProcessMap(data map[string]interface{}) error // используется
	ProcessSlice(data []string) error             // want "method \"ProcessSlice\" of interface \"DataProcessor\" is declared but not used"
	ProcessArray(data [10]int) error              // want "method \"ProcessArray\" of interface \"DataProcessor\" is declared but not used"
}

// Кейс 7: Интерфейс с указателями
type PointerHandler interface {
	HandlePointer(data *string) error // используется
	HandleDoublePtr(data **int) error // want "method \"HandleDoublePtr\" of interface \"PointerHandler\" is declared but not used"
}

// Кейс 8: Интерфейс с именованными возвращаемыми значениями
type NamedReturns interface {
	GetNamedResult() (result string, err error)            // используется
	GetMultipleNamed() (x, y int, success bool, err error) // want "method \"GetMultipleNamed\" of interface \"NamedReturns\" is declared but not used"
}

// Кейс 9: Интерфейсы с одинаковыми именами методов
type ProcessorV1 interface {
	Process(data string) error // используется
}

type ProcessorV2 interface {
	Process(data string, options map[string]interface{}) error // want "method \"Process\" of interface \"ProcessorV2\" is declared but not used"
}

// Кейс 10: Интерфейс с методами без параметров
type SimpleActions interface {
	Start()          // используется
	Stop()           // используется в горутине
	Reset()          // используется в defer и reflection
	GetStatus() bool // используется
}

// Кейс 11: Интерфейс с интерфейсными параметрами
type InterfaceParams interface {
	HandleReader(r io.Reader) error        // используется
	HandleWriter(w io.Writer) error        // want "method \"HandleWriter\" of interface \"InterfaceParams\" is declared but not used"
	HandleBoth(rw io.ReadWriter) error     // want "method \"HandleBoth\" of interface \"InterfaceParams\" is declared but not used"
	HandleCustom(custom interface{}) error // want "method \"HandleCustom\" of interface \"InterfaceParams\" is declared but not used"
}

// Кейс 12: Расширенный интерфейс с разными типами параметров
type AdvancedInterface interface {
	// Простые типы
	HandleString(s string) error  // используется
	HandleInt(i int) bool         // используется
	HandleFloat(f float64) string // используется

	// Указатели
	HandlePointer(p *string) error // используется
	HandleNilPointer() *int        // используется

	// Слайсы и массивы
	HandleSlice(slice []string) error // используется
	HandleArray(arr [5]int) error     // используется

	// Мапы
	HandleMap(m map[string]int) error          // используется
	HandleComplexMap(m map[string][]int) error // используется

	// Каналы
	HandleChannel(ch chan string) error        // используется
	HandleReadChannel(ch <-chan string) error  // используется
	HandleWriteChannel(ch chan<- string) error // используется

	// Функциональные типы
	HandleFunc(f func(string) error) error                     // используется
	HandleComplexFunc(f func(int, string) (bool, error)) error // используется

	// Интерфейсы
	HandleInterface(i interface{}) error     // используется
	HandleContext(ctx context.Context) error // используется

	// Вариативные параметры
	HandleVariadic(args ...string) error                  // используется
	HandleMixedVariadic(prefix string, args ...int) error // используется

	// Без параметров
	GetResult() bool // используется
	Clear()          // используется

	// Множественные возвращаемые значения
	GetMultiple() (string, int, error)              // используется
	GetNamedReturns() (result string, success bool) // используется
}

// Кейс 13: Интерфейс с методом, имеющим такое же имя
type AnotherReader interface {
	CustomRead(data []byte) error // want "method \"CustomRead\" of interface \"AnotherReader\" is declared but not used"
}

// Кейс 14: Интерфейс с interface{} параметрами
type AnyProcessor interface {
	Process(data interface{}) error // используется
}

type AnyHandler interface {
	Process(data interface{}, options interface{}) error // want "method \"Process\" of interface \"AnyHandler\" is declared but not used"
}

// Кейс 15: Интерфейсы для проверки вызовов на интерфейсных переменных
type StringProcessor interface {
	Process(data string) error // используется через интерфейсную переменную
}

type StringHandler interface {
	Process(data string) error // не используется, но имеет такую же сигнатуру
}

type StringConsumer interface {
	Process(data string) error // используется через интерфейсную переменную
}

// Кейс 16: Интерфейсы для проверки вызовов с разными сигнатурами
type BaseProcessor interface {
	Process() error // базовый метод
}

type ExtendedProcessor interface {
	Process() error // тот же метод
	Extra() bool    // дополнительный метод
}

// Структура, реализующая оба интерфейса
type ProcessorImpl struct{}

func (p *ProcessorImpl) Process() error {
	return nil
}

func (p *ProcessorImpl) Extra() bool {
	return true
}

// ===============================
// СТРУКТУРЫ И ИСПОЛЬЗОВАНИЕ
// ===============================

// TestPointerStruct - структура для тестирования прямых вызовов методов
// Она нужна для проверки, что линтер правильно различает:
// 1. Вызовы через интерфейс (cs.pointer.HandlePointer(&str))
// 2. Прямые вызовы на структуре (s := &TestPointerStruct{}; s.HandlePointer(&str))
type TestPointerStruct struct{}

func (t *TestPointerStruct) HandlePointer(data *string) error {
	return nil
}

func (t *TestPointerStruct) HandleDoublePtr(data **int) error {
	return nil
}

// Основная структура, использующая все интерфейсы
type ComplexService struct {
	Reader      Reader
	Writer      Writer
	Service     ServiceWithContext
	Logger      Logger
	Handler     EventHandler
	Processor   ChannelProcessor
	Data        DataProcessor
	Pointer     PointerHandler
	Named       NamedReturns
	ProcessorV1 ProcessorV1
	ProcessorV2 ProcessorV2
	Actions     SimpleActions
	Params      InterfaceParams
	Advanced    AdvancedInterface
	AnyProc     AnyProcessor
	AnyHandler  AnyHandler
	StrProc     StringProcessor
	StrHandler  StringHandler
	StrConsumer StringConsumer
	BaseProc    BaseProcessor
	ExtProc     ExtendedProcessor
	Direct      DirectProcessor
}

// ===============================
// ИНТЕРФЕЙСЫ ДЛЯ ТЕСТИРОВАНИЯ ВЫЗОВОВ МЕТОДОВ
// ===============================

type EmptyInterface interface{}

type SingleMethodInterface interface {
	Method() // want "method \"Method\" of interface \"SingleMethodInterface\" is declared but not used"
}

type MultiMethodInterface interface {
	Method1() // want "method \"Method1\" of interface \"MultiMethodInterface\" is declared but not used"
	Method2() // want "method \"Method2\" of interface \"MultiMethodInterface\" is declared but not used"
}

type AssignableInterface1 interface {
	Method() // want "method \"Method\" of interface \"AssignableInterface1\" is declared but not used"
	Extra()  // want "method \"Extra\" of interface \"AssignableInterface1\" is declared but not used"
}

type AssignableInterface2 interface {
	Method() // want "method \"Method\" of interface \"AssignableInterface2\" is declared but not used"
}

// DirectProcessor - интерфейс для тестирования прямой реализации
type DirectProcessor interface {
	ProcessDirect(data string) error
}

// Структуры для тестирования
type TestStruct struct{}

func (t *TestStruct) Method()  {}
func (t *TestStruct) Method1() {}
func (t *TestStruct) Method2() {}
func (t *TestStruct) Extra()   {}

// DirectProcessorImpl - структура, напрямую реализующая интерфейс
type DirectProcessorImpl struct{}

func (p *DirectProcessorImpl) ProcessDirect(data string) error {
	return nil
}

// DoComplexWork использует различные методы интерфейсов
func (cs *ComplexService) DoComplexWork() {
	// 1. Использование встроенных интерфейсов
	cs.Writer.CustomWrite()

	// 2. Использование интерфейса с контекстом
	ctx := context.Background()
	cs.Service.ProcessWithContext(ctx, "data")

	// 3. Использование интерфейса с вариативными параметрами
	cs.Logger.Log("info", "message", "details")

	// 4. Использование интерфейса с функциональными типами
	cs.Handler.OnEvent(func(s string) error { return nil })

	// 5. Использование интерфейса с каналами
	ch := make(chan string)
	cs.Processor.SendData(ch)

	// 6. Использование интерфейса с мапами и слайсами
	cs.Data.ProcessMap(map[string]interface{}{"key": "value"})

	// 7. Использование интерфейса с указателями
	str := "test"
	cs.Pointer.HandlePointer(&str)

	// 8. Прямой вызов метода на структуре (для тестирования различий между вызовами)
	s := &TestPointerStruct{}
	s.HandlePointer(&str)

	// 9. Использование интерфейса с именованными возвращаемыми значениями
	cs.Named.GetNamedResult()

	// 10. Использование интерфейса с одинаковыми именами методов
	cs.ProcessorV1.Process("data")

	// 11. Использование интерфейса с методами без параметров
	cs.Actions.Start()
	cs.Actions.GetStatus()

	// 12. Использование интерфейса с интерфейсными параметрами
	cs.Params.HandleReader(cs.Reader)

	// 13. Использование расширенного интерфейса
	cs.Advanced.HandleString("test")
	cs.Advanced.HandlePointer(&str)
	cs.Advanced.HandleSlice([]string{"test"})
	cs.Advanced.HandleMap(map[string]int{"key": 42})
	cs.Advanced.HandleChannel(ch)
	cs.Advanced.HandleFunc(func(s string) error { return nil })
	cs.Advanced.HandleInterface(cs.Reader)
	cs.Advanced.HandleContext(ctx)
	cs.Advanced.HandleVariadic("test1", "test2")
	cs.Advanced.GetResult()
	cs.Advanced.GetMultiple()

	// Прямая реализация интерфейса
	cs.Direct.ProcessDirect("test")
}

// ===============================
// СПЕЦИАЛЬНЫЕ КЕЙСЫ 17-22
// ===============================

// Кейс 17: Использование через type assertion
func (cs *ComplexService) TypeAssertions() {
	var iface interface{} = cs.Reader

	if r, ok := iface.(Reader); ok {
		r.CustomRead() // CustomRead используется через type assertion
	}

	if adv, ok := interface{}(cs.Advanced).(AdvancedInterface); ok { //nolint:staticcheck // намеренно для тестов
		adv.HandleInt(42)     // HandleInt используется через type assertion
		adv.HandleFloat(3.14) // HandleFloat используется через type assertion
	}
}

// Кейс 18: Использование в горутинах
func (cs *ComplexService) GoroutineUsage() {
	go func() {
		cs.Actions.Stop()              // Stop используется в горутине
		cs.Advanced.HandleNilPointer() // HandleNilPointer используется в горутине
	}()
}

// Кейс 19: Методы как значения функций
func (cs *ComplexService) MethodValues() {
	// Сохраняем методы как значения
	arrayHandler := cs.Advanced.HandleArray
	mapHandler := cs.Advanced.HandleComplexMap
	debugFunc := cs.Logger.Debug

	// Используем их
	arrayHandler([5]int{1, 2, 3, 4, 5})            // HandleArray используется как значение
	mapHandler(map[string][]int{"key": {1, 2, 3}}) // HandleComplexMap используется как значение
	debugFunc("test message")                      // Debug используется как значение функции
}

// Кейс 20: Использование в switch/type switch
func (cs *ComplexService) SwitchUsage(value interface{}) {
	switch v := value.(type) {
	case AdvancedInterface:
		// Эти методы используются в type switch
		v.HandleReadChannel(make(<-chan string))                                  // HandleReadChannel используется
		v.HandleWriteChannel(make(chan<- string))                                 // HandleWriteChannel используется
		v.HandleComplexFunc(func(int, string) (bool, error) { return true, nil }) // HandleComplexFunc используется
	case ServiceWithContext:
		v.CancelOperation(context.Background()) // CancelOperation используется в type switch
	}
}

// Кейс 21: Использование в defer
func (cs *ComplexService) DeferUsage() {
	defer cs.Advanced.GetNamedReturns() // GetNamedReturns используется в defer
	defer cs.Actions.Reset()            // Reset используется в defer

	defer func() {
		cs.Advanced.HandleMixedVariadic("prefix", 1, 2, 3) // HandleMixedVariadic используется в defer
	}()
}

// Кейс 22: Использование в reflection/interface{}
func (cs *ComplexService) ReflectionUsage() {
	// Сложный случай: метод вызывается через рефлексию
	// Наш линтер может не поймать такое использование
	var handlers []interface{} = []interface{}{
		cs.Actions.Reset,  // Reset используется через рефлексию
		cs.Advanced.Clear, // Clear используется через рефлексию
	}
	_ = handlers

	// Использование методов через interface{}
	var anyReader interface{} = cs.Reader
	if reader, ok := anyReader.(io.Reader); ok {
		reader.Read(make([]byte, 10)) // Read используется через interface{}
	}
	if reader, ok := anyReader.(Reader); ok {
		reader.CustomRead() // CustomRead используется через interface{}
	}

	// Использование методов с interface{} параметрами
	var processor interface{} = cs.ProcessorV1
	if p, ok := processor.(AnyProcessor); ok {
		p.Process("data") // Process используется через interface{}
	}

	// Использование методов через интерфейсные переменные
	var sp StringProcessor = cs.ProcessorV1
	var sc StringConsumer = sp // присваиваем переменной типа StringConsumer
	sc.Process("data")         // Process используется через интерфейсную переменную

	// Использование методов через интерфейсные переменные с разными сигнатурами
	impl := &ProcessorImpl{}
	var base BaseProcessor = impl
	var extended ExtendedProcessor = impl
	base.Process()     // Process используется через базовый интерфейс
	extended.Process() // Process используется через расширенный интерфейс
	extended.Extra()   // Extra используется только через расширенный интерфейс
}
