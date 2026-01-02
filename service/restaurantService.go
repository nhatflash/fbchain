package service

import (
	"context"
	"os"
	"strconv"
	"time"

	"github.com/ncruces/zenity"
	"github.com/nhatflash/fbchain/client"
	"github.com/nhatflash/fbchain/constant"
	"github.com/nhatflash/fbchain/enum"
	appErr "github.com/nhatflash/fbchain/error"
	"github.com/nhatflash/fbchain/helper"
	"github.com/nhatflash/fbchain/model"
	"github.com/nhatflash/fbchain/repository"
	"github.com/redis/go-redis/v9"
	"github.com/shopspring/decimal"
	"github.com/skip2/go-qrcode"
	"go.mongodb.org/mongo-driver/v2/bson"
)

type IRestaurantService interface {
	HandleCreateRestaurant(ctx context.Context, req *client.CreateRestaurantRequest, tenantId int64) (*client.RestaurantResponse, error)
	FindRestaurantsByTenantId(ctx context.Context, tenantId int64) ([]model.Restaurant, error)
	FindAllRestaurants(ctx context.Context) ([]model.Restaurant, error)
	FindRestaurantById(ctx context.Context, id int64) (*model.Restaurant, error)
	FindRestaurantImageById(ctx context.Context, id int64) (*model.RestaurantImage, error)
	FindRestaurantImagesByRestaurantId(ctx context.Context, restaurantId int64) ([]model.RestaurantImage, error)
	FindAllRestaurantImages(ctx context.Context) ([]model.RestaurantImage, error)
	HandleAddNewRestaurantItem(ctx context.Context, restaurantId int64, tenantId int64, req *client.AddRestaurantItemRequest) (*client.RestaurantItemResponse, error)
	FindItemsByRestaurantId(ctx context.Context, restaurantId int64) ([]model.RestaurantItem, error)
	FindAllRestaurantItems(ctx context.Context) ([]model.RestaurantItem, error)
	FindRestaurantItemById(ctx context.Context, id bson.ObjectID) (*model.RestaurantItem, error)
	HandleAddNewRestaurantTable(ctx context.Context, tenantId int64, restaurantId int64, req *client.AddRestaurantTableRequest) (*client.RestaurantTableResponse, error)
	FindRestaurantTableById(ctx context.Context, id int64) (*model.RestaurantTable, error)
	FindRestaurantTablesByRestaurantId(ctx context.Context, restaurantId int64) ([]model.RestaurantTable, error)
	FindAllRestaurantTables(ctx context.Context) ([]model.RestaurantTable, error)
	GetQRCodeOnRestaurantTable(ctx context.Context, tableId int64, tenantId int64, restaurantId int64) error
	HandleShowRestaurantItemsViaQRCode(ctx context.Context, tableId int64) ([]model.RestaurantItem, error)
	HandleStartTableOrderingSession(ctx context.Context, tableId int64) error
	HandleEndTableOrderingSession(ctx context.Context, tableId int64) error
	HandleCreateRestaurantOrder(ctx context.Context, tableId int64, req *client.CreateRestaurantOrderRequest) (*client.RestaurantOrderResponse, error)
	HandlePayRestaurantOrderWithCash(ctx context.Context, orderId int64, tableId int64) error
	FindRestaurantOrderById(ctx context.Context, id int64) (*model.RestaurantOrder, error)
	FindAllRestaurantOrders(ctx context.Context) ([]model.RestaurantOrder, error)
	FindAllRestaurantOrderItems(ctx context.Context) ([]model.RestaurantOrderItem, error)
	FindRestaurantOrderItemById(ctx context.Context, id int64) (*model.RestaurantOrderItem, error)
	FindRestaurantOrderItemsByOrderId(ctx context.Context, orderId int64) ([]model.RestaurantOrderItem, error)
}

type RestaurantService struct {
	RestaurantRepo 			*repository.RestaurantRepository
	SubPackageRepo 			*repository.SubPackageRepository
	RestaurantItemRepo 		*repository.RestaurantItemRepository
	RestaurantTableRepo 	*repository.RestaurantTableRepository
	RestaurantOrderRepo 	*repository.RestaurantOrderRepository
	RestaurantPaymentRepo 	*repository.RestaurantPaymentRepository
	Rdb 					*redis.Client
}

