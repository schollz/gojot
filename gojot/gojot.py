# -*- coding: utf-8 -*-
from logging import getLogger, FileHandler, Formatter, DEBUG
from os import chdir, walk, mkdir, listdir, system, remove
from os.path import isfile, join, isdir, expanduser
from subprocess import Popen, PIPE, call
from getpass import getpass
from uuid import uuid4
from threading import Thread
import json
from multiprocessing import Pool
from functools import partial
from datetime import datetime
from copy import deepcopy
import atexit

from hashlib import md5
from pick import pick
from hashids import Hashids
from termcolor import cprint
from tqdm import tqdm
import ruamel.yaml as yaml
from ruamel.yaml.comments import CommentedMap

from names import *

ALPHABET = "abcdefghijklmnopqrstuvwxyz0123456789 !@#$%^&*()-=_+"

VIMRC = """set nocompatible
set backspace=2
func! WordProcessorModeCLI()
    setlocal formatoptions=t1
    setlocal textwidth=80
    map j gj
    map k gk
    set formatprg=par
    setlocal wrap
    setlocal linebreak
    setlocal noexpandtab
    normal G$
    normal zt
    set foldcolumn=7
    highlight Normal ctermfg=black ctermbg=grey
    hi NonText ctermfg=grey guifg=grey
endfu
com! WPCLI call WordProcessorModeCLI()
"""

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


@atexit.register
def clean_files():
    try:
        remove("/tmp/temp.txt")
        cprint("Removed temp files.", "yellow")
    except:
        cprint("Exited.", "yellow")
    try:
        system(
            'mv {0}/.gnupg/gpg.conf.backup {0}/.gnupg/gpg.conf'.format(expanduser('~')))
    except:
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


