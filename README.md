# Satdress

Federated Lightning Address Server

---------------
This is an attempt to make a Cloudron app of Satdress. The below instructions is for Docker. Instructions for Cloudron will apear as soon the app is tested and working.

## How to run

1. Download the binary from the releases page (or compile with `go build` or `go get`)
2. Set the following environment variables somehow (using example values from bitmia.com):

```
PORT=17422
DOMAIN=bitmia.com
SECRET=askdbasjdhvakjvsdjasd
SITE_OWNER_URL=https://t.me/qecez
SITE_OWNER_NAME=@qecez
SITE_NAME=Bitmia
```

3. Start the app with `./satdress`
4. Serve the app to the world on your domain using whatever technique you're used to

## Get help

Maybe ask for help on https://t.me/lnurl if you're in trouble.
