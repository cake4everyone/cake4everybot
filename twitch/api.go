package twitch

import (
	"cake4everybot/webserver/oauth"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/spf13/viper"
)

func init() {

}

const (
	baseURL = "https://api.twitch.tv/helix"
)

var (
	clientID     string
	clientSecret string
	appToken     *oauth.Client
)

// Connect initializes the variables
func Connect() {
	clientID = viper.GetString("twitch.clientID")
	clientSecret = viper.GetString("twitch.clientSecret")

	appToken = oauth.New(
		"https://id.twitch.tv/oauth2/token",
		clientID,
		clientSecret,
		"",
	)
}

// GetStreamsByID gets all the streams matching the given user IDs of the streamers.
// Returns only the streams of those users that are broadcasting.
func GetStreamsByID(userIDs ...string) ([]*Stream, error) {
	if len(userIDs) == 0 {
		return []*Stream{}, nil
	}
	queryParams := map[string][]string{
		"user_id": userIDs,
		"first":   {"100"},
	}

	var streamData rawStreamData
	err := requestHelper(http.MethodGet, "/streams", queryParams, nil, &streamData)
	if err != nil {
		return []*Stream{}, fmt.Errorf("get streams by id: %v", err)
	}

	return streamData.Data, nil
}

// GetStreamsByName gets all the streams matching the given user login names of the streamers.
// Returns only the streams of those users that are broadcasting.
func GetStreamsByName(userLoginNames ...string) ([]*Stream, error) {
	if len(userLoginNames) == 0 {
		return []*Stream{}, nil
	}
	queryParams := map[string][]string{
		"user_login": userLoginNames,
		"first":      {"100"},
	}

	var streamData rawStreamData
	err := requestHelper(http.MethodGet, "/streams", queryParams, nil, &streamData)
	if err != nil {
		return []*Stream{}, fmt.Errorf("get streams by name: %v", err)
	}

	return streamData.Data, nil
}

// GetUsersByID gets all the Twitch users matching the given user IDs.
func GetUsersByID(userIDs ...string) ([]*User, error) {
	if len(userIDs) == 0 {
		return []*User{}, nil
	}
	queryParams := map[string][]string{
		"id": userIDs,
	}

	var streamData rawUserData
	err := requestHelper(http.MethodGet, "/users", queryParams, nil, &streamData)
	if err != nil {
		return []*User{}, fmt.Errorf("get users by id: %v", err)
	}

	return streamData.Data, nil
}

// GetUsersByName gets all the Twitch users matching the given user login names.
func GetUsersByName(userLoginNames ...string) ([]*User, error) {
	if len(userLoginNames) == 0 {
		return []*User{}, nil
	}
	queryParams := map[string][]string{
		"login": userLoginNames,
	}

	var streamData rawUserData
	err := requestHelper(http.MethodGet, "/users", queryParams, nil, &streamData)
	if err != nil {
		return []*User{}, fmt.Errorf("get users by name: %v", err)
	}

	return streamData.Data, nil
}

func requestHelper(method, endpoint string, queryParams map[string][]string, body io.Reader, result any) error {
	req, err := buildRequest(method, endpoint, queryParams, body)
	if err != nil {
		return err
	}
	return doRequest(req, result)
}

func buildRequest(method, endpoint string, queryParams map[string][]string, body io.Reader) (req *http.Request, err error) {
	req, err = http.NewRequest(method, baseURL+endpoint, body)
	if err != nil {
		return
	}

	var rawQueries []string
	for k, v := range queryParams {
		for _, v := range v {
			rawQueries = append(rawQueries, fmt.Sprintf("%s=%s", k, v))
		}
	}
	req.URL.RawQuery = strings.Join(rawQueries, "&")
	return
}

func doRequest(req *http.Request, result any) error {
	t, err := appToken.GenerateToken()
	if err != nil {
		return fmt.Errorf("generate token: %v", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+t)
	req.Header.Set("Client-Id", clientID)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("read response body: %v", err)
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("expected a 2xx status code, but got '%s': %s", resp.Status, body)
	}

	return json.Unmarshal(body, result)
}