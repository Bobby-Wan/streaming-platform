package db

import (
	"fmt"

	"gorm.io/gorm"
)

//StreamControllerGORM is the main db stream controller for my app
type StreamControllerGORM struct {
	ptrDb *gorm.DB
}

//NewStreamControllerGORM creates a gorm stream controller for my app
func NewStreamControllerGORM(ptrDb *gorm.DB) (*StreamControllerGORM, error) {
	if ptrDb == nil {
		return nil, fmt.Errorf("null pointer argument")
	}

	var controller StreamControllerGORM
	controller.ptrDb = ptrDb

	return &controller, nil
}

func (ptrC *StreamControllerGORM) Get(userId uint) (*ActiveStream, error) {
	if ptrC == nil {
		return nil, fmt.Errorf("null receiver")
	}
	if ptrC.ptrDb == nil {
		return nil, fmt.Errorf("null db pointer")
	}

	var stream ActiveStream
	result := ptrC.ptrDb.Where("userId=?", userId).First(&stream)
	if result.Error != nil {
		return nil, result.Error
	}

	return &stream, nil
}

func (ptrC *StreamControllerGORM) Create(stream ActiveStream) (*ActiveStream, error) {
	if ptrC == nil {
		return nil, fmt.Errorf("null receiver")
	}
	if ptrC.ptrDb == nil {
		return nil, fmt.Errorf("null db pointer")
	}

	result := ptrC.ptrDb.Create(stream)
	if result.Error != nil {
		return nil, result.Error
	}

	return &stream, nil
}

func (ptrC *StreamControllerGORM) GetByCategories(mask uint32) ([]ActiveStream, error) {
	if ptrC == nil {
		return nil, fmt.Errorf("null receiver")
	}
	if ptrC.ptrDb == nil {
		return nil, fmt.Errorf("null db pointer")
	}

	var streams []ActiveStream
	result := ptrC.ptrDb.Where("categories & ?", mask).Find(&streams)
	if result.Error != nil {
		return nil, result.Error
	}

	return streams, nil
}

func (ptrC *StreamControllerGORM) GetByViewCount(limit int) ([]ActiveStream, error) {
	if ptrC == nil {
		return nil, fmt.Errorf("null receiver")
	}
	if ptrC.ptrDb == nil {
		return nil, fmt.Errorf("null db pointer")
	}

	var streams []ActiveStream
	result := ptrC.ptrDb.Order("viewers").Limit(limit).Find(&streams)
	if result.Error != nil {
		return nil, result.Error
	}

	return streams, nil
}
