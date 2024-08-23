package helpers

import (
	"io"
	"log"
	"net/http"
)

func FetchContent(url string) (string, error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return "", err
	}

	// TODO: Move this to a separate function?
	req.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10.15; rv:128.0) Gecko/20100101 Firefox/128.0")

	client := &http.Client{}
	resp, err := client.Do(req)

	if err != nil {
		log.Println("Error get url: ", err)
		return "", err
	}

	// dispose response when leaving func
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		log.Println("Error status code: ", err)
		log.Println("Status code: ", resp.StatusCode)
		return "", err
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Println("Error readall: ", err)
		return "", err
	}

	return string(body), nil
}
