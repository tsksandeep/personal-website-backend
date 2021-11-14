import os
import sys
import time

from selenium import webdriver
from selenium.webdriver.chrome.options import Options

HEADERS = {
    "User-Agent": "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/95.0.4638.69 Safari/537.36"
}


def main(url: str, client_key: str, action: str) -> str:
    options = Options()

    options.add_argument("--headless")
    options.add_argument("--disable-gpu")
    options.add_argument("--no-sandbox")
    options.add_argument('--disable-dev-shm-usage')
    options.add_argument('--disable-gpu-sandbox')
    options.add_argument("--single-process")
    options.add_argument("start-maximized")
    options.add_experimental_option("excludeSwitches", ["enable-automation"])
    options.add_experimental_option('useAutomationExtension', False)
    options.add_argument(
        'user-agent=Mozilla/5.0 (Macintosh; Intel Mac OS X 10_14_6) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/76.0.3809.100 Safari/537.36')

    options.binary_location = os.getcwd() + "/bin/headless-chromium"

    browser = webdriver.Chrome(
        os.getcwd() + "/bin/chromedriver", options=options)

    browser.get(url)

    browser.execute_script("""
    (function () {
        var w = window,
            C = '___grecaptcha_cfg',
            cfg = (w[C] = w[C] || {}),
            N = 'grecaptcha';
        var gr = (w[N] = w[N] || {});
        gr.ready =
            gr.ready ||
            function (f) {
            (cfg['fns'] = cfg['fns'] || []).push(f);
            };
        w['__recaptcha_api'] = 'https://www.google.com/recaptcha/api2/';
        (cfg['render'] = cfg['render'] || []).push(
            '""" + client_key + """'
        );
        w['__google_recaptcha_client'] = true;
        var d = document,
            po = d.createElement('script');
        po.type = 'text/javascript';
        po.async = true;
        po.src =
            'https://www.gstatic.com/recaptcha/releases/yZguKF1TiDm6F3yJWVhmOKQ9/recaptcha__en.js';
        po.crossOrigin = 'anonymous';
        po.integrity =
            'sha384-H5arsq6vo+zXBZN6NglbXmqmJ0Fd0I03R0T+1QSVXdReK5pp2+U7XIsjY638NZyL';

        var e = d.querySelector('script[nonce]'),
            n = e && (e['nonce'] || e.getAttribute('nonce'));
        if (n) {
            po.setAttribute('nonce', n);
        }
        var s = d.getElementsByTagName('script')[0];
        s.parentNode.insertBefore(po, s);
    })();
    """)

    time.sleep(3)

    try:
        browser.execute_async_script("""
        window.tokenVal = "";
        var done = arguments[0];
        grecaptcha.execute('""" + client_key + """', { action: '""" + action + """' }).then((token) => {window.tokenVal = token; done(token)});
        """)
    except Exception:
        return "page did not render recaptcha properly"

    return browser.execute_script("return window.tokenVal")


if __name__ == "__main__":
    print(main(sys.argv[1], sys.argv[2], sys.argv[3]))
