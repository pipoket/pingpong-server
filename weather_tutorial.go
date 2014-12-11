package main

import (
    "log"
    "encoding/json"
    "net/http"
)


type weatherProvider interface {
    temperature(city string) (float64, error)
}


type weatherData struct {
    Name string `json:"name"`
    Main struct {
        Kelvin float64 `json:"temp"`
    } `json:"main"`
}


type openWeatherMap struct {}


func (w openWeatherMap) temperature(city string) (float64, error) {
    resp, err := http.Get("http://api.openweathermap.org/data/2.5/weather?q=" + city)
    if err != nil {
        return 0, err
    }

    defer resp.Body.Close()

    var d struct {
        Main struct {
            Kelvin float64 `json:"temp"`
        } `json:"main"`
    }

    if err := json.NewDecoder(resp.Body).Decode(&d); err != nil {
        return 0, err
    }

    log.Printf("openWeatherMap: %s: %.2f", city, d.Main.Kelvin)
    return d.Main.Kelvin, nil
}


type weatherUnderground struct {
    apiKey string
}


func (w weatherUnderground) temperature(city string) (float64, error) {
    resp, err := http.Get("http://api.openweathermap.org/data/2.5/weather?q=" + city)
    if err != nil {
        return 0, err
    }

    defer resp.Body.Close()

    var d struct {
        Main struct {
            Kelvin float64 `json:"temp"`
        } `json:"main"`
    }

    if err := json.NewDecoder(resp.Body).Decode(&d); err != nil {
        return 0, err
    }

    d.Main.Kelvin += 2;
    log.Printf("weatherUnderground: %s: %.2f", city, d.Main.Kelvin)
    return d.Main.Kelvin, nil
}


func temperature(city string, providers ...weatherProvider) (float64, error) {
    sum := 0.0

    for _, provider := range providers {
        k, err := provider.temperature(city)
        if err != nil {
            return 0, err
        }

        sum += k
    }

    return sum / float64(len(providers)), nil
}


type multiWeatherProvider []weatherProvider


func (w multiWeatherProvider) temperature(city string) (float64, error) {
    temps := make(chan float64, len(w))
    errs := make(chan error, len(w))

    for _, provider := range w {
        go func(p weatherProvider) {
            k, err := p.temperature(city)
            if err != nil {
                errs <- err
                return
            }

            temps <- k
        }(provider)
    }

    sum := 0.0

    for i := 0; i < len(w); i++ {
        select {
        case temp := <-temps:
            sum += temp
        case err:= <-errs:
            return 0, err
        }
    }

    return sum / float64(len(w)), nil
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


func hello(w http.ResponseWriter, r *http.Request) {
    w.Write([]byte("hello"))
}
