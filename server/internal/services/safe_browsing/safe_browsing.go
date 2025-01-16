// // internal/services/safebrowsing/service.go

package safebrowsing

// import (
// 	"context"
// 	"fmt"
// 	"time"

// 	"github.com/your-username/your-project/internal/core/ports"
// 	safebrowsing "google.golang.org/api/safebrowsing/v4"
// )

// type Service struct {
// 	client *safebrowsing.Service
// }

// func NewService() ports.SafeBrowsingService {
// 	return &Service{}
// }

// func (s *Service) Initialize(ctx context.Context) error {
// 	client, err := safebrowsing.NewService(ctx)
// 	if err != nil {
// 		return fmt.Errorf("failed to create Safe Browsing service: %w", err)
// 	}
// 	s.client = client
// 	return nil
// }

// func (s *Service) IsSafeURL(ctx context.Context, url string) (bool, error) {
// 	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
// 	defer cancel()

// 	req := &safebrowsing.ThreatUrlLookup{
// 		ThreatInfo: &safebrowsing.ThreatInfo{
// 			ThreatTypes: []string{
// 				"MALWARE",
// 				"SOCIAL_ENGINEERING",
// 				"UNWANTED_SOFTWARE",
// 				"POTENTIALLY_HARMFUL_APPLICATION",
// 			},
// 			PlatformTypes:    []string{"ANY_PLATFORM"},
// 			ThreatEntryTypes: []string{"URL"},
// 			ThreatEntries: []*safebrowsing.ThreatEntry{
// 				{Url: url},
// 			},
// 		},
// 	}

// 	resp, err := s.client.ThreatMatches.Find(req).Context(ctx).Do()
// 	if err != nil {
// 		return false, fmt.Errorf("failed to check URL safety: %w", err)
// 	}

// 	return len(resp.Matches) == 0, nil
// }
