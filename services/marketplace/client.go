package marketplace

import (
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"

	"github.com/pkg/errors"
	"github.com/sitename/sitename/model_helper"
	"github.com/sitename/sitename/services/httpservice"
)

// Client is the programmatic interface to the marketplace server API.
type Client struct {
	address    string
	httpClient *http.Client
}

// NewClient creates a client to the marketplace server at the given address.
func NewClient(address string, httpService httpservice.HTTPService) (*Client, error) {
	var httpClient *http.Client
	addressURL, err := url.Parse(address)
	if err != nil {
		return nil, errors.Wrap(err, "failed to parse marketplace address")
	}
	if addressURL.Hostname() == "localhost" || addressURL.Hostname() == "127.0.0.1" {
		httpClient = httpService.MakeClient(true)
	} else {
		httpClient = httpService.MakeClient(false)
	}

	return &Client{
		address:    address,
		httpClient: httpClient,
	}, nil
}

// GetPlugins fetches the list of plugins from the configured server.
func (c *Client) GetPlugins(request *model_helper.MarketplacePluginFilter) ([]*model_helper.BaseMarketplacePlugin, error) {
	u, err := url.Parse(c.buildURL("/api/v1/plugins"))
	if err != nil {
		return nil, err
	}

	request.ApplyToURL(u)

	resp, err := c.doGet(u.String())
	if err != nil {
		return nil, err
	}
	defer closeBody(resp)

	switch resp.StatusCode {
	case http.StatusOK:
		return model_helper.BaseMarketplacePluginsFromReader(resp.Body)
	default:
		return nil, errors.Errorf("failed with status code %d", resp.StatusCode)
	}
}

func (c *Client) GetPlugin(filter *model_helper.MarketplacePluginFilter, pluginVersion string) (*model_helper.BaseMarketplacePlugin, error) {
	plugins, err := c.GetPlugins(filter)
	if err != nil {
		return nil, err
	}
	for _, plugin := range plugins {
		if plugin.Manifest.Version == pluginVersion {
			return plugin, nil
		}
	}
	return nil, errors.New("plugin not found")
}

// closeBody ensures the Body of an http.Response is properly closed.
func closeBody(r *http.Response) {
	if r.Body != nil {
		_, _ = io.Copy(ioutil.Discard, r.Body)
		_ = r.Body.Close()
	}
}

func (c *Client) buildURL(urlPath string, args ...any) string {
	return fmt.Sprintf("%s/%s", strings.TrimRight(c.address, "/"), strings.TrimLeft(fmt.Sprintf(urlPath, args...), "/"))
}

func (c *Client) doGet(u string) (*http.Response, error) {
	return c.httpClient.Get(u)
}
