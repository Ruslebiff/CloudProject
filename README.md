# CloudProject
NTNU Cloud Technologies project 2019

# Food 
To register you have to post the body in a json format 

	# Register ingredient: cravings/food/ingredient
	Register Ingredient: 
	{
		"token":"",
		"name":"Kardemomme",
		"unit":""
	}

Unit should be either "l" or "g". Please use "l" for ingredients measured by volume and "g" for ingredients measured by weight


	# Register recipe: cravings/food/recipe
	
	Register Recipe:
	{
		"token":"",
		"recipeName":"Emils Kakoramarama2",
		"ingredients":[
			{
				"name":"cardamom",
				"quantity":5,
				"unit":"kg"
			},
			{
				"name":"milk",
				"quantity":69,
				"unit":"l"
			}
		]
	}

# Delete
	cravings/food/... 

	Delete Ingredient/Recipe: 
	{
		"token":"",
		"name":""
	}

# HandlerMeal
	mealHandler:
		Get method:
			URL: ./cravings/meal/?ingredients=cheese_milk|{70}_cardamom|{500}|{g}	{:} = optional
			'_' splits up the different ingredients
			'|' splits up the ingredient, quantity and unit (in this given order)
			if quantity is not set or is not valid, it is put to a default value (but it is highly recomended to write in all the information)

		Post method:
	[	
		{
			"name": "[ingredient name]",
			"unit": "[ingredient unit]",
			"quantity": [ingredient quantity]
		}
	]
	ex:
	[
		{
			"name": "cardamom",
			"unit": "g",
			"quantity": 500
			
		},
		{
			"name": "milk",
			"unit": "l",
			"quantity": 70
				
		}
	]
list as many ingredients with quantity and unit as you want

The user can send a post request with the payload of the 'remaining' struct of any given recipe to get the recipe for 'the next meal'. This process can be done repeatedly until the 'remaining' list is empty.

	limit: int, sets to 5 as default
	allowMissing: bool, true as default. Decides wether or not to print out recipes that are missing ingredients
	sortBy: "have"|"missing"|"remaining". have sorts in a descending order, missing and remaining sorts in an ascending order

# Webhooks

Webhooks endpoint: /cravings/webhooks/
Here you can get information about webhooks for this website

Post method:
Post is used to create new webhooks.
Use endpoint:
/cravings/webhooks/

	And send with body:
	{
	"event":"[Event name]",
	"url":"[Url name]"
	}

Get method:
Get is used to see all or one choosen webhook.
To get all webhooks, use normal endpoint:
/cravings/webhooks/

To get one webhook, use normal endpoint + choosen id for webhook:
/cravings/webhooks/[ID]

Delete method:
Delete is used to delete one webhook.
Use endpoint:
/cravings/webhooks/

	And send with body:
	{
	"id":"[ID]"
	}

# Test
Test cover = %
Test coverage can be tested by entering following command in termina: go test -cover

# Docker
A Dockerfile is included in this repository. This is tested to work with our build and the following commands. 

Example command for building docker image: 

	docker build -t cravings .

(working directory should be in the repository when executing the build command)

Example command for running the container: 

	docker run -i -t -p 8081:8080 cravings

(this will run it on port 8081 on the host machine)


# Original project plan

Short description:
We will make an API that can be used to get meal ideas from what ingredients you already have. This API could for example be used by a website or app providing a GUI to the users.

We have a database containing ingredients and recipes, including nutritional info. Each request to our API reads data from the database. When registering new recipes or ingredients to the database, it will get the nutritional info from an external API (Edamam). 
The project will use both OpenStack and Docker, and store data in Firebase. 

Example usage: 
You post what ingredients you already have, including how much of each ingredient. In return you get recipes you can make using these ingredients, and also nutritional info about it (calories, fat, carbohydrates etc.). 

Another example:
User can read every ingredient/recipe in the database including its nutrients. 
Furthermore, user can get one recipe or ingredient by name

Getting recipes:
What ingredients you have is provided either via URL or preferrably using POST request with a JSON body. Recipes that can be made using the ingredients you already have will be returned. These will also include nutritional info for the recipe.

Registration: 
You will need an auth token provided by us to get access to register new recipes or ingredients to the database. Tokens are stored in our database in a separate collection. 

Registration is done by sending a POST request to our registration handler for ingredients or recipe, including a JSON structure in body. We will provide templates for this. 

When something new is registered, we get the nutritional info for it from the Edamam API.

Webhooks: 
Webhooks for seeing what’s registered into the database through the /register/ handler. This includes both recipes and ingredients.

Statistics: 
Statistics indicates the availability of the database used in the assignment, and the website of which the program retrieves information.  In addition, it indicates time elapsed since the start of the program. Last but not least, it indicates how many recipes and ingredients are stored in the database. 


Potential expansions of project:
The user also get what he/she has left after using a recipe, and which recipes you also can make afterwards without having to buy new ingredients. The program will also calculate how many days these resources will last, considering an average male needs about 2200 kcal each day. 

Can register several recipes in one POST
Registration of recipes and ingredients could be done automatically via some external API or website. 

User requests a recipe, inserts what it has of ingredients. The system provides a “shopping list” 


# What went well and wrong 
reflection of what went well and what went wrong with the project  
After about an hour into our first meeting we had layed out a project plan of what our final product should look like.
We had good working routines, meeting as a group everyday to work.

We managed to reach all of our main goals, and added some of the potential expansions for the project.

When the user 'uses' a recipe, there is a list of all ingredients he/she has for the recipe(have), needs to complete the recipe(missing) and ingredients after making the recipe(remaining). To find recipes for the next day, the application just posts a new request with the remaining list.
We also different queries to handlerMeal (look at #handlerMeal for more information)


# Hard aspects of the project
reflection on the hard aspects of the project   
Recipes that has units in teaspoon or tablespoon values became a bigger problem fixing than expected. The solution we decided to go for was calculating how many calories there was per spoon and from there get the quantity per unit. This lead to extra lines of code only for handling spoon units.


# What we learned
what new has the group learned
We got a deeper insight in how it is to make an API database that is meant to be used by other applications. 


# Work log
total work hours dedicated to the project cumulatively by the group

Tuesday 5/11-19 group work from 13.00 to 15.00 total time 8 hours (2 hours per person) 
(12-13 was just idea discussion, not working)

Wednesday 6/11-19 group work from 09:30 - 14:30 total time 20 hours (5 hours per person)

Thursday 7/11-19 group work from 12:15 - 14:00 total time 7 hours (1 hour and 45 min per person)

Friday 8/11-19 group work from 10:00 - 17:00

Saturday 9/11-19 group work from 12:00 — 18:00

Sunday 10/11-19 group work from 12:00 - 17:00

Monday 11/11-19 group work from 10:00 - 16:00

Tuesday 12/11-19 group work from 10:30 - 17:30

Wednesday 13/11-19 group work from 10:00 - 17:00

Thursday 14/11 - 19 group work from 12:00 - 16:00

Friday 15/11 - 19 group work from 11:00 - 