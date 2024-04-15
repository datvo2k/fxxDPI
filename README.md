# FXXDPI - Use to bypass ISP blocking of some domains
This software designed to bypass DPI with techniques like DoH, ECH and split CLIENT-HELLO packet into small chunks.

## REFERENCES
[A walk-through of an SSL handshake](https://www.commandlinefanatic.com/cgi-bin/showarticle.cgi?article=art059)   
[A walkthrough of a TLS 1.3 handshake](https://commandlinefanatic.com/cgi-bin/showarticle.cgi?article=art080)   
[The Transport Layer Security (TLS) Protocol Version 1.3](https://datatracker.ietf.org/doc/html/rfc8446)   
[DNS Queries over HTTPS (DoH)](https://datatracker.ietf.org/doc/html/rfc8484)   
[Green Tunnel](https://github.com/SadeghHayeri/GreenTunnel)   
[A dive into TLS 1.3 Architecture and its advantages over TLS 1.2](https://medium.com/@akarX23/a-dive-into-tls-1-3-architecture-and-its-advantages-over-tls-1-2-2f552de24fa0)


## IDEA:
- Run proxy localhost 
- Use DNS providers like cloudflare, GoogleDNS, OpenDNS

## EXPLAIN:
### HOW TO WORK

## LOGS:
- version 0.1

## HOW TO CHECK
```
http_proxy=http://127.0.0.1:10053 wget -O - https://example.com
```

## VERSIONS:
- 7/4/2023: v0.1