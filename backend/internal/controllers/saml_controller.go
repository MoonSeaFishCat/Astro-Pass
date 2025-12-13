package controllers

import (
	"encoding/base64"
	"net/http"
	"net/url"
	"strconv"

	"astro-pass/internal/services"
	"astro-pass/internal/utils"
	"github.com/gin-gonic/gin"
)

type SAMLController struct {
	samlService *services.SAMLService
}

func NewSAMLController() *SAMLController {
	return &SAMLController{
		samlService: services.NewSAMLService(),
	}
}

// GetMetadata 获取SAML元数据
// @Summary 获取SAML元数据
// @Description 获取IdP的SAML元数据XML
// @Tags SAML
// @Produce xml
// @Param entity_id query string false "实体ID"
// @Success 200 {string} string "SAML元数据XML"
// @Router /api/saml/metadata [get]
func (c *SAMLController) GetMetadata(ctx *gin.Context) {
	entityID := ctx.Query("entity_id")
	if entityID == "" {
		entityID = "astro-pass-idp" // 默认实体ID
	}

	metadata, err := c.samlService.GenerateMetadata(entityID)
	if err != nil {
		utils.ErrorResponse(ctx, http.StatusInternalServerError, err.Error())
		return
	}

	ctx.Header("Content-Type", "application/samlmetadata+xml")
	ctx.String(http.StatusOK, metadata)
}

// HandleSSO 处理SAML SSO请求
// @Summary 处理SAML SSO请求
// @Description 处理来自SP的SAML认证请求
// @Tags SAML
// @Accept application/x-www-form-urlencoded
// @Produce html
// @Success 200 {string} string "登录页面或重定向"
// @Router /api/saml/sso [get,post]
func (c *SAMLController) HandleSSO(ctx *gin.Context) {
	var samlRequest, relayState string

	if ctx.Request.Method == "GET" {
		// HTTP-Redirect绑定
		samlRequest = ctx.Query("SAMLRequest")
		relayState = ctx.Query("RelayState")
	} else {
		// HTTP-POST绑定
		samlRequest = ctx.PostForm("SAMLRequest")
		relayState = ctx.PostForm("RelayState")
	}

	if samlRequest == "" {
		utils.BadRequest(ctx, "缺少SAMLRequest参数")
		return
	}

	// 处理SAML请求
	requestModel, err := c.samlService.ProcessAuthnRequest(samlRequest, relayState)
	if err != nil {
		utils.ErrorResponse(ctx, http.StatusBadRequest, err.Error())
		return
	}

	// 检查用户是否已认证
	userID, exists := ctx.Get("user_id")
	if !exists {
		// 用户未认证，重定向到登录页面
		loginURL := "/login?saml_request_id=" + requestModel.RequestID
		if relayState != "" {
			loginURL += "&relay_state=" + url.QueryEscape(relayState)
		}
		ctx.Redirect(http.StatusFound, loginURL)
		return
	}

	// 用户已认证，生成断言和响应
	c.generateAndSendResponse(ctx, requestModel.RequestID, userID.(uint), relayState)
}

// HandleSAMLLogin 处理SAML登录完成
// @Summary 处理SAML登录完成
// @Description 用户登录后处理SAML请求
// @Tags SAML
// @Security BearerAuth
// @Produce html
// @Param request_id query string true "SAML请求ID"
// @Param relay_state query string false "中继状态"
// @Success 200 {string} string "SAML响应表单"
// @Router /api/saml/login-complete [get]
func (c *SAMLController) HandleSAMLLogin(ctx *gin.Context) {
	userID, exists := ctx.Get("user_id")
	if !exists {
		utils.Unauthorized(ctx, "未认证")
		return
	}

	requestID := ctx.Query("request_id")
	relayState := ctx.Query("relay_state")

	if requestID == "" {
		utils.BadRequest(ctx, "缺少request_id参数")
		return
	}

	c.generateAndSendResponse(ctx, requestID, userID.(uint), relayState)
}

// generateAndSendResponse 生成并发送SAML响应
func (c *SAMLController) generateAndSendResponse(ctx *gin.Context, requestID string, userID uint, relayState string) {
	// 生成断言
	assertion, err := c.samlService.GenerateAssertion(requestID, userID)
	if err != nil {
		utils.ErrorResponse(ctx, http.StatusInternalServerError, err.Error())
		return
	}

	// 生成响应
	responseXML, err := c.samlService.GenerateResponse(assertion.AssertionID)
	if err != nil {
		utils.ErrorResponse(ctx, http.StatusInternalServerError, err.Error())
		return
	}

	// Base64编码响应
	encodedResponse := base64.StdEncoding.EncodeToString([]byte(responseXML))

	// 生成自动提交的HTML表单
	html := `<!DOCTYPE html>
<html>
<head>
    <title>SAML Response</title>
</head>
<body onload="document.forms[0].submit()">
    <form method="post" action="` + assertion.EntityID + `">
        <input type="hidden" name="SAMLResponse" value="` + encodedResponse + `" />
        <input type="hidden" name="RelayState" value="` + relayState + `" />
        <noscript>
            <p>JavaScript is disabled. Please click the button below to continue.</p>
            <input type="submit" value="Continue" />
        </noscript>
    </form>
</body>
</html>`

	ctx.Header("Content-Type", "text/html")
	ctx.String(http.StatusOK, html)
}

