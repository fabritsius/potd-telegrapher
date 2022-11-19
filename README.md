![POTD Telegrapher](https://fabritsius.github.io/potd-telegrapher/img/potd-banner.jpg)

This project gets [Wikipedia's Picture of the Day](https://en.wikipedia.org/wiki/Wikipedia:Picture_of_the_day) and publishes it on the [Telegram Channel](https://t.me/w_potd).

### How it works?

- Get the content from [Wikipedia](https://en.wikipedia.org/wiki/Wikipedia:Picture_of_the_day)
- Create a [telegra.ph](https://telegra.ph/) article
- Post an article link to the [Telegram Channel](https://t.me/w_potd)
- A **cron** job is launched every day which does all of the above

### Roadmap

- [x] Release 1.0.0
- [ ] Add an action that only creates an article
- [ ] Use the Wiki API instead of web scraping (potentially)
- [ ] Preserve all of Wikipedia's links (potentially)
- [ ] Additionally send an article to Instagram
