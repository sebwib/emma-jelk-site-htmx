package db

import "time"

type OrderStatus string

const (
	OrderStatusPlaced    OrderStatus = "PLACED"
	OrderStatusContacted OrderStatus = "CONTACTED"
	OrderStatusShipped   OrderStatus = "SHIPPED"
)

type OrderRow struct {
	UUID        string
	OrderID     string
	CreatedAt   string
	ContactedAt string
	SentAt      string
	Email       string
	PrintID     string
	Title       string
	Typ         string
	Quantity    int
	Price       float64
	Status      OrderStatus
	HasPaid     bool
}

type Order struct {
	BuyerEmail  string
	CreatedAt   string
	ContactedAt string
	SentAt      string
	OrderID     string
	Status      OrderStatus
	HasPaidAll  bool
	Rows        []OrderRow
	TotalPrice  float64
}

func (db *DB) createOrdersTable() error {
	_, err := db.Exec(`
	CREATE TABLE IF NOT EXISTS orders (
		uuid TEXT PRIMARY KEY,
		order_id TEXT,
		created_at TEXT,
		contacted_at TEXT,
		sent_at TEXT,
		email TEXT,
		print_id TEXT,
		title TEXT,
		typ TEXT,
		quantity INTEGER,
		price REAL,
		status TEXT,
		has_paid BOOLEAN
	);
	`)

	return err
}

func (db *DB) AddOrder(order OrderRow) error {
	_, err := db.Exec(`
	INSERT INTO orders (uuid, order_id, created_at, contacted_at, sent_at, email, print_id, title, typ, quantity, price, status, has_paid)
	VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?);
	`, order.UUID, order.OrderID, order.CreatedAt, order.ContactedAt, order.SentAt, order.Email, order.PrintID, order.Title, order.Typ, order.Quantity, order.Price, order.Status, order.HasPaid)
	return err
}

func buildOrderFromRows(orderRows []OrderRow) Order {
	if len(orderRows) == 0 {
		return Order{}
	}

	order := Order{
		BuyerEmail:  orderRows[0].Email,
		OrderID:     orderRows[0].OrderID,
		CreatedAt:   orderRows[0].CreatedAt,
		ContactedAt: orderRows[0].ContactedAt,
		SentAt:      orderRows[0].SentAt,
		Status:      orderRows[0].Status,

		Rows: orderRows,
	}

	hasPaidAll := true
	totalPrice := 0.0
	for _, row := range orderRows {
		totalPrice += float64(row.Quantity) * row.Price
		if !row.HasPaid {
			hasPaidAll = false
		}
	}
	order.HasPaidAll = hasPaidAll
	order.TotalPrice = totalPrice

	return order
}

func (db *DB) GetOrderByID(orderID string) (Order, error) {
	rows, err := db.Query(`SELECT uuid, order_id, created_at, contacted_at, sent_at, email, print_id, title, typ, quantity, price, status, has_paid FROM orders WHERE order_id = ?;`, orderID)
	if err != nil {
		return Order{}, err
	}
	defer rows.Close()

	var orderRows []OrderRow
	for rows.Next() {
		var orderRow OrderRow
		err := rows.Scan(&orderRow.UUID, &orderRow.OrderID, &orderRow.CreatedAt, &orderRow.ContactedAt, &orderRow.SentAt, &orderRow.Email, &orderRow.PrintID, &orderRow.Title, &orderRow.Typ, &orderRow.Quantity, &orderRow.Price, &orderRow.Status, &orderRow.HasPaid)
		if err != nil {
			return Order{}, err
		}
		orderRows = append(orderRows, orderRow)
	}

	if len(orderRows) == 0 {
		return Order{}, nil
	}

	return buildOrderFromRows(orderRows), nil
}

func (db *DB) GetAllOrders() ([]Order, error) {
	rows, err := db.Query(`SELECT uuid, order_id, created_at, contacted_at, sent_at, email, print_id, title, typ, quantity, price, status, has_paid FROM orders ORDER BY created_at DESC;`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var orders []OrderRow
	for rows.Next() {
		var order OrderRow
		err := rows.Scan(&order.UUID, &order.OrderID, &order.CreatedAt, &order.ContactedAt, &order.SentAt, &order.Email, &order.PrintID, &order.Title, &order.Typ, &order.Quantity, &order.Price, &order.Status, &order.HasPaid)
		if err != nil {
			return nil, err
		}
		orders = append(orders, order)
	}

	orderMap := make(map[string][]OrderRow)
	for _, row := range orders {
		orderMap[row.OrderID] = append(orderMap[row.OrderID], row)
	}

	var result []Order
	for _, rows := range orderMap {
		if len(rows) > 0 {
			result = append(result, buildOrderFromRows(rows))
		}
	}

	return result, nil
}

// ...existing code...

func (db *DB) UpdateOrderStatus(orderID string, status string) error {
	var err error

	currentStatus := ""
	err = db.QueryRow(`SELECT status FROM orders WHERE order_id = ? LIMIT 1;`, orderID).Scan(&currentStatus)
	if err != nil {
		return err
	}

	if currentStatus == status {
		// No change needed
		return nil
	}

	timestamp := time.Now().Format(time.RFC3339)
	if status == string(OrderStatusContacted) {
		_, err = db.Exec(`UPDATE orders SET status = ?, contacted_at = ? WHERE order_id = ?;`, status, timestamp, orderID)
	} else if status == string(OrderStatusShipped) {
		_, err = db.Exec(`UPDATE orders SET status = ?, sent_at = ? WHERE order_id = ?;`, status, timestamp, orderID)
	} else {
		_, err = db.Exec(`UPDATE orders SET status = ? WHERE order_id = ?;`, status, orderID)
	}

	return err
}

func (db *DB) UpdateOrderPaidStatus(uuid string, hasPaid bool) error {
	_, err := db.Exec(`UPDATE orders SET has_paid = ? WHERE uuid = ?;`, hasPaid, uuid)
	return err
}
