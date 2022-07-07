# Keylogger

Find your own event handler from `/dev/input/*`. You can find your own event file via command: 
```terminal
cat /proc/bus/input/devices
```

Run:
```terminal
go build . && sudo ./koki
```
