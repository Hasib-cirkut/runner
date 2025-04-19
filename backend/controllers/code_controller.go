package controllers

import (
	"net/http"
	"runner-api/services"
	"runner-api/utils"

	"github.com/gin-gonic/gin"
)

type CodeRunRequest struct {
	Code     string `json:"code" binding:"required"`
	Language string `json:"language" binding:"required"`
}

type CodeRunResponse struct {
	Stdout string `json:"stdout"`
	Stderr string `json:"stderr"`
	Error  string `json:"error,omitempty"`
}

// PingExample godoc
// @Summary ping example
// @Schemes
// @Description do ping // @Tags example
// @Accept json
// @Produce json
// @Success 200 {string} pong
// @Router /ping/ [get]
func HealthCheck(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"reply": "pong"})
}

// RunCode godoc
// @Summary Run code in a container
// @Description Execute code in a specified programming language within a container
// @Tags code
// @Accept json
// @Produce json
// @Param request body CodeRunRequest true "Code execution request"
// @Success 200 {object} CodeRunResponse
// @Failure 400 {object} map[string]string
// @Router /runcode [post]
func RunCode(c *gin.Context) {
	var req CodeRunRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})

		return
	}

	decodedCode, err := utils.Base64Decode(req.Code)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid Base64 encoding: " + err.Error()})

		return
	}

	stdout, stderr, err := services.RunCodeInContainer(decodedCode, req.Language)

	response := CodeRunResponse{
		Stdout: stdout,
		Stderr: stderr,
	}

	if err != nil {
		response.Error = err.Error()
	}

	c.JSON(http.StatusOK, response)
}

// GetSupportedLanguages godoc
// @Summary Get list of supported programming languages
// @Description Returns a list of programming languages that can be executed
// @Tags code
// @Accept json
// @Produce json
// @Success 200 {object} map[string][]string
// @Router /languages [get]
func GetSupportedLanguages(c *gin.Context) {
	languages := []string{
		"javascript",
	}

	c.JSON(http.StatusOK, gin.H{"languages": languages})
}
