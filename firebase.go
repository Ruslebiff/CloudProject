package cravings

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	firebase "firebase.google.com/go"
	"github.com/pkg/errors"
	"google.golang.org/api/iterator"
	"google.golang.org/api/option"
)

// FireBaseDB is an instance of FirestoreDatabase struct which is used in firebase.go
var FireBaseDB = FirestoreDatabase{}

// DBInit initialises the database
func DBInit() error {
	// Firebase initialisation
	FireBaseDB.Ctx = context.Background()
	// We use a service account, load credentials file that you downloaded from your project's settings menu.
	// Make sure this file is gitignored, it is the access token to the database.
	sa := option.WithCredentialsFile(FirestoreCredentials)
	app, err := firebase.NewApp(FireBaseDB.Ctx, nil, sa) //  Creates the application with its contents
	if err != nil {
		fmt.Println("Failed to initialize the firebase database when creating a new app: ", err)
	}
	//  Sets the app created to our local struct's client
	FireBaseDB.Client, err = app.Firestore(FireBaseDB.Ctx)
	if err != nil {
		fmt.Println("Failed to create app")
	}
	return err
}

// DBClose Close firebase connection
func DBClose() {
	FireBaseDB.Client.Close()
}

// DBSaveRecipe saves recipe to database
func DBSaveRecipe(r *Recipe) error { //  Creates a new document in firebase
	ref := FireBaseDB.Client.Collection(RecipeCollection).NewDoc()
	r.ID = ref.ID                        //  Asserts the recipes id to be the one given by firebase
	_, err := ref.Set(FireBaseDB.Ctx, r) //  Set the context of the document to the one of the recipe
	if err != nil {
		fmt.Println("ERROR saving recipe to recipe collection: ", err)
	}
	return nil
}

// DBSaveIngredient saves ingredient to database
func DBSaveIngredient(i *Ingredient) error { //  Creates a new document in firebase
	ref := FireBaseDB.Client.Collection(IngredientCollection).NewDoc()
	i.ID = ref.ID                        //  Asserts the ingredients id to be the one given by firebase
	_, err := ref.Set(FireBaseDB.Ctx, i) //  Set the context of the document to the one of the ingredient
	if err != nil {
		fmt.Println("ERROR saving ingredient to ingredients collection: ", err)
	}
	return nil
}

// DBSaveWebhook saves a new webhook to the database
func DBSaveWebhook(i *Webhook) error {
	ref := FireBaseDB.Client.Collection(WebhooksCollection).NewDoc()
	i.ID = ref.ID                        //  Asserts the webhooks id to be the one given by firebase
	_, err := ref.Set(FireBaseDB.Ctx, i) //  Set the context of the document to the one of the webhook
	if err != nil {
		fmt.Println("ERROR saving ingredient to ingredients collection: ", err)
	}
	return nil
}

// DBDelete deletes an entry from given collection in database by its id, either ingredient, recipe or webhook
func DBDelete(id string, collection string) error {
	_, err := FireBaseDB.Client.Collection(collection).Doc(id).Delete(FireBaseDB.Ctx)
	if err != nil {
		fmt.Println("ERROR deleting from collection: "+collection, err)
		return errors.Wrap(err, "Error in FirebaseDatabase.Delete()")
	}
	return nil
}

// DBReadRecipeByName reads a single recipe by Name
func DBReadRecipeByName(name string) (Recipe, error) {
	temp := Recipe{}                  //  Recipe to be returned
	allrec, err := DBReadAllRecipes() //  Query all the recipes
	if err != nil {
		return temp, err
	}
	//  Loops through all the recipes and checks if the parameter name is equal to one of the recipes
	for _, i := range allrec {
		if i.RecipeName == name { //  If recipe is found, return
			temp.ID = i.ID
			temp.RecipeName = i.RecipeName
			temp.Ingredients = i.Ingredients

			return temp, err
		}
	}
	return temp, err
}

