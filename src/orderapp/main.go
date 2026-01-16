package main

import (
	"log"
	"net/http"
	"strconv"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
)

// 基本数据模型

type Category struct {
	ID   int64  `json:"id"`
	Name string `json:"name"`
}

type Dish struct {
	ID            int64   `json:"id"`
	CategoryID    int64   `json:"category_id"`
	Name          string  `json:"name"`
	Price         float64 `json:"price"`
	Image         string  `json:"image"`
	IsRecommended bool    `json:"is_recommended"`
	IsAvailable   bool    `json:"is_available"`
}

type Table struct {
	ID             int64  `json:"id"`
	Name           string `json:"name"`
	Status         string `json:"status"` // FREE / USING
	CurrentOrderID *int64 `json:"current_order_id,omitempty"`
	PeopleCount    int    `json:"people_count"`
}

type CartItem struct {
	ID       int64  `json:"id"`
	TableID  int64  `json:"table_id"`
	DishID   int64  `json:"dish_id"`
	Quantity int    `json:"quantity"`
	Remark   string `json:"remark"`
}

type OrderItem struct {
	ID       int64   `json:"id"`
	OrderID  int64   `json:"order_id"`
	DishID   int64   `json:"dish_id"`
	Name     string  `json:"name"`
	Price    float64 `json:"price"`
	Quantity int     `json:"quantity"`
	Remark   string  `json:"remark"`
}

type Order struct {
	ID          int64       `json:"id"`
	TableID     int64       `json:"table_id"`
	PeopleCount int         `json:"people_count"`
	Status      string      `json:"status"` // PLACED / COOKING / SERVED / PAID / CANCELED
	TotalAmount float64     `json:"total_amount"`
	CreatedAt   time.Time   `json:"created_at"`
	Items       []OrderItem `json:"items"`
}

const (
	TableStatusFree  = "FREE"
	TableStatusUsing = "USING"

	OrderStatusPlaced   = "PLACED"
	OrderStatusCooking  = "COOKING"
	OrderStatusServed   = "SERVED"
	OrderStatusPaid     = "PAID"
	OrderStatusCanceled = "CANCELED"
)

// 内存存储（简单示例，不考虑持久化和多进程共享）

var (
	mu         sync.Mutex
	categories []Category
	dishes     []Dish
	tables     []Table
	cartItems  []CartItem
	orders     []Order

	nextCategoryID  int64 = 1
	nextDishID      int64 = 1
	nextTableID     int64 = 1
	nextCartItemID  int64 = 1
	nextOrderID     int64 = 1
	nextOrderItemID int64 = 1
)

func initData() {
	categories = []Category{
		{ID: nextCategoryID, Name: "热菜"},
		{ID: nextCategoryID + 1, Name: "凉菜"},
		{ID: nextCategoryID + 2, Name: "饮料"},
	}
	nextCategoryID += int64(len(categories))

	dishes = []Dish{
		{ID: nextDishID, CategoryID: 1, Name: "宫保鸡丁", Price: 38, Image: "/images/gongbao.jpg", IsRecommended: true, IsAvailable: true},
		{ID: nextDishID + 1, CategoryID: 1, Name: "鱼香肉丝", Price: 32, Image: "/images/yuxiang.jpg", IsRecommended: false, IsAvailable: true},
		{ID: nextDishID + 2, CategoryID: 2, Name: "拍黄瓜", Price: 16, Image: "/images/paihuanggua.jpg", IsRecommended: false, IsAvailable: true},
		{ID: nextDishID + 3, CategoryID: 3, Name: "可乐", Price: 8, Image: "/images/cola.jpg", IsRecommended: false, IsAvailable: true},
	}
	nextDishID += int64(len(dishes))

	tables = []Table{
		{ID: nextTableID, Name: "A1", Status: TableStatusFree, PeopleCount: 0},
		{ID: nextTableID + 1, Name: "A2", Status: TableStatusFree, PeopleCount: 0},
		{ID: nextTableID + 2, Name: "B1", Status: TableStatusFree, PeopleCount: 0},
	}
	nextTableID += int64(len(tables))
}

