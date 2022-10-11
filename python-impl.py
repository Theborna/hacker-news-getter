from datetime import datetime
import os
import sys
from selenium import webdriver
from selenium.webdriver.chrome.service import Service
from selenium.webdriver.chrome.options import Options
import pandas as pd

PATH = 'C:/Users/hp/chrome-driver/chromedriver.exe'
BASE_PATH = os.path.dirname(sys.executable)


def get_df(driver: webdriver.Chrome) -> pd.DataFrame:
    containers = driver.find_elements(
        by="xpath", value='//td[@class="title"]/span/a')
    titles = []
    # subtitles = []
    links = []
    for container in containers:
        title = container.text
        # subtitle = container.find_element(by='xpath', value='./a/p').text
        link = container.get_attribute('href')
        titles.append(title)
        # subtitles.append(subtitle)
        links.append(link)

    return pd.DataFrame({
        'title': titles,
        # 'subtitle': subtitles,
        'link': links,
    })


def main():
    service = Service(executable_path=PATH)
    option = Options()
    option.headless = True
    driver = webdriver.Chrome(service=service, options=option)
    website = "https://news.ycombinator.com"
    sites = ['newest','front','jobs']
    data_frames = []
    for site in sites:
        print(f'{website}/{site}')
        driver.get(f'{website}/{site}')
        headlines = get_df(driver)
        headlines['type'] = site
        data_frames.append(headlines)
        print(headlines.head())
    
    headlines = pd.concat(data_frames)
    now = datetime.now()
    time = now.strftime("%B-%d-%Y")
    path = os.path.join(BASE_PATH,f'hacker-news_{time}.csv')
    headlines.to_csv(path)
    
    driver.quit()


if __name__ == '__main__':
    main()
