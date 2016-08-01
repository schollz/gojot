import os
import os.path
import sys
import socket
import shlex
import time
import json
import hashlib
import datetime
import collections

import argparse
import getpass
import subprocess
from pkg_resources import get_distribution

__version__ = "script"
try:
    __version__ = get_distribution('sdees').version
except:
    pass  # user is using as script

REMOTE_SERVER = "www.google.com"
VIM_COMMAND = "vim +100000000 +WP -c 'cal cursor(10000000000000,5000)' -c 'startinsert'"
DATA_PATH = os.path.expanduser('~')


def is_connected():
    try:
        # see if we can resolve the host name -- tells us if there is
        # a DNS listening
        host = socket.gethostbyname(REMOTE_SERVER)
        # connect to the host -- tells us if the host is actually
        # reachable
        s = socket.create_connection((host, 80), 2)
        return True
    except:
        pass
    return False


def get_new_password():
    password1 = '?'
    password2 = '??'
    while password1 != password2:
        password1 = getpass.getpass(prompt='Enter password: ')
        password2 = getpass.getpass(prompt='Enter password again: ')
        if password1 != password2:
            print("Passwords do not match, please try again.")
    return password1


def split_entries_from_file(filename):
    text = open(filename, 'r').read() + "\n"
    return split_entries(text)


def split_entries(text):
    entries = {}
    entry = ""
    entryTime = -1
    for a in text.splitlines():
        if len(a) > 11:
            if a[4] == '-' and a[7] == '-' and a[10] == ' ':
                if len(entry.strip()) > 0:
                    if entryTime not in entries:
                        entries[entryTime] = entry
                    else:
                        entries[entryTime] += entry
                entry = ""
                try:
                    t = datetime.datetime.strptime(
                        ' '.join(a.split()[0:2]), "%Y-%m-%d %H:%M")
                    entryTime = int(
                        (t - datetime.datetime(1970, 1, 1)).total_seconds())
                except:
                    entryTime = -1
        entry += a.strip() + "\n"
    if len(entry) > 0:
        if entryTime not in entries:
            entries[entryTime] = entry
        else:
            entries[entryTime] += entry
    od = collections.OrderedDict(sorted(entries.items()))
    entryArray = []
    for entry in entries:
        entryArray.append(entries[entry])
    return entryArray


def write_entry(filename, entry, password):
    hashOfEntry = str(hashlib.sha224(entry.encode('utf-8')).hexdigest())

    if not os.path.exists(os.path.join(DATA_PATH, '.sdeestemp', filename)):
        os.makedirs(os.path.join(DATA_PATH, '.sdeestemp', filename))

    with open(os.path.join(DATA_PATH, '.sdeestemp', 'temp'), 'w') as f:
        f.write(entry)

    cmd = 'gpg -q --batch --yes --no-use-agent --passphrase %s --symmetric --cipher-algo AES256 -o %s %s' % (
        password, os.path.join(DATA_PATH, '.sdeestemp', filename, hashOfEntry), os.path.join(DATA_PATH, '.sdeestemp', 'temp'))
    process = subprocess.Popen(
        cmd, stdout=subprocess.PIPE, stderr=subprocess.PIPE, shell=True)
    output, error = process.communicate()
    if len(error) > 0:
        print(error)
        clean_up()
        sys.exit(1)
    os.remove(os.path.join(DATA_PATH, '.sdeestemp', 'temp'))
    print("Wrote %s." % hashOfEntry)


def import_file(filename):
    entries = split_entries_from_file(filename)
    password = get_new_password()
    for entry in entries:
        write_entry(filename, entry, password)


def clean_up():
    filesToClean = ['temp', 'temp2', 'temp_copy', 'temp_diff', 'tempEntry']
    for fileToClean in filesToClean:
        if os.path.exists(os.path.join(DATA_PATH, '.sdeestemp', fileToClean)):
            os.remove(os.path.join(DATA_PATH, '.sdeestemp', fileToClean))


def check_prereqs():
    prereqs = ["gpg", "rsync"]
    for prereq in prereqs:
        command = "hash " + prereq  # the shell command
        process = subprocess.Popen(
            command, stdout=subprocess.PIPE, stderr=subprocess.PIPE, shell=True)
        output, error = process.communicate()
        if len(error) > 0:
            print("Need to install " + prereq)
            sys.exit(1)


def sync_down(server):
    if is_connected() and server != None:
        print("Syncing down...")
        cmd = "rsync --ignore-errors -arq --update %s:.sdeestemp/ %s/" % (
            server, os.path.join(DATA_PATH, '.sdeestemp'))
        rsync = subprocess.Popen(
            cmd.split(), stdout=subprocess.PIPE, stderr=subprocess.STDOUT)
        stdout, nothing = rsync.communicate()
        print(stdout, nothing)


