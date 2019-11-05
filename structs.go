package cravings

//  Struct for a recipe containting ingredients
type Recpie struct {
	Ingredients []Ingredient `json:"ingredients"`
}

//  Struct for an ingredient
type Ingredient struct {
	Quantity int    `json:"quantity"`
	Unit     string `json:"unit"`
	Name     string `json:"name"`
	Calories int    `json:"kcal"`
	Weight   int    `json:"weight"`
}
