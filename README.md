MCApi
=====

An API for fetching the status of and querying Minecraft servers.

It is running at [mcapi.us](https://mcapi.us).

## Configuration

```bash
go build
./mcapi -gencfg
```

Options in config.json:
* HTTPAppHost &mdash; host and port to listen on
* RedisHost &mdash; host of redis server
* StaticFiles &mdash; path to static files
* TemplateFile &mdash; path to index file
* SentryDSN &mdash; optional sentry dsn to report errors to
* AdminKey &mdash; secret token used to get list of servers or clear the list

Rate limiting with Cloudflare requires setting the following environment variables:
* CLOUDFLARE_EMAIL &mdash; your Cloudflare account email address
* CLOUDFLARE_AUTH &mdash; your Cloudflare authentication token

Disable ratelimiting generally or with Cloudflare by using built-time variables
to update `rateLimitEnabled` and `cloudflareEnabled` or modifying the code
in `ratelimit.go#16-17`.

Setting `APPROVED_IPS` to a comma separated list of IP addresses will prevent
the rate limits from applying to those addresses.

----------

本fork在部分文件上对于中文进行了十分友好的修复。  
详见文件主体内容。  
字体文件因潜在的版权问题所以不放在仓库内，请自行下载/复制。  
index.html未进行翻译。都找到这里来了还看不懂的话个人建议使用别人已搭建好的服务器。  
比如说 mcapi.amazefcc233.com（个人服务器，不保证稳定性  

自己加的部分代码十分垃圾，都是从网上找来的。  
若您有go开发水平，建议fork源仓库然后自行修改代码=-=  
