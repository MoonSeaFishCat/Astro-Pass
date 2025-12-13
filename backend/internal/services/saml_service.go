package services

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/base64"
	"encoding/pem"
	"encoding/xml"
	"fmt"
	"math/big"
	"time"

	"astro-pass/internal/config"
	"astro-pass/internal/database"
	"astro-pass/internal/models"
	"github.com/google/uuid"
)

type SAMLService struct{}

func NewSAMLService() *SAMLService {
	return &SAMLService{}
}

// SAML XML结构定义
type SAMLMetadata struct {
	XMLName                xml.Name `xml:"urn:oasis:names:tc:SAML:2.0:metadata EntityDescriptor"`
	EntityID               string   `xml:"entityID,attr"`
	IDPSSODescriptor       *IDPSSODescriptor `xml:"IDPSSODescriptor,omitempty"`
	SPSSODescriptor        *SPSSODescriptor  `xml:"SPSSODescriptor,omitempty"`
}

type IDPSSODescriptor struct {
	XMLName              xml.Name `xml:"urn:oasis:names:tc:SAML:2.0:metadata IDPSSODescriptor"`
	ProtocolSupportEnumeration string `xml:"protocolSupportEnumeration,attr"`
	KeyDescriptor        []KeyDescriptor `xml:"KeyDescriptor"`
	SingleSignOnService  []SingleSignOnService `xml:"SingleSignOnService"`
	SingleLogoutService  []SingleLogoutService `xml:"SingleLogoutService"`
}

type SPSSODescriptor struct {
	XMLName              xml.Name `xml:"urn:oasis:names:tc:SAML:2.0:metadata SPSSODescriptor"`
	ProtocolSupportEnumeration string `xml:"protocolSupportEnumeration,attr"`
	KeyDescriptor        []KeyDescriptor `xml:"KeyDescriptor"`
	AssertionConsumerService []AssertionConsumerService `xml:"AssertionConsumerService"`
	SingleLogoutService  []SingleLogoutService `xml:"SingleLogoutService"`
}

type KeyDescriptor struct {
	XMLName xml.Name `xml:"urn:oasis:names:tc:SAML:2.0:metadata KeyDescriptor"`
	Use     string   `xml:"use,attr"`
	KeyInfo KeyInfo  `xml:"KeyInfo"`
}

type KeyInfo struct {
	XMLName xml.Name `xml:"http://www.w3.org/2000/09/xmldsig# KeyInfo"`
	X509Data X509Data `xml:"X509Data"`
}

type X509Data struct {
	XMLName xml.Name `xml:"http://www.w3.org/2000/09/xmldsig# X509Data"`
	X509Certificate string `xml:"X509Certificate"`
}

type SingleSignOnService struct {
	XMLName  xml.Name `xml:"urn:oasis:names:tc:SAML:2.0:metadata SingleSignOnService"`
	Binding  string   `xml:"Binding,attr"`
	Location string   `xml:"Location,attr"`
}

type SingleLogoutService struct {
	XMLName  xml.Name `xml:"urn:oasis:names:tc:SAML:2.0:metadata SingleLogoutService"`
	Binding  string   `xml:"Binding,attr"`
	Location string   `xml:"Location,attr"`
}

type AssertionConsumerService struct {
	XMLName  xml.Name `xml:"urn:oasis:names:tc:SAML:2.0:metadata AssertionConsumerService"`
	Binding  string   `xml:"Binding,attr"`
	Location string   `xml:"Location,attr"`
	Index    string   `xml:"index,attr"`
}

// AuthnRequest SAML认证请求
type AuthnRequest struct {
	XMLName      xml.Name `xml:"urn:oasis:names:tc:SAML:2.0:protocol AuthnRequest"`
	ID           string   `xml:"ID,attr"`
	Version      string   `xml:"Version,attr"`
	IssueInstant string   `xml:"IssueInstant,attr"`
	Destination  string   `xml:"Destination,attr"`
	Issuer       Issuer   `xml:"Issuer"`
}

type Issuer struct {
	XMLName xml.Name `xml:"urn:oasis:names:tc:SAML:2.0:assertion Issuer"`
	Value   string   `xml:",chardata"`
}

