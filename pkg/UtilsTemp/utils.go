package utils

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"reflect"
	"strconv"
)

// Type is meant to work as an interface into the service base modules to access consumer data to perform some addition support mode
type DataRetriever[T any] func(k string, v string) []T

func DynamicRESTFromTypeStruct[T any](fields map[string]bool, r *gin.Engine, dataSource DataRetriever[T]) {
	var def T

	t := reflect.TypeOf(def)

	if t.Kind() == reflect.Struct {

		for i := 0; i < t.NumField(); i++ {
			f := t.Field(i)

			if f.Type.Name() == "string" && fields[f.Name] {

				restPath := fmt.Sprintf("/%s/:id", f.Name)

				r.GET(restPath, func(ctx *gin.Context) {
					sourceData := dataSource(f.Name, ctx.Param("id"))

					if sourceData != nil {
						ctx.JSON(http.StatusOK, sourceData)
					} else {
						ctx.JSON(http.StatusNotFound, sourceData)
					}

				})
			}
		}

	} else {
		fmt.Println("Type is not struct")
	}
}

// In short I want to pivot my dataset and index by field.
func BuildIndexedDataFromStructByFilter[T any](fields map[string]bool, sourceData *[]T) map[string]map[string][]T {
	defer fmt.Println(" + Utils.BuildIndexedDataFromStructByFilter - Exit")
	fmt.Println(" + Utils.BuildIndexedDataFromStructByFilter - Enter")

	newIndex := make(map[string]map[string][]T)

	for field, fieldFlag := range fields {
		if fieldFlag {

			fmap := make(map[string][]T)

			for i := range *sourceData {
				t := (*sourceData)[i]
				v := reflect.ValueOf(t).FieldByName(field).String()

				if _, ok := fmap[v]; !ok {
					fmap[v] = []T{}
				}

				fmap[v] = append(fmap[v], t)
			}

			newIndex[field] = fmap

			fmt.Println(" - Index Len: ", strconv.Itoa(len(newIndex[field])))
		}
	}

	return newIndex
}
