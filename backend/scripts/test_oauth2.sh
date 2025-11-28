#!/bin/bash

# Astro-Pass OAuth2/OIDC 功能测试脚本
# 用于快速验证新实现的功能

set -e

BASE_URL="http://localhost:8080"
API_URL="$BASE_URL/api"

echo "🧪 Astro-Pass OAuth2/OIDC 功能测试"
echo "=================================="
echo ""

# 颜色定义
GREEN='\033[0;32m'
RED='\033[0;31m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# 测试结果统计
PASSED=0
FAILED=0

# 测试函数
test_endpoint() {
    local name=$1
    local method=$2
    local url=$3
    local data=$4
    local expected_code=$5
    
    echo -n "测试: $name ... "
    
    if [ "$method" = "GET" ]; then
        response=$(curl -s -w "\n%{http_code}" "$url")
    else
        response=$(curl -s -w "\n%{http_code}" -X "$method" "$url" -H "Content-Type: application/json" -d "$data")
    fi
    
    http_code=$(echo "$response" | tail -n1)
    body=$(echo "$response" | sed '$d')
    
    if [ "$http_code" = "$expected_code" ]; then
        echo -e "${GREEN}✓ 通过${NC} (HTTP $http_code)"
        PASSED=$((PASSED + 1))
        return 0
    else
        echo -e "${RED}✗ 失败${NC} (期望 $expected_code, 实际 $http_code)"
        echo "响应: $body"
        FAILED=$((FAILED + 1))
        return 1
    fi
}

echo "📋 测试 1: OIDC 自动发现端点"
echo "----------------------------"
test_endpoint "OIDC Discovery" "GET" "$BASE_URL/.well-known/openid-configuration" "" "200"
echo ""

echo "📋 测试 2: JWKS 公钥端点"
echo "----------------------------"
test_endpoint "JWKS Endpoint" "GET" "$API_URL/oauth2/jwks" "" "200"
echo ""

echo "📋 测试 3: 健康检查"
echo "----------------------------"
test_endpoint "Health Check" "GET" "$BASE_URL/health" "" "200"
test_endpoint "Ready Check" "GET" "$BASE_URL/ready" "" "200"
echo ""

echo "📋 测试 4: 用户注册"
echo "----------------------------"
RANDOM_USER="testuser_$(date +%s)"
REGISTER_DATA="{\"username\":\"$RANDOM_USER\",\"email\":\"$RANDOM_USER@test.com\",\"password\":\"Test123456\",\"nickname\":\"测试用户\"}"
test_endpoint "User Registration" "POST" "$API_URL/auth/register" "$REGISTER_DATA" "200"
echo ""

echo "📋 测试 5: 用户登录"
echo "----------------------------"
LOGIN_DATA="{\"username\":\"$RANDOM_USER\",\"password\":\"Test123456\"}"
LOGIN_RESPONSE=$(curl -s -X POST "$API_URL/auth/login" -H "Content-Type: application/json" -d "$LOGIN_DATA")
ACCESS_TOKEN=$(echo "$LOGIN_RESPONSE" | grep -o '"access_token":"[^"]*' | cut -d'"' -f4)

if [ -n "$ACCESS_TOKEN" ]; then
    echo -e "${GREEN}✓ 登录成功${NC}"
    echo "Access Token: ${ACCESS_TOKEN:0:20}..."
    PASSED=$((PASSED + 1))
else
    echo -e "${RED}✗ 登录失败${NC}"
    echo "响应: $LOGIN_RESPONSE"
    FAILED=$((FAILED + 1))
fi
echo ""

if [ -n "$ACCESS_TOKEN" ]; then
    echo "📋 测试 6: 创建 OAuth2 客户端"
    echo "----------------------------"
    CLIENT_DATA="{\"client_name\":\"测试应用\",\"client_uri\":\"http://localhost:3001\",\"logo_uri\":\"http://localhost:3001/logo.png\",\"redirect_uris\":[\"http://localhost:3001/callback\"]}"
    CLIENT_RESPONSE=$(curl -s -X POST "$API_URL/oauth2/clients" \
        -H "Authorization: Bearer $ACCESS_TOKEN" \
        -H "Content-Type: application/json" \
        -d "$CLIENT_DATA")
    
    CLIENT_ID=$(echo "$CLIENT_RESPONSE" | grep -o '"client_id":"[^"]*' | cut -d'"' -f4)
    CLIENT_SECRET=$(echo "$CLIENT_RESPONSE" | grep -o '"client_secret":"[^"]*' | cut -d'"' -f4)
    
    if [ -n "$CLIENT_ID" ]; then
        echo -e "${GREEN}✓ 客户端创建成功${NC}"
        echo "Client ID: $CLIENT_ID"
        echo "Client Secret: ${CLIENT_SECRET:0:20}..."
        PASSED=$((PASSED + 1))
    else
        echo -e "${RED}✗ 客户端创建失败${NC}"
        echo "响应: $CLIENT_RESPONSE"
        FAILED=$((FAILED + 1))
    fi
    echo ""
    
    if [ -n "$CLIENT_ID" ]; then
        echo "📋 测试 7: Token 内省"
        echo "----------------------------"
        INTROSPECT_RESPONSE=$(curl -s -X POST "$API_URL/oauth2/introspect" \
            -d "token=$ACCESS_TOKEN" \
            -d "client_id=$CLIENT_ID" \
            -d "client_secret=$CLIENT_SECRET")
        
        IS_ACTIVE=$(echo "$INTROSPECT_RESPONSE" | grep -o '"active":[^,}]*' | cut -d':' -f2)
        
        if [ "$IS_ACTIVE" = "true" ]; then
            echo -e "${GREEN}✓ Token 内省成功${NC}"
            echo "Token 状态: 有效"
            PASSED=$((PASSED + 1))
        else
            echo -e "${RED}✗ Token 内省失败${NC}"
            echo "响应: $INTROSPECT_RESPONSE"
            FAILED=$((FAILED + 1))
        fi
        echo ""
    fi
fi

echo "=================================="
echo "📊 测试结果统计"
echo "=================================="
echo -e "通过: ${GREEN}$PASSED${NC}"
echo -e "失败: ${RED}$FAILED${NC}"
echo "总计: $((PASSED + FAILED))"
echo ""

if [ $FAILED -eq 0 ]; then
    echo -e "${GREEN}🎉 所有测试通过！${NC}"
    exit 0
else
    echo -e "${RED}❌ 部分测试失败${NC}"
    exit 1
fi
