package server

import (
	"CryptoParser/config"
	"CryptoParser/internal/delivery/http/middleware"
	"context"
	"fmt"
	"net/http"

	// "CryptoParser/internal/entities"
	"CryptoParser/internal/usecase"
	protos "CryptoParser/pkg/proto/gen/go"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

const (
	SuccessStatus = `Success`
)

type Server struct {
	logger  *zap.Logger
	cfg     *config.ConfigModel
	serv    *gin.Engine
	Usecase *usecase.Usecase
	mdlware *middleware.Middleware
}

func NewServer(logger *zap.Logger, cfg *config.ConfigModel, uc *usecase.Usecase, mdlware *middleware.Middleware) (*Server, error) {
	return &Server{
		logger:  logger,
		cfg:     cfg,
		serv:    gin.Default(),
		Usecase: uc,
		mdlware: mdlware,
	}, nil
}

func (s *Server) OnStart(_ context.Context) error {
	s.createController()
	go func() {
		s.logger.Debug("serv started")
		if err := s.serv.Run(s.cfg.Server.Host + ":" + s.cfg.Server.Port); err != nil {
			s.logger.Error("failed to serve: " + err.Error())
		}
		return
	}()
	return nil
}

func (s *Server) OnStop(_ context.Context) error {
	s.logger.Debug("stop grps")
	//s.serv.GracefulStop()
	return nil
}

// @Summary Add Currency
// @Description add currency pair
// @ID add-currency
// @Accept  json
// @Produce  json
// @Param CurrencyAddRequest body auth_protov1.CurrencyAddRequest true "CurrencyAddRequest: currencypair"
// @Success 200 {object} auth_protov1.CurrencyAddResponse "Success"
// @Failure 400 {object} errorResponse
// @Failure 500 {object} errorResponse
// @Failure default {object} errorResponse
// @Router /currency/add [post]
func (s *Server) CurrencyAdd(ctx *gin.Context) {
	request := protos.CurrencyAddRequest{}

	if err := ctx.ShouldBindJSON(&request); err != nil {
		resp := errorResponse{Message: fmt.Sprintf("failed to unmarshar request: %v", err)}
		ctx.AbortWithStatusJSON(http.StatusBadRequest, resp)
		return
	}

	// TODO: Бизнес логика
	if err := s.Usecase.CurrencyAdd(ctx, request.CurrencyPair); err != nil {
		resp := errorResponse{Message: fmt.Sprintf("failed to add currency pair: %v", err)}
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, resp)
		return
	}

	ctx.JSON(http.StatusOK, &protos.CurrencyAddResponse{Status: SuccessStatus})
	return
}

// @Summary Remove Currency
// @Description remove currency pair
// @ID remove-currency
// @Accept  json
// @Produce  json
// @Param CurrencyRemoveRequest body auth_protov1.CurrencyRemoveRequest true "CurrencyRemoveRequest: currencypair"
// @Success 200 {object} auth_protov1.CurrencyRemoveResponse "Success"
// @Failure 400 {object} errorResponse
// @Failure 500 {object} errorResponse
// @Failure default {object} errorResponse
// @Router /currency/remove [post]
func (s *Server) CurrencyRemove(ctx *gin.Context) {
	request := protos.CurrencyRemoveRequest{}

	if err := ctx.ShouldBindJSON(&request); err != nil {
		resp := errorResponse{Message: fmt.Sprintf("failed to unmarshar request: %v", err)}
		ctx.AbortWithStatusJSON(http.StatusBadRequest, resp)
		return
	}

	// TODO: Бизнес логика
	if err := s.Usecase.CurrencyRemove(ctx, request.CurrencyPair); err != nil {
		resp := errorResponse{Message: fmt.Sprintf("failed to remove currency pair: %v", err)}
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, resp)
		return
	}

	ctx.JSON(http.StatusOK, &protos.CurrencyRemoveResponse{Status: SuccessStatus})
	return
}

// @Summary Currency Price
// @Description get currency price at specific time
// @ID get-currency-price
// @Accept  json
// @Produce  json
// @Param CurrencyPriceRequest body auth_protov1.CurrencyPriceRequest true "CurrencyPriceRequest: currencyPair, timestamp"
// @Success 200 {object} auth_protov1.CurrencyPriceResponse "Success"
// @Failure 400 {object} errorResponse
// @Failure 500 {object} errorResponse
// @Failure default {object} errorResponse
// @Router /currency/price [post]
func (s *Server) CurrencyPrice(ctx *gin.Context) {
	request := protos.CurrencyPriceRequest{}

	if err := ctx.ShouldBindJSON(&request); err != nil {
		resp := errorResponse{Message: fmt.Sprintf("failed to unmarshar request: %v", err)}
		ctx.AbortWithStatusJSON(http.StatusBadRequest, resp)
		return
	}

	// TODO: Бизнес логика
	price, err := s.Usecase.CurrencyPrice(ctx, request.CurrencyPair, request.Timestamp)
	if err != nil {
		resp := errorResponse{Message: fmt.Sprintf("failed to get price by currency pair and timestamp: %v", err)}
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, resp)
		return
	}

	ctx.JSON(http.StatusOK, &protos.CurrencyPriceResponse{Price: price})
	return
}
