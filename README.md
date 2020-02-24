# EpubFixer
**Warning: This project has been written in one night for one set of e-books which were freezing my e-reader. I don't plan to spend more time on this project unless i found new set of books which won't work for me.**

This tool out of the box works only for specific set of books. It should preserve required information but will lose unsupported optional tags.

EpubFixer fixes epub file with recreating and modifying spike and list of refitems. It needs to be modified for specific purpose but right now it gives a simple base for parsing EPUB files. Right now it's removing one specific file from each spine and remove additional identifiers.

## Usage
```
EpubFixer (file1) (file2) ... (fileN)
```
For each file '(fileN).fixed.epub' will be produced

## (For developers) How it works

For each filename _ProcessFile(filename)_ is executed. It opens zip reader for path then parse zip archive to _EPUB_ structure, patch it then serialize to epub file. There is possibility that some optional tags aren't implemented for xml files.

* parseContainer.go - Container parser implements only minimum requirements.
* parseRootfile.go - This structures probably need to be adjusted for more optional tags. In the root file there are few elements which need to be separated between Marshal and Unmarshal version(namespaces workaround, https://stackoverflow.com/questions/48609596/xml-namespace-prefix-issue-at-go).

So this version can only remove file and modify rootfiles/container. To modify or filter content of htmls, code needs to be upgraded. Probably some wrapper around _zip.File_ structure which will give option to read from raw archive or hold file in buffer which will be modified. This way should save some memory.