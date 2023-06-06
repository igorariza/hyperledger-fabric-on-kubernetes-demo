package routers

import (
	"github.com/fasthttp/router"
	"github.com/igorariza/hyperledger-fabric-on-kubernetes-demo/internal/handlers"
	"github.com/igorariza/hyperledger-fabric-on-kubernetes-demo/pkg/middleware"
)

func BlockchainRouter(r *router.Group) {
	r.POST("/blockchain/new/", middleware.Chain(handlers.CreateNetwork, middleware.TokenAuth, middleware.CORS()))
	r.POST("/blockchain/peer/new/", middleware.Chain(handlers.AddPeerToCorporation, middleware.TokenAuth, middleware.CORS()))
	r.POST("/blockchain/orderer/new/", middleware.Chain(handlers.AddOrdererToNetwork, middleware.TokenAuth, middleware.CORS()))

	r.GET("/blockchains/", middleware.Chain(handlers.GetBlockchains, middleware.TokenAuth, middleware.CORS()))
	r.GET("/blockchain/", middleware.Chain(handlers.GetBlockchainById, middleware.TokenAuth, middleware.CORS()))
	r.GET("/blockchain/corporations/{corporation_id}/peers/", middleware.Chain(handlers.GetPeersByCorporationId, middleware.TokenAuth, middleware.CORS()))
	r.GET("/blockchain/corporations/", middleware.Chain(handlers.GetCorporationsByBlockchainId, middleware.TokenAuth, middleware.CORS()))
	r.GET("/blockchain/orderers/", middleware.Chain(handlers.GetOrderersByNetworkId, middleware.TokenAuth, middleware.CORS()))
	r.GET("/blockchain/corporation/", middleware.Chain(handlers.GetCorporationById, middleware.TokenAuth, middleware.CORS()))

	r.DELETE("/blockchain/delete/network/", middleware.Chain(handlers.DeleteNetworkHlf, middleware.CORS()))
}
