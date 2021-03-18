package db

import (
	"errors"
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

// func (ptrC *StreamControllerGORM) Get(userId uint) (*ActiveStream, error) {
// 	if ptrC == nil {
// 		return nil, fmt.Errorf("null receiver")
// 	}
// 	if ptrC.ptrDb == nil {
// 		return nil, fmt.Errorf("null db pointer")
// 	}

// 	var stream ActiveStream
// 	result := ptrC.ptrDb.Where("userId=?", userId).First(&stream)
// 	if result.Error != nil {
// 		return nil, result.Error
// 	}

// 	return &stream, nil
// }

func (ptrC *StreamControllerGORM) Create(stream ActiveStream) (*ActiveStream, error) {
	if ptrC == nil {
		return nil, fmt.Errorf("null receiver")
	}
	if ptrC.ptrDb == nil {
		return nil, fmt.Errorf("null db pointer")
	}

	result := ptrC.ptrDb.Create(&stream)
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
	result := ptrC.ptrDb.Where("category & ?", mask).Find(&streams)
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

func (ptrC *StreamControllerGORM) End(username string) (*ActiveStream, error) {
	if ptrC == nil {
		return nil, errors.New("nil pointer provided")
	}
	if ptrC.ptrDb == nil {
		return nil, errors.New("nil db pointer in controller")
	}

	var stream ActiveStream

	result := ptrC.ptrDb.Where("username=?", username).First(&stream)
	if result.Error != nil {
		return nil, result.Error
	}

	stream.Active = false

	result = ptrC.ptrDb.Save(&stream)
	if result.Error != nil {
		return nil, result.Error
	}

	return &stream, nil
}

func (ptrC *StreamControllerGORM) Start(username string) (*ActiveStream, error) {
	if ptrC == nil {
		return nil, errors.New("nil pointer provided")
	}
	if ptrC.ptrDb == nil {
		return nil, errors.New("nil db pointer in controller")
	}

	var stream ActiveStream

	result := ptrC.ptrDb.Where("username=?", username).First(&stream)
	if result.Error != nil {
		return nil, result.Error
	}

	stream.Active = true

	result = ptrC.ptrDb.Save(&stream)
	if result.Error != nil {
		return nil, result.Error
	}

	return &stream, nil
}

func (ptrC *StreamControllerGORM) GetByUsername(username string) (*ActiveStream, error) {
	if ptrC == nil {
		return nil, errors.New("nil pointer provided")
	}
	if ptrC.ptrDb == nil {
		return nil, errors.New("nil db pointer in controller")
	}

	var stream ActiveStream
	result := ptrC.ptrDb.Where("username=?", username).Limit(1).Find(&stream)
	if result.Error != nil {
		return nil, errors.New("error fetching stream")
	}

	return &stream, nil
}

func (ptrC *StreamControllerGORM) Delete(userId uint) error {
	if ptrC == nil {
		return errors.New("nil pointer provided")
	}
	if ptrC.ptrDb == nil {
		return errors.New("nil db pointer in controller")
	}

	result := ptrC.ptrDb.Unscoped().Where("user_id=?", userId).Delete(ActiveStream{})
	if result.Error != nil {
		return result.Error
	}

	return nil
}

// func paginate(page int, size int) func(db *gorm.DB) *gorm.DB {
// 	return func(db *gorm.DB) *gorm.DB {
// 		if page < 1 {
// 			page = 1
// 		}

// 		switch {
// 		case size > 100:
// 			size = 100
// 		case size <= 0:
// 			size = 10
// 		}

// 		offset := (page - 1) * size
// 		return db.Offset(offset).Limit(size)
// 	}
// }
// func (ptrC *StreamControllerGORM) Paginate(page int, size int) error {
// 	if ptrC == nil {
// 		return errors.New("nil pointer provided")
// 	}
// 	if ptrC.ptrDb == nil {
// 		return errors.New("nil db pointer in controller")
// 	}
// }
