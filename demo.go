package Loomup

import (
	"appengine"
	//"appengine/user"
	"fmt"
    "net/http"
    "html/template"
    "strconv"
    "time"

)

const clientIdLen = 40

func init() {
	// Register our handlers with the http package.
	http.HandleFunc("/", root)
	http.HandleFunc("/geo", geo)
	http.HandleFunc("/post", post)
	
}

// rootTmpl is the main (and only) HTML template.
var rootTmpl = template.Must(template.ParseFiles("tmpl/root.html"))
var mainTmpl = template.Must(template.ParseFiles("tmpl/main.html"))

type resultrec struct{ Building string 
	Token string
	Count int 
	} 


func root(w http.ResponseWriter, r *http.Request) {
	err := rootTmpl.Execute(w, map[string]string{
        "token":    "",
        "me":       "",
        "location": "",
    })
	if err != nil {
		http.Error(w, err.Error(), 500)
	}

}

// root is an HTTP handler that joins or creates a Room,
// creates a new Client, and writes the HTML response.
func geo(w http.ResponseWriter, r *http.Request) {
	var room *Building
	var err error
	c := appengine.NewContext(r)
	
	/*
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
    */
    //39.9568382&longitude=-75.1812252
    
	latitude, _ :=	 strconv.ParseFloat(r.FormValue("latitude"), 64)
	longitude, _ :=  strconv.ParseFloat(r.FormValue("longitude"), 64)
	curr := Vertex{ latitude,  longitude}
	
	//read from the locatoin file and build the m map
	
//	curr := Vertex{39.95352,-75.18845}
	var m = map[string]Vertex{
		"James Creese Center":     {39.95364,-75.18866},
		"Behrakis Hall":           {39.95352,-75.18845},
	}

	//var room 
	for i, v := range m {
		if lat := v.lat - curr.lat; lat <=.2 && lat >= -.2 {
			if lng := v.lng - curr.lng; lng <=.2 && lng >= -.2 {
				room, err = getBuilding(c, i)
				if err != nil {
					http.Error(w, err.Error(), 500)
					return
				}
			}
		}
	}
	
	//buildingName = getBuildingName(lat, lon)
	// lat => 39.953534, lon => -75.188456
	// Get or create the Building.
	//room, err = getBuilding(c, "PHILLYCITYHALL")
	/*
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	*/
	if room == nil {
		//to do add
		fmt.Fprintf(w, "No building in you area." )
		return	
	}
	//fmt.Println("clientid ", u.ID)
	// Create a new Client, getting the channel token.
	userid := randId(20)
	token, err := room.AddClient(c, userid)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	count, cerr := room.GetCount(c)
	c.Infof("count", count)
	if cerr == nil{
		c.Infof("cerr", cerr)
	}
	
	// Render the HTML template.
	
	//data := struct{ Building, Token string, Count int }{room.Name , token, count}
	data := resultrec{room.Name , token, count}
	err = mainTmpl.Execute(w, data)
	if err != nil {
		http.Error(w, err.Error(), 500)
	}
}

// post broadcasts a message to a specified Room.
func post(w http.ResponseWriter, r *http.Request) {
	const layout = "3:04pm (MST)"
	posttime := time.Now().Format(layout)
	c := appengine.NewContext(r)
    //u := user.Current(c) // assumes 'login: required' set in app.yaml
	//c.Infof("clientid ", u.ID)
	message := r.FormValue("msg") + "   " + posttime
	/*
	if user.IsAdmin(c) {
		message = "ADMIN:" + message
	}
	*/
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
