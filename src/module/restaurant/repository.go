package restaurant

import (
	"github.com/atletaid/go-template/src/model"
)

type RestaurantRepository interface {
	CreateRestaurant(*model.Restaurant) (int64, error)
	FindRestaurantByID(restaurantID int64) (*model.Restaurant, error)
	FindAllRestaurants() (model.Restaurants, error)
	FindByLocation(cityName string) (model.Restaurants, error)
	DeleteRestaurantID(int64) error
}
