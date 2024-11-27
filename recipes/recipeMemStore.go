package recipes

import "errors"

var (
	ErrNotFound = errors.New("not found")
)
type MemStore struct {
	recipes map[string]Recipe
}

func NewMemStore() *MemStore {
	list := make(map[string]Recipe)
	return &MemStore{
		list,
	}
}

func (r MemStore) Add(name string, recipe Recipe) error {
	r.recipes[name] = recipe
	return nil
}

func (r MemStore) Get(name string) (Recipe, error) {
	if val, ok := r.recipes[name]; ok {
		return val, nil
	}
	return Recipe{}, ErrNotFound
}

func(r MemStore) List() (map[string]Recipe, error) {
	return r.recipes, nil
}

func(r MemStore) Update(name string, recipe Recipe) error {
	if _, ok := r.recipes[name]; ok {
		r.recipes[name] = recipe
		return nil
	}
	return ErrNotFound
}

func (r MemStore) Remove(name string) error {
	delete(r.recipes, name)
	return nil
}