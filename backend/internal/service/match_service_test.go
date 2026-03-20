package service_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/SX110903/match_app/backend/internal/domain"
	"github.com/SX110903/match_app/backend/internal/service"
)

// ---- mock IMatchRepository ----

type mockMatchRepo struct {
	swipes  map[string]*domain.Swipe // key: swiperID+":"+swipedID
	matches []*domain.Match
}

func newMockMatchRepo() *mockMatchRepo {
	return &mockMatchRepo{swipes: make(map[string]*domain.Swipe)}
}

func (m *mockMatchRepo) CreateSwipe(_ context.Context, s *domain.Swipe) error {
	m.swipes[s.SwiperID+":"+s.SwipedID] = s
	return nil
}

func (m *mockMatchRepo) GetSwipe(_ context.Context, swiperID, swipedID string) (*domain.Swipe, error) {
	s, ok := m.swipes[swiperID+":"+swipedID]
	if !ok {
		return nil, errors.New("not found")
	}
	return s, nil
}

func (m *mockMatchRepo) CreateMatch(_ context.Context, match *domain.Match) error {
	m.matches = append(m.matches, match)
	return nil
}

func (m *mockMatchRepo) GetMatchByID(_ context.Context, _ string) (*domain.Match, error) {
	return nil, errors.New("not found")
}

func (m *mockMatchRepo) GetMatchesByUserID(_ context.Context, _ string) ([]domain.MatchWithProfile, error) {
	return nil, nil
}

func (m *mockMatchRepo) GetCandidates(_ context.Context, _ string, _ *domain.UserPreferences, _, _ int) ([]domain.Candidate, error) {
	return nil, nil
}

func (m *mockMatchRepo) DeleteMatch(_ context.Context, _, _ string) error {
	return nil
}

// ---- mock IProfileRepository (mínimo) ----

type mockProfileRepoForMatch struct {
	existingIDs map[string]bool
}

func (m *mockProfileRepoForMatch) GetByUserID(_ context.Context, id string) (*domain.UserProfile, error) {
	if m.existingIDs[id] {
		return &domain.UserProfile{UserID: id}, nil
	}
	return nil, errors.New("not found")
}

func (m *mockProfileRepoForMatch) Create(_ context.Context, _ *domain.UserProfile) error { return nil }
func (m *mockProfileRepoForMatch) Update(_ context.Context, _ *domain.UserProfile) error { return nil }
func (m *mockProfileRepoForMatch) GetPreferences(_ context.Context, _ string) (*domain.UserPreferences, error) {
	return &domain.UserPreferences{MinAge: 18, MaxAge: 100, MaxDistanceKm: 50, InterestedIn: "both"}, nil
}
func (m *mockProfileRepoForMatch) UpsertPreferences(_ context.Context, _ *domain.UserPreferences) error {
	return nil
}
func (m *mockProfileRepoForMatch) AddPhoto(_ context.Context, _ *domain.UserPhoto) error { return nil }
func (m *mockProfileRepoForMatch) DeletePhoto(_ context.Context, _, _ string) error      { return nil }
func (m *mockProfileRepoForMatch) GetPhotoCount(_ context.Context, _ string) (int, error) {
	return 0, nil
}
func (m *mockProfileRepoForMatch) ReplaceInterests(_ context.Context, _ string, _ []string) error {
	return nil
}

// ---- Tests ----

func newMatchSvc(profileIDs ...string) service.IMatchService {
	ids := make(map[string]bool)
	for _, id := range profileIDs {
		ids[id] = true
	}
	return service.NewMatchService(newMockMatchRepo(), &mockProfileRepoForMatch{existingIDs: ids})
}

func TestSwipe_MutualLike_CreatesMatch(t *testing.T) {
	repo := newMockMatchRepo()
	profiles := &mockProfileRepoForMatch{existingIDs: map[string]bool{"A": true, "B": true}}
	svc := service.NewMatchService(repo, profiles)

	ctx := context.Background()

	// A le da like a B primero
	_, err := svc.Swipe(ctx, "A", "B", domain.SwipeRight)
	if err != nil {
		t.Fatalf("first swipe failed: %v", err)
	}

	// B le da like a A → debe crear match
	resp, err := svc.Swipe(ctx, "B", "A", domain.SwipeRight)
	if err != nil {
		t.Fatalf("second swipe failed: %v", err)
	}
	if !resp.IsMatch {
		t.Fatal("expected IsMatch=true but got false")
	}
	if resp.MatchID == "" {
		t.Fatal("expected non-empty MatchID")
	}
	if len(repo.matches) != 1 {
		t.Fatalf("expected 1 match created, got %d", len(repo.matches))
	}
}

func TestSwipe_UnilateralLike_NoMatch(t *testing.T) {
	repo := newMockMatchRepo()
	profiles := &mockProfileRepoForMatch{existingIDs: map[string]bool{"A": true, "B": true}}
	svc := service.NewMatchService(repo, profiles)

	resp, err := svc.Swipe(context.Background(), "A", "B", domain.SwipeRight)
	if err != nil {
		t.Fatalf("swipe failed: %v", err)
	}
	if resp.IsMatch {
		t.Fatal("expected IsMatch=false for unilateral like")
	}
	if len(repo.matches) != 0 {
		t.Fatalf("expected 0 matches, got %d", len(repo.matches))
	}
}

func TestSwipe_Self_ReturnsError(t *testing.T) {
	svc := newMatchSvc("A")
	_, err := svc.Swipe(context.Background(), "A", "A", domain.SwipeRight)
	if !errors.Is(err, domain.ErrSelfAction) {
		t.Fatalf("expected ErrSelfAction, got %v", err)
	}
}

func TestSwipe_NonExistentTarget_ReturnsNotFound(t *testing.T) {
	svc := newMatchSvc("A") // B no existe
	_, err := svc.Swipe(context.Background(), "A", "B", domain.SwipeRight)
	if !errors.Is(err, domain.ErrNotFound) {
		t.Fatalf("expected ErrNotFound, got %v", err)
	}
}

// suprimir "imported and not used" si time solo se usa en tipos
var _ = time.Now()
