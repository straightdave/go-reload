# go-reload
experiment for go code hot reloading

>"reloading" might not be accurate. It's not about how a server reloads itself but behaves different due to different plugins.

## Conclusion 1

We can dynamically change server's behaviors without recompiling or restarting.

```bash
# there are two so files 'hello1.so' and 'hello2.so',
# and currently server is using hello1.

$ curl http://localhost:8080/nominator
hello1

$ curl http://localhost:8080/hello
Hello, "Ronaldo"

# then tell server to use hello2

$ curl -XPUT http://localhost:8080/nominator -d 'hello2'
set to "hello2"

$ curl http://localhost:8080/hello
Hello, "Messi"

# then, build a new plugin, put hello3.so in the dir,
# and tell server to use hello3:

$ curl -XPUT http://localhost:8080/nominator -d 'hello3'
set to "hello3"

$ curl http://localhost:8080/hello
Hello, "Neymar"
```


