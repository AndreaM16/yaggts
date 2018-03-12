import argparse
import sys
import itertools
import os
from selenium import webdriver
from selenium.common.exceptions import NoSuchElementException
from selenium.webdriver.support.ui import Select
from time import sleep

TRENDS_URL = 'https://trends.google.it/trends/explore?'
QUERY = 'date=2015-01-01%202017-10-07&geo=GB&q='

def findAndClickByCssSelector(driver, classes):
    selector = '.'.join(classes)    
    try:
        btn = driver.find_element_by_css_selector(classes)
    except NoSuchElementException:
        return
    btn.click()
    return
filepath = os.path.join('/home/andream16/go/src/github.com/andream16/yaggts', 'multiTimeline.csv')
if os.path.exists(filepath):
    os.remove(filepath)
profile = webdriver.FirefoxProfile()
profile.set_preference('browser.download.folderList', 2)
profile.set_preference('browser.download.manager.showWhenStarting', True)
profile.set_preference('browser.download.dir', '/home/andream16/go/src/github.com/andream16/yaggts')
profile.set_preference('browser.helperApps.neverAsk.saveToDisk', 'text/csv')
driver = webdriver.Firefox(profile)
driver.get(TRENDS_URL + QUERY + sys.argv[1])
sleep(5)
findAndClickByCssSelector(driver, ".widget-actions-item-flatten .widget-actions-item.export")
driver.close()

quit()
