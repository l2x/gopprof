
### cpu

go tool pprof -text  example cpu_1464145985.pprof > /vagrant/test.text

go tool pprof -list=.* example cpu_1464145985.pprof > /vagrant/test.text

go tool pprof -pdf example cpu_1464145985.pprof > /vagrant/test.pdf

### goroutine

goroutine file

go tool pprof -list=.* 

go tool pprof -pdf example goroutine_1464145985.pprof > /vagrant/test.pdf

go tool pprof -top example goroutine_1464145985.pprof > /vagrant/test.text


### heap

heap file

go tool pprof -pdf 

go tool pprof -top

go tool pprof -list=.*



### block

block file


### trace
