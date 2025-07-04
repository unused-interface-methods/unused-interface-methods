package test

import (
	"context"
	"fmt"
	"io"
)

// Part 1-3

// Case 1: Built-in interfaces (embedded interfaces)
type Reader interface {
	io.Reader
	CustomRead() error // used via type assertion
}

type Writer interface {
	io.Writer
	CustomWrite() error // used
}

// Case 2: Interface with context
type ServiceWithContext interface {
	ProcessWithContext(ctx context.Context, data string) error // used
	CancelOperation(ctx context.Context) error                 // used in type switch
}

// Case 3: Interface with variadic parameters
type Logger interface {
	Log(level string, args ...interface{}) error // used
	Debug(args ...string)                        // used as function value
}

// Case 4: Interface with functional types
type EventHandler interface {
	OnEvent(callback func(string) error) error           // used
	OnError(handler func(error) bool) error              // want "method \"OnError\" of interface \"EventHandler\" is declared but not used"
	Subscribe(filter func(string) bool, cb func()) error // want "method \"Subscribe\" of interface \"EventHandler\" is declared but not used"
}

// Case 5: Interface with channels
type ChannelProcessor interface {
	SendData(ch chan<- string) error    // used
	ReceiveData(ch <-chan string) error // want "method \"ReceiveData\" of interface \"ChannelProcessor\" is declared but not used"
	ProcessStream(ch chan string) error // want "method \"ProcessStream\" of interface \"ChannelProcessor\" is declared but not used"
}

// Case 6: Interface with maps and slices
type DataProcessor interface {
	ProcessMap(data map[string]interface{}) error // used
	ProcessSlice(data []string) error             // want "method \"ProcessSlice\" of interface \"DataProcessor\" is declared but not used"
	ProcessArray(data [10]int) error              // want "method \"ProcessArray\" of interface \"DataProcessor\" is declared but not used"
}

// Case 7: Interface with pointers
type PointerHandler interface {
	HandlePointer(data *string) error // used
	HandleDoublePtr(data **int) error // want "method \"HandleDoublePtr\" of interface \"PointerHandler\" is declared but not used"
}

// Case 8: Interface with named return values
type NamedReturns interface {
	GetNamedResult() (result string, err error)            // used
	GetMultipleNamed() (x, y int, success bool, err error) // want "method \"GetMultipleNamed\" of interface \"NamedReturns\" is declared but not used"
}

// Case 9: Interfaces with the same method names
type ProcessorV1 interface {
	Process(data string) error // used
}

type ProcessorV2 interface {
	Process(data string, options map[string]interface{}) error // want "method \"Process\" of interface \"ProcessorV2\" is declared but not used"
}

// Case 10: Interface with methods without parameters
type SimpleActions interface {
	Start()          // used
	Stop()           // used in goroutine
	Reset()          // used in defer and reflection
	GetStatus() bool // used
}

// Case 11: Interface with interface parameters
type InterfaceParams interface {
	HandleReader(r io.Reader) error        // used
	HandleWriter(w io.Writer) error        // want "method \"HandleWriter\" of interface \"InterfaceParams\" is declared but not used"
	HandleBoth(rw io.ReadWriter) error     // want "method \"HandleBoth\" of interface \"InterfaceParams\" is declared but not used"
	HandleCustom(custom interface{}) error // want "method \"HandleCustom\" of interface \"InterfaceParams\" is declared but not used"
}

