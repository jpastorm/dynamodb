package repository

type PostRepository interface {
	Save(post *Post) (*Post, error)
	FindAll() ([]Post, error)
	FindByID(id string) (*Post, error)
	Delete(post *Post) error
}
