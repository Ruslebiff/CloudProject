package cravings

import (
	"context"
	"fmt"
	"log"

	firebase "firebase.google.com/go"
	"github.com/pkg/errors"
	"google.golang.org/api/option"
)

var FireBaseDB = FirestoreDatabase{}

func DBInit() {
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
		log.Fatalln(err)
	}

}

//  Func creates document which is to store a webhook
func DBSaveRecipe(r *Recipe) error { //  Creates a new document in firebase
	ref := FireBaseDB.Client.Collection(RecipeCollection).NewDoc()
	r.ID = r.ID                          //  Asserts the webhooks id to be the one given by firebase
	_, err := ref.Set(FireBaseDB.Ctx, r) //  Set the context of the document to the one of the webhook
	if err != nil {
		fmt.Println("ERROR saving recipe to Firestore DB: ", err)
	}
	return nil
}

//  Func deletes a webhook by an id
func DBDeleteRecipe(id string) error {
	_, err := FireBaseDB.Client.Collection(RecipeCollection).Doc(id).Delete(FireBaseDB.Ctx)
	if err != nil {
		fmt.Printf("ERROR deleting student from Firestore DB: %v\n", err)
		return errors.Wrap(err, "Error in FirebaseDatabase.Delete()")
	}
	return nil
}

//  Retrieves a document to read
func DBReadRecipeByID(id string) (Recipe, error) {
	res := Recipe{} // Gets a webhook by a specific id
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
