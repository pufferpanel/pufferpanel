set VERSION=1.0.0

xcopy /E /I /Y client\dist www
xcopy /E /I /Y assets\email email

heat dir www -out www.wxs -gg -g1 -ke -nologo -cg WwwComponent -dr APPLICATIONROOTDIRECTORY
heat dir email -out email.wxs -gg -g1 -ke -nologo -cg EmailComponent -dr APPLICATIONROOTDIRECTORY
candle *.wxs -arch x64 -dProductVersion="%VERSION%" -ext WixFirewallExtension
del www.wxs email.wxs
light *.wixobj -out pufferpanel.msi -b email -b www -ext WixUtilExtension -ext WixUIExtension -ext WixFirewallExtension -cultures:en-us
del *.wixobj *wixpdb
del /f www
del /f email