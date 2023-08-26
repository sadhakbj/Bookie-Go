package helpers

import (
	"encoding/base64"
	"encoding/json"
	"time"
)

// PaginationInfo stores the info about next and prev cursors.
type PaginationInfo struct {
	NextCursor string `json:"next_cursor"`
	PrevCursor string `json:"prev_cursor"`
}

// Cursor represents a generic map
type Cursor map[string]interface{}

// CreateCursor creates new cursor
func CreateCursor(id string, createdAt time.Time, pointsNext bool) Cursor {
	return Cursor{
		"id":          id,
		"created_at":  createdAt,
		"points_next": pointsNext,
	}
}

// GeneratePager generates the pager
func GeneratePager(next Cursor, prev Cursor) PaginationInfo {
	return PaginationInfo{
		NextCursor: encodeCursor(next),
		PrevCursor: encodeCursor(prev),
	}
}

func encodeCursor(cursor Cursor) string {
	if len(cursor) == 0 {
		return ""
	}
	serializedCursor, err := json.Marshal(cursor)
	if err != nil {
		return ""
	}
	encodedCursor := base64.StdEncoding.EncodeToString(serializedCursor)
	return encodedCursor
}

// DecodeCursor decodes the cursor
func DecodeCursor(cursor string) (Cursor, error) {
	decodedCursor, err := base64.StdEncoding.DecodeString(cursor)
	if err != nil {
		return nil, err
	}

	var cur Cursor
	if err := json.Unmarshal(decodedCursor, &cur); err != nil {
		return nil, err
	}
	return cur, nil
}
