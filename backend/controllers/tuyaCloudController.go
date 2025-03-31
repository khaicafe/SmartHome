// File: main.go
package controllers

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"go-react-app/models"
	"go-react-app/utils"
	"io/ioutil"
	"log"
	"net/http"
	"os/exec"
	"runtime"
	"sort"
	"strings"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
)

const (
	Host     = "https://openapi.tuyaus.com"
	ClientID = "gcgndj7kfxkfp4y5resp"
	Secret   = "e585ba9ab46645359a141e022a58c7e7"
	// DeviceID = "eb9e304aa409326e51odry"
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

// ////////////// actions cloud ///////
func ActionSendCommand(c *gin.Context) {
	deviceID := c.Param("id")

	var body struct {
		Code  string      `json:"code"`
		Value interface{} `json:"value"`
	}

	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid body"})
		return
	}
	valueBool, ok := body.Value.(bool)
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"error": "value must be boolean"})
		return
	}

	go SendCommand(deviceID, body.Code, valueBool)

	c.JSON(http.StatusOK, gin.H{
		"status": "Command sent",
		"device": deviceID,
		"code":   body.Code,
		"value":  body.Value,
	})
}

// func TurnOn(c *gin.Context) {
// 	go SendCommand(DeviceID, "switch_1", true)
// 	c.JSON(http.StatusOK, gin.H{"status": "Turning ON..."})
// }

// func TurnOff(c *gin.Context) {
// 	go SendCommand(DeviceID, "switch_1", false)
// 	c.JSON(http.StatusOK, gin.H{"status": "Turning OFF..."})
// }

///////// get list device //////////

type DeviceListResponse struct {
	Result struct {
		Devices []Device `json:"devices"`
	} `json:"result"`
	Success bool `json:"success"`
}

type Device struct {
	Name     string `json:"name"`
	DeviceID string `json:"id"`
	Model    string `json:"model"`
}

type DeviceFunctionsResponse struct {
	Result struct {
		Functions []struct {
			Code string `json:"code"`
			Name string `json:"name"`
			Type string `json:"type"`
		} `json:"functions"`
	} `json:"result"`
	Success bool `json:"success"`
}

// GET /api/devices
func GetDevices(c *gin.Context) {
	url := Host + "/v2.0/cloud/thing/device?page_size=20"
	req, _ := http.NewRequest("GET", url, nil)
	buildHeader(req, nil) // nhớ: phải có access_token và tính sign đúng

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	defer resp.Body.Close()

	body, _ := ioutil.ReadAll(resp.Body)
	log.Println("📦 Tuya Cloud device list:", string(body))

	var result map[string]interface{}
	if err := json.Unmarshal(body, &result); err != nil {
		c.JSON(500, gin.H{"error": "decode failed"})
		return
	}
	c.JSON(200, result)
}

// GET /api/device/:id/functions
func GetDeviceFunctions(c *gin.Context) {
	deviceID := c.Param("id")
	url := fmt.Sprintf("%s/v2.0/cloud/thing/%s/shadow/properties", Host, deviceID)

	req, _ := http.NewRequest("GET", url, nil)
	buildHeader(req, nil)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	defer resp.Body.Close()

	body, _ := ioutil.ReadAll(resp.Body)
	log.Println("📦 Device Functions:", string(body))

	var result map[string]interface{}
	if err := json.Unmarshal(body, &result); err != nil {
		c.JSON(500, gin.H{"error": "decode error"})
		return
	}
	c.JSON(200, result)
}

func MapSwitch(c *gin.Context) {
	var input models.MappedSwitch
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid data"})
		return
	}

	if err := models.DB.Create(&input).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "save failed"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "mapped successfully"})
}

func GetMappedSwitches(c *gin.Context) {
	var switches []models.MappedSwitch
	if err := models.DB.Find(&switches).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "load failed"})
		return
	}
	c.JSON(http.StatusOK, switches)
}

