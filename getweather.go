package main

import (
	"os"
	"fmt"
	"net/http"
	//"io/ioutil"
	"github.com/davecgh/go-spew/spew"
	"encoding/json"
	"strings"
)


func parseGeneric(city string) {
	resp, _ := http.Get("http://samples.openweathermap.org/data/2.5/weather?q="+city+"&appid=b6907d289e10d714a6e88b30761fae22")
	defer resp.Body.Close()

	var j interface{}
	json.NewDecoder(resp.Body).Decode(&j)
	spew.Dump(j)

}


type weatherData struct {
	Name string `json:"name"`
	Main struct {
		TempKelvin float64 `json:"temp"`
    	} `json:"main"`
}

func (wd weatherData) TempCelcius() (float64, error) {
	return wd.Main.TempKelvin - 273.15, nil
}

func fetchWeatherData(city string) (weatherData, error) {
	fmt.Println("Fetching weather ... ", city)
	resp, e := http.Get(fmt.Sprintf("http://samples.openweathermap.org/data/2.5/weather?q=%v&appid=b6907d289e10d714a6e88b30761fae22", city))
	defer resp.Body.Close()
	
	var wd weatherData
	json.NewDecoder(resp.Body).Decode(&wd)

	fmt.Println("Found it : ", city)

	return wd, e
}



func weather(cities []string) {
	weatherDataList := map[string]weatherData{}

	weatherDataChannel := make(chan struct{string; weatherData})

	for _, city := range cities {
		go func(city string) {
			wd, e := fetchWeatherData(city)
			if e != nil {fmt.Println("error: ", e)}
			weatherDataChannel <- struct{string; weatherData}{city, wd}
		}(city)
	}

	for i:=0; i<len(cities); i++ {
		pair := <- weatherDataChannel
		city := pair.string
		wd := pair.weatherData
		weatherDataList[city] = wd
		spew.Dump(wd.TempCelcius())
	}

	spew.Dump(weatherDataList)

	// Dump to json
	f, _ := os.Create("plop.json")
	defer f.Close()
	as_json, _ := json.MarshalIndent(weatherDataList, "", "    ")
	f.Write(as_json)
}

func webWeather(w http.ResponseWriter, r *http.Request) {
	cities := strings.Split(r.URL.Path, "/")[1]
	weather(strings.Split(cities, ","))
}

func main() {
	println(http.HandleFunc)
	//http.HandleFunc("/", webWeather)
	//http.ListenAndServe(":8080", nil)
	weather([]string{"Toulouse", "Berlin", "London", "Paris", "Sarayevo", "plip", "plop"})
}
