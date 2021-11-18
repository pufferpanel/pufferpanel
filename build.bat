heat dir www -out www.wxs -gg -g1 -ke -nologo -cg WwwComponent -dr APPLICATIONROOTDIRECTORY
heat dir email -out email.wxs -gg -g1 -ke -nologo -cg EmailComponent -dr APPLICATIONROOTDIRECTORY
candle *.wxs -arch x64 -dProductVersion="%VERSION%"
del www.wxs email.wxs
light *.wixobj -out pufferpanel.msi -b email -b www -ext WixUIExtension -cultures:en-us
del *.wixobj *wixpdb
pause