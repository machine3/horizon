package application

import (
	"fmt"
	"net/http"

	"g.hz.netease.com/horizon/pkg/server/route"

	"github.com/gin-gonic/gin"
)

// RegisterRoutes register routes
func RegisterRoutes(engine *gin.Engine, api *API) {
	apiGroup := engine.Group("/apis/core/v1")
	var routes = route.Routes{
		{
			Method:      http.MethodPost,
			Pattern:     fmt.Sprintf("/groups/:%v/applications", _groupIDParam),
			HandlerFunc: api.Create,
		},
		{
			Method:      http.MethodGet,
			Pattern:     fmt.Sprintf("/applications/:%v", _applicationIDParam),
			HandlerFunc: api.Get,
		},
		{
			Method:      http.MethodPut,
			Pattern:     fmt.Sprintf("/applications/:%v", _applicationIDParam),
			HandlerFunc: api.Update,
		},
		{
			Method:      http.MethodDelete,
			Pattern:     fmt.Sprintf("/applications/:%v", _applicationIDParam),
			HandlerFunc: api.Delete,
		},
	}

	route.RegisterRoutes(apiGroup, routes)
}
