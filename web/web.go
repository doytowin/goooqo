/*
 * The Clear BSD License
 *
 * Copyright (c) 2024, DoytoWin, Inc.
 * All rights reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 */

package web

import (
	"encoding/json"
	"fmt"
	. "github.com/doytowin/goooqo/core"
	"io"
	"net/http"
	"os"
	"regexp"
)

type restService[E Entity, Q Query] struct {
	DataAccess[E]
	idRgx *regexp.Regexp
}

func NewRestService[E Entity, Q Query](
	prefix string,
	dataAccess DataAccess[E],
) http.Handler {
	return &restService[E, Q]{
		DataAccess: dataAccess,
		idRgx:      regexp.MustCompile(prefix + `([\da-fA-F]+)$`),
	}
}

func BuildRestService[E Entity, Q Query](
	prefix string,
	dataAccess DataAccess[E],
) {
	s := NewRestService[E, Q](prefix, dataAccess)
	http.Handle(prefix, s)
}

func (s *restService[E, Q]) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	match := s.idRgx.FindStringSubmatch(request.URL.Path)
	var data any
	var err error
	if len(match) > 0 {
		// Process requests for /<model>/{id}
		id := match[1]
		data, err = s.process(request, id)
	} else if request.Method == "POST" {
		body, _ := io.ReadAll(request.Body)
		var entities []E
		err = json.Unmarshal(body, &entities)
		if NoError(err) {
			data, err = s.CreateMulti(request.Context(), entities)
		}
	} else {
		query := *new(Q)
		queryMap := request.URL.Query()
		ResolveQuery(queryMap, &query)
		if request.Method == "PATCH" {
			body, _ := io.ReadAll(request.Body)
			var entity E
			err = json.Unmarshal(body, &entity)
			if NoError(err) {
				data, err = s.PatchByQuery(request.Context(), entity, query)
			}
		} else if request.Method == "DELETE" {
			data, err = s.DeleteByQuery(request.Context(), query)
		} else {
			data, err = s.Page(request.Context(), query)
		}
	}
	writeResult(writer, err, data)
}

func (s *restService[E, Q]) process(request *http.Request, id string) (any, error) {
	var err error
	var data any
	switch request.Method {
	case "PUT":
		body, _ := io.ReadAll(request.Body)
		entity := *new(E)
		err = json.Unmarshal(body, &entity)
		if NoError(err) {
			err = entity.SetId(&entity, id)
			if err != nil {
				return nil, err
			}
			return s.Update(request.Context(), entity)
		}
	case "PATCH":
		body, _ := io.ReadAll(request.Body)
		entity := *new(E)
		err = json.Unmarshal(body, &entity)
		if NoError(err) {
			err = entity.SetId(&entity, id)
			if err != nil {
				return nil, err
			}
			return s.Patch(request.Context(), entity)
		}
	case "DELETE":
		return s.Delete(request.Context(), id)
	default:
		var entity *E
		entity, err = s.Get(request.Context(), id)
		if entity == nil && NoError(err) {
			err = fmt.Errorf("record not found. id: %s", id)
		} else {
			data = entity
		}
	}
	return data, err
}

func writeResult(writer http.ResponseWriter, err error, data any) {
	response := Response{Data: data, Success: NoError(err), Error: ReadError(err)}
	var bytes []byte
	if os.Getenv("web_intent") == "true" {
		bytes, err = json.MarshalIndent(response, "", "  ")
	} else {
		bytes, err = json.Marshal(response)
	}
	if NoError(err) {
		writer.Header().Set("Content-Type", "application/json; charset=UTF-8")
		_, _ = writer.Write(bytes)
	}
}
