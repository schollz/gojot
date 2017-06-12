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
from random import choice
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



def random_name():
    left = ["admiring",
                    "adoring",
                    "affectionate",
                    "agitated",
                    "amazing",
                    "angry",
                    "awesome",
                    "blissful",
                    "boring",
                    "brave",
                    "clever",
                    "cocky",
                    "compassionate",
                    "competent",
                    "condescending",
                    "confident",
                    "cranky",
                    "dazzling",
                    "determined",
                    "distracted",
                    "dreamy",
                    "eager",
                    "ecstatic",
                    "elastic",
                    "elated",
                    "elegant",
                    "eloquent",
                    "epic",
                    "fervent",
                    "festive",
                    "flamboyant",
                    "focused",
                    "friendly",
                    "frosty",
                    "gallant",
                    "gifted",
                    "goofy",
                    "gracious",
                    "happy",
                    "hardcore",
                    "heuristic",
                    "hopeful",
                    "hungry",
                    "infallible",
                    "inspiring",
                    "jolly",
                    "jovial",
                    "keen",
                    "kind",
                    "laughing",
                    "loving",
                    "lucid",
                    "mystifying",
                    "modest",
                    "musing",
                    "naughty",
                    "nervous",
                    "nifty",
                    "nostalgic",
                    "objective",
                    "optimistic",
                    "peaceful",
                    "pedantic",
                    "pensive",
                    "practical",
                    "priceless",
                    "quirky",
                    "quizzical",
                    "relaxed",
                    "reverent",
                    "romantic",
                    "sad",
                    "serene",
                    "sharp",
                    "silly",
                    "sleepy",
                    "stoic",
                    "stupefied",
                    "suspicious",
                    "tender",
                    "thirsty",
                    "trusting",
                    "unruffled",
                    "upbeat",
                    "vibrant",
                    "vigilant",
                    "vigorous",
                    "wizardly",
                    "wonderful",
                    "xenodochial",
                    "youthful",
                    "zealous",
                    "zen"]

    right = [
            # Muhammad ibn Jābir al-Ḥarrānī al-Battānī was a founding father of astronomy. https://en.wikipedia.org/wiki/Mu%E1%B8%A5ammad_ibn_J%C4%81bir_al-%E1%B8%A4arr%C4%81n%C4%AB_al-Batt%C4%81n%C4%AB
            "albattani",

            # Frances E. Allen, became the first female IBM Fellow in 1989. In 2006, she became the first female recipient of the ACM's Turing Award. https://en.wikipedia.org/wiki/Frances_E._Allen
            "allen",

            # June Almeida - Scottish virologist who took the first pictures of the rubella virus - https://en.wikipedia.org/wiki/June_Almeida
            "almeida",

            # Maria Gaetana Agnesi - Italian mathematician, philosopher, theologian and humanitarian. She was the first woman to write a mathematics handbook and the first woman appointed as a Mathematics Professor at a University. https://en.wikipedia.org/wiki/Maria_Gaetana_Agnesi
            "agnesi",

            # Archimedes was a physicist, engineer and mathematician who invented too many things to list them here. https://en.wikipedia.org/wiki/Archimedes
            "archimedes",

            # Maria Ardinghelli - Italian translator, mathematician and physicist - https://en.wikipedia.org/wiki/Maria_Ardinghelli
            "ardinghelli",

            # Aryabhata - Ancient Indian mathematician-astronomer during 476-550 CE https://en.wikipedia.org/wiki/Aryabhata
            "aryabhata",

            # Wanda Austin - Wanda Austin is the President and CEO of The Aerospace Corporation, a leading architect for the US security space programs. https://en.wikipedia.org/wiki/Wanda_Austin
            "austin",

            # Charles Babbage invented the concept of a programmable computer. https://en.wikipedia.org/wiki/Charles_Babbage.
            "babbage",

            # Stefan Banach - Polish mathematician, was one of the founders of modern functional analysis. https://en.wikipedia.org/wiki/Stefan_Banach
            "banach",

            # John Bardeen co-invented the transistor - https://en.wikipedia.org/wiki/John_Bardeen
            "bardeen",

            # Jean Bartik, born Betty Jean Jennings, was one of the original programmers for the ENIAC computer. https://en.wikipedia.org/wiki/Jean_Bartik
            "bartik",

            # Laura Bassi, the world's first female professor https://en.wikipedia.org/wiki/Laura_Bassi
            "bassi",

            # Hugh Beaver, British engineer, founder of the Guinness Book of World Records https://en.wikipedia.org/wiki/Hugh_Beaver
            "beaver",

            # Alexander Graham Bell - an eminent Scottish-born scientist, inventor, engineer and innovator who is credited with inventing the first practical telephone - https://en.wikipedia.org/wiki/Alexander_Graham_Bell
            "bell",

            # Karl Friedrich Benz - a German automobile engineer. Inventor of the first practical motorcar. https://en.wikipedia.org/wiki/Karl_Benz
            "benz",

            # Homi J Bhabha - was an Indian nuclear physicist, founding director, and professor of physics at the Tata Institute of Fundamental Research. Colloquially known as "father of Indian nuclear programme"- https://en.wikipedia.org/wiki/Homi_J._Bhabha
            "bhabha",

            # Bhaskara II - Ancient Indian mathematician-astronomer whose work on calculus predates Newton and Leibniz by over half a millennium - https://en.wikipedia.org/wiki/Bh%C4%81skara_II#Calculus
            "bhaskara",

            # Elizabeth Blackwell - American doctor and first American woman to receive a medical degree - https://en.wikipedia.org/wiki/Elizabeth_Blackwell
            "blackwell",

            # Niels Bohr is the father of quantum theory. https://en.wikipedia.org/wiki/Niels_Bohr.
            "bohr",

            # Kathleen Booth, she's credited with writing the first assembly language. https://en.wikipedia.org/wiki/Kathleen_Booth
            "booth",

            # Anita Borg - Anita Borg was the founding director of the Institute for Women and Technology (IWT). https://en.wikipedia.org/wiki/Anita_Borg
            "borg",

            # Satyendra Nath Bose - He provided the foundation for Bose–Einstein statistics and the theory of the Bose–Einstein condensate. - https://en.wikipedia.org/wiki/Satyendra_Nath_Bose
            "bose",

            # Evelyn Boyd Granville - She was one of the first African-American woman to receive a Ph.D. in mathematics; she earned it in 1949 from Yale University. https://en.wikipedia.org/wiki/Evelyn_Boyd_Granville
            "boyd",

            # Brahmagupta - Ancient Indian mathematician during 598-670 CE who gave rules to compute with zero - https://en.wikipedia.org/wiki/Brahmagupta#Zero
            "brahmagupta",

            # Walter Houser Brattain co-invented the transistor - https://en.wikipedia.org/wiki/Walter_Houser_Brattain
            "brattain",

            # Emmett Brown invented time travel. https://en.wikipedia.org/wiki/Emmett_Brown (thanks Brian Goff)
            "brown",

            # Rachel Carson - American marine biologist and conservationist, her book Silent Spring and other writings are credited with advancing the global environmental movement. https://en.wikipedia.org/wiki/Rachel_Carson
            "carson",

            # Subrahmanyan Chandrasekhar - Astrophysicist known for his mathematical theory on different stages and evolution in structures of the stars. He has won nobel prize for physics - https://en.wikipedia.org/wiki/Subrahmanyan_Chandrasekhar
            "chandrasekhar",

            # Claude Shannon - The father of information theory and founder of digital circuit design theory. (https://en.wikipedia.org/wiki/Claude_Shannon)
            "shannon",

            # Joan Clarke - Bletchley Park code breaker during the Second World War who pioneered techniques that remained top secret for decades. Also an accomplished numismatist https://en.wikipedia.org/wiki/Joan_Clarke
            "clarke",

            # Jane Colden - American botanist widely considered the first female American botanist - https://en.wikipedia.org/wiki/Jane_Colden
            "colden",

            # Gerty Theresa Cori - American biochemist who became the third woman—and first American woman—to win a Nobel Prize in science, and the first woman to be awarded the Nobel Prize in Physiology or Medicine. Cori was born in Prague. https://en.wikipedia.org/wiki/Gerty_Cori
            "cori",

            # Seymour Roger Cray was an American electrical engineer and supercomputer architect who designed a series of computers that were the fastest in the world for decades. https://en.wikipedia.org/wiki/Seymour_Cray
            "cray",

            # This entry reflects a husband and wife team who worked together:
            # Joan Curran was a Welsh scientist who developed radar and invented chaff, a radar countermeasure. https://en.wikipedia.org/wiki/Joan_Curran
            # Samuel Curran was an Irish physicist who worked alongside his wife during WWII and invented the proximity fuse. https://en.wikipedia.org/wiki/Samuel_Curran
            "curran",

            # Marie Curie discovered radioactivity. https://en.wikipedia.org/wiki/Marie_Curie.
            "curie",

            # Charles Darwin established the principles of natural evolution. https://en.wikipedia.org/wiki/Charles_Darwin.
            "darwin",

            # Leonardo Da Vinci invented too many things to list here. https://en.wikipedia.org/wiki/Leonardo_da_Vinci.
            "davinci",

            # Edsger Wybe Dijkstra was a Dutch computer scientist and mathematical scientist. https://en.wikipedia.org/wiki/Edsger_W._Dijkstra.
            "dijkstra",

            # Donna Dubinsky - played an integral role in the development of personal digital assistants (PDAs) serving as CEO of Palm, Inc. and co-founding Handspring. https://en.wikipedia.org/wiki/Donna_Dubinsky
            "dubinsky",

            # Annie Easley - She was a leading member of the team which developed software for the Centaur rocket stage and one of the first African-Americans in her field. https://en.wikipedia.org/wiki/Annie_Easley
            "easley",

            # Thomas Alva Edison, prolific inventor https://en.wikipedia.org/wiki/Thomas_Edison
            "edison",

            # Albert Einstein invented the general theory of relativity. https://en.wikipedia.org/wiki/Albert_Einstein
            "einstein",

            # Gertrude Elion - American biochemist, pharmacologist and the 1988 recipient of the Nobel Prize in Medicine - https://en.wikipedia.org/wiki/Gertrude_Elion
            "elion",

            # Douglas Engelbart gave the mother of all demos: https://en.wikipedia.org/wiki/Douglas_Engelbart
            "engelbart",

            # Euclid invented geometry. https://en.wikipedia.org/wiki/Euclid
            "euclid",

            # Leonhard Euler invented large parts of modern mathematics. https://de.wikipedia.org/wiki/Leonhard_Euler
            "euler",

            # Pierre de Fermat pioneered several aspects of modern mathematics. https://en.wikipedia.org/wiki/Pierre_de_Fermat
            "fermat",

            # Enrico Fermi invented the first nuclear reactor. https://en.wikipedia.org/wiki/Enrico_Fermi.
            "fermi",

            # Richard Feynman was a key contributor to quantum mechanics and particle physics. https://en.wikipedia.org/wiki/Richard_Feynman
            "feynman",

            # Benjamin Franklin is famous for his experiments in electricity and the invention of the lightning rod.
            "franklin",

            # Galileo was a founding father of modern astronomy, and faced politics and obscurantism to establish scientific truth.  https://en.wikipedia.org/wiki/Galileo_Galilei
            "galileo",

            # William Henry "Bill" Gates III is an American business magnate, philanthropist, investor, computer programmer, and inventor. https://en.wikipedia.org/wiki/Bill_Gates
            "gates",

            # Adele Goldberg, was one of the designers and developers of the Smalltalk language. https://en.wikipedia.org/wiki/Adele_Goldberg_(computer_scientist)
            "goldberg",

            # Adele Goldstine, born Adele Katz, wrote the complete technical description for the first electronic digital computer, ENIAC. https://en.wikipedia.org/wiki/Adele_Goldstine
            "goldstine",

            # Shafi Goldwasser is a computer scientist known for creating theoretical foundations of modern cryptography. Winner of 2012 ACM Turing Award. https://en.wikipedia.org/wiki/Shafi_Goldwasser
            "goldwasser",

            # James Golick, all around gangster.
            "golick",

            # Jane Goodall - British primatologist, ethologist, and anthropologist who is considered to be the world's foremost expert on chimpanzees - https://en.wikipedia.org/wiki/Jane_Goodall
            "goodall",

            # Lois Haibt - American computer scientist, part of the team at IBM that developed FORTRAN - https://en.wikipedia.org/wiki/Lois_Haibt
            "haibt",

            # Margaret Hamilton - Director of the Software Engineering Division of the MIT Instrumentation Laboratory, which developed on-board flight software for the Apollo space program. https://en.wikipedia.org/wiki/Margaret_Hamilton_(scientist)
            "hamilton",

            # Stephen Hawking pioneered the field of cosmology by combining general relativity and quantum mechanics. https://en.wikipedia.org/wiki/Stephen_Hawking
            "hawking",

            # Werner Heisenberg was a founding father of quantum mechanics. https://en.wikipedia.org/wiki/Werner_Heisenberg
            "heisenberg",

            # Grete Hermann was a German philosopher noted for her philosophical work on the foundations of quantum mechanics. https://en.wikipedia.org/wiki/Grete_Hermann
            "hermann",

            # Jaroslav Heyrovský was the inventor of the polarographic method, father of the electroanalytical method, and recipient of the Nobel Prize in 1959. His main field of work was polarography. https://en.wikipedia.org/wiki/Jaroslav_Heyrovsk%C3%BD
            "heyrovsky",

            # Dorothy Hodgkin was a British biochemist, credited with the development of protein crystallography. She was awarded the Nobel Prize in Chemistry in 1964. https://en.wikipedia.org/wiki/Dorothy_Hodgkin
            "hodgkin",

            # Erna Schneider Hoover revolutionized modern communication by inventing a computerized telephone switching method. https://en.wikipedia.org/wiki/Erna_Schneider_Hoover
            "hoover",

            # Grace Hopper developed the first compiler for a computer programming language and  is credited with popularizing the term "debugging" for fixing computer glitches. https://en.wikipedia.org/wiki/Grace_Hopper
            "hopper",

            # Frances Hugle, she was an American scientist, engineer, and inventor who contributed to the understanding of semiconductors, integrated circuitry, and the unique electrical principles of microscopic materials. https://en.wikipedia.org/wiki/Frances_Hugle
            "hugle",

            # Hypatia - Greek Alexandrine Neoplatonist philosopher in Egypt who was one of the earliest mothers of mathematics - https://en.wikipedia.org/wiki/Hypatia
            "hypatia",

            # Mary Jackson, American mathematician and aerospace engineer who earned the highest title within NASA's engineering department - https://en.wikipedia.org/wiki/Mary_Jackson_(engineer)
            "jackson",

            # Yeong-Sil Jang was a Korean scientist and astronomer during the Joseon Dynasty; he invented the first metal printing press and water gauge. https://en.wikipedia.org/wiki/Jang_Yeong-sil
            "jang",

            # Betty Jennings - one of the original programmers of the ENIAC. https://en.wikipedia.org/wiki/ENIAC - https://en.wikipedia.org/wiki/Jean_Bartik
            "jennings",

            # Mary Lou Jepsen, was the founder and chief technology officer of One Laptop Per Child (OLPC), and the founder of Pixel Qi. https://en.wikipedia.org/wiki/Mary_Lou_Jepsen
            "jepsen",

            # Katherine Coleman Goble Johnson - American physicist and mathematician contributed to the NASA. https://en.wikipedia.org/wiki/Katherine_Johnson
            "johnson",

            # Irène Joliot-Curie - French scientist who was awarded the Nobel Prize for Chemistry in 1935. Daughter of Marie and Pierre Curie. https://en.wikipedia.org/wiki/Ir%C3%A8ne_Joliot-Curie
            "joliot",

            # Karen Spärck Jones came up with the concept of inverse document frequency, which is used in most search engines today. https://en.wikipedia.org/wiki/Karen_Sp%C3%A4rck_Jones
            "jones",

            # A. P. J. Abdul Kalam - is an Indian scientist aka Missile Man of India for his work on the development of ballistic missile and launch vehicle technology - https://en.wikipedia.org/wiki/A._P._J._Abdul_Kalam
            "kalam",

            # Susan Kare, created the icons and many of the interface elements for the original Apple Macintosh in the 1980s, and was an original employee of NeXT, working as the Creative Director. https://en.wikipedia.org/wiki/Susan_Kare
            "kare",

            # Mary Kenneth Keller, Sister Mary Kenneth Keller became the first American woman to earn a PhD in Computer Science in 1965. https://en.wikipedia.org/wiki/Mary_Kenneth_Keller
            "keller",

            # Johannes Kepler, German astronomer known for his three laws of planetary motion - https://en.wikipedia.org/wiki/Johannes_Kepler
            "kepler",

            # Har Gobind Khorana - Indian-American biochemist who shared the 1968 Nobel Prize for Physiology - https://en.wikipedia.org/wiki/Har_Gobind_Khorana
            "khorana",

            # Jack Kilby invented silicone integrated circuits and gave Silicon Valley its name. - https://en.wikipedia.org/wiki/Jack_Kilby
            "kilby",

            # Maria Kirch - German astronomer and first woman to discover a comet - https://en.wikipedia.org/wiki/Maria_Margarethe_Kirch
            "kirch",

            # Donald Knuth - American computer scientist, author of "The Art of Computer Programming" and creator of the TeX typesetting system. https://en.wikipedia.org/wiki/Donald_Knuth
            "knuth",

            # Sophie Kowalevski - Russian mathematician responsible for important original contributions to analysis, differential equations and mechanics - https://en.wikipedia.org/wiki/Sofia_Kovalevskaya
            "kowalevski",

            # Marie-Jeanne de Lalande - French astronomer, mathematician and cataloguer of stars - https://en.wikipedia.org/wiki/Marie-Jeanne_de_Lalande
            "lalande",

            # Hedy Lamarr - Actress and inventor. The principles of her work are now incorporated into modern Wi-Fi, CDMA and Bluetooth technology. https://en.wikipedia.org/wiki/Hedy_Lamarr
            "lamarr",

            # Leslie B. Lamport - American computer scientist. Lamport is best known for his seminal work in distributed systems and was the winner of the 2013 Turing Award. https://en.wikipedia.org/wiki/Leslie_Lamport
            "lamport",

            # Mary Leakey - British paleoanthropologist who discovered the first fossilized Proconsul skull - https://en.wikipedia.org/wiki/Mary_Leakey
            "leakey",

            # Henrietta Swan Leavitt - she was an American astronomer who discovered the relation between the luminosity and the period of Cepheid variable stars. https://en.wikipedia.org/wiki/Henrietta_Swan_Leavitt
            "leavitt",

            # Daniel Lewin -  Mathematician, Akamai co-founder, soldier, 9/11 victim-- Developed optimization techniques for routing traffic on the internet. Died attempting to stop the 9-11 hijackers. https://en.wikipedia.org/wiki/Daniel_Lewin
            "lewin",

            # Ruth Lichterman - one of the original programmers of the ENIAC. https://en.wikipedia.org/wiki/ENIAC - https://en.wikipedia.org/wiki/Ruth_Teitelbaum
            "lichterman",

            # Barbara Liskov - co-developed the Liskov substitution principle. Liskov was also the winner of the Turing Prize in 2008. - https://en.wikipedia.org/wiki/Barbara_Liskov
            "liskov",

            # Ada Lovelace invented the first algorithm. https://en.wikipedia.org/wiki/Ada_Lovelace (thanks James Turnbull)
            "lovelace",

            # Auguste and Louis Lumière - the first filmmakers in history - https://en.wikipedia.org/wiki/Auguste_and_Louis_Lumi%C3%A8re
            "lumiere",

            # Mahavira - Ancient Indian mathematician during 9th century AD who discovered basic algebraic identities - https://en.wikipedia.org/wiki/Mah%C4%81v%C4%ABra_(mathematician)
            "mahavira",

            # Maria Mayer - American theoretical physicist and Nobel laureate in Physics for proposing the nuclear shell model of the atomic nucleus - https://en.wikipedia.org/wiki/Maria_Mayer
            "mayer",

            # John McCarthy invented LISP: https://en.wikipedia.org/wiki/John_McCarthy_(computer_scientist)
            "mccarthy",

            # Barbara McClintock - a distinguished American cytogeneticist, 1983 Nobel Laureate in Physiology or Medicine for discovering transposons. https://en.wikipedia.org/wiki/Barbara_McClintock
            "mcclintock",

            # Malcolm McLean invented the modern shipping container: https://en.wikipedia.org/wiki/Malcom_McLean
            "mclean",

            # Kay McNulty - one of the original programmers of the ENIAC. https://en.wikipedia.org/wiki/ENIAC - https://en.wikipedia.org/wiki/Kathleen_Antonelli
            "mcnulty",

            # Lise Meitner - Austrian/Swedish physicist who was involved in the discovery of nuclear fission. The element meitnerium is named after her - https://en.wikipedia.org/wiki/Lise_Meitner
            "meitner",

            # Carla Meninsky, was the game designer and programmer for Atari 2600 games Dodge 'Em and Warlords. https://en.wikipedia.org/wiki/Carla_Meninsky
            "meninsky",

            # Johanna Mestorf - German prehistoric archaeologist and first female museum director in Germany - https://en.wikipedia.org/wiki/Johanna_Mestorf
            "mestorf",

            # Marvin Minsky - Pioneer in Artificial Intelligence, co-founder of the MIT's AI Lab, won the Turing Award in 1969. https://en.wikipedia.org/wiki/Marvin_Minsky
            "minsky",

            # Maryam Mirzakhani - an Iranian mathematician and the first woman to win the Fields Medal. https://en.wikipedia.org/wiki/Maryam_Mirzakhani
            "mirzakhani",

            # Samuel Morse - contributed to the invention of a single-wire telegraph system based on European telegraphs and was a co-developer of the Morse code - https://en.wikipedia.org/wiki/Samuel_Morse
            "morse",

            # Ian Murdock - founder of the Debian project - https://en.wikipedia.org/wiki/Ian_Murdock
            "murdock",

            # John von Neumann - todays computer architectures are based on the von Neumann architecture. https://en.wikipedia.org/wiki/Von_Neumann_architecture
            "neumann",

            # Isaac Newton invented classic mechanics and modern optics. https://en.wikipedia.org/wiki/Isaac_Newton
            "newton",

            # Florence Nightingale, more prominently known as a nurse, was also the first female member of the Royal Statistical Society and a pioneer in statistical graphics https://en.wikipedia.org/wiki/Florence_Nightingale#Statistics_and_sanitary_reform
            "nightingale",

            # Alfred Nobel - a Swedish chemist, engineer, innovator, and armaments manufacturer (inventor of dynamite) - https://en.wikipedia.org/wiki/Alfred_Nobel
            "nobel",

            # Emmy Noether, German mathematician. Noether's Theorem is named after her. https://en.wikipedia.org/wiki/Emmy_Noether
            "noether",

            # Poppy Northcutt. Poppy Northcutt was the first woman to work as part of NASA’s Mission Control. http://www.businessinsider.com/poppy-northcutt-helped-apollo-astronauts-2014-12?op=1
            "northcutt",

            # Robert Noyce invented silicone integrated circuits and gave Silicon Valley its name. - https://en.wikipedia.org/wiki/Robert_Noyce
            "noyce",

            # Panini - Ancient Indian linguist and grammarian from 4th century CE who worked on the world's first formal system - https://en.wikipedia.org/wiki/P%C4%81%E1%B9%87ini#Comparison_with_modern_formal_systems
            "panini",

            # Ambroise Pare invented modern surgery. https://en.wikipedia.org/wiki/Ambroise_Par%C3%A9
            "pare",

            # Louis Pasteur discovered vaccination, fermentation and pasteurization. https://en.wikipedia.org/wiki/Louis_Pasteur.
            "pasteur",

            # Cecilia Payne-Gaposchkin was an astronomer and astrophysicist who, in 1925, proposed in her Ph.D. thesis an explanation for the composition of stars in terms of the relative abundances of hydrogen and helium. https://en.wikipedia.org/wiki/Cecilia_Payne-Gaposchkin
            "payne",

            # Radia Perlman is a software designer and network engineer and most famous for her invention of the spanning-tree protocol (STP). https://en.wikipedia.org/wiki/Radia_Perlman
            "perlman",

            # Rob Pike was a key contributor to Unix, Plan 9, the X graphic system, utf-8, and the Go programming language. https://en.wikipedia.org/wiki/Rob_Pike
            "pike",

            # Henri Poincaré made fundamental contributions in several fields of mathematics. https://en.wikipedia.org/wiki/Henri_Poincar%C3%A9
            "poincare",

            # Laura Poitras is a director and producer whose work, made possible by open source crypto tools, advances the causes of truth and freedom of information by reporting disclosures by whistleblowers such as Edward Snowden. https://en.wikipedia.org/wiki/Laura_Poitras
            "poitras",

            # Claudius Ptolemy - a Greco-Egyptian writer of Alexandria, known as a mathematician, astronomer, geographer, astrologer, and poet of a single epigram in the Greek Anthology - https://en.wikipedia.org/wiki/Ptolemy
            "ptolemy",

            # C. V. Raman - Indian physicist who won the Nobel Prize in 1930 for proposing the Raman effect. - https://en.wikipedia.org/wiki/C._V._Raman
            "raman",

            # Srinivasa Ramanujan - Indian mathematician and autodidact who made extraordinary contributions to mathematical analysis, number theory, infinite series, and continued fractions. - https://en.wikipedia.org/wiki/Srinivasa_Ramanujan
            "ramanujan",

            # Sally Kristen Ride was an American physicist and astronaut. She was the first American woman in space, and the youngest American astronaut. https://en.wikipedia.org/wiki/Sally_Ride
            "ride",

            # Rita Levi-Montalcini - Won Nobel Prize in Physiology or Medicine jointly with colleague Stanley Cohen for the discovery of nerve growth factor (https://en.wikipedia.org/wiki/Rita_Levi-Montalcini)
            "montalcini",

            # Dennis Ritchie - co-creator of UNIX and the C programming language. - https://en.wikipedia.org/wiki/Dennis_Ritchie
            "ritchie",

            # Wilhelm Conrad Röntgen - German physicist who was awarded the first Nobel Prize in Physics in 1901 for the discovery of X-rays (Röntgen rays). https://en.wikipedia.org/wiki/Wilhelm_R%C3%B6ntgen
            "roentgen",

            # Rosalind Franklin - British biophysicist and X-ray crystallographer whose research was critical to the understanding of DNA - https://en.wikipedia.org/wiki/Rosalind_Franklin
            "rosalind",

            # Meghnad Saha - Indian astrophysicist best known for his development of the Saha equation, used to describe chemical and physical conditions in stars - https://en.wikipedia.org/wiki/Meghnad_Saha
            "saha",

            # Jean E. Sammet developed FORMAC, the first widely used computer language for symbolic manipulation of mathematical formulas. https://en.wikipedia.org/wiki/Jean_E._Sammet
            "sammet",

            # Carol Shaw - Originally an Atari employee, Carol Shaw is said to be the first female video game designer. https://en.wikipedia.org/wiki/Carol_Shaw_(video_game_designer)
            "shaw",

            # Dame Stephanie "Steve" Shirley - Founded a software company in 1962 employing women working from home. https://en.wikipedia.org/wiki/Steve_Shirley
            "shirley",

            # William Shockley co-invented the transistor - https://en.wikipedia.org/wiki/William_Shockley
            "shockley",

            # Françoise Barré-Sinoussi - French virologist and Nobel Prize Laureate in Physiology or Medicine; her work was fundamental in identifying HIV as the cause of AIDS. https://en.wikipedia.org/wiki/Fran%C3%A7oise_Barr%C3%A9-Sinoussi
            "sinoussi",

            # Betty Snyder - one of the original programmers of the ENIAC. https://en.wikipedia.org/wiki/ENIAC - https://en.wikipedia.org/wiki/Betty_Holberton
            "snyder",

            # Frances Spence - one of the original programmers of the ENIAC. https://en.wikipedia.org/wiki/ENIAC - https://en.wikipedia.org/wiki/Frances_Spence
            "spence",

            # Richard Matthew Stallman - the founder of the Free Software movement, the GNU project, the Free Software Foundation, and the League for Programming Freedom. He also invented the concept of copyleft to protect the ideals of this movement, and enshrined this concept in the widely-used GPL (General Public License) for software. https://en.wikiquote.org/wiki/Richard_Stallman
            "stallman",

            # Michael Stonebraker is a database research pioneer and architect of Ingres, Postgres, VoltDB and SciDB. Winner of 2014 ACM Turing Award. https://en.wikipedia.org/wiki/Michael_Stonebraker
            "stonebraker",

            # Janese Swanson (with others) developed the first of the Carmen Sandiego games. She went on to found Girl Tech. https://en.wikipedia.org/wiki/Janese_Swanson
            "swanson",

            # Aaron Swartz was influential in creating RSS, Markdown, Creative Commons, Reddit, and much of the internet as we know it today. He was devoted to freedom of information on the web. https://en.wikiquote.org/wiki/Aaron_Swartz
            "swartz",

            # Bertha Swirles was a theoretical physicist who made a number of contributions to early quantum theory. https://en.wikipedia.org/wiki/Bertha_Swirles
            "swirles",

            # Nikola Tesla invented the AC electric system and every gadget ever used by a James Bond villain. https://en.wikipedia.org/wiki/Nikola_Tesla
            "tesla",

            # Ken Thompson - co-creator of UNIX and the C programming language - https://en.wikipedia.org/wiki/Ken_Thompson
            "thompson",

            # Linus Torvalds invented Linux and Git. https://en.wikipedia.org/wiki/Linus_Torvalds
            "torvalds",

            # Alan Turing was a founding father of computer science. https://en.wikipedia.org/wiki/Alan_Turing.
            "turing",

            # Varahamihira - Ancient Indian mathematician who discovered trigonometric formulae during 505-587 CE - https://en.wikipedia.org/wiki/Var%C4%81hamihira#Contributions
            "varahamihira",

            # Sir Mokshagundam Visvesvaraya - is a notable Indian engineer.  He is a recipient of the Indian Republic's highest honour, the Bharat Ratna, in 1955. On his birthday, 15 September is celebrated as Engineer's Day in India in his memory - https://en.wikipedia.org/wiki/Visvesvaraya
            "visvesvaraya",

            # Christiane Nüsslein-Volhard - German biologist, won Nobel Prize in Physiology or Medicine in 1995 for research on the genetic control of embryonic development. https://en.wikipedia.org/wiki/Christiane_N%C3%BCsslein-Volhard
            "volhard",

            # Marlyn Wescoff - one of the original programmers of the ENIAC. https://en.wikipedia.org/wiki/ENIAC - https://en.wikipedia.org/wiki/Marlyn_Meltzer
            "wescoff",

            # Andrew Wiles - Notable British mathematician who proved the enigmatic Fermat's Last Theorem - https://en.wikipedia.org/wiki/Andrew_Wiles
            "wiles",

            # Roberta Williams, did pioneering work in graphical adventure games for personal computers, particularly the King's Quest series. https://en.wikipedia.org/wiki/Roberta_Williams
            "williams",

            # Sophie Wilson designed the first Acorn Micro-Computer and the instruction set for ARM processors. https://en.wikipedia.org/wiki/Sophie_Wilson
            "wilson",

            # Jeannette Wing - co-developed the Liskov substitution principle. - https://en.wikipedia.org/wiki/Jeannette_Wing
            "wing",

            # Steve Wozniak invented the Apple I and Apple II. https://en.wikipedia.org/wiki/Steve_Wozniak
            "wozniak",

            # The Wright brothers, Orville and Wilbur - credited with inventing and building the world's first successful airplane and making the first controlled, powered and sustained heavier-than-air human flight - https://en.wikipedia.org/wiki/Wright_brothers
            "wright",

            # Rosalyn Sussman Yalow - Rosalyn Sussman Yalow was an American medical physicist, and a co-winner of the 1977 Nobel Prize in Physiology or Medicine for development of the radioimmunoassay technique. https://en.wikipedia.org/wiki/Rosalyn_Sussman_Yalow
            "yalow",

            # Ada Yonath - an Israeli crystallographer, the first woman from the Middle East to win a Nobel prize in the sciences. https://en.wikipedia.org/wiki/Ada_Yonath
            "yonath",]

    return ("{}_{}".format(choice(left),choice(right)))

# FROM https://github.com/mkaz/termgraph/blob/master/termgraph.py
def chart(labels, data):
    args = {}
    args['format'] = '{:>5.0f}'
    args['suffix'] = ''
    args['width'] = 50

    # verify data
    m = len(labels)
    if m != len(data):
        print(">> Error: Label and data array sizes don't match")
        sys.exit(1)

    # massage data
    # normalize for graph
    max = 0
    for i in range(m):
        if data[i] > max:
            max = data[i]

    step = max / args['width']
    # display graph
    for i in range(m):
        print_blocks(labels[i], data[i], step, args)

    print()


def print_blocks(label, count, step, args):
    # TODO: add flag to hide data labels
    blocks = int(count / step)
    print("{}: ".format(label), end="")
    if count < step:
        sys.stdout.write(sm_tick)
    else:
        for i in range(blocks):
            sys.stdout.write(tick)

    print(args['format'].format(count) + args['suffix'])


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
        cprint("Generating credentials...", "yellow", end='', flush=True)
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
    cprint("\nChecking credentials...", "yellow", end='', flush=True)
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
            [subject, index] = pick(["New"] + subjects, "Enter subject: ")

        if len(subjects) == 0 or subject == "New":
            subject = input("\nDocument? ")

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