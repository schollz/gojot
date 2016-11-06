About
=========

Why sdees?
--------------

**Problem:** I would like to write into a document (e.g., a journal with
time stamped entries) and also have it be available on all computers I
use, and I want it to be locally encrypted. Unfortunately, typical cloud
storage synchronization utilities (e.g. Dropbox) do not play well with
syncing a single document that is encrypted locally.

**Solution:** The basic solution is to create a lot of files - each
encrypted - where each file contains the text for one entry in the
document. You don’t need any special tool for this - you can just use a
text editor, ``gpg``, and synchronization software like ``git`` or
Dropbox.

My program, *sdees*, just makes this solution easier to attain.
*sdees* comes as a single executable file with the text-editor and
``gpg`` bulit-in - the only system requirement is the installation of
``git`` (which is pretty easy to get on any system). The other benefit
of *sdees* is that it will automatically combine all the time-stamped
entries so it appears that you are editing a single document, and it
will also resolve merges that can occur if you edit the same entry
offline on two computers.
