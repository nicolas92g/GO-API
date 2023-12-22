package REST

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"projet/api/REST/vars"
	"projet/api/commons/apiError"
)

type ApiRepository struct {
	db *sql.DB
}

func (ar *ApiRepository) Connect() apiError.ApiError {
	var err error
	ar.db, err = sql.Open("mysql", vars.BDD_LOGIN+":"+vars.BDD_PASSWORD+"@tcp(127.0.0.1:3306)/apirestpatrop")
	if err != nil {
		return apiError.BddError{Msg: "Failed to connect to the database : " + err.Error()}
	}
	return nil
}

func (ar *ApiRepository) Close() {
	_ = ar.db.Close()
}

//USER TABLE :

func (ar *ApiRepository) IsUserKeyExists(api_key string) bool {
	//get user existence
	var ret bool
	err := ar.db.QueryRow("SELECT EXISTS(SELECT user_id FROM appuser WHERE api_key=?)", api_key).Scan(&ret)
	if err != nil {
		return false
	}
	return ret
}

func (ar *ApiRepository) IsUserIdExists(id uint) bool {
	//get user existence
	var ret bool
	err := ar.db.QueryRow("SELECT EXISTS(SELECT user_id FROM appuser WHERE user_id=?)", id).Scan(&ret)
	if err != nil {
		return false
	}
	return ret
}

func (ar *ApiRepository) GetCurrentUser(api_key string) (UserModel, apiError.ApiError) {
	var ret UserModel

	err := ar.db.QueryRow("SELECT admin, user_id, api_key FROM appuser WHERE api_key=?", api_key).Scan(&ret.Admin, &ret.Id, &ret.Api_key)
	if err != nil {
		return ret, apiError.BddError{Msg: "Failed to get current user : " + err.Error()}
	}

	return ret, nil
}

func (ar *ApiRepository) GetUsers(limit uint) ([]UserModel, apiError.ApiError) {
	ret := make([]UserModel, 0)

	rows, err := ar.db.Query("SELECT user_id, admin, api_key FROM appuser LIMIT ?", limit)
	if err != nil {
		return nil, apiError.BddError{Msg: "failed to get users" + err.Error()}
	}

	for rows.Next() {
		var a UserModel
		err2 := rows.Scan(&a.Id, &a.Admin, &a.Api_key)
		if err2 != nil {
			return nil, apiError.BddError{Msg: "failed to read one user" + err2.Error()}
		}
		ret = append(ret, a)
	}
	return ret, nil
}

func (ar *ApiRepository) CreateUser(u UserModel) apiError.ApiError {
	_, err := ar.db.Exec(
		"INSERT INTO appuser(admin, api_key) VALUES(?, ?)",
		u.Admin,
		u.Api_key,
	)

	if err != nil {
		return apiError.BddError{Msg: "Failed to write user in db : " + err.Error()}
	}

	return nil
}

func (ar *ApiRepository) DeleteUser(id uint) apiError.ApiError {
	_, err := ar.db.Exec("DELETE FROM appuser WHERE user_id = ?", id)

	if err != nil {
		return apiError.BddError{Msg: "Failed to delete user in db : " + err.Error()}
	}

	return nil
}

func (ar *ApiRepository) GetUserRents(id uint) ([]RentModel, apiError.ApiError) {
	ret := make([]RentModel, 0)

	rows, err := ar.db.Query("SELECT appartement_id, user_id, date_begin, date_end, price FROM rent WHERE user_id = ?", id)
	if err != nil {
		return ret, apiError.BddError{Msg: "Failed to get user rents in the db : " + err.Error()}
	}
	for rows.Next() {
		var v RentModel
		err2 := rows.Scan(&v.Appartement_id, &v.User_id, &v.Date_begin, &v.Date_end, &v.Price)
		if err2 != nil {
			return ret, apiError.BddError{Msg: "Failed to scan one rent : " + err2.Error()}
		}
		ret = append(ret, v)
	}

	return ret, nil
}

//APPARTEMENT TABLE

