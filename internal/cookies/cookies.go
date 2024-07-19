package cookies

import "github.com/gorilla/sessions"

var Store = sessions.NewCookieStore([]byte("db-hackathon"))
