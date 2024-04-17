package models

import "moviemicroservice.com/src/modules/metadata/pkg/models"

type MovieDetails struct {
	Ratings  *float64        `json:"rating,omitempty"`
	Metadata models.MetaData `json:"metadata"`
}