func (ar *ApiRepository) IsAppartementExists(id uint) bool {
	var ret bool
	err := ar.db.QueryRow("SELECT EXISTS(SELECT area FROM appartement WHERE appartement_id = ?)", id).Scan(&ret)

	if err != nil {
		return false
	}

	return ret
}

func (ar *ApiRepository) GetAppartements(limit uint) ([]AppartementModel, apiError.ApiError) {
	ret := make([]AppartementModel, 0)

	rows, err := ar.db.Query("SELECT appartement_id, area, capacity, streetNumber, streetName, city, disponibility FROM appartement LIMIT ?", limit)
	if err != nil {
		return nil, apiError.BddError{Msg: "failed to get apartments" + err.Error()}
	}

	for rows.Next() {
		var a AppartementModel
		err2 := rows.Scan(&a.Id, &a.Area, &a.Capacity, &a.StreetNumber, &a.StreetName, &a.City, &a.Disponibility)
		if err2 != nil {
			return nil, apiError.BddError{Msg: "failed to read one apartment" + err2.Error()}
		}
		ret = append(ret, a)
	}
	return ret, nil
}

func (ar *ApiRepository) GetAppartement(id uint) (AppartementModel, apiError.ApiError) {
	var ret AppartementModel

	//get apartment
	err := ar.db.QueryRow(
		"SELECT appartement_id, area, capacity, streetNumber, streetName, city, disponibility FROM appartement WHERE appartement_id = ?",
		id).Scan(&ret.Id, &ret.Area, &ret.Capacity, &ret.StreetNumber, &ret.StreetName, &ret.City, &ret.Disponibility)
	if err != nil {
		return ret, apiError.BddError{Msg: "Failed to read appartement : " + err.Error()}
	}
	return ret, nil
}

func (ar *ApiRepository) CreateAppartement(a AppartementModel) apiError.ApiError {
	_, err := ar.db.Exec(
		"INSERT INTO appartement(area, capacity, streetNumber, streetName, city, disponibility) VALUES( ?, ?, ?, ?, ?, ?)",
		a.Area,
		a.Capacity,
		a.StreetNumber,
		a.StreetName,
		a.City,
		a.Disponibility,
	)

	if err != nil {
		return apiError.BddError{Msg: "Failed to write apartment in db : " + err.Error()}
	}

	return nil
}

func (ar *ApiRepository) DeleteAppartement(id uint) apiError.ApiError {
	_, err := ar.db.Exec("DELETE FROM appartement WHERE appartement_id = ?", id)

	if err != nil {
		return apiError.BddError{Msg: "Failed to delete apartment in db : " + err.Error()}
	}

	return nil
}

func (ar *ApiRepository) UpdateAppartement(flat AppartementModel) apiError.ApiError {
	_, err := ar.db.Exec(
		"UPDATE appartement SET area = ?, capacity = ?, streetNumber = ?, streetName = ?, city = ?, disponibility = ? WHERE appartement_id = ?",
		flat.Area,
		flat.Capacity,
		flat.StreetNumber,
		flat.StreetName,
		flat.City,
		flat.Disponibility,
		flat.Id,
	)

	if err != nil {
		return apiError.BddError{Msg: "Failed to write apartment in db : " + err.Error()}
	}

	return nil
}

//RENT TABLE

func (ar *ApiRepository) IsRentExists(rent RentModel) bool {
	var ret bool
	err := ar.db.QueryRow("SELECT EXISTS(SELECT user_id FROM rent WHERE user_id=? AND appartement_id=?)", rent.User_id, rent.Appartement_id).Scan(&ret)
	if err != nil {
		return false
	}
	return ret
}