def sync_up(server):
    clean_up()
    if is_connected() and server != None:
        print("Syncing up...")
        cmd = "rsync --ignore-errors -arq --update %s/ %s:.sdeestemp/" % (
            os.path.join(DATA_PATH, '.sdeestemp'), server)
        rsync = subprocess.Popen(
            cmd.split(), stdout=subprocess.PIPE, stderr=subprocess.STDOUT)
        stdout, nothing = rsync.communicate()
        print(stdout, nothing)


def is_encrypted(dfile):
    if not os.path.exists(os.path.join(DATA_PATH, '.sdeestemp', dfile)):
        return False
    command = "file %s" % (
        os.path.join(DATA_PATH, '.sdeestemp', dfile))
    process = subprocess.Popen(
        command, stdout=subprocess.PIPE, stderr=subprocess.PIPE, shell=True)
    output, error = process.communicate()
    if "encrypted data" in output.decode():
        return True
    else:
        return False


def set_up():
    server = None
    syncedUp = False
    parser = argparse.ArgumentParser(prog='sdees')
    parser.add_argument("-ls", "--list", help="list available files",
                        action="store_true")
    parser.add_argument("-l", "--local", help="work locally",
                        action="store_true")
    parser.add_argument("-n", "--nodate", help="don't add date",
                        action="store_true")
    parser.add_argument("-e", "--edit", help="edit full document",
                        action="store_true")
    parser.add_argument('--editor', help='change editor')
    parser.add_argument('--importfile', help='import file')
    parser.add_argument("-u", "--update", help="update sdees",
                        action="store_true")
    parser.add_argument('newfile', nargs='?', help='work on a new file')
    args = parser.parse_args()

    if not os.path.exists(os.path.join(DATA_PATH, '.sdeestemp')):
        os.makedirs(os.path.join(DATA_PATH, '.sdeestemp'))

    if args.importfile != None:
        print("Importing %s..." % args.importfile)
        import_file(args.importfile)
        sys.exit(1)

    if args.update:
        os.chdir(os.path.join(DATA_PATH, '.sdeestemp'))
        os.system('git clone https://github.com/schollz/sdees.git')
        os.chdir('sdees')
        os.system('python3 setup.py install')
        os.chdir('../')
        os.system('rm -rf sdees')
        sys.exit(1)

    if args.list:
        os.chdir(os.path.join(DATA_PATH, '.sdeestemp'))
        print("\nAvailable files:")
        os.system(
            'ls -lSht | grep -v config.json | grep -v "total "')
        sys.exit(1)

    # Try to download config.json if doesn't exist
    if not os.path.exists(os.path.join(DATA_PATH, '.sdeestemp', 'config.json')):
        server = input("Enter host@server (make sure to ssh-copy-id first): ")
        dnsaddress = server.split('@')[1]
        address = socket.gethostbyname(dnsaddress)
        server = server.replace(dnsaddress, address)
        if args.local == False:
            sync_down(server)
            syncedUp = True

    # config still doesn't exist, make it
    if not os.path.exists(os.path.join(DATA_PATH, '.sdeestemp', 'config.json')):
        newfile = args.newfile
        if newfile == None:
            newfile = input("Enter a new file name: ")
        files = []
        files.append(newfile)
        config = {"server": server, "files": files, "editor": VIM_COMMAND}
        with open(os.path.join(DATA_PATH, '.sdeestemp', 'config.json'), 'w') as f:
            f.write(json.dumps(config, indent=2))

    # Load config
    config = json.load(
        open(os.path.join(DATA_PATH, '.sdeestemp', 'config.json'), 'r'))

    if 'editor' not in config:
        config['editor'] = VIM_COMMAND
    if args.editor != None:
        print("Editor changed to %s" % args.editor)
        config['editor'] = args.editor
        if 'vim' == args.editor:
            config['editor'] = VIM_COMMAND
        with open(os.path.join(DATA_PATH, '.sdeestemp', 'config.json'), 'w') as f:
            f.write(json.dumps(config, indent=2))

    if not syncedUp and args.local == False:
        sync_down(config['server'])

    # Use the argument file as the default, or add it if doesn't exist
    if args.newfile != None:
        try:
            config['files'].remove(args.newfile)
        except:
            pass
        config['files'] = [args.newfile] + config['files']
        with open(os.path.join(DATA_PATH, '.sdeestemp', 'config.json'), 'w') as f:
            f.write(json.dumps(config, indent=2))

    if args.local == True:
        config['server'] = None
    return args, {'server': config['server'], 'file': config['files'][0], 'editor': config['editor']}


