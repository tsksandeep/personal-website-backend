import os
import sys
import requests

from bs4 import BeautifulSoup
from selenium import webdriver
from selenium.webdriver.chrome.options import Options

HEADERS = {
    "User-Agent": "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/95.0.4638.69 Safari/537.36"
}

URL = "https://www.google.com/recaptcha/api.js?render="


def main(client_key: str, action: str) -> str:
    if not (client_key and action):
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
    browser.implicitly_wait(3)

    recaptcha_code = requests.get(URL+client_key, headers=HEADERS).text
    recaptcha_code = recaptcha_code.strip(
        "/* PLEASE DO NOT COPY AND PASTE THIS CODE. */")

    RECAPTCHA_HTML = """
    <html>
        <body>
            <div id="recaptchaV3Token"></div>
            <script>
            """ + recaptcha_code + """
            function generateToken() {
                grecaptcha.ready(function () {
                grecaptcha
                    .execute('""" + client_key + """', {
                    action: '""" + action + """',})
                    .then(function (e) {
                    document.getElementById("recaptchaV3Token").innerText = e;
                    });
                });
            }
            generateToken();
            </script>
        </body>
    </html>
    """

    browser.get(f"data:text/html;charset=utf-8,{RECAPTCHA_HTML}")
    soup = BeautifulSoup(browser.page_source, 'lxml')
    browser.quit()

    return soup.find("div", {"id": "recaptchaV3Token"}).text.strip()


if __name__ == "__main__":
    print(main(sys.argv[1], sys.argv[2]))
