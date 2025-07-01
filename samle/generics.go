package test_data

// ===============================
// ДЖЕНЕРИК-ИНТЕРФЕЙСЫ
// ===============================

// 1. Простой дженерик
type SimpleRepo[T any] interface {
	Get(id string) T   // не используется
	Save(item T) error // не используется
}

// 2. Дженерик с ограничениями
type Comparable interface {
	Compare(other Comparable) int // не используется
}

type SortableRepo[T Comparable] interface {
	GetSorted() []T      // не используется
	Insert(item T) error // не используется
	Remove(item T) bool  // не используется
}

// 3. Множественные параметры типа
type Cache[K comparable, V any] interface {
	Get(key K) (V, bool) // не используется
	Set(key K, value V)  // не используется
	Delete(key K) bool   // не используется
	Keys() []K           // не используется
	Values() []V         // не используется
}

// 4. Сложные ограничения
type Serializable interface {
	Serialize() []byte        // не используется
	Deserialize([]byte) error // не используется
}

type PersistentCache[K comparable, V Serializable] interface {
	Load(key K) (V, error)      // не используется
	Store(key K, value V) error // не используется
	Persist() error             // не используется
	Restore() error             // не используется
}

// 5. Вложенные дженерики
type NestedRepo[T any] interface {
	GetMap() map[string]T                         // не используется
	GetSlice() []T                                // не используется
	GetChannel() chan T                           // не используется
	ProcessBatch(items []T) (map[string]T, error) // не используется
}

// 6. Дженерик-интерфейс из doc/GENERICS_PROBLEM.md
type GenericRepository[T any] interface {
	Get(id string) (T, error) // используется (в GetUser и ListPosts)
	Save(item T) error        // используется (в SaveUser)
	Delete(id string) error   // не используется
	List() ([]T, error)       // используется (в ListPosts)
}

// 7. Простой дженерик для тестирования
type Repository[T any] interface {
	Get(id string) (T, error) // используется (в GetUser)
	Save(item T) error        // не используется
	Delete(id string) error   // не используется
	List() ([]T, error)       // используется (в ListPosts)
}

// ===============================
// ОБЫЧНЫЙ ИНТЕРФЕЙС ДЛЯ СРАВНЕНИЯ
// ===============================

// Обычный интерфейс для сравнения с дженериками
type RegularInterface interface {
	DoSomething() error // используется (в Work)
	GetResult() string  // не используется
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
