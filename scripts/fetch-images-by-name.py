#!/usr/bin/env python3
"""Скачивает первое фото из поиска по названию товара."""
from __future__ import annotations

import json
import os
import re
import ssl
import sys
import time
import urllib.error
import urllib.parse
import urllib.request

SCRIPT_DIR = os.path.dirname(os.path.abspath(__file__))
ROOT = os.path.join(SCRIPT_DIR, "..", "frontend", "public", "images")
PRODUCTS = os.path.join(ROOT, "products")
CATEGORIES = os.path.join(ROOT, "categories")

UA = "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36"
CTX = ssl.create_default_context()

# slug -> русское название товара (из seed)
PRODUCT_NAMES: dict[str, str] = {
    "milk": "Молоко 3.2% 1л",
    "bread": "Хлеб белый нарезной",
    "eggs": "Яйца С0 10шт",
    "cheese": "Сыр Российский 45%",
    "chicken": "Куриное филе",
    "buckwheat": "Гречка",
    "rice": "Рис длиннозерный",
    "pasta": "Макароны спагетти",
    "sugar": "Сахар-песок",
    "salt": "Соль поваренная",
    "sunflower-oil": "Масло подсолнечное",
    "kefir": "Кефир 2.5%",
    "cottage-cheese": "Творог 5%",
    "yogurt": "Йогурт клубничный",
    "sausage": "Колбаса докторская",
    "hotdog": "Сосиски молочные",
    "potato": "Картофель",
    "carrot": "Морковь",
    "onion": "Лук репчатый",
    "tomato": "Помидоры",
    "cucumber": "Огурцы",
    "apple": "Яблоки",
    "banana": "Бананы",
    "orange": "Апельсины",
    "chocolate": "Шоколад молочный",
    "cookies": "Печенье овсяное",
    "ketchup": "Кетчуп",
    "mayonnaise": "Майонез",
    "peas": "Горошек консервированный",
    "tuna": "Тунец в собственном соку",
    "coca-cola": "Coca-Cola 2л",
    "sprite": "Sprite 1.5л",
    "fanta": "Fanta апельсин 1.5л",
    "juice": "Сок апельсиновый",
    "water": "Вода минеральная",
    "tea": "Чай черный",
    "coffee": "Кофе растворимый",
    "beer": "Пиво светлое",
    "wine": "Вино красное",
    "energy-drink": "Энергетик Red Bull",
    "laundry-powder": "Стиральный порошок",
    "laundry-gel": "Гель для стирки",
    "dish-soap": "Средство для мытья посуды",
    "floor-cleaner": "Средство для мытья пола",
    "bathroom-cleaner": "Средство для ванной",
    "shampoo": "Шампунь",
    "shower-gel": "Гель для душа",
    "toothpaste": "Зубная паста",
    "toilet-paper": "Туалетная бумага",
    "paper-towels": "Бумажные полотенца",
    "sponge": "Губки для мытья посуды",
    "rubber-gloves": "Резиновые перчатки",
    "bucket": "Пластиковое ведро",
    "mop": "Швабра",
    "trash-bags": "Мусорные пакеты",
    "cement": "Цемент М500 мешок",
    "sand": "Песок строительный",
    "brick": "Кирпич облицовочный",
    "plaster": "Штукатурка строительная",
    "paint": "Краска белая",
    "wallpaper": "Обои для стен",
    "laminate": "Ламинат напольный",
    "drill": "Дрель электрическая",
    "hammer": "Молоток",
    "screwdriver": "Набор отверток",
    "lamp": "Лампа настольная",
    "pillow": "Подушка",
    "blanket": "Одеяло",
    "carpet": "Ковер",
    "vase": "Ваза декоративная",
    "photo-frame": "Фото рамка",
    "wall-clock": "Часы настенные",
    "curtains": "Шторы",
    "smartphone": "Смартфон",
    "headphones": "Наушники беспроводные",
    "tablet": "Планшет",
    "tv": "Телевизор 55 дюймов",
    "laptop": "Ноутбук",
    "microwave": "Микроволновка",
    "fridge": "Холодильник",
    "washing-machine": "Стиральная машина",
    "vacuum": "Пылесос",
    "iron": "Утюг",
    "hairdryer": "Фен для волос",
    "blender": "Блендер кухонный",
    "coffee-machine": "Кофемашина",
    "electric-toothbrush": "Электрическая зубная щетка",
    "powerbank": "Powerbank зарядка",
    "usb-cable": "USB кабель",
    "keyboard": "Клавиатура компьютерная",
    "computer-mouse": "Компьютерная мышь",
    "webcam": "Веб-камера",
    "bluetooth-speaker": "Колонка Bluetooth",
    "wifi-router": "Роутер Wi-Fi",
    "ssd": "SSD накопитель 512GB",
    "flash-drive": "Флешка USB 64GB",
    "memory-card": "Карта памяти 128GB",
    "charger": "Зарядное устройство для телефона",
    "phone-case": "Чехол для телефона",
    "screen-protector": "Защитное стекло для телефона",
    "batteries": "Батарейки AA",
    "tape": "Скотч упаковочный",
    "glue": "Клей канцелярский",
    "notebook": "Блокнот А5",
    "pen": "Ручка шариковая",
    "pencils": "Карандаши",
    "scissors": "Ножницы",
    "ruler": "Линейка",
    "paperclips": "Скрепки",
    "stapler": "Степлер",
    "hole-punch": "Дырокол",
    "folder": "Папка для документов",
    "correction-fluid": "Корректирующая жидкость",
    "markers": "Маркеры",
    "glue-stick": "Клей-карандаш",
    "eraser": "Ластик",
    "pencil-sharpener": "Точилка для карандашей",
    "notebook-covers": "Обложки для тетрадей",
    "exercise-book": "Тетрадь 48 листов",
    "diary": "Дневник школьный",
    "drawing-album": "Альбом для рисования",
    "watercolors": "Краски акварельные",
    "paint-brushes": "Кисти для рисования",
    "play-dough": "Пластилин",
}

