# phillycodefest
Until Jake gets the GeoLocation, we can use this temp sol. You manually enter the GIS coord in a form.

Jake is doing task 1. below. 
Amarildo is doing 2.
Paddy is doing 3. 


1. user nav to our site, we collect their loc, and display venues around them they can begin communicating w/
2. once a user clicks on a building, we create a chat room for that building.
3. each user can be in only one chat room (or building) at any one instance, so their is a mechanism for them to go to step 1 and repeat through these general loop for interaction

We will be using golang on appengine. Here is the break down of the essential api w/ respect to each step :
1. Google Places API : https://developers.google.com/places/
 - Google Places API allows you to query for detailed place information on a variety of categories, such as: establishments, prominent points of interest, geographic locations, and more. The Google Places API is also integrated into the Google Maps API as a JavaScript Library.
 - The Google Place Autocomplete feature returns place information based on text search terms and can be used to provide autocomplete functionality for text-based geographic searches.

2. Channel Go API : https://cloud.google.com/appengine/docs/go/channel/
 - The Channel API creates a persistent connection between your application and Google servers, allowing your application to send messages to JavaScript clients in real time without the use of polling. This is useful for applications designed to update users about new information immediately. Some example use-cases include collaborative applications, multi-player games, or chat rooms.

3. The overall routing is through the appengine 'net/http' and 'html/template'
 - for presentation, maybe we can find someone w/ expr using say bootstrap...

