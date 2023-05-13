from PIL import Image
import shutil
import uuid
import os
import sys
import re

def get_files(path):
    """指定したディレクトリ内のファイルを取得する

    Args:
        path (str): ディレクトリのパス

    Returns:
        list: ファイルのリスト
    """
    files = []
    for filename in os.listdir(path):
        file_path = os.path.join(path, filename)

        # 拡張子がjpg, jpeg, pngのファイルのみカウント
        if not filename.lower().endswith(('.jpg', '.jpeg')):
            continue

        if not os.path.isfile(file_path):
            continue

        files.append(filename)

    sorted_files = sorted(files, key=lambda x: int(re.search(r"\d+", x).group()))

    return sorted_files

def combine_images(images, width, height):
    if width < height: # 縦長の場合
        total_width = width * 10 - 45  # 合成後の画像の幅（5ピクセル × 9）

        # 新しいキャンバスの作成
        result_image = Image.new("RGB", (total_width, height))

        # 画像の結合
        for j in range(10):
            offset = j * (width - 5)  # 合成時のオフセット位置
            result_image.paste(images[j], (offset, 0))
    else: # 横長の場合
        total_height = height * 10 - 45  # 合成後の画像の高さ（5ピクセル × 9）

        # 新しいキャンバスの作成
        result_image = Image.new("RGB", (width, total_height))

        # 画像の結合
        for j in range(10):
            offset = j * (height - 5)  # 合成時のオフセット位置
            result_image.paste(images[j], (0, offset))

    # 合成画像の保存
    time_uuid = uuid.uuid1()
    result_image.save(f"{output_dir}/{time_uuid}.jpg")
    print(f"saved {output_dir}/{time_uuid}.jpg")
    return time_uuid

# 分割された画像が入っているディレクトリのパス
input_dir = "./source"

# 結合した画像を保存するディレクトリのパス
output_dir = "./combined"

# 分割された画像を移動するディレクトリのパス
archive_dir = "./archive"

# input_dirの画像ファイルの名前を変更
if os.path.exists(f"{input_dir}/ダウンロード.jpeg"):
    shutil.move(f"{input_dir}/ダウンロード.jpeg", f"{input_dir}/ダウンロード (0).jpeg")

# input_dir配下のファイルを取得
file_list = get_files(input_dir)

if len(file_list) == 0:
    print("画像がありません。")
    sys.exit()

if len(file_list)%10 != 0:
    print(f"画像が {len(file_list)} 枚あります。10枚単位で合成します。")
    sys.exit()

for i in range(int(len(file_list)/10)):
    # 10枚の画像を読み込む
    images = []

    width = 0 # 合成後の画像の幅
    height = 0 # 合成後の画像の高さ

    for j in range(0, 10):
        image_path = f"{input_dir}/{file_list[i*10+j]}"  # 画像ファイルのパス

        print(f"reading {image_path}")

        if not file_list[i*10+j].lower().endswith(('.jpg', '.jpeg')):
            continue
        image = Image.open(image_path) # 画像を開く
        image_width, image_height = image.size # 1枚の画像の幅を計算

        if j == 0:
            width = image_width
            height = image_height

        if (width < height) != (image_width < image_height):
            print(f"{file_list[i*10+j]} は {file_list[i*10]} と画像サイズが異なります。")
            sys.exit()

        images.append(image)

    # 画像の結合
    print(f"combining {file_list[i*10]} ~ {file_list[i*10+9]}")
    time_uuid = combine_images(images, width, height)
    os.mkdir(f"{archive_dir}/{time_uuid}")

    # 結合を終えた画像を削除
    for j in range(0, 10):
        image_path = f"{input_dir}/{file_list[i*10+j]}"
        shutil.move(image_path, f"{archive_dir}/{time_uuid}/{file_list[i*10+j]}.jpg")
