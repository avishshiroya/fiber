package controllers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	// "log"
	"net/http"

	"github.com/gofiber/fiber/v2"
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

	jsonData, err := json.Marshal(requestBody)
	if err != nil {
		panic(err)
	}
	// Create the HTTP request
	req, err := http.NewRequest("POST", "https://openrouter.ai/api/v1/chat/completions", bytes.NewBuffer(jsonData))
	if err != nil {
		panic(err)
	}

	// Set headers
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer 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")

	// Make the request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}

	defer resp.Body.Close()
	result, err := io.ReadAll(resp.Body)

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Error in connectivity with cook.",
			"data":  err,
		})
	}
	var responseData map[string]interface{}
	if err := json.Unmarshal(result, &responseData); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to unmarshal response",
		})
	}

	// Debugging response
	// fmt.Println("API Response:", responseData)

	// Safely extract 'choices'
	choicesRaw, ok := responseData["choices"]
	if !ok {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Missing 'choices' in API response",
		})
	}

	choices, ok := choicesRaw.([]interface{})
	if !ok || len(choices) == 0 {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Invalid or empty 'choices' format",
		})
	}

	choice, ok := choices[0].(map[string]interface{})
	if !ok {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Invalid 'choice' structure",
		})
	}

	message, ok := choice["message"].(map[string]interface{})
	if !ok {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Missing 'message' in choice",
		})
	}

	content, ok := message["content"].(string)
	if !ok {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Missing 'content' in message",
		})
	}

	// matches := cleanMarkdownCodeBlock(content)
	var response Recipe
	// // json.Unmarshal([]byte(content), &response)
	if err := json.Unmarshal([]byte(content), &response); err != nil {
		fmt.Println(err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Error at response unmarshel",
		})
	}
	// encodedJSON, err := json.MarshalIndent(response, "", "  ")
	// if err != nil {
	// 	log.Fatal(err)
	// }

	return c.JSON(fiber.Map{
		"status":  200,
		"message": "Data Get Successfully",
		"data":    response,
	})

}
