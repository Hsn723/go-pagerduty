package pagerduty

import (
	"context"
	"fmt"
	"net/http"

	"github.com/google/go-querystring/query"
)

// Extension represents a single PagerDuty extension. These are addtional
// features to be used as part of the incident management process.
type Extension struct {
	APIObject
	Name             string      `json:"name"`
	EndpointURL      string      `json:"endpoint_url,omitempty"`
	ExtensionObjects []APIObject `json:"extension_objects"`
	ExtensionSchema  APIObject   `json:"extension_schema"`
	Config           interface{} `json:"config"`
}

// ListExtensionResponse represents the single response from the PagerDuty API
// when listing extensions.
type ListExtensionResponse struct {
	APIListObject
	Extensions []Extension `json:"extensions"`
}

// ListExtensionOptions are the options to use when listing extensions.
type ListExtensionOptions struct {
	APIListObject
	ExtensionObjectID string `url:"extension_object_id,omitempty"`
	ExtensionSchemaID string `url:"extension_schema_id,omitempty"`
	Query             string `url:"query,omitempty"`
}

// ListExtensions lists the extensions from the API.
//
// Deprecated: Use ListExtensionsWithContext instead.
func (c *Client) ListExtensions(o ListExtensionOptions) (*ListExtensionResponse, error) {
	return c.ListExtensionsWithContext(context.Background(), o)
}

// ListExtensionsWithContext lists the extensions from the API.
func (c *Client) ListExtensionsWithContext(ctx context.Context, o ListExtensionOptions) (*ListExtensionResponse, error) {
	v, err := query.Values(o)
	if err != nil {
		return nil, err
	}

	resp, err := c.get(ctx, "/extensions?"+v.Encode())
	if err != nil {
		return nil, err
	}

	var result ListExtensionResponse
	if err = c.decodeJSON(resp, &result); err != nil {
		return nil, err
	}

	return &result, nil
}

// CreateExtension creates a single extension.
//
// Deprecated: Use CreateExtensionWithContext instead.
func (c *Client) CreateExtension(e *Extension) (*Extension, error) {
	return c.CreateExtensionWithContext(context.Background(), e)
}

// CreateExtensionWithContext creates a single extension.
func (c *Client) CreateExtensionWithContext(ctx context.Context, e *Extension) (*Extension, error) {
	resp, err := c.post(ctx, "/extensions", e, nil)
	return getExtensionFromResponse(c, resp, err)
}

// DeleteExtension deletes an extension by its ID.
//
// Deprecated: Use DeleteExtensionWithContext instead.
func (c *Client) DeleteExtension(id string) error {
	return c.DeleteExtensionWithContext(context.Background(), id)
}

// DeleteExtensionWithContext deletes an extension by its ID.
func (c *Client) DeleteExtensionWithContext(ctx context.Context, id string) error {
	_, err := c.delete(ctx, "/extensions/"+id)
	return err
}

// GetExtension gets an extension by its ID.
//
// Deprecated: Use GetExtensionWithContext instead.
func (c *Client) GetExtension(id string) (*Extension, error) {
	return c.GetExtensionWithContext(context.Background(), id)
}

// GetExtensionWithContext gets an extension by its ID.
func (c *Client) GetExtensionWithContext(ctx context.Context, id string) (*Extension, error) {
	resp, err := c.get(ctx, "/extensions/"+id)
	return getExtensionFromResponse(c, resp, err)
}

// UpdateExtension updates an extension by its ID.
//
// Deprecated: Use UpdateExtensionWithContext instead.
func (c *Client) UpdateExtension(id string, e *Extension) (*Extension, error) {
	return c.UpdateExtensionWithContext(context.Background(), id, e)
}

// UpdateExtensionWithContext updates an extension by its ID.
func (c *Client) UpdateExtensionWithContext(ctx context.Context, id string, e *Extension) (*Extension, error) {
	resp, err := c.put(ctx, "/extensions/"+id, e, nil)
	return getExtensionFromResponse(c, resp, err)
}

func getExtensionFromResponse(c *Client, resp *http.Response, err error) (*Extension, error) {
	if err != nil {
		return nil, err
	}

	var target map[string]Extension
	if dErr := c.decodeJSON(resp, &target); dErr != nil {
		return nil, fmt.Errorf("Could not decode JSON response: %v", dErr)
	}

	const rootNode = "extension"

	t, nodeOK := target[rootNode]
	if !nodeOK {
		return nil, fmt.Errorf("JSON response does not have %s field", rootNode)
	}

	return &t, nil
}
