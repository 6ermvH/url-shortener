import random

import requests
from locust import HttpUser, between, events, task

SEED_COUNT = 100

short_urls: list[str] = []


@events.test_start.add_listener
def seed(environment, **kwargs):
    for i in range(SEED_COUNT):
        url = f"https://example.com/page/{i}"
        resp = requests.post(
            f"{environment.host}/",
            json={"url": url},
            timeout=5,
        )
        if resp.status_code == 201:
            short = resp.json().get("shortUrl", "")
            if short:
                short_urls.append(short)


class URLShortenerUser(HttpUser):
    wait_time = between(0.1, 0.5)

    @task(3)
    def shorten(self):
        url = f"https://example.com/page/{random.randint(0, 999)}"
        response = self.client.post("/", json={"url": url})
        if response.status_code == 201:
            short = response.json().get("shortUrl", "")
            if short and short not in short_urls:
                short_urls.append(short)
        else:
            response.failure(f"expected 201, got {response.status_code}")

    @task(7)
    def resolve(self):
        if not short_urls:
            return
        short = random.choice(short_urls)
        response = self.client.get(f"/{short}", allow_redirects=False)
        if response.status_code != 200:
            response.failure(f"expected 200, got {response.status_code}")
