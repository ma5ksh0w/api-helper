package helper

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"reflect"
	"strconv"
	"strings"
)

func ParseJSON(r *http.Request, dst interface{}) error {
	if r.Body == nil {
		return io.ErrUnexpectedEOF
	}

	defer r.Body.Close()

	data, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return err
	}

	return json.Unmarshal(data, dst)
}

func ParseJSONMulti(r *http.Request, dsts ...interface{}) error {
	if r.Body == nil {
		return io.ErrUnexpectedEOF
	}

	defer r.Body.Close()

	data, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return err
	}

	for _, dst := range dsts {
		if err := json.Unmarshal(data, dst); err != nil {
			return err
		}
	}

	return nil
}

func ReadAuthToken(r *http.Request) (string, error) {
	auth := r.Header.Get("Authorization")
	if auth == "" {
		return "", ErrAuthFailed
	}

	if !strings.HasPrefix(auth, "Bearer ") {
		return "", ErrInvalidArgument
	}

	return strings.TrimPrefix(auth, "Bearer "), nil
}

func WriteReadAuthTokenError(rw http.ResponseWriter, err error) error {
	switch err {
	case ErrInvalidArgument:
		return WriteError(rw, 0, http.StatusBadRequest, "invalid token")

	case ErrAuthFailed:
		return WriteError(rw, 1, http.StatusForbidden, "invalid token")

	default:
		return WriteError(rw, -1, 0, "invalid token")
	}
}

func ParseVars(vars map[string]string, dst interface{}) error {
	if dst == nil {
		return ErrInvalidArgument
	}

	v := reflect.ValueOf(dst)

	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}

	if v.Kind() != reflect.Struct {
		return ErrInvalidArgument
	}

	for i := 0; i < v.NumField(); i++ {
		field := v.Field(i)
		if !field.CanSet() {
			continue
		}

		tag := v.Type().Field(i).Tag.Get("var")
		if tag == "" {
			continue
		}

		if val, ok := vars[tag]; ok {
			switch field.Type().Kind() {
			case reflect.String:
				field.SetString(val)

			case reflect.Int:
				vi, err := strconv.Atoi(val)
				if err != nil {
					return err
				}

				field.SetInt(int64(vi))

			case reflect.Int64:
				vi, err := strconv.ParseInt(val, 10, 64)
				if err != nil {
					return err
				}

				field.SetInt(vi)

			case reflect.Float32, reflect.Float64:
				vf, err := strconv.ParseFloat(val, 64)
				if err != nil {
					return err
				}

				field.SetFloat(vf)

			case reflect.Bool:
				vb, err := strconv.ParseBool(val)
				if err != nil {
					return err
				}

				field.SetBool(vb)

			case reflect.Slice:
				items := strings.Split(val, ",")
				fmt.Println()
				switch field.Type().Elem().Name() {
				case "int":
					for _, item := range items {
						n, err := strconv.Atoi(strings.TrimSpace(item))
						if err != nil {
							return err
						}

						field.Set(reflect.Append(field, reflect.ValueOf(n)))
					}

				case "string":
					for i := range items {
						items[i] = strings.TrimSpace(items[i])
					}

					field.Set(reflect.ValueOf(items))

				case "float64":
					for _, item := range items {
						n, err := strconv.ParseFloat(strings.TrimSpace(item), 64)
						if err != nil {
							return err
						}

						field.Set(reflect.Append(field, reflect.ValueOf(n)))
					}
				}

			default:
			}
		}
	}

	return nil
}