func (ar *ApiRepository) GetAppartementRents(id uint) ([]RentModel, apiError.ApiError) {
	ret := make([]RentModel, 0)

	rows, err := ar.db.Query("SELECT appartement_id, user_id, date_begin, date_end, price FROM rent WHERE appartement_id = ?", id)
	if err != nil {
		return ret, apiError.BddError{Msg: "Failed to get appartement rents in the db : " + err.Error()}
	}
	for rows.Next() {
		var v RentModel
		err2 := rows.Scan(&v.Appartement_id, &v.User_id, &v.Date_begin, &v.Date_end, &v.Price)

		if err2 != nil {
			return ret, apiError.BddError{Msg: "Failed to scan one rent : " + err2.Error()}
		}
		ret = append(ret, v)
	}

	return ret, nil
}

func (ar *ApiRepository) DeleteRent(rent RentModel) apiError.ApiError {
	_, err := ar.db.Exec("DELETE FROM rent WHERE user_id = ? AND appartement_id=?", rent.User_id, rent.Appartement_id)

	if err != nil {
		return apiError.BddError{Msg: "Failed to delete rent in db : " + err.Error()}
	}

	return nil
}

func (ar *ApiRepository) CreateRent(rent RentModel) apiError.ApiError {
	_, err := ar.db.Exec(
		"INSERT INTO rent(appartement_id, user_id, date_begin, date_end, price) VALUES(?, ?, ?, ?, ?)",
		rent.Appartement_id,
		rent.User_id,
		rent.Date_begin,
		rent.Date_end,
		rent.Price,
	)

	if err != nil {
		return apiError.BddError{Msg: "Failed to write rent in db : " + err.Error()}
	}

	return nil
}

//OWN TABLE

func (ar *ApiRepository) IsOwnRelationExists(own OwnModel) bool {
	var ret bool
	err := ar.db.QueryRow("SELECT EXISTS(SELECT user_id FROM own WHERE user_id=? AND appartement_id=?)", own.User_id, own.Appartement_id).Scan(&ret)
	if err != nil {
		return false
	}
	return ret
}

func (ar *ApiRepository) CreateOwnRelation(own OwnModel) apiError.ApiError {
	_, err := ar.db.Exec(
		"INSERT INTO own(appartement_id, user_id) VALUES(?, ?)",
		own.Appartement_id,
		own.User_id,
	)

	if err != nil {
		return apiError.BddError{Msg: "Failed to write own in db : " + err.Error()}
	}

	return nil
}

func (ar *ApiRepository) DeleteOwnRelation(own OwnModel) apiError.ApiError {
	_, err := ar.db.Exec("DELETE FROM own WHERE user_id = ? AND appartement_id=?", own.User_id, own.Appartement_id)

	if err != nil {
		return apiError.BddError{Msg: "Failed to delete rent in db : " + err.Error()}
	}

	return nil
}

func (ar *ApiRepository) GetUserOwns(user uint) ([]OwnModel, apiError.ApiError) {

	var ret []OwnModel = make([]OwnModel, 0)

	rows, err := ar.db.Query("SELECT appartement_id FROM own WHERE user_id=?", user)

	if err != nil {
		return ret, apiError.BddError{Msg: "Failed to get user owns : " + err.Error()}
	}
	for rows.Next() {
		var element OwnModel = OwnModel{User_id: user}
		err2 := rows.Scan(&element.Appartement_id)
		if err2 != nil {
			return ret, apiError.BddError{Msg: "Failed to get user own : " + err2.Error()}
		}
		ret = append(ret, element)
	}
	return ret, nil
}

func (ar *ApiRepository) GetAppartementOwners(appartement uint) ([]OwnModel, apiError.ApiError) {

	var ret []OwnModel = make([]OwnModel, 0)

	rows, err := ar.db.Query("SELECT user_id FROM own WHERE appartement_id=?", appartement)

	if err != nil {
		return ret, apiError.BddError{Msg: "Failed to get appartement owners : " + err.Error()}
	}
	for rows.Next() {
		var element = OwnModel{Appartement_id: appartement}
		err2 := rows.Scan(&element.User_id)
		if err2 != nil {
			return ret, apiError.BddError{Msg: "Failed to get appartement owner : " + err2.Error()}
		}
		ret = append(ret, element)
	}
	return ret, nil
}