// Case 12: Extended interface with different types of parameters
type AdvancedInterface interface {
	// Simple types
	HandleString(s string) error  // used
	HandleInt(i int) bool         // used
	HandleFloat(f float64) string // used

	// Pointers
	HandlePointer(p *string) error // used
	HandleNilPointer() *int        // used

	// Slices and arrays
	HandleSlice(slice []string) error // used
	HandleArray(arr [5]int) error     // used

	// Maps
	HandleMap(m map[string]int) error          // used
	HandleComplexMap(m map[string][]int) error // used

	// Channels
	HandleChannel(ch chan string) error        // used
	HandleReadChannel(ch <-chan string) error  // used
	HandleWriteChannel(ch chan<- string) error // used

	// Function types
	HandleFunc(f func(string) error) error                     // used
	HandleComplexFunc(f func(int, string) (bool, error)) error // used

	// Interfaces
	HandleInterface(i interface{}) error     // used
	HandleContext(ctx context.Context) error // used

	// Variadic parameters
	HandleVariadic(args ...string) error                  // used
	HandleMixedVariadic(prefix string, args ...int) error // used

	// Multiple return values
	GetResult() bool // used
	Clear()          // used

	// GetNamedReturns - method with the same name
	GetNamedReturns() (result string, success bool) // used
}

// Case 13: Interface with method having the same name
type AnotherReader interface {
	CustomRead(data []byte) error // want "method \"CustomRead\" of interface \"AnotherReader\" is declared but not used"
}

// Case 14: Interface with interface{} parameters
type AnyProcessor interface {
	Process(data interface{}) error // used
}

type AnyHandler interface {
	Process(data interface{}, options interface{}) error // want "method \"Process\" of interface \"AnyHandler\" is declared but not used"
}

// Case 15: Interfaces for checking calls on interface variables
type StringProcessor interface {
	Process(data string) error // used via interface variable
}

type StringHandler interface {
	Process(data string) error // want "method \"Process\" of interface \"StringHandler\" is declared but not used"
}

type StringConsumer interface {
	Process(data string) error // used via interface variable
}

// Case 16: Interfaces for checking calls with different signatures
type BaseProcessor interface {
	Process() error // base method
}

type ExtendedProcessor interface {
	Process() error // same method
	Extra() bool    // additional method
}

// Structure implementing both interfaces
type ProcessorImpl struct{}

func (p *ProcessorImpl) Process() error {
	return nil
}

func (p *ProcessorImpl) Extra() bool {
	return true
}

// ===============================
// STRUCTURES AND USAGE
// ===============================

// TestPointerStruct - structure for testing direct method calls
// It is needed to verify that the linter correctly distinguishes:
// 1. Calls through interface (cs.pointer.HandlePointer(&str))
// 2. Direct calls on structure (s := &TestPointerStruct{}; s.HandlePointer(&str))
type TestPointerStruct struct{}

func (t *TestPointerStruct) HandlePointer(data *string) error {
	return nil
}

func (t *TestPointerStruct) HandleDoublePtr(data **int) error {
	return nil
}

// Main structure using all interfaces
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
// INTERFACES FOR TESTING METHOD CALLS
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

// DirectProcessor - interface for testing direct implementation
type DirectProcessor interface {
	ProcessDirect(data string) error
}

// Structures for testing
type TestStruct struct{}

func (t *TestStruct) Method()  {}
func (t *TestStruct) Method1() {}
func (t *TestStruct) Method2() {}
func (t *TestStruct) Extra()   {}

// DirectProcessorImpl - structure directly implementing the interface
type DirectProcessorImpl struct{}

func (p *DirectProcessorImpl) ProcessDirect(data string) error {
	return nil
}

