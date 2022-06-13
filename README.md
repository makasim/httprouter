# HTTP router

It is an adoption of [julienschmidt/httprouter](https://github.com/julienschmidt/httprouter) router optimized for serving thousands of routes.

Including:
 * Move handlers out of the tree to reduce memory consumption. 
 * Define one global handler for all valid routes.
 * Reduce pointers usage in the tree.
 * Support static and param routes (aka /foo/{bar}, /foo/bar).

Note: Some original features are removed [features](https://github.com/julienschmidt/httprouter#features) because I did not need them.

The package provides a router for [valyala/fasthttp](https://github.com/valyala/fasthttp) and Go's standard [net/http](https://pkg.go.dev/net/http).

[Benchmarks](https://github.com/makasim/go-http-routing-benchmark/tree/makasim-http-router):
```
$gotest -run=xxx -bench=HttpRouter ./ -v 
#GithubAPI Routes: 203
   HttpRouter: 37088 Bytes
   MakasimHttpRouter: 36808 Bytes
   MakasimGlobalHttpRouter: 31408 Bytes

#GPlusAPI Routes: 13
   HttpRouter: 2760 Bytes
   MakasimHttpRouter: 2936 Bytes
   MakasimGlobalHttpRouter: 2504 Bytes

#Many Routes: 10000
   HttpRouter: 954608 Bytes
   MakasimHttpRouter: 1062504 Bytes
   MakasimGlobalHttpRouter: 742992 Bytes

#ParseAPI Routes: 26
   HttpRouter: 5024 Bytes
   MakasimHttpRouter: 4864 Bytes
   MakasimGlobalHttpRouter: 4144 Bytes

#Static Routes: 157
   HttpRouter: 21680 Bytes
   MakasimHttpRouter: 22808 Bytes
   MakasimGlobalHttpRouter: 17408 Bytes

goos: darwin
goarch: amd64
pkg: github.com/julienschmidt/go-http-routing-benchmark
cpu: Intel(R) Core(TM) i7-4980HQ CPU @ 2.80GHz
BenchmarkHttpRouter_Param
BenchmarkHttpRouter_Param                     	12404929	        90.02 ns/op	      32 B/op	       1 allocs/op
BenchmarkMakasimHttpRouter_Param
BenchmarkMakasimHttpRouter_Param              	 9343818	       125.1 ns/op	      24 B/op	       1 allocs/op
BenchmarkHttpRouter_Param5
BenchmarkHttpRouter_Param5                    	 4924988	       243.5 ns/op	     160 B/op	       1 allocs/op
BenchmarkMakasimHttpRouter_Param5
BenchmarkMakasimHttpRouter_Param5             	 2125035	       565.9 ns/op	     120 B/op	       5 allocs/op
BenchmarkHttpRouter_Param20
BenchmarkHttpRouter_Param20                   	 1687647	       713.0 ns/op	     640 B/op	       1 allocs/op
BenchmarkMakasimHttpRouter_Param20
BenchmarkMakasimHttpRouter_Param20            	  552642	      2029 ns/op	     480 B/op	      20 allocs/op
BenchmarkHttpRouter_ParamWrite
BenchmarkHttpRouter_ParamWrite                	10587574	       111.0 ns/op	      32 B/op	       1 allocs/op
BenchmarkMakasimHttpRouter_ParamWrite
BenchmarkMakasimHttpRouter_ParamWrite         	 8094574	       148.9 ns/op	      24 B/op	       1 allocs/op
BenchmarkHttpRouter_GithubStatic
BenchmarkHttpRouter_GithubStatic              	32061902	        36.32 ns/op	       0 B/op	       0 allocs/op
BenchmarkMakasimHttpRouter_GithubStatic
BenchmarkMakasimHttpRouter_GithubStatic       	17407140	        63.08 ns/op	       0 B/op	       0 allocs/op
BenchmarkMakasimGlobalHttpRouter_GithubStatic
BenchmarkMakasimGlobalHttpRouter_GithubStatic 	22013314	        53.32 ns/op	       0 B/op	       0 allocs/op
BenchmarkHttpRouter_GithubParam
BenchmarkHttpRouter_GithubParam               	 5184572	       221.3 ns/op	      96 B/op	       1 allocs/op
BenchmarkHttpRouter_GithubAll
BenchmarkHttpRouter_GithubAll                 	   28773	     41894 ns/op	   13792 B/op	     167 allocs/op
BenchmarkMakasimHttpRouter_GithubAll
BenchmarkMakasimHttpRouter_GithubAll          	   16942	     77183 ns/op	    8136 B/op	     339 allocs/op
BenchmarkHttpRouter_GPlusStatic
BenchmarkHttpRouter_GPlusStatic               	60446588	        21.03 ns/op	       0 B/op	       0 allocs/op
BenchmarkMakasimHttpRouter_GPlusStatic
BenchmarkMakasimHttpRouter_GPlusStatic        	26642642	        38.89 ns/op	       0 B/op	       0 allocs/op
BenchmarkMakasimGlobalHttpRouter_GPlusStatic
BenchmarkMakasimGlobalHttpRouter_GPlusStatic  	36243447	        32.98 ns/op	       0 B/op	       0 allocs/op
BenchmarkHttpRouter_GPlusParam
BenchmarkHttpRouter_GPlusParam                	 7653328	       158.7 ns/op	      64 B/op	       1 allocs/op
BenchmarkHttpRouter_GPlus2Params
BenchmarkHttpRouter_GPlus2Params              	 6572565	       179.5 ns/op	      64 B/op	       1 allocs/op
BenchmarkHttpRouter_GPlusAll
BenchmarkHttpRouter_GPlusAll                  	  583130	      2034 ns/op	     640 B/op	      11 allocs/op
BenchmarkHttpRouter_ManyRoutes
BenchmarkHttpRouter_ManyRoutes                	26469423	        46.87 ns/op	       0 B/op	       0 allocs/op
BenchmarkMakasimHttpRouter_ManyRoutes
BenchmarkMakasimHttpRouter_ManyRoutes         	12506808	        87.37 ns/op	       0 B/op	       0 allocs/op
BenchmarkMakasimGlobalHttpRouter_ManyRoutes
BenchmarkMakasimGlobalHttpRouter_ManyRoutes   	15010459	        79.22 ns/op	       0 B/op	       0 allocs/op
BenchmarkHttpRouter_ParseStatic
BenchmarkHttpRouter_ParseStatic               	51509540	        22.12 ns/op	       0 B/op	       0 allocs/op
BenchmarkMakasimHttpRouter_ParseStatic
BenchmarkMakasimHttpRouter_ParseStatic        	28528497	        39.12 ns/op	       0 B/op	       0 allocs/op
BenchmarkMakasimGlobalHttpRouter_ParseStatic
BenchmarkMakasimGlobalHttpRouter_ParseStatic  	35581128	        33.30 ns/op	       0 B/op	       0 allocs/op
BenchmarkHttpRouter_ParseParam
BenchmarkHttpRouter_ParseParam                	 8653006	       166.0 ns/op	      64 B/op	       1 allocs/op
BenchmarkHttpRouter_Parse2Params
BenchmarkHttpRouter_Parse2Params              	 6703875	       159.8 ns/op	      64 B/op	       1 allocs/op
BenchmarkHttpRouter_ParseAll
BenchmarkHttpRouter_ParseAll                  	  435255	      2729 ns/op	     640 B/op	      16 allocs/op
BenchmarkHttpRouter_StaticAll
BenchmarkHttpRouter_StaticAll                 	  100560	     11007 ns/op	       0 B/op	       0 allocs/op
PASS
ok  	github.com/julienschmidt/go-http-routing-benchmark	43.369s
```