// Response SAML响应
type Response struct {
	XMLName      xml.Name `xml:"urn:oasis:names:tc:SAML:2.0:protocol Response"`
	ID           string   `xml:"ID,attr"`
	Version      string   `xml:"Version,attr"`
	IssueInstant string   `xml:"IssueInstant,attr"`
	Destination  string   `xml:"Destination,attr"`
	InResponseTo string   `xml:"InResponseTo,attr"`
	Issuer       Issuer   `xml:"Issuer"`
	Status       Status   `xml:"Status"`
	Assertion    *Assertion `xml:"Assertion,omitempty"`
}

type Status struct {
	XMLName    xml.Name   `xml:"urn:oasis:names:tc:SAML:2.0:protocol Status"`
	StatusCode StatusCode `xml:"StatusCode"`
}

type StatusCode struct {
	XMLName xml.Name `xml:"urn:oasis:names:tc:SAML:2.0:protocol StatusCode"`
	Value   string   `xml:"Value,attr"`
}

type Assertion struct {
	XMLName            xml.Name           `xml:"urn:oasis:names:tc:SAML:2.0:assertion Assertion"`
	ID                 string             `xml:"ID,attr"`
	Version            string             `xml:"Version,attr"`
	IssueInstant       string             `xml:"IssueInstant,attr"`
	Issuer             Issuer             `xml:"Issuer"`
	Subject            Subject            `xml:"Subject"`
	Conditions         Conditions         `xml:"Conditions"`
	AttributeStatement *AttributeStatement `xml:"AttributeStatement,omitempty"`
	AuthnStatement     AuthnStatement     `xml:"AuthnStatement"`
}

type Subject struct {
	XMLName             xml.Name            `xml:"urn:oasis:names:tc:SAML:2.0:assertion Subject"`
	NameID              NameID              `xml:"NameID"`
	SubjectConfirmation SubjectConfirmation `xml:"SubjectConfirmation"`
}

type NameID struct {
	XMLName xml.Name `xml:"urn:oasis:names:tc:SAML:2.0:assertion NameID"`
	Format  string   `xml:"Format,attr"`
	Value   string   `xml:",chardata"`
}

type SubjectConfirmation struct {
	XMLName                 xml.Name                `xml:"urn:oasis:names:tc:SAML:2.0:assertion SubjectConfirmation"`
	Method                  string                  `xml:"Method,attr"`
	SubjectConfirmationData SubjectConfirmationData `xml:"SubjectConfirmationData"`
}

type SubjectConfirmationData struct {
	XMLName      xml.Name `xml:"urn:oasis:names:tc:SAML:2.0:assertion SubjectConfirmationData"`
	InResponseTo string   `xml:"InResponseTo,attr"`
	NotOnOrAfter string   `xml:"NotOnOrAfter,attr"`
	Recipient    string   `xml:"Recipient,attr"`
}

type Conditions struct {
	XMLName              xml.Name             `xml:"urn:oasis:names:tc:SAML:2.0:assertion Conditions"`
	NotBefore            string               `xml:"NotBefore,attr"`
	NotOnOrAfter         string               `xml:"NotOnOrAfter,attr"`
	AudienceRestriction  AudienceRestriction  `xml:"AudienceRestriction"`
}

type AudienceRestriction struct {
	XMLName  xml.Name `xml:"urn:oasis:names:tc:SAML:2.0:assertion AudienceRestriction"`
	Audience Audience `xml:"Audience"`
}

type Audience struct {
	XMLName xml.Name `xml:"urn:oasis:names:tc:SAML:2.0:assertion Audience"`
	Value   string   `xml:",chardata"`
}

type AttributeStatement struct {
	XMLName   xml.Name    `xml:"urn:oasis:names:tc:SAML:2.0:assertion AttributeStatement"`
	Attribute []Attribute `xml:"Attribute"`
}

type Attribute struct {
	XMLName        xml.Name         `xml:"urn:oasis:names:tc:SAML:2.0:assertion Attribute"`
	Name           string           `xml:"Name,attr"`
	AttributeValue []AttributeValue `xml:"AttributeValue"`
}

type AttributeValue struct {
	XMLName xml.Name `xml:"urn:oasis:names:tc:SAML:2.0:assertion AttributeValue"`
	Value   string   `xml:",chardata"`
}

type AuthnStatement struct {
	XMLName      xml.Name     `xml:"urn:oasis:names:tc:SAML:2.0:assertion AuthnStatement"`
	AuthnInstant string       `xml:"AuthnInstant,attr"`
	AuthnContext AuthnContext `xml:"AuthnContext"`
}

