package api

import (
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

var client *CanvasClient
var domain string

func TestMain(m *testing.M) {
	client = NewClient("domain", "thisIsAToken")
	domain = "https://domain.instructure.com"
	returnCode := m.Run()

	os.Exit(returnCode)
}

func testServer() (*http.Client, *http.ServeMux, *httptest.Server) {
	mux := http.NewServeMux()
	server := httptest.NewServer(mux)
	client := &http.Client{}
	return client, mux, server
}

// func TestCanvasClient_GetUserProfile(t *testing.T) {
// 	defer gock.Off()

// 	expected := &User{
// 		ID:           1997,
// 		Name:         "Users Name",
// 		ShortName:    "Users Short Name",
// 		SortableName: "Sortable, Name",
// 		Title:        "Professor sortable name",
// 		Bio:          "Bio of a sortable name",
// 		PrimaryEmail: "sortable@name.email.com",
// 		LoginID:      "sortable911",
// 		SisUserID:    "sortable911",
// 		LtiUserID:    "Sortable911",
// 		AvatarURL:    "http://avatar.com",
// 		Calendar:     nil,
// 		TimeZone:     "USA USA",
// 		Locale:       "nil",
// 	}

// 	gock.New(domain).
// 		Get("/api/v1/users/1956/profile").
// 		Reply(200).
// 		JSON(expected)

// 	got, err := client.GetUserProfile(1956)
// 	assert.Nil(t, err)
// 	assert.Equal(t, expected, got)

// 	gock.New(domain).
// 		Get("/api/v1/users/1945/profile").
// 		Reply(404)

// 	expected = &User{}
// 	got, err = client.GetUserProfile(0)

// 	assert.Error(t, err)
// 	assert.Equal(t, expected, got)

// }

// func TestCanvasClient_GetDashboardPositions(t *testing.T) {
// 	defer gock.Off()

// 	expected := make(DashboardPositions)
// 	expected["course_16552"] = 0
// 	expected["course_16553"] = 4
// 	expected["course_16512"] = 2

// 	temp := &temporaryPositions{DashboardPositions: expected}

// 	gock.New(domain).
// 		Get("/api/v1/users/1945/dashboard_positions").
// 		Reply(200).
// 		JSON(temp)

// 	got, err := client.GetDashboardPositions(1945)

// 	assert.Nil(t, err)
// 	assert.Equal(t, &expected, got)

// }

// func TestCanvasClient_ActivityStream(t *testing.T) {
// 	defer gock.Off()

// 	expected := make(ActivityStream, 1)
// 	x := make(ActivityStreamItem)
// 	x["kappa"] = 123.0
// 	x["this is fine"] = "please"
// 	expected[0] = x
// 	gock.New(domain).
// 		Get("/api/v1/users/self/activity_stream").
// 		Reply(200).
// 		JSON(expected)

// 	got, err := client.GetActivityStream().Do()

// 	assert.Nil(t, err)
// 	assert.Equal(t, &expected, got)

// }
