package main

import (
	"fmt"
	"net/http"
	"io/ioutil"
	"strings"
	"regexp"
	"time"
	"strconv"
)

func main(){
	
	key2 := "295e2e49d729ecd68febc5856850ffcc"
	lon, lat := getlonlat() //get lontitude latitude for onecall api request

	time.Sleep(200 * time.Millisecond)
	// Make href2 to Get all forecast info
	url2 := "http://api.openweathermap.org/data/2.5/onecall?lat={val2}&lon={val1}&exclude=current,minutely,hourly,alerts&appid={API key}"
	url2 = strings.Replace(url2, "{val1}", lon, 1)
	url2 = strings.Replace(url2, "{val2}", lat, 1)
	url2 = strings.Replace(url2, "{API key}", key2, 1)
	// make Get request
	resp, err := http.Get(url2)
	if err != nil {
		fmt.Print(err)
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	// parse Body to find result
	deltafeel, datafeelunix := getdeltatemp(string(body)) // get min delta (feel temp, real temp)
	deltasunrise, datasunriseunix := getsunrise(string(body))


	datafeel := time.Unix(datafeelunix, 0)
	datasunrise := time.Unix(datasunriseunix, 0)
	hours := deltasunrise / 3600
	minutes := deltasunrise % 60
	year1, month1, day1 := datafeel.Date()
	year2, month2, day2 := datasunrise.Date()
	fmt.Printf("minimum delta of feel temp and real temp is %.2f on data %v %s %v\n", deltafeel, day1, month1, year1)
	fmt.Printf("Max duration of the day = %v hours and %v minutes  on data %v %s %v\n", hours, minutes, day2, month2, year2)
	exit := ""
	fmt.Println("type something to exit or CTRL-C")
	fmt.Scan(&exit)
	
}	
// function to Get lontitude and latitude
func getlonlat() (string, string) {
	
	key := "b593e03bc1c42758ff1f0758ea00311f"
	url := "http://api.openweathermap.org/data/2.5/weather?q={city name}&appid={API key}"
	city := ""
	fmt.Println("CHOSE YOUR CITY TO FIND OUT WEATHER")
	
	regexlonlat := regexp.MustCompile(`"coord":{"lon":..\.....,"lat":..\.....}`)
	regexlon := regexp.MustCompile(`"lon":..\.....`)
	regexlat := regexp.MustCompile(`"lat":..\.....`)
	var flag bool = false
	for flag != true {
		
		// make href to get API
		url = "http://api.openweathermap.org/data/2.5/weather?q={city name}&appid={API key}"
		url = strings.Replace(url, "{API key}", key, 1)
		fmt.Scan(&city)
		//fmt.Println(city)
		url = strings.Replace(url, "{city name}", city, 1)
		// fmt.Println(string(url))
		// make GET request
		resp, err := http.Get(url)
		if err != nil {
			fmt.Print(err)
		}
		defer resp.Body.Close()
		body, err := ioutil.ReadAll(resp.Body)
		// fmt.Println(string(body))
		if err != nil {
			fmt.Print(err)
		}
		parseerror := regexp.MustCompile(`{"cod":"404","message":"city not found"}`)
		error := parseerror.Find(body)
		// fmt.Println(string(error))
		if error != nil {
			fmt.Println("Wrong city, Try again")
			continue
		} else {
			flag = true
			break
			
		}
	}
	// parse Body to fing lontitude and latitude
	url = strings.Replace(url, "{city name}", city, 1)
	resp, err := http.Get(url)
		if err != nil {
			fmt.Print(err)
		}
		defer resp.Body.Close()
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			fmt.Print(err)
		}
	lonlat := regexlonlat.Find(body)	
	lon := regexlon.Find(lonlat)
	lonval := strings.Trim(string(lon), `"lon:"`)
	lat:= regexlat.Find(lonlat)
	latval := strings.Trim(string(lat), `"lat:"`)

	return lonval, latval
}


// function to Get minimum delta and data of the day, when current temp and feel temp are minimum 
func getdeltatemp(body string) (float64, int64) {
	var mindelta float64 = 100000
	var minday []byte
	array := strings.Split(string(body), `"dt":`)
	deltaregex := regexp.MustCompile(`"night":\d\d\d...`)
	dataregex := regexp.MustCompile(`^\d{10}`)
	for idx, str := range array {
		if idx == 0 {continue}

		deltaarray := deltaregex.FindAll([]byte(str), -1)
		currentday := dataregex.Find([]byte(str))
		val0 := strings.Replace(string(deltaarray[0]), `"night":`, ``, 1)
		val1 := strings.Replace(string(deltaarray[1]), `"night":`, ``, 1)
		val0 = strings.Trim(val0, `,"e`)
		val1 = strings.Trim(val1, `,"e`)
		
		floatval0, _ := strconv.ParseFloat(val0, 64)
		floatval1, _ := strconv.ParseFloat(val1, 64)
		delta := floatval0 - floatval1
		if delta < 0 {delta = -delta}
		if delta < mindelta {
			mindelta = delta
			minday = currentday
		}
	}
	intdata, _ := strconv.ParseInt(string(minday), 10, 64)
	return mindelta, intdata
}


// function to Get delta between sunrise and sunset
func getsunrise(body string) (int64, int64) {
	var maxdelta int64 = 0
	var maxday []byte
	array := strings.Split(string(body), `"dt":`)
	sunriseregex := regexp.MustCompile(`"sunrise":\d{10}`)
	sunsetregex := regexp.MustCompile(`"sunset":\d{10}`)
	dataregex := regexp.MustCompile(`^\d{10}`)
	for idx, str := range array {
		if idx == 0 {continue}
		if idx >= 5 {break}
		currentday := dataregex.Find([]byte(str))
		sunrise := sunriseregex.Find([]byte(str))
		sunset := sunsetregex.Find([]byte(str))
		aboba0 := strings.Replace(string(sunrise), `"sunrise":`, ``, 1)
		aboba1 := strings.Replace(string(sunset), `"sunset":`, ``, 1)

		val1, _ := strconv.ParseInt(aboba0, 10, 64)
		val2, _ := strconv.ParseInt(aboba1, 10, 64)
		delta := val2 - val1
		if delta > maxdelta {
			maxdelta = delta
			maxday = currentday
		}

	}
	intdata, _ := strconv.ParseInt(string(maxday), 10, 64)
	return maxdelta, intdata
}