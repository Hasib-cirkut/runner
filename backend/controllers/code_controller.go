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
func GetSupportedLanguages(c *gin.Context) {
	languages := []string{
		"javascript",
		"python",
		"go",
		"java",
		"c",
		"cpp",
	}

	c.JSON(http.StatusOK, gin.H{"languages": languages})
}
