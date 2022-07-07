package helpers

import (
	"jwt/src/requests"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

func Pagination(context *gin.Context) *requests.Pagination {
	limit := 5
	page := 1
	sort := "created_at asc"

	var searchs []requests.Search

	query := context.Request.URL.Query()

	for key, value := range query {
		queryValue := value[len(value)-1]

		switch key {
		case "limit":
			limit, _ = strconv.Atoi(queryValue)
		case "page":
			page, _ = strconv.Atoi(queryValue)
		case "sort":
			sort = queryValue
		}

		// check if query parameter key contains dot
		if strings.Contains(key, ".") {
			// split query parameter key by dot
			searchKeys := strings.Split(key, ".")

			// create search object
			search := requests.Search{Column: searchKeys[0], Action: searchKeys[1], Query: queryValue}

			// add search object to searchs array
			searchs = append(searchs, search)
		}
	}

	return &requests.Pagination{Limit: limit, Page: page, Sort: sort, Searchs: searchs}
}
