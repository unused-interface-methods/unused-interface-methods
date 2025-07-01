package test_data

// ===============================
// ДЖЕНЕРИК-ИНТЕРФЕЙСЫ
// ===============================

// 1. Простой дженерик
type SimpleRepo[T any] interface {
	Get(id string) T   // want "method \"Get\" of interface \"SimpleRepo\" is declared but not used"
	Save(item T) error // want "method \"Save\" of interface \"SimpleRepo\" is declared but not used"
}

// 2. Дженерик с ограничениями
type Comparable interface {
	Compare(other Comparable) int // want "method \"Compare\" of interface \"Comparable\" is declared but not used"
}

type SortableRepo[T Comparable] interface {
	GetSorted() []T      // want "method \"GetSorted\" of interface \"SortableRepo\" is declared but not used"
	Insert(item T) error // want "method \"Insert\" of interface \"SortableRepo\" is declared but not used"
	Remove(item T) bool  // want "method \"Remove\" of interface \"SortableRepo\" is declared but not used"
}

// 3. Множественные параметры типа
type Cache[K comparable, V any] interface {
	Get(key K) (V, bool) // want "method \"Get\" of interface \"Cache\" is declared but not used"
	Set(key K, value V)  // want "method \"Set\" of interface \"Cache\" is declared but not used"
	Delete(key K) bool   // want "method \"Delete\" of interface \"Cache\" is declared but not used"
	Keys() []K           // want "method \"Keys\" of interface \"Cache\" is declared but not used"
	Values() []V         // want "method \"Values\" of interface \"Cache\" is declared but not used"
}

// 4. Сложные ограничения
type Serializable interface {
	Serialize() []byte        // want "method \"Serialize\" of interface \"Serializable\" is declared but not used"
	Deserialize([]byte) error // want "method \"Deserialize\" of interface \"Serializable\" is declared but not used"
}

type PersistentCache[K comparable, V Serializable] interface {
	Load(key K) (V, error)      // want "method \"Load\" of interface \"PersistentCache\" is declared but not used"
	Store(key K, value V) error // want "method \"Store\" of interface \"PersistentCache\" is declared but not used"
	Persist() error             // want "method \"Persist\" of interface \"PersistentCache\" is declared but not used"
	Restore() error             // want "method \"Restore\" of interface \"PersistentCache\" is declared but not used"
}

// 5. Вложенные дженерики
type NestedRepo[T any] interface {
	GetMap() map[string]T                         // want "method \"GetMap\" of interface \"NestedRepo\" is declared but not used"
	GetSlice() []T                                // want "method \"GetSlice\" of interface \"NestedRepo\" is declared but not used"
	GetChannel() chan T                           // want "method \"GetChannel\" of interface \"NestedRepo\" is declared but not used"
	ProcessBatch(items []T) (map[string]T, error) // want "method \"ProcessBatch\" of interface \"NestedRepo\" is declared but not used"
}

// 6. Дженерик-интерфейс из doc/GENERICS_PROBLEM.md
type GenericRepository[T any] interface {
	Get(id string) (T, error) // используется (в GetUser и ListPosts)
	Save(item T) error        // используется (в SaveUser)
	Delete(id string) error   // want "method \"Delete\" of interface \"GenericRepository\" is declared but not used"
	List() ([]T, error)       // используется (в ListPosts)
}

// 7. Простой дженерик для тестирования
type Repository[T any] interface {
	Get(id string) (T, error) // используется (в GetUser)
	Save(item T) error        // want "method \"Save\" of interface \"Repository\" is declared but not used"
	Delete(id string) error   // want "method \"Delete\" of interface \"Repository\" is declared but not used"
	List() ([]T, error)       // используется (в ListPosts)
}

// ===============================
// ОБЫЧНЫЙ ИНТЕРФЕЙС ДЛЯ СРАВНЕНИЯ
// ===============================

// Обычный интерфейс для сравнения с дженериками
type RegularInterface interface {
	DoSomething() error // используется (в Work)
	GetResult() string  // want "method \"GetResult\" of interface \"RegularInterface\" is declared but not used"
}

// ===============================
// ИСПОЛЬЗОВАНИЕ
// ===============================

// Конкретные типы
type User struct {
	ID   string
	Name string
}

type Post struct {
	ID    string
	Title string
}

// Использование дженериков
type UserService struct {
	userRepo GenericRepository[User]
	repo     Repository[User]
}

type PostService struct {
	postRepo GenericRepository[Post]
	repo     Repository[Post]
}

// Использование обычного интерфейса
type Service struct {
	regular RegularInterface
}

func (s *Service) Work() {
	s.regular.DoSomething() // Этот метод используется
	// GetResult() НЕ используется
}

// Использование дженерик-методов
func (us *UserService) GetUser(id string) (*User, error) {
	// Вызов на GenericRepository[User]
	user, err := us.userRepo.Get(id)
	if err != nil {
		return nil, err
	}

	// Вызов на Repository[User]
	us.repo.Get(id)

	return &user, nil
}

func (us *UserService) SaveUser(user User) error {
	// Вызов Save на GenericRepository[User]
	return us.userRepo.Save(user)
}

func (ps *PostService) ListPosts() ([]Post, error) {
	// Вызов List на GenericRepository[Post]
	posts, err := ps.postRepo.List()
	if err != nil {
		return nil, err
	}

	// Вызов на Repository[Post]
	ps.repo.List()

	return posts, nil
}

// Delete НЕ используется ни в одном инстанцировании
// Save в Repository[T] НЕ используется
