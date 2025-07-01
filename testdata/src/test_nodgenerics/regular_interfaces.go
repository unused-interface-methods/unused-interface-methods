package test_data

// ===============================
// ОБЫЧНЫЕ ИНТЕРФЕЙСЫ ИЗ GENERICS.GO
// ===============================

// 2. Дженерик с ограничениями
type Comparable interface {
	Compare(other Comparable) int // want "method \"Compare\" of interface \"Comparable\" is declared but not used"
}

// Обычный интерфейс для сравнения с дженериками
type RegularInterface interface {
	DoSomething() error // используется (в Work)
	GetResult() string  // want "method \"GetResult\" of interface \"RegularInterface\" is declared but not used"
}

// ===============================
// ИСПОЛЬЗОВАНИЕ
// ===============================

// Использование обычного интерфейса
type Service struct {
	regular RegularInterface
}

func (s *Service) Work() {
	s.regular.DoSomething() // Этот метод используется
	// GetResult() НЕ используется
}
