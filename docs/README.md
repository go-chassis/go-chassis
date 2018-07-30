# How to write doc in local

## Prepare
1. Install Python 2.7 with zlib, libssl-dev(openssl-devel)
1. Install pip
1. Install readthe doc support https://docs.readthedocs.io/en/latest/getting_started.html
1. Install RTD module
```shell
sudo pip install sphinx_rtd_theme
```

## Generate doc

In windows
```shell
cd docs
make.bat html
```

In linux
```shell
cd docs
sphinx-autobuild . _build/html
```

## Check the result

1. See html pages in _build folder
1. Access http://127.0.0.1:8000
