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
			"name":"Kardemomme",
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