package external

import (
	"errors"
)

const (
	TASK_SERVICE_ID_ATLASSIAN = "atlassian"
	TASK_SERVICE_ID_GOOGLE    = "google"
	TASK_SERVICE_ID_SLACK     = "slack"
	TASK_SERVICE_ID_TRELLO    = "trello"

	TASK_SOURCE_ID_GCAL  = "gcal"
	TASK_SOURCE_ID_GMAIL = "gmail"
	TASK_SOURCE_ID_JIRA  = "jira"
)

type Config struct {
	GoogleLoginConfig     OauthConfigWrapper
	GoogleAuthorizeConfig OauthConfigWrapper
	Slack                 OauthConfigWrapper
	Trello                OauthConfigWrapper
	GoogleOverrideURLs    GoogleURLOverrides
	Atlassian             AtlassianConfig
}

func GetConfig() Config {
	return Config{
		GoogleLoginConfig:     getGoogleLoginConfig(),
		GoogleAuthorizeConfig: getGoogleAuthorizeConfig(),
		Slack:                 getSlackConfig(),
		Trello:                getTrelloConfig(),
		Atlassian:             AtlassianConfig{OauthConfig: getAtlassianOauthConfig()},
	}
}

type TaskServiceResult struct {
	Service TaskService
	Details TaskServiceDetails
	Sources []TaskSource
}

type TaskSourceResult struct {
	Source  TaskSource
	Details TaskSourceDetails
}

func (config Config) GetTaskServiceResult(serviceID string) (*TaskServiceResult, error) {
	nameToService := config.GetNameToService()
	result, ok := nameToService[serviceID]
	if !ok {
		return nil, errors.New("task service not found")
	}
	return &result, nil
}

func (config Config) GetTaskSourceResult(serviceID string) (*TaskSourceResult, error) {
	nameToSource := config.getNameToSource()
	result, ok := nameToSource[serviceID]
	if !ok {
		return nil, errors.New("task source not found")
	}
	return &result, nil
}

func (config Config) getNameToSource() map[string]TaskSourceResult {
	atlassianService := AtlassianService{Config: config.Atlassian}
	googleService := GoogleService{
		LoginConfig:     config.GoogleLoginConfig,
		AuthorizeConfig: config.GoogleAuthorizeConfig,
		OverrideURLs:    config.GoogleOverrideURLs,
	}
	return map[string]TaskSourceResult{
		TASK_SOURCE_ID_GMAIL: {
			Details: TaskSourceGmail,
			Source:  GmailSource{Google: googleService},
		},
		TASK_SOURCE_ID_GCAL: {
			Details: TaskSourceGoogleCalendar,
			Source:  GoogleCalendarSource{Google: googleService},
		},
		TASK_SOURCE_ID_JIRA: {
			Details: TaskSourceJIRA,
			Source:  JIRASource{Atlassian: atlassianService},
		},
	}
}

func (config Config) GetNameToService() map[string]TaskServiceResult {
	atlassianService := AtlassianService{Config: config.Atlassian}
	googleService := GoogleService{
		LoginConfig:     config.GoogleLoginConfig,
		AuthorizeConfig: config.GoogleAuthorizeConfig,
		OverrideURLs:    config.GoogleOverrideURLs,
	}
	return map[string]TaskServiceResult{
		TASK_SERVICE_ID_GOOGLE: {
			Service: googleService,
			Details: TaskServiceGoogle,
			Sources: []TaskSource{
				GmailSource{Google: googleService},
				GoogleCalendarSource{Google: googleService},
			},
		},
		TASK_SERVICE_ID_ATLASSIAN: {
			Service: atlassianService,
			Details: TaskServiceAtlassian,
			Sources: []TaskSource{JIRASource{Atlassian: atlassianService}},
		},
		TASK_SERVICE_ID_SLACK: {
			Service: SlackService{Config: config.Slack},
			Details: TaskServiceSlack,
			Sources: []TaskSource{},
		},
		TASK_SERVICE_ID_TRELLO: {
			Service: TrelloService{Config: config.Trello},
			Details: TaskServiceTrello,
			Sources: []TaskSource{},
		},
	}
}

type TaskServiceDetails struct {
	ID           string
	Name         string
	Logo         string
	IsLinkable   bool
	IsSignupable bool
}

var TaskServiceAtlassian = TaskServiceDetails{
	TASK_SERVICE_ID_ATLASSIAN,
	"Atlassian",
	"/images/jira.svg",
	true,
	false,
}
var TaskServiceGoogle = TaskServiceDetails{
	TASK_SERVICE_ID_GOOGLE,
	"Google",
	"/images/gmail.svg",
	true,
	true,
}
var TaskServiceSlack = TaskServiceDetails{
	TASK_SERVICE_ID_SLACK,
	"Slack",
	"/images/slack.svg",
	false,
	false,
}
var TaskServiceTrello = TaskServiceDetails{
	TASK_SERVICE_ID_TRELLO,
	"Trello",
	"/images/trello.svg",
	true,
	false,
}

type TaskSourceDetails struct {
	ID            string
	Name          string
	Logo          string
	IsCompletable bool
	IsReplyable   bool
}

var TaskSourceGoogleCalendar = TaskSourceDetails{
	TASK_SOURCE_ID_GCAL,
	"Google Calendar",
	"/images/gcal.svg",
	false,
	false,
}

var TaskSourceGmail = TaskSourceDetails{
	TASK_SOURCE_ID_GMAIL,
	"Gmail",
	"/images/gmail.svg",
	true,
	true,
}
var TaskSourceJIRA = TaskSourceDetails{
	TASK_SOURCE_ID_JIRA,
	"Jira",
	"/images/jira.svg",
	true,
	false,
}
