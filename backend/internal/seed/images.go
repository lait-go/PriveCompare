package seed

import (
	"log"
	"strings"

	"pricecompare/internal/models"

	"gorm.io/gorm"
)

const imgBase = "/images/products"
const catBase = "/images/categories"

func img(file string) string  { return imgBase + "/" + file + ".jpg" }
func catImg(file string) string { return catBase + "/" + file + ".jpg" }

var categoryImages = map[string]string{
	"food":         catImg("food"),
	"household":    catImg("household"),
	"construction": catImg("construction"),
	"chemistry":    catImg("chemistry"),
	"drinks":       catImg("drinks"),
	"home":         catImg("home"),
	"electronics":  catImg("electronics"),
	"other":        catImg("other"),
}

var storeLogos = map[string]string{
	"SuperMart":  "https://ui-avatars.com/api/?name=SM&background=2563eb&color=fff&size=128&bold=true&format=png",
	"FreshFood":  "https://ui-avatars.com/api/?name=FF&background=16a34a&color=fff&size=128&bold=true&format=png",
	"MegaStore":  "https://ui-avatars.com/api/?name=MS&background=dc2626&color=fff&size=128&bold=true&format=png",
	"BuildMax":   "https://ui-avatars.com/api/?name=BM&background=f59e0b&color=fff&size=128&bold=true&format=png",
	"HomeShop":   "https://ui-avatars.com/api/?name=HS&background=7c3aed&color=fff&size=128&bold=true&format=png",
	"CityMarket": "https://ui-avatars.com/api/?name=CM&background=0891b2&color=fff&size=128&bold=true&format=png",
}

