#!/usr/bin/env python3
"""Скачивает тематические изображения товаров с Wikimedia Commons."""
from __future__ import annotations

import argparse
import json
import os
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
SEARCHES_FILE = os.path.join(SCRIPT_DIR, "product-image-searches.json")
MANIFEST_FILE = os.path.join(SCRIPT_DIR, "product-image-manifest.json")

UA = "PriceCompare/1.0 (demo marketplace; contact: dev@localhost)"
CTX = ssl.create_default_context()
API = "https://commons.wikimedia.org/w/api.php"

# Проверенные вручную файлы (приоритет над поиском)
VERIFIED: dict[str, str] = {
    "milk": "Milk_glass.jpg",
    "bread": "Fresh_made_bread_05.jpg",
    "eggs": "Eggs in basket 2020 G1.jpg",
    "cheese": "Emmental cheese.jpg",
    "chicken": "Raw chicken slices.jpg",
    "buckwheat": "Buckwheat and products from it.jpg",
    "rice": "White rice.jpg",
    "pasta": "Spaghetti.jpg",
    "sugar": "White sugar cubes.jpg",
    "salt": "Salt shaker on white background.jpg",
    "sunflower-oil": "Sunflower oil.jpg",
    "kefir": "Kefir in glass.jpg",
    "cottage-cheese": "Cottage cheese 01.jpg",
    "yogurt": "Strawberry yogurt.jpg",
    "sausage": "Kielbasa.jpg",
    "hotdog": "Hot dogs with mustard.jpg",
    "potato": "Potatoes in basket.jpg",
    "carrot": "Carrots of many colors.jpg",
    "onion": "Yellow onions.jpg",
    "tomato": "Tomatoes.jpg",
    "cucumber": "Cucumber fruit.jpg",
    "apple": "Red Apple.jpg",
    "banana": "Bananas.jpg",
    "orange": "Orange-Fruit-Pieces.jpg",
    "chocolate": "Chocolate bar.jpg",
    "cookies": "Oatmeal cookies.jpg",
    "ketchup": "Heinz Tomato Ketchup.jpg",
    "mayonnaise": "Mayonnaise.jpg",
    "peas": "Green peas in a bowl.jpg",
    "tuna": "Canned tuna.jpg",
    "coca-cola": "Coca-Cola bottle.jpg",
    "sprite": "Sprite lemon lime 1.jpg",
    "fanta": "Fanta Orange Glass Bottle.jpg",
    "juice": "Orange juice 1.jpg",
    "water": "Bottled water.jpg",
    "tea": "Black tea.jpg",
    "coffee": "Instant coffee.jpg",
    "beer": "Beer glasses.jpg",
    "wine": "Red wine in glass.jpg",
    "energy-drink": "Red Bull can.jpg",
    "laundry-powder": "Laundry detergent.jpg",
    "laundry-gel": "Liquid laundry detergent.jpg",
    "dish-soap": "Dishwashing liquid.jpg",
    "floor-cleaner": "Floor cleaning.jpg",
    "bathroom-cleaner": "Bathroom cleaner.jpg",
    "shampoo": "Shampoo bottle.jpg",
    "shower-gel": "Shower gel.jpg",
    "toothpaste": "Toothpaste.jpg",
    "toilet-paper": "Toilet paper rolls.jpg",
    "paper-towels": "Paper towels.jpg",
    "sponge": "Kitchen sponge.jpg",
    "rubber-gloves": "Rubber gloves.jpg",
    "bucket": "Plastic bucket.jpg",
    "mop": "Mop.jpg",
    "trash-bags": "Garbage bags.jpg",
    "cement": "Portland cement bag.jpg",
    "sand": "Sand pile.jpg",
    "brick": "Brick.jpg",
    "plaster": "Plaster bag.jpg",
    "paint": "Paint can.jpg",
    "wallpaper": "Wallpaper rolls.jpg",
    "laminate": "Laminate flooring.jpg",
    "drill": "Power drill.jpg",
    "hammer": "Claw hammer.jpg",
    "screwdriver": "Screwdriver set.jpg",
    "lamp": "Desk lamp.jpg",
    "pillow": "Pillow.jpg",
    "blanket": "Blanket.jpg",
    "carpet": "Carpet.jpg",
    "vase": "Vase with flowers.jpg",
    "photo-frame": "Picture frame.jpg",
    "wall-clock": "Wall clock.jpg",
    "curtains": "Curtains.jpg",
    "smartphone": "Smartphone.jpg",
    "headphones": "Headphones.jpg",
    "tablet": "Tablet computer.jpg",
    "tv": "Flat screen television.jpg",
    "laptop": "Laptop.jpg",
    "microwave": "Microwave oven.jpg",
    "fridge": "Refrigerator.jpg",
    "washing-machine": "Washing machine.jpg",
    "vacuum": "Vacuum cleaner.jpg",
    "iron": "Clothes iron.jpg",
    "hairdryer": "Hair dryer.jpg",
    "blender": "Blender.jpg",
    "coffee-machine": "Espresso machine.jpg",
    "electric-toothbrush": "Electric toothbrush.jpg",
    "powerbank": "Power bank.jpg",
    "usb-cable": "USB cable.jpg",
    "keyboard": "Computer keyboard.jpg",
    "computer-mouse": "Computer mouse.jpg",
    "webcam": "Webcam.jpg",
    "bluetooth-speaker": "Bluetooth speaker.jpg",
    "wifi-router": "Wireless router.jpg",
    "ssd": "Solid state drive.jpg",
    "flash-drive": "USB flash drive.jpg",
    "memory-card": "SD memory card.jpg",
    "charger": "Phone charger.jpg",
    "phone-case": "Mobile phone case.jpg",
    "screen-protector": "Screen protector.jpg",
    "batteries": "AA batteries.jpg",
    "tape": "Adhesive tape.jpg",
    "glue": "Glue bottle.jpg",
    "notebook": "Notebook.jpg",
    "pen": "Ballpoint pen.jpg",
    "pencils": "Pencils.jpg",
    "scissors": "Scissors.jpg",
    "ruler": "Ruler.jpg",
    "paperclips": "Paper clips.jpg",
    "stapler": "Stapler.jpg",
    "hole-punch": "Hole punch.jpg",
    "folder": "Ring binder.jpg",
    "correction-fluid": "Correction fluid.jpg",
    "markers": "Markers.jpg",
    "glue-stick": "Glue stick.jpg",
    "eraser": "Eraser.jpg",
    "pencil-sharpener": "Pencil sharpener.jpg",
    "notebook-covers": "Notebook covers.jpg",
    "exercise-book": "Exercise book.jpg",
    "diary": "Diary.jpg",
    "drawing-album": "Drawing pad.jpg",
    "watercolors": "Watercolor paint set.jpg",
    "paint-brushes": "Paint brushes.jpg",
    "play-dough": "Play-Doh.jpg",
    "food": "Assorted food products.jpg",
    "household": "Cleaning supplies.jpg",
    "construction": "Construction materials.jpg",
    "chemistry": "Cleaning products.jpg",
    "drinks": "Soft drink bottles.jpg",
    "home": "Home decor.jpg",
    "electronics": "Consumer electronics.jpg",
    "other": "Office supplies.jpg",
}

