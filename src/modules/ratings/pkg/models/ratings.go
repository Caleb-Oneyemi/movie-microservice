package models

type RecordID string

type RecordType string

const (
	RecordTypeMovie = RecordType("movie")
)

type UserID string

type RatingValue int

type Rating struct {
	RecordID   string      `json:"recordId"`
	RecordType string      `json:"recordType"`
	UserID     string      `json:"userId"`
	Value      RatingValue `json:"value"`
}

type RatingEventType string

type RatingEvent struct {
	UserID     UserID          `json:"userId"`
	RecordID   RecordID        `json:"recordId"`
	RecordType RecordType      `json:"recordType"`
	Value      RatingValue     `json:"value"`
	Type       RatingEventType `json:"type"`
}

const (
	RatingEventTypePut    = "put"
	RatingEventTypeDelete = "delete"
)
