package recreation

import (
	"github.com/atletaid/go-template/src/model"
)

type RecreationRepository interface {
	CreateRecreation(recreation *model.Recreation) (int64, error)
	FindRecreationByID(recreationID int64) (*model.Recreation, error)
	FindAllRecreations() (model.Recreations, error)
	FindByLocation(cityName string) (model.Recreations, error)
	DeleteRecreation(recreationID int64) error
}
