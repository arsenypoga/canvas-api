package api

import (
	"errors"
	"fmt"
	"net/url"
	"strconv"
)

// Users is a array of a User
type Users []User

// User is user profile
type User struct {
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

// AccountUsersOptions is an interface for the lookup of account users
type AccountUsersOptions struct {
	searchTerm     string
	enrollmentType string
	sort           string
	order          string
	err            []error
}

// AccountUsersOption is an adapter for generating options
type AccountUsersOption func(*AccountUsersOptions)

// ActivityStreamOptions is an interface for the lookup of Activity Stream
type ActivityStreamOptions struct {
	onlyActiveCourses bool
}

// ActivityStream is Users activity feed
type ActivityStream struct {
	DiscussionTopics  []DiscussionTopic
	Announcements     []Announcement
	Conversations     []Conversation
	Messages          []Message
	Conferences       []Conference
	Collaborations    []Collaboration
	AssesmentRequests []AssesmentRequest
}

// ActivityStreamOption is an adapter for generating options
type ActivityStreamOption func(*ActivityStreamOptions)

type activityStreamPlaceholder []activityStreamPlaceholderItem
type activityStreamPlaceholderItem map[string]interface{}

// Message is a ActivityStream message
type Message struct {
	ID                   int64
	NotificationCategory string

	CreatedAt string
	UpdatedAt string
	Title     string
	Message   string
	ReadState bool
	CourseID  int64
	GroupID   int64
	HTMLURL   string
}

// DiscussionTopic is a ActivityStream discussion
type DiscussionTopic struct {
	ID                         int64
	TotalRootDiscussionEntries int64
	RequireInitialPost         bool

	CreatedAt string
	UpdatedAt string
	Title     string
	Message   string
	ReadState bool
	CourseID  int64
	GroupID   int64
	HTMLURL   string

	UserHasPosted         interface{}
	RootDiscussionEntries interface{}
}

// Announcement is a ActivityStream announcement
type Announcement struct {
	ID                         int64
	TotalRootDiscussionEntries int64
	ContextType                string
	RequireInitialPost         bool
	CreatedAt                  string
	UpdatedAt                  string
	Title                      string
	Message                    string
	ReadState                  bool
	CourseID                   int64
	GroupID                    int64
	HTMLURL                    string
	UserHasPosted              interface{}
	RootDiscussionEntries      interface{}
}

type Conversation struct {
	ID               int64
	Private          bool
	ParticipantCount int64

	CreatedAt      string
	UpdatedAt      string
	Title          string
	LatestMessages interface{}
	ReadState      bool
	CourseID       int64
	GroupID        int64
	HTMLURL        string
}

type Conference struct {
	ID int64

	CreatedAt string
	UpdatedAt string
	Title     string
	Message   string
	ReadState bool
	CourseID  int64
	GroupID   int64
	HTMLURL   string
}

type Collaboration struct {
	ID int64

	CreatedAt string
	UpdatedAt string
	Title     string
	Message   string
	ReadState bool
	CourseID  int64
	GroupID   int64
	HTMLURL   string
}

type AssesmentRequest struct {
	ID int64

	CreatedAt string
	UpdatedAt string
	Title     string
	Message   string
	ReadState bool
	CourseID  int64
	GroupID   int64
	HTMLURL   string
}

func SearchTerm(searchTerm string) AccountUsersOption {
	return func(auo *AccountUsersOptions) {
		auo.searchTerm = searchTerm
	}
}

func EnrollmentType(enrollmentType string) AccountUsersOption {
	return func(auo *AccountUsersOptions) {
		auo.enrollmentType = enrollmentType
	}
}

func Sort(sort string) AccountUsersOption {
	if sort != "username" && sort != "email" && sort != "sis_id" || sort != "last_login" && sort != "" {
		return func(auo *AccountUsersOptions) {
			auo.err = append(auo.err, errors.New("keyword sort can be only one of: 'username' | 'email' | 'sis_id' | 'last_login'"))
		}
	}
	return func(auo *AccountUsersOptions) {
		auo.sort = sort
	}
}

func Order(order string) AccountUsersOption {
	if order != "asc" && order != "desc" && order != "" {
		return func(auo *AccountUsersOptions) {
			auo.err = append(auo.err, errors.New("keyword order can be only one of: 'asc' | 'desc'"))
		}
	}
	return func(auo *AccountUsersOptions) {
		auo.order = order
	}
}

func (c *CanvasClient) GetAccountUsers(setters ...AccountUsersOption) (Users, error) {
	args := &AccountUsersOptions{
		searchTerm:     "",
		enrollmentType: "",
		sort:           "",
		order:          "",
		err:            nil,
	}
	u := make(Users, 0)
	for _, setter := range setters {
		setter(args)
	}

	if len(args.err) != 0 {
		return u, args.err[0]
	}

	parsedURL, err := url.Parse(fmt.Sprintf("%s/api/v1/accounts/self/users", c.ClientURL()))

	if err != nil {
		return u, err
	}

	q := parsedURL.Query()

	q.Add("search_term", args.searchTerm)
	q.Add("enrollment_type", args.enrollmentType)
	q.Add("sort", args.sort)
	q.Add("order", args.order)

	parsedURL.RawQuery = q.Encode()

	err = c.getJSON(parsedURL.String(), u)

	if err != nil {
		return u, err
	}
	return u, nil
}

// TemporaryPositions maps dasboard class name with the position index
type temporaryPositions struct {
	DashboardPositions map[string]int `json:"dashboard_positions"`
}

// DashboardPositions is a cleaned up version of a Temporary positions without a single key
type DashboardPositions map[string]int

// GetUserProfile returns user profile with the given profileID
func (c *CanvasClient) GetUserProfile(userID int64) (*User, error) {
	profile := User{}

	requestURL := fmt.Sprintf("%s/api/v1/users/%d/profile", c.ClientURL(), userID)
	err := c.getJSON(requestURL, &profile)

	if err != nil {
		return &profile, err
	}

	return &profile, nil
}

// GetDashboardPositions returns dashboard positions for a user
func (c *CanvasClient) GetDashboardPositions(userID int64) (*DashboardPositions, error) {
	temp := temporaryPositions{}
	d := make(DashboardPositions)

	requestURL := fmt.Sprintf("%s/api/v1/users/%d/dashboard_positions", c.ClientURL(), userID)
	err := c.getJSON(requestURL, &temp)

	if err != nil {
		return &d, err
	}

	d = temp.DashboardPositions

	return &d, nil
}

func OnlyActiveUsers(active bool) ActivityStreamOptions {
	return ActivityStreamOptions{
		onlyActiveCourses: active,
	}
}

func (c *CanvasClient) GetActivityStream(setters ...ActivityStreamOption) (*ActivityStream, error) {
	args := &ActivityStreamOptions{
		onlyActiveCourses: false,
	}
	s := make(activityStreamPlaceholder, 0)
	stream := ActivityStream{}
	for _, setter := range setters {
		setter(args)
	}

	parsedURL, err := url.Parse(fmt.Sprintf("%s/api/v1/users/self/activity_stream", c.ClientURL()))

	if err != nil {
		return &stream, err
	}

	q := parsedURL.Query()

	q.Add("order", strconv.FormatBool(args.onlyActiveCourses))

	parsedURL.RawQuery = q.Encode()

	err = c.getJSON(parsedURL.String(), &s)

	if err != nil {
		return &stream, err
	}

	return activityStreamFromPlaceholder(&s), nil
}

func activityStreamFromPlaceholder(placeholder *activityStreamPlaceholder) *ActivityStream {
	stream := ActivityStream{
		Announcements:    make([]Announcement, 0),
		DiscussionTopics: make([]DiscussionTopic, 0),
		Conversations:    make([]Conversation, 0),
		Messages:         make([]Message, 0),
		Conferences:      make([]Conference, 0),
	}
	for _, item := range *placeholder {
		groupID, _ := item["group_id"].(float64)
		if item["type"] == "Announcement" {
			a := Announcement{
				ID:                         int64(item["announcement_id"].(float64)),
				TotalRootDiscussionEntries: int64(item["total_root_discussion_entries"].(float64)),
				RequireInitialPost:         item["require_initial_post"].(bool),
				UserHasPosted:              item["user_has_posted"],
				RootDiscussionEntries:      item["root_discussion_entries"],
				ContextType:                item["context_type"].(string),
				CreatedAt:                  item["created_at"].(string),
				UpdatedAt:                  item["updated_at"].(string),
				Title:                      item["title"].(string),
				Message:                    item["message"].(string),
				ReadState:                  item["read_state"].(bool),
				CourseID:                   int64(item["course_id"].(float64)),
				GroupID:                    int64(groupID),
				HTMLURL:                    item["html_url"].(string),
			}
			stream.Announcements = append(stream.Announcements, a)
		} else if item["type"] == "DiscussionTopic" {
			groupID, _ := item["group_id"].(float64)

			d := DiscussionTopic{
				ID:                         int64(item["discussion_topic_id"].(float64)),
				TotalRootDiscussionEntries: int64(item["total_root_discussion_entries"].(float64)),
				RequireInitialPost:         item["require_initial_post"].(bool),
				CreatedAt:                  item["created_at"].(string),
				UpdatedAt:                  item["updated_at"].(string),
				Title:                      item["title"].(string),
				Message:                    item["message"].(string),
				ReadState:                  item["read_state"].(bool),
				CourseID:                   int64(item["course_id"].(float64)),
				GroupID:                    int64(groupID),
				HTMLURL:                    item["html_url"].(string),
			}
			stream.DiscussionTopics = append(stream.DiscussionTopics, d)
		} else if item["type"] == "Conversation" {
			c := Conversation{
				ID:               int64(item["id"].(float64)),
				CreatedAt:        item["created_at"].(string),
				UpdatedAt:        item["updated_at"].(string),
				Title:            item["title"].(string),
				LatestMessages:   item["latest_messages"],
				ReadState:        item["read_state"].(bool),
				CourseID:         int64(item["course_id"].(float64)),
				GroupID:          int64(groupID),
				HTMLURL:          item["html_url"].(string),
				Private:          item["private"].(bool),
				ParticipantCount: int64(item["participant_count"].(float64)),
			}
			stream.Conversations = append(stream.Conversations, c)
		} else if item["type"] == "Message" {
			m := Message{
				ID:                   int64(item["id"].(float64)),
				CreatedAt:            item["created_at"].(string),
				UpdatedAt:            item["updated_at"].(string),
				Title:                item["title"].(string),
				Message:              item["message"].(string),
				ReadState:            item["read_state"].(bool),
				CourseID:             int64(item["course_id"].(float64)),
				GroupID:              int64(groupID),
				HTMLURL:              item["html_url"].(string),
				NotificationCategory: item["notification_category"].(string),
			}
			stream.Messages = append(stream.Messages, m)
		} else if item["type"] == "Conference" {
			c := Conference{
				ID:        int64(item["conference_id"].(float64)),
				CreatedAt: item["created_at"].(string),
				UpdatedAt: item["updated_at"].(string),
				Title:     item["title"].(string),
				Message:   item["message"].(string),
				ReadState: item["read_state"].(bool),
				CourseID:  int64(item["course_id"].(float64)),
				GroupID:   int64(groupID),
				HTMLURL:   item["html_url"].(string),
			}
			stream.Conferences = append(stream.Conferences, c)
		} else if item["type"] == "Submission" {
			//TODO: Add this when the Sumbission interface is finished
		} else if item["type"] == "Collaboration" {
			c := Collaboration{
				ID:        int64(item["conference_id"].(float64)),
				CreatedAt: item["created_at"].(string),
				UpdatedAt: item["updated_at"].(string),
				Title:     item["title"].(string),
				Message:   item["message"].(string),
				ReadState: item["read_state"].(bool),
				CourseID:  int64(item["course_id"].(float64)),
				GroupID:   int64(groupID),
				HTMLURL:   item["html_url"].(string),
			}
			stream.Collaborations = append(stream.Collaborations, c)

		} else if item["type"] == "AssesmentRequest" {
			a := AssesmentRequest{
				ID:        int64(item["conference_id"].(float64)),
				CreatedAt: item["created_at"].(string),
				UpdatedAt: item["updated_at"].(string),
				Title:     item["title"].(string),
				Message:   item["message"].(string),
				ReadState: item["read_state"].(bool),
				CourseID:  int64(item["course_id"].(float64)),
				GroupID:   int64(groupID),
				HTMLURL:   item["html_url"].(string),
			}
			stream.AssesmentRequests = append(stream.AssesmentRequests, a)
		}
	}
	return &stream
}
