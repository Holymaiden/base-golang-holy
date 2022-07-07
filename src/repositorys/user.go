package repositorys

import (
	"fmt"
	"jwt/src/helpers"
	"jwt/src/models"
	"jwt/src/requests"
	"math"
	"strings"

	"gorm.io/gorm"
)

type UserRepository interface {
	All() []models.User
	Find(id uint64) models.User
	Create(model models.User) models.User
	Update(model models.User) models.User
	Delete(model models.User) models.User
	Pagination(*requests.Pagination) (Result, int)
}

type userConnection struct {
	connection *gorm.DB
}

func NewUserRepository(db *gorm.DB) UserRepository {
	return &userConnection{
		connection: db,
	}
}

func (u *userConnection) All() []models.User {
	var users []models.User
	u.connection.Find(&users)
	return users
}

func (u *userConnection) Find(id uint64) models.User {
	var user models.User
	u.connection.First(&user, id)
	return user
}

func (u *userConnection) Create(model models.User) models.User {
	model.Password = helpers.HashAndSalt([]byte(model.Password))
	u.connection.Create(&model)
	return model
}

func (u *userConnection) Update(model models.User) models.User {
	u.connection.Save(&model)
	return model
}

func (u *userConnection) Delete(model models.User) models.User {
	u.connection.Delete(&model)
	return model
}

func (u *userConnection) Pagination(p *requests.Pagination) (Result, int) {
	var records []models.User
	var totalRows int64

	totalRows, totalPages, fromRow, toRow := 0, 0, 0, 0

	fmt.Println("pagination.Limit: ", p.Limit)
	fmt.Println("pagination.Page: ", p.Page)

	offset := (p.Page - 1) * p.Limit

	// get data with limit, offset & order
	find := u.connection.Limit(p.Limit).Offset(offset).Order(p.Sort)

	// generate where query
	searchs := p.Searchs

	if len(searchs) > 0 {
		for _, value := range searchs {
			column := value.Column
			action := value.Action
			query := value.Query

			switch action {
			case "equals":
				whereQuery := fmt.Sprintf("%s = ?", column)
				find = find.Where(whereQuery, query)
			case "contains":
				whereQuery := fmt.Sprintf("%s LIKE ?", column)
				find = find.Where(whereQuery, "%"+query+"%")
			case "in":
				whereQuery := fmt.Sprintf("%s IN (?)", column)
				queryArray := strings.Split(query, ",")
				find = find.Where(whereQuery, queryArray)
			}
		}
	}

	find = find.Find(&records)

	// has error find data
	errFind := find.Error

	if errFind != nil {
		return Result{Error: errFind}, totalPages
	}

	p.Rows = records

	errCount := u.connection.Model(&models.User{}).Count(&totalRows).Error

	if errCount != nil {
		return Result{Error: errCount}, totalPages
	}

	p.TotalRows = totalRows

	// calculate total pages
	totalPages = int(math.Ceil(float64(totalRows) / float64(p.Limit)))

	if p.Page == 1 || p.Page == 0 {
		// set from & to row on first page
		fromRow = 1
		toRow = p.Limit
	} else {
		if p.Page <= totalPages {
			// calculate from & to row
			fromRow = ((p.Page - 1) * p.Limit) + 1
			toRow = p.Page * p.Limit
		}
	}

	if int64(toRow) > totalRows {
		// set to row with total rows
		toRow = int(totalRows)
	}

	p.FromRow = fromRow
	p.ToRow = toRow

	return Result{Result: p}, totalPages
}
