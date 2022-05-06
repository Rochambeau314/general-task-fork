package external

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/rs/zerolog/log"
	"net/http"

	"github.com/GeneralTask/task-manager/backend/config"
	"github.com/GeneralTask/task-manager/backend/constants"
	"github.com/GeneralTask/task-manager/backend/database"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"golang.org/x/oauth2"
)

type LinearService struct {
	Config       OauthConfigWrapper
	ConfigValues LinearConfigValues
}

type LinearConfigValues struct {
	UserInfoURL   *string
	TaskFetchURL  *string
	TaskUpdateURL *string
}

func getLinearConfig() *OauthConfig {
	return &OauthConfig{Config: &oauth2.Config{
		ClientID:     config.GetConfigValue("LINEAR_OAUTH_CLIENT_ID"),
		ClientSecret: config.GetConfigValue("LINEAR_OAUTH_CLIENT_SECRET"),
		RedirectURL:  config.GetConfigValue("SERVER_URL") + "link/linear/callback/",
		Scopes:       []string{"read", "write"},
		Endpoint: oauth2.Endpoint{
			AuthURL:  "https://linear.app/oauth/authorize",
			TokenURL: "https://api.linear.app/oauth/token",
		},
	}}
}

func (linear LinearService) GetLinkURL(stateTokenID primitive.ObjectID, userID primitive.ObjectID) (*string, error) {
	authURL := linear.Config.AuthCodeURL(stateTokenID.Hex(), oauth2.AccessTypeOffline, oauth2.ApprovalForce)
	return &authURL, nil
}

func (linear LinearService) GetSignupURL(stateTokenID primitive.ObjectID, forcePrompt bool) (*string, error) {
	return nil, errors.New("linear does not support signup")
}

func (linear LinearService) HandleLinkCallback(params CallbackParams, userID primitive.ObjectID) error {
	parentCtx := context.Background()
	db, dbCleanup, err := database.GetDBConnection()
	if err != nil {
		return errors.New("internal server error")
	}
	defer dbCleanup()

	extCtx, cancel := context.WithTimeout(parentCtx, constants.ExternalTimeout)
	defer cancel()
	token, err := linear.Config.Exchange(extCtx, *params.Oauth2Code)
	if err != nil {
		log.Error().Msgf("failed to fetch token from Linear: %v", err)
		return errors.New("internal server error")
	}
	log.Debug().Interface("token", token).Send()

	tokenString, err := json.Marshal(&token)
	log.Info().Msgf("token string: %s", string(tokenString))
	if err != nil {
		log.Error().Msgf("error parsing token: %v", err)
		return errors.New("internal server error")
	}

	externalAPITokenCollection := database.GetExternalTokenCollection(db)
	dbCtx, cancel := context.WithTimeout(parentCtx, constants.DatabaseTimeout)
	defer cancel()
	// TODO: add DisplayID, AccountID, etc.
	_, err = externalAPITokenCollection.UpdateOne(
		dbCtx,
		bson.M{"$and": []bson.M{{"user_id": userID}, {"service_id": TASK_SERVICE_ID_LINEAR}}},
		bson.M{"$set": &database.ExternalAPIToken{
			UserID:         userID,
			ServiceID:      TASK_SERVICE_ID_LINEAR,
			Token:          string(tokenString),
			AccountID:      "todo",
			DisplayID:      "todo",
			IsUnlinkable:   true,
			IsPrimaryLogin: false,
		}},
		options.Update().SetUpsert(true),
	)
	if err != nil {
		log.Error().Msgf("error saving token: %v", err)
		return errors.New("internal server error")
	}

	return nil
}

func (linear LinearService) HandleSignupCallback(params CallbackParams) (primitive.ObjectID, *bool, *string, error) {
	return primitive.NilObjectID, nil, nil, errors.New("linear does not support signup")
}

func getLinearHttpClient(db *mongo.Database, userID primitive.ObjectID, accountID string) *http.Client {
	return getExternalOauth2Client(db, userID, accountID, TASK_SERVICE_ID_LINEAR, getLinearConfig())
}