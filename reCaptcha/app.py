import os
import sys

from selenium import webdriver
from selenium.webdriver.chrome.options import Options

HEADERS = {
    "User-Agent": "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/95.0.4638.69 Safari/537.36"
}


def main(url: str, client_key: str, action: str) -> str:
    if not url:
        return "MISSING ARGUMENTS"

    options = Options()

    options.add_argument("--headless")
    options.add_argument("--disable-gpu")
    options.add_argument("--no-sandbox")
    options.add_argument('--disable-dev-shm-usage')
    options.add_argument('--disable-gpu-sandbox')
    options.add_argument("--single-process")
    options.add_argument('window-size=1920x1080')
    options.add_argument(
        'user-agent=Mozilla/5.0 (Macintosh; Intel Mac OS X 10_14_6) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/76.0.3809.100 Safari/537.36')

    options.binary_location = os.getcwd() + "/bin/headless-chromium"

    browser = webdriver.Chrome(
        os.getcwd() + "/bin/chromedriver", options=options)

    browser.get(url)

    browser.execute_async_script("""
    window.tokenVal = "";
    var done = arguments[0];
    grecaptcha.execute('""" + client_key + """', { action: '""" + action + """' }).then((token) => {window.tokenVal = token; done(token)});
    """)

    result = browser.execute_script("return window.tokenVal")

    return result


if __name__ == "__main__":
    print(main(sys.argv[1], sys.argv[2], sys.argv[3]))
