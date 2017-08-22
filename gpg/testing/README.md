```
gpg --gen-key
ID: "Testy McTestFace"
PW: 1234
gpg --yes --armor --recipient "Testy McTestFace" --trust-model always --encrypt hello.txt
```
