package streamelements

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

// GetChannels returns a list of all channels the current user has access to. The current user is
// defined by the used bearer token.
func (se *Streamelements) GetChannels() ([]*Channel1, error) {
	var channels = make([]*Channel1, 0)

	r, err := se.doReq(http.MethodGet, "/users/channels", []byte{}, nil)
	if err != nil {
		return channels, err
	}

	data, err := io.ReadAll(r.Body)
	if err != nil {
		return channels, err
	}

	if r.StatusCode != 200 {
		return channels, fmt.Errorf("wrong status code, expected 200 but got %d! Response data: %s", r.StatusCode, string(data))
	}

	err = json.Unmarshal(data, &channels)
	return channels, err
}

// GetChannelDetails returns details for the given channelID (streamelements ID). This only works
// when the current user has access to this channel, otherwise a 403 is returned.
//
// NOTE: Some documentation is missing from streamelements. It appears that this endpoint only
// only returns the current users channel.
func (se *Streamelements) GetChannelDetails(channelID string) (*ChannelDetails, error) {
	r, err := se.doReq(
		http.MethodGet,
		fmt.Sprintf("/channels/%s/details", channelID),
		[]byte{},
		nil,
	)
	if err != nil {
		return nil, err
	}

	data, err := io.ReadAll(r.Body)
	if err != nil {
		return nil, err
	}

	if r.StatusCode != 200 {
		return nil, fmt.Errorf("wrong status code, expected 200 but got %d! Response data: %s", r.StatusCode, string(data))
	}

	c := &ChannelDetails{}
	err = json.Unmarshal(data, c)
	return c, err
}

// GetChannel returns basic details for the given channel. This endpoint does not requiere access to
// the requested channel.
//
// "channel" can be either a streamelements channel ID or the username of a channel.
func (se *Streamelements) GetChannel(channel string) (*SimpleChannelDetails, error) {
	r, err := se.doReq(
		http.MethodGet,
		fmt.Sprintf("/channels/%s", channel),
		[]byte{},
		nil,
	)
	if err != nil {
		return nil, err
	}

	data, err := io.ReadAll(r.Body)
	if err != nil {
		return nil, err
	}

	if r.StatusCode != 200 {
		return nil, fmt.Errorf("wrong status code, expected 200 but got %d! Response data: %s", r.StatusCode, string(data))
	}

	c := &SimpleChannelDetails{}
	err = json.Unmarshal(data, c)
	return c, err
}

// AddPoints modifies (add or remove) the streamelements points for a user. When amount == 0
// AddPoints is a no-op.
//
//	channelID // the streamelements ID of the channel to add the points to
//	username  // the username to modify
//	amount    // the amount to modify, amount > 0 adds the points, amount < 0 removes the points.
func (se *Streamelements) AddPoints(channelID, username string, amount int) error {
	if amount == 0 {
		return nil
	}
	r, err := se.doReq(
		http.MethodPut,
		fmt.Sprintf("/points/%s/%s/%d", channelID, username, amount),
		[]byte("{}"),
		map[string]string{"Content-Type": "application/json"},
	)
	if err != nil {
		return err
	}

	if r.StatusCode == 200 {
		return nil
	}

	data, err := io.ReadAll(r.Body)
	if err != nil {
		return fmt.Errorf("wrong status code, expected 200 but got %d! But also failed to read data: %v", r.StatusCode, err)
	}
	return fmt.Errorf("wrong status code, expected 200 but got %d! Response data: %s", r.StatusCode, string(data))
}

// GetPoints returns the current streamelements points for a user.
//
//	channelID // the streamelements ID of the channel to get the points from
//	username  // the username to fetch
func (se *Streamelements) GetPoints(channelID, username string) (*UserPoints, error) {
	r, err := se.doReq(
		http.MethodGet,
		fmt.Sprintf("/points/%s/%s", channelID, username),
		[]byte{},
		nil,
	)
	if err != nil {
		return nil, err
	}

	data, err := io.ReadAll(r.Body)
	if err != nil {
		return nil, err
	}

	if r.StatusCode != 200 {
		return nil, fmt.Errorf("wrong status code, expected 200 but got %d! Response data: %v", r.StatusCode, string(data))
	}

	up := &UserPoints{}
	err = json.Unmarshal(data, up)
	return up, err
}
