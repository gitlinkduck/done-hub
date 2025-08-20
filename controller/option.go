package controller

import (
	"done-hub/common/config"
	"done-hub/common/utils"
	"done-hub/model"
	"done-hub/safty"
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

func GetOptions(c *gin.Context) {
	var options []*model.Option
	for k, v := range config.GlobalOption.GetAll() {
		if strings.HasSuffix(k, "Token") || strings.HasSuffix(k, "Secret") {
			continue
		}
		options = append(options, &model.Option{
			Key:   k,
			Value: utils.Interface2String(v),
		})
	}
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "",
		"data":    options,
	})
	return
}

func GetSafeTools(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "",
		"data":    safty.GetAllSafeToolsName(),
	})
	return
}

func UpdateOption(c *gin.Context) {
	var option model.Option
	err := json.NewDecoder(c.Request.Body).Decode(&option)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "无效的参数",
		})
		return
	}
	switch option.Key {
	case "GitHubOAuthEnabled":
		if option.Value == "true" && config.GitHubClientId == "" {
			c.JSON(http.StatusOK, gin.H{
				"success": false,
				"message": "无法启用 GitHub OAuth，请先填入 GitHub Client Id 以及 GitHub Client Secret！",
			})
			return
		}
	case "OIDCAuthEnabled":
		if option.Value == "true" && (config.OIDCClientId == "" || config.OIDCClientSecret == "" || config.OIDCIssuer == "" || config.OIDCScopes == "" || config.OIDCUsernameClaims == "") {
			c.JSON(http.StatusOK, gin.H{
				"success": false,
				"message": "无法启用 OIDC，请先填入OIDC信息！",
			})
			return
		}
	case "LinuxDoOAuthEnabled":
		if option.Value == "true" && (config.LinuxDoClientId == "" || config.LinuxDoClientSecret == "") {
			c.JSON(http.StatusOK, gin.H{
				"success": false,
				"message": "无法启用 LINUX DO OAuth，请先填入 LINUX DO Client Id 以及 LINUX DO Client Secret！",
			})
			return
		}
	case "LinuxDoOAuthTrustLevelEnabled":
		if option.Value == "true" && config.LinuxDoOAuthEnabled == false {
			c.JSON(http.StatusOK, gin.H{
				"success": false,
				"message": "无法启用 LINUX DO 信用等级限制，请先启用 LINUX DO OAuth ！",
			})
			return
		}
	case "LinuxDoOAuthLowestTrustLevel":
		lowestTrustLevel, err := strconv.Atoi(option.Value)
		if err != nil {
			c.JSON(http.StatusOK, gin.H{
				"success": false,
				"message": "LINUX DO 信任等级必须为数字",
			})
			return
		}
		if lowestTrustLevel < config.Basic || lowestTrustLevel > config.Leader {
			c.JSON(http.StatusOK, gin.H{
				"success": false,
				"message": "LINUX DO 信任等级必须 1~4 之间",
			})
			return
		}

	case "EmailDomainRestrictionEnabled":
		if option.Value == "true" && len(config.EmailDomainWhitelist) == 0 {
			c.JSON(http.StatusOK, gin.H{
				"success": false,
				"message": "无法启用邮箱域名限制，请先填入限制的邮箱域名！",
			})
			return
		}
	case "WeChatAuthEnabled":
		if option.Value == "true" && config.WeChatServerAddress == "" {
			c.JSON(http.StatusOK, gin.H{
				"success": false,
				"message": "无法启用微信登录，请先填入微信登录相关配置信息！",
			})
			return
		}
	case "TurnstileCheckEnabled":
		if option.Value == "true" && config.TurnstileSiteKey == "" {
			c.JSON(http.StatusOK, gin.H{
				"success": false,
				"message": "无法启用 Turnstile 校验，请先填入 Turnstile 校验相关配置信息！",
			})
			return
		}
	case "InviterRewardValue":
		value, err := strconv.Atoi(option.Value)
		if err != nil {
			c.JSON(http.StatusOK, gin.H{
				"success": false,
				"message": "充值返利值必须是有效的数字",
			})
			return
		}

		// 验证充值返利值的范围
		if config.InviterRewardType == "percentage" {
			// 百分比类型：值应在0-100之间
			if value < 0 || value > 100 {
				c.JSON(http.StatusOK, gin.H{
					"success": false,
					"message": "当充值返利类型为百分比时，返利值应在0-100之间",
				})
				return
			}
		} else {
			// 固定类型：值应>=0
			if value < 0 {
				c.JSON(http.StatusOK, gin.H{
					"success": false,
					"message": "当充值返利类型为固定时，返利值应大于等于0",
				})
				return
			}
		}
	}
	err = model.UpdateOption(option.Key, option.Value)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"success": false,
			"message": err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "",
	})
	return
}
