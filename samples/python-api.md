Language: Python
Version: 3.12
Project Structure:
└── main.py
└── pyproject.toml
└── requirements.txt
└── setup.py
└── test_main.py


## Source Files

### main.py (11 lines)
```py
from flask import Flask

app = Flask(__name__)

@app.route('/')
def hello():
    return {'message': 'Hello World'}

if __name__ == '__main__':
    app.run(debug=True)

```

### pyproject.toml (15 lines)
```toml
[tool.poetry]
name = "python-api"
version = "0.1.0"
description = ""
authors = ["Your Name <you@example.com>"]
readme = "README.md"

[tool.poetry.dependencies]
python = "^3.12"


[build-system]
requires = ["poetry-core"]
build-backend = "poetry.core.masonry.api"

```

### requirements.txt (4 lines)
```txt
Flask==2.0.1
python-dotenv==0.19.0
pytest==7.0.1

```

### setup.py (9 lines)
```py
from setuptools import setup

setup(
    name='python-api-sample',
    version='1.0.0',
    description='Sample Python API',
    packages=[''],
)

```

### test_main.py (3 lines)
```py
def test_hello():
    assert True

```
