package main

import (
  "encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
)

type Info struct {
  Response Response `json:"Response"`
}

type Response struct {
  View []View `json:"View"`
}

type View struct {
  Result []Result `json:"Result"`
}

type Result struct {
  Location Location `json:"Location"`
}

type Location struct {
  Address Address `json:"Address"`
}

type Address struct {
  AdditionalData []AdditionData `json: "AdditionData"`
}

type AdditionData struct {
  Value string `json: "value"`
  Key string `json: "key"`
}

func handler(w http.ResponseWriter, r *http.Request) {
  lat, long := getCoords(r)
  info := getJSON(lat, long)

  if(info != nil){
    displayOutput(info, w)
  }
}

func getCoords(req *http.Request) (string, string) {
  err := req.ParseForm()
  if err != nil {
    panic(err)
  }
  v := req.Form
  return v.Get("latitude"), v.Get("longitude")
}

func getJSON(lat string, long string) *Info {
  url := "https://reverse.geocoder.api.here.com/6.2/reversegeocode.json?prox="+lat+","+long+",250&mode=retrieveAddresses&maxresults=1&gen=9&app_id=th4BCanYJ2cBxYFMPiX2&app_code=y0iHOrPJZnYoml89Z0gz_Q"
  response, err := http.Get(url)
  if err != nil {
  	fmt.Print(err.Error())
  	os.Exit(1)
  }

  body, err := ioutil.ReadAll(response.Body)
  if err != nil {
  	log.Fatal(err)
    return nil
  }

  var info *Info
  jsonErr := json.Unmarshal(body, &info)
  if jsonErr != nil {
    log.Fatal(jsonErr)
    return nil
  } else {
    return info
  }
}

func displayOutput(info *Info, w http.ResponseWriter)  {
  if(len(info.Response.View) == 0){
    fmt.Fprintf(w, "Invalid Coordinates")
  } else {
    country := info.Response.View[0].Result[0].Location.Address.AdditionalData

    if(country[0].Value != "United States"){
      fmt.Fprintf(w, "Coordinates are not in US")
    } else {
      fmt.Fprintf(w, country[1].Value)
    }
  }
}

func main() {
    http.HandleFunc("/", handler)
    log.Fatal(http.ListenAndServe(":8080", nil))
}