func NewRestaurantService(rr *repository.RestaurantRepository, 
						  spr *repository.SubPackageRepository, 
						  rir *repository.RestaurantItemRepository, 
						  rtr *repository.RestaurantTableRepository, 
						  ror *repository.RestaurantOrderRepository,
						  rpr *repository.RestaurantPaymentRepository,
						  rdb *redis.Client) IRestaurantService {
	return &RestaurantService{
		RestaurantRepo: rr,
		SubPackageRepo: spr,
		RestaurantItemRepo: rir,
		RestaurantTableRepo: rtr,
		RestaurantOrderRepo: ror,
		RestaurantPaymentRepo: rpr,
		Rdb: rdb,
	}
}

func (rs *RestaurantService) HandleCreateRestaurant(ctx context.Context, req *client.CreateRestaurantRequest, tenantId int64) (*client.RestaurantResponse, error) {
	var err error
	
	if err = validateCreateRestaurantRequest(ctx, req.Name, rs.SubPackageRepo, rs.RestaurantRepo); err != nil {
		return nil, err
	}
	var s *model.SubPackage
	s, err = rs.SubPackageRepo.FindFirstSubPackage(ctx)
	if err != nil {
		return nil, err
	}
	var r *model.Restaurant
	r, err = rs.RestaurantRepo.CreateNewRestaurant(ctx, req.Name, req.Location, req.Description, req.ContactEmail, req.ContactPhone, req.PostalCode, *req.Type, req.Notes, s.Id, req.Images, tenantId)
	if err != nil {
		return nil, err
	}
	var rImgs []model.RestaurantImage
	rImgs, err = rs.RestaurantRepo.FindRestaurantImagesByRestaurantId(ctx, r.Id)
	if err != nil {
		return nil, err
	}
	return helper.MapToRestaurantResponse(r, rImgs), nil
}



func (rs *RestaurantService) FindRestaurantsByTenantId(ctx context.Context, tenantId int64) ([]model.Restaurant, error) {
	
	r, err := rs.RestaurantRepo.FindRestaurantsByTenantId(ctx, tenantId)
	if err != nil {
		return nil, err
	}
	return r, nil
}



func (rs *RestaurantService) FindAllRestaurants(ctx context.Context) ([]model.Restaurant, error) {
	r, err := rs.RestaurantRepo.FindAllRestaurants(ctx)
	if err != nil {
		return nil, err
	}
	return r, nil
}




func (rs *RestaurantService) FindRestaurantById(ctx context.Context, id int64) (*model.Restaurant, error) {
	r, err := rs.RestaurantRepo.FindRestaurantById(ctx, id)
	if err != nil {
		return nil, err
	}
	return r, nil
}




func (rs *RestaurantService) FindRestaurantImageById(ctx context.Context, id int64) (*model.RestaurantImage, error) {
	img, err := rs.RestaurantRepo.FindRestaurantImageById(ctx, id)
	if err != nil {
		return nil, err
	}
	return img, nil
}



func (rs *RestaurantService) FindRestaurantImagesByRestaurantId(ctx context.Context, restaurantId int64) ([]model.RestaurantImage, error) {
	imgs, err := rs.RestaurantRepo.FindRestaurantImagesByRestaurantId(ctx, restaurantId)
	if err != nil {
		return nil, err
	}
	return imgs, nil
}


func (rs *RestaurantService) FindAllRestaurantImages(ctx context.Context) ([]model.RestaurantImage, error) {
	imgs, err := rs.RestaurantRepo.FindAllRestaurantImages(ctx)
	if err != nil {
		return nil, err
	}
	return imgs, nil
}



