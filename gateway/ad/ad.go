/*
 * SPDX-License-Identifier: Apache-2.0
 *
 * The OpenSearch Contributors require contributions made to
 * this file be licensed under the Apache-2.0 license or a
 * compatible open source license.
 *
 * Modifications Copyright OpenSearch Contributors. See
 * GitHub history for details.
 */
/*
 * Copyright 2020 Amazon.com, Inc. or its affiliates. All Rights Reserved.
 *
 * Licensed under the Apache License, Version 2.0 (the "License").
 * You may not use this file except in compliance with the License.
 * A copy of the License is located at
 *
 *        http://www.apache.org/licenses/LICENSE-2.0
 *
 * or in the "license" file accompanying this file. This file is distributed
 * on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either
 * express or implied. See the License for the specific language governing
 * permissions and limitations under the License.
 */

package ad

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"opensearch-cli/client"
	"opensearch-cli/entity"
	gw "opensearch-cli/gateway"
	"opensearch-cli/mapper"
)

const (
	baseURL           = "_plugins/_anomaly_detection/detectors"
	startURLTemplate  = baseURL + "/%s/" + "_start"
	stopURLTemplate   = baseURL + "/%s/" + "_stop"
	searchURLTemplate = baseURL + "/_search"
	deleteURLTemplate = baseURL + "/%s"
	getURLTemplate    = baseURL + "/%s"
	updateURLTemplate = baseURL + "/%s"
)

//go:generate go run -mod=mod github.com/golang/mock/mockgen  -destination=mocks/mock_ad.go -package=mocks . Gateway

// Gateway interface to AD Plugin
type Gateway interface {
	CreateDetector(context.Context, interface{}) ([]byte, error)
	StartDetector(context.Context, string) error
	StopDetector(context.Context, string) (*string, error)
	DeleteDetector(context.Context, string) error
	SearchDetector(context.Context, interface{}) ([]byte, error)
	GetDetector(context.Context, string) ([]byte, error)
	UpdateDetector(context.Context, string, interface{}) error
}

type gateway struct {
	gw.HTTPGateway
}

// New creates new Gateway instance
func New(c *client.Client, p *entity.Profile) (Gateway, error) {
	g, err := gw.NewHTTPGateway(c, p)
	if err != nil {
		return nil, err
	}
	return &gateway{*g}, nil
}

func (g *gateway) buildCreateURL() (*url.URL, error) {
	endpoint, err := gw.GetValidEndpoint(g.Profile)
	if err != nil {
		return nil, err
	}
	endpoint.Path = baseURL
	return endpoint, nil
}

/*CreateDetector Creates an anomaly detector job.
It calls http request: POST _plugins/_anomaly_detection/detectors
Sample Input:
{
 "name": "test-detector",
 "description": "Test detector",
 "time_field": "timestamp",
 "indices": [
   "order*"
 ],
 "feature_attributes": [
   {
     "feature_name": "total_order",
     "feature_enabled": true,
     "aggregation_query": {
       "total_order": {
         "sum": {
           "field": "value"
         }
       }
     }
   }
 ],
 "filter_query": {
   "bool": {
     "filter": [
       {
         "exists": {
           "field": "value",
           "boost": 1
         }
       }
     ],
     "adjust_pure_negative": true,
     "boost": 1
   }
 },
 "detection_interval": {
   "period": {
     "interval": 1,
     "unit": "Minutes"
   }
 },
 "window_delay": {
   "period": {
     "interval": 1,
     "unit": "Minutes"
   }
 }
}*/
func (g *gateway) CreateDetector(ctx context.Context, payload interface{}) ([]byte, error) {
	createURL, err := g.buildCreateURL()
	if err != nil {
		return nil, err
	}
	detectorRequest, err := g.BuildRequest(ctx, http.MethodPost, payload, createURL.String(), gw.GetDefaultHeaders())
	if err != nil {
		return nil, err
	}
	response, err := g.Call(detectorRequest, http.StatusCreated)
	if err != nil {
		return nil, err
	}
	return response, nil
}

func (g *gateway) buildStartURL(ID string) (*url.URL, error) {
	endpoint, err := gw.GetValidEndpoint(g.Profile)
	if err != nil {
		return nil, err
	}
	endpoint.Path = fmt.Sprintf(startURLTemplate, ID)
	return endpoint, nil
}

// StartDetector Starts an anomaly detector job.
// It calls http request: POST _plugins/_anomaly_detection/detectors/<detectorId>/_start
func (g *gateway) StartDetector(ctx context.Context, ID string) error {
	startURL, err := g.buildStartURL(ID)
	if err != nil {
		return err
	}
	detectorRequest, err := g.BuildRequest(ctx, http.MethodPost, "", startURL.String(), gw.GetDefaultHeaders())
	if err != nil {
		return err
	}
	_, err = g.Call(detectorRequest, http.StatusOK)
	if err != nil {
		return err
	}
	return nil
}

func (g *gateway) buildStopURL(ID string) (*url.URL, error) {
	endpoint, err := gw.GetValidEndpoint(g.Profile)
	if err != nil {
		return nil, err
	}
	endpoint.Path = fmt.Sprintf(stopURLTemplate, ID)
	return endpoint, nil
}

