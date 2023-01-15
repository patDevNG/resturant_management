package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)


func EntryHandler() gin.HandlerFunc{
return func(c *gin.Context) {
	msg :="testing..."
	c.JSON(http.StatusOK, gin.H{"data": msg})
}
}