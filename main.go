package main

import (
	"fmt"
	"log"

	"github.com/fasthttp/router"
	"github.com/igorariza/hyperledger-fabric-on-kubernetes-demo/internal/routers"
	"github.com/igorariza/hyperledger-fabric-on-kubernetes-demo/pkg/kube-clients/k8sclient"
	"github.com/igorariza/hyperledger-fabric-on-kubernetes-demo/pkg/utils"
	"github.com/valyala/fasthttp"
)

const (
	fortyMegaBytes = 4 * 1024 * 1024 * 10
)

var (
	corsAllowHeaders     = "authorization, content-type, set-cookie, cookie, server"
	corsAllowMethods     = "HEAD,GET,POST,PUT,DELETE,OPTIONS"
	corsAllowOrigin      = "*"
	corsAllowCredentials = "true"
)

func main() {

	r := router.New()
	group := r.Group(utils.GetSecretsService("API_VERSION"))
	routers.ThePowerRouter(group)
	routers.BlockchainRouter(group)
	routers.InvitationBlockchainRouter(group)
	routers.ChannelBlockchainRouter(group)
	routers.ChaincodeBlockchainRouter(group)
	routers.SubstrateBlockchainRouter(group)

	fastHttpServer := &fasthttp.Server{
		Name:               "Blockchain Management Service",
		MaxRequestBodySize: fortyMegaBytes,
		Handler:            r.Handler,
	}

	_, err := k8sclient.GetClient()
	fmt.Println("get client", err)
	if err != nil {
		fmt.Println("Erro get client", err)
	}

	// // Listen to port and start the server
	fmt.Println("Server is running in port", utils.GetSecretsService("PORT"))
	log.Fatal(fasthttp.ListenAndServe(":"+utils.GetSecretsService("PORT"), func(ctx *fasthttp.RequestCtx) {
		fmt.Println(string(ctx.Path()))
		r.HandleOPTIONS = true
		ctx.Response.Header.Set("Access-Control-Allow-Credentials", corsAllowCredentials)
		ctx.Response.Header.Set("Access-Control-Allow-Headers", corsAllowHeaders)
		ctx.Response.Header.Set("Access-Control-Allow-Methods", corsAllowMethods)
		ctx.Response.Header.Set("Access-Control-Allow-Origin", corsAllowOrigin)
		ctx.Response.Header.Set("Content-Type", "application/json")
		r.Handler(ctx)
	}), fastHttpServer)

}
