package api

import (
	"context"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/GeneralTask/task-manager/backend/constants"
	"github.com/GeneralTask/task-manager/backend/database"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/oauth2"
	"google.golang.org/api/calendar/v3"
	"google.golang.org/api/googleapi"
)

func TestGetGoogleConfig(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		googleConfig := GetGoogleConfig()
		assert.Equal(
			t,
			"https://accounts.google.com/o/oauth2/auth?access_type=offline&client_id=786163085684-uvopl20u17kp4p2vd951odnm6f89f2f6.apps.googleusercontent.com&prompt=consent&redirect_uri=http%3A%2F%2Flocalhost%3A8080%2Flogin%2Fcallback%2F&response_type=code&scope=https%3A%2F%2Fwww.googleapis.com%2Fauth%2Fuserinfo.email+https%3A%2F%2Fwww.googleapis.com%2Fauth%2Fuserinfo.profile+https%3A%2F%2Fwww.googleapis.com%2Fauth%2Fgmail.modify+https%3A%2F%2Fwww.googleapis.com%2Fauth%2Fcalendar.events&state=state-token",
			googleConfig.AuthCodeURL("state-token", oauth2.AccessTypeOffline, oauth2.ApprovalForce),
		)
	})
}

func TestCalendar(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		standardEvent := calendar.Event{
			Created:        "2021-02-25T17:53:01.000Z",
			Summary:        "Standard Event",
			Start:          &calendar.EventDateTime{DateTime: "2021-03-06T15:00:00-05:00"},
			End:            &calendar.EventDateTime{DateTime: "2021-03-06T15:30:00-05:00"},
			HtmlLink:       "generaltask.io",
			Id:             "standard_event",
			ServerResponse: googleapi.ServerResponse{HTTPStatusCode: 0},
		}

		startTime, _ := time.Parse(time.RFC3339, "2021-03-06T15:00:00-05:00")
		endTime, _ := time.Parse(time.RFC3339, "2021-03-06T15:30:00-05:00")

		db, dbCleanup, err := database.GetDBConnection()
		assert.NoError(t, err)
		defer dbCleanup()
		userID := primitive.NewObjectID()
		standardTask := database.CalendarEvent{
			TaskBase: database.TaskBase{
				IDOrdering:    0,
				IDExternal:    "standard_event",
				IDTaskSection: constants.IDTaskSectionToday,
				Deeplink:      "generaltask.io",
				Title:         "Standard Event",
				Source:        database.TaskSourceGoogleCalendar,
				UserID:        userID,
			},
			DatetimeStart: primitive.NewDateTimeFromTime(startTime),
			DatetimeEnd:   primitive.NewDateTimeFromTime(endTime),
		}

		autoEvent := calendar.Event{
			Created:        "2021-02-25T17:53:01.000Z",
			Summary:        "Auto Event (via Clockwise)",
			Start:          &calendar.EventDateTime{DateTime: "2021-03-06T15:00:00-05:00"},
			End:            &calendar.EventDateTime{DateTime: "2021-03-06T15:30:00-05:00"},
			HtmlLink:       "generaltask.io",
			Id:             "auto_event",
			ServerResponse: googleapi.ServerResponse{HTTPStatusCode: 0},
		}

		allDayEvent := calendar.Event{
			Created:        "2021-02-25T17:53:01.000Z",
			Summary:        "All day Event",
			Start:          &calendar.EventDateTime{Date: "2021-03-06"},
			End:            &calendar.EventDateTime{Date: "2021-03-06"},
			HtmlLink:       "generaltask.io",
			Id:             "all_day_event",
			ServerResponse: googleapi.ServerResponse{HTTPStatusCode: 0},
		}

		server := getServerForTasks([]*calendar.Event{&standardEvent, &allDayEvent, &autoEvent})
		defer server.Close()

		var calendarResult = make(chan CalendarResult)
		api := &API{
			GoogleOverrideURLs: GoogleURLOverrides{CalendarFetchURL: &server.URL},
		}
		go LoadCalendarEvents(api, userID, "exampleAccountID", 0, nil, calendarResult)
		result := <-calendarResult
		assert.NoError(t, result.Error)
		assert.Equal(t, 1, len(result.CalendarEvents))
		firstTask := result.CalendarEvents[0]
		assertCalendarEventsEqual(t, &standardTask, firstTask)

		taskCollection := db.Collection("tasks")

		var calendarEventFromDB database.CalendarEvent
		err = taskCollection.FindOne(
			context.TODO(),
			bson.M{"$and": []bson.M{
				{"id_external": "standard_event"},
				{"source.name": database.TaskSourceGoogleCalendar.Name},
				{"user_id": userID},
			}},
		).Decode(&calendarEventFromDB)
		assert.NoError(t, err)
		assertCalendarEventsEqual(t, &standardTask, &calendarEventFromDB)
		assert.Equal(t, "exampleAccountID", calendarEventFromDB.SourceAccountID)
	})
	t.Run("ExistingEvent", func(t *testing.T) {
		standardEvent := calendar.Event{
			Created:        "2021-02-25T17:53:01.000Z",
			Summary:        "Standard Event",
			Start:          &calendar.EventDateTime{DateTime: "2021-03-06T15:00:00-05:00"},
			End:            &calendar.EventDateTime{DateTime: "2021-03-06T15:30:00-05:00"},
			HtmlLink:       "generaltask.io",
			Id:             "standard_event",
			ServerResponse: googleapi.ServerResponse{HTTPStatusCode: 0},
		}

		startTime, _ := time.Parse(time.RFC3339, "2021-03-06T15:00:00-05:00")
		oldEndtime, _ := time.Parse(time.RFC3339, "2021-03-06T15:35:00-05:00")
		endTime, _ := time.Parse(time.RFC3339, "2021-03-06T15:30:00-05:00")

		db, dbCleanup, err := database.GetDBConnection()
		assert.NoError(t, err)
		defer dbCleanup()
		userID := primitive.NewObjectID()
		standardTask := database.CalendarEvent{
			TaskBase: database.TaskBase{
				IDOrdering:      1,
				IDExternal:      "standard_event",
				IDTaskSection:   constants.IDTaskSectionToday,
				Deeplink:        "generaltask.io",
				Title:           "Standard Event",
				Source:          database.TaskSourceGoogleCalendar,
				UserID:          userID,
				SourceAccountID: "exampleAccountID",
			},
			DatetimeStart: primitive.NewDateTimeFromTime(startTime),
			DatetimeEnd:   primitive.NewDateTimeFromTime(oldEndtime),
		}
		database.GetOrCreateTask(db, userID, "standard_event", database.TaskSourceGoogleCalendar, standardTask)
		// Rescheduling end time along shouldn't trigger a reset like in the next test case
		standardTask.DatetimeEnd = primitive.NewDateTimeFromTime(endTime)

		autoEvent := calendar.Event{
			Created:        "2021-02-25T17:53:01.000Z",
			Summary:        "Auto Event (via Clockwise)",
			Start:          &calendar.EventDateTime{DateTime: "2021-03-06T15:00:00-05:00"},
			End:            &calendar.EventDateTime{DateTime: "2021-03-06T15:30:00-05:00"},
			HtmlLink:       "generaltask.io",
			Id:             "auto_event",
			ServerResponse: googleapi.ServerResponse{HTTPStatusCode: 0},
		}

		allDayEvent := calendar.Event{
			Created:        "2021-02-25T17:53:01.000Z",
			Summary:        "All day Event",
			Start:          &calendar.EventDateTime{Date: "2021-03-06"},
			End:            &calendar.EventDateTime{Date: "2021-03-06"},
			HtmlLink:       "generaltask.io",
			Id:             "all_day_event",
			ServerResponse: googleapi.ServerResponse{HTTPStatusCode: 0},
		}

		server := getServerForTasks([]*calendar.Event{&standardEvent, &allDayEvent, &autoEvent})
		defer server.Close()

		var calendarResult = make(chan CalendarResult)
		api := &API{
			GoogleOverrideURLs: GoogleURLOverrides{CalendarFetchURL: &server.URL},
		}
		go LoadCalendarEvents(api, userID, "exampleAccountID", 0, nil, calendarResult)
		result := <-calendarResult
		assert.NoError(t, result.Error)
		assert.Equal(t, 1, len(result.CalendarEvents))
		firstTask := result.CalendarEvents[0]
		assertCalendarEventsEqual(t, &standardTask, firstTask)

		taskCollection := db.Collection("tasks")

		var calendarEventFromDB database.CalendarEvent
		err = taskCollection.FindOne(
			context.TODO(),
			bson.M{"$and": []bson.M{
				{"id_external": "standard_event"},
				{"source.name": database.TaskSourceGoogleCalendar.Name},
				{"user_id": userID},
			}},
		).Decode(&calendarEventFromDB)
		assert.NoError(t, err)
		assertCalendarEventsEqual(t, &standardTask, &calendarEventFromDB)
		assert.Equal(t, "exampleAccountID", calendarEventFromDB.SourceAccountID)
	})
	t.Run("RescheduledEvent", func(t *testing.T) {
		standardEvent := calendar.Event{
			Created:        "2021-02-25T17:53:01.000Z",
			Summary:        "Standard Event",
			Start:          &calendar.EventDateTime{DateTime: "2021-03-06T15:00:00-05:00"},
			End:            &calendar.EventDateTime{DateTime: "2021-03-06T15:30:00-05:00"},
			HtmlLink:       "generaltask.io",
			Id:             "standard_event",
			ServerResponse: googleapi.ServerResponse{HTTPStatusCode: 0},
		}

		// Start time in DB is different to simulate rescheduling a meeting
		oldStartTime, _ := time.Parse(time.RFC3339, "2021-03-06T13:00:00-05:00")
		startTime, _ := time.Parse(time.RFC3339, "2021-03-06T15:00:00-05:00")
		endTime, _ := time.Parse(time.RFC3339, "2021-03-06T15:30:00-05:00")

		db, dbCleanup, err := database.GetDBConnection()
		assert.NoError(t, err)
		defer dbCleanup()
		userID := primitive.NewObjectID()
		standardTask := database.CalendarEvent{
			TaskBase: database.TaskBase{
				IDOrdering:      1,
				IDExternal:      "standard_event",
				IDTaskSection:   constants.IDTaskSectionToday,
				Deeplink:        "generaltask.io",
				Title:           "Standard Event",
				Source:          database.TaskSourceGoogleCalendar,
				UserID:          userID,
				SourceAccountID: "exampleAccountID",
			},
			DatetimeStart: primitive.NewDateTimeFromTime(oldStartTime),
			DatetimeEnd:   primitive.NewDateTimeFromTime(endTime),
		}
		database.GetOrCreateTask(db, userID, "standard_event", database.TaskSourceGoogleCalendar, standardTask)
		standardTask.DatetimeStart = primitive.NewDateTimeFromTime(startTime)
		// IDOrdering expected to be zero because ordering is reset upon rescheduling
		standardTask.IDOrdering = 0

		server := getServerForTasks([]*calendar.Event{&standardEvent})
		defer server.Close()

		var calendarResult = make(chan CalendarResult)
		api := &API{
			GoogleOverrideURLs: GoogleURLOverrides{CalendarFetchURL: &server.URL},
		}
		go LoadCalendarEvents(api, userID, "exampleAccountID", 0, nil, calendarResult)
		result := <-calendarResult
		assert.NoError(t, result.Error)
		assert.Equal(t, 1, len(result.CalendarEvents))
		firstTask := result.CalendarEvents[0]
		assertCalendarEventsEqual(t, &standardTask, firstTask)

		taskCollection := db.Collection("tasks")

		var calendarEventFromDB database.CalendarEvent
		err = taskCollection.FindOne(
			context.TODO(),
			bson.M{"$and": []bson.M{
				{"id_external": "standard_event"},
				{"source.name": database.TaskSourceGoogleCalendar.Name},
				{"user_id": userID},
			}},
		).Decode(&calendarEventFromDB)
		assert.NoError(t, err)
		// DB is not updated until task merge
		standardTask.IDOrdering = 1
		assertCalendarEventsEqual(t, &standardTask, &calendarEventFromDB)
		assert.Equal(t, "exampleAccountID", calendarEventFromDB.SourceAccountID)
	})
	t.Run("EmptyResult", func(t *testing.T) {
		server := getServerForTasks([]*calendar.Event{})
		api := API{
			GoogleOverrideURLs: GoogleURLOverrides{CalendarFetchURL: &server.URL},
		}
		defer server.Close()
		var calendarResult = make(chan CalendarResult)
		go LoadCalendarEvents(&api, primitive.NewObjectID(), "exampleAccountID", 0, nil, calendarResult)
		result := <-calendarResult
		assert.NoError(t, result.Error)
		assert.Equal(t, 0, len(result.CalendarEvents))
	})
}

func assertCalendarEventsEqual(t *testing.T, a *database.CalendarEvent, b *database.CalendarEvent) {
	assert.Equal(t, a.DatetimeStart, b.DatetimeStart)
	assert.Equal(t, a.DatetimeEnd, b.DatetimeEnd)
	assert.Equal(t, a.Deeplink, b.Deeplink)
	assert.Equal(t, a.IDExternal, b.IDExternal)
	assert.Equal(t, a.IDOrdering, b.IDOrdering)
	assert.Equal(t, a.IDTaskSection, b.IDTaskSection)
	assert.Equal(t, a.Title, b.Title)
	assert.Equal(t, a.Source, b.Source)
}

func getGmailArchiveServer(t *testing.T, expectedLabel string) *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		body, err := ioutil.ReadAll(r.Body)
		assert.NoError(t, err)
		assert.Equal(t, "{\"removeLabelIds\":[\""+expectedLabel+"\"]}\n", string(body))
		w.WriteHeader(200)
		w.Write([]byte(`{}`))
	}))
}
