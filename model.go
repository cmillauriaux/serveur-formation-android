package main

type Wish struct {
	WishID      string  `form:"wishID" json:"wishID" gorm:"primary_key"`
	UserID      string  `form:"userID" json:"userID"`
	Title       string  `form:"title" json:"title" binding:"required"`
	Description string  `form:"description" json:"description" binding:"required"`
	Price       float32 `form:"price" json:"price" binding:"required"`
	Rating      float32 `form:"rating" json:"rating" binding:"required"`
	Link        string  `form:"link" json:"link" binding:"required"`
}

type User struct {
	UserID   string `form:"userID" json:"userID" gorm:"primary_key"`
	Mail     string `form:"mail" json:"mail" binding:"required"`
	Password string `form:"password" json:"password" binding:"required"`
}
