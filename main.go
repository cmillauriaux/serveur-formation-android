package main

import (
	"math/rand"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

func main() {
	// Open Database
	open()

	// Initialize API
	router := gin.Default()

	// Register API
	router.POST("/subscribe/:key", newUser)
	router.POST("/connect", connectUser)
	router.GET("/vehicules", listeVehicules)
	router.POST("/vehicules/random", createRandomVehicules)
	router.GET("/vehicules/:vehiculeID", detailsVehicule)
	router.GET("/reservations/:agenceID", listeReservations)
	router.GET("/reservations/:agenceID/:reservationID", detailsReservation)
	router.POST("/reservations/:userID", newReservation)
	router.GET("agences/:agenceID", detailsAgence)
	router.PUT("agences/:agenceID/:userID", updateAgence)
	//router.PUT("/reservations/:agenceID/:reservationID", updateReservation)
	router.POST("/retours/:agenceID/:reservationID", newRetour)

	// Launch server
	router.Run(":80")
}

func newUser(c *gin.Context) {
	key := c.Param("key")
	if key != "1234567890" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "La clé est erronnée"})
		return
	}
	var json User
	if err := c.ShouldBindJSON(&json); err == nil {
		if json.Mail == "" || json.Password == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Le mail et le mot de passe sont obligatoires"})
		} else {
			userID, agenceID, err := newUserFromDB(&json)
			if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			} else {
				c.JSON(http.StatusOK, gin.H{"userID": userID, "agenceID": agenceID})
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
				c.JSON(http.StatusOK, gin.H{"userID": user.UserID, "agenceID": user.AgenceID})
			}
		}
	} else {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}
}

func createRandomVehicules(c *gin.Context) {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 10; i++ {
		vehicule := Vehicule{
			Libelle:      "Vehicule n°" + strconv.FormatInt(int64(r.Intn(1000)), 10),
			NbPlaces:     r.Intn(6) + 1,
			IsDisponible: true,
			LocationMin:  r.Intn(2) + 1,
			LocationMax:  r.Intn(6) + 3,
			TarifMax:     (r.Float32() * 100) + 100,
			TarifMin:     (r.Float32() * 100),
			TarifMoyen:   (r.Float32() * 100) + 50}
		newVehiculeFromDB(&vehicule)
	}
}

func listeReservations(c *gin.Context) {
	agenceID := c.Param("agenceID")
	c.JSON(200, gin.H{"reservations": listReservationsFromDB(agenceID)})
}

func listeVehicules(c *gin.Context) {
	c.JSON(200, gin.H{"vehicules": listVehiculesFromDB()})
}

func detailsReservation(c *gin.Context) {
	agenceID := c.Param("agenceID")
	reservationID := c.Param("reservationID")
	// TODO : Vérifier l'existence de la réservation
	c.JSON(200, gin.H{"reservation": getReservationByAgenceAndID(agenceID, reservationID)})
}

func detailsVehicule(c *gin.Context) {
	vehiculeID := c.Param("vehiculeID")
	// TODO : Vérifier l'existence du véhicule
	c.JSON(200, gin.H{"vehicule": getVehiculeByID(vehiculeID)})
}

func detailsAgence(c *gin.Context) {
	agenceID := c.Param("agenceID")
	// TODO : Vérifier l'existence du véhicule
	c.JSON(200, gin.H{"agence": getAgenceByID(agenceID)})
}

func newReservation(c *gin.Context) {
	userID := c.Param("userID")
	var json Reservation
	if err := c.ShouldBindJSON(&json); err == nil {
		reservation, err := newReservationFromDB(userID, &json)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		} else {
			c.JSON(http.StatusOK, gin.H{"reservation": reservation})
		}
	} else {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}
}

func newRetour(c *gin.Context) {
	agenceID := c.Param("agenceID")
	reservationID := c.Param("reservationID")
	var json Retour
	if err := c.ShouldBindJSON(&json); err == nil {
		retour, err := newRetourFromDB(agenceID, reservationID, &json)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		} else {
			c.JSON(http.StatusOK, gin.H{"retour": retour})
		}
	} else {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}
}

func updateAgence(c *gin.Context) {
	userID := c.Param("userID")
	agenceID := c.Param("agenceID")
	var json Agence
	if err := c.ShouldBindJSON(&json); err == nil {
		err := updateAgenceFromDB(userID, agenceID, &json)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err})
		} else {
			c.JSON(http.StatusOK, gin.H{"result": "OK"})
		}
	} else {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}
}
