# -*- coding: utf-8 -*-
from logging import getLogger, FileHandler, Formatter, DEBUG
from os import chdir, walk, mkdir, listdir, system, remove
from os.path import isfile, join, isdir
from subprocess import Popen, PIPE, call
from getpass import getpass
from uuid import uuid4
from threading import Thread
import json
from multiprocessing import Pool
from functools import partial
from datetime import datetime

from hashlib import md5
from pick import pick
from hashids import Hashids
from termcolor import cprint
from tqdm import tqdm
import ruamel.yaml as yaml
from ruamel.yaml.comments import CommentedMap

ALPHABET = "abcdefghijklmnopqrstuvwxyz"

# create logger with 'spam_application'
logger = getLogger('gojot')
logger.setLevel(DEBUG)
# create file handler which logs even debug messages
fh = FileHandler('gojot.log')
fh.setLevel(DEBUG)
# create formatter and add it to the handlers
formatter = Formatter(
    '%(asctime)s - %(funcName)10s - %(levelname)s - %(message)s')
fh.setFormatter(formatter)
# add the handlers to the logger
logger.addHandler(fh)


class MyException(Exception):
    pass


def encode_str(s, salt):
    hashids = Hashids(salt=salt)
    nums = []
    for let in s:
        nums.append(int(ALPHABET.index(let)))
    return hashids.encode(*nums)


def decode_str(s, salt):
    hashids = Hashids(salt=salt)
    new_s = ""
    for num in hashids.decode(s):
        new_s += ALPHABET[num]
    return new_s


def git_log():
    GIT_COMMIT_FIELDS = ['id', 'author_name',
                         'author_email', 'date', 'message']
    GIT_LOG_FORMAT = ['%H', '%an', '%ae', '%ad', '%s']
    GIT_LOG_FORMAT = '%x1f'.join(GIT_LOG_FORMAT) + '%x1e'

    p = Popen('git log --format="%s"' % GIT_LOG_FORMAT,
              shell=True, stdout=PIPE, stderr=PIPE)
    (log, _) = p.communicate()
    log = log.strip(b'\n\x1e').split(b"\x1e")
    log = [row.strip().split(b"\x1f") for row in log]
    log = [dict(zip(GIT_COMMIT_FIELDS, row)) for row in log]
    return log


def git_clone():
    p = Popen('git clone git@github.com:schollz/test5.git',
              shell=True, stdout=PIPE, stderr=PIPE)
    (log, logerr) = p.communicate()
    logger.debug(log)
    logger.debug(logerr)


def git_pull():
    p = Popen('git pull --rebase origin master',
              shell=True, stdout=PIPE, stderr=PIPE)
    (log, logerr) = p.communicate()
    logger.debug(log)
    logger.debug(logerr)


def decrypt(fname, passphrase):
    p = Popen('gpg --yes --passphrase "{passphrase}" --decrypt {fname}'.format(
        passphrase=passphrase, fname=fname), shell=True, stdout=PIPE, stderr=PIPE)
    (log, logerr) = p.communicate()
    logger.debug(log)
    logger.debug(logerr)
    if b"bad passphrase" in logerr:
        raise MyException("Bad passphrase")
    return log


def add_file(fname, contents, recipient):
    with open(fname, "w") as f:
        f.write(contents)
    p = Popen('gpg --yes --armor --recipient "%s" --trust-model always --encrypt  %s' %
              (recipient, fname), shell=True, stdout=PIPE, stderr=PIPE)
    (log, logerr) = p.communicate()
    logger.debug(log)
    logger.debug(logerr)
    remove(fname)
    p = Popen('git add %s.asc' % fname, shell=True, stdout=PIPE, stderr=PIPE)
    (log, logerr) = p.communicate()
    logger.debug(log)
    logger.debug(logerr)
    p = Popen('git commit -m "%s.asc"' %
              fname, shell=True, stdout=PIPE, stderr=PIPE)
    (log, logerr) = p.communicate()
    logger.debug(log)
    logger.debug(logerr)


def pick_key():
    p = Popen('gpg --list-keys', shell=True, stdout=PIPE, stderr=PIPE)
    (log, _) = p.communicate()
    keys = []
    usernames = []
    for gpg_key in log.split(b"------\n")[1].split(b"\n\n"):
        try:
            (pub, uid, sub) = gpg_key.split(b"\n")
        except:
            continue
        pub = pub[4:].strip()
        uid = uid[4:].strip()
        keys.append("{uid} {pub}".format(
            pub=pub.decode('utf-8'), uid=uid.decode('utf-8')))
        usernames.append(uid.decode('utf-8').split("<")[0].strip())
    [_, index] = pick(keys, "Pick key: ")
    return usernames[index]


