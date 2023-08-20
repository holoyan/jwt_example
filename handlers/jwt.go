package handlers

import (
	"../core"
	"../models"
	"context"
	"encoding/base64"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"log"
	"net/http"
)

func IssueTokens(res http.ResponseWriter, req *http.Request)  {
	queryParams := req.URL.Query()
	userId := queryParams.Get("guid")

	var accessToken, accessTokenId string = getAccessToken(userId)
	var refreshToken = getRefreshToken(accessTokenId)

	core.JsonResponse(res, map[string]interface{}{
		"accessToken": accessToken,
		"refreshToken": base64.StdEncoding.EncodeToString([]byte(refreshToken)),
	}, 200)
}

func UpdateTokens(res http.ResponseWriter, req *http.Request)  {
	request := core.JsonBody(req)
	accessToken := request["accessToken"]
	var refreshToken = request["refreshToken"]

	decodedJwt := core.DecodeJwt(accessToken)

	response := make(map[string]interface{})
	response["success"] = decodedJwt.Success
	response["message"] = decodedJwt.Message
	if !decodedJwt.Success {
		core.JsonResponse(res, response, 400)
		return
	}

	// get access token
	result := core.FindOne("access_tokens", "_id", decodedJwt.Jti)
	if !result.Success {
		errorResponse := map[string]interface{} {
			"message": result.Message,
			"success": result.Success,
		}
		core.JsonResponse(res, errorResponse, 400)
		return
	}
	userId := result.Result["user_id"].(string)

	// get refresh token by accessToken id and compare values
	refreshTokenResult := core.FindOne("refresh_tokens", "access_token_id", decodedJwt.Jti)
	//primitive.ObjectID.String()
	if !refreshTokenResult.Success {
		errorResponse := map[string]interface{} {
			"message": refreshTokenResult.Message,
			"success": refreshTokenResult.Success,
		}
		core.JsonResponse(res, errorResponse, 400)
		return
	}

	hashedToken := refreshTokenResult.Result["token"].(string)

	var decodedRefresh, _ = base64.StdEncoding.DecodeString(refreshToken)
	if !core.CheckBcrypt(string(decodedRefresh), hashedToken) {
		response["success"] = false
		response["message"] = "token mismatch"
		core.JsonResponse(res, response, 400)
		return
	}

	// remove old access and refresh token
	core.DeleteOne("access_tokens", "_id", decodedJwt.Jti)
	core.DeleteOne("refresh_tokens", "_id", refreshTokenResult.Result["_id"].(primitive.ObjectID).Hex())

	// get new
	var newAccessToken, accessTokenId string = getAccessToken(userId)
	var newRefreshToken = getRefreshToken(accessTokenId)

	core.JsonResponse(res, map[string]interface{}{
		"accessToken": newAccessToken,
		"refreshToken": base64.StdEncoding.EncodeToString([]byte(newRefreshToken)),
	}, 200)
}

func getAccessToken(userId string) (string, string) {
	accessToken := models.AccessToken{User_Id: userId}
	data, err := core.DB.Collection("access_tokens").InsertOne(context.Background(), accessToken)
	if err != nil {
		log.Fatal(err)
	}

	insertedID, ok := data.InsertedID.(primitive.ObjectID)
	if !ok {
		log.Fatal("InsertedID assertion failed")
	}

	insertedIDString := insertedID.Hex()
	return core.GenerateJwt(userId, insertedIDString, 24), insertedIDString
}

func getRefreshToken(accessTokenId string) string {
	tokenData := core.GenerateRandomString(64)
	refreshToken := models.RefreshToken{Access_Token_Id: accessTokenId, Token: core.Bcrypt(tokenData)}
	_, err := core.DB.Collection("refresh_tokens").InsertOne(context.Background(), refreshToken)
	if err != nil {
		log.Fatal(err)
	}

	return tokenData
}
