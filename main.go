package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
)

type PrayerTimesResponse struct {
	Code   int    `json:"code"`
	Status string `json:"status"`
	Data   struct {
		Timings struct {
			Fajr     string `json:"Fajr"`
			Sunrise  string `json:"Sunrise"`
			Dhuhr    string `json:"Dhuhr"`
			Asr      string `json:"Asr"`
			Sunset   string `json:"Sunset"`
			Maghrib  string `json:"Maghrib"`
			Isha     string `json:"Isha"`
			Imsak    string `json:"Imsak"`
			Midnight string `json:"Midnight"`
		} `json:"timings"`
		Date struct {
			Readable string `json:"readable"`
		} `json:"date"`
	} `json:"data"`
}

func main() {
	fmt.Println("Salat Reminder CLI for Bandung, Indonesia")

	for {
		prayerTimes, err := getPrayerTimes()
		if err != nil {
			fmt.Printf("Error fetching prayer times: %v\n", err)
			time.Sleep(time.Minute)
			continue
		}

		fmt.Printf("Prayer times for %s:\n", prayerTimes.Data.Date.Readable)
		fmt.Printf("Fajr: %s\n", prayerTimes.Data.Timings.Fajr)
		fmt.Printf("Sunrise: %s\n", prayerTimes.Data.Timings.Sunrise)
		fmt.Printf("Dhuhr: %s\n", prayerTimes.Data.Timings.Dhuhr)
		fmt.Printf("Asr: %s\n", prayerTimes.Data.Timings.Asr)
		fmt.Printf("Maghrib: %s\n", prayerTimes.Data.Timings.Maghrib)
		fmt.Printf("Isha: %s\n", prayerTimes.Data.Timings.Isha)

		checkAndNotifyPrayers(prayerTimes)

		time.Sleep(time.Minute)
	}
}

func getPrayerTimes() (*PrayerTimesResponse, error) {
	currentDate := time.Now().Format("02-01-2006")
	url := fmt.Sprintf("http://api.aladhan.com/v1/timingsByCity/%s?city=Bandung&country=Indonesia&method=11", currentDate)

	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var prayerTimes PrayerTimesResponse
	err = json.Unmarshal(body, &prayerTimes)
	if err != nil {
		return nil, err
	}

	return &prayerTimes, nil
}

func checkAndNotifyPrayers(prayerTimes *PrayerTimesResponse) {
	now := time.Now()
	checkPrayer(now, prayerTimes.Data.Timings.Fajr, "Fajr")
	checkPrayer(now, prayerTimes.Data.Timings.Dhuhr, "Dhuhr")
	checkPrayer(now, prayerTimes.Data.Timings.Asr, "Asr")
	checkPrayer(now, prayerTimes.Data.Timings.Maghrib, "Maghrib")
	checkPrayer(now, prayerTimes.Data.Timings.Isha, "Isha")
}

func checkPrayer(now time.Time, prayerTime string, prayerName string) {
	t, err := time.Parse("15:04", prayerTime)
	if err != nil {
		fmt.Printf("Error parsing prayer time: %v\n", err)
		return
	}

	prayerDateTime := time.Date(now.Year(), now.Month(), now.Day(), t.Hour(), t.Minute(), 0, 0, now.Location())
	if now.After(prayerDateTime) && now.Sub(prayerDateTime) < time.Minute {
		fmt.Printf("It's time for %s prayer!\n", prayerName)
	}
}
