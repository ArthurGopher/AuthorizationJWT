package handlers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type Address struct {
	City string `json:"city"`
}
type SearchRequest struct {
	Query string `json:"query"`
}
type SearchResponse struct {
	Addresses []*Address `json:"addresses"`
}
type GeocodeRequest struct {
	Lat string `json:"lat"`
	Lng string `json:"lng"`
}
type GeocodeResponse struct {
	Addresses []*Address `json:"addresses"`
}

// GeocodeHandler выполняет обратное геокодирование по координатам.
// @Summary Обратное геокодирование
// @Description Возвращает адреса по переданным широте и долготе через DaData API.
// @Tags Геокодирование
// @Accept json
// @Produce json
// @Param coordinates body GeocodeRequest true "Широта и долгота"
// @Success 200 {object} GeocodeResponse
// @Failure 400 {string} string "Неверный формат запроса"
// @Failure 500 {string} string "Ошибка вызова API"
// @Router /api/address/geocode [post]
func GeocodeHandler(w http.ResponseWriter, r *http.Request) {
	var req GeocodeRequest
	err := json.NewDecoder(r.Body).Decode(&req)

	if err != nil {
		http.Error(w, "неверный формат запроса", http.StatusBadRequest)
		return
	}
	if req.Lat == "" || req.Lng == "" {
		http.Error(w, "Поля lat и lng пустые", http.StatusBadRequest)
		return
	}

	addresses, err := callDaDataGeocode(req.Lat, req.Lng)
	if err != nil {
		http.Error(w, "ошибка вызова API", http.StatusInternalServerError)
		return
	}

	resp := GeocodeResponse{Addresses: addresses}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(resp)

}

func callDaDataGeocode(lat, lng string) ([]*Address, error) {
	apiURL := "https://suggestions.dadata.ru/suggestions/api/4_1/rs/geolocate/address"
	apiKey := "4d6ca084001e77daa4be473fdab5330a50e1db72"

	bodyReq := map[string]string{"lat": lat, "lon": lng}
	bodyReqBytes, _ := json.Marshal(bodyReq)

	req, err := http.NewRequest("POST", apiURL, bytes.NewBuffer(bodyReqBytes))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Token %s", apiKey))

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API вернул статус: %v", resp.StatusCode)
	}

	body, _ := io.ReadAll(resp.Body)
	var response struct {
		Suggestions []struct {
			Value string `json:"value"`
		} `json:"suggestions"`
	}

	err = json.Unmarshal(body, &response)
	if err != nil {
		return nil, err
	}
	var addresses []*Address

	for _, suggestion := range response.Suggestions {
		addresses = append(addresses, &Address{City: suggestion.Value})
	}

	return addresses, nil

}

// SearchHandler ищет адреса по переданному параметру query.
// @Summary Поиск адресов
// @Description Ищет адреса по переданному параметру query через DaData API.
// @Tags Адреса
// @Accept json
// @Produce json
// @Param query body SearchRequest true "Запрос с адресом"
// @Success 200 {object} SearchResponse
// @Failure 400 {string} string "Неверный формат запроса"
// @Failure 500 {string} string "Ошибка вызова API"
// @Router /api/address/search [post]
func SearchHandler(w http.ResponseWriter, r *http.Request) {
	var req SearchRequest

	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "неверный формат запроса", http.StatusBadRequest)
		return
	}
	if req.Query == "" {
		http.Error(w, "должен принимать POST запросы с параметром query", http.StatusBadRequest)
		return
	}

	addresses, err := callDaDataSearch(req.Query)
	if err != nil {
		http.Error(w, "ошибка вызова API", http.StatusInternalServerError)
		return
	}

	resp := SearchResponse{Addresses: addresses}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(resp)

}

func callDaDataSearch(query string) ([]*Address, error) {
	apiURL := "https://suggestions.dadata.ru/suggestions/api/4_1/rs/suggest/address"
	apiKey := "4d6ca084001e77daa4be473fdab5330a50e1db72"

	bodyReq := map[string]string{"query": query}
	bodyReqBytes, _ := json.Marshal(bodyReq)

	req, err := http.NewRequest("POST", apiURL, bytes.NewBuffer(bodyReqBytes))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Token %s", apiKey))

	client := &http.Client{}

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API вернул статус: %v", resp.StatusCode)
	}

	body, _ := io.ReadAll(resp.Body)
	var response struct {
		Suggestions []struct {
			Value string `json:"value"`
		} `json:"suggestions"`
	}

	err = json.Unmarshal(body, &response)
	if err != nil {
		return nil, err
	}

	var addresses []*Address

	for _, suggestion := range response.Suggestions {
		addresses = append(addresses, &Address{City: suggestion.Value})
	}

	return addresses, nil

}