func main() {
	initData()

	r := gin.Default()

	r.Use(func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		c.Writer.Header().Set("Cache-Control", "no-cache, no-store, must-revalidate")
		c.Writer.Header().Set("Pragma", "no-cache")
		c.Writer.Header().Set("Expires", "0")
		if c.Request.Method == http.MethodOptions {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}
		c.Next()
	})

	api := r.Group("/api")
	{
		api.GET("/categories", getCategories)
		api.GET("/dishes", getDishes)

		api.GET("/tables", getTables)
		api.POST("/tables/:id/open", openTable)
		api.POST("/tables/:id/close", closeTable)

		api.GET("/cart", getCart)
		api.POST("/cart/items", addCartItem)
		api.PUT("/cart/items/:id", updateCartItem)
		api.DELETE("/cart/items/:id", deleteCartItem)

		api.POST("/orders", createOrder)
		api.GET("/orders/:id", getOrder)
		api.GET("/tables/:id/current-order", getCurrentOrderForTable)
		api.POST("/orders/:id/status", updateOrderStatus)
		api.GET("/orders/:id/bill", getOrderBill)
		api.POST("/orders/:id/pay", payOrder)
	}

	if err := r.Run(":8090"); err != nil {
		log.Fatal(err)
	}
}

// Handler 实现

func getCategories(c *gin.Context) {
	mu.Lock()
	defer mu.Unlock()

	log.Printf("[getCategories] returning %d categories", len(categories))
	c.JSON(http.StatusOK, categories)
}

func getDishes(c *gin.Context) {
	mu.Lock()
	defer mu.Unlock()

	categoryIDStr := c.Query("category_id")
	if categoryIDStr == "" {
		c.JSON(http.StatusOK, dishes)
		return
	}

	categoryID, err := strconv.ParseInt(categoryIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid category_id"})
		return
	}

	var result []Dish
	for _, d := range dishes {
		if d.CategoryID == categoryID {
			result = append(result, d)
		}
	}
	c.JSON(http.StatusOK, result)
}

func getTables(c *gin.Context) {
	mu.Lock()
	defer mu.Unlock()

	c.JSON(http.StatusOK, tables)
}

func openTable(c *gin.Context) {
	idStr := c.Param("id")
	tableID, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid table id"})
		return
	}

	var req struct {
		PeopleCount int `json:"people_count"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid body"})
		return
	}

	mu.Lock()
	defer mu.Unlock()

	for i := range tables {
		if tables[i].ID == tableID {
			if tables[i].Status == TableStatusUsing {
				c.JSON(http.StatusBadRequest, gin.H{"error": "table already in use"})
				return
			}
			tables[i].Status = TableStatusUsing
			tables[i].PeopleCount = req.PeopleCount
			c.JSON(http.StatusOK, tables[i])
			return
		}
	}

	c.JSON(http.StatusNotFound, gin.H{"error": "table not found"})
}

func closeTable(c *gin.Context) {
	idStr := c.Param("id")
	tableID, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid table id"})
		return
	}

	mu.Lock()
	defer mu.Unlock()

	for i := range tables {
		if tables[i].ID == tableID {
			tables[i].Status = TableStatusFree
			tables[i].PeopleCount = 0
			tables[i].CurrentOrderID = nil
			// 清空该桌购物车
			var newCart []CartItem
			for _, item := range cartItems {
				if item.TableID != tableID {
					newCart = append(newCart, item)
				}
			}
			cartItems = newCart

			c.JSON(http.StatusOK, tables[i])
			return
		}
	}

	c.JSON(http.StatusNotFound, gin.H{"error": "table not found"})
}

// 购物车

type CartResponseItem struct {
	ID       int64   `json:"id"`
	DishID   int64   `json:"dish_id"`
	DishName string  `json:"dish_name"`
	Price    float64 `json:"price"`
	Quantity int     `json:"quantity"`
	Remark   string  `json:"remark"`
	Amount   float64 `json:"amount"`
}

type CartResponse struct {
	TableID int64              `json:"table_id"`
	Items   []CartResponseItem `json:"items"`
	Total   float64            `json:"total"`
}

func getCart(c *gin.Context) {
	tableIDStr := c.Query("table_id")
	if tableIDStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "table_id is required"})
		return
	}
	tableID, err := strconv.ParseInt(tableIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid table_id"})
		return
	}

	mu.Lock()
	defer mu.Unlock()

	var resp CartResponse
	resp.TableID = tableID

	for _, item := range cartItems {
		if item.TableID != tableID {
			continue
		}
		// 找菜品信息
		var dishName string
		var price float64
		for _, d := range dishes {
			if d.ID == item.DishID {
				dishName = d.Name
				price = d.Price
				break
			}
		}
		amount := price * float64(item.Quantity)
		resp.Items = append(resp.Items, CartResponseItem{
			ID:       item.ID,
			DishID:   item.DishID,
			DishName: dishName,
			Price:    price,
			Quantity: item.Quantity,
			Remark:   item.Remark,
			Amount:   amount,
		})
		resp.Total += amount
	}

	c.JSON(http.StatusOK, resp)
}

func addCartItem(c *gin.Context) {
	var req struct {
		TableID  int64  `json:"table_id"`
		DishID   int64  `json:"dish_id"`
		Quantity int    `json:"quantity"`
		Remark   string `json:"remark"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid body"})
		return
	}
	if req.TableID == 0 || req.DishID == 0 || req.Quantity <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid data"})
		return
	}

	mu.Lock()
	defer mu.Unlock()

	// 检查桌台是否存在
	if !tableExists(req.TableID) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "table not found"})
		return
	}

	item := CartItem{
		ID:       nextCartItemID,
		TableID:  req.TableID,
		DishID:   req.DishID,
		Quantity: req.Quantity,
		Remark:   req.Remark,
	}
	nextCartItemID++
	cartItems = append(cartItems, item)

	c.JSON(http.StatusOK, item)
}

