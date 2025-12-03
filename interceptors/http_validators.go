package interceptors

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// HttpValidateRequest is a helper function to validate request body in handlers
// Usage in handler:
//
//	var req MyRequest
//	if err := interceptors.HttpValidateRequest(c, &req); err != nil {
//	    return // error response already sent
//	}
func HttpValidateRequest(c *gin.Context, req interface{}) error {
	// Bind JSON body to struct
	if err := c.ShouldBindJSON(req); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid request body",
			"details": err.Error(),
		})
		return err
	}

	// Validate using custom validator if implemented
	if validator, ok := req.(IValidator); ok {
		if err := validator.Validate(); err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
				"error":   "Validation failed",
				"details": err.Error(),
			})
			return err
		}
	}

	return nil
}

// HttpValidateQuery is a helper function to validate query parameters
// Usage in handler:
//
//	var req MyQueryRequest
//	if err := interceptors.HttpValidateQuery(c, &req); err != nil {
//	    return // error response already sent
//	}
func HttpValidateQuery(c *gin.Context, req interface{}) error {
	// Bind query parameters to struct
	if err := c.ShouldBindQuery(req); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid query parameters",
			"details": err.Error(),
		})
		return err
	}

	// Validate using custom validator if implemented
	if validator, ok := req.(IValidator); ok {
		if err := validator.Validate(); err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
				"error":   "Validation failed",
				"details": err.Error(),
			})
			return err
		}
	}

	return nil
}

// HttpValidateUri is a helper function to validate URI parameters
// Usage in handler:
//
//	var req MyUriRequest
//	if err := interceptors.HttpValidateUri(c, &req); err != nil {
//	    return // error response already sent
//	}
func HttpValidateUri(c *gin.Context, req interface{}) error {
	// Bind URI parameters to struct
	if err := c.ShouldBindUri(req); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid URI parameters",
			"details": err.Error(),
		})
		return err
	}

	// Validate using custom validator if implemented
	if validator, ok := req.(IValidator); ok {
		if err := validator.Validate(); err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
				"error":   "Validation failed",
				"details": err.Error(),
			})
			return err
		}
	}

	return nil
}
