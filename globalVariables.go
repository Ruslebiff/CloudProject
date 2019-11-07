package cravings

import (
	"time"
)

var StartTime time.Time // start run time

const FirestoreCredentials = "./cloudproject-2a9c2-firebase-adminsdk-0om9b-bca5ed564a.json"
const RecipeCollection = "recipes"
const IngredientCollection = "ingredients"
const TokenCollection = "tokens"
const WebhooksCollection = "webhhooks"

//********************* URL ********************
var URLRegistration = "https://api.edamam.com/api/nutrition-details"
