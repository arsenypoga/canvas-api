package api

import (
	"fmt"
)

// Profile is user profile
type Profile struct {
	ID           int64             `json:"id"`
	Name         string            `json:"name"`
	ShortName    string            `json:"short_name"`
	SortableName string            `json:"sortable_name"`
	Title        string            `json:"title"`
	Bio          string            `json:"bio"`
	PrimaryEmail string            `json:"primary_email"`
	LoginID      string            `json:"login_id"`
	SisUserID    string            `json:"sis_user_id"`
	LtiUserID    string            `json:"lti_user_id"`
	AvatarURL    string            `json:"avatar_url"`
	Calendar     map[string]string `json:"calendar"`
	TimeZone     string            `json:"time_zone"`
	Locale       string            `json:"locale"`
}

// DashboardPositions maps dasboard class name with the position index
type DashboardPositions struct {
	DashboardPositions map[string]int `json:"dashboard_positions"`
}

// GetUserProfile returns user profile with the given profileID
func (c *CanvasClient) GetUserProfile(profileID int64) (*Profile, error) {
	profile := Profile{}

	requestURL := fmt.Sprintf("%s/api/v1/users/%d/profile", c.ClientURL(), profileID)
	err := c.getJSON(requestURL, &profile)

	if err != nil {
		panic(err)
	}

	return &profile, nil
}

// GetDashboardPositions returns dashboard positions for a user
func (c *CanvasClient) GetDashboardPositions(userID int64) (*DashboardPositions, error) {
	dp := DashboardPositions{}

	requestURL := fmt.Sprintf("%s/api/v1/users/%d/dashboard_positions", c.ClientURL(), userID)
	err := c.getJSON(requestURL, &dp)

	if err != nil {
		return &dp, err
	}

	return &dp, nil
}
