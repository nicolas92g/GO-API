package REST

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"net/http"
	"projet/api/commons"
	"projet/api/commons/apiError"
	"projet/api/commons/middlewares"
	"strconv"
)

type ApiController struct {
	service ApiService
	user    UserModel
	req     commons.RequestApi
	res     commons.ResponseApi
}

func (c *ApiController) checkAdminRights() bool {
	if !(c.user.Admin) {
		c.res.SendMessage("Your need to be an admin for doing this", http.StatusForbidden)
	}
	return !c.user.Admin
}

func (c *ApiController) checkMethod(method string) bool {
	if c.req.GetMethod() != method {
		c.res.SendError(apiError.BadMethod{Msg: "need to use " + method + " method for this request"})
		return true
	}
	return false
}

func (c *ApiController) GetCurrentUser() bool {
	ApiKey := c.req.GetParameters().Get("api_key")
	if ApiKey == "" {
		c.res.SendMessage("There is no api_key", http.StatusBadRequest)
		return false
	}
	var err apiError.ApiError
	c.user, err = c.service.GetCurrentUser(ApiKey)
	return !(c.res.SendError(err))
}

func (c *ApiController) getIntParameter(name string, defaultValue int) int {
	ret, err := strconv.Atoi(c.req.GetParameters().Get(name))
	if err != nil {
		ret = defaultValue
	}
	return ret
}

func (c *ApiController) Dispatch(r *http.Request, w http.ResponseWriter) {
	c.req.Create(r)
	c.res.Create(w)

	if c.res.SendError(c.service.Connect()) {
		return
	}
	defer c.service.Close()
	if !(c.GetCurrentUser()) {
		return
	}

	route, err := c.req.GetPathAt(1)
	if c.res.SendError(err) {
		return
	}

	switch route {
	case "user":
		c.res.SendError(c.userDispatch())
	case "flats":
		c.res.SendError(c.flatsDispatch())
	case "rents":
		c.res.SendError(c.rentsDispatch())
	case "owns":
		c.res.SendError(c.ownDispatch())
	default:
		c.res.SendMessage(route+" is not a valid route", http.StatusNotFound)
	}
}

func (c *ApiController) userDispatch() apiError.ApiError {
	route, err := c.req.GetPathAt(2)
	if err != nil {
		return err
	}
	switch route {
	case "list":
		c.GetUsers()
	case "create":
		c.CreateUser()
	case "delete":
		c.DeleteUser()
	case "owns":
		c.GetOwns()
	case "rents":
		c.GetRents()
	default:
		return apiError.BadURL{Msg: route + " is not a valid route"}
	}
	return nil
}

func (c *ApiController) GetUsers() {

	if c.checkAdminRights() {
		return
	}
	if c.checkMethod("GET") {
		return
	}

	limit := c.getIntParameter("limit", 100)

	users, err := c.service.GetUsers(uint(limit))
	c.res.SendError(err)

	c.res.SendContent(users, http.StatusOK)
}

func (c *ApiController) CreateUser() {

	if c.checkAdminRights() {
		return
	}
	if c.checkMethod("POST") {
		return
	}

	var user UserModel

	for {
		bytes := make([]byte, 20)
		_, _ = rand.Read(bytes)
		user.Api_key = hex.EncodeToString(bytes)

		if !(c.service.IsUserKeyExists(user.Api_key)) {
			break
		}
	}

	if c.res.SendError(middlewares.JsonMiddleware(&c.req, &user)) {
		return
	}
	if c.res.SendError(c.service.CreateUser(user)) {
		return
	}

	var userOrAdmin string
	if user.Admin {
		userOrAdmin = "admin"
	} else {
		userOrAdmin = "user"
	}
	c.res.SendMessage("this new "+userOrAdmin+" key is : "+user.Api_key, http.StatusOK)
}

func (c *ApiController) DeleteUser() {
	if c.checkAdminRights() {
		return
	}
	if c.checkMethod("GET") {
		return
	}

	id := c.getIntParameter("id", -1)
	if id == -1 {
		c.res.SendMessage("need to add id parameter", http.StatusBadRequest)
		return
	}

	if c.res.SendError(c.service.DeleteUser(uint(id))) {
		return
	}
	c.res.SendMessage("user was deleted succesfully", http.StatusOK)
}

func (c *ApiController) GetOwns() {
	if c.checkMethod("GET") {
		return
	}

	id := c.getIntParameter("id", -1)
	if id == -1 {
		id = int(c.user.Id)
	}
	//has to be admin to see other user
	if id != int(c.user.Id) && c.checkAdminRights() {
		return
	}

	owns, err := c.service.GetUserOwns(uint(id))

	if c.res.SendError(err) {
		return
	}
	c.res.SendContent(owns, http.StatusOK)
}

func (c *ApiController) GetRents() {
	if c.checkMethod("GET") {
		return
	}

	id := c.getIntParameter("id", -1)
	if id == -1 {
		id = int(c.user.Id)
	}
	//has to be admin to see other user
	if id != int(c.user.Id) && c.checkAdminRights() {
		return
	}

	rents, err := c.service.GetUserRents(uint(id))

	if c.res.SendError(err) {
		return
	}
	c.res.SendContent(rents, http.StatusOK)
}

func (c *ApiController) flatsDispatch() apiError.ApiError {
	route, err := c.req.GetPathAt(2)
	if err != nil {
		return err
	}
	switch route {
	case "list":
		c.GetAppartements()
	case "create":
		c.CreateAppartement()
	case "delete":
		c.DeleteAppartement()
	case "update":
		c.UpdateAppartement()
	case "owners":
		c.GetOwners()
	case "rents":
		c.GetFlatsRents()
	default:
		return apiError.BadURL{Msg: route + " is not a valid route"}
	}
	return nil
}

