//go:build ignore

package main

import (
	"crypto/rand"
	"database/sql"
	"encoding/base64"
	"fmt"
	"log"

	_ "github.com/go-sql-driver/mysql"
	"golang.org/x/crypto/argon2"
)

const dsn = "matchhub_user:devpassword@tcp(127.0.0.1:3307)/matchhub?parseTime=true"

func hash(pw string) string {
	salt := make([]byte, 16)
	rand.Read(salt)
	h := argon2.IDKey([]byte(pw), salt, 3, 64*1024, 4, 32)
	return fmt.Sprintf("$argon2id$v=19$m=65536,t=3,p=4$%s$%s",
		base64.RawStdEncoding.EncodeToString(salt),
		base64.RawStdEncoding.EncodeToString(h))
}

func main() {
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	if err := db.Ping(); err != nil {
		log.Fatal(err)
	}

	pw := hash("Test1234!")

	type seedU struct {
		id, email, name string
		age             int
		photo           string
	}

	users := []seedU{
		{"idddddd0-id00-0000-0000-00000000a1ce", "u_alice@test.com", "Alice", 24, "https://i.pravatar.cc/400?u=alice"},
		{"idddddd0-id00-0000-0000-00000000b4un", "u_bruno@test.com", "Bruno", 28, "https://i.pravatar.cc/400?u=bruno"},
		{"idddddd0-id00-0000-0000-0000000ca41a", "u_carla@test.com", "Carla", 26, "https://i.pravatar.cc/400?u=carla"},
		{"idddddd0-id00-0000-0000-000000dav1d0", "u_david@test.com", "David", 30, "https://i.pravatar.cc/400?u=david"},
		{"idddddd0-id00-0000-0000-00000e1ena00", "u_elena@test.com", "Elena", 22, "https://i.pravatar.cc/400?u=elena"},
		{"idddddd0-id00-0000-0000-000000f4an00", "u_fran@test.com", "Fran", 32, "https://i.pravatar.cc/400?u=fran"},
	}

	for _, u := range users {
		var n int
		db.QueryRow("SELECT COUNT(*) FROM users WHERE id=?", u.id).Scan(&n)
		if n > 0 {
			fmt.Printf("skip: %s\n", u.email)
			continue
		}

		db.Exec(`INSERT INTO users (id,email,password_hash,email_verified_at,is_admin,is_frozen,vip_level,credits)
			VALUES (?,?,?,NOW(),FALSE,FALSE,0,100)`, u.id, u.email, pw)

		profID := "pr-" + u.id[4:]
		db.Exec(`INSERT INTO user_profiles (id,user_id,name,age,bio,occupation,location,latitude,longitude)
			VALUES (?,?,?,?,'Identity test','Tester','Madrid',40.4168,-3.7038)`,
			profID, u.id, u.name, u.age)

		db.Exec(`INSERT INTO user_preferences (id,user_id,min_age,max_age,max_distance_km,interested_in)
			VALUES (?,?,18,99,9999,'both')`, "pf-"+u.id[4:], u.id)

		db.Exec(`INSERT INTO user_photos (id,user_id,url,sort_order) VALUES (?,?,?,0)`,
			"ph-"+u.id[4:], u.id, u.photo)

		fmt.Printf("  ✓ %s (%s) id=%s\n", u.email, u.name, u.id)
	}

	fmt.Println("\n-- verify --")
	rows, _ := db.Query("SELECT id, email FROM users WHERE email LIKE 'u%@test.com' ORDER BY email")
	defer rows.Close()
	for rows.Next() {
		var id, em string
		rows.Scan(&id, &em)
		fmt.Printf("  %s  %s\n", id, em)
	}
}
