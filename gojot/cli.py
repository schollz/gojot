# -*- coding: utf-8 -*-

import click
from gojot import *


@click.command()
@click.argument('repo', nargs=1)
def main(repo, args=None):
    """Console script for gojot"""
    click.echo("Replace this message by putting your code into "
               "gojot.cli.main")
    click.echo("See click documentation at http://click.pocoo.org/")
    run(repo)


if __name__ == "__main__":
    main()
