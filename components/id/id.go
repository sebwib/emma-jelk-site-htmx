package id

// Layout
const (
	ContentID          = "content"
	ModalContainerID   = "modal-container"
	BackgroundID       = "background"
	HeaderNavID        = "header-nav"
	HeaderNavMenuID    = "header-nav-menu"
	HeaderDesktopNavID = "header-desktop-nav"
)

// Sidebar
const (
	SidebarID           = "sidebar"
	FavouriteProductsID = "favourite-products"
	ChatTitlesID        = "chat-titles"
)

// Products
const (
	ProductsGridID       = "products-grid"
	CategoryFilterFormID = "category-filter-form"
	VisualizationID      = "visualization"
	AddCubeButtonID      = "add-cube"
	ReportProblemID      = "report-a-problem"
)

// Todos
const (
	TodoListID             = "todo-list"
	TodoDetailsID          = "todo-details"
	DetailInnerContainerID = "detail-inner-container"
)

const (
	ProductFormNameID        = "name"
	ProductFormDescriptionID = "description"
	ProductFormPriceID       = "price"
)

const (
	CartPage         = "cart-page"
	CartEmailInput   = "cart-email-input"
	CartSubmitButton = "cart-submit-button"
)

// Modal
const (
	ModalInnerID         = "modal-inner"
	ProblemDescriptionID = "problem-description"
	EditArtModalInner    = "edit-art-modal-inner"
	EditPrintModalInner  = "edit-print-modal-inner"
)

const (
	CategoryFormID           = "category-form"
	CategoryInputID          = "category-input"
	CategoriesAutocompleteID = "categories-autocomplete"
)

func Selector(id string) string {
	return "#" + id
}

// Helper functions for dynamic
func ProductID(id string) string {
	return "product-" + id
}

func TodoID(id string) string {
	return "todo-" + id
}

func CategoryCheckboxID(id string) string {
	return "category-" + id
}

func CategoryTagsID(name string) string {
	return name + "-tags"
}

func CartItemID(id string) string {
	return "cart-item-" + id
}

type CartSymbolMode string

const (
	CartSymbolModeDesktop CartSymbolMode = "desktop"
	CartSymbolModeMobile  CartSymbolMode = "mobile"
	CartSymbolModeCount   CartSymbolMode = "count"
)

func CartSymbolID(mode CartSymbolMode) string {
	if mode == CartSymbolModeMobile {
		return "cart-symbol-mobile"
	}
	if mode == CartSymbolModeCount {
		return "cart-symbol-count"
	}
	return "cart-symbol"
}
