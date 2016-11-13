#!/usr/bin/python3
import os
import collections
import json
from hashlib import sha256

import requests


header = """
The latest version is %(version)s, released %(date)s.

"""

template = """.. raw:: html

	<a class="download downloadBox" href="%(url)s">
	<div class="platform">%(platform)s</div>
	<div class="reqs">%(reqs)s</div>
	<div>
	<span class="filename">%(filename)s</span>
	<span class="size">(%(size)s)</span>
	</div>
	<div class="checksum">SHA256: %(sha256sum)s</div>
	</a>


"""

r = requests.get("https://api.github.com/repos/schollz/sdees/releases/latest")
data = json.loads(r.text)
version = data['name']
releaseDate = data['published_at'].split("T")[0]

with open("source/downloads.rst", "w") as f:
    f.write(header % {"version": version, 'date': releaseDate})


featured = collections.OrderedDict()
featured["win64"] = {
    "platform": "Microsoft Windows",
    "reqs": "Windows XP or later, Intel 64-bit processor",
    "size": 0,
    "sha256sum": 0,
    "url": ""}
featured["osx"] = {
    "platform": "Apple OS X",
    "reqs": "OS X 10.8 or later, Intel 64-bit processor",
    "size": 0,
    "sha256sum": 0,
    "url": ""}
featured["linux64"] = {
    "platform": "Linux x64",
    "reqs": "Linux 2.6.23 or later, Intel 64-bit processor",
    "size": 0,
    "sha256sum": 0,
    "url": ""}
featured["linux-arm"] = {"platform": "Linux ARM / Raspberry Pi",
                         "reqs": "Linux 2.6.23 or later, ARM processor",
                         "size": 0,
                         "sha256sum": 0,
                         "url": ""}
for f in featured:
    extension = ".tar.gz"
    if f == "win64":
        extension = ".zip"
    featured[f]["url"] = "https://github.com/schollz/sdees/releases/download/" + \
        data["name"] + "/sdees-" + version + "-" + f + extension
    featured[f]["filename"] = "sdees-%s-%s%s" % (version, f, extension)
    os.system("wget " + featured[f]["url"])
    featured[f]["size"] = str(
        int(os.path.getsize(featured[f]["filename"]) / 1000000)) + " MB"
    fi = open(featured[f]["filename"], 'rb')
    featured[f]["sha256sum"] = sha256(fi.read()).hexdigest()
    fi.close()
    try:
        os.remove(featured[f]["filename"])
    except:
        pass
    with open("source/downloads.rst", "a") as fi:
        fi.write(template % featured[f])