// DoComplexWork uses various methods of interfaces
func (cs *ComplexService) DoComplexWork() {
	// 1. Using built-in interfaces
	cs.Writer.CustomWrite()

	// 2. Using interface with context
	ctx := context.Background()
	cs.Service.ProcessWithContext(ctx, "data")

	// 3. Using interface with variadic parameters
	cs.Logger.Log("info", "message", "details")

	// 4. Using interface with functional types
	cs.Handler.OnEvent(func(s string) error { return nil })

	// 5. Using interface with channels
	ch := make(chan string)
	cs.Processor.SendData(ch)

	// 6. Using interface with maps and slices
	cs.Data.ProcessMap(map[string]interface{}{"key": "value"})

	// 7. Using interface with pointers
	str := "test"
	cs.Pointer.HandlePointer(&str)

	// 8. Direct method call on structure (for testing differences between calls)
	s := &TestPointerStruct{}
	s.HandlePointer(&str)

	// 9. Using interface with named return values
	cs.Named.GetNamedResult()

	// 10. Using interface with the same method names
	cs.ProcessorV1.Process("data")

	// 11. Using interface with methods without parameters
	cs.Actions.Start()
	cs.Actions.GetStatus()

	// 12. Using interface with interface parameters
	cs.Params.HandleReader(cs.Reader)

	// 13. Using extended interface
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
	cs.Advanced.GetNamedReturns()

	// Direct implementation of interface
	cs.Direct.ProcessDirect("test")
}

// Part 2-3

// Case 17: Usage through type assertion
func (cs *ComplexService) TypeAssertions() {
	var iface interface{} = cs.Reader

	if r, ok := iface.(Reader); ok {
		r.CustomRead() // CustomRead used via type assertion
	}

	if adv, ok := interface{}(cs.Advanced).(AdvancedInterface); ok { //nolint:staticcheck // intentional for tests
		adv.HandleInt(42)     // HandleInt used via type assertion
		adv.HandleFloat(3.14) // HandleFloat used via type assertion
	}
}

// Case 18: Usage in goroutines
func (cs *ComplexService) GoroutineUsage() {
	go func() {
		cs.Actions.Stop()              // Stop used in goroutine
		cs.Advanced.HandleNilPointer() // HandleNilPointer used in goroutine
	}()
}

// Case 19: Methods as function values
func (cs *ComplexService) MethodValues() {
	// Save methods as values
	arrayHandler := cs.Advanced.HandleArray
	mapHandler := cs.Advanced.HandleComplexMap
	debugFunc := cs.Logger.Debug

	// Use them
	arrayHandler([5]int{1, 2, 3, 4, 5})            // HandleArray used as value
	mapHandler(map[string][]int{"key": {1, 2, 3}}) // HandleComplexMap used as value
	debugFunc("test message")                      // Debug used as function value
}

// Case 20: Usage in switch/type switch
func (cs *ComplexService) SwitchUsage(value interface{}) {
	switch v := value.(type) {
	case AdvancedInterface:
		// These methods used in type switch
		v.HandleReadChannel(make(<-chan string))                                  // HandleReadChannel used
		v.HandleWriteChannel(make(chan<- string))                                 // HandleWriteChannel used
		v.HandleComplexFunc(func(int, string) (bool, error) { return true, nil }) // HandleComplexFunc used
	case ServiceWithContext:
		v.CancelOperation(context.Background()) // CancelOperation used in type switch
	}
}

// Case 21: Usage in defer
func (cs *ComplexService) DeferUsage() {
	defer cs.Advanced.GetNamedReturns() // GetNamedReturns used in defer
	defer cs.Actions.Reset()            // Reset used in defer

	defer func() {
		cs.Advanced.HandleMixedVariadic("prefix", 1, 2, 3) // HandleMixedVariadic used in defer
	}()
}

// Case 22: Usage in reflection/interface{}
func (cs *ComplexService) ReflectionUsage() {
	// Complex case: method called via reflection
	// Our linter may not catch such usage
	var handlers []interface{} = []interface{}{
		cs.Actions.Reset,  // Reset used via reflection
		cs.Advanced.Clear, // Clear used via reflection
	}
	_ = handlers

	// Using methods through interface{}
	var anyReader interface{} = cs.Reader
	if reader, ok := anyReader.(io.Reader); ok {
		reader.Read(make([]byte, 10)) // Read used via interface{}
	}
	if reader, ok := anyReader.(Reader); ok {
		reader.CustomRead() // CustomRead used via interface{}
	}

	// Using methods with interface{} parameters
	var processor interface{} = cs.ProcessorV1
	if p, ok := processor.(AnyProcessor); ok {
		p.Process("data") // Process used via interface{}
	}

	// Using methods through interface variables
	var sp StringProcessor = cs.ProcessorV1
	var sc StringConsumer = sp // assign to StringConsumer variable
	sc.Process("data")         // Process used via interface variable

	// Using methods through interface variables with different signatures
	impl := &ProcessorImpl{}
	var base BaseProcessor = impl
	var extended ExtendedProcessor = impl
	base.Process()     // Process used via base interface
	extended.Process() // Process used via extended interface
	extended.Extra()   // Extra used only via extended interface
}

