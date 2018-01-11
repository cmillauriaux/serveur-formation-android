package main

import (
	"errors"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	"github.com/satori/go.uuid"
)

var db *gorm.DB

func open() {
	var err error
	db, err = gorm.Open("sqlite3", "reservation.db")
	if err != nil {
		panic("failed to connect database")
	}
	db.AutoMigrate(&Agence{})
	db.AutoMigrate(&User{})
	db.AutoMigrate(&Reservation{})
	db.AutoMigrate(&Vehicule{})
}

func newUserFromDB(user *User) (string, string, error) {
	if getUserByMailFromDB(user.Mail) != nil {
		return "", "", errors.New("Le compte existe déjà")
	}

	id, _ := uuid.NewV4()

	user.UserID = id.String()
	id, _ = uuid.NewV4()
	user.AgenceID = id.String()
	db.Create(user)

	agence := Agence{AgenceID: user.AgenceID}
	db.Create(agence)

	return user.UserID, user.AgenceID, nil
}

func newVehiculeFromDB(vehicule *Vehicule) error {
	id, _ := uuid.NewV4()
	vehicule.VehiculeID = id.String()
	db.Create(vehicule)

	return nil
}

func getUserFromDB(userID string) *User {
	user := User{}
	db.Where(&User{UserID: userID}).First(&user)
	if user.UserID == "" {
		return nil
	}
	return &user
}

func getUserByMailFromDB(mail string) *User {
	user := User{}
	db.Where(&User{Mail: mail}).First(&user)
	if user.UserID == "" {
		return nil
	}
	return &user
}

func getUserByMailAndPasswordFromDB(mail string, password string) *User {
	user := User{}
	db.Where(&User{Mail: mail, Password: password}).First(&user)
	if user.UserID == "" {
		return nil
	}
	return &user
}

func listReservationsFromDB(agenceID string) []Reservation {
	reservations := make([]Reservation, 0)
	db.Where(&Reservation{AgenceID: agenceID}).Find(&reservations)
	return reservations
}

func listVehiculesFromDB() []Vehicule {
	vehicules := make([]Vehicule, 0)
	db.Find(&vehicules)
	return vehicules
}

func newReservationFromDB(userID string, reservation *Reservation) (*Reservation, error) {
	user := getUserFromDB(userID)
	if user == nil {
		return nil, errors.New("L'utilisateur n'existe pas")
	}
	vehicule := getVehiculeByID(reservation.VehiculeID)
	if vehicule == nil {
		return nil, errors.New("Le véhicule n'existe pas")
	}
	if !vehicule.IsDisponible {
		return nil, errors.New("Le véhicule n'est pas disponible")
	}
	reservation.AgenceID = user.AgenceID
	id, _ := uuid.NewV4()
	reservation.ReservationID = id.String()
	reservation.IsEncours = true
	db.Create(reservation)

	vehicule.IsDisponible = false
	db.Save(vehicule)

	return reservation, nil
}

func newRetourFromDB(agenceID string, reservationID string, retour *Retour) (*Retour, error) {
	agence := getAgenceByID(agenceID)
	if agence == nil {
		return nil, errors.New("L'agence n'existe pas")
	}
	reservation := getReservationByAgenceAndID(agenceID, reservationID)
	if reservation == nil {
		return nil, errors.New("La réservation n'existe pas ou ne correspond pas à la bonne agence")
	}

	if !reservation.IsEncours {
		return nil, errors.New("La réservation est déjà terminée")
	}

	vehicule := getVehiculeByID(reservation.VehiculeID)
	if vehicule == nil {
		return nil, errors.New("Le véhicule n'existe pas")
	}

	id, _ := uuid.NewV4()
	retour.RetourID = id.String()
	retour.ReservationID = reservationID
	reservation.IsEncours = true
	db.Create(retour)

	reservation.IsEncours = false
	db.Save(reservation)

	vehicule.IsDisponible = true
	db.Save(vehicule)

	return retour, nil
}

func getVehiculeByID(vehiculeID string) *Vehicule {
	vehicule := Vehicule{}
	db.Where(&Vehicule{VehiculeID: vehiculeID}).First(&vehicule)
	if vehicule.VehiculeID == "" {
		return nil
	}
	return &vehicule
}

func getAgenceByID(agenceID string) *Agence {
	agence := Agence{}
	db.Where(&Agence{AgenceID: agenceID}).First(&agence)
	if agence.AgenceID == "" {
		return nil
	}
	return &agence
}

func getReservationByAgenceAndID(agenceID string, reservationID string) *Reservation {
	reservation := Reservation{}
	db.Where(&Reservation{AgenceID: agenceID, ReservationID: reservationID}).First(&reservation)
	if reservation.ReservationID == "" {
		return nil
	}
	return &reservation
}

func updateReservationFromDB(agenceID string, reservationID string, reservation *Reservation) error {
	reservationDb := getReservationByAgenceAndID(agenceID, reservationID)
	if reservationDb == nil {
		return errors.New("La réservation n'existe pas")
	}
	reservationDb.DateDebut = reservation.DateDebut
	reservationDb.DateFin = reservation.DateFin
	reservationDb.IsEncours = reservation.IsEncours
	reservationDb.TarifJournalier = reservation.TarifJournalier
	db.Save(reservationDb)
	return nil
}

func updateAgenceFromDB(userID string, agenceID string, agence *Agence) error {
	user := getUserFromDB(userID)
	if user == nil {
		return errors.New("L'utilisateur n'existe pas")
	}

	agenceDb := getAgenceByID(agenceID)
	if agenceDb == nil {
		return errors.New("L'agence n'existe pas")
	}

	if user.AgenceID != agenceDb.AgenceID {
		return errors.New("L'utilisateur n'appartient pas à l'agence")
	}

	agenceDb.CodePostal = agence.CodePostal
	agenceDb.RaisonSociale = agence.RaisonSociale
	agenceDb.SIRET = agence.SIRET
	agenceDb.Ville = agence.Ville
	agenceDb.Voie = agence.Voie
	db.Save(agenceDb)
	return nil
}
