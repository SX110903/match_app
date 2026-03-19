package main

import (
    "fmt"
    "github.com/SX110903/match_app/backend/internal/auth"
    "github.com/SX110903/match_app/backend/internal/config"
    "github.com/SX110903/match_app/backend/internal/database"
)

func main() {
    cfg := config.Get()
    db, err := database.NewMySQL(cfg.Database)
    if err != nil {
        panic(err)
    }
    defer db.Close()

    var hash string
    err = db.QueryRowx("SELECT password_hash FROM users WHERE email = 'ana@test.com'").Scan(&hash)
    if err != nil {
        panic(err)
    }
    
    ok, err := auth.VerifyPassword("Password123!", hash)
    fmt.Printf("Password123! match: %v, err: %v\n", ok, err)
    
    ok2, err2 := auth.VerifyPassword("Password123", hash)
    fmt.Printf("Password123 match: %v, err: %v\n", ok2, err2)
}
