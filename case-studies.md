- The next request returns headers but not body using ICAP but works with a regular server:

'''
http_proxy=http://192.168.10.168:8080/ curl -H "If-Unmodified-Since: Wed, 02 Mar 2016 23:39:12 GMT" -H "Range: bytes=11060536-11184850" http://fg.v4.download.windowsupdate.com/d/msdownload/update/software/crup/2016/03/windows10.0-kb3141032-x64_a3c042bfeb72c3866f0dfd1d1c6cd518b1c5bfd5.psf
'''

