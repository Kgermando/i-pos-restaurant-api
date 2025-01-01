package ingredient

import (
	"fmt"
	"kgermando/i-pos-restaurant-api/database"
	"kgermando/i-pos-restaurant-api/models"
	"strconv"

	"github.com/gofiber/fiber/v2"
)

// Paginate
func GetPaginatedIngredientEntreprise(c *fiber.Ctx) error {
	db := database.DB
	codeEntreprise := c.Params("code_entreprise")

	page, err := strconv.Atoi(c.Query("page", "1"))
	if err != nil || page <= 0 {
		page = 1 // Default page number
	}
	limit, err := strconv.Atoi(c.Query("limit", "15"))
	if err != nil || limit <= 0 {
		limit = 15
	}
	offset := (page - 1) * limit

	search := c.Query("search", "")

	var dataList []models.Ingredient

	var length int64
	db.Model(dataList).Where("code_entreprise = ?", codeEntreprise).Count(&length)
	db.Where("code_entreprise = ?", codeEntreprise).
		Where("name ILIKE ?", "%"+search+"%").
		Offset(offset).
		Limit(limit).
		Order("ingredients.updated_at DESC").
		Preload("Pos").
		Find(&dataList)

	if err != nil {
		fmt.Println("error s'est produite: ", err)
		return c.Status(500).SendString(err.Error())
	}

	// Calculate total number of pages
	totalPages := len(dataList) / limit
	if remainder := len(dataList) % limit; remainder > 0 {
		totalPages++
	}
	pagination := map[string]interface{}{
		"total_pages": totalPages,
		"page":        page,
		"page_size":   limit,
		"length":      length,
	}

	return c.JSON(fiber.Map{
		"status":     "success",
		"message":    "All ingredients",
		"data":       dataList,
		"pagination": pagination,
	})
}

// Paginate
func GetPaginatedIngredient(c *fiber.Ctx) error {
	db := database.DB
	codeEntreprise := c.Params("code_entreprise")
	posId := c.Params("pos_id")

	page, err := strconv.Atoi(c.Query("page", "1"))
	if err != nil || page <= 0 {
		page = 1 // Default page number
	}
	limit, err := strconv.Atoi(c.Query("limit", "15"))
	if err != nil || limit <= 0 {
		limit = 15
	}
	offset := (page - 1) * limit

	search := c.Query("search", "")

	var dataList []models.Ingredient

	var length int64
	db.Model(dataList).Where("code_entreprise = ?", codeEntreprise).
		Where("pos_id = ?", posId).Count(&length)
	db.Where("code_entreprise = ?", codeEntreprise).
		Where("pos_id = ?", posId).
		Where("name ILIKE ?", "%"+search+"%").
		Offset(offset).
		Limit(limit).
		Order("ingredients.updated_at DESC").
		Preload("Pos").
		Find(&dataList)

	if err != nil {
		fmt.Println("error s'est produite: ", err)
		return c.Status(500).SendString(err.Error())
	}

	// Calculate total number of pages
	totalPages := len(dataList) / limit
	if remainder := len(dataList) % limit; remainder > 0 {
		totalPages++
	}
	pagination := map[string]interface{}{
		"total_pages": totalPages,
		"page":        page,
		"page_size":   limit,
		"length":      length,
	}

	return c.JSON(fiber.Map{
		"status":     "success",
		"message":    "All ingredients",
		"data":       dataList,
		"pagination": pagination,
	})
}

// Get All data
func GetAllIngredients(c *fiber.Ctx) error {
	db := database.DB
	codeEntreprise := c.Params("code_entreprise")
	posId := c.Params("pos_id")

	var data []models.Ingredient
	db.Where("code_entreprise = ?", codeEntreprise).
		Where("pos_id = ?", posId).
		Preload("Pos").
		Find(&data)
	return c.JSON(fiber.Map{
		"status":  "success",
		"message": "All ingredients",
		"data":    data,
	})
}

// Get All data by id
func GetAllIngredientBySearch(c *fiber.Ctx) error {
	db := database.DB
	codeEntreprise := c.Params("code_entreprise")
	posId := c.Params("pos_id")

	search := c.Query("search", "")

	var data []models.Ingredient
	db.Where("code_entreprise = ?", codeEntreprise).
		Where("pos_id = ?", posId).
		Where("name ILIKE ?", "%"+search+"%").
		Preload("Pos").
		Find(&data)
	return c.JSON(fiber.Map{
		"status":  "success",
		"message": "All ingredients",
		"data":    data,
	})
}

// Get one data
func GetIngredient(c *fiber.Ctx) error {
	id := c.Params("id")
	db := database.DB

	var ingredient models.Ingredient
	db.Find(&ingredient, id)
	if ingredient.Name == "" {
		return c.Status(404).JSON(
			fiber.Map{
				"status":  "error",
				"message": "No ingredient name found",
				"data":    nil,
			},
		)
	}
	return c.JSON(
		fiber.Map{
			"status":  "success",
			"message": "ingredient found",
			"data":    ingredient,
		},
	)
}

// Create data
func CreateIngredient(c *fiber.Ctx) error {
	p := &models.Ingredient{}

	if err := c.BodyParser(&p); err != nil {
		return err
	}

	database.DB.Create(p)

	return c.JSON(
		fiber.Map{
			"status":  "success",
			"message": "ingredient created success",
			"data":    p,
		},
	)
}

// Update data
func UpdateIngredient(c *fiber.Ctx) error {
	id := c.Params("id")
	db := database.DB

	type UpdateData struct {
		Name           string  `json:"name"`
		Description    string  `json:"description"` 
		Unite          string  `json:"unite"`
		PosID          uint    `json:"pos_id"`
		Signature      string  `json:"signature"`
		CodeEntreprise uint    `json:"code_entreprise"`
	}

	var updateData UpdateData

	if err := c.BodyParser(&updateData); err != nil {
		return c.Status(500).JSON(
			fiber.Map{
				"status":  "error",
				"message": "Review your iunput",
				"data":    nil,
			},
		)
	}

	ingredient := new(models.Ingredient)

	db.First(&ingredient, id)
	ingredient.Name = updateData.Name
	ingredient.Description = updateData.Description
	ingredient.Unite = updateData.Unite
	ingredient.Signature = updateData.Signature
	ingredient.PosID = updateData.PosID
	ingredient.CodeEntreprise = updateData.CodeEntreprise

	db.Save(&ingredient)

	return c.JSON(
		fiber.Map{
			"status":  "success",
			"message": "ingredient updated success",
			"data":    ingredient,
		},
	)

}

// Delete data
func DeleteIngredient(c *fiber.Ctx) error {
	id := c.Params("id")

	db := database.DB

	var ingredient models.Ingredient
	db.First(&ingredient, id)
	if ingredient.Name == "" {
		return c.Status(404).JSON(
			fiber.Map{
				"status":  "error",
				"message": "No ingredient name found",
				"data":    nil,
			},
		)
	}

	db.Delete(&ingredient)

	return c.JSON(
		fiber.Map{
			"status":  "success",
			"message": "ingredient deleted success",
			"data":    nil,
		},
	)
}
