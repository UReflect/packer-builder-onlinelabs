package onlinelabs

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
)

var (
	AccountURL, _ = url.Parse("https://account.cloud.online.net")
	APIURL, _     = url.Parse("https://api.cloud.online.net")
)

const (
	UAString = "packer-builder-onlinelabs (https://github.com/meatballhat/packer-builder-onlinelabs)"
)

// ClientInterface as documented at https://doc.cloud.online.net/api/
type ClientInterface interface {
	CreateKey(string, string) (string, error)
	DestroyKey(string) error
	CreateServer(string, string, string, []*Volume, []string) (*Server, error)
	GetServer(string) (*Server, error)
	PowerOnServer(string) error
	PowerOffServer(string) error
	DestroyServer(string) error
	CreateSnapshot(string, string, string) (*Snapshot, error)
	CreateImage(string, string, string, string) (*Image, error)
	DestroyImage(string) error
}

type Client struct {
	APIToken       string
	OrganizationID string
	AccountURL     *url.URL
	APIURL         *url.URL
}

func NewClient(apiToken, orgID string, accountURL, apiURL *url.URL) ClientInterface {
	return &Client{
		APIToken:       apiToken,
		OrganizationID: orgID,
		AccountURL:     accountURL,
		APIURL:         apiURL,
	}
}

func (c *Client) CreateKey(name, format string) (string, error) {
	// TODO: implement if/when available (?)
	return "", nil
}

func (c *Client) DestroyKey(keyID string) error {
	// TODO: implement if/when available (?)
	return nil
}

func (c *Client) CreateServer(name, org, image string, volumes []*Volume, tags []string) (*Server, error) {
	volMap := map[string]*Volume{}
	for i, vol := range volumes {
		volMap[fmt.Sprintf("%d", i)] = vol
	}

	body := &createServerParams{
		Name:         name,
		Image:        image,
		Organization: org,
		Volumes:      volMap,
		Tags:         tags,
	}

	resp, err := NewAPIRequest(c, "POST", "/servers", body)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode >= 300 {
		return nil, errFromResponse("server creation failed", resp)
	}

	srvBody := map[string]*Server{"server": &Server{}}
	err = json.NewDecoder(resp.Body).Decode(&srvBody)
	if err != nil {
		return nil, err
	}

	return srvBody["server"], nil
}

func (c *Client) GetServer(id string) (*Server, error) {
	path := fmt.Sprintf("/servers/%s", id)
	resp, err := NewAPIRequest(c, "GET", path, nil)

	if err != nil {
		return nil, err
	}
	if resp.StatusCode >= 300 {
		return nil, errFromResponse("fetching server failed", resp)
	}

	srvBody := map[string]*Server{"server": &Server{}}
	err = json.NewDecoder(resp.Body).Decode(&srvBody)
	if err != nil {
		return nil, err
	}

	return srvBody["server"], nil
}

func (c *Client) DestroyServer(id string) error {
	return nil
}

func (c *Client) PowerOnServer(id string) error {
	return c.sendAction(id, "poweron")
}

func (c *Client) PowerOffServer(id string) error {
	return c.sendAction(id, "poweroff")
}

func (c *Client) sendAction(id, action string) error {
	path := fmt.Sprintf("/servers/%s/action", id)
	resp, err := NewAPIRequest(c, "POST", path, map[string]string{"action": action})
	if err != nil {
		return err
	}
	if resp.StatusCode >= 300 {
		return errFromResponse(fmt.Sprintf("server %s failed", action), resp)
	}
	return nil
}

func (c *Client) CreateSnapshot(name, org, volumeID string) (*Snapshot, error) {
	return nil, nil
}

func (c *Client) CreateImage(org, name, arch, rootVolume string) (*Image, error) {
	return nil, nil
}

func (c *Client) DestroyImage(id string) error {
	return nil
}

func NewAPIRequest(c *Client, method, path string, body interface{}) (*http.Response, error) {
	bodyReader := bytes.NewReader([]byte(""))

	if body != nil {
		bodyBytes, err := json.Marshal(body)
		if err != nil {
			return nil, err
		}

		bodyReader = bytes.NewReader(bodyBytes)
	}

	url, err := url.Parse(c.APIURL.String())
	if err != nil {
		return nil, err
	}

	url.Path = path
	req, err := http.NewRequest(method, url.String(), bodyReader)
	if err != nil {
		return nil, err
	}

	req.Header.Set("X-Auth-Token", c.APIToken)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("User-Agent", UAString)
	return http.DefaultClient.Do(req)
}

func errFromResponse(message string, resp *http.Response) error {
	errText := ""
	b, err := ioutil.ReadAll(resp.Body)
	if err == nil {
		errText = string(b)
	}
	return fmt.Errorf("%s: %s %s", message, resp.Status, errText)
}
