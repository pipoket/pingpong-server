package main

import "net/http"


type weatherData struct {
    name string `json:"name"`
    main struct {
        kelvin float64 `json:"temp"`
    } `json:"main"`
}


func query(city string) (weatherData, error) {
    resp, err := http.Get("http://api.openweathermap.org/data/2.5/weather?q=" + city)
    if err != nil {
        return weatherData{}, err
    }

    defer resp.Body.Close()

    var d weatherData

    if err := json.NewDecoder(resp.Body).Decode(&d); err != nil {
        return weatherData{}, err
    }

    return d, nil
}


func main() {
    http.HandleFunc("/", hello)
    http.ListenAndServe(":8080", nil)
}


func hello(w http.ResponseWriter, r *http.Request) {
    w.Write([]byte("hello"))
}
