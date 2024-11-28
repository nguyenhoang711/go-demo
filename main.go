package main

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"regexp"
	"strconv"

	"github.com/demo/packer/middlewares"
	"github.com/demo/packer/recipes"
	_ "github.com/go-sql-driver/mysql"
	"github.com/gosimple/slug"
)

// bad func
// func createPointer() *int {
//     var x int
//     return &x // returning a pointer to a local variable
// }

type homeHandler struct{}

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

type RecipesHandler struct {
	store recipeStore
}

var db *sql.DB

type recipHandlerV2 struct {}

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

func NewRecipesHandlerV2() *recipHandlerV2 {
	return &recipHandlerV2{}
}

func (h *recipHandlerV2) addRecipe(w http.ResponseWriter, r *http.Request) {
	var recipe recipes.Recipe
	if err := json.NewDecoder(r.Body).Decode(&recipe); err != nil {
		InternalServerErrorHandler(w, r)
		return
	}

	if err := recipes.CreateRecipe(db, recipe.Ingredients, recipe.Name); err != nil {
		InternalServerErrorHandler(w, r)
		return
	}
		
	// set status code to 200
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("create recipe success"))
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

func (h *recipHandlerV2) gettRecipe(w http.ResponseWriter, r *http.Request) {
	matches := RecipeReWithIDV2.FindStringSubmatch(r.URL.Path)
	i, err := strconv.Atoi(matches[1])
	if err != nil {
		InternalServerErrorHandler(w, r)
		return
	}
	recipe, err := recipes.GetRecipe(db, i)
	if err != nil {
		InternalServerErrorHandler(w, r)
		return
	}
	jsonBytes, err := json.Marshal(recipe)
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

func (h *recipHandlerV2) updateRecipe(w http.ResponseWriter, r *http.Request) {
	var recipe recipes.Recipe
	if err := json.NewDecoder(r.Body).Decode(&recipe); err != nil {
		InternalServerErrorHandler(w, r)
		return
	}

	err := recipes.UpdateRecipe(db, recipe)
	if err != nil {
		InternalServerErrorHandler(w, r)
		return
	}
	jsonBytes, err := json.Marshal(recipe)
	if err != nil {
		NotFoundHandler(w, r)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write(jsonBytes)
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
	RecipeRe       = regexp.MustCompile(`^/recipes-v2/*$`)
	RecipeReWithID = regexp.MustCompile(`^/recipes-v2/([a-z0-9]+(?:-[a-z0-9]+)+)$`)
	RecipeReWithIDV2 = regexp.MustCompile(`^/recipes-v2/([a-z0-9]+(?:-[a-z0-9]+)*)$`)
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

func (h *recipHandlerV2) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch {
	case r.Method == http.MethodPost && RecipeRe.MatchString(r.URL.Path):
		h.addRecipe(w, r)
		return
	case r.Method == http.MethodGet && RecipeReWithIDV2.MatchString(r.URL.Path):
		h.gettRecipe(w, r)
		return
	case r.Method == http.MethodPut && RecipeRe.MatchString(r.URL.Path):
		h.updateRecipe(w, r)
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

	// get a database handle.
	var err error
	db, err = sql.Open("mysql", "root:123456@tcp(localhost:3306)/recipes")
	if err != nil {
		log.Panic(err)
	}

	store := recipes.NewMemStore()
	recipeController := NewRecipesHandler(store)
	recipeControllerV2 := NewRecipesHandlerV2()
	mux := http.NewServeMux()

	// create custom handler which implement ServeHTTP
	mux.Handle("/", &homeHandler{})
	mux.Handle("/recipes", middlewares.LogTimeRequestMiddleware(middlewares.LogRequestMiddleware(recipeController)))
	mux.Handle("/recipes/", middlewares.LogTimeRequestMiddleware(middlewares.LogRequestMiddleware(recipeController)))
	mux.Handle("/recipes-v2", middlewares.LogTimeRequestMiddleware(middlewares.LogRequestMiddleware(recipeControllerV2)))
	mux.Handle("/recipes-v2/", middlewares.LogTimeRequestMiddleware(middlewares.LogRequestMiddleware(recipeControllerV2)))
	// start an HTTP server listen on port 8080
	http.ListenAndServe(":8080", mux)
}
