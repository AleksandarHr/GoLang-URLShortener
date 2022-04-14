# The only existing record before running this script is {originalurl:"https://whatever.this.is", shorturl:"LofFb"}
curl http://localhost:8080/ -H 'Content-Type: text/plain;charset=UTF-8' --data-raw 'https://whatever.this.is'
# Shortened URL: LofFb --> Returns the same shortened URL on second try
curl http://localhost:8080/ -H 'Content-Type: text/plain;charset=UTF-8' --data-raw 'https://whatever.this.is'
# Shortened URL: LofFb --> Returns the same shortened URL on 2+ try

curl http://localhost:8080/LofFb
# <a href="https://whatever.this.is">See Other</a>.  --> Redirects to original long URL

curl http://localhost:8080/ -H 'Content-Type: text/plain;charset=UTF-8' --data-raw ''
# parse "": empty url 
curl http://localhost:8080/asd
# Supports only short URLs of length 5. --> See notes in ./algorithms/algorithms.go

curl http://localhost:8080/ -H 'Content-Type: text/plain;charset=UTF-8' --data-raw 'httos://someverylongdomainnamehere.com/some/very/very/long/path/here?foo=bar'
# Shortened URL: 8TAAM  --> NOTE: Not guaranteed to be the same, due to the uniqueness of the UNIX timestamp