call('clear', shell=True)
cprint("Working on schollz/test5", "green")
chdir("/tmp/")
git_thread = None
if isdir("test5"):
    cprint("Pulling the latest...", "yellow")
    chdir("test5")
    git_thread = Thread(target=git_pull)
    git_thread.start()
else:
    cprint("Cloning the latest...", "yellow")
    git_clone()
    chdir("test5")

# Check if config file exists
config = {}
if not isfile("config.asc"):
    cprint("Generating credentials...", "yellow")
    username = pick_key()
    config = {"user": username, "salt": str(uuid4())}
    add_file("config", json.dumps(config), config['user'])
    cprint("...ok.", "yellow")

# CHECK PASSPHRASE
passphrase = getpass("Passphrase: ")
cprint("Checking credentials...", "yellow")
try:
    content = decrypt("config.asc", passphrase)
    cprint("...ok.", "yellow")
except:
    cprint("...bad credentials.", "red")
    exit(1)
config = json.loads(content.decode('utf-8'))


if git_thread != None:
    git_thread.join()
    cprint("...pulled latest.", "yellow")


subjects = []
for d in [x[0] for x in walk(".")]:
    if ".git" not in d and d != ".":
        subjects.append(decode_str(d[2:], config['salt']))


subject = "New"
if len(subjects) > 0:
    [subject, index] = pick(["New"] + subjects, "Enter subject: ")
if subject == "New":
    subject = input("Document? ")
    mkdir(encode_str(subject, config['salt']))

subject = encode_str(subject, config['salt'])
chdir(subject)

all_files = [f for f in listdir(".") if isfile(join(".", f))]


# files_in_subject = [f for f in listdir(".") if isfile(join(".", f))]

# files = []
# files_names = []
# for fname in all_files:
#     if fname in files_in_subject:
#         files.append(fname)
#         files_names.append(all_files_nicename[all_files.index(fname)])

# [file_to_edit, index] = pick(["New"] + all_files_nicename, "Pick entry: ")
file_to_edit = "New"
if file_to_edit != "New":
    files = [all_files[index - 1]]
else:
    files = all_files

file_contents = []

p = Pool(4)
max_ = len(files)
with tqdm(total=max_) as pbar:
    for i, datum in tqdm(enumerate(p.imap_unordered(partial(decrypt, passphrase=passphrase), files))):
        file_contents.append(
            {"text": datum.decode('utf-8'), "meta": {"file": "asldkfjaskdf"}})
        pbar.update()


with open("/tmp/temp.txt", "wb") as f:
    for file_data in file_contents:
        f.write(yaml.dump(file_data['meta'],
                          Dumper=yaml.RoundTripDumper).encode('utf-8'))
        f.write(b"\n---\n")
        f.write(file_data['text'].encode('utf-8'))
        f.write(b"\n---\n")
    current_entry = CommentedMap()
    current_entry['time'] = str(datetime.now().strftime("%Y-%m-%d %H:%M:%S"))
    current_entry['entry'] = str(uuid4())
    f.write(yaml.round_trip_dump(current_entry).encode('utf-8'))
    f.write(b"\n---\n\n")
system("vim /tmp/temp.txt")

temp_contents = open("/tmp/temp.txt", "r").read()
for entry in temp_contents.split("---"):
    entry = entry.strip()
    m = md5()
    m.update(entry.encode('utf-8'))
    entry_hash = m.hexdigest()
    if not isfile(entry_hash + ".asc"):
        add_file(entry_hash, entry, config['user'])
remove("/tmp/temp.txt")


# Import
# gpg --import ../public_key_2017.gpg
# gpg --allow-secret-key-import --import ../private_key_2017.gpg

# Encrypt (must do one file at a time)
# gpg --yes --armor --recipient "Zackary N. Scholl" --encrypt test2.txt
# os.system('gpg --yes --armor --recipient "Zackary N. Scholl" --encrypt
# test2.txt')

# Decrypt (batch)
# gpg --yes --decrypt-files *.asc
# gpg --yes --passphrase "PASSPHRASE" --decrypt *.asc

# from pgp import *   #pip install py-pgp
# from pgp.keyserver import get_keyserver

# ks = get_keyserver('hkp://pgp.mit.edu/')
# results = ks.search('zack.scholl@gmail.com')
# print(results)
# for result in results:
# 	recipient_key = result.get()
# 	print(recipient_key.user_ids[0],recipient_key.fingerprint, recipient_key.creation_time)
# 	break
