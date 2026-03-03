from locust import HttpUser, task, between


class URLShortenerUser(HttpUser):
    wait_time = between(0.1, 0.5)

    short_url: str = ""

    def on_start(self):
        response = self.client.post(
            "/",
            json={"url": "https://example.com"},
        )
        response.raise_for_status()
        self.short_url = response.json().get("shortUrl", "")

    @task(3)
    def shorten(self):
        response = self.client.post(
            "/",
            json={"url": "https://example.com"},
        )
        if response.status_code != 201:
            response.failure(f"expected 201, got {response.status_code}")

    @task(7)
    def resolve(self):
        if not self.short_url:
            return
        response = self.client.get(f"/{self.short_url}", allow_redirects=False)
        if response.status_code != 200:
            response.failure(f"expected 200, got {response.status_code}")