CATEGORY_NAMES: dict[str, str] = {
    "food": "Пищевые продукты",
    "household": "Хозяйственные товары",
    "construction": "Строительные материалы",
    "chemistry": "Бытовая химия",
    "drinks": "Напитки",
    "home": "Товары для дома",
    "electronics": "Электроника",
    "other": "Канцтовары",
}


def http_get(url: str, timeout: int = 30) -> bytes | None:
    req = urllib.request.Request(url, headers={"User-Agent": UA, "Accept": "image/*,*/*"})
    try:
        with urllib.request.urlopen(req, context=CTX, timeout=timeout) as resp:
            return resp.read()
    except Exception:
        return None


def bing_image_urls(query: str, count: int = 8) -> list[str]:
    q = urllib.parse.quote(query)
    url = f"https://www.bing.com/images/search?q={q}&form=HDRSC2&first=1"
    req = urllib.request.Request(url, headers={"User-Agent": UA, "Accept-Language": "ru-RU,ru;q=0.9"})
    try:
        with urllib.request.urlopen(req, context=CTX, timeout=30) as resp:
            html = resp.read().decode("utf-8", errors="ignore")
    except Exception:
        return []

    urls: list[str] = []
    for m in re.finditer(r'murl&quot;:&quot;(https?://[^&]+?)&quot;', html):
        u = m.group(1).replace("\\u0026", "&")
        if u not in urls:
            urls.append(u)
        if len(urls) >= count:
            break

    if not urls:
        for m in re.finditer(r'"murl":"(https?://[^"]+)"', html):
            u = m.group(1)
            if u not in urls:
                urls.append(u)
            if len(urls) >= count:
                break
    return urls


def duckduckgo_image_urls(query: str, count: int = 8) -> list[str]:
    try:
        from duckduckgo_search import DDGS
    except ImportError:
        return []

    urls: list[str] = []
    try:
        with DDGS() as ddgs:
            for item in ddgs.images(query, max_results=count, safesearch="moderate"):
                u = item.get("image") or item.get("thumbnail")
                if u and u.startswith("http"):
                    urls.append(u)
    except Exception:
        pass
    return urls


def search_image_urls(query: str) -> list[str]:
    search_q = f"{query} фото товар"
    urls = duckduckgo_image_urls(search_q)
    if not urls:
        urls = bing_image_urls(search_q)
    if not urls:
        urls = duckduckgo_image_urls(query)
    if not urls:
        urls = bing_image_urls(query)
    return urls


def is_valid_image(data: bytes) -> bool:
    if len(data) < 3000:
        return False
    if data[:3] == b"\xff\xd8\xff":
        return True
    if data[:8] == b"\x89PNG\r\n\x1a\n":
        return True
    if data[:4] == b"RIFF" and data[8:12] == b"WEBP":
        return True
    if data[:6] in (b"GIF87a", b"GIF89a"):
        return True
    return False


def download_first_image(slug: str, name: str, dest_dir: str) -> bool:
    dest = os.path.join(dest_dir, f"{slug}.jpg")
    urls = search_image_urls(name)

    for i, img_url in enumerate(urls):
        data = http_get(img_url)
        if data and is_valid_image(data):
            with open(dest, "wb") as f:
                f.write(data)
            print(f"  ✓ {slug}.jpg ← {name!r} ({len(data)//1024} KB, src #{i+1})")
            return True
        time.sleep(0.3)

    print(f"  ✗ {slug}: не удалось скачать для {name!r}")
    return False


def main() -> int:
    os.makedirs(PRODUCTS, exist_ok=True)
    os.makedirs(CATEGORIES, exist_ok=True)

    ok = 0
    total = len(PRODUCT_NAMES) + len(CATEGORY_NAMES)

    print(f"=== Товары ({len(PRODUCT_NAMES)}) ===")
    for n, (slug, name) in enumerate(PRODUCT_NAMES.items(), 1):
        print(f"[{n}/{len(PRODUCT_NAMES)}] {name}")
        if download_first_image(slug, name, PRODUCTS):
            ok += 1
        time.sleep(1.0)

    print(f"\n=== Категории ({len(CATEGORY_NAMES)}) ===")
    for n, (slug, name) in enumerate(CATEGORY_NAMES.items(), 1):
        print(f"[{n}/{len(CATEGORY_NAMES)}] {name}")
        if download_first_image(slug, name, CATEGORIES):
            ok += 1
        time.sleep(1.0)

    print(f"\nГотово: {ok}/{total}")
    return 0 if ok == total else 1


if __name__ == "__main__":
    sys.exit(main())
