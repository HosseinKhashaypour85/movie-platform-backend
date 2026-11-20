package admin

import (
	"net/http"
	"strconv"
	"time"

	"movie-platform/internal/httpx"
	"movie-platform/modules/catalog"
	"movie-platform/modules/engagement"
	"movie-platform/modules/users"
	"movie-platform/modules/watch_history"
)

type dailyMetric struct {
	Date     string `json:"date"`
	Views    int64  `json:"views"`
	NewUsers int64  `json:"new_users"`
}
type topMovieMetric struct {
	MovieID uint   `json:"movie_id"`
	Title   string `json:"title"`
	Views   int64  `json:"views"`
}

func (h *Handler) Analytics(w http.ResponseWriter, r *http.Request) {
	days, _ := strconv.Atoi(r.URL.Query().Get("days"))
	if days == 0 {
		days = 30
	}
	if days != 7 && days != 30 && days != 90 {
		httpx.Error(w, 400, "invalid_range", "days must be 7, 30 or 90")
		return
	}
	now := time.Now().UTC()
	from := now.AddDate(0, 0, -days)
	var totalUsers, activeUsers, movies, series, views, favorites, ratings int64
	checks := []struct {
		model any
		query string
		dest  *int64
	}{
		{&users.User{}, "", &totalUsers},
		{&users.User{}, "is_active = ?", &activeUsers},
		{&catalog.Movie{}, "is_published = ?", &movies},
		{&catalog.Series{}, "is_published = ?", &series},
		{&watch_history.WatchHistory{}, "watched_at >= ?", &views},
		{&engagement.Favorite{}, "", &favorites},
		{&engagement.Rating{}, "", &ratings},
	}
	for i, item := range checks {
		var err error
		switch i {
		case 1:
			err = h.db.Model(item.model).Where(item.query, true).Count(item.dest).Error
		case 2, 3:
			err = h.db.Model(item.model).Where(item.query, true).Count(item.dest).Error
		case 4:
			err = h.db.Model(item.model).Where(item.query, from).Count(item.dest).Error
		default:
			err = h.db.Model(item.model).Count(item.dest).Error
		}
		if err != nil {
			httpx.Error(w, 500, "internal_error", "analytics query failed")
			return
		}
	}
	var daily []dailyMetric
	if err := h.db.Model(&watch_history.WatchHistory{}).Select("DATE(watched_at) AS date, COUNT(*) AS views").Where("watched_at >= ?", from).Group("DATE(watched_at)").Order("date ASC").Scan(&daily).Error; err != nil {
		httpx.Error(w, 500, "internal_error", "analytics query failed")
		return
	}
	var userDaily []struct {
		Date     string
		NewUsers int64
	}
	if err := h.db.Model(&users.User{}).Select("DATE(created_at) AS date, COUNT(*) AS new_users").Where("created_at >= ?", from).Group("DATE(created_at)").Scan(&userDaily).Error; err != nil {
		httpx.Error(w, 500, "internal_error", "analytics query failed")
		return
	}
	newUsers := make(map[string]int64, len(userDaily))
	for _, v := range userDaily {
		newUsers[v.Date] = v.NewUsers
	}
	for i := range daily {
		daily[i].NewUsers = newUsers[daily[i].Date]
	}
	var top []topMovieMetric
	if err := h.db.Table("watch_histories").Select("watch_histories.movie_id, movies.title, COUNT(*) AS views").Joins("JOIN movies ON movies.id = watch_histories.movie_id").Where("watch_histories.watched_at >= ?", from).Group("watch_histories.movie_id,movies.title").Order("views DESC").Limit(10).Scan(&top).Error; err != nil {
		httpx.Error(w, 500, "internal_error", "analytics query failed")
		return
	}
	httpx.JSON(w, 200, map[string]any{"range_days": days, "generated_at": now, "summary": map[string]int64{"users": totalUsers, "active_users": activeUsers, "published_movies": movies, "published_series": series, "views_in_range": views, "favorites": favorites, "ratings": ratings}, "daily": daily, "top_movies": top})
}
