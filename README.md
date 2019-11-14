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

# HandelerMeal
	mealHandler:
		Get method:
			URL: http://localhost:8080/cravings/meal/?ingredients=cheese_milk|{70}_cardamom|{500}|{g}	{:} = optional
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