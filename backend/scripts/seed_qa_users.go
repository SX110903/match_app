//go:build ignore

package main

import (
	"crypto/rand"
	"database/sql"
	"encoding/base64"
	"fmt"
	"log"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"golang.org/x/crypto/argon2"
)

const dsn = "matchhub_user:devpassword@tcp(127.0.0.1:3307)/matchhub?parseTime=true"

func hashPassword(password string) string {
	salt := make([]byte, 16)
	rand.Read(salt)
	hash := argon2.IDKey([]byte(password), salt, 3, 64*1024, 4, 32)
	return fmt.Sprintf("$argon2id$v=19$m=65536,t=3,p=4$%s$%s",
		base64.RawStdEncoding.EncodeToString(salt),
		base64.RawStdEncoding.EncodeToString(hash),
	)
}

type qaUser struct {
	id           string
	email        string
	name         string
	age          int
	bio          string
	occupation   string
	location     string
	lat, lon     float64
	photoURLs    []string
	interests    []string
	interestedIn string
	group        string // A, B, C, D
}

var qaUsers = []qaUser{
	// ── GROUP A: complete profile + photos ──────────────────────────────────
	{
		id: "22222222-aaaa-0001-0000-000000000001", group: "A",
		email: "qa_a1@test.com", name: "Ana García", age: 25,
		bio:          "Profesora de yoga y amante de los atardeceres.",
		occupation:   "Instructora de Yoga", location: "Madrid, España",
		lat: 40.4168, lon: -3.7038,
		photoURLs:    []string{"https://i.pravatar.cc/300?u=qa_a1_1", "https://i.pravatar.cc/300?u=qa_a1_2"},
		interests:    []string{"yoga", "viajes", "meditación"},
		interestedIn: "both",
	},
	{
		id: "22222222-aaaa-0001-0000-000000000002", group: "A",
		email: "qa_a2@test.com", name: "Bruno Díaz", age: 30,
		bio:          "Fotógrafo freelance. Cada foto cuenta una historia.",
		occupation:   "Fotógrafo", location: "Barcelona, España",
		lat: 41.3874, lon: 2.1686,
		photoURLs:    []string{"https://i.pravatar.cc/300?u=qa_a2_1"},
		interests:    []string{"fotografía", "viajes", "arte"},
		interestedIn: "both",
	},
	{
		id: "22222222-aaaa-0001-0000-000000000003", group: "A",
		email: "qa_a3@test.com", name: "Carla Ruiz", age: 27,
		bio:          "Chef repostera. Creo que el amor entra por el estómago.",
		occupation:   "Repostera", location: "Valencia, España",
		lat: 39.4699, lon: -0.3763,
		photoURLs:    []string{"https://i.pravatar.cc/300?u=qa_a3_1", "https://i.pravatar.cc/300?u=qa_a3_2"},
		interests:    []string{"cocina", "repostería", "música"},
		interestedIn: "male",
	},
	{
		id: "22222222-aaaa-0001-0000-000000000004", group: "A",
		email: "qa_a4@test.com", name: "David Mora", age: 32,
		bio:          "Músico y productor. Busco mi musa.",
		occupation:   "Productor Musical", location: "Sevilla, España",
		lat: 37.3882, lon: -5.9823,
		photoURLs:    []string{"https://i.pravatar.cc/300?u=qa_a4_1", "https://i.pravatar.cc/300?u=qa_a4_2"},
		interests:    []string{"música", "cine", "literatura"},
		interestedIn: "female",
	},
	{
		id: "22222222-aaaa-0001-0000-000000000005", group: "A",
		email: "qa_a5@test.com", name: "Elena Vega", age: 29,
		bio:          "Abogada de día, bailarina de noche.",
		occupation:   "Abogada", location: "Bilbao, España",
		lat: 43.2627, lon: -2.9253,
		photoURLs:    []string{"https://i.pravatar.cc/300?u=qa_a5_1"},
		interests:    []string{"derecho", "baile", "viajes"},
		interestedIn: "both",
	},

	// ── GROUP B: complete profile, no photos ────────────────────────────────
	{
		id: "22222222-bbbb-0001-0000-000000000011", group: "B",
		email: "qa_b1@test.com", name: "Felipe Ortega", age: 28,
		bio:          "Desarrollador backend. Fan del código limpio y el café.",
		occupation:   "Backend Developer", location: "Madrid, España",
		lat: 40.4168, lon: -3.7038,
		photoURLs:    nil,
		interests:    []string{"tecnología", "café", "ajedrez"},
		interestedIn: "both",
	},
	{
		id: "22222222-bbbb-0001-0000-000000000012", group: "B",
		email: "qa_b2@test.com", name: "Gloria Méndez", age: 26,
		bio:          "Bióloga marina. El océano es mi hogar.",
		occupation:   "Bióloga Marina", location: "Málaga, España",
		lat: 36.7213, lon: -4.4217,
		photoURLs:    nil,
		interests:    []string{"biología", "buceo", "naturaleza"},
		interestedIn: "both",
	},
	{
		id: "22222222-bbbb-0001-0000-000000000013", group: "B",
		email: "qa_b3@test.com", name: "Hugo Serrano", age: 35,
		bio:          "Empresario. Trabajo duro y juego duro.",
		occupation:   "Empresario", location: "Zaragoza, España",
		lat: 41.6488, lon: -0.8891,
		photoURLs:    nil,
		interests:    []string{"negocios", "golf", "tecnología"},
		interestedIn: "female",
	},
	{
		id: "22222222-bbbb-0001-0000-000000000014", group: "B",
		email: "qa_b4@test.com", name: "Irene Blanco", age: 23,
		bio:          "Estudiante de arquitectura. Amo los espacios bien diseñados.",
		occupation:   "Estudiante de Arquitectura", location: "Granada, España",
		lat: 37.1773, lon: -3.5986,
		photoURLs:    nil,
		interests:    []string{"arquitectura", "arte", "diseño"},
		interestedIn: "both",
	},
	{
		id: "22222222-bbbb-0001-0000-000000000015", group: "B",
		email: "qa_b5@test.com", name: "Jorge Navarro", age: 31,
		bio:          "Montañero y aventurero. Siempre buscando la próxima cumbre.",
		occupation:   "Guía de Montaña", location: "Pamplona, España",
		lat: 42.8169, lon: -1.6438,
		photoURLs:    nil,
		interests:    []string{"montaña", "escalada", "senderismo"},
		interestedIn: "both",
	},

	// ── GROUP C: incomplete profile (no bio, no occupation, no interests) ───
	{
		id: "22222222-cccc-0001-0000-000000000021", group: "C",
		email: "qa_c1@test.com", name: "Karen Pons", age: 24,
		bio: "", occupation: "", location: "Madrid, España",
		lat: 40.4168, lon: -3.7038,
		photoURLs:    []string{"https://i.pravatar.cc/300?u=qa_c1_1"},
		interests:    nil,
		interestedIn: "both",
	},
	{
		id: "22222222-cccc-0001-0000-000000000022", group: "C",
		email: "qa_c2@test.com", name: "Luis Romero", age: 29,
		bio: "", occupation: "", location: "",
		lat: 0, lon: 0,
		photoURLs:    nil,
		interests:    nil,
		interestedIn: "both",
	},
	{
		id: "22222222-cccc-0001-0000-000000000023", group: "C",
		email: "qa_c3@test.com", name: "María Santos", age: 33,
		bio: "Hola!", occupation: "", location: "Barcelona, España",
		lat: 41.3874, lon: 2.1686,
		photoURLs:    nil,
		interests:    nil,
		interestedIn: "male",
	},
	{
		id: "22222222-cccc-0001-0000-000000000024", group: "C",
		email: "qa_c4@test.com", name: "Nicolás Reyes", age: 26,
		bio: "", occupation: "Algo", location: "",
		lat: 0, lon: 0,
		photoURLs:    nil,
		interests:    []string{"música"},
		interestedIn: "both",
	},
	{
		id: "22222222-cccc-0001-0000-000000000025", group: "C",
		email: "qa_c5@test.com", name: "Olivia Cano", age: 21,
		bio: "", occupation: "", location: "Sevilla, España",
		lat: 37.3882, lon: -5.9823,
		photoURLs:    nil,
		interests:    nil,
		interestedIn: "both",
	},

	// ── GROUP D: new accounts, no profile activity ───────────────────────────
	{
		id: "22222222-dddd-0001-0000-000000000031", group: "D",
		email: "qa_d1@test.com", name: "Pedro Alonso", age: 27,
		bio: "", occupation: "", location: "",
		lat: 0, lon: 0, photoURLs: nil, interests: nil, interestedIn: "both",
	},
	{
		id: "22222222-dddd-0001-0000-000000000032", group: "D",
		email: "qa_d2@test.com", name: "Queta Sanz", age: 22,
		bio: "", occupation: "", location: "",
		lat: 0, lon: 0, photoURLs: nil, interests: nil, interestedIn: "both",
	},
	{
		id: "22222222-dddd-0001-0000-000000000033", group: "D",
		email: "qa_d3@test.com", name: "Roberto Gil", age: 34,
		bio: "", occupation: "", location: "",
		lat: 0, lon: 0, photoURLs: nil, interests: nil, interestedIn: "female",
	},
	{
		id: "22222222-dddd-0001-0000-000000000034", group: "D",
		email: "qa_d4@test.com", name: "Sara Campos", age: 28,
		bio: "", occupation: "", location: "",
		lat: 0, lon: 0, photoURLs: nil, interests: nil, interestedIn: "male",
	},
	{
		id: "22222222-dddd-0001-0000-000000000035", group: "D",
		email: "qa_d5@test.com", name: "Tomás Rubio", age: 30,
		bio: "", occupation: "", location: "",
		lat: 0, lon: 0, photoURLs: nil, interests: nil, interestedIn: "both",
	},
}

