package model

import (
	"github.com/dgrijalva/jwt-go"
	"time"
)

type User struct {
	Username string `form:"username" json:"username" binding:"required"`
	Password string `form:"password" json:"password" binding:"required"`
}
type Getquestion struct {
	Id       int    `json:"id"`
	Username string `json:"username"`
	Question string `json:"question"`
}
type Getanswer struct {
	Id         int    `json:"id"`
	Questionid string `json:"questionid"`
	Username   string `json:"username"`
	Answer     string `json:"answer"`
}
type Question struct {
	question string `form:"question" json:"question" binding:"required"`
}
type Answer struct {
	question string `form:"question" json:"question" binding:"required"`
	answer   string `form:"answer" json:"answer" binding:"required"`
}
type Message struct {
	Id        string `form:"Id" json:"id" binding:"required"`
	Text      string `form:"Text" json:"Text" binding:"required"`
	CreatedAt time.Time
}
type MyClaims struct {
	Username string `json:"username"`
	jwt.StandardClaims
}