func updateCartItem(c *gin.Context) {
	idStr := c.Param("id")
	itemID, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	var req struct {
		Quantity *int   `json:"quantity"`
		Remark   string `json:"remark"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid body"})
		return
	}

	mu.Lock()
	defer mu.Unlock()

	for i := range cartItems {
		if cartItems[i].ID == itemID {
			if req.Quantity != nil {
				if *req.Quantity <= 0 {
					c.JSON(http.StatusBadRequest, gin.H{"error": "quantity must be > 0"})
					return
				}
				cartItems[i].Quantity = *req.Quantity
			}
			if req.Remark != "" {
				cartItems[i].Remark = req.Remark
			}
			c.JSON(http.StatusOK, cartItems[i])
			return
		}
	}

	c.JSON(http.StatusNotFound, gin.H{"error": "cart item not found"})
}

func deleteCartItem(c *gin.Context) {
	idStr := c.Param("id")
	itemID, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	mu.Lock()
	defer mu.Unlock()

	var newCart []CartItem
	var deleted bool
	for _, item := range cartItems {
		if item.ID == itemID {
			deleted = true
			continue
		}
		newCart = append(newCart, item)
	}
	cartItems = newCart

	if !deleted {
		c.JSON(http.StatusNotFound, gin.H{"error": "cart item not found"})
		return
	}

	c.Status(http.StatusNoContent)
}

// 订单

func createOrder(c *gin.Context) {
	var req struct {
		TableID      int64  `json:"table_id"`
		PeopleCount  int    `json:"people_count"`
		CustomerName string `json:"customer_name"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid body"})
		return
	}
	if req.TableID == 0 || req.PeopleCount <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid data"})
		return
	}

	mu.Lock()
	defer mu.Unlock()

	// 收集该桌的购物车
	var tableCart []CartItem
	for _, item := range cartItems {
		if item.TableID == req.TableID {
			tableCart = append(tableCart, item)
		}
	}
	if len(tableCart) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "cart is empty"})
		return
	}

	order := Order{
		ID:          nextOrderID,
		TableID:     req.TableID,
		PeopleCount: req.PeopleCount,
		Status:      OrderStatusPlaced,
		CreatedAt:   time.Now(),
	}
	nextOrderID++

	// 生成订单明细
	for _, item := range tableCart {
		var dishName string
		var price float64
		for _, d := range dishes {
			if d.ID == item.DishID {
				dishName = d.Name
				price = d.Price
				break
			}
		}
		amount := price * float64(item.Quantity)
		orderItem := OrderItem{
			ID:       nextOrderItemID,
			OrderID:  order.ID,
			DishID:   item.DishID,
			Name:     dishName,
			Price:    price,
			Quantity: item.Quantity,
			Remark:   item.Remark,
		}
		nextOrderItemID++
		order.Items = append(order.Items, orderItem)
		order.TotalAmount += amount
	}

	orders = append(orders, order)

	// 清空该桌购物车
	var newCart []CartItem
	for _, item := range cartItems {
		if item.TableID != req.TableID {
			newCart = append(newCart, item)
		}
	}
	cartItems = newCart

	// 更新桌台当前订单
	for i := range tables {
		if tables[i].ID == req.TableID {
			tables[i].Status = TableStatusUsing
			tables[i].PeopleCount = req.PeopleCount
			orderID := order.ID
			tables[i].CurrentOrderID = &orderID
			break
		}
	}

	c.JSON(http.StatusOK, order)
}

