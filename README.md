Yabumi
======

Command line tool for post the message to slack.

Usage
-----

```
Usage:
  yabumi [OPTIONS] [Url]

Application Options:
  -C, --channel=     slack channel to post
  -a, --attachment   use attachment
  -t, --title=       Title text (attachment)
      --title-link=  Title link url (attachment)
  -c, --color=       Color code or 'good', 'warning', 'danger' (attachment)
  -p, --pretext=     optional text that appears above the message attachment block (attachment)
      --author-name= author_name (attachment)
      --author-link= author_link (attachment)
      --image-url=   image url (attachment)
      --thumb-url=   thumbnail image url (attachment)
      --footer=      footer text (attachment)
      --footer-icon= footer icon url (attachment)
  -f, --field=       title,value,short (attachment)
  -D, --debug        enable debug mode. do not send request, show json only
  -v, --version      show version

Help Options:
  -h, --help         Show this help message

Arguments:
  Url:               slack webhook endpoint url
```

Simple message
--------------

```
$ fortune | ./yabumi https://hooks.slack.com/services/xxx/yyy/zzzz
```

![capture](https://github.com/yteraoka/yabumi/raw/master/images/slack-capture2.png)

Attachment
----------

```
$ echo -e "むかしむかし、あるところにおじいさんとおばあさんが..." | ./yabumi \
 -a \
 --color good \
 --pretext "ぷれてきすと" \
 --author-name "昔話ボット" \
 --author-icon "https://3.bp.blogspot.com/-ixMmPNVQHYs/Ws2uyQlDBxI/AAAAAAABLRQ/FtMx5Tz8SGcKtPLC1IDbc4m9EVQ39edmQCLcBGAs/s180-c/fantasy_akuma_purson.png" \
 --title "浦島太郎" \
 --image-url https://2.bp.blogspot.com/-wZyXlcpBXBU/WvQHkT9wodI/AAAAAAABL78/p8xjPHseoM8E-FSdZDATMMSlAiy0EhxIQCLcBGAs/s180-c/mukashibanashi_ryoushi.png \
 --field "ジャンル,昔話,true" \
 --field "タグ,亀、海,true" \
 https://hooks.slack.com/services/xxx/yyy/zzzz
```

![capture](https://github.com/yteraoka/yabumi/raw/master/images/slack-capture1.png)
