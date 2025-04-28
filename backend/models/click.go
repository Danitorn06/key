package models

import "time"

type ClickLog struct {
    ID           int       `json:"id,omitempty"`
    ItemClicked  string    `json:"item_clicked"`
    ReferrerURL  string    `json:"referrer_url"`
    ClickedAt    time.Time `json:"clicked_at,omitempty"` // เวลาคลิก (ระบบ backend เติมให้)
}