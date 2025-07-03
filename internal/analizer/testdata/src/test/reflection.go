package test

import (
	"reflect"
)

// ===============================
// REFLECTION EXAMPLES
// ===============================

// Case 28: Methods called through reflection
type ReflectableInterface interface {
	PublicMethod() string                    // want "method \"PublicMethod\" of interface \"ReflectableInterface\" is declared but not used"
	AnotherMethod(arg string) error          // want "method \"AnotherMethod\" of interface \"ReflectableInterface\" is declared but not used"
	UnusedMethod() int                       // want "method \"UnusedMethod\" of interface \"ReflectableInterface\" is declared but not used"
	ReflectOnlyMethod(data interface{}) bool // want "method \"ReflectOnlyMethod\" of interface \"ReflectableInterface\" is declared but not used"
}

// Case 29: Interface for type assertions through reflection
type TypeCheckInterface interface {
	GetType() reflect.Type // used by direct call (not through reflection)
	GetValue() interface{} // used through type assertion
	UnusedGetter() string  // want "method \"UnusedGetter\" of interface \"TypeCheckInterface\" is declared but not used"
}

// Case 30: Interface with methods checked through reflection
type IntrospectableInterface interface {
	HasMethod(name string) bool  // used by direct call
	CallMethod(name string) bool // want "method \"CallMethod\" of interface \"IntrospectableInterface\" is declared but not used"
	GetMethods() []string        // want "method \"GetMethods\" of interface \"IntrospectableInterface\" is declared but not used"
}

// Case 31: Generic interface with reflection
type GenericReflectable[T any] interface {
	ReflectType() reflect.Type // want "method \"ReflectType\" of interface \"GenericReflectable\" is declared but not used"
	GetDefault() T             // want "method \"GetDefault\" of interface \"GenericReflectable\" is declared but not used"
	ProcessReflected(v T) bool // want "method \"ProcessReflected\" of interface \"GenericReflectable\" is declared but not used"
}

// ===============================
// IMPLEMENTATIONS
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
// USAGE THROUGH REFLECTION
// ===============================

func UseReflection() {
	// Case 28: Direct calls through reflection
	impl := &ReflectableImpl{}

	// Use PublicMethod through reflection
	v := reflect.ValueOf(impl)
	method := v.MethodByName("PublicMethod")
	if method.IsValid() {
		results := method.Call(nil)
		_ = results[0].String()
	}

	// Use AnotherMethod through reflection
	anotherMethod := v.MethodByName("AnotherMethod")
	if anotherMethod.IsValid() {
		args := []reflect.Value{reflect.ValueOf("test")}
		anotherMethod.Call(args)
	}

	// Use ReflectOnlyMethod only through reflection
	reflectOnly := v.MethodByName("ReflectOnlyMethod")
	if reflectOnly.IsValid() {
		args := []reflect.Value{reflect.ValueOf("data")}
		reflectOnly.Call(args)
	}

	// Case 29: Type checking through reflection
	typeChecker := &TypeCheckImpl{}

	// Check type
	actualType := reflect.TypeOf(typeChecker)
	if actualType.Implements(reflect.TypeOf((*TypeCheckInterface)(nil)).Elem()) {
		// Use GetType
		resultType := typeChecker.GetType()
		_ = resultType.Name()

		// Use GetValue through type assertion
		value := typeChecker.GetValue()
		if typed, ok := value.(*TypeCheckImpl); ok {
			_ = typed
		}
	}

	// Case 30: Method introspection
	introspectable := &IntrospectableImpl{}

	// Use HasMethod through reflection
	hasMethod := introspectable.HasMethod("PublicMethod")
	_ = hasMethod

	// Case 31: Generic with reflection
	genericImpl := &GenericReflectableImpl[string]{value: "test"}

	// Use ReflectType through reflection
	gv := reflect.ValueOf(genericImpl)
	reflectTypeMethod := gv.MethodByName("ReflectType")
	if reflectTypeMethod.IsValid() {
		results := reflectTypeMethod.Call(nil)
		reflectedType := results[0].Interface().(reflect.Type)
		_ = reflectedType.Kind()
	}

	// Use GetDefault in a regular way
	defaultValue := genericImpl.GetDefault()
	_ = defaultValue
}

// Additional function for demonstrating complex reflection
func ComplexReflectionUsage() {
	// Create a slice of interfaces
	interfaces := []interface{}{
		&ReflectableImpl{},
		&TypeCheckImpl{},
		&IntrospectableImpl{},
	}

	// Iterate and call methods through reflection
	for _, iface := range interfaces {
		v := reflect.ValueOf(iface)
		t := reflect.TypeOf(iface)

		// Check all methods
		for i := 0; i < t.NumMethod(); i++ {
			method := t.Method(i)
			methodValue := v.Method(i)

			// Call parameterless methods through reflection
			if method.Type.NumIn() == 1 { // receiver only
				if method.Type.NumOut() > 0 {
					results := methodValue.Call(nil)
					_ = results
				}
			}
		}
	}
}
