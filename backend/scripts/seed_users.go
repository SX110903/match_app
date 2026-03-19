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

func hashPassword(password string) string {
	salt := make([]byte, 16)
	rand.Read(salt)
	hash := argon2.IDKey([]byte(password), salt, 3, 64*1024, 4, 32)
	return fmt.Sprintf("$argon2id$v=19$m=65536,t=3,p=4$%s$%s",
		base64.RawStdEncoding.EncodeToString(salt),
		base64.RawStdEncoding.EncodeToString(hash),
	)
}

type seedUser struct {
	id         string
	email      string
	name       string
	age        int
	bio        string
	occupation string
	location   string
	lat, lon   float64
	photoURLs  []string // empty = no photos (edge case)
	interests  []string
	interestedIn string
}

var users = []seedUser{
	{
		id: "11111111-1111-1111-1111-111111111001",
		email: "sofia@test.com", name: "Sofía Ramírez", age: 24,
		bio: "Amante del café y los viajes. Buscando a alguien con quien explorar el mundo.",
		occupation: "Diseñadora UX", location: "Madrid, España",
		lat: 40.4168, lon: -3.7038,
		photoURLs:    []string{"https://i.pravatar.cc/300?u=sofia001", "https://i.pravatar.cc/300?u=sofia002"},
		interests:    []string{"viajes", "diseño", "yoga", "fotografía"},
		interestedIn: "both",
	},
	{
		id: "11111111-1111-1111-1111-111111111002",
		email: "miguel@test.com", name: "Miguel Torres", age: 28,
		bio: "Ingeniero de día, guitarrista de noche. Me encanta el senderismo.",
		occupation: "Ingeniero de Software", location: "Barcelona, España",
		lat: 41.3874, lon: 2.1686,
		photoURLs:    []string{"https://i.pravatar.cc/300?u=miguel001", "https://i.pravatar.cc/300?u=miguel002"},
		interests:    []string{"música", "senderismo", "tecnología", "cocina"},
		interestedIn: "both",
	},
	{
		id: "11111111-1111-1111-1111-111111111003",
		email: "valentina@test.com", name: "Valentina López", age: 26,
		bio: "Médica residente con pasión por los animales y la lectura.",
		occupation: "Médica", location: "Sevilla, España",
		lat: 37.3882, lon: -5.9823,
		photoURLs:    []string{"https://i.pravatar.cc/300?u=valentina001"},
		interests:    []string{"medicina", "libros", "animales", "running"},
		interestedIn: "male",
	},
	{
		id: "11111111-1111-1111-1111-111111111004",
		email: "diego@test.com", name: "Diego Fernández", age: 30,
		bio: "Chef profesional. Creo que la mejor cita empieza con buena comida.",
		occupation: "Chef", location: "Valencia, España",
		lat: 39.4699, lon: -0.3763,
		photoURLs:    []string{"https://i.pravatar.cc/300?u=diego001", "https://i.pravatar.cc/300?u=diego002", "https://i.pravatar.cc/300?u=diego003"},
		interests:    []string{"gastronomía", "vinos", "viajes", "cine"},
		interestedIn: "female",
	},
	{
		id: "11111111-1111-1111-1111-111111111005",
		email: "isabella@test.com", name: "Isabella Martín", age: 22,
		bio: "Estudiante de Bellas Artes. Me gusta pintar y ver series hasta tarde.",
		occupation: "Estudiante", location: "Bilbao, España",
		lat: 43.2627, lon: -2.9253,
		photoURLs:    []string{"https://i.pravatar.cc/300?u=isabella001", "https://i.pravatar.cc/300?u=isabella002"},
		interests:    []string{"arte", "series", "música", "café"},
		interestedIn: "both",
	},
	{
		id: "11111111-1111-1111-1111-111111111006",
		email: "andres@test.com", name: "Andrés Gómez", age: 32,
		bio: "Periodista deportivo. El deporte es mi vida, la escritura mi pasión.",
		occupation: "Periodista", location: "Málaga, España",
		lat: 36.7213, lon: -4.4217,
		photoURLs:    []string{"https://i.pravatar.cc/300?u=andres001"},
		interests:    []string{"fútbol", "periodismo", "ciclismo", "playa"},
		interestedIn: "female",
	},
	// Users 7 & 8 will have a pre-seeded mutual match
	{
		id: "11111111-1111-1111-1111-111111111007",
		email: "camila@test.com", name: "Camila Rodríguez", age: 27,
		bio: "Psicóloga. Me encanta escuchar historias y tomar vino tinto.",
		occupation: "Psicóloga", location: "Zaragoza, España",
		lat: 41.6488, lon: -0.8891,
		photoURLs:    []string{"https://i.pravatar.cc/300?u=camila001", "https://i.pravatar.cc/300?u=camila002"},
		interests:    []string{"psicología", "vino", "meditación", "literatura"},
		interestedIn: "both",
	},
	{
		id: "11111111-1111-1111-1111-111111111008",
		email: "javier@test.com", name: "Javier Sánchez", age: 29,
		bio: "Arquitecto y amante del diseño urbano. Busco a alguien que comparta mis paseos.",
		occupation: "Arquitecto", location: "Zaragoza, España",
		lat: 41.6500, lon: -0.8900,
		photoURLs:    []string{"https://i.pravatar.cc/300?u=javier001", "https://i.pravatar.cc/300?u=javier002"},
		interests:    []string{"arquitectura", "fotografía", "senderismo", "gastronomía"},
		interestedIn: "both",
	},
	{
		id: "11111111-1111-1111-1111-111111111009",
		email: "natalia@test.com", name: "Natalia Castro", age: 25,
		bio: "Veterinaria con tres gatos. No negociable: tienes que querer a los animales.",
		occupation: "Veterinaria", location: "Granada, España",
		lat: 37.1773, lon: -3.5986,
		photoURLs:    []string{"https://i.pravatar.cc/300?u=natalia001"},
		interests:    []string{"animales", "naturaleza", "yoga", "cocina"},
		interestedIn: "male",
	},
	// User 10: no photos (edge case test)
	{
		id: "11111111-1111-1111-1111-111111111010",
		email: "pablo@test.com", name: "Pablo Herrera", age: 31,
		bio: "Matemático. Busco a alguien que aprecie los patrones en la vida.",
		occupation: "Profesor universitario", location: "Salamanca, España",
		lat: 40.9701, lon: -5.6635,
		photoURLs:    nil, // NO PHOTOS — edge case
		interests:    []string{"matemáticas", "ajedrez", "senderismo"},
		interestedIn: "both",
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

	for _, u := range users {
		// Check if already exists
		var exists int
		db.QueryRow("SELECT COUNT(*) FROM users WHERE id = ?", u.id).Scan(&exists)
		if exists > 0 {
			fmt.Printf("  skip (exists): %s\n", u.email)
			continue
		}

		// 1. Insert user
		_, err := db.Exec(`
			INSERT INTO users (id, email, password_hash, email_verified_at, is_admin, is_frozen, vip_level, credits)
			VALUES (?, ?, ?, NOW(), FALSE, FALSE, 0, 100)`,
			u.id, u.email, password,
		)
		if err != nil {
			log.Printf("ERR insert user %s: %v", u.email, err)
			continue
		}

		// 2. Insert profile
		profileID := u.id[:8] + "-prof-0000-0000-" + u.id[len(u.id)-12:]
		_, err = db.Exec(`
			INSERT INTO user_profiles (id, user_id, name, age, bio, occupation, location, latitude, longitude)
			VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)`,
			profileID, u.id, u.name, u.age, u.bio, u.occupation, u.location, u.lat, u.lon,
		)
		if err != nil {
			log.Printf("ERR insert profile %s: %v", u.email, err)
		}

		// 3. Insert preferences
		prefsID := u.id[:8] + "-pref-0000-0000-" + u.id[len(u.id)-12:]
		_, err = db.Exec(`
			INSERT INTO user_preferences (id, user_id, min_age, max_age, max_distance_km, interested_in)
			VALUES (?, ?, 18, 45, 100, ?)`,
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

		photos := len(u.photoURLs)
		fmt.Printf("  ✓ %s (%s) age=%d photos=%d\n", u.email, u.name, u.age, photos)
	}

	// Seed mutual match between user7 (camila) and user8 (javier)
	fmt.Println("\nSeeding mutual match: camila <-> javier...")
	u7 := "11111111-1111-1111-1111-111111111007"
	u8 := "11111111-1111-1111-1111-111111111008"

	// Swipe u7 -> u8 right
	db.Exec(`INSERT IGNORE INTO swipes (id, swiper_id, swiped_id, direction, created_at)
		VALUES ('sw-7-8-000-0000-000000000001', ?, ?, 'right', NOW())`, u7, u8)
	// Swipe u8 -> u7 right
	db.Exec(`INSERT IGNORE INTO swipes (id, swiper_id, swiped_id, direction, created_at)
		VALUES ('sw-8-7-000-0000-000000000002', ?, ?, 'right', NOW())`, u8, u7)
	// Create match
	_, err2 := db.Exec(`INSERT IGNORE INTO matches (id, user1_id, user2_id, created_at)
		VALUES ('match-7-8-0000-000000000001', ?, ?, NOW())`, u7, u8)
	if err2 != nil {
		fmt.Printf("  match err: %v\n", err2)
	} else {
		fmt.Println("  ✓ match camila <-> javier created")
	}

	// Verify
	fmt.Println("\n--- Verification ---")
	var totalUsers, totalProfiles, totalPhotos int
	db.QueryRow("SELECT COUNT(*) FROM users WHERE deleted_at IS NULL").Scan(&totalUsers)
	db.QueryRow("SELECT COUNT(*) FROM user_profiles").Scan(&totalProfiles)
	db.QueryRow("SELECT COUNT(*) FROM user_photos").Scan(&totalPhotos)
	fmt.Printf("users=%d  profiles=%d  photos=%d\n", totalUsers, totalProfiles, totalPhotos)

	var userNoPhotos int
	db.QueryRow(`SELECT COUNT(*) FROM users u
		LEFT JOIN user_photos ph ON ph.user_id = u.id
		WHERE ph.id IS NULL AND u.deleted_at IS NULL`).Scan(&userNoPhotos)
	fmt.Printf("users without photos: %d\n", userNoPhotos)

	rows, _ := db.Query("SELECT email FROM users WHERE id LIKE '11111111%' ORDER BY email")
	defer rows.Close()
	for rows.Next() {
		var email string
		rows.Scan(&email)
		fmt.Printf("  seeded: %s\n", email)
	}
}
