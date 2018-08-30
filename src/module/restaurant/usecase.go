package restaurant

import (
	"log"

	"github.com/atletaid/go-template/src/model"
)

type Usecase interface {
	CreateRestaurant(restaurantName, restaurantCity, restaurantImage, restaurantDescription string, restaurantTime, restaurantPrice int, positionLat, positionLong float64) (int64, error)
	GetRestaurant(restaurantID int64) (*model.Restaurant, error)
	GetAllRestaurants() (model.Restaurants, error)
	GetRestaurantsByCity(cityName string) (model.Restaurants, error)
	DeleteRestaurantByID(restaurantID int64) error
}

type usecase struct {
	restaurantRepo RestaurantRepository
}

func NewRestaurantUsecase(
	restaurantRepo RestaurantRepository,
) Usecase {
	return &usecase{
		restaurantRepo: restaurantRepo,
	}
}

func (u *usecase) CreateRestaurant(restaurantName, restaurantCity, restaurantImage, restaurantDescription string, restaurantTime, restaurantPrice int, positionLat, positionLong float64) (int64, error) {
	newRestaurant := model.NewRestaurant(restaurantName, restaurantCity, restaurantImage, restaurantDescription, restaurantTime, restaurantPrice, positionLat, positionLong)
	restaurantID, err := u.restaurantRepo.CreateRestaurant(newRestaurant)
	if err != nil {
		log.Println(err)
		return 0, err
	}

	return restaurantID, nil
}

func (u *usecase) GetAllRestaurants() (model.Restaurants, error) {
	restaurants, err := u.restaurantRepo.FindAllRestaurants()
	if err != nil {
		log.Println(err)
		return nil, err
	}

	return restaurants, nil
}

func (u *usecase) GetRestaurant(venueID int64) (*model.Restaurant, error) {
	restaurant, err := u.restaurantRepo.FindRestaurantByID(venueID)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	return restaurant, nil
}

func (u *usecase) GetRestaurantsByCity(cityName string) (model.Restaurants, error) {
	restaurants, err := u.restaurantRepo.FindByLocation(cityName)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	return restaurants, nil
}

func (u *usecase) DeleteRestaurantByID(restaurantID int64) error {
	if _, err := u.restaurantRepo.FindRestaurantByID(restaurantID); err != nil {
		log.Println(err)
		return err
	}

	if err := u.restaurantRepo.DeleteRestaurantID(restaurantID); err != nil {
		log.Println(err)
		return err
	}

	return nil
}
