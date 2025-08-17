package dto

import (
	"net/url"
	"strings"

	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
)

// RegisterCustomValidations registers custom validation functions
func RegisterCustomValidations() {
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		v.RegisterValidation("image_link", validateImageLink)
	}
}

// validateImageLink validates that the field is a valid image URL
func validateImageLink(fl validator.FieldLevel) bool {
	link, ok := fl.Field().Interface().(string)
	if !ok {
		return false
	}

	// Check if it's a valid URL
	parsedURL, err := url.Parse(link)
	if err != nil {
		return false
	}

	// Check if it's an HTTP/HTTPS URL
	if parsedURL.Scheme != "http" && parsedURL.Scheme != "https" {
		return false
	}

	// Check if it has a valid image extension
	lowerLink := strings.ToLower(link)
	imageExtensions := []string{".jpg", ".jpeg", ".png", ".gif", ".bmp", ".webp", ".svg"}

	for _, ext := range imageExtensions {
		if strings.HasSuffix(lowerLink, ext) {
			return true
		}
	}

	// Also check if it contains image-related paths (like /images/, /photos/, etc.)
	imagePaths := []string{"/images/", "/photos/", "/img/", "/pics/", "/media/"}
	for _, path := range imagePaths {
		if strings.Contains(lowerLink, path) {
			return true
		}
	}

	return false
}
