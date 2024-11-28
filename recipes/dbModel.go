package recipes

import (
	"database/sql"
	"errors"
	"fmt"
)

func GetRecipe(db *sql.DB, recipeId int) (Recipe, error) {
	var recipe Recipe
	err := db.QueryRow("SELECT * FROM recipe WHERE id = ?", recipeId).Scan(&recipe.Id, &recipe.Name)
	if err != nil {
		return recipe, fmt.Errorf("error in query all recipes: %v", err)
	}

	// Get the ingredients
	rows2, err := db.Query("SELECT id, recipe_id, name, amount FROM ingredient WHERE recipe_id = ?", recipeId)
	if err != nil {
		return recipe, err
	}
	defer rows2.Close()
	for rows2.Next() {
		var ingredient Ingredient
		err := rows2.Scan(&ingredient.Id, &ingredient.RecipeID, &ingredient.Name, &ingredient.Amount)
		if err != nil {
			return recipe, err
		}
		recipe.Ingredients = append(recipe.Ingredients, ingredient)
	}
	return recipe, nil
}

func CreateRecipe(db *sql.DB, ingredients []Ingredient, recipeName string) error {
	tx, err := db.Begin()
	if err != nil {
		return err
	}
	if recipeName == "" {
		err := errors.New("recipe name not exists")
		return err
	}
	result, err := db.Exec("INSERT INTO recipe (name) VALUES (?)", recipeName)
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("add recipe error:: %v", err)
	}
	id, err := result.LastInsertId()
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("add recipe error:: %v", err)
	}
	for _, ingre := range ingredients {
		_, err := db.Exec("INSERT INTO ingredient (name, amount, recipe_id) VALUES (?, ?, ?)", ingre.Name, ingre.Amount, id)
		if err != nil {
			tx.Rollback()
			return fmt.Errorf("add ingredient error:: %v", err)
		}
	}
	return tx.Commit()
}
