package models

type SystemStats struct {
	ActiveJobs   int64 `json:"active_jobs"`
	DisabledJobs int64 `json:"disabled_jobs"`
	Jobs         int64 `json:"jobs"`
	ErrorCount   int64 `json:"error_count"`
	SuccessCount int64 `json:"success_count"`
}
