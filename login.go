package webAppGo

// MaxUserAttempts is the maximum number of failed login attempts to be made
// on a specific user over a given period of time
var MaxUserAttempts = 3

// MaxIPAttempts is the maximum number of failed login attempts to be made
// from a specific ip address over a given period of time
var MaxIPAttempts = 6

// LoginAttemptTime is the period of time in which the use is allowed to make
// incorrect logins within this time frame
var LoginAttemptTime = 10 * 60

// Login is for tracking successful and failed login attempts of users and ip addresses
type Login struct {
	IP        string
	UserName  string
	Timestamp int
	Attempt   bool
}
