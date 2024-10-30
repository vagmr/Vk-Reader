package books

type Book struct {
	ID       string `json:"id"`
	Title    string `json:"title"`
	Author   string `json:"author"`
	Path     string `json:"path"`
	Format   string `json:"format"`
	Progress int    `json:"progress"`
}

type BookManager struct {
	books map[string]*Book
}

func NewBookManager() *BookManager {
	return &BookManager{
		books: make(map[string]*Book),
	}
}

func (bm *BookManager) AddBook(book *Book) error {
	bm.books[book.ID] = book
	return nil
}

func (bm *BookManager) GetBook(id string) (*Book, bool) {
	book, exists := bm.books[id]
	return book, exists
} 