package REST

type AppartementModel struct {
	Id            uint   `json:"id"`
	Area          uint   `json:"area"`
	Capacity      uint   `json:"capacity"`
	StreetNumber  uint   `json:"street_number"`
	StreetName    string `json:"street_name"`
	City          string `json:"city"`
	Disponibility bool   `json:"disponibility"`
}

type UserModel struct {
	Id      uint   `json:"id"`
	Admin   bool   `json:"admin"`
	Api_key string `json:"api_key"`
}

type RentModel struct {
	User_id        int     `json:"user_id"`
	Appartement_id int     `json:"flat_id"`
	Date_begin     string  `json:"date_begin"`
	Date_end       string  `json:"date_end"`
	Price          float64 `json:"price"`
}

type OwnModel struct {
	User_id        uint `json:"user_id"`
	Appartement_id uint `json:"appartement_id"`
}