// productImages — точное соответствие названия товара и изображения
var productImages = map[string]string{
	"Молоко 3.2% 1л":             img("milk"),
	"Хлеб белый нарезной":        img("bread"),
	"Яйца С0 10шт":               img("eggs"),
	"Сыр Российский 45%":         img("cheese"),
	"Куриное филе":               img("chicken"),
	"Гречка":                     img("buckwheat"),
	"Рис длиннозерный":           img("rice"),
	"Макароны спагетти":          img("pasta"),
	"Сахар-песок":                img("sugar"),
	"Соль поваренная":            img("salt"),
	"Масло подсолнечное":         img("sunflower-oil"),
	"Кефир 2.5%":                 img("kefir"),
	"Творог 5%":                  img("cottage-cheese"),
	"Йогурт клубничный":          img("yogurt"),
	"Колбаса докторская":         img("sausage"),
	"Сосиски молочные":           img("hotdog"),
	"Картофель":                  img("potato"),
	"Морковь":                    img("carrot"),
	"Лук репчатый":               img("onion"),
	"Помидоры":                   img("tomato"),
	"Огурцы":                     img("cucumber"),
	"Яблоки":                     img("apple"),
	"Бананы":                     img("banana"),
	"Апельсины":                  img("orange"),
	"Шоколад молочный":           img("chocolate"),
	"Печенье овсяное":            img("cookies"),
	"Кетчуп":                     img("ketchup"),
	"Майонез":                    img("mayonnaise"),
	"Горошек консервированный":   img("peas"),
	"Тунец в собственном соку":   img("tuna"),
	"Coca-Cola 2л":               img("coca-cola"),
	"Sprite 1.5л":                img("sprite"),
	"Fanta апельсин 1.5л":        img("fanta"),
	"Сок апельсиновый":           img("juice"),
	"Вода минеральная":            img("water"),
	"Чай черный":                 img("tea"),
	"Кофе растворимый":           img("coffee"),
	"Пиво светлое":               img("beer"),
	"Вино красное":               img("wine"),
	"Энергетик":                  img("energy-drink"),
	"Стиральный порошок":         img("laundry-powder"),
	"Гель для стирки":            img("laundry-gel"),
	"Средство для мытья посуды":  img("dish-soap"),
	"Средство для пола":          img("floor-cleaner"),
	"Средство для ванной":        img("bathroom-cleaner"),
	"Шампунь":                    img("shampoo"),
	"Гель для душа":              img("shower-gel"),
	"Зубная паста":               img("toothpaste"),
	"Туалетная бумага 8 рул":     img("toilet-paper"),
	"Бумажные полотенца":         img("paper-towels"),
	"Губки для мытья":            img("sponge"),
	"Перчатки резиновые":         img("rubber-gloves"),
	"Ведро пластиковое":          img("bucket"),
	"Швабра":                     img("mop"),
	"Мусорные пакеты":            img("trash-bags"),
	"Цемент М500":                img("cement"),
	"Песок строительный":         img("sand"),
	"Кирпич облицовочный":        img("brick"),
	"Штукатурка":                 img("plaster"),
	"Краска белая":               img("paint"),
	"Обои":                       img("wallpaper"),
	"Ламинат":                    img("laminate"),
	"Дрель электрическая":        img("drill"),
	"Молоток":                    img("hammer"),
	"Набор отверток":             img("screwdriver"),
	"Лампа настольная":           img("lamp"),
	"Подушка":                    img("pillow"),
	"Одеяло":                     img("blanket"),
	"Ковер":                      img("carpet"),
	"Ваза декоративная":          img("vase"),
	"Фото рамка":                 img("photo-frame"),
	"Часы настенные":             img("wall-clock"),
	"Шторы":                      img("curtains"),
	"Смартфон":                   img("smartphone"),
	"Наушники беспроводные":      img("headphones"),
	"Планшет":                    img("tablet"),
	"Телевизор 55\"":             img("tv"),
	"Ноутбук":                    img("laptop"),
	"Микроволновка":              img("microwave"),
	"Холодильник":                img("fridge"),
	"Стиральная машина":          img("washing-machine"),
	"Пылесос":                    img("vacuum"),
	"Утюг":                       img("iron"),
	"Фен":                        img("hairdryer"),
	"Блендер":                    img("blender"),
	"Кофемашина":                 img("coffee-machine"),
	"Электрическая зубная щетка": img("electric-toothbrush"),
	"Powerbank":                  img("powerbank"),
	"USB-кабель":                 img("usb-cable"),
	"Клавиатура":                 img("keyboard"),
	"Мышь компьютерная":          img("computer-mouse"),
	"Веб-камера":                 img("webcam"),
	"Колонка Bluetooth":          img("bluetooth-speaker"),
	"Роутер Wi-Fi":               img("wifi-router"),
	"SSD 512GB":                  img("ssd"),
	"Флешка 64GB":                img("flash-drive"),
	"Карта памяти 128GB":         img("memory-card"),
	"Зарядное устройство":        img("charger"),
	"Чехол для телефона":         img("phone-case"),
	"Защитное стекло":            img("screen-protector"),
	"Батарейки AA":               img("batteries"),
	"Скотч упаковочный":          img("tape"),
	"Клей канцелярский":          img("glue"),
	"Блокнот А5":                 img("notebook"),
	"Ручка шариковая":            img("pen"),
	"Карандаши":                  img("pencils"),
	"Ножницы":                    img("scissors"),
	"Линейка":                    img("ruler"),
	"Скрепки":                    img("paperclips"),
	"Степлер":                    img("stapler"),
	"Дырокол":                    img("hole-punch"),
	"Папка для документов":       img("folder"),
	"Корректирующая жидкость":    img("correction-fluid"),
	"Маркеры":                    img("markers"),
	"Клей-карандаш":              img("glue-stick"),
	"Ластик":                     img("eraser"),
	"Точилка":                    img("pencil-sharpener"),
	"Обложки для тетрадей":       img("notebook-covers"),
	"Тетрадь 48л":                img("exercise-book"),
	"Дневник":                    img("diary"),
	"Альбом для рисования":       img("drawing-album"),
	"Краски акварельные":         img("watercolors"),
	"Кисти для рисования":        img("paint-brushes"),
	"Пластилин":                  img("play-dough"),
}

