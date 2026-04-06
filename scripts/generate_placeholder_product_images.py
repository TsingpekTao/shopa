from __future__ import annotations

import csv
import random
import re
from dataclasses import dataclass
from pathlib import Path

from PIL import Image, ImageDraw, ImageFilter, ImageFont


OUTPUT_ROOT = Path("C:/Users/Tao/Pictures/\u5546\u54c1")
IMAGE_SIZE = (1024, 1024)
PER_CATEGORY = 50
SEED = 20260403


@dataclass(frozen=True)
class ProductSpec:
    category: str
    title: str
    accent: tuple[int, int, int]
    style: str


def u(text: str) -> str:
    return text.encode("utf-8").decode("unicode_escape")


def load_font(size: int, bold: bool = False) -> ImageFont.FreeTypeFont | ImageFont.ImageFont:
    candidates: list[str] = []
    if bold:
        candidates.extend(
            [
                r"C:\Windows\Fonts\msyhbd.ttc",
                r"C:\Windows\Fonts\simhei.ttf",
                r"C:\Windows\Fonts\arialbd.ttf",
            ]
        )
    candidates.extend(
        [
            r"C:\Windows\Fonts\msyh.ttc",
            r"C:\Windows\Fonts\simsun.ttc",
            r"C:\Windows\Fonts\arial.ttf",
        ]
    )
    for path in candidates:
        try:
            if Path(path).exists():
                return ImageFont.truetype(path, size=size)
        except OSError:
            continue
    return ImageFont.load_default()


FONT_TITLE = load_font(42, bold=True)
FONT_META = load_font(24)
FONT_BADGE = load_font(28, bold=True)


def zh_items(items: list[str]) -> list[str]:
    return [u(item) for item in items]


