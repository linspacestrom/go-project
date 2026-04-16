package platform

import "context"

func (s *Service) GetBusinessAnalytics(ctx context.Context) (*struct {
	MentorApprovals                int64
	MentorRejections               int64
	ActiveBookings                 int64
	MentorRequestsBySkill          int64
	MentorRequestsByCategory       int64
	UserCancelsAfterMentorApproval int64
}, error) {
	return s.repo.GetBusinessAnalytics(ctx)
}

func (s *Service) GetTechnicalAnalytics(ctx context.Context) (*struct {
	BookingConflictCount int64
	FailedOutboxEvents   int64
	OutboxLagCount       int64
}, error) {
	return s.repo.GetTechnicalAnalytics(ctx)
}
