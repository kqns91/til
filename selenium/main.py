from dotenv import load_dotenv
import os
from selenium import webdriver
import time
from selenium.webdriver.common.action_chains import ActionChains

# .env ファイルを読み込む
load_dotenv()

# 環境変数の利用
browse_key = os.getenv("BROWSE_KEY")
user_name = os.getenv("USER_NAME")
password = os.getenv("PASSWORD")

# ブラウザを開く
driver = webdriver.Chrome('./chromedriver_mac_arm64/chromedriver')

# Anicritにアクセス
driver.get('https://www.anicrit.jp/photo')

# 閲覧キーを入力
driver.find_element_by_xpath('//*[@id="login_form"]/input[1]').send_keys(browse_key)
driver.find_element_by_xpath('//*[@id="login_form"]/button').click()

# ログイン画面に遷移
driver.find_element_by_xpath('//*[@id="register"]/div/a/img').click()

# ログイン
driver.find_element_by_xpath('//*[@id="login_id"]').send_keys(user_name)
driver.find_element_by_xpath('//*[@id="password"]').send_keys(password)
driver.find_element_by_xpath('//*[@id="login"]').click()

# 写真一覧ページを新規タブで開く
element = driver.find_element_by_xpath('//*[@id="aftop"]/section[1]/div/div[2]/div/a')
driver.execute_script("arguments[0].click();", element)

# 新規タブに移動
driver.switch_to.window(driver.window_handles[-1])

# リロード
driver.refresh()

time.sleep(1)

img = driver.find_element_by_xpath('//*[@id="front_group_detail_group_id-1-0"]/div[1]/a')
driver.execute_script("arguments[0].click();", img)

# 入力を待機
input('Enterで終了')

#ブラウザを閉じる
driver.quit()

