package goreplica

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"strings"
)

// Client : This is the Client implementation of GoReplica
type Client struct {
	ServiceEndpoint string       `json:"service_endpoint,omitempty"`
	ClientID        string       `json:"client_id,omitempty"`
	ClientSecret    string       `json:"secret,omitempty"`
	AccessToken     string       `json:"access_token,omitempty"`
	HTTPClient      *http.Client `json:"http_client,omitempty"`
}

// AuthResponse : Structure of the expected response from the authorization call.
type AuthResponse struct {
	AccessToken  string `json:"access_token,omitempty"`
	RefreshToken string `json:"refresh_token,omitempty"`
}

// VoiceResponse : Structure of the expected response from the call to get a listed of available voices.
type VoiceResponse struct {
	UUID string `json:"uuid,omitempty"`
	Name string `json:"name,omitempty"`
}

// SpeechResponse : Structure of the expected response from a call to get speech generated.
type SpeechResponse struct {
	UUID       string            `json:"uuid,omitempty"`
	Quality    string            `json:"quality,omitempty"`
	Duration   float32           `json:"duration,omitempty"`
	SpeakerID  string            `json:"speaker_id,omitempty"`
	Text       string            `json:"txt,omitempty"`
	BitRate    int               `json:"bit_rate,omitempty"`
	SampleRate int               `json:"sample_rate,omitempty"`
	Extension  string            `json:"extension,omitempty"`
	Extensions []string          `json:"extensions,omitempty"`
	URL        string            `json:"url,omitempty"`
	URLs       map[string]string `json:"urls,omitempty"`
}

// UnauthorizedResponse : Structure of a 401 Unauthorized response.
type UnauthorizedResponse struct {
	Reasons   []string `json:"reasons,omitempty"`
	Exception string   `json:"exception,omitempty"`
}

// BadRequestResponse : Structure of a 400 Bad Request Response.
type BadRequestResponse struct {
	ErrorCode int    `json:"error_code,omitempty"`
	Error     string `json:"error,omitempty"`
}

// Auth : Authentication endpoint. Returns a JSON response with JWT token which must be used to make calls to other endpoints.
func (replicaClient *Client) Auth() error {

	data := url.Values{}
	data.Add(CLIENTID, replicaClient.ClientID)
	data.Add(CLIENTSECRET, replicaClient.ClientSecret)

	u, _ := url.ParseRequestURI(replicaClient.ServiceEndpoint)
	u.Path = "/auth/"
	u.RawQuery = data.Encode()
	urlStr := u.String()

	request, err := http.NewRequest(http.MethodPost, urlStr, strings.NewReader(data.Encode()))
	if err != nil {
		return err
	}

	request.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	request.Header.Add("Content-Length", strconv.Itoa(len(data.Encode())))
	request.Close = true
	response, requestErr := replicaClient.HTTPClient.Do(request)
	if requestErr != nil {
		return requestErr
	}

	switch response.StatusCode {

	case 200: //Successful
		authResponse := AuthResponse{}
		authResponseBody, responseBodyErr := ioutil.ReadAll(response.Body)
		if responseBodyErr != nil {
			return responseBodyErr
		}

		marshallErr := json.Unmarshal(authResponseBody, &authResponse)
		if marshallErr != nil {
			return marshallErr
		}

		replicaClient.AccessToken = authResponse.AccessToken

		return nil

	case 401: // Unauthorized
		var unauthorizedResponse UnauthorizedResponse
		responseBody, responseBodyErr := ioutil.ReadAll(response.Body)
		if responseBodyErr != nil {
			return responseBodyErr
		}

		marshallErr := json.Unmarshal(responseBody, &unauthorizedResponse)
		if marshallErr != nil {
			return marshallErr
		}

		err := errors.New(unauthorizedResponse.Exception + " : " + strings.Join(unauthorizedResponse.Reasons, "; "))

		return err

	default:
		err := errors.New("unknown response")
		return err
	}
}

