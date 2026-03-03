from locust import HttpUser, task, between


class URLShortenerUser(HttpUser):
    wait_time = between(0.1, 0.5)

    short_url: str = ""

    def on_start(self):
        response = self.client.post(
            "/shorten",
            json={"url": "https://example.com"},
        )
        if response.ok:
            self.short_url = response.json().get("shortUrl", "")

    @task(3)
    def shorten(self):
        self.client.post(
            "/shorten",
            json={"url": "https://example.com"},
        )

    @task(7)
    def resolve(self):
        if self.short_url:
            self.client.get(f"/{self.short_url}", allow_redirects=False)
