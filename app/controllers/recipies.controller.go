package controllers

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
	openrouter "github.com/shiroyaavish/open_router"
)

type NutritionalInfo struct {
	MinCalories int `json:"min_calories"`
	MaxCalories int `json:"max_calories"`
	MinProtien  int `json:"min_protien"`
	MaxProtien  int `json:"max_protien"`
	MinCarbs    int `json:"min_carbs"`
	MaxCarbs    int `json:"max_carbs"`
	MinFat      int `json:"min_fat"`
	MaxFat      int `json:"max_fat"`
}

var recipeRequest struct {
	Ingredients     []string        `json:"ingredients"`
	NutritionalInfo NutritionalInfo `json:"nutritional_info"`
	Notes           string          `json:"notes"`
	Request         string          `json:"request"`
}

type Recipe struct {
	CookingTime     float64            `json:"cooking_time_minutes"`
	Ingredients     []string           `json:"ingredients_required"`
	NutritionalInfo map[string]float64 `json:"nutritional_info"`
	Instructions    string             `json:"instructions"`
	RecipeName      string             `json:"recipe_name"`
}

var ProductRequest struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	Prompt      string `json:"prompt"`
	Tone        string `json:"tone"`
	Language    string `json:"language"`
}

type ProductNotification struct {
	Title        string `json:"title"`
	Notification string `json:"notification"`
}

