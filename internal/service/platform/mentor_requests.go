package platform

import (
	"context"

	"github.com/google/uuid"
	"github.com/linspacestrom/go-project/internal/domain"
	"github.com/linspacestrom/go-project/internal/repository/postgres"
)

func (s *Service) CreateMentorRequest(ctx context.Context, req *domain.MentorRequest) (*domain.MentorRequest, error) {
	req.Status = domain.MentorRequestStatusPending

	created, err := s.repo.CreateMentorRequest(ctx, req)
	if err != nil {
		return nil, err
	}

	if err = s.publishEvent(ctx, "mentor_request", &created.ID, domain.EventMentorRequestCreate, map[string]any{"request_id": created.ID, "city_id": created.CityID}); err != nil {
		return nil, err
	}

	if created.SkillID != nil {
		mentors, listErr := s.repo.ListMentorSubscribersBySkill(ctx, *created.SkillID, created.CityID)
		if listErr != nil {
			return nil, listErr
		}
		for _, mentorID := range mentors {
			payload := map[string]any{
				"request_id": created.ID,
				"mentor_id":  mentorID,
				"skill_id":   *created.SkillID,
			}
			if outErr := s.publishEvent(ctx, "mentor_request", &created.ID, domain.EventMentorSkillMatched, payload); outErr != nil {
				return nil, outErr
			}
		}
	}
	if created.MentorID != nil {
		if outErr := s.publishEvent(ctx, "mentor_request", &created.ID, domain.EventMentorSkillMatched, map[string]any{"request_id": created.ID, "mentor_id": *created.MentorID}); outErr != nil {
			return nil, outErr
		}
	}

	return created, nil
}

func (s *Service) ListMentorRequests(ctx context.Context, filter postgres.MentorRequestFilter) ([]domain.MentorRequest, uint64, error) {
	return s.repo.ListMentorRequests(ctx, filter)
}

func (s *Service) UpdateMentorRequestStatus(ctx context.Context, requestID uuid.UUID, status string, mentorID *uuid.UUID) (*domain.MentorRequest, error) {
	req, err := s.repo.UpdateMentorRequestStatus(ctx, requestID, status, mentorID)
	if err != nil {
		return nil, err
	}
	if outErr := s.publishEvent(ctx, "mentor_request", &req.ID, domain.EventMentorRequestUpdate, map[string]any{"request_id": req.ID, "status": status}); outErr != nil {
		return nil, outErr
	}

	return req, nil
}

func (s *Service) SubscribeMentorSkill(ctx context.Context, mentorID, skillID uuid.UUID) error {
	return s.repo.UpsertMentorSkillSubscription(ctx, mentorID, skillID)
}

func (s *Service) UnsubscribeMentorSkill(ctx context.Context, mentorID, skillID uuid.UUID) error {
	return s.repo.DeleteMentorSkillSubscription(ctx, mentorID, skillID)
}
