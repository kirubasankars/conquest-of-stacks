user1_token=$(curl -s -X POST 'http://localhost:8080/login?username=admin&password=admin' | jq '.token' -r)
user2_token=$(curl -s -X POST 'http://localhost:8080/login?username=test&password=test' | jq '.token' -r)

curl -s -X POST 'http://localhost:8080/auth/join_game/15s' -H "Authorization:Bearer $user1_token"
curl -s -X POST 'http://localhost:8080/auth/join_game/15s' -H "Authorization:Bearer $user2_token"

