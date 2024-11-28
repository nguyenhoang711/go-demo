package recipes

// Represents a recipe
type Recipe struct {
	Id          int
	Name        string       `json:"name"`
	Ingredients []Ingredient `json:"ingredients"`
}

// Represents individual ingredients
type Ingredient struct {
	Id       int
	RecipeID int   	`json:"recipe_id"`
	Name     string `json:"name"`
	Amount   int    `json:"amount"`
}