type AuthnContext struct {
	XMLName              xml.Name `xml:"urn:oasis:names:tc:SAML:2.0:assertion AuthnContext"`
	AuthnContextClassRef string   `xml:"AuthnContextClassRef"`
}

// CreateSAMLConfig 创建SAML配置
func (s *SAMLService) CreateSAMLConfig(entityID, configType, name, description string) (*models.SAMLConfig, error) {
	// 生成证书和私钥
	cert, privateKey, err := s.generateSelfSignedCertificate(entityID)
	if err != nil {
		return nil, fmt.Errorf("生成证书失败: %v", err)
	}

	baseURL := config.Cfg.App.URL
	
	samlConfig := &models.SAMLConfig{
		EntityID:             entityID,
		Type:                 configType,
		Name:                 name,
		Description:          description,
		Status:               "active",
		IDPCertificate:       cert,
		IDPPrivateKey:        privateKey,
		IDPSSOServiceURL:     baseURL + "/api/saml/sso",
		IDPSLOServiceURL:     baseURL + "/api/saml/slo",
		SignAssertions:       true,
		EncryptAssertions:    false,
		SignRequests:         false,
	}

	if err := database.DB.Create(samlConfig).Error; err != nil {
		return nil, fmt.Errorf("创建SAML配置失败: %v", err)
	}

	return samlConfig, nil
}

// generateSelfSignedCertificate 生成自签名证书
func (s *SAMLService) generateSelfSignedCertificate(entityID string) (string, string, error) {
	// 生成私钥
	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		return "", "", err
	}

	// 创建证书模板
	template := x509.Certificate{
		SerialNumber: big.NewInt(1),
		Subject: pkix.Name{
			Organization:  []string{"Astro-Pass"},
			Country:       []string{"CN"},
			Province:      []string{""},
			Locality:      []string{""},
			StreetAddress: []string{""},
			PostalCode:    []string{""},
		},
		NotBefore:    time.Now(),
		NotAfter:     time.Now().Add(365 * 24 * time.Hour), // 1年有效期
		KeyUsage:     x509.KeyUsageKeyEncipherment | x509.KeyUsageDigitalSignature,
		ExtKeyUsage:  []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth},
		IPAddresses:  nil,
		DNSNames:     []string{entityID},
	}

	// 生成证书
	certDER, err := x509.CreateCertificate(rand.Reader, &template, &template, &privateKey.PublicKey, privateKey)
	if err != nil {
		return "", "", err
	}

	// 编码证书
	certPEM := pem.EncodeToMemory(&pem.Block{
		Type:  "CERTIFICATE",
		Bytes: certDER,
	})

	// 编码私钥
	privateKeyDER, err := x509.MarshalPKCS8PrivateKey(privateKey)
	if err != nil {
		return "", "", err
	}

	privateKeyPEM := pem.EncodeToMemory(&pem.Block{
		Type:  "PRIVATE KEY",
		Bytes: privateKeyDER,
	})

	return string(certPEM), string(privateKeyPEM), nil
}

// GenerateMetadata 生成SAML元数据
func (s *SAMLService) GenerateMetadata(entityID string) (string, error) {
	var samlConfig models.SAMLConfig
	if err := database.DB.Where("entity_id = ? AND status = ?", entityID, "active").
		First(&samlConfig).Error; err != nil {
		return "", fmt.Errorf("SAML配置不存在")
	}

	// 提取证书内容（去除PEM头尾）
	certPEM := samlConfig.IDPCertificate
	block, _ := pem.Decode([]byte(certPEM))
	if block == nil {
		return "", fmt.Errorf("无效的证书格式")
	}
	certBase64 := base64.StdEncoding.EncodeToString(block.Bytes)

	metadata := SAMLMetadata{
		EntityID: entityID,
	}

	if samlConfig.Type == "idp" {
		metadata.IDPSSODescriptor = &IDPSSODescriptor{
			ProtocolSupportEnumeration: "urn:oasis:names:tc:SAML:2.0:protocol",
			KeyDescriptor: []KeyDescriptor{
				{
					Use: "signing",
					KeyInfo: KeyInfo{
						X509Data: X509Data{
							X509Certificate: certBase64,
						},
					},
				},
			},
			SingleSignOnService: []SingleSignOnService{
				{
					Binding:  "urn:oasis:names:tc:SAML:2.0:bindings:HTTP-POST",
					Location: samlConfig.IDPSSOServiceURL,
				},
				{
					Binding:  "urn:oasis:names:tc:SAML:2.0:bindings:HTTP-Redirect",
					Location: samlConfig.IDPSSOServiceURL,
				},
			},
			SingleLogoutService: []SingleLogoutService{
				{
					Binding:  "urn:oasis:names:tc:SAML:2.0:bindings:HTTP-POST",
					Location: samlConfig.IDPSLOServiceURL,
				},
			},
		}
	}

	xmlData, err := xml.MarshalIndent(metadata, "", "  ")
	if err != nil {
		return "", fmt.Errorf("生成元数据失败: %v", err)
	}

	return xml.Header + string(xmlData), nil
}

