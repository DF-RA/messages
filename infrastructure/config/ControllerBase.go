package config

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/go-playground/validator/v10"
	"io"
	"net/http"
	"reflect"
)

func GetBody[T any](body io.ReadCloser) (T, error) {
	entity, err := decodeBody[T](body)
	if err != nil {
		return entity, getBodyError(err)
	}

	if err := validateStructByType(entity); err != nil {
		return entity, err
	}

	return entity, nil
}

func validateStructByType(entity interface{}) error {
	if entitiesReflect := reflect.ValueOf(entity); entitiesReflect.Kind() == reflect.Slice {
		for i := 0; i < entitiesReflect.Len(); i++ {
			entityReflect := entitiesReflect.Index(i).Interface()
			if err := validateStruct(entityReflect); err != nil {
				return err
			}
		}
	} else if err := validateStruct(entity); err != nil {
		return err
	}
	return nil
}

func validateStruct(entity interface{}) error {
	if err := validator.New().Struct(entity); err != nil {
		return err
	}
	for i := 0; i < reflect.ValueOf(entity).NumField(); i++ {
		field := reflect.ValueOf(entity).Field(i)
		if field.Kind() == reflect.Struct || field.Kind() == reflect.Slice {
			if err := validateStructByType(field.Interface()); err != nil {
				return err
			}
		}
	}
	return nil
}

func ResponseWithoutBody[T any](writer http.ResponseWriter, body T) {
	writer.Header().Set("Content-Type", "application/json; charset=UTF-8")
	writer.WriteHeader(http.StatusOK)
	json.NewEncoder(writer).Encode(body)
}

func ErrorResponse(writer http.ResponseWriter, httpCode int, err error) {
	writer.Header().Set("Content-Type", "application/json; charset=UTF-8")
	writer.WriteHeader(httpCode)
	json.NewEncoder(writer).Encode(map[string]string{"error": err.Error()})
}

func decodeBody[T any](body io.ReadCloser) (T, error) {
	var output T
	err := json.NewDecoder(body).Decode(&output)
	defer func(Body io.ReadCloser) {
		err = Body.Close()
	}(body)

	return output, err
}

func getBodyError(err error) error {
	if err.Error() == "EOF" {
		return errors.New("Invalid body or empty")
	}
	if unmarshalTypeError := err.(*json.UnmarshalTypeError); unmarshalTypeError != nil {
		return errors.New(fmt.Sprintf("Invalid field %s is not %s value", unmarshalTypeError.Field, unmarshalTypeError.Value))
	}
	return err
}
