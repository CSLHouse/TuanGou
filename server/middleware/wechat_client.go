package middleware

import (
	"bytes"
	"cooller/server/global"
	systemRes "cooller/server/model/system/response"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"time"
)

type WechatClient struct {
	httpClient *http.Client
}

func NewWechatClient(httpClient *http.Client) *WechatClient {
	if httpClient == nil {
		httpClient = http.DefaultClient
		httpClient.Timeout = time.Second * 5
	}

	return &WechatClient{
		httpClient: httpClient,
	}
}

func (pc *WechatClient) WXLogin(jscode string) (wxMap map[string]string, error error) {

	loginUrl := fmt.Sprintf(global.GVA_CONFIG.Wechat.SessionUrl,
		url.QueryEscape(global.GVA_CONFIG.Wechat.AppID), url.QueryEscape(global.GVA_CONFIG.Wechat.Secret), url.QueryEscape(jscode))
	fmt.Println("----loginUrl:", loginUrl)
	httpResp, err := pc.httpClient.Get(loginUrl)

	if err != nil {
		return wxMap, err
	}
	defer httpResp.Body.Close()

	if httpResp.StatusCode != http.StatusOK {
		return wxMap, fmt.Errorf("http.Status: %s", httpResp.Status)
	}
	err = json.NewDecoder(httpResp.Body).Decode(&wxMap)
	fmt.Println("----wxMap:", wxMap)

	if err != nil {
		fmt.Println("----err:", err)
	}

	return wxMap, nil
}

func (pc *WechatClient) GetWXAccessToken() (accessToken map[string]any, err error) {
	accessTokenUrl := fmt.Sprintf(global.GVA_CONFIG.Wechat.AccessTokenUrl,
		url.QueryEscape(global.GVA_CONFIG.Wechat.AppID), url.QueryEscape(global.GVA_CONFIG.Wechat.Secret))
	fmt.Println("----accessTokenUrl:", accessTokenUrl)
	httpResp, err := pc.httpClient.Get(accessTokenUrl)

	if err != nil {
		return accessToken, err
	}
	defer httpResp.Body.Close()

	if httpResp.StatusCode != http.StatusOK {
		return accessToken, fmt.Errorf("http.Status: %s", httpResp.Status)
	}
	err = json.NewDecoder(httpResp.Body).Decode(&accessToken)
	fmt.Println("----wxMap:", accessToken)

	if err != nil {
		fmt.Println("----err:", err)
	}
	return accessToken, err
}

func (pc *WechatClient) GetWXTelephone(accessToken string, code string) (telephoneData systemRes.PhoneModel, err error) {
	telephoneUrl := fmt.Sprintf(global.GVA_CONFIG.Wechat.TelephoneUrl,
		accessToken)
	fmt.Println("----telephoneUrl:", telephoneUrl)

	reqData := make(map[string]interface{})
	reqData["code"] = code
	data, err := json.Marshal(reqData)
	httpResp, err := pc.httpClient.Post(telephoneUrl,
		"application/json;", bytes.NewBuffer(data))

	if err != nil {
		return telephoneData, err
	}
	defer httpResp.Body.Close()

	if httpResp.StatusCode != http.StatusOK {
		return telephoneData, fmt.Errorf("http.Status: %s", httpResp.Status)
	}

	err = json.NewDecoder(httpResp.Body).Decode(&telephoneData)
	if err != nil {
		fmt.Println("--[GetWXTelephone]--err:", err)
		return telephoneData, err
	}
	var res map[string]any
	_ = json.NewDecoder(httpResp.Body).Decode(&res)
	fmt.Println("----res:", res)
	return telephoneData, err
}