func main() {
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		log.Fatal("open:", err)
	}
	defer db.Close()
	if err := db.Ping(); err != nil {
		log.Fatal("ping:", err)
	}
	fmt.Println("Connected to DB")

	password := hashPassword("Test1234!")
	now := time.Now()

	for _, u := range qaUsers {
		var exists int
		db.QueryRow("SELECT COUNT(*) FROM users WHERE id = ?", u.id).Scan(&exists)
		if exists > 0 {
			fmt.Printf("  skip (exists): %s\n", u.email)
			continue
		}

		// 1. Insert user
		_, err := db.Exec(`
			INSERT INTO users (id, email, password_hash, email_verified_at, is_admin, is_frozen, vip_level, credits)
			VALUES (?, ?, ?, ?, FALSE, FALSE, 0, 100)`,
			u.id, u.email, password, now,
		)
		if err != nil {
			log.Printf("ERR insert user %s: %v", u.email, err)
			continue
		}

		// 2. Insert profile (always, even if sparse)
		profileID := u.id[:8] + "-prof-0000-0000-" + u.id[len(u.id)-12:]
		var bioVal, occVal, locVal interface{}
		if u.bio != "" {
			bioVal = u.bio
		}
		if u.occupation != "" {
			occVal = u.occupation
		}
		if u.location != "" {
			locVal = u.location
		}
		var latVal, lonVal interface{}
		if u.lat != 0 || u.lon != 0 {
			latVal = u.lat
			lonVal = u.lon
		}
		_, err = db.Exec(`
			INSERT INTO user_profiles (id, user_id, name, age, bio, occupation, location, latitude, longitude)
			VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)`,
			profileID, u.id, u.name, u.age, bioVal, occVal, locVal, latVal, lonVal,
		)
		if err != nil {
			log.Printf("ERR insert profile %s: %v", u.email, err)
		}

		// 3. Insert preferences
		prefsID := u.id[:8] + "-pref-0000-0000-" + u.id[len(u.id)-12:]
		_, err = db.Exec(`
			INSERT INTO user_preferences (id, user_id, min_age, max_age, max_distance_km, interested_in)
			VALUES (?, ?, 18, 45, 500, ?)`,
			prefsID, u.id, u.interestedIn,
		)
		if err != nil {
			log.Printf("ERR insert prefs %s: %v", u.email, err)
		}

		// 4. Insert photos
		for i, url := range u.photoURLs {
			photoID := fmt.Sprintf("%s-ph%02d-0000-0000-%s", u.id[:8], i, u.id[len(u.id)-12:])
			db.Exec(`INSERT INTO user_photos (id, user_id, url, sort_order) VALUES (?, ?, ?, ?)`,
				photoID, u.id, url, i,
			)
		}

		// 5. Insert interests
		for i, interest := range u.interests {
			intID := fmt.Sprintf("%s-in%02d-0000-0000-%s", u.id[:8], i, u.id[len(u.id)-12:])
			db.Exec(`INSERT IGNORE INTO user_interests (id, user_id, interest) VALUES (?, ?, ?)`,
				intID, u.id, interest,
			)
		}

		fmt.Printf("  ✓ [%s] %s (%s) age=%d photos=%d interests=%d\n",
			u.group, u.email, u.name, u.age, len(u.photoURLs), len(u.interests))
	}

	// Verification
	fmt.Println("\n--- Verification ---")
	var totalUsers, totalQA int
	db.QueryRow("SELECT COUNT(*) FROM users WHERE deleted_at IS NULL").Scan(&totalUsers)
	db.QueryRow("SELECT COUNT(*) FROM users WHERE id LIKE '22222222%'").Scan(&totalQA)
	fmt.Printf("total_users=%d  qa_users=%d\n", totalUsers, totalQA)

	rows, _ := db.Query("SELECT id, email FROM users WHERE id LIKE '22222222%' ORDER BY id")
	defer rows.Close()
	for rows.Next() {
		var id, email string
		rows.Scan(&id, &email)
		fmt.Printf("  seeded: %s  %s\n", id, email)
	}
}