// DBReadIngredientByName reads a ingredient recipe by name
func DBReadIngredientByName(name string) (Ingredient, error) {
	alling, err := DBReadAllIngredients() // Get all ingredients
	temp := Ingredient{}
	if err != nil {
		return temp, err
	}

	for _, i := range alling {
		if i.Name == name { // If name of parameter is in DB, return
			temp.ID = i.ID
			temp.Name = i.Name
			temp.Nutrients = i.Nutrients
			temp.Calories = i.Calories
			temp.Weight = i.Weight
			return temp, err
		}
	}
	return temp, err
}

// DBReadAllRecipes reads all recipes from database
func DBReadAllRecipes() ([]Recipe, error) {
	var temprecipes []Recipe //  Slice of all recipes, iterate over these
	iter := FireBaseDB.Client.Collection(RecipeCollection).Documents(FireBaseDB.Ctx)
	for {
		recipe := Recipe{} //  Create a placeholder for the document
		doc, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			fmt.Println("Failed to iterate over: "+RecipeCollection, err)
		}
		err = doc.DataTo(&recipe) // put data into temp struct
		if err != nil {
			fmt.Println("Error when converting retrieved document to struct: ", err)
		}

		temprecipes = append(temprecipes, recipe) // add to temp array
	}
	return temprecipes, nil
}

// DBReadAllIngredients reads all ingredients from database
func DBReadAllIngredients() ([]Ingredient, error) {
	var tempingredients []Ingredient
	ingredient := Ingredient{} //  Collects the entire collection
	iter := FireBaseDB.Client.Collection(IngredientCollection).Documents(FireBaseDB.Ctx)
	for {
		doc, err := iter.Next() //  Iterates over each document
		if err == iterator.Done {
			break
		}
		if err != nil {
			log.Fatalf("Failed to iterate: %v", err)
		}
		err = doc.DataTo(&ingredient) // Put data into temp struct
		if err != nil {
			fmt.Println("Error when converting retrieved document to struct: ", err)
		}
		tempingredients = append(tempingredients, ingredient) // Append to temp array
	}
	return tempingredients, nil
}

// DBReadAllWebhooks returns all registered webhooks in the database
func DBReadAllWebhooks() ([]Webhook, error) {
	var tempWebhooks []Webhook
	Wh := Webhook{}
	iter := FireBaseDB.Client.Collection(WebhooksCollection).Documents(FireBaseDB.Ctx)
	for {
		doc, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			log.Fatalf("Failed to iterate: %v", err)
		}

		err = doc.DataTo(&Wh) // put data into temp struct
		if err != nil {
			fmt.Println("Error when converting retrieved document to struct: ", err)
		}

		tempWebhooks = append(tempWebhooks, Wh)

	}
	return tempWebhooks, nil
}

//	DBCheckAuthorization func is used in the register functions to see if the user has authorization to upload data to the DB
//  The authorization token is one which we the creators of the program has created and saved manually
//  For security purposes we have chosen not to include code which saves the token, and the token itself can be given
//  to the reviewers of this project by mail which can be found in the readme
func DBCheckAuthorization(w http.ResponseWriter, r *http.Request) (bool, []byte) {
	tempToken := Token{}                //  Loop through collection of authorization tokens
	resp, err := ioutil.ReadAll(r.Body) //  Read the body of the json posted with the authentication token
	if err != nil {
		http.Error(w, "Couldn't read request: ", http.StatusBadRequest)
	}

	err = json.Unmarshal(resp, &tempToken)
	if err != nil {
		http.Error(w, "Unable to unmarshal request body: ", http.StatusBadRequest)
	}
	//  Loop through the collection of documents containing approved tokens
	iter := FireBaseDB.Client.Collection(TokenCollection).Documents(FireBaseDB.Ctx)
	for {
		DBToken := Token{}
		doc, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			http.Error(w, "Couldn't iterate over document colleciton : ", http.StatusInternalServerError)
		}
		err = doc.DataTo(&DBToken) // put data into temp struct
		if err != nil {
			http.Error(w, "Couldn't retrieve document from collection : ", http.StatusInternalServerError)
		}
		//  If the token the user posted is in the collection, return true
		if tempToken.AuthToken == DBToken.AuthToken {
			return true, resp
		}
	}
	return false, resp
}
