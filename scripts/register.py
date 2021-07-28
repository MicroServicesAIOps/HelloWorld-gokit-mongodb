import requests
import base64


def register(host, user_list):
    for user in user_list:
        register_params = {
            "username": user,
            "password": "password",
            "email": "email"
        }
        resp = requests.post(host + "/register", json=register_params)
        print(resp.json())
        user_id = resp.json()['id']

        print(user_id)


if __name__ == '__main__':
    host = "http://yourID:8080"
    user_id_start = 1006
    user_id_end = 1020
    user_list = []
    for i in range(user_id_start, user_id_end):
        user_list.append("user" + str(i))
    register(host, user_list)