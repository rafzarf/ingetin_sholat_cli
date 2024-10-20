package main

import (
	"flag"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMain(m *testing.M) {
	// This will properly handle the test flags
	flag.Parse()
	os.Exit(m.Run())
}

func TestFetchPrayerTimes(t *testing.T) {
	// Create a mock server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Provide a sample response
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{
			"code": 200,
			"status": "OK",
			"data": {
				"timings": {
					"Fajr": "04:30",
					"Sunrise": "05:45",
					"Dhuhr": "12:00",
					"Asr": "15:30",
					"Sunset": "18:15",
					"Maghrib": "18:30",
					"Isha": "19:45",
					"Imsak": "04:20",
					"Midnight": "00:00"
				},
				"date": {
					"readable": "01 Jan 2023"
				}
			}
		}`))
	}))
	defer server.Close()

	// Replace the actual API URL with the mock server URL
	originalAPIBaseURL := apiBaseURL
	apiBaseURL = server.URL
	defer func() { apiBaseURL = originalAPIBaseURL }()

	// Set test values for city and country
	city = "TestCity"
	country = "TestCountry"

	// Call the function
	prayerTimes, err := FetchPrayerTimes()

	// Check for errors
	require.NoError(t, err, "FetchPrayerTimes should not return an error")

	// Check the returned data
	assert.Equal(t, 200, prayerTimes.Code, "Expected code 200")
	assert.Equal(t, "OK", prayerTimes.Status, "Expected status OK")
	assert.Equal(t, "04:30", prayerTimes.Data.Timings.Fajr, "Expected Fajr time 04:30")
	assert.Equal(t, "12:00", prayerTimes.Data.Timings.Dhuhr, "Expected Dhuhr time 12:00")
	assert.Equal(t, "15:30", prayerTimes.Data.Timings.Asr, "Expected Asr time 15:30")
	assert.Equal(t, "18:30", prayerTimes.Data.Timings.Maghrib, "Expected Maghrib time 18:30")
	assert.Equal(t, "19:45", prayerTimes.Data.Timings.Isha, "Expected Isha time 19:45")

	// Add more assertions for other prayer times and data as needed
}
