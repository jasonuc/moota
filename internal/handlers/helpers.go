package handlers

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"maps"
	"net/http"
	"strconv"
	"strings"

	"github.com/go-chi/chi/v5"
)

func readStringReqParam(r *http.Request, key string) (string, error) {
	value := chi.URLParam(r, key)
	if value == "" {
		return "", fmt.Errorf("missing required param")
	}
	return value, nil
}

// func readStringQueryParam(r *http.Request, key string) (string, error) {
// 	value := r.URL.Query().Get(key)
// 	if value == "" {
// 		return "", fmt.Errorf("missing required query param: %s", key)
// 	}
// 	return value, nil
// }

func readFloatQueryParam(r *http.Request, key string) (float64, error) {
	value := r.URL.Query().Get(key)
	if value == "" {
		return 0, fmt.Errorf("missing required query param: %s", key)
	}

	floatVal, err := strconv.ParseFloat(value, 64)
	if err != nil {
		return 0, err
	}

	return floatVal, nil
}

func readBoolQueryParam(r *http.Request, key string) bool {
	value := r.URL.Query().Get(key)
	if value == "" {
		return false
	}

	boolValue, err := strconv.ParseBool(value)
	if err != nil {
		return false
	}

	return boolValue
}

type envelope map[string]any

func readJSON(w http.ResponseWriter, r *http.Request, v any) error {
	r.Body = http.MaxBytesReader(w, r.Body, http.DefaultMaxHeaderBytes)

	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()

	err := dec.Decode(v)
	if err != nil {
		var invalidUnmarshalErr *json.InvalidUnmarshalError
		var unmarshalTypeErr *json.UnmarshalTypeError
		var maxBytesErr *http.MaxBytesError

		switch {
		case errors.Is(err, io.ErrUnexpectedEOF):
			return errors.New("invalid request body contains bad-json")
		case errors.As(err, &unmarshalTypeErr):
			if unmarshalTypeErr.Field != "" {
				return fmt.Errorf("invalid request body contains unsurpotted field %q", unmarshalTypeErr.Field)
			} else {
				return errors.New("invalid request body contains unsurpotted field")
			}
		case errors.As(err, &invalidUnmarshalErr):
			panic(err)
		case errors.As(err, &maxBytesErr):
			return fmt.Errorf("body must not be larger than %d bytes", maxBytesErr.Limit)
		case strings.HasPrefix(err.Error(), "json: unknown field "):
			disallowedField := strings.TrimPrefix(err.Error(), "json: unknown field ")
			return fmt.Errorf("body contains unkown key %s", disallowedField)
		default:
			return err
		}
	}

	return nil
}

func writeJSON(w http.ResponseWriter, status int, data envelope, header http.Header) error {
	js, err := json.Marshal(data)
	if err != nil {
		return err
	}
	js = append(js, '\n')

	maps.Copy(w.Header(), header)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	_, err = w.Write(js)
	if err != nil {
		return err
	}

	return nil
}
