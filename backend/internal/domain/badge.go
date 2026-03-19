package domain

const (
	BadgeNone           = "none"
	BadgeInfluencer     = "influencer"
	BadgeVerified       = "verified"
	BadgeVerifiedGov    = "verified_gov"
	VerifyCost          = 50000
	InfluencerThreshold = 100
)

var ValidBadges = map[string]bool{
	BadgeNone:        true,
	BadgeInfluencer:  true,
	BadgeVerified:    true,
	BadgeVerifiedGov: true,
}
