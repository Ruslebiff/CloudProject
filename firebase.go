package cravings

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	firebase "firebase.google.com/go"
	"github.com/pkg/errors"
	"google.golang.org/api/iterator"
	"google.golang.org/api/option"
)

// fireBaseDB is an instance of FirestoreDatabase struct which is used in firebase.go
var fireBaseDB = FirestoreDatabase{}

// DBInit initialises the database
func DBInit() error {
	// Firebase initialisation
	fireBaseDB.Ctx = context.Background()
	// We use a service account, load credentials file that you downloaded from your project's settings menu.
	// Make sure this file is gitignored, it is the access token to the database.
	sa := option.WithCredentialsFile(FirestoreCredentials)
	app, err := firebase.NewApp(fireBaseDB.Ctx, nil, sa) //  Creates the application with its contents

	if err != nil {
		fmt.Println("Failed to initialize the firebase database when creating a new app: ", err)
		return err
	}
	//  Sets the app created to our local struct's client
	fireBaseDB.Client, err = app.Firestore(fireBaseDB.Ctx)
	if err != nil {
		fmt.Println("Failed to create app")
		return err
	}

	return err
}

// DBClose Close firebase connection
func DBClose() {
	err := fireBaseDB.Client.Close()
	if err != nil {
		fmt.Println("Failed to close firebase client")
	} else {
		fmt.Println("Successfully closed firebase client")
	}
}

// DBSaveRecipe saves recipe to database
func DBSaveRecipe(r *Recipe, w http.ResponseWriter) error { //  Creates a new document in firebase
	ref := fireBaseDB.Client.Collection(RecipeCollection).NewDoc()
	r.ID = ref.ID                        //  Asserts the recipes id to be the one given by firebase
	_, err := ref.Set(fireBaseDB.Ctx, r) //  Set the context of the document to the one of the recipe

	if err != nil {
		fmt.Fprintln(w, "ERROR saving recipe to recipe collection: "+err.Error(), http.StatusInternalServerError)
		return err
	}

	return nil
}

// DBSaveIngredient saves ingredient to database
func DBSaveIngredient(i *Ingredient, w http.ResponseWriter) error { //  Creates a new document in firebase
	ref := fireBaseDB.Client.Collection(IngredientCollection).NewDoc()
	i.ID = ref.ID                        //  Asserts the ingredients id to be the one given by firebase
	_, err := ref.Set(fireBaseDB.Ctx, i) //  Set the context of the document to the one of the ingredient

	if err != nil {
		fmt.Fprintln(w, "ERROR saving ingredient to ingredients collection: "+err.Error(), http.StatusInternalServerError)
		return err
	}

	return nil
}

// DBSaveWebhook saves a new webhook to the database
func DBSaveWebhook(i *Webhook, w http.ResponseWriter) error {
	ref := fireBaseDB.Client.Collection(WebhooksCollection).NewDoc()
	i.ID = ref.ID                        //  Asserts the webhooks id to be the one given by firebase
	_, err := ref.Set(fireBaseDB.Ctx, i) //  Set the context of the document to the one of the webhook

	if err != nil {
		fmt.Fprintln(w, "ERROR saving webhook to webhooks collection: "+err.Error(), http.StatusInternalServerError)
		return err
	}

	return nil
}

// DBDelete deletes an entry from given collection in database by its id, either ingredient, recipe or webhook
func DBDelete(id string, collection string, w http.ResponseWriter) error {
	_, err := fireBaseDB.Client.Collection(collection).Doc(id).Delete(fireBaseDB.Ctx)
	if err != nil {
		fmt.Fprintln(w, "ERROR deleting from collection: "+collection+err.Error(), http.StatusBadRequest)
		return errors.Wrap(err, "Error in FirebaseDatabase.Delete()")
	}

	return nil
}

