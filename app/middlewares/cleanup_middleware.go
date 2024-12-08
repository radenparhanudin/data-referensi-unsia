package middlewares

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/gofiber/fiber/v2"
)

func CleanupMiddleware() fiber.Handler {
	folders := []string{"tmp/uploads", "tmp/exports"}
	maxAge := 5 * time.Minute

	return func(c *fiber.Ctx) error {
		// Ensure required folders exist
		for _, folder := range folders {
			if _, err := os.Stat(folder); os.IsNotExist(err) {
				if err := os.MkdirAll(folder, os.ModePerm); err != nil {
					return fmt.Errorf("failed to create folder %s: %v", folder, err)
				}
				fmt.Printf("Folder created: %s\n", folder)
			}
		}

		// Clean up old files in the folders
		go func() {
			for _, folder := range folders {
				err := filepath.Walk(folder, func(path string, info os.FileInfo, err error) error {
					if err != nil {
						return err
					}

					// Skip if it's a directory
					if info.IsDir() {
						return nil
					}

					// Check if the file is older than maxAge
					if time.Since(info.ModTime()) > maxAge {
						if err := os.Remove(path); err != nil {
							fmt.Printf("Failed to delete file: %s, error: %v\n", path, err)
						} else {
							fmt.Printf("Deleted file: %s\n", path)
						}
					}
					return nil
				})

				if err != nil {
					fmt.Printf("Error cleaning up folder %s: %v\n", folder, err)
				}
			}
		}()

		return c.Next()
	}
}
