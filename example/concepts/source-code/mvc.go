// Structure de projet Go modulaire (Gin + PostgreSQL + Generics + Auth + Validation)

// --- main.go ---
package main

import (
    "yourapp/config"
    "yourapp/handler"
    "yourapp/middleware"

    "github.com/gin-gonic/gin"
)

func main() {
    db := config.InitDB()
    r := gin.Default()

    r.GET("/users", handler.GetUsers(db))
    r.GET("/users/:id", handler.GetUserByID(db))

    auth := r.Group("/")
    auth.Use(middleware.AuthMiddleware())
    auth.POST("/users", handler.CreateUser(db))
    auth.PUT("/users/:id", handler.UpdateUser(db))
    auth.DELETE("/users/:id", handler.DeleteUser(db))

    r.Run(":8080")
}

// --- config/database.go ---
package config

import (
    "log"

    "github.com/jmoiron/sqlx"
    _ "github.com/lib/pq"
)

func InitDB() *sqlx.DB {
    db, err := sqlx.Connect("postgres", "user=postgres password=postgres dbname=testdb sslmode=disable")
    if err != nil {
        log.Fatal(err)
    }
    return db
}

// --- model/user.go ---
package model

type User struct {
    ID   int    `json:"id" db:"id"`
    Name string `json:"name" binding:"required" db:"name"`
}

// --- repository/sql_repository.go ---
package repository

import (
    "fmt"
    "reflect"

    sq "github.com/Masterminds/squirrel"
    "github.com/jmoiron/sqlx"
)

type SQLRepository[T any] struct {
    DB    *sqlx.DB
    Table string
}

func (r *SQLRepository[T]) GetByID(id int) (*T, error) {
    var t T
    query := fmt.Sprintf("SELECT * FROM %s WHERE id = $1", r.Table)
    err := r.DB.Get(&t, query, id)
    return &t, err
}

func (r *SQLRepository[T]) ListPaginated(limit, offset int) ([]T, error) {
    var items []T
    query := fmt.Sprintf("SELECT * FROM %s ORDER BY id LIMIT $1 OFFSET $2", r.Table)
    err := r.DB.Select(&items, query, limit, offset)
    return items, err
}

func (r *SQLRepository[T]) Create(entity *T) error {
    val := reflect.Indirect(reflect.ValueOf(entity))
    typ := val.Type()

    values := map[string]interface{}{}
    for i := 0; i < typ.NumField(); i++ {
        field := typ.Field(i)
        dbTag := field.Tag.Get("db")
        if dbTag == "" || dbTag == "-" || dbTag == "id" {
            continue
        }
        values[dbTag] = val.Field(i).Interface()
    }

    query, args, err := sq.Insert(r.Table).SetMap(values).PlaceholderFormat(sq.Dollar).ToSql()
    if err != nil {
        return err
    }
    _, err = r.DB.Exec(query, args...)
    return err
}

func (r *SQLRepository[T]) Update(id int, entity *T) error {
    val := reflect.Indirect(reflect.ValueOf(entity))
    typ := val.Type()
    values := map[string]interface{}{}
    for i := 0; i < typ.NumField(); i++ {
        field := typ.Field(i)
        dbTag := field.Tag.Get("db")
        if dbTag == "" || dbTag == "-" || dbTag == "id" {
            continue
        }
        values[dbTag] = val.Field(i).Interface()
    }
    query, args, err := sq.Update(r.Table).SetMap(values).Where(sq.Eq{"id": id}).PlaceholderFormat(sq.Dollar).ToSql()
    if err != nil {
        return err
    }
    _, err = r.DB.Exec(query, args...)
    return err
}

func (r *SQLRepository[T]) Delete(id int) error {
    query, args, err := sq.Delete(r.Table).Where(sq.Eq{"id": id}).PlaceholderFormat(sq.Dollar).ToSql()
    if err != nil {
        return err
    }
    _, err = r.DB.Exec(query, args...)
    return err
}

// --- middleware/auth.go ---
package middleware

import (
    "net/http"

    "github.com/gin-gonic/gin"
)

func AuthMiddleware() gin.HandlerFunc {
    return func(c *gin.Context) {
        if c.GetHeader("Authorization") != "Bearer secret-token" {
            c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
            return
        }
        c.Next()
    }
}

// --- handler/user.go ---
package handler

import (
    "net/http"
    "strconv"

    "yourapp/model"
    "yourapp/repository"

    "github.com/gin-gonic/gin"
    "github.com/jmoiron/sqlx"
)

func CreateUser(db *sqlx.DB) gin.HandlerFunc {
    return func(c *gin.Context) {
        var user model.User
        if err := c.ShouldBindJSON(&user); err != nil {
            c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
            return
        }
        repo := repository.SQLRepository[model.User]{DB: db, Table: "users"}
        err := repo.Create(&user)
        if err != nil {
            c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
            return
        }
        c.JSON(http.StatusCreated, user)
    }
}

func GetUsers(db *sqlx.DB) gin.HandlerFunc {
    return func(c *gin.Context) {
        page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
        size, _ := strconv.Atoi(c.DefaultQuery("size", "10"))
        offset := (page - 1) * size
        repo := repository.SQLRepository[model.User]{DB: db, Table: "users"}
        users, err := repo.ListPaginated(size, offset)
        if err != nil {
            c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
            return
        }
        c.JSON(http.StatusOK, users)
    }
}

func GetUserByID(db *sqlx.DB) gin.HandlerFunc {
    return func(c *gin.Context) {
        id, _ := strconv.Atoi(c.Param("id"))
        repo := repository.SQLRepository[model.User]{DB: db, Table: "users"}
        user, err := repo.GetByID(id)
        if err != nil {
            c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
            return
        }
        c.JSON(http.StatusOK, user)
    }
}

func UpdateUser(db *sqlx.DB) gin.HandlerFunc {
    return func(c *gin.Context) {
        id, _ := strconv.Atoi(c.Param("id"))
        var user model.User
        if err := c.ShouldBindJSON(&user); err != nil {
            c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
            return
        }
        repo := repository.SQLRepository[model.User]{DB: db, Table: "users"}
        err := repo.Update(id, &user)
        if err != nil {
            c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
            return
        }
        c.JSON(http.StatusOK, gin.H{"status": "updated"})
    }
}

func DeleteUser(db *sqlx.DB) gin.HandlerFunc {
    return func(c *gin.Context) {
        id, _ := strconv.Atoi(c.Param("id"))
        repo := repository.SQLRepository[model.User]{DB: db, Table: "users"}
        err := repo.Delete(id)
        if err != nil {
            c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
            return
        }
        c.JSON(http.StatusOK, gin.H{"status": "deleted"})
    }
}
