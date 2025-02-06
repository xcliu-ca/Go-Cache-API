package middleware

import (
    "encoding/json"
    "time"

    "github.com/gofiber/fiber/v2"
    "github.com/patrickmn/go-cache"
)

type Post struct {
    UserID int    `json:"userId"`
    ID     int    `json:"id"`
    Title  string `json:"title"`
    Body   string `json:"body"`
}

func CacheMiddleware(cache *cache.Cache) fiber.Handler {
    return func(c *fiber.Ctx) error {
        if c.Method() != "GET" {
            // Only cache GET requests
            return c.Next()
        }

        cacheKey := c.Path() + "?" + c.Params("id") // Generate a cache key from the request path and query parameters

        // Check if the response is already in the cache
        if cached, found := cache.Get(cacheKey); found {
            return c.JSON(cached)
        }
        err := c.Next()
        if err != nil {
            return err
        }

        var data Post
        cacheKey := c.Path() + "?" + c.Params("id")

        body := c.Response().Body()
        err = json.Unmarshal(body, &data)
        if err != nil {
            return c.JSON(fiber.Map{"error": err.Error()})
        }

        // Cache the response for 10 minutes
        cache.Set(cacheKey, data, 10*time.Minute)

        return nil
    }
}