def make_specs() -> list[ProductSpec]:
    rng = random.Random(SEED)
    categories = {
        u("\\u670d\\u88c5"): {
            "prefixes": zh_items(
                [
                    "\\u8f7b\\u901a\\u52e4",
                    "\\u7b80\\u7ea6",
                    "\\u590d\\u53e4",
                    "\\u57ce\\u5e02",
                    "\\u8212\\u9002",
                    "\\u9ad8\\u7ea7\\u611f",
                ]
            ),
            "colors": [
                (u("\\u5976\\u6cb9\\u767d"), (235, 226, 210)),
                (u("\\u96fe\\u973e\\u84dd"), (131, 149, 174)),
                (u("\\u7126\\u7cd6\\u68d5"), (163, 105, 62)),
                (u("\\u6a44\\u6984\\u7eff"), (104, 122, 86)),
                (u("\\u6a31\\u82b1\\u7c89"), (216, 166, 177)),
            ],
            "kinds": zh_items(
                [
                    "\\u5706\\u9886T\\u6064",
                    "\\u8fde\\u5e3d\\u536b\\u8863",
                    "\\u9488\\u7ec7\\u5f00\\u886b",
                    "\\u76f4\\u7b52\\u725b\\u4ed4\\u88e4",
                    "\\u534a\\u8eab\\u957f\\u88d9",
                    "\\u8f7b\\u7fbd\\u7ed2\\u670d",
                    "\\u5de5\\u88c5\\u77ed\\u88e4",
                    "\\u886c\\u886b\\u5916\\u5957",
                ]
            ),
        },
        u("\\u978b\\u5b50"): {
            "prefixes": zh_items(
                [
                    "\\u8f7b\\u76c8",
                    "\\u7ecf\\u5178",
                    "\\u6f6e\\u6d41",
                    "\\u57ce\\u5e02",
                    "\\u900f\\u6c14",
                    "\\u673a\\u80fd",
                ]
            ),
            "colors": [
                (u("\\u4e91\\u767d"), (240, 239, 234)),
                (u("\\u66dc\\u9ed1"), (74, 76, 84)),
                (u("\\u5976\\u8336\\u68d5"), (166, 128, 101)),
                (u("\\u7535\\u5149\\u84dd"), (90, 132, 206)),
                (u("\\u73ca\\u745a\\u6a59"), (205, 120, 78)),
            ],
            "kinds": zh_items(
                [
                    "\\u8001\\u7238\\u8fd0\\u52a8\\u978b",
                    "\\u8f7b\\u8dd1\\u8bad\\u7ec3\\u978b",
                    "\\u9ad8\\u5e2e\\u5e06\\u5e03\\u978b",
                    "\\u539a\\u5e95\\u4f11\\u95f2\\u978b",
                    "\\u5546\\u52a1\\u4e50\\u798f\\u978b",
                    "\\u6237\\u5916\\u5f92\\u6b65\\u978b",
                    "\\u900f\\u6c14\\u51c9\\u978b",
                    "\\u8f6f\\u5e95\\u62d6\\u978b",
                ]
            ),
        },
        u("\\u8033\\u673a"): {
            "prefixes": zh_items(
                [
                    "HiFi",
                    "\\u964d\\u566a",
                    "\\u6e38\\u620f",
                    "\\u65e0\\u7ebf",
                    "\\u957f\\u7eed\\u822a",
                    "\\u5546\\u52a1",
                ]
            ),
            "colors": [
                (u("\\u51b0\\u5ddd\\u767d"), (237, 239, 244)),
                (u("\\u66dc\\u77f3\\u9ed1"), (56, 58, 65)),
                (u("\\u661f\\u7a7a\\u84dd"), (67, 97, 164)),
                (u("\\u6696\\u91d1"), (198, 160, 95)),
                (u("\\u8584\\u8377\\u7eff"), (149, 183, 170)),
            ],
            "kinds": zh_items(
                [
                    "\\u84dd\\u7259\\u8033\\u673a",
                    "\\u5934\\u6234\\u8033\\u673a",
                    "\\u771f\\u65e0\\u7ebf\\u8033\\u585e",
                    "\\u6e38\\u620f\\u8033\\u9ea6",
                    "\\u534a\\u5165\\u8033\\u8033\\u673a",
                    "\\u7ffb\\u8bd1\\u8033\\u673a",
                    "\\u8fd0\\u52a8\\u6302\\u8033\\u8033\\u673a",
                    "\\u901a\\u52e4\\u964d\\u566a\\u8033\\u673a",
                ]
            ),
        },
        u("\\u5bb6\\u5c45"): {
            "prefixes": zh_items(
                [
                    "\\u5317\\u6b27",
                    "\\u6cbb\\u6108",
                    "\\u7b80\\u7ea6",
                    "\\u6696\\u5c45",
                    "\\u539f\\u6728",
                    "\\u751f\\u6d3b\\u611f",
                ]
            ),
            "colors": [
                (u("\\u71d5\\u9ea6\\u767d"), (232, 224, 210)),
                (u("\\u96fe\\u84dd\\u7070"), (131, 148, 160)),
                (u("\\u539f\\u6728\\u68d5"), (162, 125, 87)),
                (u("\\u53ef\\u53ef\\u68d5"), (124, 90, 67)),
                (u("\\u5976\\u674f\\u8272"), (228, 200, 172)),
            ],
            "kinds": zh_items(
                [
                    "\\u53f0\\u706f",
                    "\\u62b1\\u6795",
                    "\\u9999\\u85b0\\u673a",
                    "\\u8fb9\\u51e0",
                    "\\u9a6c\\u514b\\u676f",
                    "\\u50a8\\u7269\\u7bb1",
                    "\\u82b1\\u74f6",
                    "\\u4f11\\u95f2\\u5355\\u6905",
                ]
            ),
        },
    }

    specs: list[ProductSpec] = []
    for category, config in categories.items():
        combos: list[ProductSpec] = []
        for prefix in config["prefixes"]:
            for color_name, accent in config["colors"]:
                for kind in config["kinds"]:
                    combos.append(ProductSpec(category, f"{prefix}{color_name}{kind}", accent, kind))
        rng.shuffle(combos)
        specs.extend(combos[:PER_CATEGORY])
    return specs


