package cravings

import (
	"context"
	"fmt"
	"log"

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
	app, err := firebase.NewApp(FireBaseDB.Ctx, nil, sa)		//  Creates the application with its contents
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
		fmt.Printf("ERROR deleting from collection: %v\n", err)
		return errors.Wrap(err, "Error in FirebaseDatabase.Delete()")
	}
	return nil
}

// DBReadRecipeByName reads a single recipe by Name
func DBReadRecipeByName(name string) (Recipe, error) {
	temp := Recipe{}		//  Recipe to be returned
	allrec, err := DBReadAllRecipes()	//  Query all the recipes 
	if err != nil {
		return temp, err
	}
	//  Loops through all the recipes and checks if the parameter name is equal to one of the recipes
	for _, i := range allrec {
		if i.RecipeName == name {
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
	alling, err := DBReadAllIngredients()
	temp := Ingredient{}
	if err != nil {
		return temp, err
	}

	for _, i := range alling {
		if i.Name == name {
			fmt.Println(name)
			temp.ID = i.ID
			temp.Name = i.Name
			temp.Nutrients = i.Nutrients

			return temp, err
		}
	}
	return temp, err
}

// DBReadRecipeByID reads a single recipe by ID
func DBReadRecipeByID(id string) (Recipe, error) {
	res := Recipe{} //  Creates an empty struct for the recipe
	//  Collects that document with given id from collection from firestore
	ref, err := FireBaseDB.Client.Collection(RecipeCollection).Doc(id).Get(FireBaseDB.Ctx)
	if err != nil {
		return res, err
	}
	err = ref.DataTo(&res)
	if err != nil {
		return res, err
	}
	return res, nil
}

// DBReadIngredientByID reads a ingredient recipe by ID
func DBReadIngredientByID(id string) (Ingredient, error) {
	res := Ingredient{} //  Creates empty struct
	//  Collects that document with given id from collection from firestore
	ref, err := FireBaseDB.Client.Collection(IngredientCollection).Doc(id).Get(FireBaseDB.Ctx)
	if err != nil {
		return res, err
	}
	err = ref.DataTo(&res)
	if err != nil {
		return res, err
	}
	return res, nil
}

// DBReadAllRecipes reads all recipes from database
func DBReadAllRecipes() ([]Recipe, error) {
	var temprecipes []Recipe

	iter := FireBaseDB.Client.Collection(RecipeCollection).Documents(FireBaseDB.Ctx)
	for {
		recipe := Recipe{}
		doc, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			log.Fatalf("Failed to iterate: %v", err)
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
	ingredient := Ingredient{}		//  Collects the entire collection
	iter := FireBaseDB.Client.Collection(IngredientCollection).Documents(FireBaseDB.Ctx)
	for {
		doc, err := iter.Next()		//  Iterates over each document
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
func DBCheckAuthorization(tokenParam string) bool {
	tempToken := Token{} //  Loop through collection of authorization tokens
	iter := FireBaseDB.Client.Collection(TokenCollection).Documents(FireBaseDB.Ctx)
	for {
		doc, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			log.Fatalf("Failed to iterate: %v", err)
		}
		err = doc.DataTo(&tempToken) // put data into temp struct
		if err != nil {
			fmt.Println("Error when converting retrieved document to struct: ", err)
		} //  If the user's token is in the collection, return true

		if tempToken.AuthToken == tokenParam {
			return true
		}
	}
	return false
}
