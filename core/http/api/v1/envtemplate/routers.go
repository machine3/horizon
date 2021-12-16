package envtemplate

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
			Method:      http.MethodGet,
			Pattern:     fmt.Sprintf("/applications/:%v/envtemplates", _applicationIDParam),
			HandlerFunc: api.Get,
		},
		{
			Method:      http.MethodPost,
			Pattern:     fmt.Sprintf("/applications/:%v/envtemplates", _applicationIDParam),
			HandlerFunc: api.Update,
		},
	}

	route.RegisterRoutes(apiGroup, routes)
}