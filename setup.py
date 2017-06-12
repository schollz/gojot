#!/usr/bin/env python
# -*- coding: utf-8 -*-

from setuptools import setup

with open('README.rst') as readme_file:
    readme = readme_file.read()

with open('HISTORY.rst') as history_file:
    history = history_file.read()

requirements = [
    'Click>=6.0',
    'pick',
    'hashids',
    'termcolor',
    'tqdm',
    'ruamel.yaml',
    'humanize',
    # TODO: put package requirements here
]

test_requirements = [
    'pick',
    'hashids',
    'termcolor',
    'tqdm',
    'ruamel.yaml',
    'humanize',
    # TODO: put package test requirements here
]

setup(
    name='gojot',
    version='3.0.2',
    description="A command-line journal that is distributed and encrypted, making it easy to jot notes.",
    long_description=readme + '\n\n' + history,
    author="Zack Scholl",
    author_email='zack.scholl@gmail.com',
    url='https://github.com/schollz/gojot',
    packages=[
        'gojot',
    ],
    package_dir={'gojot':
                 'gojot'},
    entry_points={
        'console_scripts': [
            'gojot=gojot.cli:main'
        ]
    },
    include_package_data=True,
    install_requires=requirements,
    license="MIT license",
    zip_safe=False,
    keywords='gojot',
    classifiers=[
        'Development Status :: 2 - Pre-Alpha',
        'Intended Audience :: Developers',
        'License :: OSI Approved :: MIT License',
        'Natural Language :: English',
        'Programming Language :: Python :: 3.3',
        'Programming Language :: Python :: 3.4',
        'Programming Language :: Python :: 3.5',
    ],
    test_suite='tests',
    tests_require=test_requirements
)
