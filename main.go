package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func main() {
	// Open Database
	open()

	// Initialize API
	router := gin.Default()

	// Register API
	router.POST("/subscribe", newUser)
	router.POST("/connect", connectUser)
	router.GET("/liste/:userID", liste)
	router.GET("/details/:userID/:wishID", details)
	router.POST("/details/:userID", newWish)
	router.PUT("/details/:userID/:wishID", update)
	router.DELETE("/details/:userID/:wishID", delete)

	// Launch server
	router.Run(":8080")
}

func newUser(c *gin.Context) {
	var json User
	if err := c.ShouldBindJSON(&json); err == nil {
		if json.Mail == "" || json.Password == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Le mail et le mot de passe sont obligatoires"})
		} else {
			userID, err := newUserFromDB(&json)
			if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			} else {
				c.JSON(http.StatusBadRequest, gin.H{"userID": userID})
			}
		}
	} else {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}
}

func connectUser(c *gin.Context) {
	var json User
	if err := c.ShouldBindJSON(&json); err == nil {
		if json.Mail == "" || json.Password == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Le mail et le mot de passe sont obligatoires"})
		} else {
			user := getUserByMailAndPasswordFromDB(json.Mail, json.Password)
			if user == nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": "Identifiants incorrects"})
			} else {
				c.JSON(http.StatusBadRequest, gin.H{"userID": user.UserID})
			}
		}
	} else {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}
}

func liste(c *gin.Context) {
	userID := c.Param("userID")
	c.JSON(200, listWishesFromDB(userID))
}

func details(c *gin.Context) {
	userID := c.Param("userID")
	wishID := c.Param("wishID")
	c.JSON(200, getWishesByUserAndID(userID, wishID))
}

func newWish(c *gin.Context) {
	userID := c.Param("userID")
	var json Wish
	if err := c.ShouldBindJSON(&json); err == nil {
		if json.Description == "" || json.Link == "" || json.Price == 0 || json.Rating == 0 || json.Title == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Le souhait n'est pas complet"})
		} else {
			wish, err := newWishFromDB(userID, &json)
			if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": err})
			} else {
				c.JSON(http.StatusBadRequest, wish)
			}
		}
	} else {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}
}

func update(c *gin.Context) {
	userID := c.Param("userID")
	wishID := c.Param("wishID")
	var json Wish
	if err := c.ShouldBindJSON(&json); err == nil {
		if json.Description == "" || json.Link == "" || json.Price == 0 || json.Rating == 0 || json.Title == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Le souhait n'est pas complet"})
		} else {
			err := updateWishFromDB(userID, wishID, &json)
			if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": err})
			} else {
				c.Status(200)
			}
		}
	} else {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}
}

func delete(c *gin.Context) {
	userID := c.Param("userID")
	wishID := c.Param("wishID")
	err := deleteWishFromDB(userID, wishID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err})
	} else {
		c.Status(200)
	}
}