CATEGORY_SLUGS = {
    "food", "household", "construction", "chemistry",
    "drinks", "home", "electronics", "other",
}


def load_json(path: str) -> dict:
    with open(path, encoding="utf-8") as f:
        return json.load(f)


def save_json(path: str, data: dict) -> None:
    with open(path, "w", encoding="utf-8") as f:
        json.dump(data, f, ensure_ascii=False, indent=2)
        f.write("\n")


def wiki_file_url(filename: str, width: int = 400) -> str:
    encoded = urllib.parse.quote(filename.replace(" ", "_"))
    return f"https://commons.wikimedia.org/wiki/Special:FilePath/{encoded}?width={width}"


def api_get(params: dict) -> dict:
    url = API + "?" + urllib.parse.urlencode(params)
    for attempt in range(6):
        req = urllib.request.Request(url, headers={"User-Agent": UA})
        try:
            with urllib.request.urlopen(req, context=CTX, timeout=45) as resp:
                return json.load(resp)
        except urllib.error.HTTPError as exc:
            if exc.code == 429 and attempt < 5:
                time.sleep(3 * (attempt + 1))
                continue
            raise


def search_wikimedia(query: str) -> str | None:
    data = api_get({
        "action": "query",
        "generator": "search",
        "gsrsearch": f"filetype:bitmap {query}",
        "gsrnamespace": "6",
        "gsrlimit": "5",
        "prop": "imageinfo",
        "iiprop": "url|mime",
        "iiurlwidth": "400",
        "format": "json",
    })
    pages = data.get("query", {}).get("pages", {})
    for page in sorted(pages.values(), key=lambda p: p.get("index", 99)):
        ii = page.get("imageinfo", [{}])[0]
        mime = ii.get("mime", "")
        if not mime.startswith("image/"):
            continue
        if ii.get("thumburl") or ii.get("url"):
            return page.get("title", "").replace("File:", "")
    return None


def http_get(url: str, timeout: int = 60) -> bytes | None:
    for attempt in range(5):
        req = urllib.request.Request(url, headers={"User-Agent": UA})
        try:
            with urllib.request.urlopen(req, context=CTX, timeout=timeout) as resp:
                return resp.read()
        except urllib.error.HTTPError as exc:
            if exc.code == 429 and attempt < 4:
                time.sleep(2 ** attempt + 1)
                continue
            return None
        except Exception:
            return None
    return None


