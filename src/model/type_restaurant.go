package model

import (
	"time"
)

type Restaurant struct {
	RestaurantID          int64     `json:"restaurant_id"`
	RestaurantName        string    `json:"restaurant_name"`
	RestaurantTimeMinute  int       `json:"restaurant_time_minute"`
	RestaurantPrice       int       `json:"restaurant_price"`
	PositionLat           float64   `json:"position_lat"`
	PositionLong          float64   `json:"position_long"`
	RestaurantCity        string    `json:"restaurant_city"`
	RestaurantImage       string    `json:"restaurant_image"`
	RestaurantDescription string    `json:"restaurant_description"`
	CreatedAt             time.Time `json:"created_at"`
}

type Restaurants []*Restaurant

func NewRestaurant(restaurantName, restaurantCity, restaurantImage, restaurantDescription string, restaurantTime, restaurantPrice int, positionLat, positionLong float64) *Restaurant {
	return &Restaurant{
		RestaurantName:        restaurantName,
		RestaurantTimeMinute:  restaurantTime,
		RestaurantPrice:       restaurantPrice,
		PositionLat:           positionLat,
		PositionLong:          positionLong,
		RestaurantCity:        restaurantCity,
		RestaurantImage:       restaurantImage,
		RestaurantDescription: restaurantDescription,
	}
}
