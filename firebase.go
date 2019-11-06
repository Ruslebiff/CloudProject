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

// FireBaseDB is an instance of FirestoreDatabase struct ------ Flytt til globalfil senere
var FireBaseDB = FirestoreDatabase{}

// DBInit initialises the database
func DBInit() error {
	// Firebase initialisation
	FireBaseDB.Ctx = context.Background()
	// We use a service account, load credentials file that you downloaded from your project's settings menu.
	// Make sure this file is gitignored, it is the access token to the database.
	sa := option.WithCredentialsFile(FirestoreCredentials)
	app, err := firebase.NewApp(FireBaseDB.Ctx, nil, sa)
	if err != nil {
		fmt.Println("Failed to initialize the firebase database when creating a new app: ", err)
	}

	FireBaseDB.Client, err = app.Firestore(FireBaseDB.Ctx)
	if err != nil {
		fmt.Println("Failed to create app")
	}
	return err
}

// DBSaveRecipe saves recipe to database
func DBSaveRecipe(r *Recipe) error { //  Creates a new document in firebase
	ref := FireBaseDB.Client.Collection(RecipeCollection).NewDoc()
	r.ID = ref.ID                        //  Asserts the webhooks id to be the one given by firebase
	_, err := ref.Set(FireBaseDB.Ctx, r) //  Set the context of the document to the one of the webhook
	if err != nil {
		fmt.Println("ERROR saving recipe to recipe collection: ", err)
	}
	return nil
}

// DBSaveIngredients saves ingredient to database
func DBSaveIngredient(i *Ingredient) error { //  Creates a new document in firebase
	ref := FireBaseDB.Client.Collection(IngredientCollection).NewDoc()
	i.ID = ref.ID                        //  Asserts the webhooks id to be the one given by firebase
	_, err := ref.Set(FireBaseDB.Ctx, i) //  Set the context of the document to the one of the webhook
	if err != nil {
		fmt.Println("ERROR saving ingredient to ingredients collection: ", err)
	}
	return nil
}

// DBDelete deletes an entry from given collection in database
func DBDelete(id string, collection string) error {
	_, err := FireBaseDB.Client.Collection(collection).Doc(id).Delete(FireBaseDB.Ctx)
	if err != nil {
		fmt.Printf("ERROR deleting from collection: %v\n", err)
		return errors.Wrap(err, "Error in FirebaseDatabase.Delete()")
	}
	return nil
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
	recipe := Recipe{}
	iter := FireBaseDB.Client.Collection(RecipeCollection).Documents(FireBaseDB.Ctx)
	for {
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
	ingredient := Ingredient{}
	iter := FireBaseDB.Client.Collection(IngredientCollection).Documents(FireBaseDB.Ctx)
	for {
		doc, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			log.Fatalf("Failed to iterate: %v", err)
		}
		err = doc.DataTo(&ingredient) // put data into temp struct
		if err != nil {
			fmt.Println("Error when converting retrieved document to struct: ", err)
		}

		temprecipes = append(tempingredients, ingredient) // add to temp array

	}
	return tempingredients, nil
}
