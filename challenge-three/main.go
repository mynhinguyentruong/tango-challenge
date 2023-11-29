package main

import (
	"fmt"
	"strings"
  "strconv"
  "log"
  "net/http"
  "bytes"
	"encoding/json"
  "io"
)

type TangoGuessResponse struct {
  Status string `json:"status"`
	CatchAll map[string]interface{} `json:"-"`
}

func submitCorrectAnswer(arr []int) {
  url := "https://interview.tangohq.com/guess-escpos"
	contentType := "application/json"
  bearerToken := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzZXNzaW9uSWQiOiIyMGQ5ZjJlYi0yOGVjLTQxYWYtODQyMy1lY2ZjZjQ3MmRiNWQiLCJtZXNzYWdlIjoiRGVjb2RpbmcgdGhlIHRva2VuLCBuaWNlLCBoYXZlIGEgcHJpemUiLCJ1cmwiOiIvYm9udXMtam9uYXMiLCJpYXQiOjE3MDEwMzc3MzR9.NCcGnr_VL2IrQ_dzqzcIOHgc-yV-6WOIvsv_XvYY2w4" 

  client := &http.Client{}

  payload, err := json.Marshal(map[string][]int{"myGuess": arr})
  if err != nil {
    log.Fatal("Error encoding JSON:", err)
  }

  req, err := http.NewRequest("POST", url, bytes.NewBuffer(payload))
  if err != nil {
    log.Fatal("Error creating HTTP request:", err)
  }

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
    if err := json.Unmarshal(body, &response.CatchAll); err != nil {
      log.Fatal("Error unmarshalling JSON response:", err)
    }      

    log.Println("Result:", response)

  } else {
    log.Fatalf("Unexpected response status code: %d", res.StatusCode)
  }
  
}

func main() {
  var result []int

 	// Define a map to store ASCII control characters and their values
	controlCharMap := map[string]int{
		"NUL": 0, "SOH": 1, "STX": 2, "ETX": 3,
		"EOT": 4, "ENQ": 5, "ACK": 6, "BEL": 7,
		"BS": 8, "HT": 9, "LF": 10, "VT": 11,
		"FF": 12, "CR": 13, "SO": 14, "SI": 15,
		"DLE": 16, "DC1": 17, "DC2": 18, "DC3": 19,
		"DC4": 20, "NAK": 21, "SYN": 22, "ETB": 23,
		"CAN": 24, "EM": 25, "SUB": 26, "ESC": 27,
		"FS": 28, "GS": 29, "RS": 30, "US": 31,
	}

	// Result of manually lookup ESC/POS® Command
	input := "Hello Tango ESC J 1 I'm excited to join ESC J 2 GS V 65 0"

	// Split input into words
	words := strings.Fields(input)

	// Convert each word to its corresponding ASCII control character value
	for _, word := range words {
    // Lookup if word is in the map
		value, found := controlCharMap[word]
		if found {
			fmt.Printf("Found special command %s: %d\n", word, value)
      result = append(result, value)
		} else {
      fmt.Printf("Looking at char: %v\n", word)

      // Convert to int if possible "10" -> 10
      s, err := strconv.Atoi(word) 
      if err == nil {
        result = append(result, s)
	    } else {
        // Convert string to ASCII codes
        for _, char := range word {
          fmt.Printf("Converting %v to integer: %v\n", char, int(char))
          result = append(result, int(char))
        }
      }
		}
	}

  	// Print the resulting ESCPOS sequence
	fmt.Println("ESCPOS Sequence:")
	fmt.Println(result)

  submitCorrectAnswer(result)

}

