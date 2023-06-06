package middleware

import (
	"encoding/json"
	"fmt"
	"net/http"

	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/valyala/fasthttp"
)

var (
	corsAllowHeaders     = "authorization, content-type, set-cookie, cookie, server"
	corsAllowMethods     = "HEAD,GET,POST,PUT,DELETE,OPTIONS"
	corsAllowOrigin      = "*"
	corsAllowCredentials = "true"
)

type Middleware func(fasthttp.RequestHandler) fasthttp.RequestHandler

func CORS() Middleware {
	return func(next fasthttp.RequestHandler) fasthttp.RequestHandler {
		return func(ctx *fasthttp.RequestCtx) {
			fmt.Println("Cors Init")
			ctx.Response.Header.Set("Access-Control-Allow-Credentials", corsAllowCredentials)
			ctx.Response.Header.Set("Access-Control-Allow-Headers", corsAllowHeaders)
			ctx.Response.Header.Set("Access-Control-Allow-Methods", corsAllowMethods)
			ctx.Response.Header.Set("Access-Control-Allow-Origin", corsAllowOrigin)
			ctx.Response.Header.Set("Content-Type", "application/json")
			next(ctx)
		}
	}

}

type Body struct {
	RolesIds []string `json:"roles_ids"`
}

func Authorization(permission string) Middleware {
	var (
		// uID            uint32
		// projectId      string
		// organizationId string
		err  error
		body *Body
	)

	return func(next fasthttp.RequestHandler) fasthttp.RequestHandler {
		return func(ctx *fasthttp.RequestCtx) {
			token := ctx.Request.Header.Peek("Authorization")
			ApiKey := ctx.Request.Header.Peek("x-api-key")

			// checkUser := &clientUser.CheckUserData{
			// 	Permission: permission,
			// }

			if token != nil && len(string(token)) > 7 {
				token := string(token)
				token = token[7:]
				//user, err := clientToken.GetOneTokenCCP(token, nil)
				if err != nil {
					ctx.SetStatusCode(http.StatusUnauthorized)
					//ctx.Response.SetBodyString(utils.NewErrorCcp(err).Error())
					return
				}
				// ctx.SetUserValue("UserId", user.UserIdDex)
				// ctx.SetUserValue("UserEmail", user.Email)
				// ctx.SetUserValue("UserName", user.FirstName)

			} else if len(ApiKey) > 0 {
				ctx.SetUserValue("ApiKey", string(ApiKey))
				// apiKey, err := clientApikey.GetOneApiKey(&pbApikeys.GetOneApiKeyRequest{
				// 	Value: string(ApiKey),
				// })

				if err != nil {
					ctx.Response.SetStatusCode(http.StatusUnauthorized)
					//err := utils.NewErrorCcp(err)
					ctx.Response.SetBodyString(err.Error())
					return
				}
				// if apiKey.UserId != "" {
				// 	ctx.SetUserValue("UserIdApiKey", apiKey.UserId)
				// }
			} else {
				ctx.Response.SetStatusCode(http.StatusUnauthorized)
				st := status.New(codes.Unauthenticated, "token or apikeys is empty")
				detail := &errdetails.BadRequest_FieldViolation{
					Field:       "user_id",
					Description: "UserId dex not found",
				}

				st, err = st.WithDetails(detail)
				//err := utils.NewErrorCcp(st.Err())
				ctx.Response.SetBodyString(err.Error())
				return
			}
			/*
				CHECK PERMISSIONS
			*/
			if len(ctx.PostBody()) > 0 {
				if err := json.Unmarshal(ctx.PostBody(), &body); err != nil {

					ctx.Response.SetStatusCode(http.StatusBadRequest)
					st := status.New(codes.InvalidArgument, "Error unmarshaling")
					detail := &errdetails.ErrorInfo{
						Reason:   "Error unmarshaling",
						Domain:   "ccp-controlplane-service",
						Metadata: map[string]string{"error": err.Error()},
					}

					st, err = st.WithDetails(detail)
					if err != nil {
						return
					}
					//ctx.Response.SetBodyString(utils.NewErrorCcp(st.Err()).Error())
					return
				}

				if body.RolesIds != nil && len(body.RolesIds) > 0 {
					//checkUser.RolesIds = body.RolesIds
				}
			} else {
				if ctx.UserValue("RoleId") != nil {
					//roleId := ctx.UserValue("RoleId").(string)
					//checkUser.RolesIds = []string{roleId}
				}
			}

			if ctx.UserValue("UId") != nil {
				// uID, err = utils.StringToUint32(ctx.UserValue("UId").(string))
				// if err != nil {
				// 	//bylogs.LogErr(err)
				// 	ctx.Response.SetStatusCode(http.StatusBadRequest)
				// 	ctx.Response.SetBodyString(fmt.Sprintf("{\"error\": \"%s\"}", utils.NewErrorCcp(err)))
				// 	return
				// }
				// checkUser.UserId = uID

			}
			if ctx.UserValue("OrganizationId") != nil {
				//organizationId = ctx.UserValue("OrganizationId").(string)
				//checkUser.OrganizationId = organizationId
			}
			if ctx.UserValue("ProjectId") != nil {
				//projectId = ctx.UserValue("ProjectId").(string)
				//checkUser.ProjectId = projectId
			}

			if ctx.UserValue("UserId") != nil {
				//checkUser.UserIdAdmin = ctx.UserValue("UserId").(string)
			} else if ctx.UserValue("ApiKey") != nil {
				//checkUser.ApiKeyValue = ctx.UserValue("ApiKey").(string)
			}

			//PERMISSIONS CHECK
			//_, err := clientUser.CheckUser(checkUser)
			// if err != nil {
			// 	bylogs.LogErr("Error checking user permission", err)
			// 	ctx.Response.SetStatusCode(http.StatusUnauthorized)
			// 	ctx.Response.SetBodyString(utils.NewErrorCcp(err).Error())
			// 	return
			// }
			next(ctx)
		}
	}
}

