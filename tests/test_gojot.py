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
    datas = gojot.parse_entries("""
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
    assert datas[1] == {'hash': 'f0e124f2a81affb35281aa971ff5c318', 'meta': {
        'entry': 'cool_entry2', 'time': '2017-06-07 09:36:03'}, 'text': 'Some other text'}


def test_command_line_interface():
    runner = CliRunner()
    result = runner.invoke(cli.main)
    # assert result.exit_code == 0
    assert 'Missing argument' in result.output
    help_result = runner.invoke(cli.main, ['--help'])
    assert help_result.exit_code == 0
    assert 'Show this message and exit.' in help_result.output
