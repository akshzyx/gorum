package post

import "time"

type PaginatedPosts struct {
	Data       []map[string]interface{} `json:"data"`
	NextCursor *time.Time               `json:"next_cursor"`
	HasMore    bool                     `json:"has_more"`
}
