package tablebox

import (
	"fmt"
	"kgermando/i-pos-restaurant-api/database"
	"kgermando/i-pos-restaurant-api/models"
	"strconv"

	"github.com/gofiber/fiber/v2"
)

// Paginate
func GetPaginatedTableBoxEntreprise(c *fiber.Ctx) error {
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

	var dataList []models.TableBox

	var length int64
	db.Model(dataList).Where("code_entreprise = ?", codeEntreprise).Count(&length)
	db.Where("code_entreprise = ?", codeEntreprise).
		Where("name ILIKE ? OR numero ILIKE ?", "%"+search+"%", "%"+search+"%").
		Offset(offset).
		Limit(limit).
		Order("table_box.updated_at DESC").
		Preload("Commandes").
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
		"message":    "All commandes",
		"data":       dataList,
		"pagination": pagination,
	})
}

// Paginate
func GetPaginatedTableBox(c *fiber.Ctx) error {
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

	var dataList []models.TableBox

	var length int64
	db.Model(dataList).Where("code_entreprise = ?", codeEntreprise).
	Where("pos_id = ?", posId).Count(&length)
	db.Where("code_entreprise = ?", codeEntreprise).
		Where("pos_id = ?", posId).
		Where("name ILIKE ? OR numero ILIKE ?", "%"+search+"%", "%"+search+"%").
		Offset(offset).
		Limit(limit).
		Order("table_box.updated_at DESC").
		Preload("Commandes").
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
		"message":    "All commandes",
		"data":       dataList,
		"pagination": pagination,
	})
}

// Get All data
func GetAllTableBox(c *fiber.Ctx) error {
	codeEntreprise := c.Params("code_entreprise")
	db := database.DB

	var data []models.TableBox
	db.Where("code_entreprise = ?", codeEntreprise).Preload("Commandes").Find(&data)
	return c.JSON(fiber.Map{
		"status":  "success",
		"message": "All tableBox",
		"data":    data,
	})
}

// Get one data
func GetTableBox(c *fiber.Ctx) error {
	id := c.Params("id")
	db := database.DB

	var tableBox models.TableBox
	db.Preload("Commandes").Find(&tableBox, id)
	if tableBox.Name == "" {
		return c.Status(404).JSON(
			fiber.Map{
				"status":  "error",
				"message": "No tableBox name found",
				"data":    nil,
			},
		)
	}
	return c.JSON(
		fiber.Map{
			"status":  "success",
			"message": "tableBox found",
			"data":    tableBox,
		},
	)
}

// Create data
func CreateTableBox(c *fiber.Ctx) error {
	p := &models.TableBox{}

	if err := c.BodyParser(&p); err != nil {
		return err
	}

	database.DB.Create(p)

	return c.JSON(
		fiber.Map{
			"status":  "success",
			"message": "tableBox created success",
			"data":    p,
		},
	)
}

// Update data
func UpdateTableBox(c *fiber.Ctx) error {
	id := c.Params("id")
	db := database.DB

	type UpdateData struct {
		PosID          uint   `json:"pos_id"`
		Name           string `json:"name"`
		Numero         int    `json:"numero"`
		Status         string `json:"status"`
		Signature      string `json:"signature"`
		CodeEntreprise uint   `json:"code_entreprise"`
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

	tableBox := new(models.TableBox)

	db.First(&tableBox, id)
	tableBox.PosID = updateData.PosID
	tableBox.Name = updateData.Name
	tableBox.Numero = updateData.Numero
	tableBox.Status = updateData.Status
	tableBox.Signature = updateData.Signature
	tableBox.CodeEntreprise = updateData.CodeEntreprise

	db.Save(&tableBox)

	return c.JSON(
		fiber.Map{
			"status":  "success",
			"message": "tableBox updated success",
			"data":    tableBox,
		},
	)

}

// Delete data
func DeleteTableBox(c *fiber.Ctx) error {
	id := c.Params("id")

	db := database.DB

	var tableBox models.TableBox
	db.First(&tableBox, id)
	if tableBox.Name == "" {
		return c.Status(404).JSON(
			fiber.Map{
				"status":  "error",
				"message": "No tableBox name found",
				"data":    nil,
			},
		)
	}

	db.Delete(&tableBox)

	return c.JSON(
		fiber.Map{
			"status":  "success",
			"message": "tableBox deleted success",
			"data":    nil,
		},
	)
}
