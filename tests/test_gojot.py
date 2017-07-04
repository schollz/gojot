#!/usr/bin/env python
# -*- coding: utf-8 -*-

"""
test_gojot
----------------------------------

Tests for `gojot` module.
"""

import pytest

from contextlib import contextmanager
from click.testing import CliRunner
from os.path import isfile, isdir
from shutil import rmtree

from gojot import gojot
from gojot import cli


@pytest.fixture
def response():
    """Sample pytest fixture.
    See more at: http://doc.pytest.org/en/latest/fixture.html
    """
    # import requests
    # return requests.get('https://github.com/audreyr/cookiecutter-pypackage')


def test_content(response):
    """Sample pytest test function with the pytest fixture as an argument.
    """
    # from bs4 import BeautifulSoup
    # assert 'GitHub' in BeautifulSoup(response.content).title.string


def test_gojot():
    assert gojot.encode_str('teststring', 'testsalt') == '3of9iviQfPizfBuLhBS0'
    assert gojot.decode_str('3of9iviQfPizfBuLhBS0', 'testsalt') == 'teststring'
    assert '_' in gojot.random_name()
    gojot.git_clone('https://github.com/schollz/test5.git')
    assert isdir('test5') == True
    rmtree('test5')
    datas = gojot.parse_entries({'salt':'asldkfjalkdsfj'},"""
---
time: '2017-06-07 09:36:02'
entry: cool_entry1
---
Some text
---
time: '2017-06-07 09:36:03'
entry: cool_entry2
---
Some other text
""")
    assert len(datas) == 2
    assert datas[1] == {'hash': '911e2ea1b3fc4ffba9ed6ca14acaa506', 'meta': {
            'entry': 'cool_entry2', 'time': '2017-06-07 09:36:03'}, 'text': 'Some other text'}
