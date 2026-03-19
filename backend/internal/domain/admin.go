package domain

import "time"

type AdminLog struct {
	ID        string     `db:"id"         json:"id"`
	AdminID   string     `db:"admin_id"   json:"admin_id"`
	TargetID  *string    `db:"target_id"  json:"target_id,omitempty"`
	Action    string     `db:"action"     json:"action"`
	Details   *string    `db:"details"    json:"details,omitempty"`
	CreatedAt time.Time  `db:"created_at" json:"created_at"`
}

type NotificationSettings struct {
	UserID      string `db:"user_id"      json:"user_id"`
	NewMatches  bool   `db:"new_matches"  json:"new_matches"`
	NewMessages bool   `db:"new_messages" json:"new_messages"`
	NewsUpdates bool   `db:"news_updates" json:"news_updates"`
	Marketing   bool   `db:"marketing"    json:"marketing"`
}

type PrivacySettings struct {
	UserID           string `db:"user_id"            json:"user_id"`
	ShowOnlineStatus bool   `db:"show_online_status" json:"show_online_status"`
	ShowLastSeen     bool   `db:"show_last_seen"     json:"show_last_seen"`
	ShowDistance     bool   `db:"show_distance"      json:"show_distance"`
	IncognitoMode    bool   `db:"incognito_mode"     json:"incognito_mode"`
}
