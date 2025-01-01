package commande

import (
	"fmt"
	"kgermando/i-pos-restaurant-api/database"
	"kgermando/i-pos-restaurant-api/models"
	"strconv" 

	"github.com/gofiber/fiber/v2"
)

// Paginate
func GetPaginatedComposition(c *fiber.Ctx) error {
	db := database.DB
	platId := c.Params("plat_id")

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

	var dataList []models.Composition

	var length int64
	db.Model(dataList).Where("plat_id = ?", platId).Count(&length)
	db.Where("plat_id = ?", platId).
		Joins("JOIN plats ON compositions.plat_id=plats.id").
		Where("plats.name ILIKE ? OR plats.reference ILIKE ?", "%"+search+"%", "%"+search+"%").
		Offset(offset).
		Limit(limit).
		Order("compositions.created_at DESC").
		Preload("Plat").
		Preload("Ingredient"). 
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
		"message":    "All compositions",
		"data":       dataList,
		"pagination": pagination,
	})
}

// Get data
func GetCompositionMargeBeneficiaire(c *fiber.Ctx) error {
	db := database.DB
	platId := c.Params("plat_id")

	var data models.Composition

	db.Model(data).Where("plat_id = ?", platId).Preload("Plat").Preload("Ingredient").Last(&data)

	return c.JSON(fiber.Map{
		"status":  "success",
		"message": "Total qty compositions",
		"data":    data,
	})
}

// Get Total data
func GetTotalComposition(c *fiber.Ctx) error {
	db := database.DB
	platId := c.Params("plat_id")

	var data []models.Composition
	var totalQty int64

	db.Model(data).Where("plat_id = ?", platId).Select("SUM(quantity)").Scan(&totalQty)

	return c.JSON(fiber.Map{
		"status":  "success",
		"message": "Total qty compositions",
		"data":    totalQty,
	})
}

// Get All data
func GetAllCompositions(c *fiber.Ctx) error {
	db := database.DB
	codeEntreprise := c.Params("code_entreprise")
	posId := c.Params("pos_id")

	var data []models.Composition
	db.Where("code_entreprise = ?", codeEntreprise).
	Where("pos_id = ?", posId).
	Preload("Ingredient").
	Find(&data)
	return c.JSON(fiber.Map{
		"status":  "success",
		"message": "All compositions",
		"data":    data,
	})
}

// Get one data
func GetComposition(c *fiber.Ctx) error {
	id := c.Params("id")
	db := database.DB

	var composition models.Composition
	db.Find(&composition, id)
	if composition.PlatID == 0 {
		return c.Status(404).JSON(
			fiber.Map{
				"status":  "error",
				"message": "No composition name found",
				"data":    nil,
			},
		)
	}
	return c.JSON(
		fiber.Map{
			"status":  "success",
			"message": "composition found",
			"data":    composition,
		},
	)
}

// Create data
func CreateComposition(c *fiber.Ctx) error {
	p := &models.Composition{}

	if err := c.BodyParser(&p); err != nil {
		return err
	}

	database.DB.Create(p)

	return c.JSON(
		fiber.Map{
			"status":  "success",
			"message": "composition created success",
			"data":    p,
		},
	)
}

// Update data
func UpdateComposition(c *fiber.Ctx) error {
	id := c.Params("id")
	db := database.DB

	type UpdateData struct {
		PlatID         uint       `json:"plat_id"`
		IngredientID   uint       `json:"ingredient_id"`
		Quantity       uint64     `json:"quantity"` 
		Signature      string     `json:"signature"`
		CodeEntreprise uint       `json:"code_entreprise"`
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

	composition := new(models.Composition)

	db.First(&composition, id)
	composition.PlatID = updateData.PlatID
	composition.IngredientID = updateData.IngredientID
	composition.Quantity = updateData.Quantity
	composition.Signature = updateData.Signature

	db.Save(&composition)

	return c.JSON(
		fiber.Map{
			"status":  "success",
			"message": "composition updated success",
			"data":    composition,
		},
	)

}

// Delete data
func DeleteComposition(c *fiber.Ctx) error {
	id := c.Params("id")

	db := database.DB

	var composition models.Composition
	db.First(&composition, id)
	if composition.PlatID == 0 {
		return c.Status(404).JSON(
			fiber.Map{
				"status":  "error",
				"message": "No composition name found",
				"data":    nil,
			},
		)
	}

	db.Delete(&composition)

	return c.JSON(
		fiber.Map{
			"status":  "success",
			"message": "composition deleted success",
			"data":    nil,
		},
	)
}