// CreateSAMLConfigRequest 创建SAML配置请求
type CreateSAMLConfigRequest struct {
	EntityID    string `json:"entity_id" binding:"required"`
	Type        string `json:"type" binding:"required,oneof=idp sp"`
	Name        string `json:"name" binding:"required"`
	Description string `json:"description"`
}

// CreateSAMLConfig 创建SAML配置
// @Summary 创建SAML配置
// @Description 创建新的SAML IdP或SP配置
// @Tags SAML
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param request body CreateSAMLConfigRequest true "SAML配置信息"
// @Success 201 {object} map[string]interface{}
// @Router /api/admin/saml/configs [post]
func (c *SAMLController) CreateSAMLConfig(ctx *gin.Context) {
	var req CreateSAMLConfigRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(ctx, "请求参数错误: "+err.Error())
		return
	}

	config, err := c.samlService.CreateSAMLConfig(req.EntityID, req.Type, req.Name, req.Description)
	if err != nil {
		utils.ErrorResponse(ctx, http.StatusBadRequest, err.Error())
		return
	}

	utils.SuccessWithMessage(ctx, "SAML配置创建成功", gin.H{
		"id":        config.ID,
		"entity_id": config.EntityID,
		"type":      config.Type,
		"name":      config.Name,
		"status":    config.Status,
	})
}

// GetSAMLConfigs 获取SAML配置列表
// @Summary 获取SAML配置列表
// @Description 获取所有SAML配置
// @Tags SAML
// @Security BearerAuth
// @Produce json
// @Success 200 {object} map[string]interface{}
// @Router /api/admin/saml/configs [get]
func (c *SAMLController) GetSAMLConfigs(ctx *gin.Context) {
	configs, err := c.samlService.GetSAMLConfigs()
	if err != nil {
		utils.ErrorResponse(ctx, http.StatusInternalServerError, err.Error())
		return
	}

	// 格式化响应数据（隐藏敏感信息）
	configData := make([]gin.H, 0, len(configs))
	for _, config := range configs {
		configData = append(configData, gin.H{
			"id":                     config.ID,
			"entity_id":              config.EntityID,
			"type":                   config.Type,
			"name":                   config.Name,
			"description":            config.Description,
			"status":                 config.Status,
			"idp_sso_service_url":    config.IDPSSOServiceURL,
			"idp_slo_service_url":    config.IDPSLOServiceURL,
			"sign_assertions":        config.SignAssertions,
			"encrypt_assertions":     config.EncryptAssertions,
			"sign_requests":          config.SignRequests,
			"created_at":             config.CreatedAt,
			"updated_at":             config.UpdatedAt,
		})
	}

	utils.Success(ctx, gin.H{
		"configs": configData,
		"total":   len(configData),
	})
}

// UpdateSAMLConfigRequest 更新SAML配置请求
type UpdateSAMLConfigRequest struct {
	Name              string `json:"name"`
	Description       string `json:"description"`
	Status            string `json:"status"`
	SignAssertions    *bool  `json:"sign_assertions"`
	EncryptAssertions *bool  `json:"encrypt_assertions"`
	SignRequests      *bool  `json:"sign_requests"`
}

// UpdateSAMLConfig 更新SAML配置
// @Summary 更新SAML配置
// @Description 更新指定的SAML配置
// @Tags SAML
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param id path string true "配置ID"
// @Param request body UpdateSAMLConfigRequest true "更新信息"
// @Success 200 {object} map[string]interface{}
// @Router /api/admin/saml/configs/{id} [put]
func (c *SAMLController) UpdateSAMLConfig(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		utils.BadRequest(ctx, "无效的配置ID")
		return
	}

	var req UpdateSAMLConfigRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(ctx, "请求参数错误: "+err.Error())
		return
	}

	// 构建更新数据
	updates := make(map[string]interface{})
	if req.Name != "" {
		updates["name"] = req.Name
	}
	if req.Description != "" {
		updates["description"] = req.Description
	}
	if req.Status != "" {
		updates["status"] = req.Status
	}
	if req.SignAssertions != nil {
		updates["sign_assertions"] = *req.SignAssertions
	}
	if req.EncryptAssertions != nil {
		updates["encrypt_assertions"] = *req.EncryptAssertions
	}
	if req.SignRequests != nil {
		updates["sign_requests"] = *req.SignRequests
	}

	if len(updates) == 0 {
		utils.BadRequest(ctx, "没有提供更新数据")
		return
	}

	err = c.samlService.UpdateSAMLConfig(uint(id), updates)
	if err != nil {
		utils.ErrorResponse(ctx, http.StatusBadRequest, err.Error())
		return
	}

	utils.SuccessWithMessage(ctx, "SAML配置更新成功", nil)
}

// DeleteSAMLConfig 删除SAML配置
// @Summary 删除SAML配置
// @Description 删除指定的SAML配置
// @Tags SAML
// @Security BearerAuth
// @Produce json
// @Param id path string true "配置ID"
// @Success 200 {object} map[string]interface{}
// @Router /api/admin/saml/configs/{id} [delete]
func (c *SAMLController) DeleteSAMLConfig(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		utils.BadRequest(ctx, "无效的配置ID")
		return
	}

	err = c.samlService.DeleteSAMLConfig(uint(id))
	if err != nil {
		utils.ErrorResponse(ctx, http.StatusBadRequest, err.Error())
		return
	}

	utils.SuccessWithMessage(ctx, "SAML配置删除成功", nil)
}