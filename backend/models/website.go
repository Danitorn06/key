package models

import "time"

type Website struct {
    ID             int       `json:"id"`
    AffiliatorName string    `json:"affiliator_name,omitempty"` // ใช้แค่ตอน frontend ส่งมา
    WebsiteURL     string    `json:"website_url"`
    CreatedAt      time.Time `json:"created_at,omitempty"`     // Response เท่านั้น
}

