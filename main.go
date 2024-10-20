package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

// PrayerTimesResponse represents the structure of the JSON response from the AlAdhan API.
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

// Constants for API configuration
const (
	apiBaseURL = "http://api.aladhan.com/v1/timingsByCity"
	city       = "Bandung"
	country    = "Indonesia"
	method     = "11" // Calculation method for Bandung
)

func main() {
	fmt.Println("Salat Reminder CLI for Bandung, Indonesia")

	for {
		prayerTimes, err := fetchPrayerTimes()
		if err != nil {
			handleError("Error fetching prayer times", err)
			time.Sleep(time.Minute)
			continue
		}

		displayPrayerTimes(prayerTimes)
		checkAndNotifyPrayers(prayerTimes)

		time.Sleep(time.Minute)
	}
}

// fetchPrayerTimes retrieves prayer times from the AlAdhan API.
func fetchPrayerTimes() (*PrayerTimesResponse, error) {
	currentDate := time.Now().Format("02-01-2006")
	url := fmt.Sprintf("%s/%s?city=%s&country=%s&method=%s", apiBaseURL, currentDate, city, country, method)

	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
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

// displayPrayerTimes prints formatted prayer times to the console.
func displayPrayerTimes(prayerTimes *PrayerTimesResponse) {
	fmt.Printf("\nPrayer times for %s:\n", prayerTimes.Data.Date.Readable)
	fmt.Printf("%-10s %s\n", "Fajr:", prayerTimes.Data.Timings.Fajr)
	fmt.Printf("%-10s %s\n", "Sunrise:", prayerTimes.Data.Timings.Sunrise)
	fmt.Printf("%-10s %s\n", "Dhuhr:", prayerTimes.Data.Timings.Dhuhr)
	fmt.Printf("%-10s %s\n", "Asr:", prayerTimes.Data.Timings.Asr)
	fmt.Printf("%-10s %s\n", "Maghrib:", prayerTimes.Data.Timings.Maghrib)
	fmt.Printf("%-10s %s\n", "Isha:", prayerTimes.Data.Timings.Isha)
}

// checkAndNotifyPrayers checks if it's time for any prayer and notifies if so.
func checkAndNotifyPrayers(prayerTimes *PrayerTimesResponse) {
	now := time.Now()
	checkPrayer(now, prayerTimes.Data.Timings.Fajr, "Fajr")
	checkPrayer(now, prayerTimes.Data.Timings.Dhuhr, "Dhuhr")
	checkPrayer(now, prayerTimes.Data.Timings.Asr, "Asr")
	checkPrayer(now, prayerTimes.Data.Timings.Maghrib, "Maghrib")
	checkPrayer(now, prayerTimes.Data.Timings.Isha, "Isha")
}

// checkPrayer checks if it's time for a specific prayer and prints a notification if so.
func checkPrayer(now time.Time, prayerTime string, prayerName string) {
	prayerDateTime, err := parsePrayerTime(now, prayerTime)
	if err != nil {
		handleError("Error parsing prayer time", err)
		return
	}

	if isWithinNotificationWindow(now, prayerDateTime) {
		fmt.Printf("\nIt's time for %s prayer!\n", prayerName)
	}
}

// parsePrayerTime converts a prayer time string to a time.Time object.
func parsePrayerTime(date time.Time, timeStr string) (time.Time, error) {
	t, err := time.Parse("15:04", timeStr)
	if err != nil {
		return time.Time{}, err
	}
	return time.Date(date.Year(), date.Month(), date.Day(), t.Hour(), t.Minute(), 0, 0, date.Location()), nil
}

// isWithinNotificationWindow checks if the current time is within one minute after the prayer time.
func isWithinNotificationWindow(now, prayerTime time.Time) bool {
	return now.After(prayerTime) && now.Sub(prayerTime) < time.Minute
}

// handleError prints an error message to the console.
func handleError(message string, err error) {
	fmt.Printf("%s: %v\n", message, err)
}
