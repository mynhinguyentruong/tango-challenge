package main

import (
  "log"
  "net/http"
  "bytes"
	"encoding/json"
  "io"
  "strings"
  "errors"
)

type Hint struct {
  Jaro float64 `json:"jaro"` 
  PositionAndCharacter []bool `json:"positionAndCharacter"`
  Character []bool `json:"character"`
}

type TangoGuessResponse struct {
  Status string `json:"status"`
	CatchAll map[string]interface{} `json:"-"`
  Hint Hint `json:"hint, omitempty"` 
}

// bearerToken := eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzZXNzaW9uSWQiOiIyMGQ5ZjJlYi0yOGVjLTQxYWYtODQyMy1lY2ZjZjQ3MmRiNWQiLCJtZXNzYWdlIjoiRGVjb2RpbmcgdGhlIHRva2VuLCBuaWNlLCBoYXZlIGEgcHJpemUiLCJ1cmwiOiIvYm9udXMtam9uYXMiLCJpYXQiOjE3MDEwMzc3MzR9.NCcGnr_VL2IrQ_dzqzcIOHgc-yV-6WOIvsv_XvYY2w4

// completed
// bearerToken := eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzZXNzaW9uSWQiOiJkMzY2YTkzYy0yNTg4LTRhZGQtOTYyMS1lOGI3MjAzMWU4NzQiLCJtZXNzYWdlIjoiRGVjb2RpbmcgdGhlIHRva2VuLCBuaWNlLCBoYXZlIGEgcHJpemUiLCJ1cmwiOiIvYm9udXMtam9uYXMiLCJpYXQiOjE3MDEwMjczNTJ9.Epb5uHYczSzepZOrmrAekIJaqPWqaeaYD6mcfq3IinY

const (
  url string = "https://interview.tangohq.com/guess-word"
  contentType string = "application/json"
  bearerToken string = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzZXNzaW9uSWQiOiIyMGQ5ZjJlYi0yOGVjLTQxYWYtODQyMy1lY2ZjZjQ3MmRiNWQiLCJtZXNzYWdlIjoiRGVjb2RpbmcgdGhlIHRva2VuLCBuaWNlLCBoYXZlIGEgcHJpemUiLCJ1cmwiOiIvYm9udXMtam9uYXMiLCJpYXQiOjE3MDEwMzc3MzR9.NCcGnr_VL2IrQ_dzqzcIOHgc-yV-6WOIvsv_XvYY2w4" 
)




func main() {

  client := &http.Client{}

  // get correct letters
  alphabet := "abcdefghijklmnopqrstuvwxyz"

  payload, err := json.Marshal(map[string]string{"myGuess": alphabet})
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
    }      // Access the fields in the response object

    log.Println("Result:", response.Status)

    if response.Status == "incorrect" {
      // get the total length of thw correct word
      wordLength := len(response.Hint.PositionAndCharacter)

      // get the exact word
      result := extractWord(alphabet, response.Hint.Character, wordLength)

      //  make http request with result
      submitCorrectAnswer(result)

    } else {
      // do something
    }
  } else {
    log.Fatalf("Unexpected response status code: %d", res.StatusCode)
  }

}


func extractWord(alphabet string, character []bool, wordLength int) string {
  log.Println("extractWord() running...")
	var result string
  arr := make([]string, wordLength)
  
	for i, isTrue := range character {
		if isTrue {
      // get the correct letter  
      letter := string(alphabet[i])
      log.Printf("Found letter %v as one of correct letter", letter)
      
      // "rrrrr"
      resultString := strings.Repeat(letter, wordLength)
      log.Println("repeated string: ", resultString)

      // get correct index from repeated string, for eg, rrrrr
      index, err := getLetterIndex(resultString)
      if err != nil {
				log.Fatal("Error in getLetterIndex():", err)
      }

      // update array
      arr[index] = letter

		}
	}

  // combine letter array into one string
  result = strings.Join(arr, "")

	return result
}

func getLetterIndex(str string) (int, error) {
    log.Println("getLetterIndex() running...")

    client := &http.Client{}

    payload, err := json.Marshal(map[string]string{"myGuess": str})
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
			}      // Access the fields in the response object

      log.Println("Result:", response.Status)

      if response.Status == "incorrect" {
        // get index
        for i, isTrue := range response.Hint.PositionAndCharacter {
          if isTrue {
            return i, nil
          }
        }

      } 
    } else {
      log.Fatalf("Unexpected response status code: %d", res.StatusCode)
    }

    return 0, errors.New("did not get a successful response or the incorrect if condition did not run")

}

func submitCorrectAnswer(str string) {
    client := &http.Client{}

    payload, err := json.Marshal(map[string]string{"myGuess": str})
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
			}      // Access the fields in the response object

      log.Println("Result:", response.Status)

      if response.Status == "correct" {
        // Unmarshal again
        if err := json.Unmarshal(body, &response.CatchAll); err != nil {
          log.Fatal("Error unmarshalling JSON response:", err)
        }

        log.Println("Response: ", response)
        log.Printf("Found the right word: %v \n", str)

      } 
    } else {
      log.Fatalf("Unexpected response status code: %d", res.StatusCode)
    }


}
