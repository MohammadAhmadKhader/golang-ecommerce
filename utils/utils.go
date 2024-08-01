package utils

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"reflect"
	"runtime"
	"strings"

	"github.com/go-playground/validator/v10"
	"github.com/go-sql-driver/mysql"
	"github.com/mohammadahmadkhader/golang-ecommerce/config"
	myDB "github.com/mohammadahmadkhader/golang-ecommerce/db"
)
const GenericErrMessage = "An unexpected error has occurred, please try again later!"
var Validate = validator.New()

func ParseJSON(r *http.Request, payload any) error {
	if r.Body == nil {
		return fmt.Errorf("request body is empty")
	}

	return json.NewDecoder(r.Body).Decode(payload)
}

func WriteJSON(w http.ResponseWriter, status int, v any) error {
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(status)

	return json.NewEncoder(w).Encode(v)
}

func WriteError(w http.ResponseWriter, status int, err error) {
	stackTrace := logCaptureStackTrace()
	returnedErr := err.Error()
	

	errObj := make(map[string]any, 0)
	var mySqlError *mysql.MySQLError
	var validationErrMsg validator.ValidationErrors
	var UnmarshalTypeErr *json.UnmarshalTypeError

	if config.Envs.Env == "production" {
		if status >= 500 {
			returnedErr = GenericErrMessage
			errObj["error"] = GenericErrMessage
			errObj["statusCode"] = status

		} else if errors.As(err, &mySqlError) {
			errObj["error"] = GenericErrMessage
			errObj["statusCode"] = 500
			log.Println(err.Error())

		} else if errors.As(err, &validationErrMsg) {
			errObj["error"] = validationErrMsgHandler(err.(validator.ValidationErrors))
			errObj["statusCode"] = status

		} else if errors.As(err, &UnmarshalTypeErr){
			errObj["error"] = UnmarshalErrMsgHandler(err.(*json.UnmarshalTypeError))
			errObj["statusCode"] = 400
			
		} else {
			errObj["error"] = returnedErr
			errObj["statusCode"] = status
		}

	} else {
		errObj["error"] = returnedErr
		errObj["statusCode"] = status
		errObj["stackTrace"] = stackTrace
	}

	WriteJSON(w, status, errObj)
}

func StartMySqlDB() (*sql.DB, error) {
	dbConfig := mysql.Config{
		User:                 config.Envs.DBUser,
		Passwd:               config.Envs.DBPassword,
		Addr:                 config.Envs.DBAddress,
		DBName:               config.Envs.DBName,
		Net:                  "tcp",
		AllowNativePasswords: true,
		ParseTime:            true,
	}

	db, err := myDB.NewMySqlStart(dbConfig)
	if err != nil {
		return nil, err
	}
	
	return db, nil
}

func logCaptureStackTrace() string {
	var filteredStackTraceSlice []string

	buf := make([]byte, 1<<16)
	n := runtime.Stack(buf, false)
	var fullSrackTrace = string(buf[:n])

	lines := strings.Split(fullSrackTrace, "\n")

	for _, line := range lines {
		if strings.Contains(line, "ecommerceAPI") {
			filteredStackTraceSlice = append(filteredStackTraceSlice, line)
		}
		if strings.Contains(line, "golang-ecommerce") {
			filteredStackTraceSlice = append(filteredStackTraceSlice, line)
		}
	}

	filteredStackTraceAsString := strings.Join(filteredStackTraceSlice, "\n")
	log.Printf("\nstack Trace: %v\n", filteredStackTraceAsString)

	return filteredStackTraceAsString
}

// returns first validation error, it's responsible for handling validator validation errors.
func validationErrMsgHandler(errors validator.ValidationErrors) string {
		var message string;
		
		switch errors[0].Tag() {
			case "required":
				message = fmt.Sprintf("%s is required", errors[0].Field())
			case "email":
				message = "invalid email"
			case "gte":
				message = fmt.Sprintf("%s must be greater than or equal to %s", errors[0].Field(), errors[0].Param())
			case "lte":
				message = fmt.Sprintf("%s must be less than or equal to %s", errors[0].Field(), errors[0].Param())
			case "gt":
				message = fmt.Sprintf("%s must be greater than %s", errors[0].Field(), errors[0].Param())
			case "lt":
				message = fmt.Sprintf("%s must be less than %s", errors[0].Field(), errors[0].Param())
			case "min":
				if(errors[0].Kind() == reflect.String){
					message = fmt.Sprintf("%s minimum length allowed is %s", errors[0].Field(), errors[0].Param())
				} else {
					message = fmt.Sprintf("%s minimum allowed is %s", errors[0].Field(), errors[0].Param())
				}
			case "max":
				if(errors[0].Kind() == reflect.String){
					message = fmt.Sprintf("%s maximum length allowed is %s", errors[0].Field(), errors[0].Param())
				} else {
					message = fmt.Sprintf("%s maximum allowed is %s", errors[0].Field(), errors[0].Param())
				}
			default:
				message = fmt.Sprintf("%s is invalid", errors[0].Field())
		}
		
		return message
}

// handles invalid json parsing type error and returns error for production.
func UnmarshalErrMsgHandler(error *json.UnmarshalTypeError) string {
	errMsg := fmt.Sprintf("%v is type %v can't be equal to %v",error.Field,error.Type,error.Value)
	return errMsg
}