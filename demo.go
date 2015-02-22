package Demo

import (
	"appengine"
	"appengine/user"
	"fmt"
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
	var room *Building
	var err error
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
	curr := Vertex{39.95352,-75.18845}
	var m = map[string]Vertex{
		"James Creese Center":     {39.95364,-75.18866},
		"Behrakis Hall":           {39.95352,-75.18845},
	}
	for i, v := range m {
		if lat := v.lat - curr.lat; lat <=.0002 && lat >= -.0002 {
			if lng := v.lng - curr.lng; lng <=.0002 && lng >= -.0002 {
				room, err = getBuilding(c, i)
				if err != nil {
					http.Error(w, err.Error(), 500)
					return
				}
			}
		}
	}
	
	/*
	//buildingName = getBuildingName(lat, lon)
	// lat => 39.953534, lon => -75.188456
	// Get or create the Building.
	room, err := getBuilding(c, "PHILLYCITYHALL")
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}*/
	
	fmt.Println("clientid ", u.ID)
	// Create a new Client, getting the channel token.
	token, err := room.AddClient(c, u.ID)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	// Render the HTML template.
	data := struct{ Building, Token string }{room.Name, token}
	err = rootTmpl.Execute(w, data)
	if err != nil {
		http.Error(w, err.Error(), 500)
	}
}

// post broadcasts a message to a specified Room.
func post(w http.ResponseWriter, r *http.Request) {
	c := appengine.NewContext(r)
    u := user.Current(c) // assumes 'login: required' set in app.yaml
	c.Infof("clientid ", u.ID)
	message := r.FormValue("msg")
	
	if user.IsAdmin(c) {
		message = "ADMIN:" + message
	}
	// Get the room.
	room, err := getBuilding(c, r.FormValue("room"))
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	
	// Send the message to the clients in the room.
	err = room.Send(c, message)
	if err != nil {
		http.Error(w, err.Error(), 500)
	}
}
