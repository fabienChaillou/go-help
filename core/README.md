## List all bases langage to golang

Based tthis [source](http://le-go-par-l-exemple.keiruaprod.fr/)

tools
Add file into all folders
```shell
ls -d */ | xargs -I {} touch {}/README.md

## list all folder
for i in $(ls -d */); do echo ${i%}; done 
```

search file name
`find ~/ -type f -name "postgis-2.0.0"`
