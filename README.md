This is a very concise example of a
[Tigertonic](https://github.com/rcrowley/go-tigertonic) service. The `basic`
version is just the bare bones of a service; the `instrumented` version is the
exact same thing, but with some metrics added.

The purpose of this code is to illustrate how easy it is to add various
instrumentation to a Tigertonic service while maintaining very good separation
of concerns.

## Sample Output

```bash
# insert a book
curl -XPUT -H "Content-Type: application/json" -d '{"author":"Pushkin", "title":"Eugene Onegin"}'  localhost:34334/books/1234
# retrieve teh same book
curl localhost:34334/books/1234
```
```json
{"Author":"Pushkin","Title":"Eugene Onegin"}
```

## Motivation

I wrote this code for [this blog
entry](http://blog.mihasya.com/2014/02/07/tt-metrics.html). It's also a nice and
simple illustration of some core concepts in Tigertonic.
