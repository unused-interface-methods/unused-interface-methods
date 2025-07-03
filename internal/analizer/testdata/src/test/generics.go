package test

// ===============================
// GENERIC INTERFACES
// ===============================

// 1. Simple generic
type SimpleRepo[T any] interface {
	Get(id string) T   // want "method \"Get\" of interface \"SimpleRepo\" is declared but not used"
	Save(item T) error // want "method \"Save\" of interface \"SimpleRepo\" is declared but not used"
}

// 2. Generic with constraints
type Comparable interface {
	Compare(other Comparable) int // want "method \"Compare\" of interface \"Comparable\" is declared but not used"
}

type SortableRepo[T Comparable] interface {
	GetSorted() []T      // want "method \"GetSorted\" of interface \"SortableRepo\" is declared but not used"
	Insert(item T) error // want "method \"Insert\" of interface \"SortableRepo\" is declared but not used"
	Remove(item T) bool  // want "method \"Remove\" of interface \"SortableRepo\" is declared but not used"
}

// 3. Multiple type parameters
type Cache[K comparable, V any] interface {
	Get(key K) (V, bool) // want "method \"Get\" of interface \"Cache\" is declared but not used"
	Set(key K, value V)  // want "method \"Set\" of interface \"Cache\" is declared but not used"
	Delete(key K) bool   // want "method \"Delete\" of interface \"Cache\" is declared but not used"
	Keys() []K           // want "method \"Keys\" of interface \"Cache\" is declared but not used"
	Values() []V         // want "method \"Values\" of interface \"Cache\" is declared but not used"
}

// 4. Complex constraints
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

// 5. Nested generics
type NestedRepo[T any] interface {
	GetMap() map[string]T                         // want "method \"GetMap\" of interface \"NestedRepo\" is declared but not used"
	GetSlice() []T                                // want "method \"GetSlice\" of interface \"NestedRepo\" is declared but not used"
	GetChannel() chan T                           // want "method \"GetChannel\" of interface \"NestedRepo\" is declared but not used"
	ProcessBatch(items []T) (map[string]T, error) // want "method \"ProcessBatch\" of interface \"NestedRepo\" is declared but not used"
}

// 6. Generic interface
type GenericRepository[T any] interface {
	Get(id string) (T, error) // used (in GetUser and ListPosts)
	Save(item T) error        // used (in SaveUser)
	Delete(id string) error   // want "method \"Delete\" of interface \"GenericRepository\" is declared but not used"
	List() ([]T, error)       // used (in ListPosts)
}

// 7. Simple generic for testing
type Repository[T any] interface {
	Get(id string) (T, error) // used (in GetUser)
	Save(item T) error        // want "method \"Save\" of interface \"Repository\" is declared but not used"
	Delete(id string) error   // want "method \"Delete\" of interface \"Repository\" is declared but not used"
	List() ([]T, error)       // used (in ListPosts)
}

// ===============================
// REGULAR INTERFACE FOR COMPARISON
// ===============================

// Regular interface for comparison with generics
type RegularInterface interface {
	DoSomething() error // used (in Work)
	GetResult() string  // want "method \"GetResult\" of interface \"RegularInterface\" is declared but not used"
}

// ===============================
// USAGE
// ===============================

// Concrete types
type User struct {
	ID   string
	Name string
}

type Post struct {
	ID    string
	Title string
}

// Using generics
type UserService struct {
	userRepo GenericRepository[User]
	repo     Repository[User]
}

type PostService struct {
	postRepo GenericRepository[Post]
	repo     Repository[Post]
}

// Using regular interface
type Service struct {
	regular RegularInterface
}

func (s *Service) Work() {
	s.regular.DoSomething() // This method is used
	// GetResult() is NOT used
}

// Using generic methods
func (us *UserService) GetUser(id string) (*User, error) {
	// Call on GenericRepository[User]
	user, err := us.userRepo.Get(id)
	if err != nil {
		return nil, err
	}

	// Call on Repository[User]
	us.repo.Get(id)

	return &user, nil
}

func (us *UserService) SaveUser(user User) error {
	// Call Save on GenericRepository[User]
	return us.userRepo.Save(user)
}

func (ps *PostService) ListPosts() ([]Post, error) {
	// Call List on GenericRepository[Post]
	posts, err := ps.postRepo.List()
	if err != nil {
		return nil, err
	}

	// Call on Repository[Post]
	ps.repo.List()

	return posts, nil
}

// Delete is NOT used in any instantiation
// Save in Repository[T] is NOT used
