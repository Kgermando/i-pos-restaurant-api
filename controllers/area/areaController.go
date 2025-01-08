package area

import (
	"fmt"
	"kgermando/i-pos-restaurant-api/database"
	"kgermando/i-pos-restaurant-api/models"
	"strconv"

	"github.com/gofiber/fiber/v2"
)

// Paginate
func GetPaginatedArea(c *fiber.Ctx) error {
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

	var dataList []models.Area

	var length int64
	db.Model(dataList).Where("code_entreprise = ?", codeEntreprise).Count(&length)
	db.Where("code_entreprise = ?", codeEntreprise).
		Where("name ILIKE ? OR province ILIKE ?", "%"+search+"%", "%"+search+"%").
		Offset(offset).
		Limit(limit).
		Order("areas.updated_at DESC").
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
		"message":    "All areas",
		"data":       dataList,
		"pagination": pagination,
	})
}

// Get All data
func GetAllAreas(c *fiber.Ctx) error {
	codeEntreprise := c.Params("code_entreprise")
	db := database.DB

	var data []models.Area
	db.Where("code_entreprise = ?", codeEntreprise).Find(&data)
	return c.JSON(fiber.Map{
		"status":  "success",
		"message": "All areas",
		"data":    data,
	})
}

// Get one data
func GetArea(c *fiber.Ctx) error {
	id := c.Params("id")
	db := database.DB

	var area models.Area
	db.Find(&area, id)
	if area.Name == "" {
		return c.Status(404).JSON(
			fiber.Map{
				"status":  "error",
				"message": "No area found",
				"data":    nil,
			},
		)
	}
	return c.JSON(
		fiber.Map{
			"status":  "success",
			"message": "area found",
			"data":    area,
		},
	)
}

// Create data
func CreateArea(c *fiber.Ctx) error {
	p := &models.Area{}

	if err := c.BodyParser(&p); err != nil {
		return err
	}

	database.DB.Create(p)

	return c.JSON(
		fiber.Map{
			"status":  "success",
			"message": "area created success",
			"data":    p,
		},
	)
}

// Update data
func UpdateArea(c *fiber.Ctx) error {
	id := c.Params("id")
	db := database.DB

	type UpdateData struct {
		Name           string `json:"name"`
		Province       string `json:"province"`
		Signature      string `json:"signature"`
		CodeEntreprise uint64 `json:"code_entreprise"`
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

	area := new(models.Area)

	db.First(&area, id)
	area.Name = updateData.Name
	area.Province = updateData.Province
	area.Signature = updateData.Signature
	area.CodeEntreprise = updateData.CodeEntreprise

	db.Save(&area)

	return c.JSON(
		fiber.Map{
			"status":  "success",
			"message": "area updated success",
			"data":    area,
		},
	)

}

// Delete data
func DeleteArea(c *fiber.Ctx) error {
	id := c.Params("id")

	db := database.DB

	var area models.Area
	db.First(&area, id)
	if area.Name == "" {
		return c.Status(404).JSON(
			fiber.Map{
				"status":  "error",
				"message": "No area found",
				"data":    nil,
			},
		)
	}

	db.Delete(&area)

	return c.JSON(
		fiber.Map{
			"status":  "success",
			"message": "area deleted success",
			"data":    nil,
		},
	)
}
