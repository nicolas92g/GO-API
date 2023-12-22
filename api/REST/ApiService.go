package REST

import (
	"projet/api/commons/apiError"
)

type ApiService struct {
	repos ApiRepository
}

func (as *ApiService) Connect() apiError.ApiError {
	return as.repos.Connect()
}

func (as *ApiService) Close() {
	as.repos.Close()
}

// USER

func (as *ApiService) GetCurrentUser(api_key string) (UserModel, apiError.ApiError) {

	if !(as.repos.IsUserKeyExists(api_key)) {
		return UserModel{}, apiError.ValueNotFound{Msg: "Your key is not valid"}
	}

	return as.repos.GetCurrentUser(api_key)
}

func (as *ApiService) GetUsers(limit uint) ([]UserModel, apiError.ApiError) {
	return as.repos.GetUsers(limit)
}

func (as *ApiService) CreateUser(u UserModel) apiError.ApiError {
	if as.repos.IsUserIdExists(u.Id) {
		return apiError.AlreadyExists{Msg: "This user_id already exists"}
	}

	return as.repos.CreateUser(u)
}

func (as *ApiService) IsUserKeyExists(api_key string) bool {
	return as.repos.IsUserKeyExists(api_key)
}

func (as *ApiService) DeleteUser(id uint) apiError.ApiError {
	if !(as.repos.IsUserIdExists(id)) {
		return apiError.ValueNotFound{Msg: "The user id was not found"}
	}
	return as.repos.DeleteUser(id)
}

func (as *ApiService) GetUserRents(id uint) ([]RentModel, apiError.ApiError) {
	if !(as.repos.IsUserIdExists(id)) {
		return []RentModel{}, apiError.ValueNotFound{Msg: "The user id was not found"}
	}
	return as.repos.GetUserRents(id)
}

// APPARTEMENT

func (as *ApiService) GetAppartements(limit uint) ([]AppartementModel, apiError.ApiError) {
	return as.repos.GetAppartements(limit)
}

func (as *ApiService) GetAppartement(id uint) (AppartementModel, apiError.ApiError) {

	//check that it exists
	if !(as.repos.IsAppartementExists(id)) {
		return AppartementModel{}, apiError.ValueNotFound{Msg: "Apartment id doesn't exists"}
	}

	return as.repos.GetAppartement(id)
}

func (as *ApiService) CreateAppartement(a AppartementModel) apiError.ApiError {
	return as.repos.CreateAppartement(a)
}

func (as *ApiService) DeleteAppartement(id uint) apiError.ApiError {
	//check that it exists
	if !(as.repos.IsAppartementExists(id)) {
		return apiError.ValueNotFound{Msg: "apartment id doesn't exists"}
	}
	return as.repos.DeleteAppartement(id)
}

func (as *ApiService) UpdateAppartement(flat AppartementModel) apiError.ApiError {
	if !(as.repos.IsAppartementExists(flat.Id)) {
		return apiError.ValueNotFound{Msg: "apartment id doesn't exists"}
	}
	return as.repos.UpdateAppartement(flat)
}

//owns

func (as *ApiService) GetUserOwns(user uint) ([]OwnModel, apiError.ApiError) {
	if !(as.repos.IsUserIdExists(user)) {
		return []OwnModel{}, apiError.ValueNotFound{Msg: "The user id was not found"}
	}
	return as.repos.GetUserOwns(user)
}

func (as *ApiService) IsOwnRelationExists(own OwnModel) bool {
	return as.repos.IsOwnRelationExists(own)
}

func (as *ApiService) GetAppartementOwners(appartement uint) ([]OwnModel, apiError.ApiError) {
	//check that it exists
	if !(as.repos.IsAppartementExists(appartement)) {
		return []OwnModel{}, apiError.ValueNotFound{Msg: "Apartment id doesn't exists"}
	}
	return as.repos.GetAppartementOwners(appartement)
}

func (as *ApiService) CreateOwnRelation(own OwnModel) apiError.ApiError {
	if as.IsOwnRelationExists(own) {
		return apiError.AlreadyExists{Msg: "This own relation already exists"}
	}
	return as.repos.CreateOwnRelation(own)
}

func (as *ApiService) DeleteOwnRelation(own OwnModel) apiError.ApiError {
	if !(as.IsOwnRelationExists(own)) {
		return apiError.ValueNotFound{Msg: "This own relation doesn't exists"}
	}
	return as.repos.DeleteOwnRelation(own)
}

// rent
func (as *ApiService) IsRentExists(rent RentModel) bool {
	return as.repos.IsRentExists(rent)
}

func (as *ApiService) GetAppartementRents(id uint) ([]RentModel, apiError.ApiError) {
	//check that it exists
	if !(as.repos.IsAppartementExists(id)) {
		return []RentModel{}, apiError.ValueNotFound{Msg: "Apartment id doesn't exists"}
	}
	return as.repos.GetAppartementRents(id)
}

func (as *ApiService) DeleteRent(rent RentModel) apiError.ApiError {
	//check that it exists
	if !(as.repos.IsRentExists(rent)) {
		return apiError.ValueNotFound{Msg: "this rent doesn't exists"}
	}
	return as.repos.DeleteRent(rent)
}

func (as *ApiService) CreateRent(rent RentModel) apiError.ApiError {
	if as.IsRentExists(rent) {
		return apiError.AlreadyExists{Msg: "Rent already exists"}
	}
	return as.repos.CreateRent(rent)
}
