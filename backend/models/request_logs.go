package models

import "time"

type RequestLog struct {
    ID          int       `json:"id,omitempty"`
    AffiliatorID int      `json:"affiliator_id"`
    Endpoint    string    `json:"endpoint"`
    Method      string    `json:"method"`
    Parameters  string    `json:"parameters"`
    RequestedAt time.Time `json:"requested_at,omitempty"`
}