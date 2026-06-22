package seed

import (
	"fmt"
	"log"
	"math/rand"
	"os"
	"strings"
	"time"

	"pricecompare/internal/models"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

func Run(db *gorm.DB) error {
	var count int64
	db.Model(&models.Category{}).Count(&count)
	if count > 0 {
		log.Println("Database already seeded, skipping...")
		return nil
	}

	log.Println("Seeding database...")

	categories := []models.Category{
		{Name: "Пищевые продукты", Slug: "food", Description: "Продукты питания и бакалея", Image: CategoryImage("food")},
		{Name: "Хозяйственные товары", Slug: "household", Description: "Товары для дома и хозяйства", Image: CategoryImage("household")},
		{Name: "Строительные товары", Slug: "construction", Description: "Материалы и инструменты", Image: CategoryImage("construction")},
		{Name: "Бытовая химия", Slug: "chemistry", Description: "Средства для уборки и гигиены", Image: CategoryImage("chemistry")},
		{Name: "Напитки", Slug: "drinks", Description: "Безалкогольные и алкогольные напитки", Image: CategoryImage("drinks")},
		{Name: "Товары для дома", Slug: "home", Description: "Мебель, декор и аксессуары", Image: CategoryImage("home")},
		{Name: "Электроника", Slug: "electronics", Description: "Гаджеты и техника", Image: CategoryImage("electronics")},
		{Name: "Другое", Slug: "other", Description: "Прочие товары", Image: CategoryImage("other")},
	}
	if err := db.Create(&categories).Error; err != nil {
		return err
	}

	stores := []models.Store{
		{Name: "SuperMart", Logo: StoreLogo("SuperMart"), Rating: 4.5, Description: "Крупная сеть супермаркетов с широким ассортиментом", DeliveryCost: 199, MinOrder: 1000},
		{Name: "FreshFood", Logo: StoreLogo("FreshFood"), Rating: 4.7, Description: "Свежие продукты и фермерские товары", DeliveryCost: 149, MinOrder: 800},
		{Name: "MegaStore", Logo: StoreLogo("MegaStore"), Rating: 4.3, Description: "Гипермаркет с низкими ценами", DeliveryCost: 249, MinOrder: 1500},
		{Name: "BuildMax", Logo: StoreLogo("BuildMax"), Rating: 4.4, Description: "Строительные материалы и инструменты", DeliveryCost: 299, MinOrder: 2000},
		{Name: "HomeShop", Logo: StoreLogo("HomeShop"), Rating: 4.6, Description: "Товары для дома и интерьера", DeliveryCost: 179, MinOrder: 1200},
		{Name: "CityMarket", Logo: StoreLogo("CityMarket"), Rating: 4.2, Description: "Городской маркет с быстрой доставкой", DeliveryCost: 99, MinOrder: 500},
	}
	if err := db.Create(&stores).Error; err != nil {
		return err
	}

	adminPass := os.Getenv("DEMO_ADMIN_PASSWORD")
	if adminPass == "" {
		adminPass = "admin123"
	}
	userPass := os.Getenv("DEMO_USER_PASSWORD")
	if userPass == "" {
		userPass = "user123"
	}
	adminEmail := os.Getenv("DEMO_ADMIN_EMAIL")
	if adminEmail == "" {
		adminEmail = "admin@pricecompare.ru"
	}
	userEmail := os.Getenv("DEMO_USER_EMAIL")
	if userEmail == "" {
		userEmail = "user@pricecompare.ru"
	}

	adminHash, _ := bcrypt.GenerateFromPassword([]byte(adminPass), bcrypt.DefaultCost)
	userHash, _ := bcrypt.GenerateFromPassword([]byte(userPass), bcrypt.DefaultCost)

	users := []models.User{
		{Email: adminEmail, PasswordHash: string(adminHash), Name: "Администратор", Role: models.RoleAdmin},
		{Email: userEmail, PasswordHash: string(userHash), Name: "Демо Пользователь", Role: models.RoleUser},
	}
	if err := db.Create(&users).Error; err != nil {
		return err
	}

	productData := generateProducts(categories)
	if err := db.CreateInBatches(productData, 50).Error; err != nil {
		return err
	}

	var products []models.Product
	db.Find(&products)

	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	var prices []models.ProductPrice

	for _, product := range products {
		storeCount := 3 + rng.Intn(4)
		selectedStores := pickRandomStores(stores, storeCount, rng)

		for _, store := range selectedStores {
			basePrice := 50.0 + rng.Float64()*2000
			basePrice = float64(int(basePrice*100)) / 100

			var discountPrice *float64
			var discountPercent *float64
			if rng.Float64() < 0.25 {
				discount := basePrice * (0.05 + rng.Float64()*0.25)
				dp := basePrice - discount
				dp = float64(int(dp*100)) / 100
				discountPrice = &dp
				pct := (discount / basePrice) * 100
				pct = float64(int(pct*10)) / 10
				discountPercent = &pct
			}

			prices = append(prices, models.ProductPrice{
				ProductID:       product.ID,
				StoreID:         store.ID,
				Price:           basePrice,
				DiscountPrice:   discountPrice,
				DiscountPercent: discountPercent,
				Popularity:      rng.Intn(1000),
				InStock:         rng.Float64() > 0.05,
			})
		}
	}

	if err := db.CreateInBatches(prices, 100).Error; err != nil {
		return err
	}

	log.Printf("Seeded: %d categories, %d stores, %d products, %d prices, %d users",
		len(categories), len(stores), len(products), len(prices), len(users))
	return nil
}

func pickRandomStores(stores []models.Store, count int, rng *rand.Rand) []models.Store {
	indices := rng.Perm(len(stores))
	result := make([]models.Store, count)
	for i := 0; i < count; i++ {
		result[i] = stores[indices[i]]
	}
	return result
}

func generateProducts(categories []models.Category) []models.Product {
	brands := []string{"Простоквашино", "Danone", "Coca-Cola", "Nestle", "Unilever", "Bosch", "Samsung", "Philips", "IKEA", "L'Oreal", "Colgate", "Persil", "Макфа", "Heinz", "Barilla"}
	units := []string{"шт", "кг", "л", "уп", "пачка"}

	productNames := []struct {
		name     string
		category int
		brand    string
		unit     string
		weight   string
	}{
		{"Молоко 3.2% 1л", 0, "Простоквашино", "л", "1 л"},
		{"Хлеб белый нарезной", 0, "Макфа", "шт", "400 г"},
		{"Яйца С0 10шт", 0, "Роскар", "уп", "10 шт"},
		{"Сыр Российский 45%", 0, "Простоквашино", "кг", "200 г"},
		{"Куриное филе", 0, "Петелинка", "кг", "1 кг"},
		{"Гречка", 0, "Макфа", "кг", "900 г"},
		{"Рис длиннозерный", 0, "Макфа", "кг", "900 г"},
		{"Макароны спагетти", 0, "Barilla", "уп", "500 г"},
		{"Сахар-песок", 0, "Русский сахар", "кг", "1 кг"},
		{"Соль поваренная", 0, "Экстра", "уп", "1 кг"},
		{"Масло подсолнечное", 0, "Слобода", "л", "1 л"},
		{"Кефир 2.5%", 0, "Простоквашино", "л", "900 мл"},
		{"Творог 5%", 0, "Danone", "кг", "200 г"},
		{"Йогурт клубничный", 0, "Danone", "шт", "125 г"},
		{"Колбаса докторская", 0, "Черкизово", "кг", "400 г"},
		{"Сосиски молочные", 0, "Черкизово", "уп", "400 г"},
		{"Картофель", 0, "Фермер", "кг", "2 кг"},
		{"Морковь", 0, "Фермер", "кг", "1 кг"},
		{"Лук репчатый", 0, "Фермер", "кг", "1 кг"},
		{"Помидоры", 0, "Фермер", "кг", "500 г"},
		{"Огурцы", 0, "Фермер", "кг", "500 г"},
		{"Яблоки", 0, "Фермер", "кг", "1 кг"},
		{"Бананы", 0, "Chiquita", "кг", "1 кг"},
		{"Апельсины", 0, "Фермер", "кг", "1 кг"},
		{"Шоколад молочный", 0, "Nestle", "шт", "100 г"},
		{"Печенье овсяное", 0, "Юбилейное", "уп", "300 г"},
		{"Кетчуп", 0, "Heinz", "шт", "500 г"},
		{"Майонез", 0, "Calve", "шт", "400 г"},
		{"Горошек консервированный", 0, "Bonduelle", "шт", "400 г"},
		{"Тунец в собственном соку", 0, "Bonduelle", "шт", "185 г"},
		{"Coca-Cola 2л", 4, "Coca-Cola", "л", "2 л"},
		{"Sprite 1.5л", 4, "Coca-Cola", "л", "1.5 л"},
		{"Fanta апельсин 1.5л", 4, "Coca-Cola", "л", "1.5 л"},
		{"Сок апельсиновый", 4, "J7", "л", "1 л"},
		{"Вода минеральная", 4, "Aqua Minerale", "л", "1.5 л"},
		{"Чай черный", 4, "Lipton", "уп", "100 пак"},
		{"Кофе растворимый", 4, "Nescafe", "шт", "95 г"},
		{"Пиво светлое", 4, "Baltika", "л", "0.5 л"},
		{"Вино красное", 4, "Fanagoria", "л", "0.75 л"},
		{"Энергетик", 4, "Red Bull", "шт", "250 мл"},
		{"Стиральный порошок", 3, "Persil", "кг", "3 кг"},
		{"Гель для стирки", 3, "Persil", "л", "1.3 л"},
		{"Средство для мытья посуды", 3, "Fairy", "шт", "900 мл"},
		{"Средство для пола", 3, "Mr Proper", "л", "1 л"},
		{"Средство для ванной", 3, "Domestos", "шт", "750 мл"},
		{"Шампунь", 3, "Head & Shoulders", "шт", "400 мл"},
		{"Гель для душа", 3, "Nivea", "шт", "250 мл"},
		{"Зубная паста", 3, "Colgate", "шт", "100 мл"},
		{"Туалетная бумага 8 рул", 3, "Zewa", "уп", "8 шт"},
		{"Бумажные полотенца", 3, "Zewa", "уп", "2 шт"},
		{"Губки для мытья", 1, "Vileda", "уп", "3 шт"},
		{"Перчатки резиновые", 1, "Vileda", "уп", "1 пара"},
		{"Ведро пластиковое", 1, "Vileda", "шт", "10 л"},
		{"Швабра", 1, "Vileda", "шт", "1 шт"},
		{"Мусорные пакеты", 1, "Vileda", "уп", "30 шт"},
		{"Цемент М500", 2, "Евроцемент", "кг", "50 кг"},
		{"Песок строительный", 2, "BuildMax", "кг", "25 кг"},
		{"Кирпич облицовочный", 2, "BuildMax", "шт", "1 шт"},
		{"Штукатурка", 2, "Knauf", "кг", "25 кг"},
		{"Краска белая", 2, "Tikkurila", "л", "9 л"},
		{"Обои", 2, "Erismann", "рулон", "10 м"},
		{"Ламинат", 2, "Quick-Step", "уп", "2 м²"},
		{"Дрель электрическая", 2, "Bosch", "шт", "1 шт"},
		{"Молоток", 2, "Stanley", "шт", "1 шт"},
		{"Набор отверток", 2, "Stanley", "уп", "6 шт"},
		{"Лампа настольная", 5, "IKEA", "шт", "1 шт"},
		{"Подушка", 5, "IKEA", "шт", "50x70"},
		{"Одеяло", 5, "IKEA", "шт", "200x200"},
		{"Ковер", 5, "IKEA", "шт", "160x230"},
		{"Ваза декоративная", 5, "IKEA", "шт", "30 см"},
		{"Фото рамка", 5, "IKEA", "шт", "20x30"},
		{"Часы настенные", 5, "IKEA", "шт", "30 см"},
		{"Шторы", 5, "IKEA", "шт", "140x260"},
		{"Смартфон", 6, "Samsung", "шт", "128 ГБ"},
		{"Наушники беспроводные", 6, "Samsung", "шт", "1 шт"},
		{"Планшет", 6, "Samsung", "шт", "64 ГБ"},
		{"Телевизор 55\"", 6, "Samsung", "шт", "55\""},
		{"Ноутбук", 6, "Samsung", "шт", "15.6\""},
		{"Микроволновка", 6, "Samsung", "шт", "23 л"},
		{"Холодильник", 6, "Samsung", "шт", "350 л"},
		{"Стиральная машина", 6, "Samsung", "шт", "7 кг"},
		{"Пылесос", 6, "Philips", "шт", "1 шт"},
		{"Утюг", 6, "Philips", "шт", "1 шт"},
		{"Фен", 6, "Philips", "шт", "1 шт"},
		{"Блендер", 6, "Philips", "шт", "1 шт"},
		{"Кофемашина", 6, "Philips", "шт", "1 шт"},
		{"Электрическая зубная щетка", 6, "Philips", "шт", "1 шт"},
		{"Powerbank", 6, "Samsung", "шт", "10000 mAh"},
		{"USB-кабель", 6, "Samsung", "шт", "1 м"},
		{"Клавиатура", 6, "Logitech", "шт", "1 шт"},
		{"Мышь компьютерная", 6, "Logitech", "шт", "1 шт"},
		{"Веб-камера", 6, "Logitech", "шт", "1080p"},
		{"Колонка Bluetooth", 6, "JBL", "шт", "1 шт"},
		{"Роутер Wi-Fi", 6, "TP-Link", "шт", "1 шт"},
		{"SSD 512GB", 6, "Samsung", "шт", "512 ГБ"},
		{"Флешка 64GB", 6, "Samsung", "шт", "64 ГБ"},
		{"Карта памяти 128GB", 6, "Samsung", "шт", "128 ГБ"},
		{"Зарядное устройство", 6, "Samsung", "шт", "25W"},
		{"Чехол для телефона", 7, "Spigen", "шт", "1 шт"},
		{"Защитное стекло", 7, "Spigen", "шт", "1 шт"},
		{"Батарейки AA", 7, "Duracell", "уп", "4 шт"},
		{"Скотч упаковочный", 7, "Scotch", "шт", "48 мм"},
		{"Клей канцелярский", 7, "Pritt", "шт", "40 г"},
		{"Блокнот А5", 7, "Brauberg", "шт", "48 л"},
		{"Ручка шариковая", 7, "Pilot", "уп", "10 шт"},
		{"Карандаши", 7, "Faber-Castell", "уп", "12 шт"},
		{"Ножницы", 7, "Faber-Castell", "шт", "1 шт"},
		{"Линейка", 7, "Faber-Castell", "шт", "30 см"},
		{"Скрепки", 7, "Attache", "уп", "100 шт"},
		{"Степлер", 7, "Attache", "шт", "1 шт"},
		{"Дырокол", 7, "Attache", "шт", "1 шт"},
		{"Папка для документов", 7, "Attache", "шт", "1 шт"},
		{"Корректирующая жидкость", 7, "Erich Krause", "шт", "20 мл"},
		{"Маркеры", 7, "Stabilo", "уп", "4 шт"},
		{"Клей-карандаш", 7, "Pritt", "шт", "22 г"},
		{"Ластик", 7, "Faber-Castell", "шт", "1 шт"},
		{"Точилка", 7, "Faber-Castell", "шт", "1 шт"},
		{"Обложки для тетрадей", 7, "Attache", "уп", "10 шт"},
		{"Тетрадь 48л", 7, "Brauberg", "шт", "48 л"},
		{"Дневник", 7, "Brauberg", "шт", "96 л"},
		{"Альбом для рисования", 7, "Brauberg", "шт", "20 л"},
		{"Краски акварельные", 7, "Faber-Castell", "уп", "12 цв"},
		{"Кисти для рисования", 7, "Faber-Castell", "уп", "5 шт"},
		{"Пластилин", 7, "Faber-Castell", "уп", "10 цв"},
	}

	var products []models.Product
	for i, p := range productNames {
		catIdx := p.category
		if catIdx >= len(categories) {
			catIdx = len(categories) - 1
		}
		brand := p.brand
		if brand == "" {
			brand = brands[i%len(brands)]
		}
		unit := p.unit
		if unit == "" {
			unit = units[i%len(units)]
		}
		products = append(products, models.Product{
			Name:         p.name,
			Description:  fmt.Sprintf("Качественный товар «%s» от бренда %s. %s.", p.name, brand, p.weight),
			CategoryID:   categories[catIdx].ID,
			Image:        ProductImage(p.name, categories[catIdx].Slug, i),
			Unit:         unit,
			Brand:        brand,
			WeightVolume: p.weight,
			Barcode:      fmt.Sprintf("4607%08d", 10000000+i),
		})
	}

	for i := len(productNames); i < 120; i++ {
		catIdx := i % len(categories)
		brand := brands[i%len(brands)]
		unit := units[i%len(units)]
		name := fmt.Sprintf("Товар %s %d", categories[catIdx].Name, i-len(productNames)+1)
		products = append(products, models.Product{
			Name:         name,
			Description:  fmt.Sprintf("Описание товара «%s» от бренда %s.", name, brand),
			CategoryID:   categories[catIdx].ID,
			Image:        ProductImage(name, categories[catIdx].Slug, i),
			Unit:         unit,
			Brand:        brand,
			WeightVolume: fmt.Sprintf("%d %s", (i%10)+1, unit),
			Barcode:      fmt.Sprintf("4607%08d", 20000000+i),
		})
	}

	_ = strings.Join(brands, ",")
	return products
}
