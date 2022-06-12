# HTTP router

It is an adoption of [julienschmidt/httprouter](https://github.com/julienschmidt/httprouter) router.

It is optimized for serving thousands of routes, including:
 * Move handlers out of the tree to reduce memory consumption. 
 * Define one global handler for all valid routes.
 * Reduce pointers usage in the tree.
 * Remove some original [features](https://github.com/julienschmidt/httprouter#features) because I did not need them.

The package provides a router for [valyala/fasthttp](https://github.com/valyala/fasthttp) and Go's standard [net/http](https://pkg.go.dev/net/http).