func CreateRecipies(c *fiber.Ctx) error {
	if err := c.BodyParser(&recipeRequest); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	// Prepare the message for the API
	requestBody := map[string]interface{}{
		"model": "openrouter/quasar-alpha",
		"messages": []map[string]string{
			{
				"role":    "system",
				"content": "You are a trained cook who creates tasty and healthy recipes using only the ingredients available in the user's fridge. Your goal is to suggest creative and nutritious dishes based on the provided list of ingredients. Always ensure the recipe is practical, easy to follow, and uses only what's available in the input.",
			},
			{
				"role":    "user",
				"content": fmt.Sprintf("Generate a recipe using the following JSON data which contains the ingredients available in my fridge:\n\n%+v", recipeRequest),
			},
		},
		"response_format": map[string]interface{}{
			"type": "json_schema",
			"json_schema": map[string]interface{}{
				"name":   "recipe_data",
				"strict": true,
				"schema": map[string]interface{}{
					"type": "object",
					"properties": map[string]interface{}{
						"recipe_name": map[string]string{
							"type":        "string",
							"description": "The name of the recipe",
						},
						"ingredients_required": map[string]interface{}{
							"type":        "array",
							"description": "List of ingredients required for the recipe",
							"items": map[string]string{
								"type": "string",
							},
						},
						"nutritional_info": map[string]interface{}{
							"type":        "object",
							"description": "Nutritional information of the recipe",
							"properties": map[string]interface{}{
								"calories": map[string]string{
									"type":        "number",
									"description": "Calories in the recipe",
								},
								"protein": map[string]string{
									"type":        "number",
									"description": "Grams of protein in the recipe",
								},
								"carbs": map[string]string{
									"type":        "number",
									"description": "Grams of carbs in the recipe",
								},
								"fat": map[string]string{
									"type":        "number",
									"description": "Grams of fat in the recipe",
								},
							},
							"required":             []string{"calories", "protein", "carbs", "fat"},
							"additionalProperties": false,
						},
						"cooking_time_minutes": map[string]string{
							"type":        "number",
							"description": "Cooking time in minutes",
						},
						"instructions": map[string]string{
							"type":        "string",
							"description": "Step-by-step instructions for the recipe",
						},
					},
					"required":             []string{"recipe_name", "ingredients_required", "nutritional_info", "cooking_time_minutes", "instructions"},
					"additionalProperties": false,
				},
			},
		},
	}
	var response Recipe
	totalToken := openrouter.QuasarAlpha(requestBody, "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", &response)
	// if err != nil {
	// 	log.Println("Error:", err)
	// } else {
	// 	fmt.Println("Response:", response)
	// }

	return c.JSON(fiber.Map{
		"status":    200,
		"message":   "Data Get Successfully",
		"usedToken": totalToken,
		"data":      response,
	})

}
func CreateNotification(c *fiber.Ctx) error {
	if err := c.BodyParser(&ProductRequest); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	// Prepare the message for the API
	requestBody := map[string]interface{}{
		"model": "qwen/qwen3-30b-a3b:free",
		"messages": []map[string]string{
			{
				"role":    "system",
				"content": "You are the best marketing agent. who can create the mobile notifications with creative thought like attachment of the emotionally , romantically with fun & joy with the noitifaction to end user. And You have the mastery in the creativity thoughts. You can go as well with the festivals , current news and etc . within come in the one month. You can provide the multiple type of notification length between minimum 7 - maximum 9 and format in which is like given. most important things about you , your given all notification always unique.",
			},
			{
				"role":    "user",
				"content": fmt.Sprintf("Generate the notification title and message without reasoning with help proided data specially using prompt . provide the response in the json like [{title :'',notification:''}] :\n\n%+v", ProductRequest),
			},
		},
		"response_format": map[string]interface{}{
			"type": "json_schema",
			"json_schema": map[string]interface{}{
				"name":   "notification_data",
				"strict": true,
				"schema": map[string]interface{}{
					"type": "array",
					"items": map[string]interface{}{
						"properties": map[string]interface{}{
							"title": map[string]interface{}{
								"type":        "string",
								"description": "title of the notification in 30 - 40 character with unique and creativly and use the memes also.",
							},
							"notification": map[string]interface{}{
								"type":        "string",
								"description": "the description about sale , product and etc . using the provided message and create within 100 -110 charcters",
							},
						},
					},
					"required":             []string{"title", "notification"},
					"additionalProperties": true,
				},
			},
		},
	}
	var notification []ProductNotification
	tokensOfNotification := openrouter.QuasarAlpha(requestBody, "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", &notification)
	// if err != nil {
	// 	log.Println("Error:", err)
	// } else {
	// 	fmt.Println("Response:", response)
	// }
	notificationSelector := map[string]interface{}{
		"model": "qwen/qwen3-30b-a3b:free",
		"messages": []map[string]string{
			{
				"role":    "system",
				"content": "You are the head of marketing department. your work is select the most relatable notifications in the list which is given by your juniors. also , check the user prompt and match the requirements and notification list . then, select the appropriate notifications and return it if needed you can modified it.your strategies about you first rate the notifications. and approved the notification which has rate between 8 to 10. If notifications not appropriat then you create the 3 to 4 notifications. with the some creative ideas.",
			},
			{
				"role":    "user",
				"content": fmt.Sprintf("Check the given notification requirements \n\n%+v \n with the also given notification list:\n\n%+v \n\n and also return list which approved by you and without reasoning and any other data. provide the response in the json like [{title :'',notification:''}]", ProductRequest, notification),
			},
		},
		"response_format": map[string]interface{}{
			"type": "json_schema",
			"json_schema": map[string]interface{}{
				"name":   "approved_notification_data",
				"strict": true,
				"schema": map[string]interface{}{
					"type": "array",
					"items": map[string]interface{}{
						"type": "object",
						"properties": map[string]interface{}{
							"title": map[string]interface{}{
								"type":        "string",
								"description": "title of the notification...",
							},
							"notification": map[string]interface{}{
								"type":        "string",
								"description": "notification details...",
							},
						},
						"required":             []string{"title", "notification"},
						"additionalProperties": true,
					},
				},
			},
		},
	}
	var response []ProductNotification

	approvedNotificationTokens := openrouter.QuasarAlpha(notificationSelector, "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", &response)

	return c.JSON(fiber.Map{
		"status":    200,
		"message":   "Data Get Successfully",
		"usedToken": tokensOfNotification + approvedNotificationTokens,
		"notification":notification,
		"data":      response,
	})

}
