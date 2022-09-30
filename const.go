package main

// Routes
const (
	PING_ENDPOINT          = "/ping"
	AUTHORIZATION_ENDPOINT = "/auth/authorize"
	TOKEN_ENDPOINT         = "/auth/token"
	LOGIN_ENDPOINT         = "/login"
	CONSENT_ENDPOINT       = "/consent"
)

// Files
const (
	LOGIN_TEMPLATE   = "login.gtpl"
	CONSENT_TEMPLATE = "consent.gtpl"
)

// Session
const (
	AUTHORIZATION_SESSION  = "authorization_session"
	AUTHENTICATION_SESSION = "authentication_session"
	TICKET                 = "ticket"
	USER_ID                = "user_id"
	CLIENT_ID              = "client_id"
	SCOPES                 = "scopes"
)

const (
	AUTHLETE_API = "AuthleteApi"
)