// StopDetector Stops an anomaly detector job.
// It calls http request: POST _plugins/_anomaly_detection/detectors/<detectorId>/_stop
func (g *gateway) StopDetector(ctx context.Context, ID string) (*string, error) {
	stopURL, err := g.buildStopURL(ID)
	if err != nil {
		return nil, err
	}
	detectorRequest, err := g.BuildRequest(ctx, http.MethodPost, "", stopURL.String(), gw.GetDefaultHeaders())
	if err != nil {
		return nil, err
	}
	res, err := g.Call(detectorRequest, http.StatusOK)
	if err != nil {
		return nil, err
	}
	return mapper.StringToStringPtr(fmt.Sprintf("%s", res)), nil
}

func (g *gateway) buildSearchURL() (*url.URL, error) {
	endpoint, err := gw.GetValidEndpoint(g.Profile)
	if err != nil {
		return nil, err
	}
	endpoint.Path = searchURLTemplate
	return endpoint, nil
}

/*SearchDetector Returns all anomaly detectors for a search query.
It calls http request: POST _plugins/_anomaly_detection/detectors/_search
sample input
Sample Input:
{
 "query": {
   "match": {
     "name": "test-detector"
   }
 }*/
func (g *gateway) SearchDetector(ctx context.Context, payload interface{}) ([]byte, error) {
	searchURL, err := g.buildSearchURL()
	if err != nil {
		return nil, err
	}
	searchRequest, err := g.BuildRequest(ctx, http.MethodPost, payload, searchURL.String(), gw.GetDefaultHeaders())
	if err != nil {
		return nil, err
	}
	response, err := g.Call(searchRequest, http.StatusOK)
	if err != nil {
		return nil, err
	}
	return response, nil
}

func (g *gateway) buildDeleteURL(ID string) (*url.URL, error) {
	endpoint, err := gw.GetValidEndpoint(g.Profile)
	if err != nil {
		return nil, err
	}
	endpoint.Path = fmt.Sprintf(deleteURLTemplate, ID)
	return endpoint, nil
}

// DeleteDetector Deletes a detector based on the detector_id.
// It calls http request: DELETE _plugins/_anomaly_detection/detectors/<detectorId>
func (g *gateway) DeleteDetector(ctx context.Context, ID string) error {
	deleteURL, err := g.buildDeleteURL(ID)
	if err != nil {
		return err
	}
	detectorRequest, err := g.BuildRequest(ctx, http.MethodDelete, "", deleteURL.String(), gw.GetDefaultHeaders())
	if err != nil {
		return err
	}
	_, err = g.Call(detectorRequest, http.StatusOK)
	if err != nil {
		return err
	}
	return nil
}

func (g *gateway) buildGetURL(ID string) (*url.URL, error) {
	endpoint, err := gw.GetValidEndpoint(g.Profile)
	if err != nil {
		return nil, err
	}
	endpoint.Path = fmt.Sprintf(getURLTemplate, ID)
	return endpoint, nil
}

// GetDetector Returns all information about a detector based on the detector_id.
// It calls http request: GET _plugins/_anomaly_detection/detectors/<detectorId>
func (g *gateway) GetDetector(ctx context.Context, ID string) ([]byte, error) {
	getURL, err := g.buildGetURL(ID)
	if err != nil {
		return nil, err
	}
	detectorRequest, err := g.BuildRequest(ctx, http.MethodGet, "", getURL.String(), gw.GetDefaultHeaders())
	if err != nil {
		return nil, err
	}
	response, err := g.Call(detectorRequest, http.StatusOK)
	if err != nil {
		return nil, err
	}
	return response, nil
}

func (g *gateway) buildUpdateURL(ID string) (*url.URL, error) {
	endpoint, err := gw.GetValidEndpoint(g.Profile)
	if err != nil {
		return nil, err
	}
	endpoint.Path = fmt.Sprintf(updateURLTemplate, ID)
	return endpoint, nil
}

/*UpdateDetector Updates a detector with any changes, including the description or adding or removing of features.
It calls http request: PUT _plugins/_anomaly_detection/detectors/<detectorId>
Sample Input:
{
 "name": "test-detector",
 "description": "Test detector",
 "time_field": "timestamp",
 "indices": [
   "order*"
 ],
 "feature_attributes": [
   {
     "feature_name": "total_order",
     "feature_enabled": true,
     "aggregation_query": {
       "total_order": {
         "sum": {
           "field": "value"
         }
       }
     }
   }
 ],
 "filter_query": {
   "bool": {
     "filter": [
       {
         "exists": {
           "field": "value",
           "boost": 1
         }
       }
     ],
     "adjust_pure_negative": true,
     "boost": 1
   }
 },
 "detection_interval": {
   "period": {
     "interval": 10,
     "unit": "Minutes"
   }
 },
 "window_delay": {
   "period": {
     "interval": 1,
     "unit": "Minutes"
   }
 }
}*/
func (g *gateway) UpdateDetector(ctx context.Context, ID string, payload interface{}) error {
	updateURL, err := g.buildUpdateURL(ID)
	if err != nil {
		return err
	}
	detectorRequest, err := g.BuildRequest(ctx, http.MethodPut, payload, updateURL.String(), gw.GetDefaultHeaders())
	if err != nil {
		return err
	}
	_, err = g.Call(detectorRequest, http.StatusOK)
	if err != nil {
		return err
	}
	return nil
}