// keywordImages — для дополнительных товаров по ключевым словам в названии
var keywordImages = []struct {
	keys []string
	file string
}{
	{[]string{"молоко"}, "milk"},
	{[]string{"кефир"}, "kefir"},
	{[]string{"творог"}, "cottage-cheese"},
	{[]string{"йогурт"}, "yogurt"},
	{[]string{"сыр"}, "cheese"},
	{[]string{"хлеб"}, "bread"},
	{[]string{"яйц"}, "eggs"},
	{[]string{"куриц", "филе"}, "chicken"},
	{[]string{"колбас"}, "sausage"},
	{[]string{"сосиск"}, "hotdog"},
	{[]string{"греч"}, "buckwheat"},
	{[]string{"рис"}, "rice"},
	{[]string{"макарон", "спагетти"}, "pasta"},
	{[]string{"сахар"}, "sugar"},
	{[]string{"соль"}, "salt"},
	{[]string{"масло"}, "sunflower-oil"},
	{[]string{"картоф"}, "potato"},
	{[]string{"морков"}, "carrot"},
	{[]string{"лук"}, "onion"},
	{[]string{"помидор"}, "tomato"},
	{[]string{"огурц"}, "cucumber"},
	{[]string{"яблок"}, "apple"},
	{[]string{"банан"}, "banana"},
	{[]string{"апельсин"}, "orange"},
	{[]string{"шоколад"}, "chocolate"},
	{[]string{"печень"}, "cookies"},
	{[]string{"кетчуп"}, "ketchup"},
	{[]string{"майонез"}, "mayonnaise"},
	{[]string{"горош"}, "peas"},
	{[]string{"тунец"}, "tuna"},
	{[]string{"cola", "coca"}, "coca-cola"},
	{[]string{"sprite"}, "sprite"},
	{[]string{"fanta"}, "fanta"},
	{[]string{"сок"}, "juice"},
	{[]string{"вода"}, "water"},
	{[]string{"чай"}, "tea"},
	{[]string{"кофе"}, "coffee"},
	{[]string{"пиво"}, "beer"},
	{[]string{"вино"}, "wine"},
	{[]string{"энергетик"}, "energy-drink"},
	{[]string{"стиральн", "порошок"}, "laundry-powder"},
	{[]string{"гель для стирки"}, "laundry-gel"},
	{[]string{"мытья посуды"}, "dish-soap"},
	{[]string{"для пола"}, "floor-cleaner"},
	{[]string{"для ванной"}, "bathroom-cleaner"},
	{[]string{"шампун"}, "shampoo"},
	{[]string{"душа"}, "shower-gel"},
	{[]string{"зубная", "зубн"}, "toothpaste"},
	{[]string{"туалетн"}, "toilet-paper"},
	{[]string{"полотенц"}, "paper-towels"},
	{[]string{"губк"}, "sponge"},
	{[]string{"перчат"}, "rubber-gloves"},
	{[]string{"ведро"}, "bucket"},
	{[]string{"швабр"}, "mop"},
	{[]string{"мусор"}, "trash-bags"},
	{[]string{"цемент"}, "cement"},
	{[]string{"песок"}, "sand"},
	{[]string{"кирпич"}, "brick"},
	{[]string{"штукатур"}, "plaster"},
	{[]string{"краск"}, "paint"},
	{[]string{"обои"}, "wallpaper"},
	{[]string{"ламинат"}, "laminate"},
	{[]string{"дрель"}, "drill"},
	{[]string{"молоток"}, "hammer"},
	{[]string{"отверт"}, "screwdriver"},
	{[]string{"ламп"}, "lamp"},
	{[]string{"подуш"}, "pillow"},
	{[]string{"одеял"}, "blanket"},
	{[]string{"ковер", "ковёр"}, "carpet"},
	{[]string{"ваза"}, "vase"},
	{[]string{"рамк"}, "photo-frame"},
	{[]string{"часы"}, "wall-clock"},
	{[]string{"штор"}, "curtains"},
	{[]string{"смартфон", "телефон"}, "smartphone"},
	{[]string{"наушник"}, "headphones"},
	{[]string{"планшет"}, "tablet"},
	{[]string{"телевизор"}, "tv"},
	{[]string{"ноутбук"}, "laptop"},
	{[]string{"микроволнов"}, "microwave"},
	{[]string{"холодильник"}, "fridge"},
	{[]string{"стиральн"}, "washing-machine"},
	{[]string{"пылесос"}, "vacuum"},
	{[]string{"утюг"}, "iron"},
	{[]string{"фен"}, "hairdryer"},
	{[]string{"блендер"}, "blender"},
	{[]string{"кофемашин"}, "coffee-machine"},
	{[]string{"щетк"}, "electric-toothbrush"},
	{[]string{"powerbank"}, "powerbank"},
	{[]string{"usb", "кабель"}, "usb-cable"},
	{[]string{"клавиатур"}, "keyboard"},
	{[]string{"мыш"}, "computer-mouse"},
	{[]string{"веб-камер", "камер"}, "webcam"},
	{[]string{"колонк"}, "bluetooth-speaker"},
	{[]string{"роутер", "wi-fi"}, "wifi-router"},
	{[]string{"ssd"}, "ssd"},
	{[]string{"флеш"}, "flash-drive"},
	{[]string{"карта памяти"}, "memory-card"},
	{[]string{"заряд"}, "charger"},
	{[]string{"чехол"}, "phone-case"},
	{[]string{"стекл"}, "screen-protector"},
	{[]string{"батарей"}, "batteries"},
	{[]string{"скотч"}, "tape"},
	{[]string{"блокнот"}, "notebook"},
	{[]string{"ручк"}, "pen"},
	{[]string{"карандаш"}, "pencils"},
	{[]string{"ножниц"}, "scissors"},
	{[]string{"линейк"}, "ruler"},
	{[]string{"скреп"}, "paperclips"},
	{[]string{"степлер"}, "stapler"},
	{[]string{"дырокол"}, "hole-punch"},
	{[]string{"папк"}, "folder"},
	{[]string{"корректир"}, "correction-fluid"},
	{[]string{"маркер"}, "markers"},
	{[]string{"клей-карандаш"}, "glue-stick"},
	{[]string{"клей"}, "glue"},
	{[]string{"ластик"}, "eraser"},
	{[]string{"точил"}, "pencil-sharpener"},
	{[]string{"обложк"}, "notebook-covers"},
	{[]string{"тетрад"}, "exercise-book"},
	{[]string{"дневник"}, "diary"},
	{[]string{"альбом"}, "drawing-album"},
	{[]string{"акварель", "краски"}, "watercolors"},
	{[]string{"кист"}, "paint-brushes"},
	{[]string{"пластилин"}, "play-dough"},
}

