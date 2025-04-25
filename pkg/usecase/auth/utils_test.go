package auth

import (
	"product_recommendation/mocks"
)

func ClearMock(mockAuthRepo *mocks.AuthRepository, mockDB *mocks.DBInterface) {
	if mockAuthRepo != nil {
		mockAuthRepo.ExpectedCalls = nil
		mockAuthRepo.Calls = nil
	}
	if mockDB != nil {
		mockDB.ExpectedCalls = nil
		mockDB.Calls = nil
	}
}