func (c *ApiController) GetAppartements() {
	if c.checkMethod("GET") {
		return
	}

	limit := c.getIntParameter("limit", 100)

	users, err := c.service.GetAppartements(uint(limit))
	c.res.SendError(err)

	c.res.SendContent(users, http.StatusOK)
}

func (c *ApiController) CreateAppartement() {
	if c.checkMethod("POST") {
		return
	}

	var flat AppartementModel

	if c.res.SendError(middlewares.JsonMiddleware(&c.req, &flat)) {
		return
	}
	if c.res.SendError(c.service.CreateAppartement(flat)) {
		return
	}

	c.res.SendMessage("Your flat was created successfully", http.StatusOK)
}

func (c *ApiController) DeleteAppartement() {
	if c.checkAdminRights() {
		return
	}
	if c.checkMethod("GET") {
		return
	}

	id := c.getIntParameter("id", -1)
	if id == -1 {
		c.res.SendMessage("need to add id parameter", http.StatusBadRequest)
		return
	}

	if c.res.SendError(c.service.DeleteAppartement(uint(id))) {
		return
	}
	c.res.SendMessage("flat was deleted successfully", http.StatusOK)
}

func (c *ApiController) UpdateAppartement() {
	if c.checkMethod("POST") {
		return
	}

	id := c.getIntParameter("id", -1)
	if id == -1 {
		c.res.SendMessage("need to add id parameter", http.StatusBadRequest)
		return
	}

	var flat AppartementModel

	if c.res.SendError(middlewares.JsonMiddleware(&c.req, &flat)) {
		return
	}
	flat.Id = uint(id)
	if c.res.SendError(c.service.UpdateAppartement(flat)) {
		return
	}

	c.res.SendMessage("Your flat was updated successfully", http.StatusOK)
}

func (c *ApiController) GetOwners() {
	if c.checkAdminRights() {
		return
	}
	if c.checkMethod("GET") {
		return
	}

	id := c.getIntParameter("id", -1)
	if id == -1 {
		c.res.SendMessage("need to add id parameter", http.StatusBadRequest)
		return
	}

	users, err := c.service.GetAppartementOwners(uint(id))
	if c.res.SendError(err) {
		return
	}

	c.res.SendContent(users, http.StatusOK)
}

func (c *ApiController) GetFlatsRents() {
	if c.checkMethod("GET") {
		return
	}

	id := c.getIntParameter("id", -1)
	if id == -1 {
		c.res.SendMessage("need to add id parameter", http.StatusBadRequest)
		return
	}

	users, err := c.service.GetAppartementRents(uint(id))
	if c.res.SendError(err) {
		return
	}

	c.res.SendContent(users, http.StatusOK)
}

func (c *ApiController) rentsDispatch() apiError.ApiError {
	route, err := c.req.GetPathAt(2)
	if err != nil {
		return err
	}
	switch route {
	case "create":
		c.CreateRent()
	case "delete":
		c.DeleteRent()
	default:
		return apiError.BadURL{Msg: route + " is not a valid route"}
	}
	return nil
}

func (c *ApiController) CreateRent() {
	if c.checkMethod("POST") {
		return
	}

	var rent RentModel

	if c.res.SendError(middlewares.JsonMiddleware(&c.req, &rent)) {
		return
	}
	if c.res.SendError(c.service.CreateRent(rent)) {
		return
	}

	c.res.SendMessage("Your rent was created successfully", http.StatusOK)
}

func (c *ApiController) DeleteRent() {

	if c.checkMethod("GET") {
		return
	}

	var rent RentModel

	rent.User_id = c.getIntParameter("user_id", -1)
	if rent.User_id == -1 {
		c.res.SendMessage("need to add user_id parameter", http.StatusBadRequest)
		return
	}
	rent.Appartement_id = c.getIntParameter("flat_id", -1)
	if rent.Appartement_id == -1 {
		c.res.SendMessage("need to add flat_id parameter", http.StatusBadRequest)
		return
	}

	if c.res.SendError(c.service.DeleteRent(rent)) {
		return
	}
	c.res.SendMessage("rent was deleted successfully", http.StatusOK)
}

func (c *ApiController) ownDispatch() apiError.ApiError {
	route, err := c.req.GetPathAt(2)
	if err != nil {
		return err
	}
	switch route {
	case "create":
		c.CreateOwn()
	case "delete":
		c.DeleteOwn()
	default:
		return apiError.BadURL{Msg: route + " is not a valid route"}
	}
	return nil
}

func (c *ApiController) CreateOwn() {
	if c.checkMethod("POST") {
		return
	}

	var own OwnModel

	if c.res.SendError(middlewares.JsonMiddleware(&c.req, &own)) {
		return
	}
	fmt.Print(own)
	if c.res.SendError(c.service.CreateOwnRelation(own)) {
		return
	}

	c.res.SendMessage("Your rent was created successfully", http.StatusOK)
}

func (c *ApiController) DeleteOwn() {
	if c.checkAdminRights() {
		return
	}
	if c.checkMethod("GET") {
		return
	}

	var own OwnModel

	user_id := c.getIntParameter("user_id", -1)
	if user_id == -1 {
		c.res.SendMessage("need to add user_id parameter", http.StatusBadRequest)
		return
	}
	flat_id := c.getIntParameter("flat_id", -1)
	if flat_id == -1 {
		c.res.SendMessage("need to add flat_id parameter", http.StatusBadRequest)
		return
	}

	own.User_id = uint(user_id)
	own.Appartement_id = uint(flat_id)

	if c.res.SendError(c.service.DeleteOwnRelation(own)) {
		return
	}
	c.res.SendMessage("own was deleted successfully", http.StatusOK)
}
