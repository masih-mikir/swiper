package recreation

import (
	"log"

	"github.com/atletaid/go-template/src/model"
)

type Usecase interface {
	CreateRecreation(recreationName, recreationCity, recreationImage, recreationDescription string, recrationTime, recreationPrice int, positionLat, positionLong float64) (int64, error)
	GetRecreation(recreationID int64) (*model.Recreation, error)
	GetAllRecrations() (model.Recreations, error)
	GetRecreationsByCity(cityName string) (model.Recreations, error)
	DeleteRecreationByID(recrationID int64) error
}

type usecase struct {
	recreationRepo RecreationRepository
}

func NewRecreationUsecase(
	recreationRepo RecreationRepository,
) Usecase {
	return &usecase{
		recreationRepo: recreationRepo,
	}
}

func (u *usecase) CreateRecreation(recreationName, recreationCity, recreationImage, recreationDescription string, recrationTime, recreationPrice int, positionLat, positionLong float64) (int64, error) {
	newRecreation := model.NewRecreation(recreationName, recreationCity, recreationImage, recreationDescription, recrationTime, recreationPrice, positionLat, positionLong)
	recreationID, err := u.recreationRepo.CreateRecreation(newRecreation)
	if err != nil {
		log.Println(err)
		return 0, err
	}

	return recreationID, nil
}

func (u *usecase) GetAllRecrations() (model.Recreations, error) {
	recreations, err := u.recreationRepo.FindAllRecreations()
	if err != nil {
		log.Println(err)
		return nil, err
	}

	return recreations, nil
}

func (u *usecase) GetRecreation(venueID int64) (*model.Recreation, error) {
	recreation, err := u.recreationRepo.FindRecreationByID(venueID)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	return recreation, nil
}

func (u *usecase) GetRecreationsByCity(cityName string) (model.Recreations, error) {
	recreations, err := u.recreationRepo.FindByLocation(cityName)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	return recreations, nil
}

func (u *usecase) DeleteRecreationByID(recreationID int64) error {
	if _, err := u.recreationRepo.FindRecreationByID(recreationID); err != nil {
		log.Println(err)
		return err
	}

	if err := u.recreationRepo.DeleteRecreation(recreationID); err != nil {
		log.Println(err)
		return err
	}

	return nil
}
