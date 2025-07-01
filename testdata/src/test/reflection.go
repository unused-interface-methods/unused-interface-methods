package test_data

import (
	"reflect"
)

// ===============================
// ПРИМЕРЫ С РЕФЛЕКСИЕЙ
// ===============================

// Кейс 28: Методы, вызываемые через рефлексию
type ReflectableInterface interface {
	PublicMethod() string                    // want "method \"PublicMethod\" of interface \"ReflectableInterface\" is declared but not used"
	AnotherMethod(arg string) error          // want "method \"AnotherMethod\" of interface \"ReflectableInterface\" is declared but not used"
	UnusedMethod() int                       // want "method \"UnusedMethod\" of interface \"ReflectableInterface\" is declared but not used"
	ReflectOnlyMethod(data interface{}) bool // want "method \"ReflectOnlyMethod\" of interface \"ReflectableInterface\" is declared but not used"
}

// Кейс 29: Интерфейс для type assertions через рефлексию
type TypeCheckInterface interface {
	GetType() reflect.Type // используется прямым вызовом (не через рефлексию)
	GetValue() interface{} // используется через type assertion
	UnusedGetter() string  // want "method \"UnusedGetter\" of interface \"TypeCheckInterface\" is declared but not used"
}

// Кейс 30: Интерфейс с методами, проверяемыми через рефлексию
type IntrospectableInterface interface {
	HasMethod(name string) bool  // используется прямым вызовом
	CallMethod(name string) bool // want "method \"CallMethod\" of interface \"IntrospectableInterface\" is declared but not used"
	GetMethods() []string        // want "method \"GetMethods\" of interface \"IntrospectableInterface\" is declared but not used"
}

// Кейс 31: Дженерик интерфейс с рефлексией
type GenericReflectable[T any] interface {
	ReflectType() reflect.Type // want "method \"ReflectType\" of interface \"GenericReflectable\" is declared but not used"
	GetDefault() T             // want "method \"GetDefault\" of interface \"GenericReflectable\" is declared but not used"
	ProcessReflected(v T) bool // want "method \"ProcessReflected\" of interface \"GenericReflectable\" is declared but not used"
}

// ===============================
// РЕАЛИЗАЦИИ
// ===============================

type ReflectableImpl struct{}

func (r *ReflectableImpl) PublicMethod() string                    { return "public" }
func (r *ReflectableImpl) AnotherMethod(arg string) error          { return nil }
func (r *ReflectableImpl) UnusedMethod() int                       { return 42 }
func (r *ReflectableImpl) ReflectOnlyMethod(data interface{}) bool { return true }

type TypeCheckImpl struct{}

func (t *TypeCheckImpl) GetType() reflect.Type { return reflect.TypeOf(t) }
func (t *TypeCheckImpl) GetValue() interface{} { return t }
func (t *TypeCheckImpl) UnusedGetter() string  { return "unused" }

type IntrospectableImpl struct{}

func (i *IntrospectableImpl) HasMethod(name string) bool {
	t := reflect.TypeOf(i)
	_, ok := t.MethodByName(name)
	return ok
}
func (i *IntrospectableImpl) CallMethod(name string) bool { return false }
func (i *IntrospectableImpl) GetMethods() []string        { return []string{} }

type GenericReflectableImpl[T any] struct {
	value T
}

func (g *GenericReflectableImpl[T]) ReflectType() reflect.Type { return reflect.TypeOf(g.value) }
func (g *GenericReflectableImpl[T]) GetDefault() T             { var zero T; return zero }
func (g *GenericReflectableImpl[T]) ProcessReflected(v T) bool { return true }

// ===============================
// ИСПОЛЬЗОВАНИЕ ЧЕРЕЗ РЕФЛЕКСИЮ
// ===============================

func UseReflection() {
	// Кейс 28: Прямые вызовы через рефлексию
	impl := &ReflectableImpl{}

	// Используем PublicMethod через рефлексию
	v := reflect.ValueOf(impl)
	method := v.MethodByName("PublicMethod")
	if method.IsValid() {
		results := method.Call(nil)
		_ = results[0].String()
	}

	// Используем AnotherMethod через рефлексию
	anotherMethod := v.MethodByName("AnotherMethod")
	if anotherMethod.IsValid() {
		args := []reflect.Value{reflect.ValueOf("test")}
		anotherMethod.Call(args)
	}

	// Используем ReflectOnlyMethod только через рефлексию
	reflectOnly := v.MethodByName("ReflectOnlyMethod")
	if reflectOnly.IsValid() {
		args := []reflect.Value{reflect.ValueOf("data")}
		reflectOnly.Call(args)
	}

	// Кейс 29: Type checking через рефлексию
	typeChecker := &TypeCheckImpl{}

	// Проверяем тип
	actualType := reflect.TypeOf(typeChecker)
	if actualType.Implements(reflect.TypeOf((*TypeCheckInterface)(nil)).Elem()) {
		// Используем GetType
		resultType := typeChecker.GetType()
		_ = resultType.Name()

		// Используем GetValue через type assertion
		value := typeChecker.GetValue()
		if typed, ok := value.(*TypeCheckImpl); ok {
			_ = typed
		}
	}

	// Кейс 30: Интроспекция методов
	introspectable := &IntrospectableImpl{}

	// Используем HasMethod через рефлексию
	hasMethod := introspectable.HasMethod("PublicMethod")
	_ = hasMethod

	// Кейс 31: Дженерик с рефлексией
	genericImpl := &GenericReflectableImpl[string]{value: "test"}

	// Используем ReflectType через рефлексию
	gv := reflect.ValueOf(genericImpl)
	reflectTypeMethod := gv.MethodByName("ReflectType")
	if reflectTypeMethod.IsValid() {
		results := reflectTypeMethod.Call(nil)
		reflectedType := results[0].Interface().(reflect.Type)
		_ = reflectedType.Kind()
	}

	// Используем GetDefault обычным способом
	defaultValue := genericImpl.GetDefault()
	_ = defaultValue
}

// Дополнительная функция для демонстрации сложной рефлексии
func ComplexReflectionUsage() {
	// Создаем слайс интерфейсов
	interfaces := []interface{}{
		&ReflectableImpl{},
		&TypeCheckImpl{},
		&IntrospectableImpl{},
	}

	// Итерируемся и вызываем методы через рефлексию
	for _, iface := range interfaces {
		v := reflect.ValueOf(iface)
		t := reflect.TypeOf(iface)

		// Проверяем все методы
		for i := 0; i < t.NumMethod(); i++ {
			method := t.Method(i)
			methodValue := v.Method(i)

			// Вызываем методы без параметров через рефлексию
			if method.Type.NumIn() == 1 { // только receiver
				if method.Type.NumOut() > 0 {
					results := methodValue.Call(nil)
					_ = results
				}
			}
		}
	}
}
