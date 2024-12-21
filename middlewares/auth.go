package middlewares

import (
	"kgermando/i-pos-restaurant-api/utils"

	"github.com/gofiber/fiber/v2"
)

func IsAuthenticated(c *fiber.Ctx) error {

	cookie := c.Cookies("token")

	if _, err := utils.VerifyJwt(cookie); err != nil {
		c.Status(fiber.StatusUnauthorized)
		return c.JSON(fiber.Map{
			"message": "unauthenticated",
		})
	}

	return c.Next()
}

// func IsAuthenticated(c *fiber.Ctx) error {

// 	// cookie := c.Cookies("token")
// 	tokenString := c.Get("Authorization")

// 	if tokenString == "" {
// 		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "No token present"})
// 	}

// 	parts := strings.Split(tokenString, " ")
// 	fmt.Println("parts", parts)
// 	if len(parts) != 2 || parts[0] != "Bearer" {
// 		return c.Status(fiber.StatusUnauthorized).SendString("Invalid authozation header format")
// 	}

// 	token := parts[1]

// 	if _, err := utils.VerifyJwt(token); err != nil {

// 		return c.JSON(fiber.Map{
// 			"message":      "unauthenticated",
// 			"Unauthorized": c.Status(fiber.StatusUnauthorized),
// 		})
// 	}

// 	return c.Next()
// }

 