def safe_stem(text: str) -> str:
    text = re.sub(r"\s+", "_", text.strip())
    text = re.sub(r'[\\/:*?"<>|]+', "_", text)
    return text[:48]


def gradient_background(base: tuple[int, int, int]) -> Image.Image:
    # Draw on a small canvas first and upscale, which is much faster than
    # painting every pixel on a 1024x1024 image 200 times.
    width, height = 96, 96
    image = Image.new("RGB", (width, height), "white")
    pixels = image.load()
    for y in range(height):
        for x in range(width):
            rx = x / max(width - 1, 1)
            ry = y / max(height - 1, 1)
            mix = 0.35 * rx + 0.65 * ry
            r = int(247 + (base[0] - 247) * mix * 0.34)
            g = int(244 + (base[1] - 244) * mix * 0.31)
            b = int(241 + (base[2] - 241) * mix * 0.28)
            pixels[x, y] = (r, g, b)
    return image.resize(IMAGE_SIZE, Image.Resampling.BICUBIC).convert("RGBA")


def add_shadow(image: Image.Image) -> None:
    shadow = Image.new("RGBA", image.size, (0, 0, 0, 0))
    shadow_draw = ImageDraw.Draw(shadow)
    shadow_draw.rounded_rectangle((66, 66, 958, 958), radius=48, fill=(86, 56, 28, 65))
    image.alpha_composite(shadow.filter(ImageFilter.GaussianBlur(radius=24)))


def draw_stage(draw: ImageDraw.ImageDraw) -> None:
    draw.rounded_rectangle((64, 64, 960, 960), radius=48, fill=(255, 251, 246))
    draw.ellipse((700, 116, 930, 346), fill=(249, 233, 220))
    draw.ellipse((678, 182, 926, 430), fill=(246, 239, 231))
    draw.rounded_rectangle((136, 650, 888, 718), radius=34, fill=(243, 232, 221))


def draw_badge(draw: ImageDraw.ImageDraw, text: str, accent: tuple[int, int, int]) -> None:
    text_box = draw.textbbox((0, 0), text, font=FONT_BADGE)
    width = text_box[2] - text_box[0] + 34
    badge = tuple(max(0, item - 18) for item in accent)
    draw.rounded_rectangle((74, 72, 74 + width, 120), radius=22, fill=badge)
    draw.text((92, 82), text, font=FONT_BADGE, fill=(255, 255, 255))


def draw_meta(draw: ImageDraw.ImageDraw, spec: ProductSpec) -> None:
    draw.text((86, 772), spec.title, font=FONT_TITLE, fill=(59, 39, 28))
    draw.text((88, 834), f"{spec.category} | {spec.style}", font=FONT_META, fill=(113, 86, 67))
    draw.text((88, 872), u("\\u5546\\u57ce\\u6f14\\u793a\\u7528\\u5360\\u4f4d\\u5546\\u54c1\\u56fe"), font=FONT_META, fill=(136, 110, 88))
    draw.rounded_rectangle((86, 918, 312, 972), radius=22, fill=(250, 244, 238))
    draw.text((106, 934), "AI Placeholder", font=FONT_META, fill=(126, 97, 74))


