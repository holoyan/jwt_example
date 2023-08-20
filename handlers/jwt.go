package handlers

import (
	"../core"
	"../models"
	"context"
	"encoding/base64"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
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

	coll := core.DB.Collection("access_tokens")

	var result map[string]string
	objectID, err1 := primitive.ObjectIDFromHex(decodedJwt.Jti)
	filter := bson.M{"_id": objectID}
	if err1 != nil {
		log.Fatal(err1)
	}

	err := coll.FindOne(context.Background(), filter).Decode(&result)
	var userId string = result["user_id"]
	if err != nil {
		// ErrNoDocuments means that the filter did not match any documents in the collection
		response["success"] = false
		if err == mongo.ErrNoDocuments {
			response["message"] = "Not found"
			core.JsonResponse(res, response, 404)
			return
		}
		response["message"] = "Something went wrong"
		core.JsonResponse(res, response, 400)
		return
	}

	// get refresh token by accessToken id and compare values

	rColl := core.DB.Collection("refresh_tokens")

	var refreshResult map[string]string
	refreshFilter := bson.M{"access_token_id": decodedJwt.Jti}

	err5 := rColl.FindOne(context.Background(), refreshFilter).Decode(&refreshResult)
	if err5 != nil {
		// ErrNoDocuments means that the filter did not match any documents in the collection
		response["success"] = false
		if err5 == mongo.ErrNoDocuments {
			response["message"] = "refresh Not found"
			core.JsonResponse(res, response, 404)
			return
		}
		response["message"] = "refresh token: Something went wrong"
		core.JsonResponse(res, response, 400)
		log.Fatal(err5)
		return
	}
	var hashedToken string = refreshResult["token"]

	var decodedRefresh, _ = base64.StdEncoding.DecodeString(refreshToken)
	if !core.CheckBcrypt(string(decodedRefresh), hashedToken) {
		response["success"] = false
		response["message"] = "token mismatch"
		core.JsonResponse(res, response, 400)
		return
	}

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
