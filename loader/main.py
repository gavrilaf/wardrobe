import wikipediaapi
import requests

# Press ⌃R to execute it or replace it with your code.
# Press Double ⇧ to search everywhere for classes, files, tool windows, actions, and settings.


base_url = "http://127.0.0.1:8443/api/v1/fo"


def check_resp(resp):
    if not resp.ok:
        print("request error: {}\n{}".format(resp.status_code, resp.text))
        raise Exception()


def upload(title, text, tags):
    json = {"name": title, "content_type": "text/plain", "author": "test", "source": "wiki", "tags": tags}
    resp = requests.post(base_url, json)

    check_resp(resp)

    resp_json = resp.json()
    fo_id = resp_json["id"]

    upload_url = "{}/{}".format(base_url, fo_id)
    files = {'file': ('file', text, "text/plain")}

    resp = requests.put(upload_url, files=files)

    check_resp(resp)

    print("file object {} {} created and uploaded".format(title, fo_id))
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

            categories = list(page.categories.keys())
            tags = list(filter(valid_tag, map(lambda s: s[9:], categories)))

            upload(page.title, page.text, tags)

            self.visited.add(page.title)
            self.loaded_count += 1

            for link in page.links:
                self.queue.append(link)


if __name__ == '__main__':
    loader = Loader()
    loader.start('Python_(programming_language)')