func Chain(f fasthttp.RequestHandler, middlewares ...Middleware) fasthttp.RequestHandler {
	for _, m := range middlewares {
		f = m(f)
	}

	return f
}

func TokenAuth(next fasthttp.RequestHandler) fasthttp.RequestHandler {
	return func(ctx *fasthttp.RequestCtx) {
		fmt.Println("Token auth")
		token := ctx.Request.Header.Peek("Authorization")
		ApiKey := ctx.Request.Header.Peek("x-api-key")

		if token != nil && len(string(token)) > 7 {
			token := string(token)
			token = token[7:]

			// user, err := clientToken.GetOneTokenCCP(token, nil)
			// if err != nil {
			// 	ctx.SetStatusCode(http.StatusUnauthorized)
			// 	ctx.Write([]byte(utils.NewErrorCcp(err).Error()))
			// 	return
			// }

			// ctx.SetUserValue("UserId", user.UserIdDex)
			// ctx.SetUserValue("UserEmail", user.Email)
			// ctx.SetUserValue("UserName", user.FirstName)
			// ctx.SetUserValue("EmailVerified", user.EmailVerified)

		} else if len(ApiKey) > 7 {
			ctx.SetUserValue("ApiKey", string(ApiKey))

			// apiKey, err := clientApikey.GetOneApiKey(&pbApikeys.GetOneApiKeyRequest{
			// 	Uuid: string(ApiKey),
			// })

			// if err != nil {
			// 	ctx.Response.SetStatusCode(http.StatusUnauthorized)
			// 	err := utils.NewErrorCcp(err)
			// 	ctx.Response.SetBodyString(err.Error())
			// 	return
			// }
			// if apiKey.UserId != "" {
			// 	ctx.SetUserValue("UserIdApiKey", apiKey.UserId)
			// }

		} else {
			ctx.Response.SetStatusCode(http.StatusUnauthorized)
			st := status.New(codes.Unauthenticated, "token or apikeys is empty")
			detail := &errdetails.BadRequest_FieldViolation{
				Field:       "user_id",
				Description: "UserId dex not found",
			}

			st, _ = st.WithDetails(detail)
			//err := utils.NewErrorCcp(st.Err())
			//ctx.Response.SetBodyString(err.Error())
			return
		}

		next(ctx)
	}
}

func WithoutToken(next fasthttp.RequestHandler) fasthttp.RequestHandler {
	return func(ctx *fasthttp.RequestCtx) {
		token := ctx.Request.Header.Peek("Authorization")
		ApiKey := ctx.Request.Header.Peek("x-api-key")
		ctx.Response.Header.Set("Content-Type", "application/json")

		if token != nil && string(token) != "" || ApiKey != nil && string(ApiKey) != "" {
			ctx.Response.SetStatusCode(http.StatusBadRequest)
			st := status.New(codes.Unauthenticated, "Token denied")
			detail := &errdetails.ErrorInfo{
				Reason: "use of token and x-api-key is not allowed",
				Domain: "WithoutToken",
			}
			st, err := st.WithDetails(detail)
			if err != nil {
				ctx.Response.SetBodyString(fmt.Sprintf("{\"error\": \"%s\"}", "use of token and x-api-key is not allowed"))
				return
			}
			//err = utils.NewErrorCcp(st.Err())
			ctx.Response.SetBodyString(err.Error())
			return
		}
		next(ctx)
	}

}
