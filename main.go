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

type Response map[string]string

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
  bearerToken := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzZXNzaW9uSWQiOiI1ZmEwMjQxZi0yZTcxLTQ1NjMtYmU5Zi05Mzk3MzExODU5YWMiLCJtZXNzYWdlIjoiRGVjb2RpbmcgdGhlIHRva2VuLCBuaWNlLCBoYXZlIGEgcHJpemUiLCJ1cmwiOiIvYm9udXMtam9uYXMiLCJpYXQiOjE3MDEwODk5ODd9.EEO78UMsxlayBKnXhvM4wFU74x0ld0yDKOab0iK2l7s" 

  client := &http.Client{}

  for {
    randomNumber = rand.Intn(max - min + 1) + min
    numberOfGuesses += 1

    payload, err := json.Marshal(map[string]int{"myGuess": randomNumber})
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
			if err := json.Unmarshal(body, &response); err != nil {
				log.Fatal("Error unmarshalling JSON response:", err)
			}      

      if response.Status == "higher" {
        min = randomNumber + 1
      } else if response.Status == "lower" {
        max = randomNumber - 1
      } else {
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


// Can try doing this way for goroutine, spawn thread in main()
    // go getTangoResponse(client, url, ch, randomNumberOne, contentType, bearerToken)
    // go getTangoResponse(client, url, ch, randomNumberTwo, contentType, bearerToken)
    //
    // firstValue := <- ch 
    // secondValue := <-ch
    // if firstValue.Status  == "correct" {
    //   log.Println("Response: ", firstValue)
    //   log.Printf("Found the right number: %v \nNumber of Guesses: %v", firstValue.Number, numberOfGuesses)
    //
    //   break
    // }
    //
    // if secondValue.Status == "correct" {
    //   log.Println("Response: ", secondValue)
    //   log.Printf("Found the right number: %v \nNumber of Guesses: %v", secondValue.Number, numberOfGuesses)
    //
    //   break
    // }
    //
    // if firstValue.Status == "higher" &&  secondValue.Status == "higher" {
    //   log.Println("both were higher")
    //   if firstValue.Number > secondValue.Number {
    //     min = firstValue.Number + 1
    //   } else {
    //     min = secondValue.Number + 1
    //   }
    //   // update min to the highest number
    // } else if firstValue.Status == "lower" && secondValue.Status == "lower" {
    //   log.Println("bother were lower")
    //   // update max to the lowest number
    //   if firstValue.Number < secondValue.Number {
    //     max = firstValue.Number - 1
    //   } else {
    //     max = secondValue.Number - 1
    //   }
    //
    // } else {
    //   log.Println("one low, one high")
    //   // which is the case that one is lower one is higher
    //   // update min and max
    //
    //   if firstValue.Number < secondValue.Number {
    //
    //     min = firstValue.Number
    //     max = secondValue.Number
    //   } else {
    //     min = secondValue.Number
    //     max = firstValue.Number
    //   }
    // }

// func getTangoResponse(client *http.Client, url string, ch chan<- Value, number int, contentType string, bearerToken string) {
//     log.Println("Guessing:", number)
//
//     payload, err := json.Marshal(map[string]int{"myGuess": number})
//     if err != nil {
//       log.Fatal("Error encoding JSON:", err)
//     }
//
//     req, err := http.NewRequest("POST", url, bytes.NewBuffer(payload))
//     if err != nil {
//       log.Fatal("Error creating HTTP request:", err)
//     }
//
//     req.Header.Set("Content-Type", contentType)
//     req.Header.Set("Authorization", "Bearer "+bearerToken)
//
//     res, err := client.Do(req)
//     if err != nil {
//       log.Fatal("Error making POST request:", err)
//     }
//
//     if res.StatusCode >= 200 && res.StatusCode < 300 {
//       var response TangoGuessResponse
//
//       // Decode the response body into the Response struct
//       if err := json.NewDecoder(res.Body).Decode(&response); err != nil {
//         log.Fatal("Error decoding JSON response:", err)
//       }
//
//       // Access the fields in the response object
//       log.Println("Result:", response.Status)
//       val := Value{number, response.Status}
//
//       ch <- val 
//     } else {
//       // if status code is 502
//       // retry instead of crashing the program
//       log.Fatalf("Unexpected response status code: %d", res.StatusCode)
//     }
// }
