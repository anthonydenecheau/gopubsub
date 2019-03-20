
## ORACLE 
https://github.com/go-goracle/goracle/tree/v2.12.3

https://blogs.oracle.com/developers/how-to-connect-a-go-program-to-oracle-database-using-goracle

go get gopkg.in/goracle.v2

Pré-requis Windows:
* oracle client 
dézipper instantclient-basic-windows.x64-11.2.0.4.0.zip
Ajouter au PATH: C:\developpement\app\oracle\instantclient_11_2

* *mingw-w64 >= gcc 7.2.0 
Installer mingw-w64-install.exe
Créer une variable d'environnement MINGWPATH  C:\Program Files\mingw-w64\x86_64-8.1.0-posix-seh-rt_v6-rev0\mingw64
Ajouter au PATH: %MINGWPATH%\bin

## RUN 
go run main.go --config C:\Users\adenecheau\go\src\github.com\anthonydenecheau\gopubsub\.gopubsub.yaml publisher

go run main.go --config C:\Users\adenecheau\go\src\github.com\anthonydenecheau\gopubsub\.gopubsub.yaml subscriber