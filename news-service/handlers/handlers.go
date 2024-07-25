package handlers

import (
	"context"
	"net/http"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"news-service/database"
	"news-service/models"
)

func EditNews(c *fiber.Ctx) error {
	id := c.Params("Id")
	newsID, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "Invalid news ID"})
	}

	var news models.News
	if err := c.BodyParser(&news); err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "Invalid input"})
	}

	// Create a map to store the update fields
	updateFields := make(map[string]interface{})
	if news.Title != "" {
		updateFields["Title"] = news.Title
	}
	if news.Content != "" {
		updateFields["Content"] = news.Content
	}

	// Update logic
	if len(updateFields) > 0 {
		setClause := ""
		params := []interface{}{}
		paramID := 1
		for k, v := range updateFields {
			setClause += k + " = $" + strconv.Itoa(paramID) + ", "
			params = append(params, v)
			paramID++
		}
		setClause = setClause[:len(setClause)-2] // Remove the trailing comma and space
		params = append(params, newsID)          // Append the news ID at the end for the WHERE clause

		query := "UPDATE News SET " + setClause + " WHERE Id = $" + strconv.Itoa(paramID)
		_, err := database.DB.Exec(context.Background(), query, params...)
		if err != nil {
			return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to update news"})
		}
	}

	// Handle categories
	if len(news.Categories) > 0 {
		// Delete existing categories
		_, err := database.DB.Exec(context.Background(), "DELETE FROM NewsCategories WHERE NewsId = $1", newsID)
		if err != nil {
			return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to update categories"})
		}

		// Insert new categories
		for _, categoryID := range news.Categories {
			_, err := database.DB.Exec(context.Background(), "INSERT INTO NewsCategories (NewsId, CategoryId) VALUES ($1, $2)", newsID, categoryID)
			if err != nil {
				return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to insert categories"})
			}
		}
	}

	return c.JSON(fiber.Map{"success": true})
}

func ListNews(c *fiber.Ctx) error {
	rows, err := database.DB.Query(context.Background(), "SELECT Id, Title, Content FROM News")
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to fetch news"})
	}
	defer rows.Close()

	newsList := []models.News{}
	for rows.Next() {
		var news models.News
		if err := rows.Scan(&news.ID, &news.Title, &news.Content); err != nil {
			return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to scan news"})
		}

		// Fetch categories for each news item
		categoryRows, err := database.DB.Query(context.Background(), "SELECT CategoryId FROM NewsCategories WHERE NewsId = $1", news.ID)
		if err != nil {
			return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to fetch categories"})
		}
		defer categoryRows.Close()

		categories := []int64{}
		for categoryRows.Next() {
			var categoryID int64
			if err := categoryRows.Scan(&categoryID); err != nil {
				return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to scan categories"})
			}
			categories = append(categories, categoryID)
		}

		news.Categories = categories
		newsList = append(newsList, news)
	}

	return c.JSON(fiber.Map{"success": true, "news": newsList})
}