func (rs *RestaurantService) HandleAddNewRestaurantItem(ctx context.Context, restaurantId int64, tenantId int64, req *client.AddRestaurantItemRequest) (*client.RestaurantItemResponse, error) {
	var err error 
	var r *model.Restaurant
	
	r, err = rs.RestaurantRepo.FindRestaurantById(ctx, restaurantId)
	if err != nil {
		return nil, err
	}

	if r.TenantId != tenantId {
		return nil, appErr.UnauthorizedError("You are not allowed to add new item on this restaurant.")
	}

	var price bson.Decimal128
	price, err = bson.ParseDecimal128(req.Price)
	if err != nil {
		return nil, appErr.BadRequestError("Invalid price.")
	}

	item := &model.RestaurantItem{
		Name: req.Name,
		Description: req.Description,
		Price: price,
		Type: req.Type,
		Status: enum.ITEM_AVAILABLE,
		Image: req.Image,
		Notes: req.Notes,
		RestaurantId: restaurantId,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	var newItem *model.RestaurantItem
	newItem, err = rs.RestaurantItemRepo.AddNewRestaurantItem(ctx, item)
	if err != nil {
		return nil, err
	}

	return helper.MapToRestaurantItemResponse(newItem), nil
}


func (rs *RestaurantService) FindItemsByRestaurantId(ctx context.Context, restaurantId int64) ([]model.RestaurantItem, error) {
	items, err := rs.RestaurantItemRepo.FindRestaurantItemsByRestaurantId(ctx, restaurantId)
	if err != nil {
		return nil, err
	}
	return items, nil
}


func (rs *RestaurantService) FindAllRestaurantItems(ctx context.Context) ([]model.RestaurantItem, error) {
	items, err := rs.RestaurantItemRepo.FindAllRestaurantItems(ctx)
	if err != nil {
		return nil, err
	}
	return items, nil
}


func (rs *RestaurantService) FindRestaurantItemById(ctx context.Context, id bson.ObjectID) (*model.RestaurantItem, error) {
	item, err := rs.RestaurantItemRepo.FindRestaurantItemById(ctx, id)
	if err != nil {
		return nil, err
	}
	return item, nil
}




func (rs *RestaurantService) HandleAddNewRestaurantTable(ctx context.Context, tenantId int64, restaurantId int64, req *client.AddRestaurantTableRequest) (*client.RestaurantTableResponse, error) {
	restaurant, err := rs.RestaurantRepo.FindRestaurantById(ctx, restaurantId)
	if err != nil {
		return nil, err
	}
	if restaurant.TenantId != tenantId {
		return nil, appErr.ForbiddenError("You are not allow to perform this action.")
	}
	var label string
	if req.Label == nil || *req.Label == "" {
		count, err := rs.RestaurantTableRepo.CountRestaurantTableByRestaurantId(ctx, restaurantId)
		if err != nil {
			return nil, err
		}
		count++
		label = strconv.Itoa(count)
	} else {
		label = *req.Label
	}
	table, err := rs.RestaurantTableRepo.AddNewRestaurantTable(ctx, restaurantId, label, req.Notes)
	if err != nil {
		return nil, err
	}
	return helper.MapToRestaurantTableResponse(table), nil
}


func (rs *RestaurantService) FindRestaurantTableById(ctx context.Context, id int64) (*model.RestaurantTable, error) {
	t, err := rs.RestaurantTableRepo.FindRestaurantTableById(ctx, id)
	if err != nil {
		return nil, err
	}
	return t, nil
}


func (rs *RestaurantService) FindRestaurantTablesByRestaurantId(ctx context.Context, restaurantId int64) ([]model.RestaurantTable, error) {
	tables, err := rs.RestaurantTableRepo.FindRestaurantTablesByRestaurantId(ctx, restaurantId)
	if err != nil {
		return nil, err
	}
	return tables, nil
}


func (rs *RestaurantService) FindAllRestaurantTables(ctx context.Context) ([]model.RestaurantTable, error) {
	tables, err := rs.RestaurantTableRepo.FindAllRestaurantTables(ctx)
	if err != nil {
		return nil, err
	}
	return tables, nil
}


func (rs *RestaurantService) GetQRCodeOnRestaurantTable(ctx context.Context, tableId int64, tenantId int64, restaurantId int64) error {
	var err error
	var table *model.RestaurantTable
	table, err = rs.FindRestaurantTableById(ctx, tableId)
	if err != nil {
		return err
	}
	
	var restaurant *model.Restaurant
	restaurant, err = rs.FindRestaurantById(ctx, restaurantId)
	if err != nil {
		return err
	}

	if table.RestaurantId != restaurantId || restaurant.TenantId != tenantId {
		return appErr.ForbiddenError("You are not allowed to perform action on table or restaurant that does not belong to you.")
	}

	baseUrl := os.Getenv("BASE_URL")
	tblIdStr := strconv.FormatInt(tableId, 10)
	url := baseUrl + "/" + tblIdStr
	fileName := "qrcode-" + tblIdStr + ".png"
	var path string
	path, err = zenity.SelectFileSave(
		zenity.Title("Save Table QR Code"),
		zenity.Filename(fileName),
		zenity.FileFilters{
			{Name: "PNG Images", Patterns: []string{"*.png"}},
		},
	)
	if err != nil {
		if err == zenity.ErrCanceled {
			return appErr.BadRequestError("QR code PNG image saving has been canceled.")
		}
		return err
	}
	err = qrcode.WriteFile(url, qrcode.Medium, 256, path)
	if err != nil {
		return err
	}
	return nil
}


// In development stage, location handling on the client devices is not viable
func (rs *RestaurantService) HandleShowRestaurantItemsViaQRCode(ctx context.Context, tableId int64) ([]model.RestaurantItem, error) {
	var err error
	var table *model.RestaurantTable
	table, err = rs.FindRestaurantTableById(ctx, tableId)
	if err != nil {
		return nil, err
	}
	var items []model.RestaurantItem
	items, err = rs.FindItemsByRestaurantId(ctx, table.RestaurantId)
	if err != nil {
		return nil, err
	}
	return items, nil
}


func (rs *RestaurantService) HandleStartTableOrderingSession(ctx context.Context, tableId int64) error {
	var err error
	_, err = rs.FindRestaurantTableById(ctx, tableId)
	if err != nil {
		return err
	}
	tableIdStr := strconv.FormatInt(tableId, 10)
	sessionKey := constant.RestaurantOrderSessionKey + tableIdStr
	duration := time.Duration(30) * time.Minute
	
	if err = rs.Rdb.Set(ctx, sessionKey, "true", duration).Err(); err != nil {
		return err
	}
	return nil
}


func (rs *RestaurantService) HandleEndTableOrderingSession(ctx context.Context, tableId int64) error {
	tableIdStr := strconv.FormatInt(tableId, 10)
	sessionKey := constant.RestaurantOrderSessionKey + tableIdStr
	var err error
	var exists int64
	exists, err = rs.Rdb.Exists(ctx, sessionKey).Result()
	if err != nil {
		return err
	}
	if exists == 1 {
		rs.Rdb.Del(ctx, sessionKey)
		return nil
	}
	return appErr.NotFoundError("No session or session already closed for table: " + tableIdStr)
}



func (rs *RestaurantService) HandleCreateRestaurantOrder(ctx context.Context, tableId int64, req *client.CreateRestaurantOrderRequest) (*client.RestaurantOrderResponse, error) {
	var err error

	// checking the tableId still in the session
	tableIdStr := strconv.FormatInt(tableId, 10)
	sessionKey := constant.RestaurantOrderSessionKey + tableIdStr
	var exists int64
	exists, err = rs.Rdb.Exists(ctx, sessionKey).Result()
	if err != nil {
		return nil, err
	}
	if exists != 1 {
		return nil, appErr.BadRequestError("This table session for order item has been expired.")
	}

	var table *model.RestaurantTable
	table, err = rs.FindRestaurantTableById(ctx, tableId)
	if err != nil {
		return nil, err
	}

	// process the item requests
	quantities := make(map[string]int)
	var reqIds []string
	for _, rq := range req.Items {
		reqIds = append(reqIds, rq.ItemId)
		quantities[rq.ItemId] = *rq.Quantity
	}

	// get the items with the restaurantId that associated with the table
	var items []model.RestaurantItem
	items, err = rs.RestaurantItemRepo.FindRestaurantItemsByListIds(ctx, table.RestaurantId, reqIds)
	if err != nil {
		return nil, err
	}

	if len(items) != len(reqIds) {
		return nil, appErr.BadRequestError("Some of your items in your order does not belong to this restaurant.")
	}

	// get the total amount
	var amount decimal.Decimal
	amount, err = decimal.NewFromString("0.00")
	if err != nil {
		return nil, err
	}
	var oItems []model.RestaurantOrderItem
	for _, item := range items {
		itemId := item.Id.Hex()
		q := quantities[itemId]
		quantity := decimal.NewFromInt(int64(q))
		priceStr := item.Price.String()

		var price decimal.Decimal
		price, err = decimal.NewFromString(priceStr)
		if err != nil {
			return nil, err
		}
		price = price.Mul(quantity)
		amount = amount.Add(price)
		oItem := model.RestaurantOrderItem{
			ItemId: itemId,
			Quantity: q,
			Total: price,
		}
		oItems = append(oItems, oItem)
	}

	// create order model for passing data into repository
	order := model.RestaurantOrder{
		RestaurantId: table.RestaurantId,
		TableId: tableId,
		Amount: amount,
		Notes: req.Notes,
	}

	var rOrder *model.RestaurantOrder
	rOrder, err = rs.RestaurantOrderRepo.CreateInitialRestaurantOrder(ctx, &order, oItems)
	if err != nil {
		return nil, err
	}

	return helper.MapToRestaurantOrderResponse(rOrder), nil
}


func (rs *RestaurantService) HandlePayRestaurantOrderWithCash(ctx context.Context, orderId int64, tableId int64) error {
	var err error
	var o *model.RestaurantOrder
	o, err = rs.FindRestaurantOrderById(ctx, orderId)
	if err != nil {
		return err
	}
	var t *model.RestaurantTable
	t, err = rs.FindRestaurantTableById(ctx, tableId)
	if err != nil {
		return err
	}

	if t.Id != o.TableId {
		return appErr.BadRequestError("The order does not belong with the requested table.")
	}
	expiration := time.Duration(-15)*time.Minute
	if o.CreatedAt.After(time.Now().Add(expiration)) || o.Status != enum.R_ORDER_PENDING {
		return appErr.BadRequestError("This order has been expired or has not been allowed to pay.")
	}
	
	if err = rs.RestaurantPaymentRepo.HandleCashPayment(ctx, orderId, o.Amount); err != nil {
		return err
	}
	tblIdStr := strconv.FormatInt(o.TableId, 10)
	sessionKey := constant.RestaurantOrderSessionKey + tblIdStr
	rs.Rdb.Del(ctx, sessionKey)
	return nil
}


func validateCreateRestaurantRequest(ctx context.Context, name string, subPackageRepo *repository.SubPackageRepository, resRepo *repository.RestaurantRepository) error {
	var err error
	var exist bool
	exist, err = subPackageRepo.AnySubPackageExists(ctx)
	if err != nil {
		return err
	}
	if !exist {
		return appErr.NotFoundError("There is no subscription available in the system.")
	}
	exist, err = resRepo.IsRestaurantNameExist(ctx, name)
	if err != nil {
		return err
	}
	if exist {
		return appErr.BadRequestError("Restaurant with this requested name is already exist.")
	}
	return nil
}



func (rs *RestaurantService) FindRestaurantOrderById(ctx context.Context, id int64) (*model.RestaurantOrder, error) {
	o, err := rs.RestaurantOrderRepo.FindRestaurantOrderById(ctx, id)
	if err != nil {
		return nil, err
	}
	return o, nil
}

func (rs *RestaurantService) FindAllRestaurantOrders(ctx context.Context) ([]model.RestaurantOrder, error) {
	o, err := rs.RestaurantOrderRepo.FindAllRestaurantOrders(ctx)
	if err != nil {
		return nil, err
	}
	return o, nil
}

func (rs *RestaurantService) FindAllRestaurantOrderItems(ctx context.Context) ([]model.RestaurantOrderItem, error) {
	i, err := rs.RestaurantOrderRepo.FindAllRestaurantOrderItems(ctx)
	if err != nil {
		return nil, err
	}
	return i, nil
}

func (rs *RestaurantService) FindRestaurantOrderItemById(ctx context.Context, id int64) (*model.RestaurantOrderItem, error) {
	i, err := rs.RestaurantOrderRepo.FindRestaurantOrderItemById(ctx, id)
	if err != nil {
		return nil, err
	}
	return i, nil
}

func (rs *RestaurantService) FindRestaurantOrderItemsByOrderId(ctx context.Context, orderId int64) ([]model.RestaurantOrderItem, error) {
	i, err := rs.RestaurantOrderRepo.FindRestaurantOrderItemsByOrderId(ctx, orderId)
	if err != nil {
		return nil, err
	}
	return i, nil
}