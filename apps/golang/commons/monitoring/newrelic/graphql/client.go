package graphql

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"html/template"
	"io"
	"net/http"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/utr1903/opentelemetry-kubernetes-demo/apps/golang/commons/logger"
)

const (
	GRAPHQL_SUBSTITUTING_TEMPLATE_VARIABLES            = "substituting template variables"
	GRAPHQL_SUBSTITUTING_TEMPLATE_VARIABLES_HAS_FAILED = "substituting template variables has failed"
	GRAPHQL_EXECUTING_REQUEST                          = "executing request"
	GRAPHQL_EXECUTING_REQUEST_HAS_FAILED               = "executing request has failed"
	GRAPHQL_CREATING_PAYLOAD_HAS_FAILED                = "creating payload has failed"
	GRAPHQL_CREATING_HTTP_REQUEST_HAS_FAILED           = "creating http request has failed"
	GRAPHQL_PERFORMING_HTTP_REQUEST_HAS_FAILED         = "performing payload has failed"
	GRAPHQL_READING_HTTP_RESPONSE_BODY_HAS_FAILED      = "reading response body has failed"
	GRAPHQL_RESPONSE_HAS_RETURNED_NOT_OK_STATUS_CODE   = "response has returned not ok status code"
	GRAPHQL_PARSING_HTTP_RESPONSE_BODY_HAS_FAILED      = "parsing response body has failed"
)

type graphQlRequestPayload struct {
	Query string `json:"query"`
}

type GraphQlClient struct {
	logger                  *logger.Logger
	HhttpClient             *http.Client
	newrelicGraphQlEndpoint string
	newrelicUserApiKey      string
}

func New(
	logger *logger.Logger,
	newrelicGraphQlEndpoint string,
	newrelicUserApiKey string,
) *GraphQlClient {
	return &GraphQlClient{
		logger:                  logger,
		HhttpClient:             &http.Client{Timeout: time.Duration(30 * time.Second)},
		newrelicGraphQlEndpoint: newrelicGraphQlEndpoint,
		newrelicUserApiKey:      newrelicUserApiKey,
	}
}

func (c *GraphQlClient) Execute(
	ctx context.Context,
	queryTemplateName string,
	queryTemplate string,
	queryVariables any,
	result any,
) error {

	// Substitute variables within query
	query, err := c.substituteTemplateQuery(ctx, queryTemplateName, queryTemplate, queryVariables)
	if err != nil {
		return err
	}

	// Create payload
	payload, err := c.createPayload(ctx, query)
	if err != nil {
		return err
	}

	// Create request
	req, err := http.NewRequest(
		http.MethodPost,
		c.newrelicGraphQlEndpoint,
		payload,
	)
	if err != nil {
		c.logger.Log(logrus.ErrorLevel, ctx, "", GRAPHQL_CREATING_HTTP_REQUEST_HAS_FAILED)
		return err
	}

	// Add headers
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Api-Key", c.newrelicUserApiKey)

	// Perform HTTP request
	res, err := c.HhttpClient.Do(req)
	if err != nil {
		c.logger.Log(logrus.ErrorLevel, ctx, "", GRAPHQL_PERFORMING_HTTP_REQUEST_HAS_FAILED)
		return err
	}
	defer res.Body.Close()

	// Read HTTP response
	body, err := io.ReadAll(res.Body)
	if err != nil {
		c.logger.Log(logrus.ErrorLevel, ctx, "", GRAPHQL_READING_HTTP_RESPONSE_BODY_HAS_FAILED)
		return err
	}

	// Check if call was successful
	if res.StatusCode != http.StatusOK {
		c.logger.Log(logrus.ErrorLevel, ctx, "", GRAPHQL_RESPONSE_HAS_RETURNED_NOT_OK_STATUS_CODE+": "+string(body))
		return errors.New(GRAPHQL_RESPONSE_HAS_RETURNED_NOT_OK_STATUS_CODE)
	}

	err = json.Unmarshal(body, result)
	if err != nil {
		c.logger.Log(logrus.ErrorLevel, ctx, "", GRAPHQL_PARSING_HTTP_RESPONSE_BODY_HAS_FAILED)
		return err
	}

	return nil
}

func (c *GraphQlClient) substituteTemplateQuery(
	ctx context.Context,
	queryTemplateName string,
	queryTemplate string,
	queryVariables any,
) (
	*string,
	error,
) {
	// Parse query template
	c.logger.Log(logrus.InfoLevel, ctx, "", GRAPHQL_SUBSTITUTING_TEMPLATE_VARIABLES)

	t, err := template.New(queryTemplateName).Parse(queryTemplate)
	if err != nil {
		c.logger.Log(logrus.ErrorLevel, ctx, "", GRAPHQL_SUBSTITUTING_TEMPLATE_VARIABLES_HAS_FAILED)
		return nil, err
	}

	// Write substituted query template into buffer
	c.logger.Log(logrus.InfoLevel, ctx, "", GRAPHQL_EXECUTING_REQUEST)

	buf := new(bytes.Buffer)
	err = t.Execute(buf, queryVariables)
	if err != nil {
		c.logger.Log(logrus.ErrorLevel, ctx, "", GRAPHQL_SUBSTITUTING_TEMPLATE_VARIABLES_HAS_FAILED)
		return nil, err
	}

	// Return substituted query as string
	str := buf.String()
	return &str, nil
}

func (c *GraphQlClient) createPayload(
	ctx context.Context,
	query *string,
) (
	*bytes.Buffer,
	error,
) {

	// Create JSON data
	payload, err := json.Marshal(&graphQlRequestPayload{
		Query: *query,
	})
	if err != nil {
		c.logger.Log(logrus.ErrorLevel, ctx, "", GRAPHQL_CREATING_PAYLOAD_HAS_FAILED)
		return nil, err
	}
	return bytes.NewBuffer(payload), nil
}
