import wikipediaapi
import requests
from datetime import datetime, timezone

base_url = "http://127.0.0.1:8443/api/v1/info_objects"


def check_resp(resp):
    if not resp.ok:
        print("request error: {}\n{}".format(resp.status_code, resp.text))
        raise Exception()


def upload(page):
    categories = list(page.categories.keys())
    tags = list(filter(valid_tag, map(lambda s: s[9:], categories)))

    published = datetime.now(timezone.utc).isoformat()

    json = {"name": page.title, "published": published, "author": "test", "source": "wiki", "tags": tags}
    resp = requests.post(base_url, json)

    check_resp(resp)

    resp_json = resp.json()
    obj_id = resp_json["id"]

    upload_url = "{}/{}/files".format(base_url, obj_id)

    files = {'file': ('content.txt', page.text, "text/plain")}
    resp = requests.post(upload_url, files=files)

    check_resp(resp)

    print("info object {} {} created and uploaded".format(page.title, obj_id))
    print(tags)


def valid_tag(s):
    ls = s.lower()
    return not ("articles" in ls or
                "wikidata" in ls or
                "links" in ls or
                "help desk" in ls or
                "dmy" in ls or
                "mdy" in ls or
                "cs1" in ls)


class Loader:
    def __init__(self):
        self.wiki = wikipediaapi.Wikipedia('en', extract_format=wikipediaapi.ExtractFormat.WIKI)
        self.loaded_count = 0
        self.visited = set()
        self.queue = []

    def start(self, page_name):
        self.queue.append(page_name)
        self.do_load()

    def do_load(self):
        while self.queue:
            page_name = self.queue.pop(0)

            if page_name in self.visited:
                return

            if self.loaded_count > 100:
                return

            page = self.wiki.page(page_name)
            if not page.exists():
                continue

            upload(page)

            self.visited.add(page.title)
            self.loaded_count += 1

            for link in page.links:
                self.queue.append(link)


if __name__ == '__main__':
    loader = Loader()
    loader.start('Python_(programming_language)')


