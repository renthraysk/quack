# QUACK

This package produces and consumes QPACK encoded headers.

It uses canonical (Content-Type) header name style of HTTP/1 & go/http. 
However it does ensure they are lower cased when they need to be present as 
literal names in header frames to comply with HTTP/3 RFC9114. 
Also will canonicalise incoming header names to HTTP/1 camel case style.

## Notes

- QPACK's static table is used to generate pure go code. 
