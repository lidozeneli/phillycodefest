package pcf

import (
//	"html/template"
	"net/http"
	"fmt"

//	"appengine"
//	"appengine/datastore"
)
type Geo struct {
	lat, long float64
}
type User struct {
//	FirstName      string
//	LastName       string
//	UserName       string
	InsideBuilding bool
	CurrentPlace   Geo
	
}

func geo(w http.ResponseWriter, r *http.Request) {

}

func root(w http.ResponseWriter, r *http.Request) { 
//	c := appengine.NewContext(r)
	
	fmt.Fprintf(w, "<!DOCTYPE html> <html><head><title>Go Philly-{CodeFest}</title></head> <body> <p>Click the button to get your coordinates.</p> <button onclick='getLocation()'>Try It</button> <p id='demo'></p> <script> var x = document.getElementById('demo'); function getLocation() { if (navigator.geolocation) { navigator.geolocation.getCurrentPosition(showPosition); } else {  x.innerHTML = 'Geolocation is not supported by this browser.'; } } function showPosition(position) { var latlon = position.coords.latitude + ',' + position.coords.longitude;   x.innerHTML = 'Latitude: ' + position.coords.latitude +  '<br>Longitude: ' + position.coords.longitude;	 } </script> <form action='/geo' method='post'> lat:<br> <input type='text' name='lat'> long:<br> <input type='text' name='long'> </form></body></html>")
}

func init() {
	http.HandleFunc("/", root)
	http.HandleFunc("/geo",geo)
}
