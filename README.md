# GSStat - Game Server Stat

* https://github.com/battlesrv/gsstat

## Steam

* https://developer.valvesoftware.com/wiki/Server_queries

## Minecraft

* https://wiki.vg/Query

## Use

```
$ make
$ gsstat m -addr 46.174.49.30:25817
{"hostname":"State27","game_type":"SMP","game_id":"MINECRAFT","version":"1.13.2","plugins":"CraftBukkit on Bukkit 1.13.2-R0.1-SNAPSHOT","map":"World","numplayers":9,"maxplayers":21,"hostport":25817,"hostip":"46.174.49.30","players":["Ebatelhuya","FearPavor","killzzz","Banchivilly","Vityalox","kopatych","Xasle","Shaleniy","Gogetamon"]}
```

## Simple test

```
echo -en '\xFE\x01' | nc 1.2.3.4 12345
```

## License

* [GPL-3.0](./LICENSE)

## Author

* Konstantin Kruglov
* kruglovk@gmail.com
