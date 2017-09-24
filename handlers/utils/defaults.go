package utils

import "fmt"
import "net/http"
import "io/ioutil"

const (
	templateFailedFormat      = "Template failed to compile: %s"
	sessionFailedFormat       = "Failed to retrieve session info: %s"
	tokenExchangeFailedFormat = "Failed to exchange token: %s"
	discordErrorFormat        = "Discord sent back an error: %s"
	ioErrorFormat             = "Failed to read request body: %s"
	jsonErrorFormat           = "Failed to unmarshal JSON: %s"
	apiErrorFormat            = "Failed to query radio-api: %s"
	discordURLFormat          = "https://discordapp.com/api/oauth2/authorize?client_id=%s&scope=identify&response_type=code"
	tokenURLFormat            = "https://discordapp.com/api/oauth2/token?grant_type=authorization_code&client_id=%s&client_secret=%s&code=%s"
	listRadiosFormat          = "%s/radios?limit=%s&offset=%s&state=%s"
	createRadiosFormat        = "%s/radios"
	getRadioFormat            = "%s/radios/%s"
	// SessionName refers to the name of the session
	SessionName = "discord"
	// DiscordMeURL refers to the API endpoint for Discord user data
	DiscordMeURL = "https://discordapp.com/api/v6/users/@me"
	// NotAuthorized is returned when a user tries to access a page they don't have the permissions for
	NotAuthorized = "You are not permitted to view this page."
)

// TemplateFailed is returned when a template fails to parse
func TemplateFailed(e error) string {
	return fmt.Sprintf(templateFailedFormat, e)
}

// SessionFailed is returned when a session fails to retrieve
func SessionFailed(e error) string {
	return fmt.Sprintf(sessionFailedFormat, e)
}

// TokenExchangeFailed is returned when a token exchange with Discord fails
func TokenExchangeFailed(e error) string {
	return fmt.Sprintf(tokenExchangeFailedFormat, e)
}

// DiscordError is returned when the Discord OAuth fails
func DiscordError(e string) string {
	return fmt.Sprintf(discordErrorFormat, e)
}

// IOError is returned when an IO error occurs
func IOError(e error) string {
	return fmt.Sprintf(ioErrorFormat, e)
}

// JSONError is returned when an error occurs demarshaling JSON
func JSONError(e error) string {
	return fmt.Sprintf(jsonErrorFormat, e)
}

// APIError is returned when the dab-radio API returns an error
func APIError(e error, resp *http.Response) string {
	if e != nil {
		return fmt.Sprintf(apiErrorFormat, e)
	}
	if resp == nil {
		return "unknown error"
	}
	body, _ := ioutil.ReadAll(resp.Body)
	if string(body) != "" {
		return fmt.Sprintf("%s: %s", resp.Status, string(body))
	}
	return resp.Status

}

// DiscordAuthURL returns the authentication URL
func DiscordAuthURL(clientID string) string {
	return fmt.Sprintf(discordURLFormat, clientID)
}

// DiscordTokenURL returns the token url
func DiscordTokenURL(clientID string, clientSecret string, code string) string {
	return fmt.Sprintf(tokenURLFormat, clientID, clientSecret, code)
}

// ListRadiosURL returns the URL to query for radios with a given state
func ListRadiosURL(apiBase string, limit string, offset string, state string) string {
	return fmt.Sprintf(listRadiosFormat, apiBase, limit, offset, state)
}

// CreateRadiosURL returns the URL to create a station
func CreateRadiosURL(apiBase string) string {
	return fmt.Sprintf(createRadiosFormat, apiBase)
}

// GetRadioURL returns the URL to query a radio
func GetRadioURL(apiBase string, name string) string {
	return fmt.Sprintf(getRadioFormat, apiBase, name)
}
