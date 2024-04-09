# ruuvimetrics

This program maintains Ruuvi sensor values and exposes them via HTTP for a
Prometheus-like consumer who accepts [the text-based exposition
format](https://github.com/prometheus/docs/blob/main/content/docs/instrumenting/exposition_formats.md#text-based-format).
The program composes nicely with [ruuviscan](https://github.com/susji/ruuviscan)
and [ruuviparse](https://github.com/susji/ruuviparse). You can feed input data
to ruuvimetrics directly from Bluetooth Low Energy (BLE) reception (ruuviscan)
or with raw BLE Advertisement packets (ruuviparse). The idea is that those
programs generate sensor values as JSON to standard output and this program
reads sensor values from standard input.

# Usage

`ruuvimetrics` has some command-line arguments. To see them, invoke

    $ ruuvimetrics -h

For local BLE scans, invoke

    $ ruuviscan | ruuvimetrics

and then wait for values to be aggregated and query them with some HTTP client:

    $ curl http://localhost:9900/metrics

For parsing and exposing Ruuvi sensor values via MQTT, you probably have to
figure out how your transmitters or BLE proxies send the values, but it might
look something like
[this](https://github.com/susji/ruuviparse#integrating-with-mqtt-messaging):

    $ mosquitto_sub \
        -h mqtt.example.com \
        -t ruuvi/001 \
        --cert ruuvi-listener.cert.pem \
        --key ruuvi-listener.private.key \
        --cafile ca.crt \
        -i ruuvi-listener \
        | jq --raw-output --unbuffered \
          '.ads[].ad | ascii_downcase | select(.[8:14] == "ff9904") | .[14:]' \
        | ruuviparse \
        | ruuvimetrics
