# Instructions
1. Execute `<go run main.go>` in the root of folder
2. Execute `<curl http://localhost:8080/<shortURL>>` where *\<shortURL\>* is any 5-symbol string. This will query the DB for a corresponding long URL and redirect, if it exists.
3. Execute `<curl http://localhost:8080/ -H 'Content-Type: text/plain;charset=UTF-8' --data-raw '<longURL>'>` where *\<longURL\>* is any valid URL string. This will query the DB for a record with the same long URL. If it exists, it will return the corresponding short URL. Otherwise, it will run the shortening algorithm, insert the new record in the db, and return the short URL.

# Notes
1. **MongoDB Connection** - I am using MongoDB Atlas. Everything is setup - required variables for the connection are stored in and read from the .env file and I have granted access to connection requests coming from *any* IP address. There shouldn't be any problems, I hope. 
2. **Base58 Encoding** - We allow only digits and latin letters, with the exception for the characters 'o', 'O', 'I', '1', resulting in base58 encoding. For more information: https://en.bitcoin.it/wiki/Base58Check_encoding
3. **5-char short URLs** - We fix the shortened URL length to exactly 5 characters (as opposed to at most 5). The reason is:
    - (To simplify the implementation)
    - Because I had misread the result of 58^5 as being more than 1 billion (which was the required number of URLs supported). Unfortunately, I found out about it a bit too late to fix it. Sorry about that.
4. **Tests** - I tried, unsuccessfully, to write unit tests for the two handles. My understanding is that you cannot unit-test the mongo-db driver without mocking the db locally in some way. After some time reading up on it, I decided to drop it. There is a very basic test script (*testScript.sh*) instead.