// DBReadRecipeByName reads a single recipe by Name
func DBReadRecipeByName(name string, w http.ResponseWriter) (Recipe, error) {
	temp := Recipe{}                   //  Recipe to be returned
	allrec, err := DBReadAllRecipes(w) //  Query all the recipes

	if err != nil {
		fmt.Fprintln(w, "Error retrieving recipes from database "+err.Error(), http.StatusInternalServerError)
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
func DBReadIngredientByName(name string, w http.ResponseWriter) (Ingredient, error) {
	alling, err := DBReadAllIngredients(w) // Get all ingredients
	temp := Ingredient{}

	if err != nil {
		fmt.Fprintln(w, "Error retrieving recipes from database "+err.Error(), http.StatusInternalServerError)
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
func DBReadAllRecipes(w http.ResponseWriter) ([]Recipe, error) {
	var temprecipes []Recipe //  Slice of all recipes, iterate over these

	iter := fireBaseDB.Client.Collection(RecipeCollection).Documents(fireBaseDB.Ctx)

	for {
		recipe := Recipe{} //  Create a placeholder for the document
		doc, err := iter.Next()

		if err == iterator.Done {
			break
		}

		if err != nil {
			fmt.Fprintln(w, "Failed to iterate "+err.Error(), http.StatusInternalServerError)
			return temprecipes, err
		}

		err = doc.DataTo(&recipe) // put data into temp struct

		if err != nil {
			fmt.Fprintln(w, "Error when converting retrieved document to struct: "+err.Error(), http.StatusInternalServerError)

			return temprecipes, err
		}

		temprecipes = append(temprecipes, recipe) // add to temp array
	}

	return temprecipes, nil
}

// DBReadAllIngredients reads all ingredients from database
func DBReadAllIngredients(w http.ResponseWriter) ([]Ingredient, error) {
	var tempingredients []Ingredient

	ingredient := Ingredient{} //  Collects the entire collection
	iter := fireBaseDB.Client.Collection(IngredientCollection).Documents(fireBaseDB.Ctx)

	for {
		doc, err := iter.Next() //  Iterates over each document
		if err == iterator.Done {
			break
		}

		if err != nil {
			fmt.Fprintln(w, "Failed to iterate "+err.Error(), http.StatusInternalServerError)
			return tempingredients, err
		}

		err = doc.DataTo(&ingredient) // Put data into temp struct
		if err != nil {
			fmt.Fprintln(w, "Error when converting retrieved document to struct: "+err.Error(), http.StatusInternalServerError)
			return tempingredients, err
		}

		tempingredients = append(tempingredients, ingredient) // Append to temp array
	}

	return tempingredients, nil
}

// DBReadAllWebhooks returns all registered webhooks in the database
func DBReadAllWebhooks(w http.ResponseWriter) ([]Webhook, error) {
	var tempWebhooks []Webhook

	Wh := Webhook{}

	iter := fireBaseDB.Client.Collection(WebhooksCollection).Documents(fireBaseDB.Ctx)

	for {
		doc, err := iter.Next()
		if err == iterator.Done {
			break
		}

		if err != nil {
			fmt.Fprintln(w, "Failed to iterate: "+err.Error(), http.StatusInternalServerError)
			return tempWebhooks, err
		}

		err = doc.DataTo(&Wh) // put data into temp struct
		if err != nil {
			fmt.Fprintln(w, "Error when converting retrieved document to struct: "+err.Error(), http.StatusInternalServerError)
			return tempWebhooks, err
		}

		tempWebhooks = append(tempWebhooks, Wh)
	}

	return tempWebhooks, nil
}

// DBCheckAuthorization func is used in the register functions
// to see if the user has authorization to upload data to the DB
// The authorization token is one which we the creators of the program has created and saved manually
// For security purposes we have chosen not to include code which saves the token, and the token itself can be given
// to the reviewers of this project by mail which can be found in the readme
func DBCheckAuthorization(w http.ResponseWriter, r *http.Request) (bool, []byte) {
	tempToken := Token{}                //  Loop through collection of authorization tokens
	resp, err := ioutil.ReadAll(r.Body) //  Read the body of the json posted with the authentication token

	if err != nil {
		fmt.Fprintln(w, "Couldn't read request: "+err.Error(), http.StatusBadRequest)
	}

	err = json.Unmarshal(resp, &tempToken)
	if err != nil {
		fmt.Fprintln(w, "Unable to unmarshal request body: "+err.Error(), http.StatusBadRequest)
	}
	//  Loop through the collection of documents containing approved tokens
	iter := fireBaseDB.Client.Collection(TokenCollection).Documents(fireBaseDB.Ctx)

	for {
		DBToken := Token{}
		doc, err := iter.Next()

		if err == iterator.Done {
			break
		}

		if err != nil {
			fmt.Fprintln(w, "Couldn't iterate over document collection: "+err.Error(), http.StatusInternalServerError)
		}

		err = doc.DataTo(&DBToken) // put data into temp struct
		if err != nil {
			fmt.Fprintln(w, "Couldn't retrieve document from collection: "+err.Error(), http.StatusInternalServerError)
		}
		//  If the token the user posted is in the collection, return true
		if tempToken.AuthToken == DBToken.AuthToken {
			return true, resp
		}
	}

	return false, resp
}