def main(args=None):
    print("sdees, version " + __version__)
    password = None
    check_prereqs()
    args, config = set_up()
    clean_up()

    # If encrypted, do the decryption
    if is_encrypted(config['file']):
        passwordNotAccepted = True
        while passwordNotAccepted:
            password = getpass.getpass(prompt='Enter passphrase: ')
            cmd = 'gpg -q --no-use-agent --passphrase %s -d -o %s %s' % (password, os.path.join(
                DATA_PATH, '.sdeestemp', 'temp'), os.path.join(DATA_PATH, '.sdeestemp', config['file']))
            process = subprocess.Popen(
                cmd, stdout=subprocess.PIPE, stderr=subprocess.PIPE, shell=True)
            output, error = process.communicate()
            if len(error) == 0:
                passwordNotAccepted = False
    else:
        # Copy file to temp if not encrypted
        cmd = "cp %s %s" % (os.path.join(DATA_PATH, '.sdeestemp', config['file']), os.path.join(
            DATA_PATH, '.sdeestemp', 'temp'))
        process = subprocess.Popen(
            cmd, stdout=subprocess.PIPE, stderr=subprocess.PIPE, shell=True)
        output, error = process.communicate()

    # Copy file for diffing
    cmd = "cp %s %s" % (os.path.join(DATA_PATH, '.sdeestemp', 'temp'), os.path.join(
        DATA_PATH, '.sdeestemp', 'temp_copy'))
    process = subprocess.Popen(
        cmd, stdout=subprocess.PIPE, stderr=subprocess.PIPE, shell=True)
    output, error = process.communicate()

    timeString = time.strftime("%Y-%m-%d %H:%M  ")
    if args.nodate:
        timeString = ""
    if args.edit:
        # Add new entry directly to the file
        with open(os.path.join(DATA_PATH, '.sdeestemp', 'temp'), 'a') as f:
            f.write("\n\n" + timeString)
        # Open it in editor to write
        os.system("%s %s" % (config['editor'],
                             os.path.join(DATA_PATH, '.sdeestemp', 'temp')))
    else:
        # Add new entry in a seperate file
        with open(os.path.join(DATA_PATH, '.sdeestemp', 'tempEntry'), 'a') as f:
            f.write(timeString)
        # Open it in editor to write
        os.system("%s %s" % (config['editor'],
                             os.path.join(DATA_PATH, '.sdeestemp', 'tempEntry')))
        # append the entry to the file
        with open(os.path.join(DATA_PATH, '.sdeestemp', 'temp'), 'a') as f:
            with open(os.path.join(DATA_PATH, '.sdeestemp', 'tempEntry'), 'r') as f2:
                tempEntry = f2.read()
                if len(tempEntry) < 22:
                    print("No data appended.")
                else:
                    f.write("\n\n" + tempEntry)

    # Write a diff
    cmd = "diff %s %s" % (os.path.join(DATA_PATH, '.sdeestemp',
                                       'temp_copy'), os.path.join(DATA_PATH, '.sdeestemp', 'temp'))
    process = subprocess.Popen(
        cmd, stdout=subprocess.PIPE, stderr=subprocess.PIPE, shell=True)
    output, error = process.communicate()
    diffFile = config['file'] + '.' + str(hashlib.sha224(output).hexdigest())
    if len(output) > 1:
        with open(os.path.join(DATA_PATH, '.sdeestemp', 'temp_diff'), 'w') as f:
            f.write(output.decode())
    else:
        diffFile = ""

    # Encrypt main file
    if password == None:
        password = get_new_password()
    cmd = 'gpg -q --no-use-agent --passphrase %s --symmetric --cipher-algo AES256 -o %s %s' % (
        password, os.path.join(DATA_PATH, '.sdeestemp', 'temp2'), os.path.join(DATA_PATH, '.sdeestemp', 'temp'))
    process = subprocess.Popen(
        cmd, stdout=subprocess.PIPE, stderr=subprocess.PIPE, shell=True)
    output, error = process.communicate()
    if len(error) > 0:
        print(error)
        clean_up()
        sys.exit(1)

    # Encrypt diff file
    if len(diffFile) > 0:
        if password == None:
            password = getpass.getpass(prompt='Enter password: ')
        cmd = 'gpg -q --no-use-agent --passphrase %s --symmetric --cipher-algo AES256 -o %s %s' % (
            password, os.path.join(DATA_PATH, '.sdeestemp', 'diffs', diffFile), os.path.join(DATA_PATH, '.sdeestemp', 'temp_diff'))
        process = subprocess.Popen(
            cmd, stdout=subprocess.PIPE, stderr=subprocess.PIPE, shell=True)
        output, error = process.communicate()
        if len(error) > 0:
            print(error)
            clean_up()
            sys.exit(1)

    # overwrite the main file
    os.system('mv %s %s' % (os.path.join(DATA_PATH, '.sdeestemp', 'temp2'),
                            os.path.join(DATA_PATH, '.sdeestemp', config['file'])))

    sync_up(config['server'])
    clean_up()

if __name__ == "__main__":
    # import_file('test.txt')
    clean_up()
    main()