def draw_clothing(draw: ImageDraw.ImageDraw, spec: ProductSpec) -> None:
    base = spec.accent
    dark = tuple(max(0, c - 28) for c in base)
    light = tuple(min(255, c + 22) for c in base)
    if u("\\u88e4") in spec.style:
        left = [(360, 230), (470, 230), (504, 646), (432, 662), (396, 412), (378, 666), (308, 654), (334, 310)]
        right = [(554, 230), (664, 230), (690, 310), (716, 654), (646, 666), (628, 412), (592, 662), (520, 646)]
        draw.polygon(left, fill=base)
        draw.polygon(right, fill=base)
        draw.line((512, 232, 512, 646), fill=dark, width=6)
        draw.line((352, 320, 412, 602), fill=light, width=4)
        draw.line((672, 320, 612, 602), fill=light, width=4)
    elif u("\\u88d9") in spec.style:
        draw.polygon([(302, 238), (722, 238), (804, 648), (220, 648)], fill=base)
        draw.line((338, 284, 760, 284), fill=dark, width=8)
        for x in range(280, 744, 56):
            draw.line((x, 252, x + 70, 646), fill=light, width=4)
    else:
        torso = [(368, 212), (656, 212), (710, 612), (314, 612)]
        left = [(314, 254), (214, 338), (252, 606), (346, 572)]
        right = [(710, 254), (810, 338), (772, 606), (678, 572)]
        draw.polygon(torso, fill=base)
        draw.polygon(left, fill=base)
        draw.polygon(right, fill=base)
        if u("\\u536b\\u8863") in spec.style or u("\\u8fde\\u5e3d") in spec.style:
            draw.ellipse((428, 156, 596, 286), fill=dark)
            draw.ellipse((458, 180, 566, 266), fill=(255, 248, 240))
        else:
            draw.arc((430, 144, 594, 292), start=22, end=158, fill=(248, 242, 234), width=18)
        draw.line((512, 286, 512, 590), fill=light, width=5)
        draw.line((384, 350, 640, 350), fill=dark, width=4)
    draw.ellipse((328, 668, 696, 728), fill=(220, 208, 195))


def draw_shoes(draw: ImageDraw.ImageDraw, spec: ProductSpec) -> None:
    base = spec.accent
    dark = tuple(max(0, c - 36) for c in base)
    sole = (232, 231, 232)
    left = [(272, 530), (350, 470), (494, 450), (598, 494), (662, 560), (654, 608), (372, 630), (286, 594)]
    right = [(374, 648), (448, 592), (596, 574), (700, 618), (750, 676), (738, 718), (472, 736), (386, 706)]
    draw.polygon(left, fill=base)
    draw.polygon(right, fill=tuple(min(255, c + 12) for c in base))
    draw.polygon([(280, 588), (664, 572), (662, 622), (280, 618)], fill=sole)
    draw.polygon([(382, 702), (744, 684), (742, 730), (382, 736)], fill=sole)
    for y in (518, 544, 570):
        draw.line((410, y, 560, y - 10), fill=dark, width=6)
    for y in (638, 664, 688):
        draw.line((518, y, 672, y - 12), fill=dark, width=6)


def draw_earphones(draw: ImageDraw.ImageDraw, spec: ProductSpec) -> None:
    base = spec.accent
    dark = tuple(max(0, c - 32) for c in base)
    metal = (226, 228, 232)
    if u("\\u5934\\u6234") in spec.style or u("\\u8033\\u9ea6") in spec.style:
        draw.arc((292, 164, 732, 560), start=205, end=335, fill=dark, width=34)
        draw.rounded_rectangle((242, 382, 408, 652), radius=70, fill=base)
        draw.rounded_rectangle((618, 382, 784, 652), radius=70, fill=base)
        draw.rounded_rectangle((286, 430, 364, 604), radius=32, fill=(246, 240, 232))
        draw.rounded_rectangle((662, 430, 740, 604), radius=32, fill=(246, 240, 232))
        draw.rounded_rectangle((426, 646, 600, 706), radius=30, fill=metal)
    else:
        draw.rounded_rectangle((332, 288, 694, 684), radius=84, fill=(250, 246, 241), outline=(224, 216, 208), width=8)
        draw.rounded_rectangle((382, 350, 642, 620), radius=56, fill=base)
        draw.rounded_rectangle((252, 308, 402, 562), radius=70, fill=base)
        draw.rounded_rectangle((622, 308, 772, 562), radius=70, fill=base)
        draw.ellipse((298, 346, 356, 404), fill=dark)
        draw.ellipse((668, 346, 726, 404), fill=dark)
        draw.rounded_rectangle((450, 692, 574, 722), radius=14, fill=(224, 216, 208))
    draw.ellipse((338, 728, 688, 786), fill=(222, 211, 199))


