// Package stats implements and describes short URLs service statistics.
package stats

//go:generate  mockgen -source=./internal/domain/stats/stats.go -destination=./internal/domain/stats/mock/mock_stats.go -package=statsMock

// UserStats describes the implementation of the short URLs service user statistics.
type UserStats interface {
	Count() int
}

// URLStats describes the implementation of the short URLs service url statistics.
type URLStats interface {
	Count() int
}

// Collection describes the implementation of the short URLs service statistics.
type Collection interface {
	User() UserStats
	URL() URLStats
}

type user struct {
	count *int
}

type url struct {
	count *int
}

type collection struct {
	userStats UserStats
	urlStats  URLStats
}

// UserStat describes the implementation of the user statistic type.
type UserStat func(u *user)

// URLStat describes the implementation of the url statistic type.
type URLStat func(ul *url)

// Count implements getting user count.
func (u *user) Count() int {
	return *u.count
}

// Count implements getting url count.
func (ul *url) Count() int {
	return *ul.count
}

// User implements getting user statistic.
func (c *collection) User() UserStats {
	return c.userStats
}

// URL implements getting url statistic.
func (c *collection) URL() URLStats {
	return c.urlStats
}

// UserCount implements setting user count.
func UserCount(count int) UserStat {
	return func(u *user) {
		u.count = &count
	}
}

// URLCount implements setting url count.
func URLCount(count int) URLStat {
	return func(ul *url) {
		ul.count = &count
	}
}

// NewUserStats implements the creation of the user statistic type.
func NewUserStats(stats ...UserStat) UserStats {
	u := new(user)
	for _, stat := range stats {
		stat(u)
	}
	return u
}

// NewURLStats implements the creation of the url statistic type.
func NewURLStats(stats ...URLStat) URLStats {
	ul := new(url)
	for _, stat := range stats {
		stat(ul)
	}
	return ul
}

// NewCollectionStats implements the creation of the short URLs service statistics.
func NewCollectionStats(user UserStats, url UserStats) *collection {
	return &collection{
		user,
		url,
	}
}