func getOrder(c *gin.Context) {
	idStr := c.Param("id")
	orderID, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	mu.Lock()
	defer mu.Unlock()

	for _, o := range orders {
		if o.ID == orderID {
			c.JSON(http.StatusOK, o)
			return
		}
	}

	c.JSON(http.StatusNotFound, gin.H{"error": "order not found"})
}

func getCurrentOrderForTable(c *gin.Context) {
	idStr := c.Param("id")
	tableID, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid table id"})
		return
	}

	mu.Lock()
	defer mu.Unlock()

	var currentOrderID *int64
	for _, t := range tables {
		if t.ID == tableID {
			currentOrderID = t.CurrentOrderID
			break
		}
	}
	if currentOrderID == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "no current order for this table"})
		return
	}

	for _, o := range orders {
		if o.ID == *currentOrderID {
			c.JSON(http.StatusOK, o)
			return
		}
	}

	c.JSON(http.StatusNotFound, gin.H{"error": "order not found"})
}

func updateOrderStatus(c *gin.Context) {
	idStr := c.Param("id")
	orderID, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	var req struct {
		Status string `json:"status"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid body"})
		return
	}

	mu.Lock()
	defer mu.Unlock()

	for i := range orders {
		if orders[i].ID == orderID {
			orders[i].Status = req.Status
			c.JSON(http.StatusOK, orders[i])
			return
		}
	}

	c.JSON(http.StatusNotFound, gin.H{"error": "order not found"})
}

func getOrderBill(c *gin.Context) {
	idStr := c.Param("id")
	orderID, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	mu.Lock()
	defer mu.Unlock()

	for _, o := range orders {
		if o.ID == orderID {
			c.JSON(http.StatusOK, gin.H{
				"order_id":     o.ID,
				"table_id":     o.TableID,
				"total_amount": o.TotalAmount,
				"status":       o.Status,
			})
			return
		}
	}

	c.JSON(http.StatusNotFound, gin.H{"error": "order not found"})
}

func payOrder(c *gin.Context) {
	idStr := c.Param("id")
	orderID, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	var req struct {
		PayMethod     string `json:"pay_method"`
		TransactionID string `json:"transaction_id"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		// 允许空 body，使用默认方式
	}

	mu.Lock()
	defer mu.Unlock()

	for i := range orders {
		if orders[i].ID == orderID {
			orders[i].Status = OrderStatusPaid

			// 支付完成后，清理桌台当前订单（但不强制清台，交给 closeTable 再清）
			for j := range tables {
				if tables[j].ID == orders[i].TableID {
					tables[j].CurrentOrderID = nil
					break
				}
			}

			c.JSON(http.StatusOK, gin.H{
				"order":          orders[i],
				"pay_method":     req.PayMethod,
				"transaction_id": req.TransactionID,
			})
			return
		}
	}

	c.JSON(http.StatusNotFound, gin.H{"error": "order not found"})
}

// 辅助函数

func tableExists(id int64) bool {
	for _, t := range tables {
		if t.ID == id {
			return true
		}
	}
	return false
}