// GetVoices : Endpoint listing all available voices (speakers) for the calling client.
func (replicaClient *Client) GetVoices() (map[string]string, error) {

	u, _ := url.ParseRequestURI(replicaClient.ServiceEndpoint)
	u.Path = "/voice/"
	urlStr := u.String()

	request, err := http.NewRequest(http.MethodGet, urlStr, nil)
	if err != nil {
		return nil, err
	}

	if replicaClient.AccessToken != "" {

		request.Header.Set("Content-Type", "application/json")
		request.Header.Set("Connection", "close")
		request.Header.Set("Authorization", "Bearer "+replicaClient.AccessToken)
		request.Close = true

		response, requestErr := replicaClient.HTTPClient.Do(request)
		if requestErr != nil {
			return nil, requestErr
		}

		switch response.StatusCode {

		case 200: // Successful
			voiceMap := make(map[string]string)

			var voiceResponse []VoiceResponse
			voiceResponseBody, responseBodyErr := ioutil.ReadAll(response.Body)
			if responseBodyErr != nil {
				return nil, responseBodyErr
			}

			marshallErr := json.Unmarshal(voiceResponseBody, &voiceResponse)
			if marshallErr != nil {
				return nil, marshallErr
			}

			for i := 0; i < len(voiceResponse)-1; i++ {
				voiceMap[voiceResponse[i].UUID] = voiceResponse[i].Name
			}

			return voiceMap, nil

		case 401: // Unauthorized

			var unauthorizedResponse UnauthorizedResponse
			responseBody, responseBodyErr := ioutil.ReadAll(response.Body)
			if responseBodyErr != nil {
				return nil, responseBodyErr
			}

			marshallErr := json.Unmarshal(responseBody, &unauthorizedResponse)
			if marshallErr != nil {
				return nil, marshallErr
			}

			err := errors.New(unauthorizedResponse.Exception + " : " + strings.Join(unauthorizedResponse.Reasons, "; "))

			return nil, err

		default:
			err := errors.New("unknown response")

			return nil, err
		}
	}

	tokenErr := errors.New("authorization token is missing, make sure you get permission first")

	return nil, tokenErr
}

// GetSpeech : Endpoint listing all available voices (speakers) for the calling client.
func (replicaClient *Client) GetSpeech(text, speakerID string, bitRate, sampleRate int, ext SpeechExtension) (map[string]string, error) {

	data := url.Values{}
	data.Add(TEXT, text)
	data.Add(SPEAKERID, speakerID)
	data.Add(EXTENSION, string(ext))

	if bitRate > 0 {
		data.Add(BITRATE, strconv.Itoa(bitRate))
	}

	if sampleRate > 0 {
		data.Add(SAMPLERATE, strconv.Itoa(sampleRate))
	}

	u, _ := url.ParseRequestURI(replicaClient.ServiceEndpoint)
	u.Path = "/speech/"
	u.RawQuery = data.Encode()
	urlStr := u.String()

	request, err := http.NewRequest(http.MethodGet, urlStr, strings.NewReader(data.Encode()))
	if err != nil {
		return nil, err
	}

	if replicaClient.AccessToken != "" {

		request.Header.Set("Content-Type", "application/json")
		request.Header.Set("Connection", "close")
		request.Header.Set("Authorization", "Bearer "+replicaClient.AccessToken)
		request.Close = true

		response, requestErr := replicaClient.HTTPClient.Do(request)
		if requestErr != nil {
			return nil, requestErr
		}

		switch response.StatusCode {

		case 200: // Success

			var speechResponse SpeechResponse
			speechResponseBody, responseBodyErr := ioutil.ReadAll(response.Body)
			if responseBodyErr != nil {
				return nil, responseBodyErr
			}

			marshallErr := json.Unmarshal(speechResponseBody, &speechResponse)
			if marshallErr != nil {
				return nil, marshallErr
			}

			return speechResponse.URLs, nil

		case 400: // Bad Request

			var badRequestResponse BadRequestResponse
			responseBody, responseBodyErr := ioutil.ReadAll(response.Body)
			if responseBodyErr != nil {
				return nil, responseBodyErr
			}

			marshallErr := json.Unmarshal(responseBody, &badRequestResponse)
			if marshallErr != nil {
				return nil, marshallErr
			}

			err := errors.New(strconv.Itoa(badRequestResponse.ErrorCode) + " : " + badRequestResponse.Error)

			return nil, err

		case 401: // Unauthorized

			var unauthorizedResponse UnauthorizedResponse
			responseBody, responseBodyErr := ioutil.ReadAll(response.Body)
			if responseBodyErr != nil {
				return nil, responseBodyErr
			}

			marshallErr := json.Unmarshal(responseBody, &unauthorizedResponse)
			if marshallErr != nil {
				return nil, marshallErr
			}

			err := errors.New(unauthorizedResponse.Exception + " : " + strings.Join(unauthorizedResponse.Reasons, "; "))

			return nil, err

		default:
			err := errors.New("unknown response")

			return nil, err
		}
	}

	tokenErr := errors.New("authorization token is missing, make sure you get permission first")
	return nil, tokenErr
}
