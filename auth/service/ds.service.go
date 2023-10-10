package service

import (
	"auth/structz"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/jinzhu/gorm"
	"github.com/joho/godotenv"

	//db "dsGrid-2/database"

	"github.com/gofiber/fiber/v2"
)

func ResponseObj(w http.ResponseWriter, err interface{}, code int) {
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(err)
}

func Response(context *fiber.Ctx, msg interface{}, success bool, data interface{}) error {
	return context.Status(fiber.StatusOK).JSON(ResponseObj)
}

func DsGrid(context *fiber.Ctx) error {
	request := &structz.GridRequest{}
	var gridResponse structz.GridResponse
	var rowData []map[string]interface{} = make([]map[string]interface{}, 0)

	err := context.BodyParser(request)

	if err != nil {
		return Response(context, "Утгуудын зөв эсэхийг шалгана уу", false, "")
	}

	// println(request.Code)
	// println(request.StartRow)
	// println(request.EndRow)

	// var name string = "config"
	// println(name)
	tableName := "vw_user_status"
	query := "select * from " + tableName
	query += whereSql(request) + OrderBySql(request)
	query += " limit 10"
	fmt.Println(query)
	//

	er := godotenv.Load()

	if er != nil {
		log.Fatal("Error loading .env file")
	}
	username := os.Getenv("databaseUser")
	password := os.Getenv("databasePassword")
	databaseName := os.Getenv("databaseName")
	databaseHost := os.Getenv("databaseHost")

	dbURI := fmt.Sprintf("host=%s user=%s dbname=%s sslmode=disable password=%s", databaseHost, username, databaseName, password)

	conn, err := gorm.Open("postgres", dbURI)

	if err != nil {
		fmt.Println("error", err)
		panic(err)
	}
	fmt.Println("connection succesfully established", conn)
	//

	conn.Scopes(Paginate(&gridResponse, context, conn))
	conn.Find(&rowData)
	gridResponse.Data = rowData

	return context.Status(fiber.StatusOK).JSON(gridResponse)
	//return gridresponse in GridResponse struct
}

func OrderBySql(request *structz.GridRequest) string {
	if request.SortModel == nil {
		return ""
	}
	order := " order by "
	for index, data := range request.SortModel {
		if index != 0 {
			order += ", "
		}
		order += data.ColId + " " + data.Sort
	}
	return order
	//"return "order by , + ,"
}

func whereSql(request *structz.GridRequest) string {
	var whereArray []string

	for key := range request.FilterModel {
		result := GetFilters(request.FilterModel[key], key)
		whereArray = append(whereArray, result)
	}
	if len(whereArray) == 0 {
		return ""
	}
	where := " where "
	for i, data := range whereArray {
		if i != 0 {
			where += " and "
		}
		where += data
	}
	return where
	// return where filtertype + and ft + and ...
}

func GetGroupColumns(request *structz.GridRequest) string {
	return ""
}

func GetFilters(filterModel structz.FilterModel, key string) string {
	if filterModel.FilterType == "text" {
		return TextFilter(filterModel, key)
	}
	if filterModel.FilterType == "number" {
		return NumberFilter(filterModel, key)
	}
	if filterModel.FilterType == "date" {
		return DateFilter(filterModel, key)
	}
	if filterModel.FilterType == "set" {
		return SetFilter(filterModel, key)
	}
	return ""
	// return filtertype(text) + key ""
}

func SetFilter(filterModel structz.FilterModel, columnName string) string {
	return columnName + " IN " + " (" + filterModel.Filter + ") "
	// column name + IN + json.filter
}

func TextFilter(filterModel structz.FilterModel, columnName string) string {
	if filterModel.Type == "equals" {
		return "LOWER(" + columnName + ") = LOWER('" + filterModel.Filter + "')"
	} else if filterModel.Type == "notEqual" {
		return "LOWER(" + columnName + ") <> LOWER('" + filterModel.Filter + "')"
	} else if filterModel.Type == "startsWith" {
		return "LOWER(" + columnName + ") LIKE LOWER('" + filterModel.Filter + "')"
	} else if filterModel.Type == "endsWith" {
		return "LOWER(" + columnName + ") LIKE LOWER('" + filterModel.Filter + "')"
	} else if filterModel.Type == "contains" {
		return "LOWER(" + columnName + ") LIKE LOWER('" + filterModel.Filter + "')"
	} else if filterModel.Type == "notContains" {
		return "LOWER(" + columnName + ") NOT LIKE LOWER('" + filterModel.Filter + "')"
	}

	return ""
	// LOWER(columenName sign[=, <>, Like, Not LIKE] .Filter[datatype])
}

func NumberFilter(filterModel structz.FilterModel, columnName string) string {
	if filterModel.Type == "inRange" {
		return columnName + " BETWEEN " + filterModel.Filter + " AND " + filterModel.Filter + " "
	}

	return columnName + " " + OperatorMap(filterModel.Type) + " " + filterModel.Filter + " "
	// return columnName Between .Filter AND .Filter
	// return columnName (=, <>, <,>...) .Filter
}

func DateFilter(filterModel structz.FilterModel, columnName string) string {
	if filterModel.FilterType == "inRange" {
		return columnName + " BETWEEN '" + filterModel.FilterType + "' AND '" + filterModel.FilterType + "' "
	}
	return columnName + " " + OperatorMap(filterModel.FilterType) + " '" + filterModel.FilterType + "' "
	// return columnName BETWEEN .FilterType AND FilterType
	// return columnName = .FilterType
}

func OperatorMap(key string) string {
	if key == "equals" {
		return "="
	} else if key == "notEqual" {
		return "<>"
	} else if key == "startsWith" {
		return "<"
	} else if key == "endsWith" {
		return "<="
	} else if key == "contains" {
		return ">"
	} else if key == "notContains" {
		return ">="
	}
	return "="
	// return =, <>, = etc
}

func Paginate(response *structz.GridResponse, c *fiber.Ctx, db *gorm.DB) func(db *gorm.DB) *gorm.DB {
	totalRows := int(db.RowsAffected)

	pageNumber, err := strconv.Atoi(c.Query("page", "0"))
	if err != nil {
		pageNumber = 0
	}
	pageSize, err := strconv.Atoi(c.Query("size", ""))
	if err != nil {
		pageSize = 10
	}
	if pageSize == 0 {
		pageSize = totalRows
	}

	offset := pageNumber * pageSize
	fmt.Println(totalRows, "  ", pageNumber, "   ", pageSize, "  ", "sorted", "   ")
	response.FillPaginate(totalRows, pageNumber, pageSize)

	return func(db *gorm.DB) *gorm.DB {
		return db.Offset(offset).Limit(pageSize)
	}
	//printlin - rows(int) pageNumber(int) pageSize(int) sorted
	// return database data with limit of pageSize
}
