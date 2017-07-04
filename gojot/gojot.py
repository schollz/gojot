# -*- coding: utf-8 -*-
from logging import getLogger, FileHandler, Formatter, DEBUG
from os import chdir, walk, mkdir, listdir, system, remove, makedirs, getcwd
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
import sys

from hashlib import md5
from pick import pick
from hashids import Hashids
from termcolor import cprint
from tqdm import tqdm
import ruamel.yaml as yaml
from ruamel.yaml.comments import CommentedMap
from humanize import intcomma

from .random_name import random_name
from .chart import chart

ALPHABET = "abcdefghijklmnopqrstuvwxyz0123456789 !@#$%^&*()-=_+"

tick = '▇'
sm_tick = '|'

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
    highlight Normal ctermfg=grey ctermbg=black
    hi NonText ctermfg=black guifg=black
endfu
com! WPCLI call WordProcessorModeCLI()
"""
VIMRC2 = """set nocompatible
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
    set foldcolumn=7
    normal G$
    highlight Normal ctermfg=grey ctermbg=black
    hi NonText ctermfg=black guifg=black
endfu
com! WPCLI call WordProcessorModeCLI()
"""

if isfile('gojot.log'):
    remove('gojot.log')
    
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

HOMEDIR = join(expanduser("~"),".cache","gojot")


class MyException(Exception):
    pass


@atexit.register
def clean_files():
    try:
        remove("/tmp/temp.txt")
        cprint("Removed temp files.", "yellow")
    except:
        cprint("Exited.", "yellow")
    if isfile('{0}/.gnupg/gpg.conf.backup'.format(expanduser('~'))):
        system('mv {0}/.gnupg/gpg.conf.backup {0}/.gnupg/gpg.conf'.format(expanduser('~')))

def setup_cache():
    if not isdir(HOMEDIR):
        makedirs(HOMEDIR)

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

def git_get_remote_origin_url():
    
    p = Popen('git config --get remote.origin.url',
              shell=True, stdout=PIPE, stderr=PIPE)
    (log, logerr) = p.communicate()
    logger.debug(log)
    logger.debug(logerr)
    if b'does not exist' in logerr:
        raise MyException("repo does not exist")
    return log.decode('utf-8').strip()

def git_clone(repo):
    p = Popen('git clone ' + repo,
              shell=True, stdout=PIPE, stderr=PIPE)
    (log, logerr) = p.communicate()
    logger.debug(log)
    logger.debug(logerr)
    if b'does not exist' in logerr:
        raise MyException("repo does not exist")
    if b'Could not resolve' in logerr:
        raise MyException("unable to connect")

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
    if b'Could not resolve' in logerr:
        raise MyException("unable to connect")

def decrypt(fname, passphrase):
    p = Popen('gpg --yes --passphrase "{passphrase}" --decrypt {fname}'.format(
        passphrase=passphrase, fname=fname), shell=True, stdout=PIPE, stderr=PIPE)
    (log, logerr) = p.communicate()
    # logger.debug(log)
    logger.debug(logerr)
    if b"bad passphrase" in logerr:
        raise MyException("Bad passphrase")
    if b"secret key not" in logerr:
        raise MyException("Secret key not available to decrypt")
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
    (log, logerr) = p.communicate()
    if log == b'':
        raise MyException("need to create gpg key")
    logger.debug(log)
    logger.debug(logerr)
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


def parse_entries(config, entry_data):
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
            m.update(entry.encode('utf-8').strip()+config['salt'].encode('utf-8'))
            entry_hash = m.hexdigest()
            data['hash'] = entry_hash
            data['text'] = entry.strip()
            datas.append(deepcopy(data))
    return datas


def fix_gpg_conf():
    try:
        gpg_conf = open('{0}/.gnupg/gpg.conf'.format(expanduser('~')), 'r').read()
    except:
        return
    if 'no-tty' not in gpg_conf:
        try:
            system(
            'cp {0}/.gnupg/gpg.conf {0}/.gnupg/gpg.conf.backup'.format(expanduser('~')))
        except:
            return
    with open('{0}/.gnupg/gpg.conf'.format(expanduser('~')), 'a') as f:
        f.write("\nno-tty")


def init(repo):
    current_dir = getcwd()
    setup_cache()
    fix_gpg_conf()
    call('clear', shell=True)
    chdir(HOMEDIR)

    # Determine which repo
    if repo != None:
        repo_dir = repo.split("/")[-1].split(".git")[0].strip()        
    else:
        repo_dirs = [d for d in listdir('.') if isdir(join('.', d))]
        repos = []
        for repo_dir in repo_dirs:
            chdir(repo_dir)
            repos.append(git_get_remote_origin_url())
            chdir('..')
        if len(repo_dirs) == 1:
            repo_dir = repo_dirs[0]
            repo = repos[0]
        elif len(repo_dirs) == 0:
            repo = input("Repo? (e.g. https://github.com/USER/REPO.git) ")
            repo_dir = repo.split("/")[-1].split(".git")[0].strip()
        else:
            [_,index] = pick(["New"]+repos,"Which repo? ")
            if index == 0:
                repo = input("Repo? (e.g. https://github.com/USER/REPO.git) ")
                repo_dir = repo.split("/")[-1].split(".git")[0].strip()
            else:
                repo = repos[index-1]
                repo_dir = repo_dirs[index-1]    
    cprint("Working on %s" % repo, "green")


    git_thread = None
    if isdir(repo_dir):
        cprint("Pulling the latest...", "yellow")
        chdir(repo_dir)
        git_thread = Thread(target=git_pull)
        git_thread.start()
    else:
        cprint("Cloning...", "yellow")
        try:
            git_clone(repo)
        except BaseException as e:
            cprint(str(e), "red")
            exit(1)
        chdir(repo_dir)


    # Check if config file exists
    config = {}
    if not isfile("config.asc"):
        cprint("Generating credentials...", "yellow", end='')
        try:
            username = pick_key()
        except:
            print("""
