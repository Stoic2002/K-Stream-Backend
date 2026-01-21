package analytics

type DashboardStats struct {
	TotalUsers    int `json:"total_users"`
	TotalDramas   int `json:"total_dramas"`
	TotalEpisodes int `json:"total_episodes"`
	TotalViews    int `json:"total_views"`
}
