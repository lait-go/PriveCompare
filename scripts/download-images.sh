#!/bin/bash
# Скачивает тематические изображения товаров (Pexels, проверенные ID)
set -e
SCRIPT_DIR="$(cd "$(dirname "$0")" && pwd)"
if [ -d "$SCRIPT_DIR/../frontend/public/images" ]; then
  DIR="$SCRIPT_DIR/../frontend/public/images"
elif [ -d "$SCRIPT_DIR/../public/images" ]; then
  DIR="$SCRIPT_DIR/../public/images"
else
  DIR="${IMAGES_DIR:-$SCRIPT_DIR/../frontend/public/images}"
fi
PRODUCTS="$DIR/products"
CATEGORIES="$DIR/categories"
mkdir -p "$PRODUCTS" "$CATEGORIES"

dl() {
  local file="$1" id="$2"
  curl -sL "https://images.pexels.com/photos/${id}/pexels-photo-${id}.jpeg?auto=compress&cs=tinysrgb&w=400&h=400&fit=crop" \
    -o "${PRODUCTS}/${file}.jpg"
  echo "  ✓ ${file}.jpg"
}

dl_cat() {
  local file="$1" id="$2"
  curl -sL "https://images.pexels.com/photos/${id}/pexels-photo-${id}.jpeg?auto=compress&cs=tinysrgb&w=400&h=300&fit=crop" \
    -o "${CATEGORIES}/${file}.jpg"
  echo "  ✓ categories/${file}.jpg"
}

echo "Категории..."
dl_cat food 264636
dl_cat household 6474593
dl_cat construction 259861
dl_cat chemistry 421032
dl_cat drinks 1556688
dl_cat home 1571460
dl_cat electronics 788946
dl_cat other 112898

echo "Продукты питания..."
dl milk 773253
dl bread 1775043
dl eggs 209206
dl cheese 821003
dl chicken 106139
dl buckwheat 4110255
dl rice 723198
dl pasta 1279330
dl sugar 6580225
dl salt 277776
dl sunflower-oil 189355
dl kefir 1628088
dl cottage-cheese 8962414
dl yogurt 3735663
dl sausage 4198032
dl hotdog 1639563
dl potato 247117
dl carrot 143133
dl onion 4198092
dl tomato 533280
dl cucumber 4198093
dl apple 102104
dl banana 61127
dl orange 327098
dl chocolate 4198094
dl cookies 230325
dl ketchup 4198095
dl mayonnaise 4198096
dl peas 4198097
dl tuna 248444

echo "Напитки..."
dl coca-cola 1556688
dl sprite 4198098
dl fanta 4198099
dl juice 20204662
dl water 416528
dl tea 230477
dl coffee 894695
dl beer 1267328
dl wine 340013
dl energy-drink 2884115

echo "Бытовая химия..."
dl laundry-powder 421032
dl laundry-gel 3993446
dl dish-soap 4198100
dl floor-cleaner 4198101
dl bathroom-cleaner 4198102
dl shampoo 4198103
dl shower-gel 4198104
dl toothpaste 4198105
dl toilet-paper 582492
dl paper-towels 582494

echo "Хозяйственные..."
dl sponge 4198106
dl rubber-gloves 4198107
dl bucket 4198108
dl mop 4198111
dl trash-bags 4198112

echo "Строительные..."
dl cement 259861
dl sand 259862
dl brick 259863
dl plaster 259864
dl paint 259865
dl wallpaper 259866
dl laminate 259867
dl drill 259868
dl hammer 259869
dl screwdriver 259870

echo "Для дома..."
dl lamp 1571460
dl pillow 1571453
dl blanket 1571458
dl carpet 4198113
dl vase 4198114
dl photo-frame 4198115
dl wall-clock 4198116
dl curtains 4198117

echo "Электроника..."
dl smartphone 788946
dl headphones 2983101
dl tablet 4041392
dl tv 4042804
dl laptop 4042802
dl microwave 4042803
dl fridge 4042805
dl washing-machine 6474535
dl vacuum 6474536
dl iron 6474537
dl hairdryer 6474542
dl blender 6474543
dl coffee-machine 6474544
dl electric-toothbrush 6474551
dl powerbank 6474552
dl usb-cable 6474553
dl keyboard 6474554
dl computer-mouse 6474555
dl webcam 6474556
dl bluetooth-speaker 6474557
dl wifi-router 6474564
dl ssd 6474565
dl flash-drive 6474566
dl memory-card 6474567
dl charger 6474568
dl phone-case 6474569
dl screen-protector 6474573
dl batteries 6474574

echo "Канцтовары..."
dl tape 6474581
dl glue 112898
dl notebook 112899
dl pen 112900
dl pencils 112901
dl scissors 112902
dl ruler 3992938
dl paperclips 3992939
dl stapler 3992940
dl hole-punch 3992942
dl folder 3992943
dl correction-fluid 3992944
dl markers 3992945
dl glue-stick 3992946
dl eraser 3992947
dl pencil-sharpener 3992948
dl notebook-covers 3992949
dl exercise-book 3992950
dl diary 3992951
dl drawing-album 3992952
dl watercolors 3992953
dl paint-brushes 3992955
dl play-dough 3992956

echo "Готово: $(ls -1 "$PRODUCTS" | wc -l | tr -d ' ') изображений товаров"