def git_clone(repo):
    p = Popen('git clone ' + repo,
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

def git_push():
    p = Popen('git push -u origin master',
              shell=True, stdout=PIPE, stderr=PIPE)
    (log, logerr) = p.communicate()
    logger.debug(log)
    logger.debug(logerr)



def decrypt(fname, passphrase):
    p = Popen('gpg --yes --passphrase "{passphrase}" --decrypt {fname}'.format(
        passphrase=passphrase, fname=fname), shell=True, stdout=PIPE, stderr=PIPE)
    (log, logerr) = p.communicate()
    # logger.debug(log)
    logger.debug(logerr)
    if b"bad passphrase" in logerr:
        raise MyException("Bad passphrase")
    return log


def add_file(fname, contents, recipient, add_to_git=True):
    with open(fname, "w") as f:
        f.write(contents)
    p = Popen('gpg --yes --armor --recipient "%s" --trust-model always --encrypt  %s' %
              (recipient, fname), shell=True, stdout=PIPE, stderr=PIPE)
    (log, logerr) = p.communicate()
    logger.debug(log)
    logger.debug(logerr)
    remove(fname)
    if not add_to_git:
        return
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


def parse_entries(entry_data):
    """
    Data should be in the format:

    ---

    yaml

    ---

    text


    """
    datas = []
    data = {'meta': {}}
    for i, entry in enumerate(entry_data.split("---")):
        entry = entry.strip()
        if i % 2 == 1:
            data['meta'] = yaml.load(entry, Loader=yaml.Loader)
        elif data['meta'] != None and len(data['meta']) != 0:
            m = md5()
            m.update(entry.encode('utf-8').strip())
            entry_hash = m.hexdigest()
            data['hash'] = entry_hash
            data['text'] = entry.strip()
            datas.append(deepcopy(data))
    return datas


def fix_gpg_conf():
    gpg_conf = open('{0}/.gnupg/gpg.conf'.format(expanduser('~')), 'r').read()
    if 'no-tty' not in gpg_conf:
        system(
            'cp {0}/.gnupg/gpg.conf {0}/.gnupg/gpg.conf.backup'.format(expanduser('~')))
    with open('{0}/.gnupg/gpg.conf'.format(expanduser('~')), 'a') as f:
        f.write("\nno-tty")


def init(repo):
    repo_dir = repo.split("/")[-1].split(".git")[0].strip()
    fix_gpg_conf()
    call('clear', shell=True)

    cprint("Working on %s" % repo, "green")
    chdir("/tmp/")
    git_thread = None
    if isdir(repo_dir):
        cprint("Pulling the latest...", "yellow")
        chdir(repo_dir)
        git_thread = Thread(target=git_pull)
        git_thread.start()
    else:
        cprint("Cloning the latest...", "yellow")
        git_clone(repo)
        chdir(repo_dir)

    # Check if config file exists
    config = {}
    if not isfile("config.asc"):
        cprint("Generating credentials...", "yellow", end='', flush=True)
        username = pick_key()
        config = {"user": username, "salt": str(uuid4())}
        add_file("config", json.dumps(config), config['user'])
        cprint("...ok.", "yellow")

    # Check passphrase
    passphrase = getpass("\nPassphrase? ")
    cprint("\nChecking credentials...", "yellow", end='', flush=True)
    try:
        content = decrypt("config.asc", passphrase)
        cprint("...ok.", "yellow")
    except:
        cprint("...bad credentials.", "red")
        exit(1)
    config = json.loads(content.decode('utf-8'))

    # Wait for pulling to finish before continuing
    if git_thread != None:
        git_thread.join()
        cprint("...pulled latest.", "yellow")

    config['passphrase'] = passphrase
    return config


def import_file(config, temp_contents):
    for entry in parse_entries(temp_contents):
        if len(entry['text'].strip()) < 2:
            continue
        if "document" not in entry['meta']:
            entry['meta']['document'] = "imported"
        if not isfile(join(encode_str(entry['meta']['document'],config['salt']),entry['hash'] + '.asc')):
            entry['meta']['last_modified'] = str(
                datetime.now().strftime("%Y-%m-%d %H:%M:%S"))
            if "time" not in entry['meta']:
                entry['meta']['time'] = str(
                    datetime.now().strftime("%Y-%m-%d %H:%M:%S"))
            if "entry" not in entry['meta']:
                entry['meta']['entry'] = random_name()
            encoded_subject = encode_str(
                entry['meta']['document'], config['salt'])
            if not isdir(encoded_subject):
                mkdir(encoded_subject)
            entry_text = "---\n\n" + \
                yaml.dump(entry['meta'],  Dumper=yaml.RoundTripDumper) + \
                "\n---\n" + entry['text'].strip()
            cprint("Adding {}".format(entry['hash']),'green')
            add_file(join(encoded_subject, entry[
                     'hash']), entry_text.strip(), config['user'])


def get_file_contents(config, encoded_subject):
    if not isdir(encoded_subject):
        return {}
    all_files = []
    for filename in [f for f in listdir(
            encoded_subject) if isfile(join(encoded_subject, f))]:
        if "file_contents.json.asc" not in filename:
            all_files.append(join(encoded_subject, filename))

    file_contents = {}
    if isfile(join(encoded_subject, 'file_contents.json.asc')):
        logger.debug("Using cache")
        file_contents_string = decrypt(
            join(encoded_subject, 'file_contents.json.asc'), config['passphrase'])
        file_contents = json.loads(file_contents_string.decode('utf-8'))
        known_files = []
        for f in file_contents:
            m = md5()
            m.update(file_contents[f]['text'].strip().encode('utf-8'))
            fname = m.hexdigest() + ".asc"
            known_files.append(join(encoded_subject, fname))
        # Update files to only get ones that aren't accounted for
        logger.debug(known_files)
        logger.debug(all_files)
        logger.debug(len(known_files))
        all_files = list(set(all_files) - set(known_files))
        logger.debug(len(all_files))


    cprint("Getting latest entries...","yellow")
    p = Pool(4)
    max_ = len(all_files)
    with tqdm(total=max_) as pbar:
        for i, datum in tqdm(enumerate(p.imap_unordered(partial(decrypt, passphrase=config['passphrase']), all_files))):
            pbar.update()
            if len(datum) == 0:
                continue
            pieces = datum.decode('utf-8').split('---')
            data = {}
            data['meta'] = yaml.load(pieces[1], Loader=yaml.Loader)
            data['text'] = pieces[2]
            key = data['meta']['time'] + data['meta']['entry']
            if data['meta']['time'] in file_contents:
                if data['meta']['last_modified'] < file_contents[key]['meta']['last_modified']:
                    continue
            file_contents[key] = data
    add_file(join(encoded_subject, "file_contents.json"), json.dumps(
        file_contents), config['user'], add_to_git=False)
    cprint("\n...ok.","yellow")
    return file_contents


def run_import(repo, fname):
    contents = open(fname, 'r').read()
    config = init(repo)
    import_file(config, contents)


def run(repo, subject):
    config = init(repo)

    # Decode subjects
    subjects = []
    for d in [x[0] for x in walk(".")]:
        if ".git" not in d and d != ".":
            subjects.append(decode_str(d[2:], config['salt']))

    if subject == None:
        if len(subjects) > 0:
            [subject, index] = pick(["New"] + subjects, "Enter subject: ")
        
        if len(subjects) == 0 or subject == "New":
            subject = input("\nDocument? ")

    encoded_subject = encode_str(subject, config['salt'])

    file_contents = get_file_contents(config, encoded_subject)
    date_strings = sorted(file_contents.keys())
    with open("/tmp/temp.txt", "wb") as f:
        for date_str in date_strings:
            file_data = file_contents[date_str]
            f.write(b"\n---\n")
            f.write(yaml.dump(file_data['meta'],
                              Dumper=yaml.RoundTripDumper).encode('utf-8'))
            f.write(b"---\n")
            f.write(file_data['text'].encode('utf-8'))
            f.write(b"\n")
        current_entry = CommentedMap()
        current_entry['time'] = str(
            datetime.now().strftime("%Y-%m-%d %H:%M:%S"))
        current_entry['entry'] = str(random_name())
        current_entry['document'] = subject
        f.write(b"\n---\n")
        f.write(yaml.round_trip_dump(current_entry).encode('utf-8'))
        f.write(b"---\n\n")

    with open("/tmp/vimrc.config", "w") as f:
        f.write(VIMRC)

    system("vim -u /tmp/vimrc.config -c WPCLI +startinsert /tmp/temp.txt")

    import_file(config, open("/tmp/temp.txt", 'r').read())

    cprint("Pushing...","yellow", end='', flush=True)
    try:
        git_push()
        cprint("...ok.","yellow")
    except:
        cprint("...oh well.","red")


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
