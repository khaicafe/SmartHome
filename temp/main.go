// File: main.go
package main

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"sort"
	"strings"
	"time"
)

const (
	Host     = "https://openapi.tuyaus.com"
	ClientID = "gcgndj7kfxkfp4y5resp"
	Secret   = "e585ba9ab46645359a141e022a58c7e7"
	DeviceID = "eb9e304aa409326e51odry"
)

var Token string

type TokenResponse struct {
	Result struct {
		AccessToken  string `json:"access_token"`
		ExpireTime   int    `json:"expire_time"`
		RefreshToken string `json:"refresh_token"`
		UID          string `json:"uid"`
	} `json:"result"`
	Success bool  `json:"success"`
	T       int64 `json:"t"`
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("❗ Vui lòng truyền đối số: on hoặc off")
		return
	}

	action := os.Args[1] == "on"

	GetToken()
	SendCommand(DeviceID, "switch_1", action)
}

func GetToken() {
	method := "GET"
	body := []byte(``)
	req, _ := http.NewRequest(method, Host+"/v1.0/token?grant_type=1", bytes.NewReader(body))
	buildHeader(req, body)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Println(err)
		return
	}
	defer resp.Body.Close()
	bs, _ := ioutil.ReadAll(resp.Body)
	ret := TokenResponse{}
	json.Unmarshal(bs, &ret)
	log.Println("🔐 Token response:", string(bs))

	if v := ret.Result.AccessToken; v != "" {
		Token = v
	}
}

func SendCommand(deviceId, code string, value bool) {
	method := "POST"
	bodyMap := map[string]interface{}{
		"commands": []map[string]interface{}{
			{
				"code":  code,
				"value": value,
			},
		},
	}
	bodyBytes, _ := json.Marshal(bodyMap)

	req, _ := http.NewRequest(method, Host+"/v1.0/devices/"+deviceId+"/commands", bytes.NewReader(bodyBytes))
	req.Header.Set("Content-Type", "application/json")
	buildHeader(req, bodyBytes)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Println(err)
		return
	}
	defer resp.Body.Close()
	bs, _ := ioutil.ReadAll(resp.Body)
	log.Println("🔁 Phản hồi từ thiết bị:", string(bs))
}

func buildHeader(req *http.Request, body []byte) {
	req.Header.Set("client_id", ClientID)
	req.Header.Set("sign_method", "HMAC-SHA256")

	ts := fmt.Sprint(time.Now().UnixNano() / 1e6)
	req.Header.Set("t", ts)

	if Token != "" {
		req.Header.Set("access_token", Token)
	}

	sign := buildSign(req, body, ts)
	req.Header.Set("sign", sign)
}

func buildSign(req *http.Request, body []byte, t string) string {
	headers := getHeaderStr(req)
	urlStr := getUrlStr(req)
	contentSha256 := Sha256(body)
	stringToSign := req.Method + "\n" + contentSha256 + "\n" + headers + "\n" + urlStr
	signStr := ClientID + Token + t + stringToSign
	sign := strings.ToUpper(HmacSha256(signStr, Secret))
	return sign
}

func Sha256(data []byte) string {
	sha256Contain := sha256.New()
	sha256Contain.Write(data)
	return hex.EncodeToString(sha256Contain.Sum(nil))
}

func getUrlStr(req *http.Request) string {
	url := req.URL.Path
	keys := make([]string, 0, 10)

	query := req.URL.Query()
	for key := range query {
		keys = append(keys, key)
	}
	if len(keys) > 0 {
		url += "?"
		sort.Strings(keys)
		for _, keyName := range keys {
			value := query.Get(keyName)
			url += keyName + "=" + value + "&"
		}
	}

	if url[len(url)-1] == '&' {
		url = url[:len(url)-1]
	}
	return url
}

func getHeaderStr(req *http.Request) string {
	signHeaderKeys := req.Header.Get("Signature-Headers")
	if signHeaderKeys == "" {
		return ""
	}
	keys := strings.Split(signHeaderKeys, ":")
	headers := ""
	for _, key := range keys {
		headers += key + ":" + req.Header.Get(key) + "\n"
	}
	return headers
}

func HmacSha256(message string, secret string) string {
	key := []byte(secret)
	h := hmac.New(sha256.New, key)
	h.Write([]byte(message))
	sha := hex.EncodeToString(h.Sum(nil))
	return sha
}
