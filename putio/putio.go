package putio

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"

	"github.com/antonholmquist/jason"
)

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
	url := "https://api.put.io/v2" + uri + "?oauth_token=" + t + "&" + arguments
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
