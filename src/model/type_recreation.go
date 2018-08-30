package model

import (
	"time"
)

type Recreation struct {
	RecreationID          int64     `json:"recreation_id"`
	RecreationName        string    `json:"recreation_name"`
	RecreationTimeMinute  int       `json:"recreation_time_minute"`
	RecreationPrice       int       `json:"recreation_price"`
	PositionLat           float64   `json:"position_lat"`
	PositionLong          float64   `json:"position_long"`
	RecreationCity        string    `json:"recreation_city"`
	RecreationImage       string    `json:"recreation_image"`
	RecreationDescription string    `json:"recreation_description"`
	CreatedAt             time.Time `json:"created_at"`
}

type Recreations []*Recreation

func NewRecreation(recreationName, recreationCity, recreationImage, recreationDescription string, recrationTime, recreationPrice int, positionLat, positionLong float64) *Recreation {
	return &Recreation{
		RecreationName:        recreationName,
		RecreationTimeMinute:  recrationTime,
		RecreationPrice:       recreationPrice,
		PositionLat:           positionLat,
		PositionLong:          positionLong,
		RecreationCity:        recreationCity,
		RecreationImage:       recreationImage,
		RecreationDescription: recreationDescription,
	}
}
