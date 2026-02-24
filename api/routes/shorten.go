package routes

import (
	"os"
	"strconv"
	"time"

	"github.com/Silify/URLShortener/database"
	"github.com/Silify/URLShortener/helpers"
	"github.com/asaskevich/govalidator"
	"github.com/gofiber/fiber/v3"
	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
)

type request struct {
	URL         string `json:"url"`
	CustomShort string `json:"short"`
	Expiry      int    `json:"expiry"` // hours
}

type response struct {
	URL            string `json:"url"`
	CustomShort    string `json:"short"`
	XRateRemaining int    `json:"rate_limit"`
	XRateLimitRest int64  `json:"rate_limit_reset"` // minutes
}

func ShortenURL(c fiber.Ctx) error {
	body := new(request)

	if err := c.Bind().Body(body); err != nil {
		return c.Status(fiber.StatusBadRequest).
			JSON(fiber.Map{"error": "cannot parse JSON"})
	}

	// Rate limit
	r2 := database.Client1

	val, err := r2.Get(database.Ctx, c.IP()).Result()
	if err == redis.Nil {
		_ = r2.Set(database.Ctx, c.IP(), os.Getenv("API_QUOTA"), 30*time.Minute).Err()
	} else if err != nil {
		return c.Status(fiber.StatusInternalServerError).
			JSON(fiber.Map{"error": "internal server error"})
	} else {
		vaInt, _ := strconv.Atoi(val)
		if vaInt <= 0 {
			limit, _ := r2.TTL(database.Ctx, c.IP()).Result()
			return c.Status(fiber.StatusServiceUnavailable).JSON(fiber.Map{
				"error":            "rate limit exceeded",
				"rate_limit_reset": int64(limit.Minutes()),
			})
		}
	}

	if !govalidator.IsURL(body.URL) {
		return c.Status(fiber.StatusBadRequest).
			JSON(fiber.Map{"error": "invalid URL"})
	}

	if !helpers.RemoveDomainError(body.URL) {
		return c.Status(fiber.StatusServiceUnavailable).
			JSON(fiber.Map{"error": "service unavailable"})
	}

	body.URL, _ = helpers.EnforceHTTP(body.URL)

	var id string
	if body.CustomShort == "" {
		id = uuid.New().String()[:6]
	} else {
		id = body.CustomShort
	}

	r := database.Client0

	val, _ = r.Get(database.Ctx, id).Result()
	if val != "" {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"error": "URL custom short is already in use",
		})
	}

	if body.Expiry == 0 {
		body.Expiry = 24
	}

	expiration := time.Duration(body.Expiry) * time.Hour

	if err := r.Set(database.Ctx, id, body.URL, expiration).Err(); err != nil {
		return c.Status(fiber.StatusInternalServerError).
			JSON(fiber.Map{"error": "internal server error"})
	}

	r2.Decr(database.Ctx, c.IP())

	remaining, _ := r2.Get(database.Ctx, c.IP()).Result()
	ttl, _ := r2.TTL(database.Ctx, c.IP()).Result()

	resp := response{
		URL:            body.URL,
		CustomShort:    os.Getenv("DOMAIN") + "/" + id,
		XRateRemaining: atoiSafe(remaining),
		XRateLimitRest: int64(ttl.Minutes()),
	}

	return c.Status(fiber.StatusCreated).JSON(resp)
}

func atoiSafe(s string) int {
	i, _ := strconv.Atoi(s)
	return i
}