Do you have a GPG key?

You can import your keys using

    gpg --import public.key
    gpg --allow-secret-key-import --import /tmp/private.key

You can make a new GPG key using

    apt-get install rng-tools
    gpg --gen-key

""")
            exit(1)
        config = {"user": username, "salt": str(uuid4())}
        add_file("config", json.dumps(config), config['user'])
        cprint("...ok.", "yellow")

    # Check passphrase
    passphrase = getpass("\nPassphrase? ")
    cprint("\nChecking credentials...", "yellow", end='')
    try:
        content = decrypt("config.asc", passphrase)
        cprint("...ok.", "yellow")
    except BaseException as e:
        cprint(str(e), "red")
        exit(1)
    config = json.loads(content.decode('utf-8'))

    # Wait for pulling to finish before continuing
    if git_thread != None:
        git_thread.join()
        cprint("...pulled latest.", "yellow")

    config['passphrase'] = passphrase
    config['repo'] = repo 
    config['repo_dir'] = repo_dir
    config['current_dir'] = current_dir
    config['output_file'] = '/tmp/temp.txt'
    return config


def import_file(config, temp_contents):
    entry_updates = 0
    for entry in parse_entries(config, temp_contents):
        if len(entry['text'].strip()) < 2:
            continue
        if "document" not in entry['meta']:
            entry['meta']['document'] = "imported"
        if not isfile(join(encode_str(entry['meta']['document'], config['salt']), entry['hash'] + '.asc')):
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
            cprint("Updating {}".format(entry['hash']), 'green')
            add_file(join(encoded_subject, entry[
                     'hash']), entry_text.strip(), config['user'])
            entry_updates += 1
    return entry_updates


def get_file_contents(config, encoded_subject):
    if not isdir(encoded_subject):
        return {}
    all_files = []
    for filename in [f for f in listdir(
            encoded_subject) if isfile(join(encoded_subject, f))]:
        if "file_contents.json.asc" not in filename:
            all_files.append(join(encoded_subject, filename))

    all_file_contents = []
    if isfile(join(encoded_subject, 'file_contents.json.asc')):
        logger.debug("Using cache")
        file_contents_string = decrypt(
            join(encoded_subject, 'file_contents.json.asc'), config['passphrase'])
        all_file_contents = json.loads(file_contents_string.decode('utf-8'))
        known_files = []
        for f in all_file_contents:
            m = md5()
            m.update(f['text'].strip().encode('utf-8')+config['salt'].encode('utf-8'))
            fname = m.hexdigest() + ".asc"
            known_files.append(join(encoded_subject, fname))
        # Update files to only get ones that aren't accounted for
        logger.debug(known_files)
        logger.debug(all_files)
        logger.debug(len(known_files))
        all_files = list(set(all_files) - set(known_files))
        logger.debug(len(all_files))

    cprint("Getting latest entries...", "yellow")
    p = Pool(8)
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
            all_file_contents.append(data)
    add_file(join(encoded_subject, "file_contents.json"), json.dumps(
        all_file_contents), config['user'], add_to_git=False)
    cprint("\n...ok.", "yellow")

    keys_to_ignore = []
    file_contents = {}
    for data in all_file_contents:    
        key = data['meta']['time']
        if key in file_contents:
            if data['meta']['last_modified'] < file_contents[key]['meta']['last_modified']:
                continue
        if data['text'].strip() == 'ignore document' or data['text'].strip() == 'ignore entry':
            keys_to_ignore.append(key)
        file_contents[key] = data
    for key in keys_to_ignore:
        if key in file_contents:
            file_contents.pop(key,None)
    return file_contents


def run_import(repo, fname):
    contents = open(fname, 'r').read()
    config = init(repo)
    import_file(config, contents)


def print_stats(file_contents):
    dates = sorted(file_contents.keys())
    extracted = {}
    for d in dates:
        t = '-'.join(file_contents[d]['meta']['time'].split()[0].split('-')[0:2])
        if t not in extracted:
            extracted[t] = 0
        extracted[t] += len(file_contents[d]['text'].split())
    labels = []
    word_count = []
    total_words = 0
    for d in sorted(extracted.keys()):
        labels.append(d)
        word_count.append(extracted[d])
        total_words += extracted[d]
    chart(labels,word_count)
    print('\n{} words total'.format(intcomma(total_words)))

def run(repo, subject, load_all=False, edit_one=False, export=False, show_stats=False):
    config = init(repo)

    # Decode subjects
    subjects = []
    for d in [x[0] for x in walk(".")]:
        if ".git" not in d and d != ".":
            subjects.append(decode_str(d[2:], config['salt']))

    if subject == None:
        if len(subjects) > 0:
            [subject, index] = pick(["New"] + subjects, "Choose document: ")

        if len(subjects) == 0 or subject == "New":
            subject = input("\nDocument? ")
        if subject == "":
            subject = "notes"
        subject = subject.lower()

    encoded_subject = encode_str(subject, config['salt'])

    file_contents = {}
    if load_all or edit_one or export or show_stats:
        file_contents = get_file_contents(config, encoded_subject)
    date_strings = sorted(file_contents.keys())
    if edit_one:
        if len(file_contents) == 0:
            cprint("There are no entries to edit in {}".format(subject),"red")
            exit(1)
        title_strings = []
        for date_str in date_strings:
            title_strings.append("{} {}".format(date_str.split()[0],file_contents[date_str]['meta']['entry']))
        [_, index] = pick(title_strings, "Pick entry: ")
        date_strings = [date_strings[index]]

    if export:
        config['output_file'] = join(config['current_dir'],subject + ".txt")

    with open(config['output_file'], "wb") as f:
        for date_str in date_strings:
            file_data = file_contents[date_str]
            f.write(b"\n---\n")
            f.write(yaml.dump(file_data['meta'],
                              Dumper=yaml.RoundTripDumper).encode('utf-8'))
            f.write(b"---\n")
            f.write(file_data['text'].encode('utf-8'))
            f.write(b"\n")
        if not edit_one and not export and not show_stats:
            if len(date_strings) > 0:
                f.write(b"\n")
            current_entry = CommentedMap()
            current_entry['time'] = str(
                datetime.now().strftime("%Y-%m-%d %H:%M:%S"))
            entry = input("\nEntry? (enter for default) ").strip()
            if len(entry) == 0:
                entry = str(random_name())
            current_entry['entry'] = entry
            current_entry['document'] = subject
            f.write(b"---\n")
            f.write(yaml.round_trip_dump(current_entry).encode('utf-8'))
            f.write(b"---\n\n\n")

    if export:
        cprint("Wrote %s" % config['output_file'],"green")
        exit(0)

    if show_stats:
        print_stats(file_contents)
        exit(0)

    with open("/tmp/vimrc.config", "w") as f:
        if load_all:
            f.write(VIMRC)
        else:
            f.write(VIMRC2)


    system("vim -u /tmp/vimrc.config -c WPCLI +startinsert /tmp/temp.txt")

    entry_updates = import_file(config, open("/tmp/temp.txt", 'r').read())

    if entry_updates > 0:
        cprint("Pushing...", "yellow", end='', flush=True)
        try:
            git_push()
            cprint("...ok.", "yellow")
        except:
            cprint("...oh well.", "red")
    else:
        cprint("No updates to push.","yellow")

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

if __name__ == "__main__":
    run(None, "notes", load_all=False,
                  edit_one=False, export=False, show_stats=True)