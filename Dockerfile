FROM golang:1.15

RUN apt-get update
RUN apt-get install -y build-essential zlib1g-dev libncurses5-dev libgdbm-dev libnss3-dev libssl-dev libreadline-dev libffi-dev wget gconf-service libasound2 libatk1.0-0 libcairo2 libcups2 libfontconfig1 libgdk-pixbuf2.0-0 libgtk-3-0 libnspr4 libpango-1.0-0 libxss1 fonts-liberation libappindicator1 lsb-release xdg-utils unzip

RUN mkdir python3.6-install && cd python3.6-install
RUN wget https://www.python.org/ftp/python/3.6.0/Python-3.6.0.tar.xz
RUN tar xvf Python-3.6.0.tar.xz
RUN cd Python-3.6.0/ && ./configure && make altinstall
RUN cd ../.. && rm -rf python3.6-install/

WORKDIR /app
COPY . .

RUN unzip bin/chromedriver.zip -d bin/
RUN unzip bin/headless-chromium.zip -d bin/
RUN rm bin/chromedriver.zip bin/headless-chromium.zip
RUN chmod 755 bin/chromedriver

RUN python3.6 -m pip install --upgrade pip
RUN python3.6 -m pip install -r reCaptcha/requirements.txt

RUN go build -o main .

ENTRYPOINT ["/app/main"]