func UpdateMappedSwitch(c *gin.Context) {
	id := c.Param("id")
	var input struct {
		Name string `json:"name"`
		IP   string `json:"ip"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid input"})
		return
	}

	var mapped models.MappedSwitch
	if err := models.DB.First(&mapped, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "not found"})
		return
	}

	mapped.Name = input.Name
	mapped.IP = input.IP

	if err := models.DB.Save(&mapped).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "update failed"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "updated"})
}

func DeleteMappedSwitch(c *gin.Context) {
	id := c.Param("id")

	if err := models.DB.Delete(&models.MappedSwitch{}, id).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "delete failed"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "deleted"})
}

// ///////// job go func ////////////
var (
	// pingInterval    = 10 * time.Second      // thời gian giữa mỗi vòng
	// pingConsistency = 5                     // số lần ping mỗi IP
	// maxConcurrent   = 10                    // giới hạn số switch xử lý song song
	switchState     = make(map[string]bool) // lưu trạng thái hiện tại
	switchStateLock = sync.RWMutex{}
)

// ✅ Hàm để reset từ bên ngoài
func ResetSwitchState() {
	switchStateLock.Lock()
	defer switchStateLock.Unlock()
	switchState = make(map[string]bool)
	log.Println("🔄 Đã reset switchState từ UI")
}
func ResetSwitchStateHandler(c *gin.Context) {
	ResetSwitchState()
	c.JSON(http.StatusOK, gin.H{"message": "Đã reset trạng thái switch thành công"})
}

// Gọi từ main.go
func StartPingLoop() {
	go func() {
		for {
			// ⚠️ Lấy lại giá trị mỗi vòng để cập nhật nếu user đã thay đổi từ UI
			pingInterval := utils.GetDurationSetting("pingInterval", 10*time.Second)
			pingConsistency := utils.GetIntSetting("pingConsistency", 5)
			maxConcurrent := utils.GetIntSetting("maxConcurrent", 10)

			log.Println("🔁 Bắt đầu vòng kiểm tra IP...")

			var switches []models.MappedSwitch
			if err := models.DB.Find(&switches).Error; err != nil {
				log.Println("❌ Lỗi DB:", err)
				time.Sleep(pingInterval)
				continue
			}

			var wg sync.WaitGroup
			sem := make(chan struct{}, maxConcurrent)

			for _, s := range switches {
				wg.Add(1)
				sem <- struct{}{} // chiếm 1 slot

				go func(s models.MappedSwitch) {
					defer wg.Done()
					defer func() { <-sem }() // trả lại slot

					result := PingHost(s.IP, pingConsistency)
					key := fmt.Sprintf("%s_%s", s.DeviceID, s.Code)
					log.Printf("📡 %s (%s) → Ping tổng hợp: %v", s.Name, s.IP, result)

					switchStateLock.RLock()
					prev, existed := switchState[key]
					switchStateLock.RUnlock()

					if !existed || prev != result {
						switchStateLock.Lock()
						switchState[key] = result
						switchStateLock.Unlock()
						if result {
							log.Printf("✅ %s → Trạng thái mới: ON → gửi lệnh", key)
							go SendCommand(s.DeviceID, s.Code, true)
						} else {
							log.Printf("⛔ %s → Trạng thái mới: OFF → gửi lệnh", key)
							go SendCommand(s.DeviceID, s.Code, false)
						}
					} else {
						log.Printf("⚠️  %s → Không thay đổi trạng thái (%v) → bỏ qua", key, result)
					}
				}(s)
			}

			wg.Wait()
			log.Println("🕓 Kết thúc vòng kiểm tra. Chờ vòng tiếp theo...\n")
			time.Sleep(pingInterval)
		}
	}()
}

func PingHost(ip string, count int) bool {
	successCount := 0

	for i := 0; i < count; i++ {
		var cmd *exec.Cmd
		if runtime.GOOS == "windows" {
			cmd = exec.Command("ping", "-n", "1", "-w", "1000", ip)
		} else {
			cmd = exec.Command("ping", "-c", "1", "-W", "1", ip)
		}

		out, err := cmd.CombinedOutput()
		output := string(out)

		log.Printf("📥 Ping lần %d đến %s: %v\n", i+1, ip, err == nil)
		log.Printf("📄 Output: %s\n", strings.TrimSpace(output))

		if err == nil && (strings.Contains(output, "ttl=") || strings.Contains(output, "TTL=")) {
			successCount++
		}

		time.Sleep(1000 * time.Millisecond)
	}

	log.Printf("✅ Tổng kết ping %s: %d/%d thành công", ip, successCount, count)

	return successCount == count
}
