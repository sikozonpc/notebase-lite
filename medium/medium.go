package medium

import t "github.com/sikozonpc/notebase/types"

type Medium interface {
	SendInsights(u *t.User, insights []*t.DailyInsight) error
}
