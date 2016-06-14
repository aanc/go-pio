package putio

import (
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"strconv"

	"github.com/antonholmquist/jason"
)

const putioAPIURL = "https://api.put.io/v2"

// Config structure is used to store put.io connection info
type Config struct {
	token string
}

func check(e error) {
	if e != nil {
		fmt.Printf("%s", e)
		panic(e)
	}
}

// SetToken sets the token
func (c *Config) SetToken(t string) {
	(*c).token = t
}

// Send sends a request to put.io API using the given token
func send(t string, uri string, arguments string) (*jason.Object, error) {
	url := putioAPIURL + uri + "?oauth_token=" + t + "&" + arguments
	response, err := http.Get(url)
	check(err)
	defer response.Body.Close()

	contents, err := ioutil.ReadAll(response.Body)
	check(err)

	v, err := jason.NewObjectFromBytes([]byte(contents))

	status, _ := v.GetString("status")
	if status == "ERROR" {
		fmt.Println(`
Error during put.io API access, please check your internet connectivity and your
token configuration. See --help for more info.`)
		os.Exit(1)
	}

	return v, err
}

// List lists files in a given directory ID
func (c *Config) List(fileID int64) (*jason.Object, error) {
	v, err := send(c.token, "/files/list", "parent_id="+strconv.FormatInt(fileID, 10))
	return v, err
}

// AccountInfo returns a json object describing the account information
// as described at https://api.put.io/v2/docs/account.html#info
func (c *Config) AccountInfo() (*jason.Object, error) {
	v, err := send(c.token, "/account/info", "")
	return v, err
}

// Transfers returns a json object containing transfers list, as described
// at https://api.put.io/v2/docs/transfers.html#list
func (c *Config) Transfers() (*jason.Object, error) {
	v, err := send(c.token, "/transfers/list", "")
	return v, err
}

// GetDownloadLink returns the download URL of a given fileID as a string
func (c *Config) GetDownloadLink(fileID int64) (string, error) {
	reqURL := putioAPIURL + "/files/" + strconv.FormatInt(fileID, 10) + "/download" + "?oauth_token=" + c.token

	// put.io returns a first response containing a redirect URL when a download
	// is requested. We need to catch that redirect, as we only want the URL, not
	// the whole file.

	// Custom redirect error
	var RedirectAttemptedError = errors.New("redirect")

	// Custom http client, so we can use the redirect error
	client := &http.Client{
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return RedirectAttemptedError
		},
	}

	// Requesting the file download using the custom client
	response, err := client.Head(reqURL)

	// Checking if we get the redirect error
	if urlError, ok := err.(*url.Error); ok && urlError.Err == RedirectAttemptedError {
		err = nil
	}
	check(err)
	defer response.Body.Close()

	// Extracting download link from headers
	return response.Header.Get("Location"), err
}
