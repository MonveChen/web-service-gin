/*
 * @Author: Monve
 * @Date: 2023-07-24 15:25:19
 * @LastEditors: Monve
 * @LastEditTime: 2023-07-25 15:10:31
 * @FilePath: /web-service-gin/api/token/token.go
 */
package token

import (
	"database/sql"
	"fmt"
	"net/http"
	"web-service-gin/internal/app/chains"
	"web-service-gin/internal/pkg/pgpool"

	"github.com/ethereum/go-ethereum/common"
	"github.com/gin-gonic/gin"
)

type DetailRequest struct {
	ChainId string `json:"chainId" example:"eth"`
	Token   string `json:"token" example:"0x429D83Bb0DCB8cdd5311e34680ADC8B12070a07f"`
}

type DetailResponse struct {
	Id       uint8  `json:"id" example:"1"`
	Symbol   string `json:"symbol" example:"PLTC"`
	Decimals uint8  `json:"decimals" example:"18"`
}

// @BasePath /api/v1

// @Security BearerAuth
// @Tags Token
// @Summary 获取详情
// @Schemes
// @Description 获取token详情
// @Accept json
// @Produce json
// @Param token query DetailRequest true "Token and chain"
// @Router /token/detail [get]
func DetialHandler(c *gin.Context) {

	chainId := c.DefaultQuery("chainId", "eth")
	token := c.Query("token")
	if chainId != "eth" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "only support eth chain"})
		return
	}

	ethAddress := common.HexToAddress(token)
	// Check if the address is the zero address (invalid address)
	if ethAddress == common.HexToAddress("0x0000000000000000000000000000000000000000") {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "token is invalid",
		})
		return
	}

	//从pg数据库查询
	db := pgpool.GetDBConn()
	defer db.Close()

	var db_data DetailResponse
	err := db.QueryRow("SELECT id,symbol, decimals FROM public.contact_token WHERE chainId = $1 AND token = $2;", chainId, token).Scan(&db_data.Id, &db_data.Symbol, &db_data.Decimals)
	if err == sql.ErrNoRows {
		fmt.Println("未找到匹配的记录")
	} else if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Error when query by sql",
		})
		return
	} else {
		c.JSON(http.StatusOK, db_data)
		return
	}

	//从合约中获取
	info, err := chains.EthContractInfo(token)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to fetch token information.",
		})
		return
	}

	//插入到数据库
	err = db.QueryRow("INSERT INTO public.contact_token (chainId, token, symbol, decimals) VALUES ($1, $2, $3, $4) RETURNING id",
		chainId, token, info.Symbol, info.Decimals).Scan(&db_data.Id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to fetch insert information.",
		})
		return
	}
	db_data.Decimals = info.Decimals
	db_data.Symbol = info.Symbol
	c.JSON(http.StatusOK, db_data)
}

type AddRequest struct {
	ChainId  string `json:"chainId" example:"eth"`
	Token    string `json:"token" example:"0x429D83Bb0DCB8cdd5311e34680ADC8B12070a07f"`
	Symbol   string `json:"symbol" example:"PLTC"`
	Decimals uint8  `json:"decimals" example:"18"`
}

// @Security BearerAuth
// @Tags Token
// @Summary 添加token信息
// @Schemes
// @Description 添加token信息,仅owener角色用户可用
// @Accept json
// @Produce json
// @Param token body AddRequest true "token info"
// @Router /token/add [post]
func AddHandler(c *gin.Context) {

	var jsonData AddRequest
	if err := c.ShouldBindJSON(&jsonData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON format"})
		return
	}

	//从pg数据库查询
	db := pgpool.GetDBConn()
	defer db.Close()

	//插入到数据库
	_, err := db.Exec("INSERT INTO public.contact_token (chainId, token, symbol, decimals) VALUES ($1, $2, $3, $4)",
		jsonData.ChainId, jsonData.Token, jsonData.Symbol, jsonData.Decimals)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}
	c.String(http.StatusOK, "success")
}

type DelRequest struct {
	Id string `json:"id" example:"1"`
}

// @Security BearerAuth
// @Tags Token
// @Summary 删除token信息
// @Schemes
// @Description 删除token信息,仅owener角色用户可用
// @Accept json
// @Produce json
// @Param token body DelRequest true "id"
// @Router /token/delete [post]
func DeleteHandler(c *gin.Context) {

	var jsonData DelRequest
	if err := c.ShouldBindJSON(&jsonData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON format"})
		return
	}

	db := pgpool.GetDBConn()
	defer db.Close()

	//删除
	_, err := db.Exec("DELETE FROM public.contact_token WHERE id = $1", jsonData.Id)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.String(http.StatusOK, "success")
}

type EditRequest struct {
	Id       string `json:"id" example:"1"`
	ChainId  string `json:"chainId" example:"eth"`
	Token    string `json:"token" example:"0x429D83Bb0DCB8cdd5311e34680ADC8B12070a07f"`
	Symbol   string `json:"symbol" example:"PLTC"`
	Decimals uint8  `json:"decimals" example:"18"`
}

// @Security BearerAuth
// @Tags Token
// @Summary 修改token信息，,仅owener角色用户可用
// @Schemes
// @Description 修改token信息
// @Accept json
// @Produce json
// @Param token body EditRequest true "id"
// @Router /token/edit [post]
func EditHandler(c *gin.Context) {

	var jsonData EditRequest
	if err := c.ShouldBindJSON(&jsonData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON format"})
		return
	}

	db := pgpool.GetDBConn()
	defer db.Close()
	//更新
	_, err := db.Exec("UPDATE public.contact_token SET chainId = $1, token = $2, symbol = $3, decimals=$4 WHERE id = $5", jsonData.ChainId, jsonData.Token, jsonData.Symbol, jsonData.Decimals, jsonData.Id)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.String(http.StatusOK, "success")
}