var categoryProductPool = map[string][]string{
	"food":         {"milk", "bread", "eggs", "cheese", "chicken", "potato", "apple"},
	"household":    {"sponge", "bucket", "mop", "rubber-gloves", "trash-bags"},
	"construction": {"cement", "brick", "hammer", "drill", "paint", "laminate"},
	"chemistry":    {"laundry-powder", "shampoo", "dish-soap", "toothpaste", "toilet-paper"},
	"drinks":       {"coca-cola", "juice", "water", "tea", "coffee", "beer"},
	"home":         {"lamp", "pillow", "blanket", "carpet", "curtains", "vase"},
	"electronics":  {"smartphone", "laptop", "headphones", "tv", "fridge", "vacuum"},
	"other":        {"notebook", "pen", "pencils", "scissors", "markers", "glue"},
}

func ProductImage(name, categorySlug string, index int) string {
	if u, ok := productImages[name]; ok {
		return u
	}
	lower := strings.ToLower(name)
	for _, kw := range keywordImages {
		for _, k := range kw.keys {
			if strings.Contains(lower, k) {
				return img(kw.file)
			}
		}
	}
	pool := categoryProductPool[categorySlug]
	if len(pool) == 0 {
		pool = categoryProductPool["other"]
	}
	return img(pool[index%len(pool)])
}

func CategoryImage(slug string) string {
	if u, ok := categoryImages[slug]; ok {
		return u
	}
	return catImg("other")
}

func StoreLogo(name string) string {
	if logo, ok := storeLogos[name]; ok {
		return logo
	}
	return "https://ui-avatars.com/api/?name=" + name + "&background=64748b&color=fff&size=128&bold=true&format=png"
}

func UpdateImages(db *gorm.DB) error {
	var categories []models.Category
	if err := db.Find(&categories).Error; err != nil {
		return err
	}
	for _, c := range categories {
		if err := db.Model(&c).Update("image", CategoryImage(c.Slug)).Error; err != nil {
			return err
		}
	}

	var stores []models.Store
	if err := db.Find(&stores).Error; err != nil {
		return err
	}
	for _, s := range stores {
		if err := db.Model(&s).Update("logo", StoreLogo(s.Name)).Error; err != nil {
			return err
		}
	}

	slugByID := make(map[uint]string, len(categories))
	for _, c := range categories {
		slugByID[c.ID] = c.Slug
	}

	var products []models.Product
	if err := db.Find(&products).Error; err != nil {
		return err
	}
	for i, prod := range products {
		imageURL := ProductImage(prod.Name, slugByID[prod.CategoryID], i)
		if err := db.Model(&prod).Update("image", imageURL).Error; err != nil {
			return err
		}
	}

	log.Printf("Updated images for %d categories, %d stores, %d products",
		len(categories), len(stores), len(products))
	return nil
}
