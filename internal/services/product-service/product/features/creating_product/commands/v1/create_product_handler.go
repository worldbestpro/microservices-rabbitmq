package v1

import (
	"context"
	"encoding/json"
	"github.com/meysamhadeli/shop-golang-microservices/internal/pkg/grpc"
	"github.com/meysamhadeli/shop-golang-microservices/internal/pkg/logger"
	"github.com/meysamhadeli/shop-golang-microservices/internal/pkg/mapper"
	"github.com/meysamhadeli/shop-golang-microservices/internal/pkg/rabbitmq"
	"github.com/meysamhadeli/shop-golang-microservices/internal/services/product-service/config"
	"github.com/meysamhadeli/shop-golang-microservices/internal/services/product-service/product/contracts/data"
	"github.com/meysamhadeli/shop-golang-microservices/internal/services/product-service/product/features/creating_product/dtos/v1"
	v12 "github.com/meysamhadeli/shop-golang-microservices/internal/services/product-service/product/features/creating_product/events/v1"
	"github.com/meysamhadeli/shop-golang-microservices/internal/services/product-service/product/models"
)

type CreateProductHandler struct {
	log                logger.ILogger
	cfg                *config.Config
	repository         data.ProductRepository
	rabbitmqPublisher  rabbitmq.IPublisher
	IdentityGrpcClient grpc.GrpcClient
}

func NewCreateProductHandler(log logger.ILogger, cfg *config.Config, repository data.ProductRepository,
	rabbitmqPublisher rabbitmq.IPublisher, identityGrpcClient grpc.GrpcClient) *CreateProductHandler {
	return &CreateProductHandler{log: log, cfg: cfg, repository: repository, rabbitmqPublisher: rabbitmqPublisher, IdentityGrpcClient: identityGrpcClient}
}

func (c *CreateProductHandler) Handle(ctx context.Context, command *CreateProduct) (*v1.CreateProductResponseDto, error) {

	// simple call grpcClient
	//identityGrpcClient := identity_service.NewIdentityServiceClient(c.IdentityGrpcClient.GetGrpcConnection())
	//user, err := identityGrpcClient.GetUserById(ctx, &identity_service.GetUserByIdReq{UserId: "1"})
	//if err != nil {
	//	return nil, err
	//}

	//c.log.Infof("userId: %s", user.User.UserId)

	product := &models.Product{
		ProductId:   command.ProductID,
		Name:        command.Name,
		Description: command.Description,
		Price:       command.Price,
		CreatedAt:   command.CreatedAt,
	}

	createdProduct, err := c.repository.CreateProduct(ctx, product)
	if err != nil {
		return nil, err
	}

	evt, err := mapper.Map[*v12.ProductCreated](createdProduct)
	if err != nil {
		return nil, err
	}

	err = c.rabbitmqPublisher.PublishMessage(ctx, evt)
	if err != nil {
		return nil, err
	}

	response := &v1.CreateProductResponseDto{ProductId: product.ProductId}
	bytes, _ := json.Marshal(response)

	c.log.Info("CreateProductResponseDto", string(bytes))

	return response, nil
}
