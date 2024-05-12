package medium

import types "github.com/sikozonpc/notebase/types"

type Medium interface {
	SendInsights(u *types.User, insights []*types.DailyInsight, authToken string) error
}
