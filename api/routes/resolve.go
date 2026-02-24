package routes

import (
	"log"

	"github.com/Silify/URLShortener/database"
	"github.com/gofiber/fiber/v3"
	"github.com/redis/go-redis/v9"
)

func ResolveURL(c fiber.Ctx) error {
	shortCode := c.Params("url")
	if shortCode == "" {
		return c.Status(fiber.StatusBadRequest).
			JSON(fiber.Map{"error": "missing short code"})
	}

	r := database.Client0

	value, err := r.Get(database.Ctx, shortCode).Result()
	if err == redis.Nil {
		return c.Status(fiber.StatusNotFound).
			JSON(fiber.Map{"error": "short url not found"})
	}
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).
			JSON(fiber.Map{"error": "internal server error"})
	}

	// Increment global counter
	rInr := database.Client1

	if err := rInr.Incr(database.Ctx, "counter").Err(); err != nil {
		log.Println(err)
	}

	return c.Redirect().Status(fiber.StatusFound).To(value)
}
