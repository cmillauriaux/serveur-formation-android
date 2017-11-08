package main

type Agence struct {
	AgenceID      string `form:"agenceID" json:"agenceID" gorm:"primary_key"`
	RaisonSociale string `form:"raisonSociale" json:"raisonSociale" binding:"required"`
	SIRET         string `form:"siret" json:"siret" binding:"required"`
	Voie          string `form:"voie" json:"voie" binding:"required"`
	CodePostal    string `form:"codePostal" json:"codePostal" binding:"required"`
	Ville         string `form:"ville" json:"ville" binding:"required"`
}

type Reservation struct {
	ReservationID   string  `form:"reservationID" json:"reservationID" gorm:"primary_key"`
	VehiculeID      string  `form:"vehiculeID" json:"vehiculeID" binding:"required"`
	AgenceID        string  `form:"agenceID" json:"agenceID" binding:"required"`
	DateDebut       int     `form:"dateDebut" json:"dateDebut" binding:"required"`
	DateFin         int     `form:"dateFin" json:"dateFin" binding:"required"`
	TarifJournalier float32 `form:"tarifJournalier" json:"tarifJournalier" binding:"required"`
	IsEncours       bool    `form:"isEnCours" json:"isEnCours"`
}

type Vehicule struct {
	VehiculeID   string  `form:"vehiculeID" json:"vehiculeID" gorm:"primary_key"`
	NbPlaces     int     `form:"nbPlaces" json:"nbPlaces" binding:"required"`
	LocationMin  int     `form:"locationMin" json:"locationMin" binding:"required"`
	LocationMax  int     `form:"locationMax" json:"locationMax" binding:"required"`
	TarifMin     float32 `form:"tarifMin" json:"tarifMin" binding:"required"`
	TarifMax     float32 `form:"tarifMax" json:"tarifMax" binding:"required"`
	TarifMoyen   float32 `form:"tarifMoyen" json:"tarifMoyen" binding:"required"`
	IsDisponible bool    `form:"isDisponible" json:"isDisponible" binding:"required"`
}

type User struct {
	UserID   string `form:"userID" json:"userID" gorm:"primary_key"`
	AgenceID string `form:"agenceID" json:"agenceID"`
	Mail     string `form:"mail" json:"mail" binding:"required"`
	Password string `form:"password" json:"password" binding:"required"`
}