// Part 3-3

// Case 23: Nested and extending interfaces
type BaseIO interface {
	Read(p []byte) (n int, err error) // used (through io.Reader)
	Close() error                     // used
}

type ExtendedIO interface {
	BaseIO
	Seek(offset int64, whence int) (int64, error) // want "method \"Seek\" of interface \"ExtendedIO\" is declared but not used"
}

// Case 24: Interface with method overriding built-in type
type Stringer interface {
	String() string // used (standard interface)
}

type CustomStringer interface {
	String(format string) string // want "method \"String\" of interface \"CustomStringer\" is declared but not used"
}

// Case 25: Chain of interfaces
type First interface {
	FirstMethod() // used
}

type Second interface {
	First
	SecondMethod() // want "method \"SecondMethod\" of interface \"Second\" is declared but not used"
}

type Third interface {
	Second
	ThirdMethod() // used
}

// Case 26: Interfaces used as anonymous fields in structures
type Greeter interface {
	Greet() string // used
}

type Speaker interface {
	Speak() // want "method \"Speak\" of interface \"Speaker\" is declared but not used"
}

// Case 27: Using method through assignment to another interface
type Source interface {
	ReadSource() string // used
}

type Destination interface {
	ReadSource() string // used (same method)
}

// ===============================
// STRUCTURES AND IMPLEMENTATIONS
// ===============================

// For Case 23
type File struct{}

func (f *File) Read(p []byte) (n int, err error) { return 0, nil }
func (f *File) Close() error                     { return nil }
func (f *File) Seek(offset int64, whence int) (int64, error) {
	return 0, nil
}

// For Case 24
type DataItem struct{}

func (d DataItem) String() string                   { return "data" }
func (d DataItem) StringWithFormat(f string) string { return "formatted_data" }

// For Case 25
type ChainImpl struct{}

func (c *ChainImpl) FirstMethod()  {}
func (c *ChainImpl) SecondMethod() {}
func (c *ChainImpl) ThirdMethod()  {}

// For Case 26
type GreeterImpl struct{}

func (g *GreeterImpl) Greet() string { return "hello" }
func (g *GreeterImpl) Speak()        {}

type Robot struct {
	Greeter // used
}

// For Case 27
type DataSource struct{}

func (ds *DataSource) ReadSource() string { return "source data" }

// ===============================
// USAGE FUNCTIONS
// ===============================

func UseMoreInterfaces() {
	// Case 23: Nested interfaces
	var closer io.Closer = &File{}
	closer.Close() // Uses BaseIO.Close

	var reader io.Reader = &File{}
	reader.Read(nil) // Uses BaseIO.Read

	// Case 24: Using Stringer through fmt
	dataItem := DataItem{}
	fmt.Println(dataItem) // Implicitly calls String()

	// Case 25: Chain of interfaces
	var t Third = &ChainImpl{}
	t.FirstMethod()
	t.ThirdMethod()

	// Case 26: Anonymous fields
	robot := Robot{Greeter: &GreeterImpl{}}
	robot.Greet() // Uses Greeter.Greet

	// Case 27: Interface assignment
	var src Source = &DataSource{}
	var dst Destination = src // Assignment that makes method used for both
	dst.ReadSource()
}