def verify_file_exists(filename: str) -> bool:
    data = http_get(wiki_file_url(filename, 200), timeout=30)
    return data is not None and len(data) >= 1500


def download_file(local_name: str, filename: str, dest_dir: str, width: int = 400) -> bool:
    dest = os.path.join(dest_dir, f"{local_name}.jpg")
    data = http_get(wiki_file_url(filename, width))
    if data is None or len(data) < 2000:
        return False
    with open(dest, "wb") as f:
        f.write(data)
    return True


def resolve_manifest(searches: dict[str, str], delay: float = 1.2) -> dict[str, str]:
    manifest: dict[str, str] = {}
    items = list(searches.items())
    for i, (slug, query) in enumerate(items, 1):
        if slug in VERIFIED:
            manifest[slug] = VERIFIED[slug]
            print(f"[{i}/{len(items)}] {slug}: verified → {VERIFIED[slug]}")
            continue

        time.sleep(delay)
        found = search_wikimedia(query)
        if found:
            manifest[slug] = found
            print(f"[{i}/{len(items)}] {slug}: search → {found}")
        else:
            print(f"[{i}/{len(items)}] {slug}: NOT FOUND for '{query}'")
    return manifest


def build_manifest(force_resolve: bool = False) -> dict[str, str]:
    searches = load_json(SEARCHES_FILE)

    if force_resolve:
        print("Resolving Wikimedia filenames via API...")
        manifest = resolve_manifest(searches)
        save_json(MANIFEST_FILE, manifest)
        return manifest

    if os.path.exists(MANIFEST_FILE):
        manifest = load_json(MANIFEST_FILE)
        missing = [s for s in searches if s not in manifest]
        if not missing:
            return manifest
        print(f"Manifest incomplete ({len(missing)} missing), filling from verified/search...")

    manifest = dict(VERIFIED)
    for slug in searches:
        if slug not in manifest:
            manifest[slug] = VERIFIED.get(slug, "")
    manifest = {k: v for k, v in manifest.items() if v}
    save_json(MANIFEST_FILE, manifest)
    return manifest


def download_all(manifest: dict[str, str]) -> int:
    os.makedirs(PRODUCTS, exist_ok=True)
    os.makedirs(CATEGORIES, exist_ok=True)

    ok = 0
    failed: list[str] = []

    for slug, filename in manifest.items():
        dest = CATEGORIES if slug in CATEGORY_SLUGS else PRODUCTS
        label = "category" if slug in CATEGORY_SLUGS else "product"

        if download_file(slug, filename, dest):
            print(f"  ✓ [{label}] {slug}.jpg ← {filename}")
            ok += 1
            time.sleep(0.8)
            continue

        # Попробовать поиск как fallback
        searches = load_json(SEARCHES_FILE)
        query = searches.get(slug)
        if query:
            time.sleep(2.0)
            alt = search_wikimedia(query)
            if alt and alt != filename and download_file(slug, alt, dest):
                print(f"  ✓ [{label}] {slug}.jpg ← {alt} (search fallback)")
                manifest[slug] = alt
                ok += 1
                time.sleep(0.8)
                continue

        print(f"  ✗ [{label}] {slug}: failed ({filename})")
        failed.append(slug)
        time.sleep(0.8)

    if failed:
        save_json(MANIFEST_FILE, manifest)
        print(f"\nFailed ({len(failed)}): {', '.join(failed)}")
    return ok


def verify_manifest(manifest: dict[str, str]) -> None:
    bad = []
    for slug, filename in manifest.items():
        if not verify_file_exists(filename):
            bad.append((slug, filename))
    if bad:
        print(f"Invalid manifest entries: {len(bad)}")
        for slug, fn in bad[:20]:
            print(f"  - {slug}: {fn}")
    else:
        print(f"All {len(manifest)} manifest entries verified.")


def main() -> int:
    parser = argparse.ArgumentParser(description="Fetch product images from Wikimedia Commons")
    parser.add_argument("--resolve", action="store_true", help="Rebuild manifest via Wikimedia API")
    parser.add_argument("--verify", action="store_true", help="Verify manifest URLs exist")
    parser.add_argument("--download-only", action="store_true", help="Download using existing manifest")
    args = parser.parse_args()

    if args.resolve:
        searches = load_json(SEARCHES_FILE)
        manifest = resolve_manifest(searches)
        save_json(MANIFEST_FILE, manifest)
        print(f"\nSaved manifest: {len(manifest)} entries → {MANIFEST_FILE}")
        return 0

    manifest = build_manifest(force_resolve=False)

    if args.verify:
        verify_manifest(manifest)
        return 0

    print(f"\nDownloading {len(manifest)} images...")
    ok = download_all(manifest)
    print(f"\nDone: {ok}/{len(manifest)} images downloaded")
    return 0 if ok == len(manifest) else 1


if __name__ == "__main__":
    sys.exit(main())
