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