def draw_home(draw: ImageDraw.ImageDraw, spec: ProductSpec) -> None:
    base = spec.accent
    dark = tuple(max(0, c - 34) for c in base)
    light = tuple(min(255, c + 22) for c in base)
    if u("\\u53f0\\u706f") in spec.style:
        draw.ellipse((422, 650, 602, 708), fill=(221, 210, 198))
        draw.rectangle((494, 286, 530, 652), fill=dark)
        draw.polygon([(364, 300), (660, 300), (584, 502), (440, 502)], fill=base)
        draw.ellipse((462, 522, 560, 612), fill=light)
    elif u("\\u62b1\\u6795") in spec.style:
        draw.rounded_rectangle((288, 278, 736, 710), radius=92, fill=base)
        draw.rounded_rectangle((344, 334, 680, 654), radius=76, fill=light)
        draw.line((362, 360, 668, 640), fill=dark, width=8)
        draw.line((662, 360, 360, 650), fill=dark, width=8)
    elif u("\\u9a6c\\u514b\\u676f") in spec.style:
        draw.ellipse((342, 662, 704, 722), fill=(220, 210, 198))
        draw.rounded_rectangle((372, 314, 636, 662), radius=56, fill=base)
        draw.arc((590, 394, 760, 580), start=260, end=110, fill=dark, width=20)
        draw.rectangle((444, 278, 564, 336), fill=light)
    elif u("\\u8fb9\\u51e0") in spec.style:
        draw.ellipse((306, 286, 718, 406), fill=base)
        for x in (378, 470, 574, 654):
            draw.rectangle((x, 404, x + 18, 718), fill=dark)
        draw.ellipse((300, 706, 722, 760), fill=(221, 210, 198))
    elif u("\\u5355\\u6905") in spec.style:
        draw.rounded_rectangle((334, 300, 690, 566), radius=74, fill=base)
        draw.rounded_rectangle((376, 520, 650, 654), radius=50, fill=light)
        for line in [(404, 650, 382, 784), (610, 650, 632, 784), (462, 650, 450, 784), (548, 650, 560, 784)]:
            draw.line(line, fill=dark, width=12)
    else:
        draw.rounded_rectangle((350, 250, 676, 710), radius=72, fill=base)
        draw.ellipse((416, 202, 612, 292), fill=light)
        draw.ellipse((378, 680, 642, 736), fill=(223, 212, 200))
        draw.line((512, 300, 512, 640), fill=dark, width=8)


def draw_product(spec: ProductSpec, destination: Path) -> None:
    image = gradient_background(spec.accent)
    add_shadow(image)
    draw = ImageDraw.Draw(image)
    draw_stage(draw)
    draw_badge(draw, spec.category, spec.accent)
    if spec.category == u("\\u670d\\u88c5"):
        draw_clothing(draw, spec)
    elif spec.category == u("\\u978b\\u5b50"):
        draw_shoes(draw, spec)
    elif spec.category == u("\\u8033\\u673a"):
        draw_earphones(draw, spec)
    else:
        draw_home(draw, spec)
    draw_meta(draw, spec)
    image.convert("RGB").save(destination, quality=95)


def main() -> None:
    specs = make_specs()
    OUTPUT_ROOT.mkdir(parents=True, exist_ok=True)
    category_dirs: dict[str, Path] = {}
    for category in sorted({spec.category for spec in specs}):
        category_dir = OUTPUT_ROOT / category
        category_dir.mkdir(parents=True, exist_ok=True)
        category_dirs[category] = category_dir

    manifest_path = OUTPUT_ROOT / "manifest.csv"
    counters = {category: 0 for category in category_dirs}
    with manifest_path.open("w", encoding="utf-8-sig", newline="") as file:
        writer = csv.writer(file)
        writer.writerow(["category", "title", "file_path"])
        for spec in specs:
            counters[spec.category] += 1
            filename = f"{counters[spec.category]:03d}_{safe_stem(spec.title)}.png"
            output_file = category_dirs[spec.category] / filename
            draw_product(spec, output_file)
            writer.writerow([spec.category, spec.title, str(output_file)])

    print(f"generated {len(specs)} images into {OUTPUT_ROOT}")


if __name__ == "__main__":
    main()
