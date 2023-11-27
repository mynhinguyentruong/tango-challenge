package main

import (
  "log"
  "math/rand"
  "net/http"
  "bytes"
	"encoding/json"
  "io"

)

type TangoGuessResponse struct {
  Status string `json:"status"`
	CatchAll map[string]interface{} `json:"-"`
}

type Value struct {
  Number int
  Status string
}

func main() {
  var min, max, randomNumber, numberOfGuesses int
  max = 100000000
  min = 0
  numberOfGuesses = 0

	url := "https://interview.tangohq.com/guess"
	contentType := "application/json"
  bearerToken := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzZXNzaW9uSWQiOiIyMGQ5ZjJlYi0yOGVjLTQxYWYtODQyMy1lY2ZjZjQ3MmRiNWQiLCJtZXNzYWdlIjoiRGVjb2RpbmcgdGhlIHRva2VuLCBuaWNlLCBoYXZlIGEgcHJpemUiLCJ1cmwiOiIvYm9udXMtam9uYXMiLCJpYXQiOjE3MDEwMzc3MzR9.NCcGnr_VL2IrQ_dzqzcIOHgc-yV-6WOIvsv_XvYY2w4" 

  client := &http.Client{}

  for {
    randomNumber = rand.Intn(max - min + 1) + min
    numberOfGuesses += 1

    // create payload
    payload, err := json.Marshal(map[string]int{"myGuess": randomNumber})
    if err != nil {
      log.Fatal("Error encoding JSON:", err)
    }

    // create post http request
    req, err := http.NewRequest("POST", url, bytes.NewBuffer(payload))
    if err != nil {
      log.Fatal("Error creating HTTP request:", err)
    }

    // set request header
    req.Header.Set("Content-Type", contentType)
    req.Header.Set("Authorization", "Bearer "+bearerToken)

    res, err := client.Do(req)
    if err != nil {
      log.Fatal("Error making POST request:", err)
    }
    defer res.Body.Close()

    // Check for a successful response (status code 2xx)
    if res.StatusCode >= 200 && res.StatusCode < 300 {
      var response TangoGuessResponse

 			// Read the response body into a byte slice
			body, err := io.ReadAll(res.Body)
			if err != nil {
				log.Fatal("Error reading response body:", err)
			}

			// Unmarshal the response body into the TangoGuessResponse struct
			if err := json.Unmarshal(body, &response); err != nil {
				log.Fatal("Error unmarshalling JSON response:", err)
			}      

      if response.Status == "higher" {
        min = randomNumber + 1
      } else if response.Status == "lower" {
        max = randomNumber - 1
      } else {
        // Unmarshal the response body again when status is correct to get the next challenge
        if err := json.Unmarshal(body, &response.CatchAll); err != nil {
          log.Fatal("Error unmarshalling JSON response:", err)
        }

        log.Println("Response: ", response)
        log.Printf("Found the right number: %v \nNumber of Guesses: %v", randomNumber, numberOfGuesses)

        break
      }
    } else {
      log.Fatalf("Unexpected response status code: %d", res.StatusCode)
    }
	}
}

