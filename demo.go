package Demo

import (
	"appengine"
	"appengine/user"
	
    "net/http"
    "html/template"
)

const clientIdLen = 40

func init() {
	// Register our handlers with the http package.
	http.HandleFunc("/", root)
	http.HandleFunc("/post", post)
}

// rootTmpl is the main (and only) HTML template.
var rootTmpl = template.Must(template.ParseFiles("tmpl/root.html"))



// root is an HTTP handler that joins or creates a Room,
// creates a new Client, and writes the HTML response.
func root(w http.ResponseWriter, r *http.Request) {

	c := appengine.NewContext(r)
    u := user.Current(c) // assumes 'login: required' set in app.yaml
    if u == nil {
        url, err := user.LoginURL(c, r.URL.String())
        if err != nil {
            http.Error(w, err.Error(), http.StatusInternalServerError)
            return
        }
        w.Header().Set("Location", url)
        w.WriteHeader(http.StatusFound)
        return
    }
	
	// Get or create the Room.
	room, err := getRoom(c, "PHILLYCITYHALL")
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	// Create a new Client, getting the channel token.
	token, err := room.AddClient(c, u.ID)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	// Render the HTML template.
	data := struct{ Room, Token string }{room.Name, token}
	err = rootTmpl.Execute(w, data)
	if err != nil {
		http.Error(w, err.Error(), 500)
	}
}

// post broadcasts a message to a specified Room.
func post(w http.ResponseWriter, r *http.Request) {
	c := appengine.NewContext(r)
	// Get the room.
	room, err := getRoom(c, r.FormValue("room"))
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	// Send the message to the clients in the room.
	err = room.Send(c, r.FormValue("msg"))
	if err != nil {
		http.Error(w, err.Error(), 500)
	}
}
