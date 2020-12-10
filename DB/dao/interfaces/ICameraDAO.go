package interfaces

import 

type ICameraDAO interface {
	Create(u *models.Camera) error
	//Update(u *models.User) error
	Delete(i int) error
	//GetById(i int) (models.User, error)
	GetAll() ([]models.Camera, error)
}
