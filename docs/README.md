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
pip install recommonmark
pip install sphinx-markdown-tables

cd docs
make.bat html
```

In linux
```shell
pip install recommonmark
pip install sphinx-markdown-tables

cd docs

sphinx-build . _build/html
or
make html
```

## Check the result

1. See html pages in _build/html folder
2. Open this file in your web browser to see docs.
