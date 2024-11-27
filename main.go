package main

import (
	"encoding/json"
	"net/http"
	"regexp"

	"github.com/demo/packer/recipes"
	"github.com/gosimple/slug"
)

// bad func
// func createPointer() *int {
//     var x int
//     return &x // returning a pointer to a local variable
// }

type homeHandler struct {}

func (h *homeHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("this is my home page"))
}

func InternalServerErrorHandler(w http.ResponseWriter, r *http.Request) {
    w.WriteHeader(http.StatusInternalServerError)
    w.Write([]byte("500 Internal Server Error"))
}

func NotFoundHandler(w http.ResponseWriter, r *http.Request) {
    w.WriteHeader(http.StatusNotFound)
    w.Write([]byte("404 Not Found"))
}

type RecipesHandler struct{
    store recipeStore
}

type recipeStore interface {
    Add(name string, recipe recipes.Recipe) error
    Get(name string) (recipes.Recipe, error)
    Update(name string, recipe recipes.Recipe) error
    List() (map[string]recipes.Recipe, error)
    Remove(name string) error
}

func NewRecipesHandler(s recipeStore) *RecipesHandler {
    return &RecipesHandler{
        store: s,
    }
}

func (h *RecipesHandler) CreateRecipe(w http.ResponseWriter, r *http.Request) {
    var recipe recipes.Recipe
    if err := json.NewDecoder(r.Body).Decode(&recipe); err != nil {
        InternalServerErrorHandler(w, r)
        return
    }

    resourceID := slug.Make(recipe.Name)
    if err := h.store.Add(resourceID, recipe); err != nil {
        InternalServerErrorHandler(w, r)
        return
    }

    // set status code to 200
    w.WriteHeader(http.StatusOK)
}
func (h *RecipesHandler) ListRecipes(w http.ResponseWriter, r *http.Request) {
    recipes, err := h.store.List()
    if err != nil {
        InternalServerErrorHandler(w, r)
        return
    }
    jsonBytes, err := json.Marshal(recipes)
    if err != nil {
        NotFoundHandler(w, r)
        return
    }
    w.WriteHeader(http.StatusOK)
    w.Write(jsonBytes)
}

func (h *RecipesHandler) GetRecipe(w http.ResponseWriter, r *http.Request) {
    matches := RecipeReWithID.FindStringSubmatch(r.URL.Path)
    recipe, err := h.store.Get(matches[1])
    if err != nil {
        if err == recipes.ErrNotFound {
            NotFoundHandler(w, r)
            return
        }
        InternalServerErrorHandler(w, r)
        return
    }
    jsonBytes, err := json.Marshal(recipe)
    if err != nil {
        InternalServerErrorHandler(w, r)
        return
    }
    w.WriteHeader(http.StatusOK)
    w.Write(jsonBytes)
}
func (h *RecipesHandler) UpdateRecipe(w http.ResponseWriter, r *http.Request) {
    matches := RecipeReWithID.FindStringSubmatch(r.URL.Path)
    if len(matches) < 2 {
        InternalServerErrorHandler(w, r)
        return
    }
    var recipe recipes.Recipe
    if err := json.NewDecoder(r.Body).Decode(&recipe); err != nil {
        InternalServerErrorHandler(w, r)
        return
    }
    if err := h.store.Update(matches[1], recipe); err != nil {
        if err == recipes.ErrNotFound {
            NotFoundHandler(w, r)
            return
        }

        InternalServerErrorHandler(w, r)
        return
    }
    
    w.WriteHeader(http.StatusOK)
}
func (h *RecipesHandler) DeleteRecipe(w http.ResponseWriter, r *http.Request) {
    matches := RecipeReWithID.FindStringSubmatch(r.URL.Path)
    if len(matches) < 2 {
        InternalServerErrorHandler(w, r)
        return
    }

    if err := h.store.Remove(matches[1]); err != nil {
        InternalServerErrorHandler(w, r)
        return
    }

    w.WriteHeader(http.StatusOK)  
}

var (
    RecipeRe       = regexp.MustCompile(`^/recipes/*$`)
    RecipeReWithID = regexp.MustCompile(`^/recipes/([a-z0-9]+(?:-[a-z0-9]+)+)$`)
)
func (h *RecipesHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
    switch {
    case r.Method == http.MethodPost && RecipeRe.MatchString(r.URL.Path):
        h.CreateRecipe(w, r)
        return
    case r.Method == http.MethodGet && RecipeRe.MatchString(r.URL.Path):
        h.ListRecipes(w, r)
        return
    case r.Method == http.MethodGet && RecipeReWithID.MatchString(r.URL.Path):
        h.GetRecipe(w, r)
        return
    case r.Method == http.MethodPut && RecipeReWithID.MatchString(r.URL.Path):
        h.UpdateRecipe(w, r)
        return
    case r.Method == http.MethodDelete && RecipeReWithID.MatchString(r.URL.Path):
        h.DeleteRecipe(w, r)
        return
    default:
        return
    }
}

func main() {
	// //print 'ping' until click any key
	// var input string
	// fmt.Scanln(&input)
	// concurrency.CheckTimeoutWithSelect()
    store := recipes.NewMemStore()
    recipeController := NewRecipesHandler(store)
	mux := http.NewServeMux()

	// create custom handler which implement ServeHTTP
	mux.Handle("/", &homeHandler{})
    mux.Handle("/recipes", LogTimeRequestMiddleware(LogRequestMiddleware(recipeController)))
    mux.Handle("/recipes/", LogTimeRequestMiddleware(LogRequestMiddleware(recipeController)))
	// start an HTTP server listen on port 8080
	http.ListenAndServe(":8080", mux)
}