// ProcessAuthnRequest 处理SAML认证请求
func (s *SAMLService) ProcessAuthnRequest(samlRequest string, relayState string) (*models.SAMLRequest, error) {
	// 解码SAML请求
	decodedRequest, err := base64.StdEncoding.DecodeString(samlRequest)
	if err != nil {
		return nil, fmt.Errorf("解码SAML请求失败: %v", err)
	}

	// 解析XML
	var authnRequest AuthnRequest
	if err := xml.Unmarshal(decodedRequest, &authnRequest); err != nil {
		return nil, fmt.Errorf("解析SAML请求失败: %v", err)
	}

	// 保存请求
	requestModel := &models.SAMLRequest{
		RequestID:   authnRequest.ID,
		Type:        "AuthnRequest",
		EntityID:    authnRequest.Issuer.Value,
		RelayState:  relayState,
		RequestData: string(decodedRequest),
		Status:      "pending",
		ExpiresAt:   time.Now().Add(10 * time.Minute), // 10分钟过期
	}

	if err := database.DB.Create(requestModel).Error; err != nil {
		return nil, fmt.Errorf("保存SAML请求失败: %v", err)
	}

	return requestModel, nil
}

// GenerateAssertion 生成SAML断言
func (s *SAMLService) GenerateAssertion(requestID string, userID uint) (*models.SAMLAssertion, error) {
	// 获取请求信息
	var samlRequest models.SAMLRequest
	if err := database.DB.Where("request_id = ? AND status = ?", requestID, "pending").
		First(&samlRequest).Error; err != nil {
		return nil, fmt.Errorf("SAML请求不存在或已处理")
	}

	// 获取用户信息
	var user models.User
	if err := database.DB.First(&user, userID).Error; err != nil {
		return nil, fmt.Errorf("用户不存在")
	}

	// 获取SAML配置
	var samlConfig models.SAMLConfig
	if err := database.DB.Where("type = ? AND status = ?", "idp", "active").
		First(&samlConfig).Error; err != nil {
		return nil, fmt.Errorf("IdP配置不存在")
	}

	// 生成断言ID
	assertionID := "_" + uuid.New().String()
	now := time.Now()
	notOnOrAfter := now.Add(5 * time.Minute) // 5分钟有效期

	// 构建断言
	assertion := Assertion{
		ID:           assertionID,
		Version:      "2.0",
		IssueInstant: now.UTC().Format(time.RFC3339),
		Issuer: Issuer{
			Value: samlConfig.EntityID,
		},
		Subject: Subject{
			NameID: NameID{
				Format: "urn:oasis:names:tc:SAML:1.1:nameid-format:emailAddress",
				Value:  user.Email,
			},
			SubjectConfirmation: SubjectConfirmation{
				Method: "urn:oasis:names:tc:SAML:2.0:cm:bearer",
				SubjectConfirmationData: SubjectConfirmationData{
					InResponseTo: requestID,
					NotOnOrAfter: notOnOrAfter.UTC().Format(time.RFC3339),
					Recipient:    samlRequest.EntityID,
				},
			},
		},
		Conditions: Conditions{
			NotBefore:    now.UTC().Format(time.RFC3339),
			NotOnOrAfter: notOnOrAfter.UTC().Format(time.RFC3339),
			AudienceRestriction: AudienceRestriction{
				Audience: Audience{
					Value: samlRequest.EntityID,
				},
			},
		},
		AttributeStatement: &AttributeStatement{
			Attribute: []Attribute{
				{
					Name: "email",
					AttributeValue: []AttributeValue{
						{Value: user.Email},
					},
				},
				{
					Name: "username",
					AttributeValue: []AttributeValue{
						{Value: user.Username},
					},
				},
				{
					Name: "displayName",
					AttributeValue: []AttributeValue{
						{Value: user.Nickname},
					},
				},
			},
		},
		AuthnStatement: AuthnStatement{
			AuthnInstant: now.UTC().Format(time.RFC3339),
			AuthnContext: AuthnContext{
				AuthnContextClassRef: "urn:oasis:names:tc:SAML:2.0:ac:classes:PasswordProtectedTransport",
			},
		},
	}

	// 序列化断言
	assertionXML, err := xml.MarshalIndent(assertion, "", "  ")
	if err != nil {
		return nil, fmt.Errorf("序列化断言失败: %v", err)
	}

	// 保存断言
	assertionModel := &models.SAMLAssertion{
		AssertionID:   assertionID,
		RequestID:     requestID,
		UserID:        userID,
		EntityID:      samlRequest.EntityID,
		AssertionData: string(assertionXML),
		Status:        "active",
		ExpiresAt:     notOnOrAfter,
	}

	if err := database.DB.Create(assertionModel).Error; err != nil {
		return nil, fmt.Errorf("保存断言失败: %v", err)
	}

	// 更新请求状态
	samlRequest.Status = "processed"
	samlRequest.UserID = userID
	database.DB.Save(&samlRequest)

	return assertionModel, nil
}

