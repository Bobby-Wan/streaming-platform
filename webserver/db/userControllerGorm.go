package db

import (
	"fmt"

	"gorm.io/gorm"
)

//UserControllerGORM is the main db user controller for my app
type UserControllerGORM struct {
	ptrDb *gorm.DB
}

//NewUserControllerGORM creates a gorm user controller for my app
func NewUserControllerGORM(ptrDb *gorm.DB) (*UserControllerGORM, error) {
	if ptrDb == nil {
		return nil, fmt.Errorf("null pointer argument")
	}

	var controller UserControllerGORM
	controller.ptrDb = ptrDb

	return &controller, nil
}

func (ptrC *UserControllerGORM) Get(username string) (*User, error) {
	if ptrC == nil {
		return nil, fmt.Errorf("null receiver")
	}
	if ptrC.ptrDb == nil {
		return nil, fmt.Errorf("null db pointer")
	}

	var user User
	result := ptrC.ptrDb.Where("username=?", username).First(&user)
	if result.Error != nil {
		return nil, result.Error
	}

	return &user, nil
}

func (ptrC *UserControllerGORM) GetByEmail(email string) (*User, error) {
	if ptrC == nil {
		return nil, fmt.Errorf("null receiver")
	}
	if ptrC.ptrDb == nil {
		return nil, fmt.Errorf("null db pointer")
	}

	var user User
	result := ptrC.ptrDb.Where("email=?", email).First(&user)
	if result.Error != nil {
		return nil, result.Error
	}

	return &user, nil
}

func (ptrC *UserControllerGORM) Create(user User) (*User, error) {
	if ptrC == nil {
		return nil, fmt.Errorf("null receiver")
	}
	if ptrC.ptrDb == nil {
		return nil, fmt.Errorf("null db pointer")
	}

	result := ptrC.ptrDb.Create(&user)
	if result.Error != nil {
		return nil, result.Error
	}

	return &user, nil
}

func (ptrC *UserControllerGORM) Delete(id uint) error {
	if ptrC == nil {
		return fmt.Errorf("null receiver")
	}
	if ptrC.ptrDb == nil {
		return fmt.Errorf("null db pointer")
	}

	result := ptrC.ptrDb.Delete(&User{}, id)
	if result.Error != nil {
		return result.Error
	}

	return nil
}

func (ptrC *UserControllerGORM) Update(user User) error {
	if ptrC == nil {
		return fmt.Errorf("null receiver")
	}
	if ptrC.ptrDb == nil {
		return fmt.Errorf("null db pointer")
	}

	result := ptrC.ptrDb.Save(&user)
	if result.Error != nil {
		return result.Error
	}

	return nil
}
