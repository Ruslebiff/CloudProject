# CloudProject
NTNU Cloud Technologies project 2019

# Register
Register Ingredient: 
{
	"token":"",
	"name":"Kardemomme",
	"unit":""
}
Unit should be either "l" or "g". Please use "l" for ingredients measured by volume and "g" for ingredients measured by weight

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
Delete Ingredient/Recipe: 
{
	"token":"",
	"name":""
}


mealHandler:
	Get method:
		URL: http://localhost:8080/cravings/meal/?ingredients=cheese_milk|{70}_cardamom|{500}|{g}	{:} = optional
		'_' splits up the different ingredients
		'|' splits up the ingredient, quantity and unit (in this given order)
		if quantity is not set or is not valid, it is put to a default value

	Post method:
[	
	{
		"name": "cheese",
		"unit": "g"
		"quantity": 1
	},
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