// GenerateResponse 生成SAML响应
func (s *SAMLService) GenerateResponse(assertionID string) (string, error) {
	// 获取断言
	var assertion models.SAMLAssertion
	if err := database.DB.Where("assertion_id = ? AND status = ?", assertionID, "active").
		Preload("User").
		First(&assertion).Error; err != nil {
		return "", fmt.Errorf("断言不存在或已失效")
	}

	// 获取请求信息
	var samlRequest models.SAMLRequest
	if err := database.DB.Where("request_id = ?", assertion.RequestID).First(&samlRequest).Error; err != nil {
		return "", fmt.Errorf("关联的SAML请求不存在")
	}

	// 获取SAML配置
	var samlConfig models.SAMLConfig
	if err := database.DB.Where("type = ? AND status = ?", "idp", "active").
		First(&samlConfig).Error; err != nil {
		return "", fmt.Errorf("IdP配置不存在")
	}

	responseID := "_" + uuid.New().String()
	now := time.Now()

	// 解析断言XML
	var assertionObj Assertion
	if err := xml.Unmarshal([]byte(assertion.AssertionData), &assertionObj); err != nil {
		return "", fmt.Errorf("解析断言失败: %v", err)
	}

	// 构建响应
	response := Response{
		ID:           responseID,
		Version:      "2.0",
		IssueInstant: now.UTC().Format(time.RFC3339),
		InResponseTo: assertion.RequestID,
		Issuer: Issuer{
			Value: samlConfig.EntityID,
		},
		Status: Status{
			StatusCode: StatusCode{
				Value: "urn:oasis:names:tc:SAML:2.0:status:Success",
			},
		},
		Assertion: &assertionObj,
	}

	// 序列化响应
	responseXML, err := xml.MarshalIndent(response, "", "  ")
	if err != nil {
		return "", fmt.Errorf("序列化响应失败: %v", err)
	}

	// 标记断言为已消费
	assertion.Status = "consumed"
	database.DB.Save(&assertion)

	return xml.Header + string(responseXML), nil
}

// GetSAMLConfigs 获取SAML配置列表
func (s *SAMLService) GetSAMLConfigs() ([]models.SAMLConfig, error) {
	var configs []models.SAMLConfig
	if err := database.DB.Where("status != ?", "deleted").
		Order("created_at DESC").
		Find(&configs).Error; err != nil {
		return nil, fmt.Errorf("获取SAML配置失败: %v", err)
	}
	return configs, nil
}

// UpdateSAMLConfig 更新SAML配置
func (s *SAMLService) UpdateSAMLConfig(id uint, updates map[string]interface{}) error {
	if err := database.DB.Model(&models.SAMLConfig{}).
		Where("id = ?", id).
		Updates(updates).Error; err != nil {
		return fmt.Errorf("更新SAML配置失败: %v", err)
	}
	return nil
}

// DeleteSAMLConfig 删除SAML配置
func (s *SAMLService) DeleteSAMLConfig(id uint) error {
	if err := database.DB.Model(&models.SAMLConfig{}).
		Where("id = ?", id).
		Update("status", "deleted").Error; err != nil {
		return fmt.Errorf("删除SAML配置失败: %v", err)
	}
	return nil
}